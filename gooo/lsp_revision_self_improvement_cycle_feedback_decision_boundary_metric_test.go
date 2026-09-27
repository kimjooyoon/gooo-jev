package gooo

import "testing"

func lspBoundaryMetricInput(t *testing.T) RevisionSelfImprovementCycleFeedbackDecisionBoundaryMetricObservation {
	t.Helper()
	cycle, bridge, link := cycleFeedbackDecisionBoundaryInputs(t)
	boundary, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionBoundary(cycle, bridge, link)
	if err != nil {
		t.Fatalf("observe cycle feedback decision boundary: %v", err)
	}
	metric, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetric(boundary)
	if err != nil {
		t.Fatalf("observe boundary metric: %v", err)
	}
	return metric
}

func TestObserveLSPRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetricBindsProjection(t *testing.T) {
	metric := lspBoundaryMetricInput(t)
	got, err := ObserveLSPRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetric(metric)
	if err != nil {
		t.Fatalf("observe lsp boundary metric: %v", err)
	}
	if got.Status != "BOUND" || got.MetricName != "provenance-link-count" || got.MetricSignal != "boundary-complete" {
		t.Fatalf("unexpected lsp boundary metric: %#v", got)
	}
	if got.LinkedStageCount != 3 || got.RequiredStageCount != 3 || got.LinkedDigestCount != 3 || got.BoundaryDigest != metric.BoundaryDigest {
		t.Fatalf("projection lost boundary metric evidence: %#v", got)
	}
	if !got.ReadOnly || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("projection safety flags = %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate lsp boundary metric: %v", err)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetricPreservesUnknown(t *testing.T) {
	metric := lspBoundaryMetricInput(t)
	metric.MetricDigest = digestString("tampered")
	got, err := ObserveLSPRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetric(metric)
	if err == nil {
		t.Fatal("expected tampered boundary metric error")
	}
	if got.Status != "UNKNOWN" || got.MissingStage != "lsp-revision-self-improvement-cycle-feedback-decision-boundary-metric-input" {
		t.Fatalf("unexpected unknown lsp boundary metric: %#v", got)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetricIsDeterministic(t *testing.T) {
	metric := lspBoundaryMetricInput(t)
	first, err := ObserveLSPRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetric(metric)
	if err != nil {
		t.Fatalf("first lsp boundary metric: %v", err)
	}
	second, err := ObserveLSPRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetric(metric)
	if err != nil {
		t.Fatalf("second lsp boundary metric: %v", err)
	}
	if first != second {
		t.Fatalf("lsp boundary metrics differ: %#v != %#v", first, second)
	}
}
