package gooo

import (
	"fmt"
	"strings"
)

type HoverResult struct {
	Status           string
	MissingStage     string
	SourceDigest     string
	IRDigest         string
	Requested        Position
	SymbolName       string
	Kind             SymbolKind
	SymbolDigest     string
	Declaration      string
	Contents         string
	HoverDigest      string
	NonExecuting     bool
	NonAuthorizing   bool
}

func Hover(source string, position Position) HoverResult {
	snapshot := Analyze(source)
	result := HoverResult{
		Status:         "UNKNOWN",
		MissingStage:   "lsp-hover",
		SourceDigest:   snapshot.SourceDigest,
		IRDigest:       snapshot.IRDigest,
		Requested:       position,
		NonExecuting:    true,
		NonAuthorizing:  true,
	}
	if snapshot.Status != "BOUND" {
		result.MissingStage = snapshot.MissingStage
		if result.MissingStage == "" {
			result.MissingStage = "lsp-hover-analysis"
		}
		return result
	}
	if position.Line < 1 || position.Column < 1 {
		result.MissingStage = "lsp-hover-position"
		return result
	}
	declaration, ok := sourceLine(source, position.Line)
	if !ok {
		result.MissingStage = "lsp-hover-position"
		return result
	}
	for _, symbol := range snapshot.Symbols {
		if symbol.Position.Line != position.Line {
			continue
		}
		result.Status = "BOUND"
		result.MissingStage = ""
		result.SymbolName = symbol.Name
		result.Kind = symbol.Kind
		result.SymbolDigest = symbol.Digest
		result.Declaration = declaration
		result.Contents = fmt.Sprintf("%s %s: %s", symbol.Kind, symbol.Name, declaration)
		result.HoverDigest = digestHoverResult(result)
		return result
	}
	result.MissingStage = "lsp-hover-symbol"
	return result
}

func (h HoverResult) Validate() error {
	if h.Status != "BOUND" && h.Status != "UNKNOWN" {
		return fmt.Errorf("hover status %q is invalid", h.Status)
	}
	if !validDigest(h.SourceDigest) {
		return fmt.Errorf("hover source digest is invalid")
	}
	if !h.NonExecuting || !h.NonAuthorizing {
		return fmt.Errorf("hover must remain non-executing and non-authorizing")
	}
	if h.Requested.Line < 1 || h.Requested.Column < 1 {
		return fmt.Errorf("hover requested position must be positive")
	}
	if h.Status == "UNKNOWN" {
		if h.MissingStage == "" {
			return fmt.Errorf("unknown hover must retain a missing stage")
		}
		return nil
	}
	if h.MissingStage != "" || h.SymbolName == "" || h.Kind == "" {
		return fmt.Errorf("bound hover is incomplete")
	}
	if !validDigest(h.SymbolDigest) || h.Declaration == "" || h.Contents == "" {
		return fmt.Errorf("bound hover evidence is incomplete")
	}
	if !validDigest(h.HoverDigest) || digestHoverResult(h) != h.HoverDigest {
		return fmt.Errorf("hover digest does not match its fields")
	}
	return nil
}

func sourceLine(source string, line int) (string, bool) {
	lines := strings.Split(source, "\n")
	if line < 1 || line > len(lines) {
		return "", false
	}
	return lines[line-1], true
}

func digestHoverResult(result HoverResult) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%d|%d|%s|%s|%s|%s|%s|%t|%t",
		result.Status,
		result.MissingStage,
		result.SourceDigest,
		result.IRDigest,
		result.Requested.Line,
		result.Requested.Column,
		result.SymbolName,
		result.Kind,
		result.SymbolDigest,
		result.Declaration,
		result.Contents,
		result.NonExecuting,
		result.NonAuthorizing,
	))
}
