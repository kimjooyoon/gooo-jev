package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedback
// converts an external review outcome into non-authorizing feedback.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedback struct {
	OutcomeDigest  string `json:"outcome_digest"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

func DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedback(outcome DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcome) (DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedback, error) {
	digest, err := Digest(outcome)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedback{}, err
	}
	feedback := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedback{
		OutcomeDigest:  digest,
		Status:        "hold",
		NonAuthorizing: true,
	}
	switch outcome.Status {
	case "accepted-for-analysis":
		feedback.Status = "analysis-ready"
	case "rejected":
		feedback.Status = "rejected"
	}
	return feedback, nil
}

