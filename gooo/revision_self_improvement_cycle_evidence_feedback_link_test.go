package gooo

import "testing"

func cycleEvidenceCoverageFeedbackInputs(t *testing.T) (
	RevisionSelfImprovementCycleFeedbackBridgeObservation,
	RevisionSelfImprovementCycleEvidenceCoverageMetricObservation,
	RevisionSelfImprovementCycleMetricFeedback,
) {
	t.Helper()
	cycle, metric, feedback := cycleFeedbackBridgeInputs(t)
	bridge, err := ObserveRevisionSelfImprovementCycleFeedbackBridge(cycle, metric, feedback)
	if err != nil {
		t.Fatalf("observe cycle feedback bridge: %v", err)
	}
	coverage, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageMetric(bridge)
	if err != nil {
		t.Fatalf("observe cycle evidence coverage metric: %v", err)
	}
	return bridge, coverage, feedback
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackBindsLinks(t *testing.T) {
	bridge, coverage, feedback := cycleEvidenceCoverageFeedbackInputs(t)
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedback(bridge, coverage, feedback)
	if err != nil {
		t.Fatalf("observe cycle evidence coverage feedback: %v", err)
	}
	if got.Status != "BOUND" ||
		got.LinkSignal != "cycle-evidence-feedback-linked" ||
		!got.FeedbackAligned {
		t.Fatalf("unexpected cycle evidence coverage feedback: %#v", got)
	}
	if got.BridgeDigest != bridge.ObservationDigest ||
		got.CoverageMetricDigest != coverage.ObservationDigest ||
		got.FeedbackDigest != feedback.FeedbackDigest {
		t.Fatalf("link observation lost stage digests: %#v", got)
	}
	if got.SourceDigest != bridge.SourceDigest ||
		got.CandidateSourceDigest != bridge.CandidateSourceDigest ||
		got.GeneratedIRDigest != bridge.GeneratedIRDigest ||
		got.ReverseObservationDigest != bridge.ReverseObservationDigest {
		t.Fatalf("link observation lost evidence digests: %#v", got)
	}
	if got.CoverageSignal != "evidence-complete" ||
		got.FeedbackSignal != feedback.FeedbackSignal ||
		got.FeedbackReason != feedback.FeedbackReason ||
		!got.RequiresObservation || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("link observation lost feedback semantics: %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate cycle evidence coverage feedback: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackPreservesUnknownLink(t *testing.T) {
	bridge, coverage, feedback := cycleEvidenceCoverageFeedbackInputs(t)
	feedback.CycleDigest = digestString("tampered")
	feedback.FeedbackDigest = digestRevisionSelfImprovementCycleMetricFeedback(feedback)
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedback(bridge, coverage, feedback)
	if err == nil {
		t.Fatal("expected feedback cycle link error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "revision-self-improvement-cycle-evidence-coverage-feedback-feedback-cycle-link" {
		t.Fatalf("unexpected unknown cycle evidence coverage feedback: %#v", got)
	}
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackIsDeterministic(t *testing.T) {
	bridge, coverage, feedback := cycleEvidenceCoverageFeedbackInputs(t)
	first, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedback(bridge, coverage, feedback)
	if err != nil {
		t.Fatalf("first cycle evidence coverage feedback: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedback(bridge, coverage, feedback)
	if err != nil {
		t.Fatalf("second cycle evidence coverage feedback: %v", err)
	}
	if first != second {
		t.Fatalf("cycle evidence coverage feedback differs: %#v != %#v", first, second)
	}
}
