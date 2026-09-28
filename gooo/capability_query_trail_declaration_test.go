package gooo

import "testing"

func TestDiscoverCapabilityQueryTrailForDeclarationFocusesOverview(t *testing.T) {
	trail := DiscoverCapabilityQueryTrailForDeclaration(
		"What can gooo do with this declaration?",
		"input invoice\ntransform normalize",
	)
	if err := trail.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if !containsCapabilityQuerySuggestion(trail.Response.SuggestedQueries, "What IR can this .gooo declaration produce?") {
		t.Fatalf("suggestions = %v, missing declaration-aware IR question", trail.Response.SuggestedQueries)
	}
}

func TestDiscoverCapabilityQueryTrailForDeclarationPreservesExplicitQuestion(t *testing.T) {
	trail := DiscoverCapabilityQueryTrailForDeclaration(
		"How do I inspect reverse observation evidence?",
		"input invoice\ntransform normalize",
	)
	if err := trail.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if trail.Intent != "capability_request" {
		t.Fatalf("intent = %q, want capability_request", trail.Intent)
	}
}