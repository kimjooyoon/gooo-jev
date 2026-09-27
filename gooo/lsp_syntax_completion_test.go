package gooo

import "testing"

func TestCompleteSyntaxSuggestsKeywordsForPartialSource(t *testing.T) {
	response := CompleteSyntax("package support_triage\n", "n")
	if len(response.Items) != 1 || response.Items[0].Label != "namespace" {
		t.Fatalf("items = %#v, want namespace", response.Items)
	}
	if response.Status != "UNKNOWN" || len(response.Diagnostics) == 0 {
		t.Fatalf("response = %#v, want unknown with diagnostics", response)
	}
	if err := response.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestCompleteSyntaxSuggestsTopLevelKeywordsForEmptySource(t *testing.T) {
	response := CompleteSyntax("", "")
	if len(response.Items) != 4 {
		t.Fatalf("items = %#v, want four top-level keywords", response.Items)
	}
	if err := response.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}
