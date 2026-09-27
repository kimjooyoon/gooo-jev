package gooo

import "testing"

func TestRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetricPreservesUnknown(t *testing.T) {
	reverseObservation := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation(
		RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation{},
	)
	closureMetric := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetric(reverseObservation)

	if closureMetric.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN", closureMetric.Status)
	}
	if closureMetric.MissingStage == "" {
		t.Fatal("missing stage must be preserved")
	}
	if !closureMetric.NonExecuting || !closureMetric.NonAuthorizing {
		t.Fatal("closure metric must remain non-executing and non-authorizing")
	}
	if err := closureMetric.Validate(); err != nil {
		t.Fatalf("closure metric should validate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetricRejectsTamperedObservation(t *testing.T) {
	reverseObservation := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation(
		RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation{},
	)
	reverseObservation.ObservationDigest = digestString("tampered")
	closureMetric := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetric(reverseObservation)

	if closureMetric.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN", closureMetric.Status)
	}
	if closureMetric.MissingStage == "" {
		t.Fatal("missing stage must remain unresolved")
	}
}