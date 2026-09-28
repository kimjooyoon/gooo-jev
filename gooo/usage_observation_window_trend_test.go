package gooo

import (
	"testing"
	"time"
)

func TestCompareUsageObservationWindowsPreservesCoverageAndDirection(t *testing.T) {
	source := "package support\nnamespace triage\nentity ticket\nproperty\n"
	discovery := DiscoverUsage(source, "pro")
	plan, err := PlanUsageActions(discovery)
	if err != nil {
		t.Fatalf("plan usage actions: %v", err)
	}
	previousObservation, err := ObserveUsageAction(
		plan,
		"syntax_completion:edit_declaration",
		UsageObservationConfirmed,
		"usage_quality",
		1,
		plan.PlanDigest,
		time.Date(2026, 9, 28, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("observe previous action: %v", err)
	}
	previousUnknown, err := ObserveUsageAction(
		plan,
		"canonical_generation:complete_declaration",
		UsageObservationUnknown,
		"usage_quality",
		0,
		plan.PlanDigest,
		time.Date(2026, 9, 28, 0, 0, 1, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("observe previous unknown: %v", err)
	}
	currentObservation, err := ObserveUsageAction(
		plan,
		"syntax_completion:edit_declaration",
		UsageObservationConfirmed,
		"usage_quality",
		2,
		plan.PlanDigest,
		time.Date(2026, 9, 28, 0, 0, 2, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("observe current action: %v", err)
	}

	previous, err := MeasureUsageObservationWindow([]UsageObservation{previousObservation, previousUnknown}, "usage_quality")
	if err != nil {
		t.Fatalf("measure previous window: %v", err)
	}
	current, err := MeasureUsageObservationWindow([]UsageObservation{currentObservation}, "usage_quality")
	if err != nil {
		t.Fatalf("measure current window: %v", err)
	}
	trend, err := CompareUsageObservationWindows(
		previous,
		current,
		time.Date(2026, 9, 28, 0, 1, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("compare usage windows: %v", err)
	}
	if err := trend.Validate(); err != nil {
		t.Fatalf("validate usage trend: %v", err)
	}
	if trend.Direction != UsageTrendRising || trend.KnownMeanDelta != 1 {
		t.Fatalf("unexpected usage trend: %#v", trend)
	}
	if trend.PreviousCoverage != 0.5 || trend.CurrentCoverage != 1 || trend.PreviousUnknownCount != 1 || trend.CurrentUnknownCount != 0 {
		t.Fatalf("coverage and unknown counts were lost: %#v", trend)
	}
}

