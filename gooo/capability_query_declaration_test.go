package gooo

import "testing"

func TestDiscoverCapabilityQueryForDeclarationFocusesSuggestions(t *testing.T) {
	response := DiscoverCapabilityQueryForDeclaration("module billing\ninput invoice\ntransform normalize")
	if err := response.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if len(response.SuggestedQueries) == 0 {
		t.Fatal("expected declaration-focused suggestions")
	}
	if !containsCapabilityQuerySuggestion(response.SuggestedQueries, "What IR can this .gooo declaration produce?") {
		t.Fatalf("suggestions = %v, missing IR guidance", response.SuggestedQueries)
	}
	if !containsCapabilityQuerySuggestion(response.SuggestedQueries, "How do I generate a canonical .gooo declaration?") {
		t.Fatalf("suggestions = %v, missing generation guidance", response.SuggestedQueries)
	}
}

func TestDiscoverCapabilityQueryForDeclarationKeepsGenericGuidanceWithoutDeclaration(t *testing.T) {
	response := DiscoverCapabilityQueryForDeclaration("")
	if err := response.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if len(response.SuggestedQueries) == 0 {
		t.Fatal("expected generic suggestions")
	}
}

func containsCapabilityQuerySuggestion(suggestions []string, want string) bool {
	for _, suggestion := range suggestions {
		if suggestion == want {
			return true
		}
	}
	return false
}