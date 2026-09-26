package decision

// DecisionConfidenceEvidenceLedgerReplayObservation records whether an existing
// transition and receipt link are sufficient to describe a replay boundary.
// It does not replay, execute, or authorize any action.
type DecisionConfidenceEvidenceLedgerReplayObservation struct {
	TransitionDigest      string `json:"transition_digest"`
	ExecutionReceiptDigest string `json:"execution_receipt_digest"`
	Status                string `json:"status"`
	NonAuthorizing        bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerReplay(transition DecisionConfidenceEvidenceLedgerTransition, link DecisionConfidenceEvidenceLedgerReceiptLink) (DecisionConfidenceEvidenceLedgerReplayObservation, error) {
	transitionDigest, err := Digest(transition)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReplayObservation{}, err
	}
	observation := DecisionConfidenceEvidenceLedgerReplayObservation{
		TransitionDigest:       transitionDigest,
		ExecutionReceiptDigest: link.ExecutionReceiptDigest,
		Status:                 "unknown",
		NonAuthorizing:         true,
	}
	if transition.Status == "extended" && link.Status == "linked" && link.ExecutionReceiptDigest != "" {
		observation.Status = "replayable"
	}
	return observation, nil
}
