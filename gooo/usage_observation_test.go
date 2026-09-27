package gooo

import (
	"testing"
	"time"
)

func TestObserveUsageActionKeepsDeferredWorkUnknown(t *testing.T) {
	source := "package support\nnamespace triage\nentity ticket\nproperty\n"
	discovery := DiscoverUsage(source, "pro")
	plan, err := PlanUsageActions(discovery)
	if err != nil {
		t.Fatalf("plan usage actions: %v", err)
	}

	ready, err := ObserveUsageAction(
		plan,
		"syntax_completion:edit_declaration",
		UsageObservationConfirmed,
		"completion_quality",
		1,
		plan.PlanDigest,
		time.Date(2026, 9, 28, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("observe ready action: %v", err)
	}
	if err := ready.ValidateAgainst(plan); err != nil {
		t.Fatalf("validate ready observation: %v", err)
	}

	if _, err := ObserveUsageAction(
		plan,
		"canonical_generation:complete_declaration",
		UsageObservationConfirmed,
		"generation_quality",
		1,
		plan.PlanDigest,
		time.Date(2026, 9, 28, 0, 0, 0, 0, time.UTC),
	); err == nil {
		t.Fatal("expected deferred action confirmation to be rejected")
	}

	unknown, err := ObserveUsageAction(
		plan,
		"canonical_generation:complete_declaration",
		UsageObservationUnknown,
		"generation_quality",
		0,
		plan.PlanDigest,
		time.Date(2026, 9, 28, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("observe deferred action: %v", err)
	}
	if unknown.Kind != UsageObservationUnknown || !unknown.NonAuthorizing {
		t.Fatalf("unexpected deferred observation: %#v", unknown)
	}
}

