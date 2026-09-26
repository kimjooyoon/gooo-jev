package decision

// DecisionConfidenceEvidenceLedgerChangeOutcomeReviewOutcome records external
// review of a change-outcome counterexample without authorizing a change.
type DecisionConfidenceEvidenceLedgerChangeOutcomeReviewOutcome struct {
	HandoffDigest  string `json:"handoff_digest"`
	Decision       string `json:"decision"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerChangeOutcomeReviewOutcome(handoff DecisionConfidenceEvidenceLedgerChangeOutcomeReviewHandoff, decision string) (DecisionConfidenceEvidenceLedgerChangeOutcomeReviewOutcome, error) {
	digest, err := Digest(handoff)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerChangeOutcomeReviewOutcome{}, err
	}
	outcome := DecisionConfidenceEvidenceLedgerChangeOutcomeReviewOutcome{
		HandoffDigest:  digest,
		Decision:       decision,
		Status:         "unknown",
		NonAuthorizing: true,
	}
	if handoff.Status != "ready-for-external-review" {
		return outcome, nil
	}
	switch decision {
	case "accepted-for-analysis", "rejected":
		outcome.Status = decision
	}
	return outcome, nil
}
