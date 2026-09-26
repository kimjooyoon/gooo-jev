package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricFeedbackPreservesReviewBoundary(t *testing.T) {
	feedback, err := DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricFeedback(DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricComparison{
		PreviousReceiptDigest: "previous",
		CurrentReceiptDigest:  "current",
		Delta:                 "changed",
		NonAuthorizing:        true,
	})
	if err != nil {
		t.Fatalf("derive changed feedback: %v", err)
	}
	if feedback.Status != "review-required" || feedback.ComparisonDigest == "" || !feedback.NonAuthorizing {
		t.Fatalf("unexpected changed feedback: %#v", feedback)
	}

	feedback, err = DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricFeedback(DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricComparison{
		PreviousReceiptDigest: "same",
		CurrentReceiptDigest:  "same",
		Delta:                 "unchanged",
		NonAuthorizing:        true,
	})
	if err != nil {
		t.Fatalf("derive unchanged feedback: %v", err)
	}
	if feedback.Status != "observation-only" || !feedback.NonAuthorizing {
		t.Fatalf("unexpected unchanged feedback: %#v", feedback)
	}

	feedback, err = DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricFeedback(DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricComparison{
		Delta:          "unknown",
		NonAuthorizing: true,
	})
	if err != nil {
		t.Fatalf("derive unknown feedback: %v", err)
	}
	if feedback.Status != "hold" || !feedback.NonAuthorizing {
		t.Fatalf("unknown integrity result escaped hold: %#v", feedback)
	}
}
