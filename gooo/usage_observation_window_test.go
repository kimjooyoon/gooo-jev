package gooo

import (
	"testing"
	"time"
)

func TestMeasureUsageObservationWindowExcludesUnknownFromMean(t *testing.T) {
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
		"usage_quality",
		1,
		plan.PlanDigest,
		time.Date(2026, 9, 28, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("observe ready action: %v", err)
	}
	unknown, err := ObserveUsageAction(
		plan,
		"canonical_generation:complete_declaration",
		UsageObservationUnknown,
		"usage_quality",
		0,
		plan.PlanDigest,
		time.Date(2026, 9, 28, 0, 0, 1, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("observe unknown action: %v", err)
	}

	window, err := MeasureUsageObservationWindow([]UsageObservation{ready, unknown}, "usage_quality")
	if err != nil {
		t.Fatalf("measure usage window: %v", err)
	}
	if err := window.Validate(); err != nil {
		t.Fatalf("validate usage window: %v", err)
	}
	if window.ObservationCount != 2 || window.KnownCount != 1 || window.UnknownCount != 1 {
		t.Fatalf("unexpected observation counts: %#v", window)
	}
	if window.KnownMean != 1 || window.ObservedCoverage != 0.5 {
		t.Fatalf("unknown observation contaminated metrics: %#v", window)
	}
	if !window.NonExecuting || !window.NonAuthorizing {
		t.Fatal("usage window must remain non-executing and non-authorizing")
	}
}

