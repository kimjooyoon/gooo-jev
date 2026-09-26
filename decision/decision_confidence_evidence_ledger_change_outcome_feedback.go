package decision

// DecisionConfidenceEvidenceLedgerChangeOutcomeFeedback preserves outcome
// changes as review input without treating them as improvement.
type DecisionConfidenceEvidenceLedgerChangeOutcomeFeedback struct {
	PlanDigest      string `json:"plan_digest"`
	Delta           string `json:"delta"`
	Status          string `json:"status"`
	NonAuthorizing  bool   `json:"non_authorizing"`
}

func DeriveDecisionConfidenceEvidenceLedgerChangeOutcomeFeedback(comparison DecisionConfidenceEvidenceLedgerChangeOutcomeComparison) (DecisionConfidenceEvidenceLedgerChangeOutcomeFeedback, error) {
	digest, err := Digest(comparison)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerChangeOutcomeFeedback{}, err
	}
	feedback := DecisionConfidenceEvidenceLedgerChangeOutcomeFeedback{
		PlanDigest:     digest,
		Delta:          comparison.Delta,
		Status:         "hold",
		NonAuthorizing: true,
	}
	switch comparison.Delta {
	case "rolled-back", "changed":
		feedback.Status = "review-required"
	case "unchanged":
		feedback.Status = "observation-only"
	}
	return feedback, nil
}
