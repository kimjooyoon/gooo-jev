package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackCounterexample
// preserves review-required comparison feedback without creating a candidate.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackCounterexample struct {
	FeedbackDigest string `json:"feedback_digest"`
	Delta          string `json:"delta"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackCounterexample(feedback DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedback) (DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackCounterexample, error) {
	digest, err := Digest(feedback)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackCounterexample{}, err
	}
	counterexample := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackCounterexample{
		FeedbackDigest: digest,
		Delta:          feedback.Delta,
		Status:         "unknown",
		NonAuthorizing: true,
	}
	switch feedback.Status {
	case "review-required":
		counterexample.Status = "counterexample"
	case "observation-only":
		counterexample.Status = "no-counterexample"
	}
	return counterexample, nil
}
