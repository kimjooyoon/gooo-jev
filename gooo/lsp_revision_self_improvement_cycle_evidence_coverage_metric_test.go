package gooo

import "testing"

func lspEvidenceCoverageMetricInput(t *testing.T) RevisionSelfImprovementCycleEvidenceCoverageMetricObservation {
	t.Helper()
	cycle, metric, feedback := cycleFeedbackBridgeInputs(t)
	bridge, err := ObserveRevisionSelfImprovementCycleFeedbackBridge(cycle, metric, feedback)
	if err != nil {
		t.Fatalf("observe cycle feedback bridge: %v", err)
	}
	coverage, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageMetric(bridge)
	if err != nil {
		t.Fatalf("observe evidence coverage metric: %v", err)
	}
	return coverage
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageMetricBindsProjection(t *testing.T) {
	metric := lspEvidenceCoverageMetricInput(t)
	got, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageMetric(metric)
	if err != nil {
		t.Fatalf("observe lsp evidence coverage metric: %v", err)
	}
	if got.Status != "BOUND" || got.MetricName != "source-ir-generation-reverse-observation-coverage" || got.CoverageSignal != "evidence-complete" {
		t.Fatalf("unexpected lsp evidence coverage metric: %#v", got)
	}
	if got.LinkedEvidenceCount != 4 || got.RequiredEvidenceCount != 4 || got.GeneratedIRDigest != metric.GeneratedIRDigest || got.ReverseObservationDigest != metric.ReverseObservationDigest {
		t.Fatalf("projection lost evidence coverage: %#v", got)
	}
	if !got.ReadOnly || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("projection safety flags = %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate lsp evidence coverage metric: %v", err)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageMetricPreservesUnknown(t *testing.T) {
	metric := lspEvidenceCoverageMetricInput(t)
	metric.ReverseObservationDigest = digestString("tampered")
	got, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageMetric(metric)
	if err == nil {
		t.Fatal("expected tampered coverage metric error")
	}
	if got.Status != "UNKNOWN" || got.MissingStage != "lsp-revision-self-improvement-cycle-evidence-coverage-metric-input" {
		t.Fatalf("unexpected unknown lsp evidence coverage metric: %#v", got)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageMetricIsDeterministic(t *testing.T) {
	metric := lspEvidenceCoverageMetricInput(t)
	first, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageMetric(metric)
	if err != nil {
		t.Fatalf("first lsp evidence coverage metric: %v", err)
	}
	second, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageMetric(metric)
	if err != nil {
		t.Fatalf("second lsp evidence coverage metric: %v", err)
	}
	if first != second {
		t.Fatalf("lsp evidence coverage metrics differ: %#v != %#v", first, second)
	}
}
