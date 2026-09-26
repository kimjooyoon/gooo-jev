package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricCounterexample
// preserves review-required integrity feedback without creating a candidate.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricCounterexample struct {
	FeedbackDigest string `json:"feedback_digest"`
	Delta          string `json:"delta"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricCounterexample(feedback DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricFeedback) (DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricCounterexample, error) {
	digest, err := Digest(feedback)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricCounterexample{}, err
	}
	counterexample := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricCounterexample{
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
