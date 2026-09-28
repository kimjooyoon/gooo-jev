package gooo

import "testing"

func TestDiscoverCapabilityQueryWithContractSketchDeclaration(t *testing.T) {
	source := "module usage_replan\ninput usage_observation\nstatus BOUND\nfield status\nconstraint no_execution"
	response := DiscoverCapabilityQueryWithDeclaration("Can gooo analyze this .gooo declaration?", source)
	if err := response.Validate(); err != nil {
		t.Fatalf("validate contract sketch capability response: %v", err)
	}
	if response.Status != CapabilityQueryAvailable {
		t.Fatalf("status = %q, want %q", response.Status, CapabilityQueryAvailable)
	}
	if response.Declaration == nil || len(response.Declaration.ObservedSignals) != 5 {
		t.Fatalf("observed signals = %#v, want five contract sketch signals", response.Declaration)
	}
	if len(response.Declaration.ObservedCapabilities) != 1 || response.Declaration.ObservedCapabilities[0] != "declaration_analysis" {
		t.Fatalf("observed capabilities = %#v, want declaration_analysis", response.Declaration.ObservedCapabilities)
	}
}
