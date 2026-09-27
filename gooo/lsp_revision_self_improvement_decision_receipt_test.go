package gooo

import "testing"

func lspDecisionReceiptInput(t *testing.T) RevisionSelfImprovementDecisionReceiptObservation {
	t.Helper()
	replan, nextIteration := provenanceNextIterationInputs(t)
	boundary, err := ObserveRevisionSelfImprovementProvenanceNextIteration(replan, nextIteration)
	if err != nil {
		t.Fatalf("observe provenance next iteration: %v", err)
	}
	receipt, err := ObserveRevisionSelfImprovementProvenanceDecisionReceipt(
		boundary,
		digestString("next-iteration-question"),
		revisionSelfImprovementDecisionReceiptRequireReplan,
		"next iteration boundary is not aligned",
	)
	if err != nil {
		t.Fatalf("observe provenance decision receipt: %v", err)
	}
	return receipt
}

func TestObserveLSPRevisionSelfImprovementDecisionReceiptBindsProjection(t *testing.T) {
	receipt := lspDecisionReceiptInput(t)
	got, err := ObserveLSPRevisionSelfImprovementDecisionReceipt(receipt)
	if err != nil {
		t.Fatalf("observe lsp decision receipt: %v", err)
	}
	if got.Status != "BOUND" || got.MissingStage != "" {
		t.Fatalf("unexpected lsp decision receipt status: %#v", got)
	}
	if got.DecisionReceiptDigest != receipt.ObservationDigest ||
		got.DecisionSignal != receipt.DecisionSignal ||
		got.QuestionDigest != receipt.QuestionDigest {
		t.Fatalf("projection lost receipt evidence: %#v", got)
	}
	if !got.ReadOnly || !got.NonExecuting || !got.NonAuthorizing || !got.DecisionAligned {
		t.Fatalf("projection safety flags = %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate lsp decision receipt: %v", err)
	}
}

func TestObserveLSPRevisionSelfImprovementDecisionReceiptPreservesUnknown(t *testing.T) {
	receipt := lspDecisionReceiptInput(t)
	receipt.ObservationDigest = digestString("tampered")
	got, err := ObserveLSPRevisionSelfImprovementDecisionReceipt(receipt)
	if err == nil {
		t.Fatal("expected tampered receipt error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "lsp-revision-self-improvement-decision-receipt-input" {
		t.Fatalf("unexpected unknown lsp decision receipt: %#v", got)
	}
}

func TestObserveLSPRevisionSelfImprovementDecisionReceiptIsDeterministic(t *testing.T) {
	receipt := lspDecisionReceiptInput(t)
	first, err := ObserveLSPRevisionSelfImprovementDecisionReceipt(receipt)
	if err != nil {
		t.Fatalf("first projection: %v", err)
	}
	second, err := ObserveLSPRevisionSelfImprovementDecisionReceipt(receipt)
	if err != nil {
		t.Fatalf("second projection: %v", err)
	}
	if first != second {
		t.Fatalf("projections differ: %#v != %#v", first, second)
	}
}