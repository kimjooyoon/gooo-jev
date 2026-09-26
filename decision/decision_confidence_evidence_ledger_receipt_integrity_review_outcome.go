package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewOutcome records external
// review of receipt integrity evidence without authorizing a change.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewOutcome struct {
	HandoffDigest  string `json:"handoff_digest"`
	Decision       string `json:"decision"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewOutcome(handoff DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewHandoff, decision string) (DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewOutcome, error) {
	digest, err := Digest(handoff)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewOutcome{}, err
	}
	outcome := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewOutcome{
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
