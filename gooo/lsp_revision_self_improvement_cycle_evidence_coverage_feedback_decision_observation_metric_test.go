package gooo

import "testing"

func lspEvidenceFeedbackDecisionObservationMetricInput(t *testing.T) RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricObservation {
	t.Helper()
	decision := evidenceFeedbackDecisionObservationMetricInput(t)
	metric, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(decision)
	if err != nil {
		t.Fatalf("observe decision observation metric: %v", err)
	}
	return metric
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricBindsProjection(t *testing.T) {
	metric := lspEvidenceFeedbackDecisionObservationMetricInput(t)
	got, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(metric)
	if err != nil {
		t.Fatalf("observe lsp decision observation metric: %v", err)
	}
	if got.Status != "BOUND" ||
		got.MetricName != "evidence-feedback-decision-observation-provenance-link-count" ||
		got.MetricSignal != "decision-observation-complete" ||
		got.LinkedStageCount != 2 ||
		got.RequiredStageCount != 2 ||
		got.LinkedDigestCount != 3 {
		t.Fatalf("unexpected lsp decision observation metric: %#v", got)
	}
	if got.DecisionObservationDigest != metric.DecisionObservationDigest ||
		got.MetricDigest != metric.MetricDigest {
		t.Fatalf("projection lost metric provenance: %#v", got)
	}
	if !got.ReadOnly || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("projection safety flags = %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate lsp decision observation metric: %v", err)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricPreservesUnknown(t *testing.T) {
	metric := lspEvidenceFeedbackDecisionObservationMetricInput(t)
	metric.MetricDigest = digestString("tampered")
	got, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(metric)
	if err == nil {
		t.Fatal("expected tampered metric error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-input" {
		t.Fatalf("unexpected unknown lsp decision observation metric: %#v", got)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricIsDeterministic(t *testing.T) {
	metric := lspEvidenceFeedbackDecisionObservationMetricInput(t)
	first, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(metric)
	if err != nil {
		t.Fatalf("first lsp decision observation metric: %v", err)
	}
	second, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(metric)
	if err != nil {
		t.Fatalf("second lsp decision observation metric: %v", err)
	}
	if first != second {
		t.Fatalf("lsp decision observation metrics differ: %#v != %#v", first, second)
	}
}
