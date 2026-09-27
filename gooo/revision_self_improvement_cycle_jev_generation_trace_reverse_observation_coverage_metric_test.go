package gooo

import "testing"

func TestRevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricPreservesUnknown(t *testing.T) {
	metric := ObserveRevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetric(
		RevisionSelfImprovementCycleJEVGenerationTraceReverseObservation{},
	)

	if metric.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN", metric.Status)
	}
	if metric.RequiredDigestCount != 6 || metric.LinkedDigestCount != 0 || metric.CoverageMilli != 0 {
		t.Fatalf("coverage = %d/%d at %d milli, want 0/6 at 0 milli", metric.LinkedDigestCount, metric.RequiredDigestCount, metric.CoverageMilli)
	}
	if metric.MissingStage == "" {
		t.Fatal("missing stage must be preserved")
	}
	if err := metric.Validate(); err != nil {
		t.Fatalf("metric should validate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricRejectsTamperedInput(t *testing.T) {
	input := RevisionSelfImprovementCycleJEVGenerationTraceReverseObservation{
		Status:         "BOUND",
		MetricName:     jevGenerationTraceReverseMetricName,
		ReverseSignal:  jevGenerationTraceReverseComplete,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	metric := ObserveRevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetric(input)

	if metric.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN for tampered input", metric.Status)
	}
	if metric.MissingStage == "" {
		t.Fatal("tampered input must preserve an unresolved stage")
	}
	if err := metric.Validate(); err != nil {
		t.Fatalf("metric should validate after preserving UNKNOWN: %v", err)
	}
}
