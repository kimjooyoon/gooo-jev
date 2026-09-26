package decision

// DecisionConfidenceEvidenceLedgerChangePlanObservation joins review and write-set
// evidence without creating or applying a change plan.
type DecisionConfidenceEvidenceLedgerChangePlanObservation struct {
	ReviewOutcomeDigest string `json:"review_outcome_digest"`
	WriteSetDigest      string `json:"write_set_digest"`
	PlanDigest          string `json:"plan_digest"`
	Status              string `json:"status"`
	NonAuthorizing      bool   `json:"non_authorizing"`
}

type decisionConfidenceEvidenceLedgerChangePlanIdentity struct {
	ReviewOutcomeDigest string `json:"review_outcome_digest"`
	WriteSetDigest      string `json:"write_set_digest"`
}

func ObserveDecisionConfidenceEvidenceLedgerChangePlan(review DecisionConfidenceEvidenceLedgerExternalApplyReviewOutcome, writeSet DecisionConfidenceEvidenceLedgerWriteSetObservation) (DecisionConfidenceEvidenceLedgerChangePlanObservation, error) {
	identity := decisionConfidenceEvidenceLedgerChangePlanIdentity{
		ReviewOutcomeDigest: review.DispositionDigest,
		WriteSetDigest:      writeSet.WriteSetDigest,
	}
	planDigest, err := Digest(identity)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerChangePlanObservation{}, err
	}
	observation := DecisionConfidenceEvidenceLedgerChangePlanObservation{
		ReviewOutcomeDigest: review.DispositionDigest,
		WriteSetDigest:      writeSet.WriteSetDigest,
		PlanDigest:          planDigest,
		Status:              "unknown",
		NonAuthorizing:      true,
	}
	if review.Status == "accepted-for-apply-review" && writeSet.Status == "write-set-observed" {
		observation.Status = "plan-observed"
	} else if review.Status == "rejected" || writeSet.Status == "mismatch" {
		observation.Status = "rejected"
	}
	return observation, nil
}
