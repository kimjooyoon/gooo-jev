package decision

// DecisionConfidenceEvidenceLedgerReplayReviewHandoff is a non-executing handoff
// for external review of a replay counterexample.
type DecisionConfidenceEvidenceLedgerReplayReviewHandoff struct {
	CounterexampleDigest string `json:"counterexample_digest"`
	Status              string `json:"status"`
	NonAuthorizing      bool   `json:"non_authorizing"`
}

func DeriveDecisionConfidenceEvidenceLedgerReplayReviewHandoff(counterexample DecisionConfidenceEvidenceLedgerReplayCounterexample) (DecisionConfidenceEvidenceLedgerReplayReviewHandoff, error) {
	digest, err := Digest(counterexample)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReplayReviewHandoff{}, err
	}
	handoff := DecisionConfidenceEvidenceLedgerReplayReviewHandoff{
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
