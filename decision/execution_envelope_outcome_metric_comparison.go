package decision

// ExecutionEnvelopeOutcomeMetricComparison records a bounded change between
// two one-hot review outcome observations.
type ExecutionEnvelopeOutcomeMetricComparison struct {
	PreviousOutcomeDigest string `json:"previous_outcome_digest"`
	CurrentOutcomeDigest  string `json:"current_outcome_digest"`
	Delta                 string `json:"delta"`
	NonAuthorizing        bool   `json:"non_authorizing"`
}

// CompareExecutionEnvelopeOutcomeMetric compares accepted counts only after
// both observations are complete, valid, and non-authorizing.
func CompareExecutionEnvelopeOutcomeMetric(previous, current ExecutionEnvelopeOutcomeMetric) ExecutionEnvelopeOutcomeMetricComparison {
	comparison := ExecutionEnvelopeOutcomeMetricComparison{
		PreviousOutcomeDigest: previous.OutcomeDigest,
		CurrentOutcomeDigest:  current.OutcomeDigest,
		Delta:                 "UNKNOWN",
		NonAuthorizing:        true,
	}
	if !previous.NonAuthorizing || !current.NonAuthorizing || previous.OutcomeDigest == "" || current.OutcomeDigest == "" {
		comparison.NonAuthorizing = false
		return comparison
	}
	previousTotal := previous.AcceptedCount + previous.RejectedCount + previous.HoldCount
	currentTotal := current.AcceptedCount + current.RejectedCount + current.HoldCount
	if previousTotal != 1 || currentTotal != 1 {
		comparison.Delta = "inconclusive"
		return comparison
	}
	switch {
	case current.AcceptedCount > previous.AcceptedCount:
		comparison.Delta = "increased"
	case current.AcceptedCount < previous.AcceptedCount:
		comparison.Delta = "declined"
	default:
		comparison.Delta = "unchanged"
	}
	return comparison
}
