package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityCounterexample preserves a
// changed receipt integrity state as a reviewable counterexample.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityCounterexample struct {
	FeedbackDigest string `json:"feedback_digest"`
	Delta          string `json:"delta"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityCounterexample(feedback DecisionConfidenceEvidenceLedgerReceiptIntegrityFeedback) (DecisionConfidenceEvidenceLedgerReceiptIntegrityCounterexample, error) {
	digest, err := Digest(feedback)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReceiptIntegrityCounterexample{}, err
	}
	counterexample := DecisionConfidenceEvidenceLedgerReceiptIntegrityCounterexample{
		FeedbackDigest: digest,
		Delta:          feedback.Delta,
		Status:         "unknown",
		NonAuthorizing: true,
	}
	switch feedback.Status {
	case "review-required":
		counterexample.Status = "counterexample"
	case "observation-only":
		counterexample.Status = "no-counterexample"
	}
	return counterexample, nil
}
