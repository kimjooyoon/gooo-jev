package gooo

import "testing"

func provenanceReplanInputs(t *testing.T) (RevisionSelfImprovementProvenanceReverseDecisionObservation, RevisionSelfImprovementExecutionPlanObservation) {
	t.Helper()
	feedback, reverse, decision := provenanceReverseDecisionInputs(t)
	provenanceDecision, err := ObserveRevisionSelfImprovementProvenanceReverseDecision(feedback, reverse, decision)
	if err != nil { t.Fatalf("ObserveRevisionSelfImprovementProvenanceReverseDecision() error = %v", err) }
	executionDecision, planObservation, plan := selfImprovementExecutionPlanInputs(t)
	executionPlan, err := ObserveRevisionSelfImprovementExecutionPlan(executionDecision, planObservation, plan)
	if err != nil { t.Fatalf("ObserveRevisionSelfImprovementExecutionPlan() error = %v", err) }
	return provenanceDecision, executionPlan
}

func TestObserveRevisionSelfImprovementProvenanceReplanBindsNextPlan(t *testing.T) {
	decision, plan := provenanceReplanInputs(t)
	result, err := ObserveRevisionSelfImprovementProvenanceReplan(decision, plan)
	if err != nil { t.Fatalf("ObserveRevisionSelfImprovementProvenanceReplan() error = %v", err) }
	if result.Status != "BOUND" || result.ReplanSignal != "provenance-replan-mismatch" || result.NextPlanDigest != plan.PlanDigest || result.PlanningSignal != plan.PlanningSignal { t.Fatalf("unexpected provenance replan: %#v", result) }
	if result.ProvenanceDecisionDigest != decision.ObservationDigest || result.ExecutionPlanObservationDigest != plan.ObservationDigest { t.Fatalf("provenance replan lost links: %#v", result) }
	if err := result.Validate(); err != nil { t.Fatalf("Validate() error = %v", err) }
}

func TestObserveRevisionSelfImprovementProvenanceReplanRetainsDecisionFailure(t *testing.T) {
	decision, plan := provenanceReplanInputs(t)
	plan.DecisionDigest = digestString("other-decision")
	plan.ObservationDigest = digestRevisionSelfImprovementExecutionPlan(plan)
	result, err := ObserveRevisionSelfImprovementProvenanceReplan(decision, plan)
	if err == nil { t.Fatal("ObserveRevisionSelfImprovementProvenanceReplan() error = nil, want decision link failure") }
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-provenance-replan-decision-link" { t.Fatalf("unexpected unknown decision replan: %#v", result) }
}

func TestObserveRevisionSelfImprovementProvenanceReplanRetainsRequirementFailure(t *testing.T) {
	decision, plan := provenanceReplanInputs(t)
	plan.RequiresInspection = !plan.RequiresInspection
	plan.ObservationDigest = digestRevisionSelfImprovementExecutionPlan(plan)
	result, err := ObserveRevisionSelfImprovementProvenanceReplan(decision, plan)
	if err == nil { t.Fatal("ObserveRevisionSelfImprovementProvenanceReplan() error = nil, want requirement link failure") }
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-provenance-replan-requirement-link" { t.Fatalf("unexpected unknown requirement replan: %#v", result) }
}

func TestObserveRevisionSelfImprovementProvenanceReplanIsDeterministic(t *testing.T) {
	decision, plan := provenanceReplanInputs(t)
	first, err := ObserveRevisionSelfImprovementProvenanceReplan(decision, plan)
	if err != nil { t.Fatalf("first ObserveRevisionSelfImprovementProvenanceReplan() error = %v", err) }
	second, err := ObserveRevisionSelfImprovementProvenanceReplan(decision, plan)
	if err != nil { t.Fatalf("second ObserveRevisionSelfImprovementProvenanceReplan() error = %v", err) }
	if first.ObservationDigest != second.ObservationDigest { t.Fatal("same provenance decision and execution plan produced different replan digest") }
}
