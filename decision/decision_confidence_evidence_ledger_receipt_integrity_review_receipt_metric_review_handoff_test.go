package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewHandoffIsNonExecuting(t *testing.T) {
	counterexample := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricCounterexample{
		FeedbackDigest: "feedback-digest",
		Delta:          "changed",
		Status:         "counterexample",
		NonAuthorizing: true,
	}
	handoff, err := DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewHandoff(counterexample)
	if err != nil {
		t.Fatalf("derive changed review handoff: %v", err)
	}
	if handoff.Status != "ready-for-external-review" || handoff.CounterexampleDigest == "" || !handoff.NonAuthorizing {
		t.Fatalf("unexpected changed review handoff: %#v", handoff)
	}

	counterexample.Status = "no-counterexample"
	handoff, err = DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewHandoff(counterexample)
	if err != nil {
		t.Fatalf("derive observation handoff: %v", err)
	}
	if handoff.Status != "observation-only" || !handoff.NonAuthorizing {
		t.Fatalf("unexpected observation handoff: %#v", handoff)
	}

	counterexample.Status = "unknown"
	handoff, err = DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewHandoff(counterexample)
	if err != nil {
		t.Fatalf("derive held handoff: %v", err)
	}
	if handoff.Status != "hold" || !handoff.NonAuthorizing {
		t.Fatalf("unknown counterexample escaped hold: %#v", handoff)
	}
}
