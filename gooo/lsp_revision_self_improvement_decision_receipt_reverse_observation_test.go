package gooo

import "testing"

func TestObserveLSPRevisionSelfImprovementDecisionReceiptReverseObservationBindsProjection(t *testing.T) {
	receipt, reverse := decisionReceiptReverseObservationInput(t)
	boundary, err := ObserveRevisionSelfImprovementDecisionReceiptReverseObservation(receipt, reverse)
	if err != nil {
		t.Fatalf("observe decision receipt reverse boundary: %v", err)
	}
	got, err := ObserveLSPRevisionSelfImprovementDecisionReceiptReverseObservation(boundary)
	if err != nil {
		t.Fatalf("observe lsp decision receipt reverse projection: %v", err)
	}
	if got.Status != "BOUND" || got.MissingStage != "" {
		t.Fatalf("unexpected lsp reverse projection: %#v", got)
	}
	if got.DecisionReceiptDigest != boundary.DecisionReceiptDigest ||
		got.ReverseObservationDigest != boundary.ReverseObservationDigest ||
		got.DecisionFeedbackSignal != boundary.DecisionFeedbackSignal {
		t.Fatalf("projection lost reverse evidence: %#v", got)
	}
	if !got.ReadOnly || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("projection safety flags = %#v", got)
	}
}

func TestObserveLSPRevisionSelfImprovementDecisionReceiptReverseObservationPreservesUnknown(t *testing.T) {
	receipt, reverse := decisionReceiptReverseObservationInput(t)
	boundary, err := ObserveRevisionSelfImprovementDecisionReceiptReverseObservation(receipt, reverse)
	if err != nil {
		t.Fatalf("observe decision receipt reverse boundary: %v", err)
	}
	boundary.ObservationDigest = digestString("tampered")
	got, err := ObserveLSPRevisionSelfImprovementDecisionReceiptReverseObservation(boundary)
	if err == nil {
		t.Fatal("expected tampered reverse boundary error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "lsp-revision-self-improvement-decision-receipt-reverse-input" {
		t.Fatalf("unexpected unknown lsp reverse projection: %#v", got)
	}
}

func TestObserveLSPRevisionSelfImprovementDecisionReceiptReverseObservationIsDeterministic(t *testing.T) {
	receipt, reverse := decisionReceiptReverseObservationInput(t)
	boundary, err := ObserveRevisionSelfImprovementDecisionReceiptReverseObservation(receipt, reverse)
	if err != nil {
		t.Fatalf("observe decision receipt reverse boundary: %v", err)
	}
	first, err := ObserveLSPRevisionSelfImprovementDecisionReceiptReverseObservation(boundary)
	if err != nil {
		t.Fatalf("first lsp reverse projection: %v", err)
	}
	second, err := ObserveLSPRevisionSelfImprovementDecisionReceiptReverseObservation(boundary)
	if err != nil {
		t.Fatalf("second lsp reverse projection: %v", err)
	}
	if first != second {
		t.Fatalf("lsp reverse projections differ: %#v != %#v", first, second)
	}
}