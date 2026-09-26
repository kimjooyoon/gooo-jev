package decision

import "testing"

func TestCompareDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricComparisonPreservesDirection(t *testing.T) {
	previous := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric{
		FeedbackDigest: "previous",
		RejectedCount:  1,
		Status:         "rejected",
		NonAuthorizing: true,
	}
	current := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric{
		FeedbackDigest:     "current",
		AnalysisReadyCount: 1,
		Status:             "analysis-ready",
		NonAuthorizing:     true,
	}
	comparison := CompareDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric(previous, current)
	if comparison.Delta != "increased" || !comparison.NonAuthorizing {
		t.Fatalf("review readiness increase was not preserved: %#v", comparison)
	}

	current.AnalysisReadyCount = 0
	current.RejectedCount = 1
	current.Status = "rejected"
	comparison = CompareDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric(previous, current)
	if comparison.Delta != "unchanged" || !comparison.NonAuthorizing {
		t.Fatalf("unchanged review state was not preserved: %#v", comparison)
	}

	previous.AnalysisReadyCount = 1
	previous.RejectedCount = 0
	previous.Status = "analysis-ready"
	current.AnalysisReadyCount = 0
	current.RejectedCount = 0
	current.HoldCount = 1
	current.Status = "hold"
	comparison = CompareDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric(previous, current)
	if comparison.Delta != "declined" || !comparison.NonAuthorizing {
		t.Fatalf("review readiness decline was not preserved: %#v", comparison)
	}

	current.HoldCount = 2
	comparison = CompareDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric(previous, current)
	if comparison.Delta != "inconclusive" || !comparison.NonAuthorizing {
		t.Fatalf("invalid one-hot state escaped inconclusive: %#v", comparison)
	}
}
