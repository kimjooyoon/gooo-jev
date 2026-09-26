package decision

// DecisionConfidenceEvidenceLedgerReplayCounterexample records a review-required
// replay signal as a counterexample observation, not as a candidate.
type DecisionConfidenceEvidenceLedgerReplayCounterexample struct {
	SignalDigest   string `json:"signal_digest"`
	Outcome        string `json:"outcome"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerReplayCounterexample(signal DecisionConfidenceEvidenceLedgerReplayReviewSignal) (DecisionConfidenceEvidenceLedgerReplayCounterexample, error) {
	digest, err := Digest(signal)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReplayCounterexample{}, err
	}
	observation := DecisionConfidenceEvidenceLedgerReplayCounterexample{
		SignalDigest:   digest,
		Outcome:        signal.Outcome,
		Status:         "unknown",
		NonAuthorizing: true,
	}
	switch signal.Status {
	case "review-required":
		observation.Status = "counterexample"
	case "observation-only":
		observation.Status = "no-counterexample"
	}
	return observation, nil
}
