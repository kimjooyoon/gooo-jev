package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedbackPreservesGuardrails(t *testing.T) {
	comparison := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparison{
		PreviousFeedbackDigest: "previous",
		CurrentFeedbackDigest:  "current",
		Delta:                  "declined",
		NonAuthorizing:         true,
	}
	feedback, err := DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedback(comparison)
	if err != nil {
		t.Fatalf("derive declined feedback: %v", err)
	}
	if feedback.Status != "review-required" || feedback.ComparisonDigest == "" || !feedback.NonAuthorizing {
		t.Fatalf("unexpected declined feedback: %#v", feedback)
	}

	comparison.Delta = "increased"
	feedback, err = DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedback(comparison)
	if err != nil {
		t.Fatalf("derive increased feedback: %v", err)
	}
	if feedback.Status != "observation-only" || !feedback.NonAuthorizing {
		t.Fatalf("increased state escaped observation-only: %#v", feedback)
	}

	comparison.Delta = "inconclusive"
	feedback, err = DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewFeedbackMetricComparisonFeedback(comparison)
	if err != nil {
		t.Fatalf("derive inconclusive feedback: %v", err)
	}
	if feedback.Status != "hold" || !feedback.NonAuthorizing {
		t.Fatalf("inconclusive state escaped hold: %#v", feedback)
	}
}
