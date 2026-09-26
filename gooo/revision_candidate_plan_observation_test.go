package gooo

import "testing"

func candidatePlanObservationInputs(t *testing.T) (RevisionCandidateObservation, RevisionApplicationPlan) {
	t.Helper()
	baseline := summaryForReplacement(t, "lineage")
	candidate := summaryForReplacement(t, "lineage-expanded")
	comparison, err := ObserveRevisionImprovement(baseline, candidate)
	if err != nil {
		t.Fatalf("ObserveRevisionImprovement() error = %v", err)
	}
	_, materialization := selectedCandidateInputs(t)
	candidateObservation, err := ObserveRevisionCandidateForImprovement(comparison, materialization)
	if err != nil {
		t.Fatalf("ObserveRevisionCandidateForImprovement() error = %v", err)
	}
	binding, edit := boundApplicationPlanInputs(t)
	plan, err := PlanRevisionApplication(validContract, binding, edit)
	if err != nil {
		t.Fatalf("PlanRevisionApplication() error = %v", err)
	}
	if candidateObservation.CandidateDigest != plan.CandidateDigest {
		t.Fatalf("fixture candidate mismatch: observation=%s plan=%s", candidateObservation.CandidateDigest, plan.CandidateDigest)
	}
	return candidateObservation, plan
}

func TestObserveRevisionCandidatePlanBindsPlanEvidence(t *testing.T) {
	observation, plan := candidatePlanObservationInputs(t)
	result, err := ObserveRevisionCandidatePlan(observation, plan)
	if err != nil {
		t.Fatalf("ObserveRevisionCandidatePlan() error = %v", err)
	}
	if result.Status != "BOUND" || !result.PlanObserved || result.PlanSignal == "" {
		t.Fatalf("unexpected revision candidate plan observation: %#v", result)
	}
	if result.CandidateDigest != plan.CandidateDigest || result.PlanDigest != plan.PlanDigest || result.EditDigest != plan.EditDigest {
		t.Fatalf("plan observation lost plan links: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionCandidatePlanRetainsCandidateFailure(t *testing.T) {
	observation, plan := candidatePlanObservationInputs(t)
	observation.ResultDigest = digestString("tampered")
	result, err := ObserveRevisionCandidatePlan(observation, plan)
	if err == nil {
		t.Fatal("ObserveRevisionCandidatePlan() error = nil, want candidate failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-candidate-plan-observation-candidate" {
		t.Fatalf("unexpected unknown candidate plan observation: %#v", result)
	}
}

func TestObserveRevisionCandidatePlanRetainsPlanFailure(t *testing.T) {
	observation, plan := candidatePlanObservationInputs(t)
	plan.PlanDigest = digestString("tampered")
	result, err := ObserveRevisionCandidatePlan(observation, plan)
	if err == nil {
		t.Fatal("ObserveRevisionCandidatePlan() error = nil, want plan failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-candidate-plan-observation-plan" {
		t.Fatalf("unexpected unknown plan observation: %#v", result)
	}
}

func TestObserveRevisionCandidatePlanRetainsLinkFailure(t *testing.T) {
	observation, plan := candidatePlanObservationInputs(t)
	observation.CandidateDigest = digestString("other-candidate")
	observation.ResultDigest = digestRevisionCandidateObservation(observation)
	result, err := ObserveRevisionCandidatePlan(observation, plan)
	if err == nil {
		t.Fatal("ObserveRevisionCandidatePlan() error = nil, want link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-candidate-plan-observation-link" {
		t.Fatalf("unexpected unknown plan link observation: %#v", result)
	}
}

func TestObserveRevisionCandidatePlanIsDeterministic(t *testing.T) {
	observation, plan := candidatePlanObservationInputs(t)
	first, err := ObserveRevisionCandidatePlan(observation, plan)
	if err != nil {
		t.Fatalf("first ObserveRevisionCandidatePlan() error = %v", err)
	}
	second, err := ObserveRevisionCandidatePlan(observation, plan)
	if err != nil {
		t.Fatalf("second ObserveRevisionCandidatePlan() error = %v", err)
	}
	if first.ObservationDigest != second.ObservationDigest {
		t.Fatal("same candidate observation and plan produced different digest")
	}
}