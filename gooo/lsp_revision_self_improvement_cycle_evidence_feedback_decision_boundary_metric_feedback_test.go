package gooo

import "testing"

func lspEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackInput(t *testing.T) RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackObservation {
	t.Helper()
	metric := evidenceFeedbackDecisionBoundaryMetricInput(t)
	feedback, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback(metric)
	if err != nil {
		t.Fatalf("observe evidence feedback decision boundary metric feedback: %v", err)
	}
	return feedback
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackBindsProjection(t *testing.T) {
	feedback := lspEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackInput(t)
	got, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback(feedback)
	if err != nil {
		t.Fatalf("observe lsp evidence feedback decision boundary metric feedback: %v", err)
	}
	if got.Status != "BOUND" ||
		got.MetricName != "evidence-feedback-decision-boundary-provenance-link-count" ||
		got.MetricSignal != "boundary-complete" ||
		got.FeedbackSignal != "observe" {
		t.Fatalf("unexpected lsp evidence feedback decision boundary metric feedback: %#v", got)
	}
	if got.BoundaryMetricDigest != feedback.BoundaryMetricDigest ||
		got.FeedbackDigest != feedback.FeedbackDigest ||
		!got.RequiresObservation {
		t.Fatalf("projection lost feedback provenance: %#v", got)
	}
	if !got.ReadOnly || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("projection safety flags = %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate lsp evidence feedback decision boundary metric feedback: %v", err)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackPreservesUnknown(t *testing.T) {
	feedback := lspEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackInput(t)
	feedback.FeedbackDigest = digestString("tampered")
	got, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback(feedback)
	if err == nil {
		t.Fatal("expected tampered feedback error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback-input" {
		t.Fatalf("unexpected unknown lsp evidence feedback decision boundary metric feedback: %#v", got)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackIsDeterministic(t *testing.T) {
	feedback := lspEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackInput(t)
	first, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback(feedback)
	if err != nil {
		t.Fatalf("first lsp evidence feedback decision boundary metric feedback: %v", err)
	}
	second, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback(feedback)
	if err != nil {
		t.Fatalf("second lsp evidence feedback decision boundary metric feedback: %v", err)
	}
	if first != second {
		t.Fatalf("lsp evidence feedback decision boundary metric feedback differs: %#v != %#v", first, second)
	}
}
