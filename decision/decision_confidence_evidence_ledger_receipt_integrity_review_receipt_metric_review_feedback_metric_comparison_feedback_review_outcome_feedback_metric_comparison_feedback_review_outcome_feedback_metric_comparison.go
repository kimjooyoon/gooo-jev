package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparison
// compares review outcome feedback observations without inferring improvement.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparison struct {
	PreviousFeedbackDigest string `json:"previous_feedback_digest"`
	CurrentFeedbackDigest  string `json:"current_feedback_digest"`
	Delta                  string `json:"delta"`
	NonAuthorizing         bool   `json:"non_authorizing"`
}

func CompareDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric(previous, current DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric) DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparison {
	comparison := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparison{
		PreviousFeedbackDigest: previous.FeedbackDigest,
		CurrentFeedbackDigest:  current.FeedbackDigest,
		Delta:                  "unknown",
		NonAuthorizing:         true,
	}
	if previous.FeedbackDigest == "" || current.FeedbackDigest == "" || !previous.NonAuthorizing || !current.NonAuthorizing {
		return comparison
	}
	previousTotal := previous.AnalysisReadyCount + previous.RejectedCount + previous.HoldCount
	currentTotal := current.AnalysisReadyCount + current.RejectedCount + current.HoldCount
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
