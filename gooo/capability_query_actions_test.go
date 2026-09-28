package gooo

import (
	"strings"
	"testing"
)

func TestSuggestCapabilityQueryActionsForAvailableDeclaration(t *testing.T) {
	hover := DiscoverCapabilityQueryHover("What can gooo do?", "entity user")
	actions, err := SuggestCapabilityQueryActions(hover)
	if err != nil {
		t.Fatalf("suggest actions: %v", err)
	}
	if err := ValidateCapabilityQueryActions(hover, actions); err != nil {
		t.Fatalf("validate actions: %v", err)
	}
	if len(actions) != 1 || actions[0].CapabilityID != "declaration_analysis" {
		t.Fatalf("actions = %#v, want declaration_analysis", actions)
	}
	if actions[0].NextOperation != "inspect_diagnostics" {
		t.Fatalf("next operation = %q, want inspect_diagnostics", actions[0].NextOperation)
	}
}

func TestSuggestCapabilityQueryActionsPreservesExternalBoundary(t *testing.T) {
	hover := DiscoverCapabilityQueryHover("What security boundary does gooo need?", "policy network")
	actions, err := SuggestCapabilityQueryActions(hover)
	if err != nil {
		t.Fatalf("suggest boundary action: %v", err)
	}
	if err := ValidateCapabilityQueryActions(hover, actions); err != nil {
		t.Fatalf("validate boundary actions: %v", err)
	}
	if len(actions) != 1 || actions[0].Kind != CapabilityQueryActionRequireExternalBoundary {
		t.Fatalf("actions = %#v, want one external-boundary action", actions)
	}
	if actions[0].NextOperation != "provide_explicit_external_boundary" {
		t.Fatalf("next operation = %q, want explicit external boundary", actions[0].NextOperation)
	}
}

func TestSuggestCapabilityQueryActionsPreservesUnknown(t *testing.T) {
	hover := DiscoverCapabilityQueryHover("What is the weather?", "entity user")
	actions, err := SuggestCapabilityQueryActions(hover)
	if err != nil {
		t.Fatalf("suggest clarification action: %v", err)
	}
	if err := ValidateCapabilityQueryActions(hover, actions); err != nil {
		t.Fatalf("validate clarification action: %v", err)
	}
	if len(actions) != 1 || actions[0].Kind != CapabilityQueryActionAskClarifyingQuestion || actions[0].CapabilityID != "" {
		t.Fatalf("actions = %#v, want one clarification action", actions)
	}
}

func TestCapabilityQueryActionRejectsTamperedEvidence(t *testing.T) {
	hover := DiscoverCapabilityQueryHover("Can gooo inspect provenance?", "observe user")
	actions, err := SuggestCapabilityQueryActions(hover)
	if err != nil {
		t.Fatalf("suggest actions: %v", err)
	}
	actions[0].EvidenceDigest = strings.Repeat("0", 64)
	if err := ValidateCapabilityQueryActions(hover, actions); err == nil {
		t.Fatal("tampered action evidence should be rejected")
	}
}
