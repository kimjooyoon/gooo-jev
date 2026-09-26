package decision

// DecisionConfidenceEvidenceLedgerChangeOutcome records an externally observed
// change outcome without interpreting application as improvement.
type DecisionConfidenceEvidenceLedgerChangeOutcome struct {
	PlanDigest      string `json:"plan_digest"`
	Outcome         string `json:"outcome"`
	Status          string `json:"status"`
	NonAuthorizing  bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerChangeOutcome(plan DecisionConfidenceEvidenceLedgerChangePlanObservation, outcome string) (DecisionConfidenceEvidenceLedgerChangeOutcome, error) {
	observation := DecisionConfidenceEvidenceLedgerChangeOutcome{
		PlanDigest:     plan.PlanDigest,
		Outcome:        outcome,
		Status:         "unknown",
		NonAuthorizing: true,
	}
	if plan.Status != "plan-observed" || plan.PlanDigest == "" {
		return observation, nil
	}
	switch outcome {
	case "applied", "aborted", "rolled-back", "not-applied":
		observation.Status = outcome
	}
	return observation, nil
}
