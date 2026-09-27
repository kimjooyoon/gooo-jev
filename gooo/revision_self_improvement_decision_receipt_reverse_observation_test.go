package gooo

import "testing"

func TestObserveRevisionSelfImprovementDecisionReceiptReverseObservationBindsEvidence(t *testing.T) {
	receipt := lspDecisionReceiptInput(t)
	_, reverse := provenanceReverseObservationInputs(t)
	got, err := ObserveRevisionSelfImprovementDecisionReceiptReverseObservation(receipt, reverse)
	if err != nil {
		t.Fatalf("observe decision receipt reverse boundary: %v", err)
	}
	if got.Status != "BOUND" || got.MissingStage != "" {
		t.Fatalf("unexpected decision receipt reverse boundary: %#v", got)
	}
	if got.DecisionReceiptDigest != receipt.ObservationDigest ||
		got.ReverseObservationDigest != reverse.ObservationDigest {
		t.Fatalf("reverse boundary lost evidence: %#v", got)
	}
	if got.DecisionFeedbackSignal != "decision-receipt-reverse-mismatch" &&
		got.DecisionFeedbackSignal != "decision-receipt-reverse-aligned" {
		t.Fatalf("unexpected relationship signal: %#v", got)
	}
	if !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("reverse boundary safety flags = %#v", got)
	}
}

func TestObserveRevisionSelfImprovementDecisionReceiptReverseObservationPreservesUnknown(t *testing.T) {
	receipt := lspDecisionReceiptInput(t)
	_, reverse := provenanceReverseObservationInputs(t)
	reverse.ObservationDigest = digestString("tampered")
	got, err := ObserveRevisionSelfImprovementDecisionReceiptReverseObservation(receipt, reverse)
	if err == nil {
		t.Fatal("expected reverse observation error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "revision-self-improvement-decision-receipt-reverse-observation" {
		t.Fatalf("unexpected unknown reverse boundary: %#v", got)
	}
}

func TestObserveRevisionSelfImprovementDecisionReceiptReverseObservationIsDeterministic(t *testing.T) {
	receipt := lspDecisionReceiptInput(t)
	_, reverse := provenanceReverseObservationInputs(t)
	first, err := ObserveRevisionSelfImprovementDecisionReceiptReverseObservation(receipt, reverse)
	if err != nil {
		t.Fatalf("first reverse boundary: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementDecisionReceiptReverseObservation(receipt, reverse)
	if err != nil {
		t.Fatalf("second reverse boundary: %v", err)
	}
	if first != second {
		t.Fatalf("reverse boundaries differ: %#v != %#v", first, second)
	}
}