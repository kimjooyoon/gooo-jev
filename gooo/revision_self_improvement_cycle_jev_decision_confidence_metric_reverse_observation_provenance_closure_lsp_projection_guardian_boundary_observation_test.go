package gooo

import "testing"

func TestRevisionSelfImprovementCycleJEVGuardianBoundaryObservationDefersUnknown(t *testing.T) {
	reverseObservation, _ := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation(
		RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation{},
	)
	closureMetric := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetric(reverseObservation)
	projection := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjection(closureMetric)
	guardian := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjectionGuardianBoundaryObservation(projection)

	if guardian.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN", guardian.Status)
	}
	if guardian.DecisionState != jevGuardianBoundaryDecisionDefer {
		t.Fatalf("decision state = %q, want DEFER", guardian.DecisionState)
	}
	if guardian.MissingStage == "" {
		t.Fatal("missing stage must be preserved")
	}
	if !guardian.NonExecuting || !guardian.NonAuthorizing {
		t.Fatal("guardian observation must not execute or authorize")
	}
	if err := guardian.Validate(); err != nil {
		t.Fatalf("guardian observation should validate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVGuardianBoundaryObservationRejectsAuthorizationMutation(t *testing.T) {
	reverseObservation, _ := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation(
		RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation{},
	)
	closureMetric := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetric(reverseObservation)
	projection := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjection(closureMetric)
	guardian := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjectionGuardianBoundaryObservation(projection)
	guardian.DecisionState = "ALLOW"

	if err := guardian.Validate(); err == nil {
		t.Fatal("authorization mutation must be rejected")
	}
}