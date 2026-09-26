package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewHandoff
// carries integrity counterexamples to external review without executing them.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewHandoff struct {
	CounterexampleDigest string `json:"counterexample_digest"`
	Status              string `json:"status"`
	NonAuthorizing      bool   `json:"non_authorizing"`
}

func DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewHandoff(counterexample DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricCounterexample) (DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewHandoff, error) {
	digest, err := Digest(counterexample)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewHandoff{}, err
	}
	handoff := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewHandoff{
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
