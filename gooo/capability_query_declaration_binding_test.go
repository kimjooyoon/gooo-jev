package gooo

import "testing"

func TestDiscoverCapabilityQueryWithDeclarationUsesExplicitBindings(t *testing.T) {
	response := DiscoverCapabilityQueryWithDeclaration(
		"What can gooo do with this declaration?",
		"observe Usage\n",
	)
	if response.Status != CapabilityQueryAvailable || len(response.Capabilities) != 3 {
		t.Fatalf("unexpected declaration-bound response: %+v", response)
	}
	if len(response.Declaration.ObservedCapabilities) != 3 || !hasCapabilityQuery(response.Capabilities, "provenance") {
		t.Fatalf("explicit declaration bindings were not preserved: %+v", response.Declaration)
	}
	if err := response.Validate(); err != nil {
		t.Fatalf("declaration-bound response should validate: %v", err)
	}
}

func TestDiscoverCapabilityQueryWithDeclarationPreservesUnknownBindingGap(t *testing.T) {
	response := DiscoverCapabilityQueryWithDeclaration("What can gooo do?", "module Usage\n")
	if response.Status != CapabilityQueryUnknown || response.FirstMismatch != "declaration_capability_binding" || response.MissingStage != "declaration_capability_binding" || len(response.Capabilities) != 0 {
		t.Fatalf("unmapped declaration should remain unknown: %+v", response)
	}
	if err := response.Validate(); err != nil {
		t.Fatalf("unknown declaration-bound response should validate: %v", err)
	}
}

func TestDiscoverCapabilityQueryWithDeclarationRejectsTamperedObservedCapabilities(t *testing.T) {
	response := DiscoverCapabilityQueryWithDeclaration("What can gooo do?", "observe Usage\n")
	response.Declaration.ObservedCapabilities[0] = "tampered"
	if err := response.Validate(); err == nil {
		t.Fatal("tampered observed capability binding should fail validation")
	}
}