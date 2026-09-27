package gooo

import "testing"

func lspEvidenceCoverageFeedbackDecisionBoundaryMetricInput(t *testing.T) RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricObservation {
	t.Helper()
	cycle, evidenceFeedback, boundary := cycleEvidenceCoverageFeedbackDecisionBoundaryInputs(t)
	gatedBoundary, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(cycle, evidenceFeedback, boundary)
	if err != nil {
		t.Fatalf("observe cycle evidence feedback decision boundary: %v", err)
	}
	metric, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric(gatedBoundary)
	if err != nil {
		t.Fatalf("observe evidence feedback decision boundary metric: %v", err)
	}
	return metric
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricBindsProjection(t *testing.T) {
	metric := lspEvidenceCoverageFeedbackDecisionBoundaryMetricInput(t)
	got, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric(metric)
	if err != nil {
		t.Fatalf("observe lsp evidence feedback decision boundary metric: %v", err)
	}
	if got.Status != "BOUND" ||
		got.MetricName != "evidence-feedback-decision-boundary-provenance-link-count" ||
		got.MetricSignal != "boundary-complete" {
		t.Fatalf("unexpected lsp evidence feedback decision boundary metric: %#v", got)
	}
	if got.LinkedStageCount != 4 ||
		got.RequiredStageCount != 4 ||
		got.LinkedDigestCount != 4 ||
		got.BoundaryDigest != metric.BoundaryDigest {
		t.Fatalf("projection lost four-link metric: %#v", got)
	}
	if !got.ReadOnly || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("projection safety flags = %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate lsp evidence feedback decision boundary metric: %v", err)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricPreservesUnknown(t *testing.T) {
	metric := lspEvidenceCoverageFeedbackDecisionBoundaryMetricInput(t)
	metric.MetricDigest = digestString("tampered")
	got, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric(metric)
	if err == nil {
		t.Fatal("expected tampered evidence feedback metric error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-input" {
		t.Fatalf("unexpected unknown lsp evidence feedback decision boundary metric: %#v", got)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricIsDeterministic(t *testing.T) {
	metric := lspEvidenceCoverageFeedbackDecisionBoundaryMetricInput(t)
	first, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric(metric)
	if err != nil {
		t.Fatalf("first lsp evidence feedback decision boundary metric: %v", err)
	}
	second, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric(metric)
	if err != nil {
		t.Fatalf("second lsp evidence feedback decision boundary metric: %v", err)
	}
	if first != second {
		t.Fatalf("lsp evidence feedback decision boundary metrics differ: %#v != %#v", first, second)
	}
}
