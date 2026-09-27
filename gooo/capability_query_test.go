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
	if response.Status != CapabilityQueryUnknown || response.FirstMismatch != "query" || response.MissingStage != "capability_catalog" || len(response.Suggestions) == 0 {
		t.Fatalf("unexpected unknown response: %+v", response)
	}
	if err := response.Validate(); err != nil {
		t.Fatalf("unknown response should validate: %v", err)
	}
}

func TestDiscoverCapabilityQueryRejectsTamperedDigest(t *testing.T) {
	response := DiscoverCapabilityQuery("show provenance and generation")
	response.Capabilities[0].Description = "tampered"
	if err := response.Validate(); err == nil {
		t.Fatal("tampered response should fail validation")
	}
}

func hasCapabilityQuery(capabilities []CapabilityQueryCapability, id string) bool {
	for _, capability := range capabilities {
		if capability.ID == id {
			return true
		}
	}
	return false
}
