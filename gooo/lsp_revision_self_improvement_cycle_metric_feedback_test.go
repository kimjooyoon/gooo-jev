package gooo

import "testing"

func lspCycleMetricFeedbackInput(t *testing.T) RevisionSelfImprovementCycleMetricFeedback {
	t.Helper()
	metric := cycleMetricObservation(t)
	feedback, err := ObserveRevisionSelfImprovementCycleMetricFeedback(metric)
	if err != nil {
		t.Fatalf("observe cycle metric feedback: %v", err)
	}
	return feedback
}

func TestObserveLSPRevisionSelfImprovementCycleMetricFeedbackBindsProjection(t *testing.T) {
	feedback := lspCycleMetricFeedbackInput(t)
	got, err := ObserveLSPRevisionSelfImprovementCycleMetricFeedback(feedback)
	if err != nil {
		t.Fatalf("observe lsp cycle metric feedback: %v", err)
	}
	if got.Status != "BOUND" || got.MissingStage != "" {
		t.Fatalf("unexpected lsp cycle metric feedback: %#v", got)
	}
	if got.CycleMetricDigest != feedback.CycleMetricDigest ||
		got.FeedbackSignal != feedback.FeedbackSignal ||
		got.FeedbackReason != feedback.FeedbackReason ||
		got.MetricDelta != feedback.MetricDelta {
		t.Fatalf("projection lost metric feedback evidence: %#v", got)
	}
	if !got.ReadOnly || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("projection safety flags = %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate lsp cycle metric feedback: %v", err)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleMetricFeedbackPreservesUnknown(t *testing.T) {
	feedback := lspCycleMetricFeedbackInput(t)
	feedback.FeedbackDigest = digestString("tampered")
	got, err := ObserveLSPRevisionSelfImprovementCycleMetricFeedback(feedback)
	if err == nil {
		t.Fatal("expected tampered feedback error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "lsp-revision-self-improvement-cycle-metric-feedback-input" {
		t.Fatalf("unexpected unknown lsp metric feedback: %#v", got)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleMetricFeedbackIsDeterministic(t *testing.T) {
	feedback := lspCycleMetricFeedbackInput(t)
	first, err := ObserveLSPRevisionSelfImprovementCycleMetricFeedback(feedback)
	if err != nil {
		t.Fatalf("first lsp cycle metric feedback: %v", err)
	}
	second, err := ObserveLSPRevisionSelfImprovementCycleMetricFeedback(feedback)
	if err != nil {
		t.Fatalf("second lsp cycle metric feedback: %v", err)
	}
	if first != second {
		t.Fatalf("lsp metric feedback projections differ: %#v != %#v", first, second)
	}
}