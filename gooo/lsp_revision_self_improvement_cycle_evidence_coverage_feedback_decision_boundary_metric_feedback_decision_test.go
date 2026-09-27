package gooo

import "testing"

func lspEvidenceFeedbackDecisionObservationInput(t *testing.T) RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionObservation {
	t.Helper()
	feedback := evidenceFeedbackDecisionObservationInput(t)
	decision, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecision(feedback)
	if err != nil {
		t.Fatalf("observe evidence feedback decision: %v", err)
	}
	return decision
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionBindsProjection(t *testing.T) {
	decision := lspEvidenceFeedbackDecisionObservationInput(t)
	got, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecision(decision)
	if err != nil {
		t.Fatalf("observe lsp evidence feedback decision: %v", err)
	}
	if got.Status != "BOUND" ||
		got.MetricSignal != "boundary-complete" ||
		got.FeedbackSignal != "observe" ||
		got.DecisionSignal != "observe-next-cycle" ||
		!got.DecisionAligned ||
		got.RequiresReplan {
		t.Fatalf("unexpected lsp evidence feedback decision: %#v", got)
	}
	if got.FeedbackDigest != decision.FeedbackDigest ||
		got.DecisionDigest != decision.DecisionDigest ||
		!got.RequiresObservation {
		t.Fatalf("projection lost decision provenance: %#v", got)
	}
	if !got.ReadOnly || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("projection safety flags = %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate lsp evidence feedback decision: %v", err)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionPreservesUnknown(t *testing.T) {
	decision := lspEvidenceFeedbackDecisionObservationInput(t)
	decision.DecisionDigest = digestString("tampered")
	got, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecision(decision)
	if err == nil {
		t.Fatal("expected tampered decision error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback-decision-input" {
		t.Fatalf("unexpected unknown lsp evidence feedback decision: %#v", got)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionIsDeterministic(t *testing.T) {
	decision := lspEvidenceFeedbackDecisionObservationInput(t)
	first, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecision(decision)
	if err != nil {
		t.Fatalf("first lsp evidence feedback decision: %v", err)
	}
	second, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecision(decision)
	if err != nil {
		t.Fatalf("second lsp evidence feedback decision: %v", err)
	}
	if first != second {
		t.Fatalf("lsp evidence feedback decisions differ: %#v != %#v", first, second)
	}
}
