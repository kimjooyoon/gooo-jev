package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeDoesNotAuthorize(t *testing.T) {
	handoff := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewHandoff{
		CounterexampleDigest: "counterexample",
		Status:              "ready-for-external-review",
		NonAuthorizing:      true,
	}
	outcome, err := ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcome(handoff, "accepted-for-analysis")
	if err != nil {
		t.Fatalf("observe accepted review: %v", err)
	}
	if outcome.Status != "accepted-for-analysis" || outcome.HandoffDigest == "" || !outcome.NonAuthorizing {
		t.Fatalf("unexpected accepted review: %#v", outcome)
	}

	outcome, err = ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcome(handoff, "approved")
	if err != nil {
		t.Fatalf("observe invalid review: %v", err)
	}
	if outcome.Status != "unknown" || !outcome.NonAuthorizing {
		t.Fatalf("invalid decision escaped unknown: %#v", outcome)
	}

	handoff.Status = "hold"
	outcome, err = ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcome(handoff, "rejected")
	if err != nil {
		t.Fatalf("observe held review: %v", err)
	}
	if outcome.Status != "unknown" || !outcome.NonAuthorizing {
		t.Fatalf("held handoff escaped unknown: %#v", outcome)
	}
}
