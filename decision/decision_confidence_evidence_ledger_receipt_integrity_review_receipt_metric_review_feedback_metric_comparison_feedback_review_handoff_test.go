package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewHandoffIsNonExecuting(t *testing.T) {
	counterexample := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackCounterexample{
		FeedbackDigest: "feedback",
		Delta:          "declined",
		Status:         "counterexample",
		NonAuthorizing: true,
	}
	handoff, err := DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewHandoff(counterexample)
	if err != nil {
		t.Fatalf("derive counterexample handoff: %v", err)
	}
	if handoff.Status != "ready-for-external-review" || handoff.CounterexampleDigest == "" || !handoff.NonAuthorizing {
		t.Fatalf("unexpected counterexample handoff: %#v", handoff)
	}

	counterexample.Status = "no-counterexample"
	handoff, err = DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewHandoff(counterexample)
	if err != nil {
		t.Fatalf("derive observation handoff: %v", err)
	}
	if handoff.Status != "observation-only" || !handoff.NonAuthorizing {
		t.Fatalf("unexpected observation handoff: %#v", handoff)
	}

	counterexample.Status = "unknown"
	handoff, err = DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewHandoff(counterexample)
	if err != nil {
		t.Fatalf("derive held handoff: %v", err)
	}
	if handoff.Status != "hold" || !handoff.NonAuthorizing {
		t.Fatalf("unknown counterexample escaped hold: %#v", handoff)
	}
}
