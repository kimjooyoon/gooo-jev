package gooo

import "testing"

func TestRevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjectionReverseObservationPreservesUnknown(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjectionReverseObservation(
		RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjection{},
	)

	if observation.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN", observation.Status)
	}
	if observation.MissingStage == "" || observation.MissingStageIndex != -1 {
		t.Fatalf("missing stage = %q/%d, want unresolved stage index -1", observation.MissingStage, observation.MissingStageIndex)
	}
	if observation.CoverageMilli != 0 || observation.CoverageBand != "unknown" {
		t.Fatalf("coverage = %d milli/%q, want 0 milli/unknown", observation.CoverageMilli, observation.CoverageBand)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("observation should validate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjectionReverseObservationRejectsTamperedProjection(t *testing.T) {
	input := RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjection{
		Status:          "BOUND",
		MetricName:      jevGenerationTraceCoverageMetricName,
		MissingStageIndex: 0,
		CoverageMilli:    1000,
		CoverageBand:    "high",
		NonExecuting:    true,
		NonAuthorizing:  true,
	}
	observation := ObserveRevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjectionReverseObservation(input)

	if observation.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN for tampered projection", observation.Status)
	}
	if observation.MissingStage == "" {
		t.Fatal("tampered projection must preserve an unresolved stage")
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("observation should validate after preserving UNKNOWN: %v", err)
	}
}
