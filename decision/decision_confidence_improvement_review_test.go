package decision

import (
	"testing"
	"time"
)

func TestReviewDecisionConfidenceImprovementSignalRequiresIndependentEvidence(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	signal := DecisionConfidenceImprovementSignal{
		Schema:            DecisionConfidenceImprovementSignalSchemaV1,
		AssessmentDigest:  "assessment-digest",
		ObservationDigest: "observation-digest",
		ReplayDigest:      "replay-digest",
		Status:            DecisionConfidenceCandidate,
		Action:            "request-additional-confidence-evidence",
		MissingStage:      "confidence-threshold",
		NonAuthorizing:    true,
		ObservedAt:        now,
	}
	if err := signal.assignDigest(); err != nil {
		t.Fatalf("signal.assignDigest() error = %v", err)
	}
	receipt, err := ReviewDecisionConfidenceImprovementSignal(signal, "reviewer-digest", "proposed-change-digest", true, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("ReviewDecisionConfidenceImprovementSignal() error = %v", err)
	}
	if receipt.Decision != DecisionConfidenceImprovementApproved || !receipt.NonExecution {
		t.Fatalf("receipt = %#v", receipt)
	}
	if err := receipt.Validate(); err != nil {
		t.Fatalf("receipt Validate() error = %v", err)
	}
}

func TestReviewDecisionConfidenceImprovementSignalPreservesRejection(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	signal := DecisionConfidenceImprovementSignal{Schema: DecisionConfidenceImprovementSignalSchemaV1, AssessmentDigest: "assessment", ObservationDigest: "observation", ReplayDigest: "replay", Status: DecisionConfidenceCandidate, Action: "reobserve-reverse-output", MissingStage: "reverse-output", NonAuthorizing: true, ObservedAt: now}
	if err := signal.assignDigest(); err != nil {
		t.Fatalf("signal.assignDigest() error = %v", err)
	}
	receipt, err := ReviewDecisionConfidenceImprovementSignal(signal, "reviewer", "change", false, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("ReviewDecisionConfidenceImprovementSignal() error = %v", err)
	}
	if receipt.Decision != DecisionConfidenceImprovementRejected {
		t.Fatalf("rejected receipt = %#v", receipt)
	}
	if err := receipt.Validate(); err != nil {
		t.Fatalf("rejected receipt Validate() error = %v", err)
	}
}

func TestReviewDecisionConfidenceImprovementSignalPreservesNoChange(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	signal := DecisionConfidenceImprovementSignal{Schema: DecisionConfidenceImprovementSignalSchemaV1, AssessmentDigest: "assessment", ObservationDigest: "observation", ReplayDigest: "replay", Status: DecisionConfidenceNoChange, Action: "retain-observed-provenance", NonAuthorizing: true, ObservedAt: now}
	if err := signal.assignDigest(); err != nil {
		t.Fatalf("signal.assignDigest() error = %v", err)
	}
	receipt, err := ReviewDecisionConfidenceImprovementSignal(signal, "", "", false, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("ReviewDecisionConfidenceImprovementSignal() error = %v", err)
	}
	if receipt.Decision != DecisionConfidenceImprovementNoChange {
		t.Fatalf("no-change receipt = %#v", receipt)
	}
	if err := receipt.Validate(); err != nil {
		t.Fatalf("no-change receipt Validate() error = %v", err)
	}
}

func TestReviewDecisionConfidenceImprovementSignalRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	receipt := DecisionConfidenceImprovementReviewReceipt{Schema: DecisionConfidenceImprovementReviewSchemaV1, SignalDigest: "signal", ReviewerDigest: "reviewer", ProposedChangeDigest: "change", Decision: DecisionConfidenceImprovementApproved, NonExecution: true, ReviewedAt: now}
	if err := receipt.assignDigest(); err != nil {
		t.Fatalf("receipt.assignDigest() error = %v", err)
	}
	receipt.ProposedChangeDigest = "tampered-change"
	if err := receipt.Validate(); err == nil {
		t.Fatal("tampered improvement review receipt unexpectedly validated")
	}
}
