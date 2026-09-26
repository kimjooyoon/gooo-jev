package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetric is a
// one-hot integrity observation, not an authorization or improvement score.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetric struct {
	ReceiptDigest  string `json:"receipt_digest"`
	VerifiedCount  uint64 `json:"verified_count"`
	MismatchCount  uint64 `json:"mismatch_count"`
	UnknownCount   uint64 `json:"unknown_count"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetric(verification DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptVerification) DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetric {
	metric := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetric{
		ReceiptDigest:  verification.ReceiptDigest,
		Status:         "unknown",
		NonAuthorizing: true,
	}
	switch verification.Status {
	case "verified":
		metric.VerifiedCount = 1
		metric.Status = "verified"
	case "mismatch":
		metric.MismatchCount = 1
		metric.Status = "mismatch"
	default:
		metric.UnknownCount = 1
	}
	return metric
}
