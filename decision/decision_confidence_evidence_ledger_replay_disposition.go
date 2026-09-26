package decision

// DecisionConfidenceEvidenceLedgerReplayDisposition is a non-executing boundary
// between replay evidence and any separately authorized replay operation.
type DecisionConfidenceEvidenceLedgerReplayDisposition struct {
	ObservationDigest string `json:"observation_digest"`
	Status            string `json:"status"`
	NonAuthorizing    bool   `json:"non_authorizing"`
}

func DeriveDecisionConfidenceEvidenceLedgerReplayDisposition(observation DecisionConfidenceEvidenceLedgerReplayObservation) (DecisionConfidenceEvidenceLedgerReplayDisposition, error) {
	digest, err := Digest(observation)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReplayDisposition{}, err
	}
	disposition := DecisionConfidenceEvidenceLedgerReplayDisposition{
		ObservationDigest: digest,
		Status:            "hold",
		NonAuthorizing:    true,
	}
	if observation.Status == "replayable" {
		disposition.Status = "ready-for-external-replay"
	}
	return disposition, nil
}
