package decision

// DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric
// records review outcome feedback as a one-hot observation without measuring improvement.
type DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric struct {
	FeedbackDigest     string `json:"feedback_digest"`
	AnalysisReadyCount uint64 `json:"analysis_ready_count"`
	RejectedCount      uint64 `json:"rejected_count"`
	HoldCount          uint64 `json:"hold_count"`
	Status             string `json:"status"`
	NonAuthorizing     bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric(feedback DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedback) DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric {
	metric := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric{
		FeedbackDigest: feedback.OutcomeDigest,
		Status:         "hold",
		NonAuthorizing: true,
	}
	switch feedback.Status {
	case "analysis-ready":
		metric.AnalysisReadyCount = 1
		metric.Status = "analysis-ready"
	case "rejected":
		metric.RejectedCount = 1
		metric.Status = "rejected"
	default:
		metric.HoldCount = 1
	}
	return metric
}
