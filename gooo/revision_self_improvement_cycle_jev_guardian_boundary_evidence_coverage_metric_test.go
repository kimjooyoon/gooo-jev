package gooo

import "testing"

func TestRevisionSelfImprovementCycleJEVGuardianBoundaryEvidenceCoverageMetricPreservesUnknown(t *testing.T) {
	reverseObservation, _ := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation(
		RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation{},
	)
	closureMetric := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetric(reverseObservation)
	projection := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjection(closureMetric)
	guardian := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjectionGuardianBoundaryObservation(projection)
	metric := ObserveRevisionSelfImprovementCycleJEVGuardianBoundaryEvidenceCoverageMetric(guardian)

	if metric.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN", metric.Status)
	}
	if metric.RequiredDigestCount != 3 || metric.LinkedDigestCount != 0 || metric.CoverageMilli != 0 {
		t.Fatalf("coverage = %d/%d at %d milli, want 0/3 at 0 milli", metric.LinkedDigestCount, metric.RequiredDigestCount, metric.CoverageMilli)
	}
	if metric.MissingStage == "" {
		t.Fatal("missing stage must be preserved")
	}
	if err := metric.Validate(); err != nil {
		t.Fatalf("metric should validate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVGuardianBoundaryEvidenceCoverageMetricRejectsTamperedGuardian(t *testing.T) {
	reverseObservation, _ := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation(
		RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation{},
	)
	closureMetric := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetric(reverseObservation)
	projection := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjection(closureMetric)
	guardian := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjectionGuardianBoundaryObservation(projection)
	guardian.DecisionState = "ALLOW"
	metric := ObserveRevisionSelfImprovementCycleJEVGuardianBoundaryEvidenceCoverageMetric(guardian)

	if metric.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN", metric.Status)
	}
	if metric.MissingStage == "" {
		t.Fatal("tampered guardian must preserve an unresolved stage")
	}
}