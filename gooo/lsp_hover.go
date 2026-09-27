package gooo

import (
	"fmt"
	"strings"
)

type HoverResult struct {
	Status               string
	MissingStage         string
	MissingStageIndex    int
	SourceDigest         string
	IRDigest             string
	EvidencePrefixDigest string
	Requested            Position
	SymbolName           string
	Kind                 SymbolKind
	SymbolDigest         string
	Declaration          string
	Contents             string
	HoverDigest          string
	NonExecuting         bool
	NonAuthorizing       bool
}

func Hover(source string, position Position) HoverResult {
	snapshot := Analyze(source)
	result := HoverResult{
		Status:               "UNKNOWN",
		MissingStage:         "lsp-hover",
		MissingStageIndex:    hoverMissingStageIndex("lsp-hover"),
		SourceDigest:         snapshot.SourceDigest,
		IRDigest:             snapshot.IRDigest,
		EvidencePrefixDigest: hoverEvidencePrefixDigest(source, position),
		Requested:            position,
		NonExecuting:         true,
		NonAuthorizing:       true,
	}
	if snapshot.Status != "BOUND" {
		missingStage := snapshot.MissingStage
		if missingStage == "" {
			missingStage = "lsp-hover-analysis"
		}
		result.setMissingStage(missingStage)
		return result
	}
	if position.Line < 1 || position.Column < 1 {
		result.setMissingStage("lsp-hover-position")
		return result
	}
	declaration, ok := sourceLine(source, position.Line)
	if !ok {
		result.setMissingStage("lsp-hover-position")
		return result
	}
	if position.Column > len([]rune(declaration))+1 {
		result.setMissingStage("lsp-hover-position")
		return result
	}
	for _, symbol := range snapshot.Symbols {
		if symbol.Position.Line != position.Line {
			continue
		}
		result.Status = "BOUND"
		result.MissingStage = ""
		result.MissingStageIndex = -1
		result.SymbolName = symbol.Name
		result.Kind = symbol.Kind
		result.SymbolDigest = symbol.Digest
		result.Declaration = declaration
		result.Contents = fmt.Sprintf("%s %s: %s", symbol.Kind, symbol.Name, declaration)
		result.HoverDigest = digestHoverResult(result)
		return result
	}
	result.setMissingStage("lsp-hover-symbol")
	return result
}

func (h *HoverResult) setMissingStage(stage string) {
	h.MissingStage = stage
	h.MissingStageIndex = hoverMissingStageIndex(stage)
}

func (h HoverResult) Validate() error {
	if h.Status != "BOUND" && h.Status != "UNKNOWN" {
		return fmt.Errorf("hover status %q is invalid", h.Status)
	}
	if !validDigest(h.SourceDigest) {
		return fmt.Errorf("hover source digest is invalid")
	}
	if h.MissingStage == "" && h.MissingStageIndex != -1 {
		return fmt.Errorf("bound hover must use missing stage index -1")
	}
	if h.MissingStage != "" && h.MissingStageIndex < 0 {
		return fmt.Errorf("unknown hover must retain a missing stage index")
	}
	if h.EvidencePrefixDigest != "" && !validDigest(h.EvidencePrefixDigest) {
		return fmt.Errorf("hover evidence prefix digest is invalid")
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

func hoverMissingStageIndex(stage string) int {
	switch stage {
	case "syntax":
		return 0
	case "lsp-hover-analysis":
		return 1
	case "lsp-hover-position":
		return 2
	case "lsp-hover-symbol":
		return 3
	case "lsp-hover":
		return 4
	default:
		return 5
	}
}

func hoverEvidencePrefixDigest(source string, position Position) string {
	if position.Line < 1 {
		return ""
	}
	lines := strings.Split(source, "\n")
	if position.Line > len(lines) {
		return ""
	}
	return digestString(strings.Join(lines[:position.Line], "\n"))
}

func digestHoverResult(result HoverResult) string {
	return digestString(fmt.Sprintf("%s|%s|%d|%s|%s|%s|%d|%d|%s|%s|%s|%s|%s|%t|%t",
		result.Status,
		result.MissingStage,
		result.MissingStageIndex,
		result.SourceDigest,
		result.IRDigest,
		result.EvidencePrefixDigest,
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
