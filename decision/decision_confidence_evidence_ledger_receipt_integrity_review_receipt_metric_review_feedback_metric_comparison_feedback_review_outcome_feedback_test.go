package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedbackRemainsNonAuthorizing(t *testing.T) {
	outcome := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcome{
		HandoffDigest:  "handoff",
		Decision:       "accepted-for-analysis",
		Status:         "accepted-for-analysis",
		NonAuthorizing: true,
	}
	feedback, err := DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedback(outcome)
	if err != nil {
		t.Fatalf("derive accepted feedback: %v", err)
	}
	if feedback.Status != "analysis-ready" || feedback.OutcomeDigest == "" || !feedback.NonAuthorizing {
		t.Fatalf("unexpected accepted feedback: %#v", feedback)
	}

	outcome.Status = "rejected"
	outcome.Decision = "rejected"
	feedback, err = DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedback(outcome)
	if err != nil {
		t.Fatalf("derive rejected feedback: %v", err)
	}
	if feedback.Status != "rejected" || !feedback.NonAuthorizing {
		t.Fatalf("unexpected rejected feedback: %#v", feedback)
	}

	outcome.Status = "unknown"
	feedback, err = DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackReviewOutcomeFeedback(outcome)
	if err != nil {
		t.Fatalf("derive held feedback: %v", err)
	}
	if feedback.Status != "hold" || !feedback.NonAuthorizing {
		t.Fatalf("unknown outcome escaped hold: %#v", feedback)
	}
}

