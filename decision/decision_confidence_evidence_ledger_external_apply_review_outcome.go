package decision

// DecisionConfidenceEvidenceLedgerExternalApplyReviewOutcome records an external
// review of a verified write-set without authorizing its application.
type DecisionConfidenceEvidenceLedgerExternalApplyReviewOutcome struct {
	DispositionDigest string `json:"disposition_digest"`
	Decision          string `json:"decision"`
	Status            string `json:"status"`
	NonAuthorizing    bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerExternalApplyReviewOutcome(disposition DecisionConfidenceEvidenceLedgerWriteSetDisposition, decision string) (DecisionConfidenceEvidenceLedgerExternalApplyReviewOutcome, error) {
	digest, err := Digest(disposition)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerExternalApplyReviewOutcome{}, err
	}
	outcome := DecisionConfidenceEvidenceLedgerExternalApplyReviewOutcome{
		DispositionDigest: digest,
		Decision:          decision,
		Status:            "unknown",
		NonAuthorizing:    true,
	}
	if disposition.Status != "ready-for-external-apply-review" {
		return outcome, nil
	}
	switch decision {
	case "accepted-for-apply-review", "rejected":
		outcome.Status = decision
	}
	return outcome, nil
}
