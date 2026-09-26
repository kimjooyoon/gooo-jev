package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcome
// records external review without authorizing a change.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcome struct {
	HandoffDigest  string `json:"handoff_digest"`
	Decision       string `json:"decision"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcome(handoff DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewHandoff, decision string) (DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcome, error) {
	digest, err := Digest(handoff)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcome{}, err
	}
	outcome := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcome{
		HandoffDigest:  digest,
		Decision:       decision,
		Status:         "unknown",
		NonAuthorizing: true,
	}
	if handoff.Status != "ready-for-external-review" {
		return outcome, nil
	}
	switch decision {
	case "accepted-for-analysis", "rejected":
		outcome.Status = decision
	}
	return outcome, nil
}
