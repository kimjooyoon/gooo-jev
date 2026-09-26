package decision

// ExecutionEnvelopeFeedbackMetricComparison records a bounded change between
// two one-hot provenance feedback observations.
type ExecutionEnvelopeFeedbackMetricComparison struct {
	PreviousFeedbackDigest string `json:"previous_feedback_digest"`
	CurrentFeedbackDigest  string `json:"current_feedback_digest"`
	Delta                  string `json:"delta"`
	NonAuthorizing         bool   `json:"non_authorizing"`
}

// CompareExecutionEnvelopeFeedbackMetric compares analysis-ready counts only
// after both metrics are valid one-hot observations.
func CompareExecutionEnvelopeFeedbackMetric(previous, current ExecutionEnvelopeFeedbackMetric) ExecutionEnvelopeFeedbackMetricComparison {
	comparison := ExecutionEnvelopeFeedbackMetricComparison{
		PreviousFeedbackDigest: previous.FeedbackDigest,
		CurrentFeedbackDigest:  current.FeedbackDigest,
		Delta:                  "UNKNOWN",
		NonAuthorizing:         true,
	}
	if !previous.NonAuthorizing || !current.NonAuthorizing || previous.FeedbackDigest == "" || current.FeedbackDigest == "" {
		comparison.NonAuthorizing = false
		return comparison
	}
	previousTotal := previous.AnalysisReadyCount + previous.RejectedCount + previous.HoldCount + previous.ObservationOnlyCount
	currentTotal := current.AnalysisReadyCount + current.RejectedCount + current.HoldCount + current.ObservationOnlyCount
	if previousTotal != 1 || currentTotal != 1 {
		comparison.Delta = "inconclusive"
		return comparison
	}
	switch {
	case current.AnalysisReadyCount > previous.AnalysisReadyCount:
		comparison.Delta = "increased"
	case current.AnalysisReadyCount < previous.AnalysisReadyCount:
		comparison.Delta = "declined"
	default:
		comparison.Delta = "unchanged"
	}
	return comparison
}
