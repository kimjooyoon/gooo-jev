package gooo

import "testing"

func TestRevisionSelfImprovementCycleJEVExecutionPlanObservationPreservesUnknown(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVExecutionPlanObservation(
		RevisionSelfImprovementCycleJEVExecutionPlanEvidence{},
	)

	if observation.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN", observation.Status)
	}
	if observation.MissingStage == "" {
		t.Fatal("missing stage must be preserved")
	}
	if observation.PlanSignal != jevExecutionPlanUnknown {
		t.Fatalf("signal = %q, want unknown", observation.PlanSignal)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("observation should validate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVExecutionPlanObservationRejectsTamperedEvidence(t *testing.T) {
	input := RevisionSelfImprovementCycleJEVExecutionPlanEvidence{
		Status:           "BOUND",
		ContractName:     jevExecutionPlanContractName,
		NonExecuting:     true,
		NonAuthorizing:   true,
	}
	observation := ObserveRevisionSelfImprovementCycleJEVExecutionPlanObservation(input)

	if observation.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN for tampered evidence", observation.Status)
	}
	if observation.MissingStage == "" {
		t.Fatal("tampered evidence must preserve an unresolved stage")
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("observation should validate after preserving UNKNOWN: %v", err)
	}
}
