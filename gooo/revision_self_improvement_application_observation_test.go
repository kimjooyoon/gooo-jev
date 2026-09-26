package gooo

import "testing"

func selfImprovementApplicationInputs(t *testing.T) (RevisionSelfImprovementIteration, RevisionCandidatePlanObservation, RevisionApplicationObservation) {
	t.Helper()
	source, history, feedback := selfImprovementIterationInputs(t)
	iteration, err := ObserveRevisionSelfImprovementIteration(source, history, feedback)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementIteration() error = %v", err)
	}
	candidateObservation, plan := candidatePlanObservationInputs(t)
	planObservation, err := ObserveRevisionCandidatePlan(candidateObservation, plan)
	if err != nil {
		t.Fatalf("ObserveRevisionCandidatePlan() error = %v", err)
	}
	applicationObservation, _, _, _ := selfImprovementReceiptInputs(t)
	return iteration, planObservation, applicationObservation
}

func TestObserveRevisionSelfImprovementApplicationBindsReceipt(t *testing.T) {
	iteration, planObservation, applicationObservation := selfImprovementApplicationInputs(t)
	result, err := ObserveRevisionSelfImprovementApplication(iteration, planObservation, applicationObservation)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementApplication() error = %v", err)
	}
	if result.Status != "BOUND" || !result.ApplicationObserved ||
		result.ApplicationDigest != applicationObservation.ApplicationDigest ||
		result.CandidateDigest != applicationObservation.CandidateDigest {
		t.Fatalf("unexpected self-improvement application observation: %#v", result)
	}
	if result.IterationDigest != iteration.IterationDigest ||
		result.PlanObservationDigest != planObservation.ObservationDigest ||
		result.FeedbackSignal != iteration.FeedbackSignal {
		t.Fatalf("application observation lost iteration links: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementApplicationRetainsIterationFailure(t *testing.T) {
	iteration, planObservation, applicationObservation := selfImprovementApplicationInputs(t)
	iteration.IterationDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementApplication(iteration, planObservation, applicationObservation)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementApplication() error = nil, want iteration failure")
	}
	if result.Status != "UNKNOWN" ||
		result.MissingStage != "revision-self-improvement-application-iteration" {
		t.Fatalf("unexpected unknown iteration application: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementApplicationRetainsPlanLinkFailure(t *testing.T) {
	iteration, planObservation, applicationObservation := selfImprovementApplicationInputs(t)
	applicationObservation.PlanObservationDigest = digestString("other-plan-observation")
	applicationObservation.ObservationDigest = digestRevisionApplicationObservation(applicationObservation)
	result, err := ObserveRevisionSelfImprovementApplication(iteration, planObservation, applicationObservation)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementApplication() error = nil, want plan link failure")
	}
	if result.Status != "UNKNOWN" ||
		result.MissingStage != "revision-self-improvement-application-plan-link" {
		t.Fatalf("unexpected unknown plan-link application: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementApplicationRetainsSourceLinkFailure(t *testing.T) {
	iteration, planObservation, applicationObservation := selfImprovementApplicationInputs(t)
	applicationObservation.SourceDigest = digestString("other-source")
	applicationObservation.ObservationDigest = digestRevisionApplicationObservation(applicationObservation)
	result, err := ObserveRevisionSelfImprovementApplication(iteration, planObservation, applicationObservation)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementApplication() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" ||
		result.MissingStage != "revision-self-improvement-application-source-link" {
		t.Fatalf("unexpected unknown source-link application: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementApplicationIsDeterministic(t *testing.T) {
	iteration, planObservation, applicationObservation := selfImprovementApplicationInputs(t)
	first, err := ObserveRevisionSelfImprovementApplication(iteration, planObservation, applicationObservation)
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementApplication() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementApplication(iteration, planObservation, applicationObservation)
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementApplication() error = %v", err)
	}
	if first.ObservationDigest != second.ObservationDigest {
		t.Fatal("same iteration and application produced different observation digest")
	}
}