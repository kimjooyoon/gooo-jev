package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewHandoffIsNonExecuting(t *testing.T) {
	counterexample := DecisionConfidenceEvidenceLedgerReceiptIntegrityCounterexample{
		FeedbackDigest: "feedback-digest",
		Delta:          "changed",
		Status:         "counterexample",
		NonAuthorizing: true,
	}
	handoff, err := DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewHandoff(counterexample)
	if err != nil {
		t.Fatalf("derive integrity review handoff: %v", err)
	}
	if handoff.Status != "ready-for-external-review" || handoff.CounterexampleDigest == "" || !handoff.NonAuthorizing {
		t.Fatalf("unexpected integrity review handoff: %#v", handoff)
	}

	counterexample.Status = "unknown"
	handoff, err = DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewHandoff(counterexample)
	if err != nil {
		t.Fatalf("derive held handoff: %v", err)
	}
	if handoff.Status != "hold" || !handoff.NonAuthorizing {
		t.Fatalf("unknown counterexample escaped hold: %#v", handoff)
	}
}
