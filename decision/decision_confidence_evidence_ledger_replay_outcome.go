package decision

// DecisionConfidenceEvidenceLedgerReplayOutcomeObservation records an externally
// reported replay outcome without interpreting it as improvement.
type DecisionConfidenceEvidenceLedgerReplayOutcomeObservation struct {
	DispositionDigest string `json:"disposition_digest"`
	Outcome          string `json:"outcome"`
	Status           string `json:"status"`
	NonAuthorizing   bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerReplayOutcome(disposition DecisionConfidenceEvidenceLedgerReplayDisposition, outcome string) (DecisionConfidenceEvidenceLedgerReplayOutcomeObservation, error) {
	digest, err := Digest(disposition)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReplayOutcomeObservation{}, err
	}
	observation := DecisionConfidenceEvidenceLedgerReplayOutcomeObservation{
		DispositionDigest: digest,
		Outcome:           outcome,
		Status:            "unknown",
		NonAuthorizing:    true,
	}
	if disposition.Status != "ready-for-external-replay" {
		return observation, nil
	}
	switch outcome {
	case "reproduced", "diverged", "rolled-back", "not-applied":
		observation.Status = outcome
	}
	return observation, nil
}
