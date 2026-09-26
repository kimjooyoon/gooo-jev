package decision

// DecisionConfidenceEvidenceLedgerChangeOutcomeReviewHandoff is a non-executing
// handoff for review of a change-outcome counterexample.
type DecisionConfidenceEvidenceLedgerChangeOutcomeReviewHandoff struct {
	CounterexampleDigest string `json:"counterexample_digest"`
	Status              string `json:"status"`
	NonAuthorizing      bool   `json:"non_authorizing"`
}

func DeriveDecisionConfidenceEvidenceLedgerChangeOutcomeReviewHandoff(counterexample DecisionConfidenceEvidenceLedgerChangeOutcomeCounterexample) (DecisionConfidenceEvidenceLedgerChangeOutcomeReviewHandoff, error) {
	digest, err := Digest(counterexample)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerChangeOutcomeReviewHandoff{}, err
	}
	handoff := DecisionConfidenceEvidenceLedgerChangeOutcomeReviewHandoff{
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
