package gooo

import "testing"

func selfImprovementLifecycleInputs(t *testing.T) (RevisionSelfImprovementExecutionApplicationObservation, RevisionSelfImprovementOutcomeObservation) {
	t.Helper()
	decision, planObservation, plan := selfImprovementExecutionPlanInputs(t)
	executionPlan, err := ObserveRevisionSelfImprovementExecutionPlan(decision, planObservation, plan)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementExecutionPlan() error = %v", err)
	}
	iteration, candidatePlanObservation, applicationObservation := selfImprovementApplicationInputs(t)
	application, err := ObserveRevisionSelfImprovementApplication(iteration, candidatePlanObservation, applicationObservation)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementApplication() error = %v", err)
	}
	executionApplication, err := ObserveRevisionSelfImprovementExecutionApplication(executionPlan, application)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementExecutionApplication() error = %v", err)
	}
	_, metricsBinding, _, generationAssessment := selfImprovementReceiptInputs(t)
	outcome, err := ObserveRevisionSelfImprovementOutcome(iteration, application, metricsBinding, generationAssessment)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementOutcome() error = %v", err)
	}
	return executionApplication, outcome
}

func TestObserveRevisionSelfImprovementLifecycleBindsOutcome(t *testing.T) {
	executionApplication, outcome := selfImprovementLifecycleInputs(t)
	result, err := ObserveRevisionSelfImprovementLifecycle(executionApplication, outcome)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementLifecycle() error = %v", err)
	}
	if result.Status != "BOUND" || result.LifecycleSignal != "plan-application-outcome-bound" ||
		result.MetricsBindingDigest == "" ||
		result.GeneratedIRDigest != outcome.GeneratedIRDigest ||
		result.PlanDigest != executionApplication.PlanDigest {
		t.Fatalf("unexpected self-improvement lifecycle: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementLifecycleRetainsExecutionApplicationFailure(t *testing.T) {
	executionApplication, outcome := selfImprovementLifecycleInputs(t)
	executionApplication.ObservationDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementLifecycle(executionApplication, outcome)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementLifecycle() error = nil, want execution application failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-lifecycle-execution-application" {
		t.Fatalf("unexpected unknown execution application lifecycle: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementLifecycleRetainsOutcomeFailure(t *testing.T) {
	executionApplication, outcome := selfImprovementLifecycleInputs(t)
	outcome.OutcomeDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementLifecycle(executionApplication, outcome)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementLifecycle() error = nil, want outcome failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-lifecycle-outcome" {
		t.Fatalf("unexpected unknown outcome lifecycle: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementLifecycleRetainsSourceLinkFailure(t *testing.T) {
	executionApplication, outcome := selfImprovementLifecycleInputs(t)
	outcome.SourceDigest = digestString("other-source")
	outcome.OutcomeDigest = digestRevisionSelfImprovementOutcomeObservation(outcome)
	result, err := ObserveRevisionSelfImprovementLifecycle(executionApplication, outcome)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementLifecycle() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-lifecycle-source-link" {
		t.Fatalf("unexpected unknown source lifecycle: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementLifecycleRetainsApplicationLinkFailure(t *testing.T) {
	executionApplication, outcome := selfImprovementLifecycleInputs(t)
	outcome.PlanDigest = digestString("other-plan")
	outcome.OutcomeDigest = digestRevisionSelfImprovementOutcomeObservation(outcome)
	result, err := ObserveRevisionSelfImprovementLifecycle(executionApplication, outcome)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementLifecycle() error = nil, want application link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-lifecycle-application-link" {
		t.Fatalf("unexpected unknown application lifecycle: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementLifecycleIsDeterministic(t *testing.T) {
	executionApplication, outcome := selfImprovementLifecycleInputs(t)
	first, err := ObserveRevisionSelfImprovementLifecycle(executionApplication, outcome)
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementLifecycle() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementLifecycle(executionApplication, outcome)
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementLifecycle() error = %v", err)
	}
	if first.ObservationDigest != second.ObservationDigest {
		t.Fatal("same plan, application, and outcome produced different lifecycle digest")
	}
}
