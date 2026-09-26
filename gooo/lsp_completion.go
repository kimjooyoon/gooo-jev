package gooo

import (
	"fmt"
	"strings"
)

// CompletionItem is a provenance-linked language completion candidate.
type CompletionItem struct {
	Label  string
	Kind   SymbolKind
	Detail string
	Digest string
}

// CompletionResponse is a non-executing completion projection.
type CompletionResponse struct {
	Status         string
	MissingStage   string
	SourceDigest   string
	Items          []CompletionItem
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
		Diagnostics:    append([]Diagnostic(nil), snapshot.Diagnostics...),
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if snapshot.Status != "BOUND" {
		return response
	}
	for _, symbol := range snapshot.Symbols {
		if strings.HasPrefix(symbol.Name, prefix) {
			response.Items = append(response.Items, CompletionItem{
				Label: symbol.Name,
				Kind: symbol.Kind,
				Detail: string(symbol.Kind),
				Digest: symbol.Digest,
			})
		}
	}
	return response
}

// Validate checks completion candidates and diagnostics remain source-bound.
func (r CompletionResponse) Validate() error {
	if r.Status != "BOUND" && r.Status != "UNKNOWN" {
		return fmt.Errorf("completion status %q is invalid", r.Status)
	}
	if r.SourceDigest == "" {
		return fmt.Errorf("completion source digest is required")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("completion must remain non-executing and non-authorizing")
	}
	seen := map[string]struct{}{}
	for _, item := range r.Items {
		if item.Label == "" || item.Detail == "" || item.Digest == "" {
			return fmt.Errorf("completion item is incomplete")
		}
		if item.Kind != EntitySymbol && item.Kind != ActivitySymbol {
			return fmt.Errorf("completion item kind is invalid")
		}
		if _, exists := seen[item.Label]; exists {
			return fmt.Errorf("completion item %q is duplicated", item.Label)
		}
		if !validDigest(item.Digest) {
			return fmt.Errorf("completion item %q digest is invalid", item.Label)
		}
		seen[item.Label] = struct{}{}
	}
	for _, diagnostic := range r.Diagnostics {
		if diagnostic.SourceDigest != r.SourceDigest || diagnostic.Stage == "" || diagnostic.Message == "" {
			return fmt.Errorf("completion diagnostic is not source-bound")
		}
	}
	if r.Status == "BOUND" && len(r.Diagnostics) != 0 {
		return fmt.Errorf("bound completion cannot have diagnostics")
	}
	if r.Status == "UNKNOWN" && len(r.Diagnostics) == 0 {
		return fmt.Errorf("unknown completion must retain diagnostics")
	}
	return nil
}
