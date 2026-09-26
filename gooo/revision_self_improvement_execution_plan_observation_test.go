package gooo

import "testing"

func selfImprovementExecutionPlanInputs(t *testing.T) (RevisionSelfImprovementDecisionObservation, RevisionCandidatePlanObservation, RevisionApplicationPlan) {
	t.Helper()
	decisionIteration, outcome, feedback := selfImprovementDecisionInputs(t)
	decision, err := ObserveRevisionSelfImprovementDecision(decisionIteration, outcome, feedback)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementDecision() error = %v", err)
	}
	candidateObservation, plan := candidatePlanObservationInputs(t)
	planObservation, err := ObserveRevisionCandidatePlan(candidateObservation, plan)
	if err != nil {
		t.Fatalf("ObserveRevisionCandidatePlan() error = %v", err)
	}
	return decision, planObservation, plan
}

func TestObserveRevisionSelfImprovementExecutionPlanBindsPlan(t *testing.T) {
	decision, planObservation, plan := selfImprovementExecutionPlanInputs(t)
	result, err := ObserveRevisionSelfImprovementExecutionPlan(decision, planObservation, plan)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementExecutionPlan() error = %v", err)
	}
	if result.Status != "BOUND" || !result.PlanObserved ||
		result.DecisionSignal != decision.DecisionSignal ||
		result.PlanDigest != plan.PlanDigest {
		t.Fatalf("unexpected self-improvement execution plan: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementExecutionPlanRetainsDecisionFailure(t *testing.T) {
	decision, planObservation, plan := selfImprovementExecutionPlanInputs(t)
	decision.DecisionDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementExecutionPlan(decision, planObservation, plan)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementExecutionPlan() error = nil, want decision failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-execution-plan-decision" {
		t.Fatalf("unexpected unknown decision plan: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementExecutionPlanRetainsObservationFailure(t *testing.T) {
	decision, planObservation, plan := selfImprovementExecutionPlanInputs(t)
	planObservation.ObservationDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementExecutionPlan(decision, planObservation, plan)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementExecutionPlan() error = nil, want observation failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-execution-plan-observation" {
		t.Fatalf("unexpected unknown plan observation: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementExecutionPlanRetainsSourceLinkFailure(t *testing.T) {
	decision, planObservation, plan := selfImprovementExecutionPlanInputs(t)
	decision.SourceDigest = digestString("other-source")
	decision.DecisionDigest = digestRevisionSelfImprovementDecision(decision)
	result, err := ObserveRevisionSelfImprovementExecutionPlan(decision, planObservation, plan)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementExecutionPlan() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-execution-plan-source-link" {
		t.Fatalf("unexpected unknown source-link plan: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementExecutionPlanRetainsPlanLinkFailure(t *testing.T) {
	decision, planObservation, plan := selfImprovementExecutionPlanInputs(t)
	planObservation.CandidateDigest = digestString("other-candidate")
	planObservation.ObservationDigest = digestRevisionCandidatePlanObservation(planObservation)
	result, err := ObserveRevisionSelfImprovementExecutionPlan(decision, planObservation, plan)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementExecutionPlan() error = nil, want plan link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-execution-plan-plan-link" {
		t.Fatalf("unexpected unknown plan-link plan: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementExecutionPlanIsDeterministic(t *testing.T) {
	decision, planObservation, plan := selfImprovementExecutionPlanInputs(t)
	first, err := ObserveRevisionSelfImprovementExecutionPlan(decision, planObservation, plan)
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementExecutionPlan() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementExecutionPlan(decision, planObservation, plan)
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementExecutionPlan() error = %v", err)
	}
	if first.ObservationDigest != second.ObservationDigest {
		t.Fatal("same decision and plan produced different observation digest")
	}
}
