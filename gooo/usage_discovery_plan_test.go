package gooo

import "testing"

func TestDiscoverUsageAndPlanPreserveCapabilityState(t *testing.T) {
	source := "package support\nnamespace triage\nentity ticket\nproperty\n"

	discovery, err := DiscoverUsage(source, "pro")
	if err != nil {
		t.Fatalf("discover usage: %v", err)
	}
	if err := discovery.Validate(); err != nil {
		t.Fatalf("validate discovery: %v", err)
	}

	plan, err := PlanUsageActions(discovery)
	if err != nil {
		t.Fatalf("plan usage actions: %v", err)
	}
	if err := plan.Validate(); err != nil {
		t.Fatalf("validate usage plan: %v", err)
	}

	var readySyntax, deferredGeneration bool
	for _, action := range plan.Actions {
		if action.CapabilityID == "syntax_completion" && action.State == "READY" {
			readySyntax = true
		}
		if action.CapabilityID == "canonical_generation" && action.State == "DEFERRED" {
			deferredGeneration = true
		}
	}
	if !readySyntax {
		t.Fatal("expected syntax completion to be ready")
	}
	if !deferredGeneration {
		t.Fatal("expected canonical generation to remain deferred")
	}
	if !plan.NonExecuting || !plan.NonAuthorizing {
		t.Fatal("usage plan must remain non-executing and non-authorizing")
	}
}

