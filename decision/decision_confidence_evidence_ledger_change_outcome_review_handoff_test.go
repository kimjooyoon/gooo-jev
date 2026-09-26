package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerChangeOutcomeReviewHandoffIsNonExecuting(t *testing.T) {
	counterexample := DecisionConfidenceEvidenceLedgerChangeOutcomeCounterexample{
		FeedbackDigest: "feedback-digest",
		Delta:          "rolled-back",
		Status:         "counterexample",
		NonAuthorizing: true,
	}
	handoff, err := DeriveDecisionConfidenceEvidenceLedgerChangeOutcomeReviewHandoff(counterexample)
	if err != nil {
		t.Fatalf("derive review handoff: %v", err)
	}
	if handoff.Status != "ready-for-external-review" || handoff.CounterexampleDigest == "" || !handoff.NonAuthorizing {
		t.Fatalf("unexpected review handoff: %#v", handoff)
	}

	counterexample.Status = "unknown"
	handoff, err = DeriveDecisionConfidenceEvidenceLedgerChangeOutcomeReviewHandoff(counterexample)
	if err != nil {
		t.Fatalf("derive hold handoff: %v", err)
	}
	if handoff.Status != "hold" || !handoff.NonAuthorizing {
		t.Fatalf("unknown counterexample escaped hold: %#v", handoff)
	}
}
