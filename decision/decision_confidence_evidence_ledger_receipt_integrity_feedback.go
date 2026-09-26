package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityFeedback preserves receipt
// verification transitions as review input without selecting a change.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityFeedback struct {
	ComparisonDigest string `json:"comparison_digest"`
	Delta            string `json:"delta"`
	Status           string `json:"status"`
	NonAuthorizing   bool   `json:"non_authorizing"`
}

func DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityFeedback(comparison DecisionConfidenceEvidenceLedgerCandidateReviewReceiptMetricComparison) (DecisionConfidenceEvidenceLedgerReceiptIntegrityFeedback, error) {
	digest, err := Digest(comparison)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReceiptIntegrityFeedback{}, err
	}
	feedback := DecisionConfidenceEvidenceLedgerReceiptIntegrityFeedback{
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
