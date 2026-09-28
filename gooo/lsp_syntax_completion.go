package gooo

import (
	"fmt"
	"sort"
	"strings"
)

type SyntaxCompletionItem struct {
	Label      string
	Kind       string
	Detail     string
	ItemDigest string
}

type SyntaxCompletionResponse struct {
	Status         string
	MissingStage   string
	SourceDigest   string
	IRDigest       string
	Prefix         string
	Items          []SyntaxCompletionItem
	ItemsDigest    string
	Diagnostics    []Diagnostic
	NonExecuting   bool
	NonAuthorizing bool
}

// CompleteSyntax provides useful keyword suggestions while a declaration is
// still incomplete. It never executes or authorizes a declaration.
func CompleteSyntax(source, prefix string) SyntaxCompletionResponse {
	snapshot := Analyze(source)
	response := SyntaxCompletionResponse{
		Status:         snapshot.Status,
		MissingStage:   snapshot.MissingStage,
		SourceDigest:   snapshot.SourceDigest,
		IRDigest:       snapshot.IRDigest,
		Prefix:         prefix,
		Diagnostics:    append([]Diagnostic(nil), snapshot.Diagnostics...),
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	labels := syntaxCompletionLabels(source)
	for _, label := range labels {
		if !strings.HasPrefix(label, prefix) {
			continue
		}
		item := SyntaxCompletionItem{
			Label:  label,
			Kind:   "keyword",
			Detail: "gooo declaration keyword",
		}
		item.ItemDigest = digestString(fmt.Sprintf("%s|%s|%s", item.Label, item.Kind, item.Detail))
		response.Items = append(response.Items, item)
	}
	if snapshot.Status == "BOUND" {
		symbols := Complete(source, prefix)
		for _, item := range symbols.Items {
			response.Items = append(response.Items, SyntaxCompletionItem{
				Label: item.Label, Kind: string(item.Kind), Detail: item.Detail,
				ItemDigest: item.ItemDigest,
			})
		}
	}
	sort.Slice(response.Items, func(left, right int) bool {
		return response.Items[left].Label < response.Items[right].Label
	})
	response.ItemsDigest = digestSyntaxCompletionItems(response)
	return response
}

func (response SyntaxCompletionResponse) Validate() error {
	if response.Status != "BOUND" && response.Status != "UNKNOWN" {
		return fmt.Errorf("syntax completion status %q is invalid", response.Status)
	}
	if !validDigest(response.SourceDigest) {
		return fmt.Errorf("syntax completion source digest is invalid")
	}
	if !response.NonExecuting || !response.NonAuthorizing {
		return fmt.Errorf("syntax completion crossed a capability boundary")
	}
	if response.Status == "UNKNOWN" && len(response.Diagnostics) == 0 {
		return fmt.Errorf("unknown syntax completion lost diagnostics")
	}
	seen := map[string]struct{}{}
	for _, item := range response.Items {
		if item.Label == "" || item.Detail == "" || !validDigest(item.ItemDigest) {
			return fmt.Errorf("syntax completion item is incomplete")
		}
		if item.Kind != "keyword" && item.Kind != string(EntitySymbol) && item.Kind != string(ActivitySymbol) {
			return fmt.Errorf("syntax completion item kind is invalid")
		}
		if _, exists := seen[item.Label]; exists {
			return fmt.Errorf("syntax completion item %q is duplicated", item.Label)
		}
		seen[item.Label] = struct{}{}
		expected := digestString(fmt.Sprintf("%s|%s|%s", item.Label, item.Kind, item.Detail))
		if item.Kind != "keyword" {
			expected = item.ItemDigest
		}
		if expected != item.ItemDigest {
			return fmt.Errorf("syntax completion item %q digest is invalid", item.Label)
		}
	}
	for _, diagnostic := range response.Diagnostics {
		if diagnostic.SourceDigest != response.SourceDigest || diagnostic.Stage == "" || diagnostic.Message == "" {
			return fmt.Errorf("syntax completion diagnostic is not source-bound")
		}
	}
	if digestSyntaxCompletionItems(response) != response.ItemsDigest {
		return fmt.Errorf("syntax completion items digest does not match")
	}
	return nil
}

func syntaxCompletionLabels(source string) []string {
	lines := strings.Split(strings.TrimSpace(source), "\n")
	last := ""
	if len(lines) > 0 {
		last = strings.TrimSpace(lines[len(lines)-1])
	}
	fields := strings.Fields(last)
	if len(fields) == 0 {
		return []string{"package", "namespace", "entity", "activity"}
	}
	switch fields[0] {
	case "package":
		return []string{"namespace"}
	case "namespace":
		return []string{"entity", "activity"}
	case "entity":
		return []string{"property", "activity"}
	case "property":
		return []string{"property", "activity"}
	case "activity":
		return []string{"activity", "entity"}
	default:
		return []string{"package", "namespace", "entity", "activity"}
	}
}

func digestSyntaxCompletionItems(response SyntaxCompletionResponse) string {
	value := fmt.Sprintf("%s|%s|%s", response.SourceDigest, response.IRDigest, response.Prefix)
	for _, item := range response.Items {
		value += "|" + item.ItemDigest
	}
	return digestString(value)
}
