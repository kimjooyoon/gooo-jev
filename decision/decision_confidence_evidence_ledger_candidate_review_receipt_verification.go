package decision

// DecisionConfidenceEvidenceLedgerCandidateReviewReceiptVerification checks the
// receipt identity without authorizing or replaying anything.
type DecisionConfidenceEvidenceLedgerCandidateReviewReceiptVerification struct {
	ReceiptDigest  string `json:"receipt_digest"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

type decisionConfidenceEvidenceLedgerCandidateReviewReceiptVerificationIdentity struct {
	SignalDigest        string `json:"signal_digest"`
	ReviewOutcomeDigest  string `json:"review_outcome_digest"`
}

func VerifyDecisionConfidenceEvidenceLedgerCandidateReviewReceipt(receipt DecisionConfidenceEvidenceLedgerCandidateReviewReceipt) DecisionConfidenceEvidenceLedgerCandidateReviewReceiptVerification {
	verification := DecisionConfidenceEvidenceLedgerCandidateReviewReceiptVerification{
		ReceiptDigest:  receipt.ReceiptDigest,
		Status:         "unknown",
		NonAuthorizing: true,
	}
	if receipt.ReceiptDigest == "" || !receipt.NonAuthorizing {
		return verification
	}
	identity := decisionConfidenceEvidenceLedgerCandidateReviewReceiptVerificationIdentity{
		SignalDigest:       receipt.SignalDigest,
		ReviewOutcomeDigest: receipt.ReviewOutcomeDigest,
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
