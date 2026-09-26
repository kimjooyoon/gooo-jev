package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptVerification
// checks review receipt identity without authorizing analysis or changes.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptVerification struct {
	ReceiptDigest  string `json:"receipt_digest"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

type decisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptVerificationIdentity struct {
	HandoffDigest string `json:"handoff_digest"`
	OutcomeDigest string `json:"outcome_digest"`
}

func VerifyDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceipt(receipt DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceipt) DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptVerification {
	verification := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptVerification{
		ReceiptDigest:  receipt.ReceiptDigest,
		Status:         "unknown",
		NonAuthorizing: true,
	}
	if receipt.ReceiptDigest == "" || !receipt.NonAuthorizing {
		return verification
	}
	identity := decisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptVerificationIdentity{
		HandoffDigest: receipt.HandoffDigest,
		OutcomeDigest: receipt.OutcomeDigest,
	}
	expected, err := Digest(identity)
	if err != nil {
		return verification
	}
	if expected != receipt.ReceiptDigest {
		verification.Status = "mismatch"
		return verification
	}
	verification.Status = "verified"
	return verification
}
