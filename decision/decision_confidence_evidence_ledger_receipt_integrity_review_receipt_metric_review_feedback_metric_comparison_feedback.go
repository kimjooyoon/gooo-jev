package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedback
// preserves comparison outcomes as review input without selecting a change.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedback struct {
	ComparisonDigest string `json:"comparison_digest"`
	Delta            string `json:"delta"`
	Status           string `json:"status"`
	NonAuthorizing   bool   `json:"non_authorizing"`
}

func DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedback(comparison DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparison) (DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedback, error) {
	digest, err := Digest(comparison)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedback{}, err
	}
	feedback := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedback{
		ComparisonDigest: digest,
		Delta:            comparison.Delta,
		Status:           "hold",
		NonAuthorizing:   true,
	}
	switch comparison.Delta {
	case "declined":
		feedback.Status = "review-required"
	case "increased", "unchanged":
		feedback.Status = "observation-only"
	}
	return feedback, nil
}
