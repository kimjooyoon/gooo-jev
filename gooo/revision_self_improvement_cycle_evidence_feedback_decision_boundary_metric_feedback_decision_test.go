package gooo

import "testing"

func evidenceFeedbackDecisionObservationInput(t *testing.T) RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackObservation {
	t.Helper()
	metric := evidenceFeedbackDecisionBoundaryMetricInput(t)
	feedback, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback(metric)
	if err != nil {
		t.Fatalf("observe evidence feedback decision boundary metric feedback: %v", err)
	}
	return feedback
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionBindsObservation(t *testing.T) {
	feedback := evidenceFeedbackDecisionObservationInput(t)
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecision(feedback)
	if err != nil {
		t.Fatalf("observe evidence feedback decision: %v", err)
	}
	if got.Status != "BOUND" ||
		got.MetricSignal != "boundary-complete" ||
		got.FeedbackSignal != "observe" ||
		got.DecisionSignal != "observe-next-cycle" ||
		got.DecisionReason != "complete boundary feedback permits next observation" ||
		!got.DecisionAligned ||
		got.RequiresReplan {
		t.Fatalf("unexpected evidence feedback decision: %#v", got)
	}
	if got.FeedbackDigest != feedback.FeedbackDigest ||
		!got.RequiresObservation ||
		!got.NonExecuting ||
		!got.NonAuthorizing {
		t.Fatalf("decision lost feedback provenance or safety: %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate evidence feedback decision: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionPreservesUnknown(t *testing.T) {
	feedback := evidenceFeedbackDecisionObservationInput(t)
	feedback.FeedbackDigest = digestString("tampered")
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecision(feedback)
	if err == nil {
		t.Fatal("expected tampered feedback error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback-decision-input" {
		t.Fatalf("unexpected unknown evidence feedback decision: %#v", got)
	}
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionIsDeterministic(t *testing.T) {
	feedback := evidenceFeedbackDecisionObservationInput(t)
	first, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecision(feedback)
	if err != nil {
		t.Fatalf("first evidence feedback decision: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecision(feedback)
	if err != nil {
		t.Fatalf("second evidence feedback decision: %v", err)
	}
	if first != second {
		t.Fatalf("evidence feedback decisions differ: %#v != %#v", first, second)
	}
}
