package gooo

import "fmt"

type DefinitionResult struct {
	Status             string
	MissingStage       string
	SourceDigest       string
	IRDigest           string
	Requested          Position
	SymbolName         string
	Kind               SymbolKind
	SymbolDigest       string
	DefinitionPosition Position
	DefinitionDigest   string
	NonExecuting       bool
	NonAuthorizing     bool
}

// Definition resolves a source position to the declaration already indexed by Analyze.
func Definition(source string, position Position) DefinitionResult {
	snapshot := Analyze(source)
	result := DefinitionResult{
		Status:         "UNKNOWN",
		MissingStage:   "lsp-definition",
		SourceDigest:   snapshot.SourceDigest,
		IRDigest:       snapshot.IRDigest,
		Requested:      position,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if snapshot.Status != "BOUND" {
		result.MissingStage = snapshot.MissingStage
		if result.MissingStage == "" {
			result.MissingStage = "lsp-definition-analysis"
		}
		return result
	}
	if position.Line < 1 || position.Column < 1 {
		result.MissingStage = "lsp-definition-position"
		return result
	}
	for _, symbol := range snapshot.Symbols {
		if symbol.Position.Line != position.Line ||
			position.Column < symbol.Position.Column ||
			position.Column > symbol.Position.Column+len(symbol.Name) {
			continue
		}
		result.Status = "BOUND"
		result.MissingStage = ""
		result.SymbolName = symbol.Name
		result.Kind = symbol.Kind
		result.SymbolDigest = symbol.Digest
		result.DefinitionPosition = symbol.Position
		result.DefinitionDigest = digestDefinitionResult(result)
		return result
	}
	result.MissingStage = "lsp-definition-symbol"
	return result
}

func (d DefinitionResult) Validate() error {
	if d.Status != "BOUND" && d.Status != "UNKNOWN" {
		return fmt.Errorf("definition status %q is invalid", d.Status)
	}
	if !validDigest(d.SourceDigest) {
		return fmt.Errorf("definition source digest is invalid")
	}
	if !d.NonExecuting || !d.NonAuthorizing {
		return fmt.Errorf("definition must remain non-executing and non-authorizing")
	}
	if d.Requested.Line < 1 || d.Requested.Column < 1 {
		return fmt.Errorf("definition requested position must be positive")
	}
	if d.Status == "UNKNOWN" {
		if d.MissingStage == "" {
			return fmt.Errorf("unknown definition must retain a missing stage")
		}
		return nil
	}
	if d.MissingStage != "" || d.SymbolName == "" || d.Kind == "" ||
		!validDigest(d.IRDigest) || !validDigest(d.SymbolDigest) ||
		d.DefinitionPosition.Line < 1 || d.DefinitionPosition.Column < 1 ||
		!validDigest(d.DefinitionDigest) {
		return fmt.Errorf("bound definition is incomplete")
	}
	if d.Kind != EntitySymbol && d.Kind != ActivitySymbol {
		return fmt.Errorf("definition kind is invalid")
	}
	if digestDefinitionResult(d) != d.DefinitionDigest {
		return fmt.Errorf("definition digest does not match its fields")
	}
	return nil
}

func digestDefinitionResult(result DefinitionResult) string {
	return digestString(fmt.Sprintf("%s|%s|%d|%d|%s|%s|%s|%d|%d|%t|%t",
		result.SourceDigest,
		result.IRDigest,
		result.Requested.Line,
		result.Requested.Column,
		result.SymbolName,
		result.Kind,
		result.SymbolDigest,
		result.DefinitionPosition.Line,
		result.DefinitionPosition.Column,
		result.NonExecuting,
		result.NonAuthorizing,
	))
}
