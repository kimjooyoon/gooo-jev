package gooo

import "testing"

func TestRevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjectionPreservesUnknown(t *testing.T) {
	projection := ObserveRevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjection(
		RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetric{},
	)

	if projection.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN", projection.Status)
	}
	if projection.MissingStage == "" || projection.MissingStageIndex != -1 {
		t.Fatalf("missing stage = %q/%d, want unresolved stage index -1", projection.MissingStage, projection.MissingStageIndex)
	}
	if projection.CoverageMilli != 0 || projection.CoverageBand != "unknown" {
		t.Fatalf("coverage = %d milli/%q, want 0 milli/unknown", projection.CoverageMilli, projection.CoverageBand)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("projection should validate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjectionRejectsTamperedMetric(t *testing.T) {
	input := RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetric{
		Status:       "BOUND",
		MetricName:   jevGenerationTraceCoverageMetricName,
		CoverageMilli: 1000,
		CoverageBand: "high",
		NonExecuting: true,
		NonAuthorizing: true,
	}
	projection := ObserveRevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjection(input)

	if projection.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN for tampered metric", projection.Status)
	}
	if projection.MissingStage == "" {
		t.Fatal("tampered metric must preserve an unresolved stage")
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("projection should validate after preserving UNKNOWN: %v", err)
	}
}
