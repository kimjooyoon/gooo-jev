package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewHandoff
// carries counterexamples to external review without executing them.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewHandoff struct {
	CounterexampleDigest string `json:"counterexample_digest"`
	Status              string `json:"status"`
	NonAuthorizing      bool   `json:"non_authorizing"`
}

func DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewHandoff(counterexample DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackCounterexample) (DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewHandoff, error) {
	digest, err := Digest(counterexample)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewHandoff{}, err
	}
	handoff := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewHandoff{
		CounterexampleDigest: digest,
		Status:              "hold",
		NonAuthorizing:      true,
	}
	switch counterexample.Status {
	case "counterexample":
		handoff.Status = "ready-for-external-review"
	case "no-counterexample":
		handoff.Status = "observation-only"
	}
	return handoff, nil
}
