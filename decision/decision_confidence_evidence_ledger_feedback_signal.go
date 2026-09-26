package decision

// DecisionConfidenceEvidenceLedgerFeedbackSignal preserves metric comparison
// outcomes as feedback without selecting or applying a change.
type DecisionConfidenceEvidenceLedgerFeedbackSignal struct {
	ComparisonDigest string `json:"comparison_digest"`
	Delta            string `json:"delta"`
	Status           string `json:"status"`
	NonAuthorizing   bool   `json:"non_authorizing"`
}

func DeriveDecisionConfidenceEvidenceLedgerFeedbackSignal(comparison DecisionConfidenceEvidenceLedgerCandidateReviewMetricComparison) (DecisionConfidenceEvidenceLedgerFeedbackSignal, error) {
	digest, err := Digest(comparison)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerFeedbackSignal{}, err
	}
	signal := DecisionConfidenceEvidenceLedgerFeedbackSignal{
		ComparisonDigest: digest,
		Delta:            comparison.Delta,
		Status:           "unknown",
		NonAuthorizing:   true,
	}
	switch comparison.Delta {
	case "declined", "inconclusive", "unknown":
		signal.Status = "review-required"
	case "increased", "unchanged":
		signal.Status = "observation-only"
	}
	return signal, nil
}
