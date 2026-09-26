package gooo

import "testing"

func revisionApplicationObservationInputs(t *testing.T) (RevisionCandidatePlanObservation, RevisionApplicationReceipt) {
	t.Helper()
	candidateObservation, plan := candidatePlanObservationInputs(t)
	planObservation, err := ObserveRevisionCandidatePlan(candidateObservation, plan)
	if err != nil {
		t.Fatalf("ObserveRevisionCandidatePlan() error = %v", err)
	}
	application, err := ApplyRevision(validContract, plan.SourceDigest, plan.Candidate, plan.Edit)
	if err != nil {
		t.Fatalf("ApplyRevision() error = %v", err)
	}
	receipt, err := ObserveRevisionApplicationReceipt(plan, application)
	if err != nil {
		t.Fatalf("ObserveRevisionApplicationReceipt() error = %v", err)
	}
	return planObservation, receipt
}

func TestObserveRevisionApplicationBindsReceipt(t *testing.T) {
	planObservation, receipt := revisionApplicationObservationInputs(t)
	result, err := ObserveRevisionApplication(planObservation, receipt)
	if err != nil {
		t.Fatalf("ObserveRevisionApplication() error = %v", err)
	}
	if result.Status != "BOUND" || !result.ApplicationObserved || result.ApplicationSignal != "application-observed" {
		t.Fatalf("unexpected revision application observation: %#v", result)
	}
	if result.PlanDigest != receipt.PlanDigest ||
		result.ApplicationDigest != receipt.ApplicationDigest ||
		result.ProposedIRDigest != receipt.ProposedIRDigest {
		t.Fatalf("application observation lost receipt links: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionApplicationRetainsPlanFailure(t *testing.T) {
	planObservation, receipt := revisionApplicationObservationInputs(t)
	planObservation.PlanDigest = digestString("other-plan")
	planObservation.ObservationDigest = digestRevisionCandidatePlanObservation(planObservation)
	result, err := ObserveRevisionApplication(planObservation, receipt)
	if err == nil {
		t.Fatal("ObserveRevisionApplication() error = nil, want plan link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-application-observation-link" {
		t.Fatalf("unexpected unknown plan observation: %#v", result)
	}
}

func TestObserveRevisionApplicationRetainsReceiptFailure(t *testing.T) {
	planObservation, receipt := revisionApplicationObservationInputs(t)
	receipt.ReceiptDigest = digestString("tampered")
	result, err := ObserveRevisionApplication(planObservation, receipt)
	if err == nil {
		t.Fatal("ObserveRevisionApplication() error = nil, want receipt failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-application-observation-receipt" {
		t.Fatalf("unexpected unknown receipt observation: %#v", result)
	}
}

func TestObserveRevisionApplicationIsDeterministic(t *testing.T) {
	planObservation, receipt := revisionApplicationObservationInputs(t)
	first, err := ObserveRevisionApplication(planObservation, receipt)
	if err != nil {
		t.Fatalf("first ObserveRevisionApplication() error = %v", err)
	}
	second, err := ObserveRevisionApplication(planObservation, receipt)
	if err != nil {
		t.Fatalf("second ObserveRevisionApplication() error = %v", err)
	}
	if first.ObservationDigest != second.ObservationDigest {
		t.Fatal("same plan observation and receipt produced different digest")
	}
}