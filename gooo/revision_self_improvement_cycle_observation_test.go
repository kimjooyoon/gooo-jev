package gooo

import "testing"

func selfImprovementCycleInputs(t *testing.T) (
	RevisionSelfImprovementProvenanceNextIterationObservation,
	RevisionSelfImprovementDecisionReceiptObservation,
	RevisionSelfImprovementDecisionReceiptReverseObservation,
) {
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
		t.Fatalf("observe decision receipt: %v", err)
	}
	feedback, reverseObservation := provenanceReverseObservationInputs(t)
	provenanceReverse, err := ObserveRevisionSelfImprovementProvenanceReverseObservation(feedback, reverseObservation)
	if err != nil {
		t.Fatalf("observe provenance reverse observation: %v", err)
	}
	reverse, err := ObserveRevisionSelfImprovementDecisionReceiptReverseObservation(receipt, provenanceReverse)
	if err != nil {
		t.Fatalf("observe decision receipt reverse observation: %v", err)
	}
	return boundary, receipt, reverse
}

func TestObserveRevisionSelfImprovementCycleBindsReplanRequiredCycle(t *testing.T) {
	nextIteration, receipt, reverse := selfImprovementCycleInputs(t)
	got, err := ObserveRevisionSelfImprovementCycle(nextIteration, receipt, reverse)
	if err != nil {
		t.Fatalf("observe self-improvement cycle: %v", err)
	}
	if got.Status != "BOUND" || got.MissingStage != "" {
		t.Fatalf("unexpected cycle status: %#v", got)
	}
	if got.CycleSignal != "self-improvement-cycle-replan-required" {
		t.Fatalf("cycle signal = %q, want replan-required", got.CycleSignal)
	}
	if got.NextIterationObservationDigest != nextIteration.ObservationDigest ||
		got.DecisionReceiptDigest != receipt.ObservationDigest ||
		got.DecisionReceiptReverseDigest != reverse.ObservationDigest {
		t.Fatalf("cycle lost stage links: %#v", got)
	}
	if got.SignalsAligned {
		t.Fatal("mismatched boundary must not claim a closed cycle")
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate cycle: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCyclePreservesUnknownLink(t *testing.T) {
	nextIteration, receipt, reverse := selfImprovementCycleInputs(t)
	receipt.InputObservationDigest = digestString("tampered")
	receipt.ObservationDigest = digestRevisionSelfImprovementDecisionReceipt(receipt)
	got, err := ObserveRevisionSelfImprovementCycle(nextIteration, receipt, reverse)
	if err == nil {
		t.Fatal("expected decision receipt link error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "revision-self-improvement-cycle-decision-receipt-link" {
		t.Fatalf("unexpected unknown cycle: %#v", got)
	}
}

func TestObserveRevisionSelfImprovementCycleIsDeterministic(t *testing.T) {
	nextIteration, receipt, reverse := selfImprovementCycleInputs(t)
	first, err := ObserveRevisionSelfImprovementCycle(nextIteration, receipt, reverse)
	if err != nil {
		t.Fatalf("first cycle observation: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCycle(nextIteration, receipt, reverse)
	if err != nil {
		t.Fatalf("second cycle observation: %v", err)
	}
	if first != second {
		t.Fatalf("cycle observations differ: %#v != %#v", first, second)
	}
}