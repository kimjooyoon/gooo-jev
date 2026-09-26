package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonPreservesDirection(t *testing.T) {
	previous := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetric{
		FeedbackDigest:     "previous",
		AnalysisReadyCount: 0,
		RejectedCount:      1,
		NonAuthorizing:     true,
	}
	current := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetric{
		FeedbackDigest:     "current",
		AnalysisReadyCount: 1,
		NonAuthorizing:     true,
	}
	comparison := CompareDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetric(previous, current)
	if comparison.Delta != "increased" || !comparison.NonAuthorizing {
		t.Fatalf("review readiness increase was not preserved: %#v", comparison)
	}

	current.AnalysisReadyCount = 0
	current.RejectedCount = 1
	comparison = CompareDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetric(previous, current)
	if comparison.Delta != "unchanged" || !comparison.NonAuthorizing {
		t.Fatalf("unchanged review state was not preserved: %#v", comparison)
	}

	previous.AnalysisReadyCount = 1
	previous.RejectedCount = 0
	current.AnalysisReadyCount = 0
	current.HoldCount = 1
	comparison = CompareDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetric(previous, current)
	if comparison.Delta != "declined" || !comparison.NonAuthorizing {
		t.Fatalf("review readiness decline was not preserved: %#v", comparison)
	}

	current.HoldCount = 2
	comparison = CompareDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetric(previous, current)
	if comparison.Delta != "inconclusive" || !comparison.NonAuthorizing {
		t.Fatalf("invalid one-hot state escaped inconclusive: %#v", comparison)
	}
}
