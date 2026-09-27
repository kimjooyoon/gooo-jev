package gooo

import "testing"

func TestRevisionSelfImprovementCycleJEVGenerationTraceObservationPreservesUnknown(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVGenerationTraceObservation(
		RevisionSelfImprovementCycleJEVGenerationTraceEvidence{},
	)

	if observation.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN", observation.Status)
	}
	if observation.MissingStage == "" {
		t.Fatal("missing stage must be preserved")
	}
	if observation.ContractDigest != "" || observation.IRDigest != "" || observation.GeneratedArtifactDigest != "" {
		t.Fatal("unknown observation must not fabricate generation evidence")
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("observation should validate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVGenerationTraceObservationRejectsUnboundEvidence(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVGenerationTraceObservation(
		RevisionSelfImprovementCycleJEVGenerationTraceEvidence{Status: "BOUND"},
	)

	if observation.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN", observation.Status)
	}
	if observation.MissingStage == "" {
		t.Fatal("unbound evidence must preserve an unresolved stage")
	}
}