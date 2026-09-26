package decision

// DecisionConfidenceEvidenceLedgerReplayReviewSignal marks replay observations
// that need review without creating a candidate or claiming improvement.
type DecisionConfidenceEvidenceLedgerReplayReviewSignal struct {
	ObservationDigest string `json:"observation_digest"`
	Outcome           string `json:"outcome"`
	Status            string `json:"status"`
	NonAuthorizing    bool   `json:"non_authorizing"`
}

func DeriveDecisionConfidenceEvidenceLedgerReplayReviewSignal(observation DecisionConfidenceEvidenceLedgerReplayOutcomeObservation) (DecisionConfidenceEvidenceLedgerReplayReviewSignal, error) {
	digest, err := Digest(observation)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReplayReviewSignal{}, err
	}
	signal := DecisionConfidenceEvidenceLedgerReplayReviewSignal{
		ObservationDigest: digest,
		Outcome:           observation.Outcome,
		Status:            "unknown",
		NonAuthorizing:    true,
	}
	switch observation.Status {
	case "reproduced":
		signal.Status = "observation-only"
	case "diverged", "rolled-back", "not-applied", "unknown":
		signal.Status = "review-required"
	}
	return signal, nil
}
