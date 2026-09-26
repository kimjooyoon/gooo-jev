package decision

// DecisionConfidenceEvidenceLedgerCandidateReviewMetricComparison compares two
// observations without calling either direction an improvement.
type DecisionConfidenceEvidenceLedgerCandidateReviewMetricComparison struct {
	PreviousSignalDigest string `json:"previous_signal_digest"`
	CurrentSignalDigest  string `json:"current_signal_digest"`
	Delta                string `json:"delta"`
	NonAuthorizing       bool   `json:"non_authorizing"`
}

func CompareDecisionConfidenceEvidenceLedgerCandidateReviewMetric(previous, current DecisionConfidenceEvidenceLedgerCandidateReviewMetric) DecisionConfidenceEvidenceLedgerCandidateReviewMetricComparison {
	comparison := DecisionConfidenceEvidenceLedgerCandidateReviewMetricComparison{
		PreviousSignalDigest: previous.SignalDigest,
		CurrentSignalDigest:  current.SignalDigest,
		Delta:                "unknown",
		NonAuthorizing:       true,
	}
	if previous.SignalDigest == "" || current.SignalDigest == "" || !previous.NonAuthorizing || !current.NonAuthorizing {
		return comparison
	}
	previousTotal := previous.ReadyCount + previous.RejectedCount + previous.HoldCount
	currentTotal := current.ReadyCount + current.RejectedCount + current.HoldCount
	if previousTotal != 1 || currentTotal != 1 {
		comparison.Delta = "inconclusive"
		return comparison
	}
	switch {
	case current.ReadyCount > previous.ReadyCount:
		comparison.Delta = "increased"
	case current.ReadyCount < previous.ReadyCount:
		comparison.Delta = "declined"
	default:
		comparison.Delta = "unchanged"
	}
	return comparison
}
