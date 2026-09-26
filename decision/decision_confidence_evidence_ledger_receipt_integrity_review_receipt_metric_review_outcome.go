package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewOutcome
// records external review without authorizing a change.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewOutcome struct {
	HandoffDigest  string `json:"handoff_digest"`
	Decision       string `json:"decision"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewOutcome(handoff DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewHandoff, decision string) (DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewOutcome, error) {
	digest, err := Digest(handoff)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewOutcome{}, err
	}
	outcome := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewOutcome{
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
