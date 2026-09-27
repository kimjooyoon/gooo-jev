package gooo

import "testing"

func lspReverseObservationProvenanceClosureMetricInput(t *testing.T) RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric {
	t.Helper()
	metric := evidenceFeedbackDecisionObservationReverseClosureInput(t)
	closure, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric(metric)
	if err != nil {
		t.Fatalf("observe reverse observation provenance closure metric: %v", err)
	}
	return closure
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetricBindsProjection(t *testing.T) {
	metric := lspReverseObservationProvenanceClosureMetricInput(t)
	got, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric(metric)
	if err != nil {
		t.Fatalf("observe lsp reverse observation provenance closure metric: %v", err)
	}
	if got.Status != "BOUND" ||
		got.LinkedDigestCount != 5 ||
		got.RequiredDigestCount != 5 ||
		got.MetricSignal != "reverse-observation-provenance-complete" {
		t.Fatalf("unexpected lsp reverse observation provenance closure metric: %#v", got)
	}
	if got.MetricObservationDigest != metric.MetricObservationDigest ||
		got.MetricDigest != metric.MetricDigest ||
		got.ReconstructedMetricDigest != metric.ReconstructedMetricDigest ||
		got.DecisionObservationDigest != metric.DecisionObservationDigest ||
		got.ReverseObservationDigest != metric.ReverseObservationDigest ||
		got.ClosureDigest != metric.ClosureDigest {
		t.Fatalf("projection lost closure provenance: %#v", got)
	}
	if !got.ReadOnly || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("projection safety flags = %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate lsp reverse observation provenance closure metric: %v", err)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetricPreservesUnknown(t *testing.T) {
	metric := lspReverseObservationProvenanceClosureMetricInput(t)
	metric.ClosureDigest = digestString("tampered")
	got, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric(metric)
	if err == nil {
		t.Fatal("expected tampered closure metric error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation-provenance-closure-metric-input" {
		t.Fatalf("unexpected unknown lsp reverse observation provenance closure metric: %#v", got)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetricIsDeterministic(t *testing.T) {
	metric := lspReverseObservationProvenanceClosureMetricInput(t)
	first, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric(metric)
	if err != nil {
		t.Fatalf("first lsp reverse observation provenance closure metric: %v", err)
	}
	second, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric(metric)
	if err != nil {
		t.Fatalf("second lsp reverse observation provenance closure metric: %v", err)
	}
	if first != second {
		t.Fatalf("lsp reverse observation provenance closure metrics differ: %#v != %#v", first, second)
	}
}
