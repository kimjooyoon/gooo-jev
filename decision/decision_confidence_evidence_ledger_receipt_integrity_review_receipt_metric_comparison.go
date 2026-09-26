package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricComparison
// compares integrity observations without inferring improvement.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricComparison struct {
	PreviousReceiptDigest string `json:"previous_receipt_digest"`
	CurrentReceiptDigest  string `json:"current_receipt_digest"`
	Delta                 string `json:"delta"`
	NonAuthorizing        bool   `json:"non_authorizing"`
}

func CompareDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetric(previous, current DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetric) DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricComparison {
	comparison := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricComparison{
		PreviousReceiptDigest: previous.ReceiptDigest,
		CurrentReceiptDigest:  current.ReceiptDigest,
		Delta:                 "unknown",
		NonAuthorizing:        true,
	}
	if previous.ReceiptDigest == "" || current.ReceiptDigest == "" || !previous.NonAuthorizing || !current.NonAuthorizing {
		return comparison
	}
	previousTotal := previous.VerifiedCount + previous.MismatchCount + previous.UnknownCount
	currentTotal := current.VerifiedCount + current.MismatchCount + current.UnknownCount
	if previousTotal != 1 || currentTotal != 1 {
		comparison.Delta = "inconclusive"
		return comparison
	}
	if previous.Status == current.Status {
		comparison.Delta = "unchanged"
		return comparison
	}
	comparison.Delta = "changed"
	return comparison
}
