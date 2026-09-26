package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackCounterexamplePreservesBoundary(t *testing.T) {
	feedback := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedback{
		ComparisonDigest: "comparison",
		Delta:            "declined",
		Status:           "review-required",
		NonAuthorizing:   true,
	}
	counterexample, err := ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackCounterexample(feedback)
	if err != nil {
		t.Fatalf("observe review-required counterexample: %v", err)
	}
	if counterexample.Status != "counterexample" || counterexample.FeedbackDigest == "" || !counterexample.NonAuthorizing {
		t.Fatalf("unexpected review-required counterexample: %#v", counterexample)
	}

	feedback.Status = "observation-only"
	counterexample, err = ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackCounterexample(feedback)
	if err != nil {
		t.Fatalf("observe observation-only counterexample: %v", err)
	}
	if counterexample.Status != "no-counterexample" || !counterexample.NonAuthorizing {
		t.Fatalf("unexpected observation-only counterexample: %#v", counterexample)
	}

	feedback.Status = "hold"
	counterexample, err = ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackCounterexample(feedback)
	if err != nil {
		t.Fatalf("observe held counterexample: %v", err)
	}
	if counterexample.Status != "unknown" || !counterexample.NonAuthorizing {
		t.Fatalf("held feedback escaped unknown: %#v", counterexample)
	}
}
