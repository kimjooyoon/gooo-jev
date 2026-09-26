package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewHandoff
// carries comparison counterexamples to external review without executing them.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewHandoff struct {
	CounterexampleDigest string `json:"counterexample_digest"`
	Status              string `json:"status"`
	NonAuthorizing      bool   `json:"non_authorizing"`
}

func DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewHandoff(counterexample DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackCounterexample) (DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewHandoff, error) {
	digest, err := Digest(counterexample)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewHandoff{}, err
	}
	handoff := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewHandoff{
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
