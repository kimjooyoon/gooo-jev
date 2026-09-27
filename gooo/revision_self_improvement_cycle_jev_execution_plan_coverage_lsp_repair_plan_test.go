package gooo

import "testing"

func TestRevisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlanPreservesUnknown(t *testing.T) {
	plan := ObserveRevisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlan(
		RevisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricLSPProjection{},
	)

	if plan.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN", plan.Status)
	}
	if plan.MissingStage == "" || plan.TargetStage != plan.MissingStage {
		t.Fatalf("target stage = %q/%q, want preserved missing stage", plan.TargetStage, plan.MissingStage)
	}
	if plan.RepairAction != "resolve_missing_provenance" {
		t.Fatalf("repair action = %q, want resolve_missing_provenance", plan.RepairAction)
	}
	if err := plan.Validate(); err != nil {
		t.Fatalf("repair plan should validate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlanRejectsTamperedProjection(t *testing.T) {
	input := RevisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricLSPProjection{
		Status:         "BOUND",
		MetricName:     jevExecutionPlanCoverageMetricName,
		CoverageMilli:  1000,
		CoverageBand:   "high",
		CoverageSignal: jevExecutionPlanCoverageLSPComplete,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	plan := ObserveRevisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlan(input)

	if plan.Status != "UNKNOWN" {
		t.Fatalf("status = %q, want UNKNOWN for tampered projection", plan.Status)
	}
	if plan.MissingStage == "" {
		t.Fatal("tampered projection must preserve an unresolved stage")
	}
	if err := plan.Validate(); err != nil {
		t.Fatalf("repair plan should validate after preserving UNKNOWN: %v", err)
	}
}
