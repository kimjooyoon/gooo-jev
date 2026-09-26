package decision

import "testing"

func TestObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricPreservesOneHotState(t *testing.T) {
	feedback := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedback{
		OutcomeDigest: "analysis",
		Status:        "analysis-ready",
	}
	metric := ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric(feedback)
	if metric.FeedbackDigest != "analysis" || metric.AnalysisReadyCount != 1 || metric.Status != "analysis-ready" || !metric.NonAuthorizing {
		t.Fatalf("analysis-ready metric was not preserved: %#v", metric)
	}

	feedback.OutcomeDigest = "rejected"
	feedback.Status = "rejected"
	metric = ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric(feedback)
	if metric.RejectedCount != 1 || metric.AnalysisReadyCount != 0 || metric.Status != "rejected" || !metric.NonAuthorizing {
		t.Fatalf("rejected metric was not preserved: %#v", metric)
	}

	feedback.OutcomeDigest = "hold"
	feedback.Status = "unknown"
	metric = ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric(feedback)
	if metric.HoldCount != 1 || metric.AnalysisReadyCount != 0 || metric.RejectedCount != 0 || metric.Status != "hold" || !metric.NonAuthorizing {
		t.Fatalf("unknown metric did not remain hold: %#v", metric)
	}
}
