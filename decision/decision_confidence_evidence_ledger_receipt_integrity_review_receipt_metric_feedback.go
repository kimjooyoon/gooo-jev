package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricFeedback
// preserves integrity metric changes as review input without selecting a change.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricFeedback struct {
	ComparisonDigest string `json:"comparison_digest"`
	Delta            string `json:"delta"`
	Status           string `json:"status"`
	NonAuthorizing   bool   `json:"non_authorizing"`
}

func DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricFeedback(comparison DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricComparison) (DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricFeedback, error) {
	digest, err := Digest(comparison)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricFeedback{}, err
	}
	feedback := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricFeedback{
		ComparisonDigest: digest,
		Delta:            comparison.Delta,
		Status:           "hold",
		NonAuthorizing:   true,
	}
	switch comparison.Delta {
	case "changed":
		feedback.Status = "review-required"
	case "unchanged":
		feedback.Status = "observation-only"
	}
	return feedback, nil
}
