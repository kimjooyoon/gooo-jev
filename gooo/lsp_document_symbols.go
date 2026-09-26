package gooo

import "fmt"

type DocumentSymbolResult struct {
	Status         string
	MissingStage   string
	SourceDigest   string
	IRDigest       string
	Symbols        []Symbol
	SymbolsDigest  string
	NonExecuting   bool
	NonAuthorizing bool
}

// DocumentSymbols projects the analyzed declaration index for editor tooling.
func DocumentSymbols(source string) DocumentSymbolResult {
	snapshot := Analyze(source)
	result := DocumentSymbolResult{
		Status:         "UNKNOWN",
		MissingStage:   "lsp-document-symbols",
		SourceDigest:   snapshot.SourceDigest,
		IRDigest:       snapshot.IRDigest,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if snapshot.Status != "BOUND" {
		result.MissingStage = snapshot.MissingStage
		if result.MissingStage == "" {
			result.MissingStage = "lsp-document-symbols-analysis"
		}
		return result
	}
	result.Status = "BOUND"
	result.MissingStage = ""
	result.Symbols = append([]Symbol(nil), snapshot.Symbols...)
	result.SymbolsDigest = digestDocumentSymbols(result)
	return result
}

func (s DocumentSymbolResult) Validate() error {
	if s.Status != "BOUND" && s.Status != "UNKNOWN" {
		return fmt.Errorf("document symbols status %q is invalid", s.Status)
	}
	if !validDigest(s.SourceDigest) {
		return fmt.Errorf("document symbols source digest is invalid")
	}
	if !s.NonExecuting || !s.NonAuthorizing {
		return fmt.Errorf("document symbols must remain non-executing and non-authorizing")
	}
	if s.Status == "UNKNOWN" {
		if s.MissingStage == "" {
			return fmt.Errorf("unknown document symbols must retain a missing stage")
		}
		return nil
	}
	if s.MissingStage != "" || !validDigest(s.IRDigest) || !validDigest(s.SymbolsDigest) {
		return fmt.Errorf("bound document symbols are incomplete")
	}
	for index, symbol := range s.Symbols {
		if symbol.Name == "" || symbol.Position.Line < 1 || symbol.Position.Column < 1 {
			return fmt.Errorf("document symbol %d is incomplete", index)
		}
		if symbol.Kind != EntitySymbol && symbol.Kind != ActivitySymbol {
			return fmt.Errorf("document symbol %d kind is invalid", index)
		}
		if !validDigest(symbol.Digest) {
			return fmt.Errorf("document symbol %d digest is invalid", index)
		}
	}
	if digestDocumentSymbols(s) != s.SymbolsDigest {
		return fmt.Errorf("document symbols digest does not match its fields")
	}
	return nil
}

func digestDocumentSymbols(result DocumentSymbolResult) string {
	value := fmt.Sprintf("%s|%s", result.SourceDigest, result.IRDigest)
	for _, symbol := range result.Symbols {
		value += fmt.Sprintf("|%s|%s|%d|%d|%s",
			symbol.Name,
			symbol.Kind,
			symbol.Position.Line,
			symbol.Position.Column,
			symbol.Digest,
		)
	}
	return digestString(value)
}
