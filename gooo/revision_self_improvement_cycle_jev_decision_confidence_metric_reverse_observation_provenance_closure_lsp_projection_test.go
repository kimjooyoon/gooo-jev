package gooo

import "testing"

func TestRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjectionPreservesUnknown(t *testing.T) {
	reverseObservation, _ := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation(
		RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation{},
	)
	closureMetric := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetric(reverseObservation)
	projection := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjection(closureMetric)

	if projection.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN", projection.Status)
	}
	if projection.MissingStage == "" || projection.MissingStageIndex != -1 {
		t.Fatalf("missing stage = %q index = %d, want unresolved stage and -1", projection.MissingStage, projection.MissingStageIndex)
	}
	if projection.EvidencePrefixDigest == "" {
		t.Fatal("unknown projection should preserve the available evidence prefix")
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("projection should validate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjectionRejectsTampering(t *testing.T) {
	reverseObservation, _ := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation(
		RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation{},
	)
	closureMetric := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetric(reverseObservation)
	closureMetric.ObservationDigest = digestString("tampered")
	projection := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjection(closureMetric)

	if projection.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN", projection.Status)
	}
	if projection.MissingStage == "" {
		t.Fatal("tampering must preserve an unresolved stage")
	}
}