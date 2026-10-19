package gooo

import "testing"

func TestDiscoverCapabilityQueryOverview(t *testing.T) {
	response := DiscoverCapabilityQuery("What can gooo do?")
	if response.Status != CapabilityQueryAvailable || len(response.Capabilities) == 0 {
		t.Fatalf("unexpected overview response: %+v", response)
	}
	if !hasCapabilityQuery(response.Capabilities, "canonical_generation") || !hasCapabilityQuery(response.Capabilities, "provenance") {
		t.Fatalf("overview omitted useful capabilities: %+v", response.Capabilities)
	}
	canonical, ok := capabilityQueryByID(response.Capabilities, "canonical_generation")
	if !ok || canonical.ExampleQuery == "" {
		t.Fatalf("overview omitted an actionable example query: %+v", response.Capabilities)
	}
	if err := response.Validate(); err != nil {
		t.Fatalf("overview should validate: %v", err)
	}
}

func TestDiscoverCapabilityQueryPreservesExternalBoundary(t *testing.T) {
	response := DiscoverCapabilityQuery("Can gooo execute and authorize a network workload?")
	if response.Status != CapabilityQueryDeferred || response.FirstMismatch != "external_boundary" {
		t.Fatalf("unexpected deferred response: %+v", response)
	}
	if err := response.Validate(); err != nil {
		t.Fatalf("deferred response should validate: %v", err)
	}
}

func TestDiscoverCapabilityQueryUnknownRetainsSuggestions(t *testing.T) {
	response := DiscoverCapabilityQuery("quantum breakfast compiler")
	if response.Status != CapabilityQueryUnknown || response.FirstMismatch != "query" || response.MissingStage != "capability_catalog" || len(response.Suggestions) == 0 || len(response.SuggestedQueries) == 0 || !hasCapabilityQuerySuggestion(response.SuggestedQueries, "How do I generate a canonical .gooo declaration?") {
		t.Fatalf("unexpected unknown response: %+v", response)
	}
	if err := response.Validate(); err != nil {
		t.Fatalf("unknown response should validate: %v", err)
	}
}

func TestDiscoverCapabilityQueryRejectsTamperedDigest(t *testing.T) {
	response := DiscoverCapabilityQuery("show provenance and generation")
	response.Capabilities[0].ExampleQuery = "tampered"
	if err := response.Validate(); err == nil {
		t.Fatal("tampered response should fail validation")
	}
}

func TestDiscoverCapabilityQueryBindsDeclarationObservation(t *testing.T) {
	response := DiscoverCapabilityQueryWithDeclaration(
		"What can gooo do with this declaration?",
		"entity Usage\noperation observe\n",
	)
	if response.Declaration == nil || !response.Declaration.Bound || response.Declaration.SourceDigest == "" {
		t.Fatalf("declaration observation was not bound: %+v", response.Declaration)
	}
	if len(response.Declaration.ObservedSignals) != 2 || response.Declaration.ObservedSignals[0] != "entity" || response.Declaration.ObservedSignals[1] != "operation" {
		t.Fatalf("unexpected declaration signals: %+v", response.Declaration.ObservedSignals)
	}
	if err := response.Validate(); err != nil {
		t.Fatalf("declaration-bound response should validate: %v", err)
	}
}

func TestDiscoverCapabilityQueryRejectsTamperedDeclarationBinding(t *testing.T) {
	response := DiscoverCapabilityQueryWithDeclaration("show provenance", "entity Usage")
	response.Declaration.SourceDigest = "sha256:tampered"
	if err := response.Validate(); err == nil {
		t.Fatal("tampered declaration binding should fail validation")
	}
}

func TestDiscoverCapabilityQueryPrefersScopedCapabilityOverOverview(t *testing.T) {
	response := DiscoverCapabilityQuery("What can gooo do with provenance?")
	if response.Status != CapabilityQueryAvailable || len(response.Capabilities) != 1 || response.Capabilities[0].ID != "provenance" {
		t.Fatalf("specific capability should win over overview: %+v", response)
	}
	if err := response.Validate(); err != nil {
		t.Fatalf("specific capability response should validate: %v", err)
	}
}

func TestDiscoverCapabilityQueryNormalizesUnicodePunctuation(t *testing.T) {
	response := DiscoverCapabilityQuery("  CAN—GOOO? INSPECT… PROVENANCE!  ")
	if response.Status != CapabilityQueryAvailable || len(response.Capabilities) != 1 || response.Capabilities[0].ID != "provenance" {
		t.Fatalf("unicode punctuation should not change scoped matching: %+v", response)
	}
	if response.Query != "CAN—GOOO? INSPECT… PROVENANCE!" {
		t.Fatalf("original query should remain evidence-visible: %q", response.Query)
	}
}

func TestDiscoverCapabilityQueryRejectsSubstringAlias(t *testing.T) {
	response := DiscoverCapabilityQuery("provenanceful")
	if response.Status != CapabilityQueryUnknown || len(response.Capabilities) != 0 {
		t.Fatalf("substring alias should not match: %+v", response)
	}
}

func TestDiscoverCapabilityQueryRoutesProvenanceAndReverseObservation(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		status CapabilityQueryState
		ids    []string
	}{
		{name: "bare provenance", query: "provenance", status: CapabilityQueryAvailable, ids: []string{"provenance"}},
		{name: "provenance chain", query: "inspect provenance chain", status: CapabilityQueryAvailable, ids: []string{"provenance"}},
		{name: "explicit reverse", query: "How do I inspect reverse observation evidence?", status: CapabilityQueryAvailable, ids: []string{"reverse_observation"}},
		{name: "origin alias", query: "origin", status: CapabilityQueryAvailable, ids: []string{"reverse_observation"}},
		{name: "mixed observations", query: "show provenance and reverse observation", status: CapabilityQueryAvailable, ids: []string{"provenance", "reverse_observation"}},
		{name: "external boundary", query: "show provenance and authorize a network workload", status: CapabilityQueryDeferred, ids: []string{"execution_authorization", "provenance", "security_boundary"}},
		{name: "substring", query: "provenanceful", status: CapabilityQueryUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := DiscoverCapabilityQuery(tt.query)
			if response.Status != tt.status || len(response.Capabilities) != len(tt.ids) {
				t.Fatalf("unexpected capability routing: %+v", response)
			}
			for i, id := range tt.ids {
				if response.Capabilities[i].ID != id {
					t.Fatalf("capability[%d] = %q, want %q", i, response.Capabilities[i].ID, id)
				}
			}
			if err := response.Validate(); err != nil {
				t.Fatalf("routed response should validate: %v", err)
			}
		})
	}
}

func hasCapabilityQuery(capabilities []CapabilityQueryCapability, id string) bool {
	_, ok := capabilityQueryByID(capabilities, id)
	return ok
}

func hasCapabilityQuerySuggestion(suggestions []string, expected string) bool {
	for _, suggestion := range suggestions {
		if suggestion == expected {
			return true
		}
	}
	return false
}

func capabilityQueryByID(capabilities []CapabilityQueryCapability, id string) (CapabilityQueryCapability, bool) {
	for _, capability := range capabilities {
		if capability.ID == id {
			return capability, true
		}
	}
	return CapabilityQueryCapability{}, false
}

func TestCapabilityQuerySuggestedExamplesAreDiscoverable(t *testing.T) {
	for _, entry := range capabilityQueryCatalog {
		if !entry.Safe {
			continue
		}
		response := DiscoverCapabilityQuery(entry.ExampleQuery)
		if response.Status != CapabilityQueryAvailable || len(response.Capabilities) == 0 {
			t.Fatalf("catalog example is not discoverable: id=%q query=%q response=%+v", entry.ID, entry.ExampleQuery, response)
		}
		if err := response.Validate(); err != nil {
			t.Fatalf("catalog example should validate: id=%q error=%v", entry.ID, err)
		}
	}
}
