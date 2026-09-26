package gooo

import "testing"

func selfImprovementExecutionApplicationInputs(t *testing.T) (RevisionSelfImprovementExecutionPlanObservation, RevisionSelfImprovementApplicationObservation) {
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
	return executionPlan, application
}

func TestObserveRevisionSelfImprovementExecutionApplicationBindsApplication(t *testing.T) {
	plan, application := selfImprovementExecutionApplicationInputs(t)
	result, err := ObserveRevisionSelfImprovementExecutionApplication(plan, application)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementExecutionApplication() error = %v", err)
	}
	if result.Status != "BOUND" || !result.PlanObserved || !result.ApplicationObserved ||
		result.PlanDigest != application.PlanDigest ||
		result.ApplicationDigest != application.ApplicationDigest {
		t.Fatalf("unexpected execution application observation: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementExecutionApplicationRetainsPlanFailure(t *testing.T) {
	plan, application := selfImprovementExecutionApplicationInputs(t)
	plan.ObservationDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementExecutionApplication(plan, application)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementExecutionApplication() error = nil, want plan failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-execution-application-plan" {
		t.Fatalf("unexpected unknown plan application: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementExecutionApplicationRetainsApplicationFailure(t *testing.T) {
	plan, application := selfImprovementExecutionApplicationInputs(t)
	application.ObservationDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementExecutionApplication(plan, application)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementExecutionApplication() error = nil, want application failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-execution-application-observation" {
		t.Fatalf("unexpected unknown application: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementExecutionApplicationRetainsSourceLinkFailure(t *testing.T) {
	plan, application := selfImprovementExecutionApplicationInputs(t)
	application.SourceDigest = digestString("other-source")
	application.ObservationDigest = digestRevisionSelfImprovementApplicationObservation(application)
	result, err := ObserveRevisionSelfImprovementExecutionApplication(plan, application)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementExecutionApplication() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-execution-application-source-link" {
		t.Fatalf("unexpected unknown source link: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementExecutionApplicationRetainsPlanLinkFailure(t *testing.T) {
	plan, application := selfImprovementExecutionApplicationInputs(t)
	application.PlanDigest = digestString("other-plan")
	application.ObservationDigest = digestRevisionSelfImprovementApplicationObservation(application)
	result, err := ObserveRevisionSelfImprovementExecutionApplication(plan, application)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementExecutionApplication() error = nil, want plan link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-execution-application-plan-link" {
		t.Fatalf("unexpected unknown plan link: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementExecutionApplicationRetainsSignalLinkFailure(t *testing.T) {
	plan, application := selfImprovementExecutionApplicationInputs(t)
	application.FeedbackSignal = "observe"
	application.ObservationDigest = digestRevisionSelfImprovementApplicationObservation(application)
	result, err := ObserveRevisionSelfImprovementExecutionApplication(plan, application)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementExecutionApplication() error = nil, want signal link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-execution-application-signal-link" {
		t.Fatalf("unexpected unknown signal link: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementExecutionApplicationIsDeterministic(t *testing.T) {
	plan, application := selfImprovementExecutionApplicationInputs(t)
	first, err := ObserveRevisionSelfImprovementExecutionApplication(plan, application)
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementExecutionApplication() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementExecutionApplication(plan, application)
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementExecutionApplication() error = %v", err)
	}
	if first.ObservationDigest != second.ObservationDigest {
		t.Fatal("same execution plan and application produced different digest")
	}
}
