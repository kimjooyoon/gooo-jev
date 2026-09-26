package decision

// DecisionConfidenceEvidenceLedgerExternalReviewOutcome records an external
// review decision without authorizing a candidate, replay, or code change.
type DecisionConfidenceEvidenceLedgerExternalReviewOutcome struct {
	HandoffDigest  string `json:"handoff_digest"`
	Decision       string `json:"decision"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerExternalReviewOutcome(handoff DecisionConfidenceEvidenceLedgerReplayReviewHandoff, decision string) (DecisionConfidenceEvidenceLedgerExternalReviewOutcome, error) {
	digest, err := Digest(handoff)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerExternalReviewOutcome{}, err
	}
	outcome := DecisionConfidenceEvidenceLedgerExternalReviewOutcome{
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
