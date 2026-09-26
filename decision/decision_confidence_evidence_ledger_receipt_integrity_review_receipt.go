package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceipt is a
// content-addressed record of integrity review evidence only.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceipt struct {
	HandoffDigest   string `json:"handoff_digest"`
	OutcomeDigest   string `json:"outcome_digest"`
	ReceiptDigest   string `json:"receipt_digest"`
	Status          string `json:"status"`
	NonAuthorizing  bool   `json:"non_authorizing"`
}

type decisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptIdentity struct {
	HandoffDigest string `json:"handoff_digest"`
	OutcomeDigest string `json:"outcome_digest"`
}

func BuildDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceipt(handoff DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewHandoff, outcome DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewOutcome) (DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceipt, error) {
	identity := decisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptIdentity{
		HandoffDigest: handoff.CounterexampleDigest,
		OutcomeDigest: outcome.HandoffDigest,
	}
	digest, err := Digest(identity)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceipt{}, err
	}
	receipt := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceipt{
		HandoffDigest:  handoff.CounterexampleDigest,
		OutcomeDigest:  outcome.HandoffDigest,
		ReceiptDigest:  digest,
		Status:         "hold",
		NonAuthorizing: true,
	}
	if handoff.Status == "ready-for-external-review" && outcome.Status == "accepted-for-analysis" {
		receipt.Status = "ready-for-integrity-analysis"
	} else if outcome.Status == "rejected" {
		receipt.Status = "rejected"
	}
	return receipt, nil
}
