package gooo

import (
	"fmt"
	"strings"
)

// CompletionItem is a provenance-linked language completion candidate.
type CompletionItem struct {
	Label        string
	Kind         SymbolKind
	Detail       string
	Position     Position
	SymbolDigest string
	Digest       string
	ItemDigest   string
}

// CompletionResponse is a non-executing completion projection.
type CompletionResponse struct {
	Status         string
	MissingStage   string
	SourceDigest   string
	IRDigest       string
	Prefix         string
	Items          []CompletionItem
	ItemsDigest    string
	Diagnostics    []Diagnostic
	NonExecuting   bool
	NonAuthorizing bool
}

// Complete returns deterministic declaration completions for a source prefix.
func Complete(source, prefix string) CompletionResponse {
	snapshot := Analyze(source)
	response := CompletionResponse{
		Status:         snapshot.Status,
		MissingStage:   snapshot.MissingStage,
		SourceDigest:   snapshot.SourceDigest,
		IRDigest:       snapshot.IRDigest,
		Prefix:         prefix,
		Diagnostics:    append([]Diagnostic(nil), snapshot.Diagnostics...),
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if snapshot.Status != "BOUND" {
		return response
	}
	for _, symbol := range snapshot.Symbols {
		if strings.HasPrefix(symbol.Name, prefix) {
			item := CompletionItem{
				Label:        symbol.Name,
				Kind:         symbol.Kind,
				Detail:       string(symbol.Kind),
				Position:     symbol.Position,
				SymbolDigest: symbol.Digest,
				Digest:       symbol.Digest,
			}
			item.ItemDigest = digestCompletionItem(item)
			response.Items = append(response.Items, item)
		}
	}
	response.ItemsDigest = digestCompletionItems(response)
	return response
}

// Validate checks completion candidates and diagnostics remain source-bound.
func (r CompletionResponse) Validate() error {
	if r.Status != "BOUND" && r.Status != "UNKNOWN" {
		return fmt.Errorf("completion status %q is invalid", r.Status)
	}
	if !validDigest(r.SourceDigest) {
		return fmt.Errorf("completion source digest is invalid")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("completion must remain non-executing and non-authorizing")
	}
	if r.Status == "UNKNOWN" {
		if len(r.Diagnostics) == 0 || r.MissingStage == "" {
			return fmt.Errorf("unknown completion must retain a missing stage and diagnostics")
		}
		return nil
	}
	if r.MissingStage != "" || !validDigest(r.IRDigest) || !validDigest(r.ItemsDigest) {
		return fmt.Errorf("bound completion provenance is incomplete")
	}
	seen := map[string]struct{}{}
	for index, item := range r.Items {
		if item.Label == "" || item.Detail == "" || !validDigest(item.Digest) || !validDigest(item.SymbolDigest) ||
			!validDigest(item.ItemDigest) {
			return fmt.Errorf("completion item %d is incomplete", index)
		}
		if item.SymbolDigest != item.Digest {
			return fmt.Errorf("completion item %q symbol digest is not retained", item.Label)
		}
		if item.Position.Line < 1 || item.Position.Column < 1 {
			return fmt.Errorf("completion item %q position is invalid", item.Label)
		}
		if item.Kind != EntitySymbol && item.Kind != ActivitySymbol {
			return fmt.Errorf("completion item kind is invalid")
		}
		if _, exists := seen[item.Label]; exists {
			return fmt.Errorf("completion item %q is duplicated", item.Label)
		}
		if digestCompletionItem(item) != item.ItemDigest {
			return fmt.Errorf("completion item %q digest does not match its fields", item.Label)
		}
		seen[item.Label] = struct{}{}
	}
	for _, diagnostic := range r.Diagnostics {
		if diagnostic.SourceDigest != r.SourceDigest || diagnostic.Stage == "" || diagnostic.Message == "" {
			return fmt.Errorf("completion diagnostic is not source-bound")
		}
	}
	if len(r.Diagnostics) != 0 {
		return fmt.Errorf("bound completion cannot have diagnostics")
	}
	if digestCompletionItems(r) != r.ItemsDigest {
		return fmt.Errorf("completion items digest does not match its fields")
	}
	return nil
}

func digestCompletionItem(item CompletionItem) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%d|%d|%s",
		item.Label,
		item.Kind,
		item.Detail,
		item.Position.Line,
		item.Position.Column,
		item.SymbolDigest,
	))
}

func digestCompletionItems(response CompletionResponse) string {
	value := fmt.Sprintf("%s|%s|%s", response.SourceDigest, response.IRDigest, response.Prefix)
	for _, item := range response.Items {
		value += "|" + item.ItemDigest
	}
	return digestString(value)
}
