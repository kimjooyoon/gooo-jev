package gooo

import "testing"

func cycleEvidenceCoverageFeedbackDecisionBoundaryInputs(t *testing.T) (
	RevisionSelfImprovementCycleObservation,
	RevisionSelfImprovementCycleEvidenceCoverageFeedbackObservation,
	RevisionSelfImprovementCycleFeedbackDecisionBoundaryObservation,
) {
	t.Helper()
	cycle, bridge, decisionLink := cycleFeedbackDecisionBoundaryInputs(t)
	boundary, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionBoundary(cycle, bridge, decisionLink)
	if err != nil {
		t.Fatalf("observe cycle feedback decision boundary: %v", err)
	}
	_, _, feedback := cycleFeedbackBridgeInputs(t)
	coverage, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageMetric(bridge)
	if err != nil {
		t.Fatalf("observe cycle evidence coverage metric: %v", err)
	}
	evidenceFeedback, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedback(bridge, coverage, feedback)
	if err != nil {
		t.Fatalf("observe cycle evidence coverage feedback: %v", err)
	}
	return cycle, evidenceFeedback, boundary
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryBindsEvidence(t *testing.T) {
	cycle, evidenceFeedback, boundary := cycleEvidenceCoverageFeedbackDecisionBoundaryInputs(t)
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(cycle, evidenceFeedback, boundary)
	if err != nil {
		t.Fatalf("observe cycle evidence feedback decision boundary: %v", err)
	}
	if got.Status != "BOUND" ||
		got.BoundarySignal != "cycle-evidence-feedback-decision-boundary-linked" ||
		!got.EvidenceFeedbackAligned ||
		!got.BoundaryAligned {
		t.Fatalf("unexpected cycle evidence feedback decision boundary: %#v", got)
	}
	if got.CycleDigest != cycle.ObservationDigest ||
		got.BridgeDigest != evidenceFeedback.BridgeDigest ||
		got.EvidenceFeedbackDigest != evidenceFeedback.ObservationDigest ||
		got.DecisionBoundaryDigest != boundary.ObservationDigest {
		t.Fatalf("boundary lost provenance links: %#v", got)
	}
	if got.CoverageSignal != "evidence-complete" ||
		got.FeedbackSignal != evidenceFeedback.FeedbackSignal ||
		got.FeedbackReason != evidenceFeedback.FeedbackReason ||
		got.DecisionSignal != boundary.DecisionSignal ||
		!got.NonExecuting ||
		!got.NonAuthorizing {
		t.Fatalf("boundary lost evidence or decision semantics: %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate cycle evidence feedback decision boundary: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryPreservesUnknown(t *testing.T) {
	cycle, evidenceFeedback, boundary := cycleEvidenceCoverageFeedbackDecisionBoundaryInputs(t)
	evidenceFeedback.CoverageMetricDigest = digestString("tampered")
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(cycle, evidenceFeedback, boundary)
	if err == nil {
		t.Fatal("expected tampered evidence feedback error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-evidence-feedback" {
		t.Fatalf("unexpected unknown cycle evidence feedback decision boundary: %#v", got)
	}
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryIsDeterministic(t *testing.T) {
	cycle, evidenceFeedback, boundary := cycleEvidenceCoverageFeedbackDecisionBoundaryInputs(t)
	first, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(cycle, evidenceFeedback, boundary)
	if err != nil {
		t.Fatalf("first cycle evidence feedback decision boundary: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(cycle, evidenceFeedback, boundary)
	if err != nil {
		t.Fatalf("second cycle evidence feedback decision boundary: %v", err)
	}
	if first != second {
		t.Fatalf("cycle evidence feedback decision boundaries differ: %#v != %#v", first, second)
	}
}
