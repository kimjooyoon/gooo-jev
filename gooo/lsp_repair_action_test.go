package gooo

import "testing"

func lspRepairActionSnapshot(status, missingStage string) LanguageSnapshot {
	sourceDigest := digestString("lsp-repair-action-source")
	snapshot := LanguageSnapshot{
		Status:         status,
		MissingStage:   missingStage,
		SourceDigest:   sourceDigest,
		IRDigest:       digestString("lsp-repair-action-ir"),
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if status == "UNKNOWN" {
		snapshot.Diagnostics = []Diagnostic{{
			Stage:        missingStage,
			Message:      "provenance is incomplete",
			Position:     Position{Line: 1, Column: 1},
			Severity:     SeverityError,
			SourceDigest: sourceDigest,
		}}
	}
	return snapshot
}

func lspRepairActionPlan(status, missingStage string) RevisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlan {
	plan := RevisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlan{
		Status:         status,
		MissingStage:   missingStage,
		PlanName:       jevExecutionPlanCoverageRepairPlanName,
		TargetStage:    missingStage,
		RepairAction:   "resolve_missing_provenance",
		PlanSignal:     jevExecutionPlanCoverageRepairPlanUnknown,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if status == "BOUND" {
		plan.MissingStage = ""
		plan.TargetStage = ""
		plan.RepairAction = "none"
		plan.PlanSignal = jevExecutionPlanCoverageRepairPlanComplete
		plan.EvidencePrefixDigest = digestString("lsp-repair-action-prefix")
		plan.ProjectionDigest = digestString("lsp-repair-action-projection")
	}
	plan.ObservationDigest = revisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlanDigest(plan)
	return plan
}

func TestObserveLSPRepairActionProjectionDeferred(t *testing.T) {
	projection := ObserveLSPRepairActionProjection(
		lspRepairActionSnapshot("UNKNOWN", "reverse_observation"),
		lspRepairActionPlan("UNKNOWN", "reverse_observation"),
	)
	if projection.Status != "BOUND" ||
		projection.ActionSignal != lspRepairActionDeferred ||
		projection.TargetStage != "reverse_observation" {
		t.Fatalf("expected deferred repair action, got %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("expected valid deferred repair action: %v", err)
	}
}

func TestObserveLSPRepairActionProjectionComplete(t *testing.T) {
	projection := ObserveLSPRepairActionProjection(
		lspRepairActionSnapshot("BOUND", ""),
		lspRepairActionPlan("BOUND", ""),
	)
	if projection.Status != "BOUND" || projection.ActionSignal != lspRepairActionComplete {
		t.Fatalf("expected complete repair action, got %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("expected valid complete repair action: %v", err)
	}
}

func TestObserveLSPRepairActionProjectionPreservesUnknownPlan(t *testing.T) {
	plan := lspRepairActionPlan("UNKNOWN", "reverse_observation")
	plan.ObservationDigest = ""
	projection := ObserveLSPRepairActionProjection(
		lspRepairActionSnapshot("UNKNOWN", "reverse_observation"),
		plan,
	)
	if projection.Status != "UNKNOWN" || projection.MissingStage != "lsp-repair-action-plan" {
		t.Fatalf("expected invalid plan to remain unknown, got %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("expected valid unknown repair action: %v", err)
	}
}

func TestLSPRepairActionProjectionRejectsTampering(t *testing.T) {
	projection := ObserveLSPRepairActionProjection(
		lspRepairActionSnapshot("UNKNOWN", "reverse_observation"),
		lspRepairActionPlan("UNKNOWN", "reverse_observation"),
	)
	projection.Title = "execute command"
	if err := projection.Validate(); err == nil {
		t.Fatal("expected title tampering to be rejected")
	}
}