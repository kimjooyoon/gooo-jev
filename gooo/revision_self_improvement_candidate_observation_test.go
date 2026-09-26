package gooo

import "testing"

func selfImprovementCandidateInputs(t *testing.T) (RevisionSelfImprovementIteration, RevisionCandidatePlanObservation) {
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
	return iteration, planObservation
}

func TestObserveRevisionSelfImprovementCandidateBindsPlan(t *testing.T) {
	iteration, planObservation := selfImprovementCandidateInputs(t)
	result, err := ObserveRevisionSelfImprovementCandidate(iteration, planObservation)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementCandidate() error = %v", err)
	}
	if result.Status != "BOUND" || !result.CandidatePlanObserved ||
		result.CandidateDigest != planObservation.CandidateDigest ||
		result.PlanDigest != planObservation.PlanDigest {
		t.Fatalf("unexpected self-improvement candidate observation: %#v", result)
	}
	if result.IterationDigest != iteration.IterationDigest ||
		result.FeedbackSignal != iteration.FeedbackSignal {
		t.Fatalf("candidate observation lost iteration links: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementCandidateRetainsIterationFailure(t *testing.T) {
	iteration, planObservation := selfImprovementCandidateInputs(t)
	iteration.IterationDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementCandidate(iteration, planObservation)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementCandidate() error = nil, want iteration failure")
	}
	if result.Status != "UNKNOWN" ||
		result.MissingStage != "revision-self-improvement-candidate-iteration" {
		t.Fatalf("unexpected unknown iteration candidate: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementCandidateRetainsPlanFailure(t *testing.T) {
	iteration, planObservation := selfImprovementCandidateInputs(t)
	planObservation.ObservationDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementCandidate(iteration, planObservation)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementCandidate() error = nil, want plan failure")
	}
	if result.Status != "UNKNOWN" ||
		result.MissingStage != "revision-self-improvement-candidate-plan" {
		t.Fatalf("unexpected unknown plan candidate: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementCandidateRetainsSourceLinkFailure(t *testing.T) {
	iteration, planObservation := selfImprovementCandidateInputs(t)
	planObservation.SourceDigest = digestString("other-source")
	planObservation.ObservationDigest = digestRevisionCandidatePlanObservation(planObservation)
	result, err := ObserveRevisionSelfImprovementCandidate(iteration, planObservation)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementCandidate() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" ||
		result.MissingStage != "revision-self-improvement-candidate-source-link" {
		t.Fatalf("unexpected unknown source-link candidate: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementCandidateIsDeterministic(t *testing.T) {
	iteration, planObservation := selfImprovementCandidateInputs(t)
	first, err := ObserveRevisionSelfImprovementCandidate(iteration, planObservation)
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementCandidate() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCandidate(iteration, planObservation)
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementCandidate() error = %v", err)
	}
	if first.ObservationDigest != second.ObservationDigest {
		t.Fatal("same iteration and candidate plan produced different observation digest")
	}
}