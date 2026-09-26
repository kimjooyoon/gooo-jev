package decision

// DecisionConfidenceEvidenceLedgerChangeOutcomeComparison compares external
// change outcomes without inferring improvement.
type DecisionConfidenceEvidenceLedgerChangeOutcomeComparison struct {
	PreviousPlanDigest string `json:"previous_plan_digest"`
	CurrentPlanDigest  string `json:"current_plan_digest"`
	Delta              string `json:"delta"`
	NonAuthorizing     bool   `json:"non_authorizing"`
}

func CompareDecisionConfidenceEvidenceLedgerChangeOutcome(previous, current DecisionConfidenceEvidenceLedgerChangeOutcome) DecisionConfidenceEvidenceLedgerChangeOutcomeComparison {
	comparison := DecisionConfidenceEvidenceLedgerChangeOutcomeComparison{
		PreviousPlanDigest: previous.PlanDigest,
		CurrentPlanDigest:  current.PlanDigest,
		Delta:              "unknown",
		NonAuthorizing:     true,
	}
	if previous.PlanDigest == "" || current.PlanDigest == "" || previous.PlanDigest != current.PlanDigest || !previous.NonAuthorizing || !current.NonAuthorizing {
		return comparison
	}
	if previous.Status == current.Status {
		comparison.Delta = "unchanged"
		return comparison
	}
	if previous.Status == "applied" && current.Status == "rolled-back" {
		comparison.Delta = "rolled-back"
		return comparison
	}
	allowed := func(status string) bool {
		switch status {
		case "applied", "aborted", "rolled-back", "not-applied":
			return true
		default:
			return false
		}
	}
	if allowed(previous.Status) && allowed(current.Status) {
		comparison.Delta = "changed"
	}
	return comparison
}
