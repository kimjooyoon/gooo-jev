package decision

// DecisionConfidenceEvidenceLedgerChangeOutcomeCounterexample preserves a
// review-required outcome as a counterexample without creating a candidate.
type DecisionConfidenceEvidenceLedgerChangeOutcomeCounterexample struct {
	FeedbackDigest string `json:"feedback_digest"`
	Delta          string `json:"delta"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerChangeOutcomeCounterexample(feedback DecisionConfidenceEvidenceLedgerChangeOutcomeFeedback) (DecisionConfidenceEvidenceLedgerChangeOutcomeCounterexample, error) {
	digest, err := Digest(feedback)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerChangeOutcomeCounterexample{}, err
	}
	counterexample := DecisionConfidenceEvidenceLedgerChangeOutcomeCounterexample{
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
