package gooo

import (
	"testing"
	"time"
)

func TestProposeUsageReplanFromConfirmedWindow(t *testing.T) {
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
		t.Fatalf("observe confirmed action: %v", err)
	}
	window, err := MeasureUsageObservationWindow([]UsageObservation{ready}, "usage_quality")
	if err != nil {
		t.Fatalf("measure usage window: %v", err)
	}

	proposal := ProposeUsageReplan(window)
	if err := proposal.Validate(); err != nil {
		t.Fatalf("validate proposal: %v", err)
	}
	if proposal.Status != "BOUND" || proposal.Disposition != UsageReplanProposeNextIteration {
		t.Fatalf("unexpected confirmed proposal: %#v", proposal)
	}
	if proposal.NextOperation != "measure_next_usage_window" || !proposal.NonExecuting || !proposal.NonAuthorizing {
		t.Fatalf("proposal crossed execution boundary: %#v", proposal)
	}
}

func TestProposeUsageReplanHoldsUnknownWindow(t *testing.T) {
	source := "package support\nnamespace triage\nentity ticket\nproperty\n"
	discovery := DiscoverUsage(source, "pro")
	plan, err := PlanUsageActions(discovery)
	if err != nil {
		t.Fatalf("plan usage actions: %v", err)
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
	window, err := MeasureUsageObservationWindow([]UsageObservation{unknown}, "usage_quality")
	if err != nil {
		t.Fatalf("measure usage window: %v", err)
	}

	proposal := ProposeUsageReplan(window)
	if err := proposal.Validate(); err != nil {
		t.Fatalf("validate unknown proposal: %v", err)
	}
	if proposal.Status != "UNKNOWN" || proposal.Disposition != UsageReplanHoldForEvidence {
		t.Fatalf("unknown usage was promoted: %#v", proposal)
	}
	if proposal.FirstMismatch != "INSUFFICIENT_KNOWN_USAGE_EVIDENCE" || proposal.MissingStageIndex != 0 {
		t.Fatalf("unknown boundary was not preserved: %#v", proposal)
	}
}

func TestProposeUsageReplanReviewsWindowWithoutConfirmation(t *testing.T) {
	source := "package support\nnamespace triage\nentity ticket\nproperty\n"
	discovery := DiscoverUsage(source, "pro")
	plan, err := PlanUsageActions(discovery)
	if err != nil {
		t.Fatalf("plan usage actions: %v", err)
	}
	refuted, err := ObserveUsageAction(
		plan,
		"syntax_completion:edit_declaration",
		UsageObservationRefuted,
		"usage_quality",
		0,
		plan.PlanDigest,
		time.Date(2026, 9, 28, 0, 0, 2, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("observe refuted action: %v", err)
	}
	window, err := MeasureUsageObservationWindow([]UsageObservation{refuted}, "usage_quality")
	if err != nil {
		t.Fatalf("measure usage window: %v", err)
	}

	proposal := ProposeUsageReplan(window)
	if err := proposal.Validate(); err != nil {
		t.Fatalf("validate review proposal: %v", err)
	}
	if proposal.Status != "BOUND" || proposal.Disposition != UsageReplanReviewNoConfirmedOutcome {
		t.Fatalf("refuted window proposed execution: %#v", proposal)
	}
}

func TestProposeUsageReplanPreservesInvalidWindowAsUnknown(t *testing.T) {
	window := UsageObservationWindow{
		EvidenceDigest:   "tampered-window",
		ObservationCount: 1,
		KnownCount:       1,
		ConfirmedCount:   1,
		ObservedCoverage: 1,
		NonExecuting:     true,
		NonAuthorizing:   true,
	}

	proposal := ProposeUsageReplan(window)
	if err := proposal.Validate(); err != nil {
		t.Fatalf("validate invalid-window proposal: %v", err)
	}
	if proposal.Status != "UNKNOWN" || proposal.FirstMismatch != "WINDOW_INTEGRITY" || proposal.Disposition != UsageReplanHoldForEvidence {
		t.Fatalf("invalid window was promoted: %#v", proposal)
	}

	proposal.EvidenceDigest = "tampered-proposal"
	if err := proposal.Validate(); err == nil {
		t.Fatal("tampered proposal was accepted")
	}
}