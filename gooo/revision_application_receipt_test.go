package gooo

import "testing"

func appliedRevisionReceiptInputs(t *testing.T) (RevisionApplicationPlan, RevisionApplication) {
	t.Helper()
	binding, edit := boundApplicationPlanInputs(t)
	plan, err := PlanRevisionApplication(validContract, binding, edit)
	if err != nil {
		t.Fatalf("PlanRevisionApplication() error = %v", err)
	}
	application, err := ApplyRevision(validContract, plan.SourceDigest, plan.Candidate, plan.Edit)
	if err != nil {
		t.Fatalf("ApplyRevision() error = %v", err)
	}
	return plan, application
}

func TestObserveRevisionApplicationReceiptBindsPlanAndApplication(t *testing.T) {
	plan, application := appliedRevisionReceiptInputs(t)
	receipt, err := ObserveRevisionApplicationReceipt(plan, application)
	if err != nil {
		t.Fatalf("ObserveRevisionApplicationReceipt() error = %v", err)
	}
	if receipt.Status != "BOUND" ||
		receipt.SourceDigest != plan.SourceDigest ||
		receipt.InputIRDigest != plan.InputIRDigest ||
		receipt.CandidateDigest != plan.CandidateDigest ||
		receipt.EditDigest != plan.EditDigest ||
		receipt.ProposedIRDigest != application.ProposedIRDigest {
		t.Fatalf("unexpected application receipt: %#v", receipt)
	}
	if err := receipt.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionApplicationReceiptRetainsApplicationFailure(t *testing.T) {
	plan, application := appliedRevisionReceiptInputs(t)
	application.ApplicationDigest = digestString("tampered")
	receipt, err := ObserveRevisionApplicationReceipt(plan, application)
	if err == nil {
		t.Fatal("ObserveRevisionApplicationReceipt() error = nil, want application failure")
	}
	if receipt.Status != "UNKNOWN" || receipt.MissingStage != "revision-application-receipt-application" {
		t.Fatalf("unexpected unknown receipt: %#v", receipt)
	}
}

func TestObserveRevisionApplicationReceiptRetainsLinkFailure(t *testing.T) {
	plan, application := appliedRevisionReceiptInputs(t)
	application.InputIRDigest = digestString("other-input-ir")
	application.ApplicationDigest = digestRevisionApplication(application)
	receipt, err := ObserveRevisionApplicationReceipt(plan, application)
	if err == nil {
		t.Fatal("ObserveRevisionApplicationReceipt() error = nil, want link failure")
	}
	if receipt.Status != "UNKNOWN" || receipt.MissingStage != "revision-application-receipt-link" {
		t.Fatalf("unexpected unknown receipt: %#v", receipt)
	}
}

func TestObserveRevisionApplicationReceiptIsDeterministic(t *testing.T) {
	plan, application := appliedRevisionReceiptInputs(t)
	first, err := ObserveRevisionApplicationReceipt(plan, application)
	if err != nil {
		t.Fatalf("first ObserveRevisionApplicationReceipt() error = %v", err)
	}
	second, err := ObserveRevisionApplicationReceipt(plan, application)
	if err != nil {
		t.Fatalf("second ObserveRevisionApplicationReceipt() error = %v", err)
	}
	if first.ReceiptDigest != second.ReceiptDigest {
		t.Fatal("same plan and application produced different receipt digest")
	}
}
