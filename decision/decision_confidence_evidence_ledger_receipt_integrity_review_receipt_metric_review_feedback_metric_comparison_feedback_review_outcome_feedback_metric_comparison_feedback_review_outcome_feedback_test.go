package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewOutcomeFeedback
// preserves external review as analysis input without selecting a change.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewOutcomeFeedback struct {
	OutcomeDigest  string `json:"outcome_digest"`
	Decision       string `json:"decision"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

func DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewOutcomeFeedback(outcome DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewOutcome) (DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewOutcomeFeedback, error) {
	digest, err := Digest(outcome)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewOutcomeFeedback{}, err
	}
	feedback := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewOutcomeFeedback{
		OutcomeDigest:  digest,
		Decision:       outcome.Decision,
		Status:         "hold",
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
