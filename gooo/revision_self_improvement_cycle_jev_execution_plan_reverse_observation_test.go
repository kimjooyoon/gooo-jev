package gooo

import "testing"

func TestRevisionSelfImprovementCycleJEVExecutionPlanReverseObservationPreservesUnknown(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVExecutionPlanReverseObservation(
		RevisionSelfImprovementCycleJEVExecutionPlanObservation{},
	)

	if observation.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN", observation.Status)
	}
	if observation.MissingStage == "" {
		t.Fatal("missing stage must be preserved")
	}
	if observation.PlanSignal != jevExecutionPlanReverseUnknown {
		t.Fatalf("signal = %q, want unknown", observation.PlanSignal)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("observation should validate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVExecutionPlanReverseObservationRejectsTamperedObservation(t *testing.T) {
	input := RevisionSelfImprovementCycleJEVExecutionPlanObservation{
		Status:       "BOUND",
		MetricName:   jevExecutionPlanMetricName,
		ContractName: jevExecutionPlanContractName,
		PlanSignal:   jevExecutionPlanComplete,
		NonExecuting: true,
		NonAuthorizing: true,
	}
	observation := ObserveRevisionSelfImprovementCycleJEVExecutionPlanReverseObservation(input)

	if observation.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN for tampered observation", observation.Status)
	}
	if observation.MissingStage == "" {
		t.Fatal("tampered observation must preserve an unresolved stage")
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("observation should validate after preserving UNKNOWN: %v", err)
	}
}
