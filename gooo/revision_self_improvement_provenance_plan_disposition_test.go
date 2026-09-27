package gooo

import "testing"

func provenancePlanDispositionInputs(
	t *testing.T,
) (RevisionSelfImprovementIterationProvenanceObservation, RevisionSelfImprovementExecutionPlanObservation) {
	t.Helper()
	iteration, history, bridge := iterationProvenanceInputs(t, false)
	provenance, err := ObserveRevisionSelfImprovementIterationProvenance(iteration, history, bridge)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementIterationProvenance() error = %v", err)
	}
	decision, planObservation, plan := selfImprovementExecutionPlanInputs(t)
	executionPlan, err := ObserveRevisionSelfImprovementExecutionPlan(decision, planObservation, plan)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementExecutionPlan() error = %v", err)
	}
	return provenance, executionPlan
}

func TestObserveRevisionSelfImprovementProvenancePlanDispositionBindsStable(t *testing.T) {
	iteration, plan := provenancePlanDispositionInputs(t)
	result, err := ObserveRevisionSelfImprovementProvenancePlanDisposition(iteration, plan)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenancePlanDisposition() error = %v", err)
	}
	if result.Status != "BOUND" ||
		result.DispositionSignal != "provenance-observe-plan" ||
		result.DecisionSignal != "observe" ||
		!result.PlanObserved {
		t.Fatalf("unexpected provenance plan disposition: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementProvenancePlanDispositionRetainsIterationFailure(t *testing.T) {
	iteration, plan := provenancePlanDispositionInputs(t)
	iteration.ObservationDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementProvenancePlanDisposition(iteration, plan)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementProvenancePlanDisposition() error = nil, want iteration failure")
	}
	if result.Status != "UNKNOWN" ||
		result.MissingStage != "revision-self-improvement-provenance-plan-disposition-iteration" {
		t.Fatalf("unexpected unknown iteration disposition: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementProvenancePlanDispositionRetainsSourceLinkFailure(t *testing.T) {
	iteration, plan := provenancePlanDispositionInputs(t)
	plan.SourceDigest = digestString("other-source")
	plan.ObservationDigest = digestRevisionSelfImprovementExecutionPlan(plan)
	result, err := ObserveRevisionSelfImprovementProvenancePlanDisposition(iteration, plan)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementProvenancePlanDisposition() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" ||
		result.MissingStage != "revision-self-improvement-provenance-plan-disposition-source-link" {
		t.Fatalf("unexpected unknown source-link disposition: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementProvenancePlanDispositionIsDeterministic(t *testing.T) {
	iteration, plan := provenancePlanDispositionInputs(t)
	first, err := ObserveRevisionSelfImprovementProvenancePlanDisposition(iteration, plan)
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementProvenancePlanDisposition() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementProvenancePlanDisposition(iteration, plan)
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementProvenancePlanDisposition() error = %v", err)
	}
	if first.ObservationDigest != second.ObservationDigest {
		t.Fatal("same provenance and plan produced different disposition digest")
	}
}
