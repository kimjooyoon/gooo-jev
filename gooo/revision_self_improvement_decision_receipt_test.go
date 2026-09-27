package gooo

import "testing"

func TestObserveRevisionSelfImprovementProvenanceDecisionReceiptBindsReplan(t *testing.T) {
	replan, nextIteration := provenanceNextIterationInputs(t)
	input, err := ObserveRevisionSelfImprovementProvenanceNextIteration(replan, nextIteration)
	if err != nil {
		t.Fatalf("observe provenance next iteration: %v", err)
	}
	got, err := ObserveRevisionSelfImprovementProvenanceDecisionReceipt(
		input,
		digestString("next-iteration-question"),
		revisionSelfImprovementDecisionReceiptRequireReplan,
		"next iteration boundary is not aligned",
	)
	if err != nil {
		t.Fatalf("observe provenance decision receipt: %v", err)
	}
	if got.Status != revisionSelfImprovementDecisionReceiptBoundStatus {
		t.Fatalf("status = %q, want %q", got.Status, revisionSelfImprovementDecisionReceiptBoundStatus)
	}
	if got.MissingStage != "" {
		t.Fatalf("missing stage = %q, want empty", got.MissingStage)
	}
	if got.InputObservationDigest != input.ObservationDigest {
		t.Fatalf("input digest = %q, want %q", got.InputObservationDigest, input.ObservationDigest)
	}
	if got.DecisionSignal != revisionSelfImprovementDecisionReceiptRequireReplan {
		t.Fatalf("decision signal = %q, want require-replan", got.DecisionSignal)
	}
	if got.ReceiptSignal != revisionSelfImprovementDecisionReceiptAcceptedSignal {
		t.Fatalf("receipt signal = %q, want accepted", got.ReceiptSignal)
	}
	if !got.DecisionAligned || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("decision receipt flags = %#v", got)
	}
}

func TestObserveRevisionSelfImprovementProvenanceDecisionReceiptPreservesUnknown(t *testing.T) {
	replan, nextIteration := provenanceNextIterationInputs(t)
	input, err := ObserveRevisionSelfImprovementProvenanceNextIteration(replan, nextIteration)
	if err != nil {
		t.Fatalf("observe provenance next iteration: %v", err)
	}
	got, err := ObserveRevisionSelfImprovementProvenanceDecisionReceipt(
		input,
		digestString("next-iteration-question"),
		revisionSelfImprovementDecisionReceiptAllowNextIteration,
		"incorrectly permits a mismatched transition",
	)
	if err == nil {
		t.Fatal("expected decision alignment error")
	}
	if got.Status != revisionSelfImprovementDecisionReceiptUnknownStatus {
		t.Fatalf("status = %q, want %q", got.Status, revisionSelfImprovementDecisionReceiptUnknownStatus)
	}
	if got.MissingStage != revisionSelfImprovementDecisionReceiptMissingAlignment {
		t.Fatalf("missing stage = %q, want %q", got.MissingStage, revisionSelfImprovementDecisionReceiptMissingAlignment)
	}
	if got.DecisionAligned {
		t.Fatal("unknown decision receipt must not claim alignment")
	}
	if !got.NonExecuting || !got.NonAuthorizing {
		t.Fatal("unknown decision receipt must preserve safety flags")
	}
}

func TestObserveRevisionSelfImprovementProvenanceDecisionReceiptRejectsTamperedInput(t *testing.T) {
	replan, nextIteration := provenanceNextIterationInputs(t)
	input, err := ObserveRevisionSelfImprovementProvenanceNextIteration(replan, nextIteration)
	if err != nil {
		t.Fatalf("observe provenance next iteration: %v", err)
	}
	input.ObservationDigest = digestString("tampered")
	got, err := ObserveRevisionSelfImprovementProvenanceDecisionReceipt(
		input,
		digestString("next-iteration-question"),
		revisionSelfImprovementDecisionReceiptRequireReplan,
		"tampered input must remain unknown",
	)
	if err == nil {
		t.Fatal("expected input validation error")
	}
	if got.Status != revisionSelfImprovementDecisionReceiptUnknownStatus {
		t.Fatalf("status = %q, want %q", got.Status, revisionSelfImprovementDecisionReceiptUnknownStatus)
	}
	if got.MissingStage != revisionSelfImprovementDecisionReceiptMissingInput {
		t.Fatalf("missing stage = %q, want %q", got.MissingStage, revisionSelfImprovementDecisionReceiptMissingInput)
	}
}

func TestObserveRevisionSelfImprovementProvenanceDecisionReceiptIsDeterministic(t *testing.T) {
	replan, nextIteration := provenanceNextIterationInputs(t)
	input, err := ObserveRevisionSelfImprovementProvenanceNextIteration(replan, nextIteration)
	if err != nil {
		t.Fatalf("observe provenance next iteration: %v", err)
	}
	first, err := ObserveRevisionSelfImprovementProvenanceDecisionReceipt(
		input,
		digestString("next-iteration-question"),
		revisionSelfImprovementDecisionReceiptRequireReplan,
		"stable reason",
	)
	if err != nil {
		t.Fatalf("first observation: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementProvenanceDecisionReceipt(
		input,
		digestString("next-iteration-question"),
		revisionSelfImprovementDecisionReceiptRequireReplan,
		"stable reason",
	)
	if err != nil {
		t.Fatalf("second observation: %v", err)
	}
	if first != second {
		t.Fatalf("observations differ: %#v != %#v", first, second)
	}
}