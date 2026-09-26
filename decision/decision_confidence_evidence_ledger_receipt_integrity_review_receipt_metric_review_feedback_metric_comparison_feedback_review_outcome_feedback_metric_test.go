package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetricIsOneHot(t *testing.T) {
	ready := ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric(DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedback{
		OutcomeDigest:  "outcome",
		Status:         "analysis-ready",
		NonAuthorizing: true,
	})
	if ready.AnalysisReadyCount != 1 || ready.RejectedCount != 0 || ready.HoldCount != 0 || ready.Status != "analysis-ready" || !ready.NonAuthorizing {
		t.Fatalf("unexpected analysis-ready metric: %#v", ready)
	}

	rejected := ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric(DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedback{
		OutcomeDigest:  "outcome",
		Status:         "rejected",
		NonAuthorizing: true,
	})
	if rejected.AnalysisReadyCount != 0 || rejected.RejectedCount != 1 || rejected.HoldCount != 0 || rejected.Status != "rejected" || !rejected.NonAuthorizing {
		t.Fatalf("unexpected rejected metric: %#v", rejected)
	}

	hold := ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackMetric(DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedback{
		OutcomeDigest:  "outcome",
		Status:         "hold",
		NonAuthorizing: true,
	})
	if hold.AnalysisReadyCount != 0 || hold.RejectedCount != 0 || hold.HoldCount != 1 || hold.Status != "hold" || !hold.NonAuthorizing {
		t.Fatalf("unexpected hold metric: %#v", hold)
	}
}
