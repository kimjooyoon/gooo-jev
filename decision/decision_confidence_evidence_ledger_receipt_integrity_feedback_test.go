package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReceiptIntegrityFeedbackPreservesChange(t *testing.T) {
	feedback, err := DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityFeedback(DecisionConfidenceEvidenceLedgerCandidateReviewReceiptMetricComparison{
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

	feedback, err = DeriveDecisionConfidenceEvidenceLedgerReceiptIntegrityFeedback(DecisionConfidenceEvidenceLedgerCandidateReviewReceiptMetricComparison{
		PreviousReceiptDigest: "same",
		CurrentReceiptDigest:  "same",
		Delta:                 "unchanged",
		NonAuthorizing:        true,
	})
	if err != nil {
		t.Fatalf("derive unchanged feedback: %v", err)
	}
	if feedback.Status != "observation-only" || !feedback.NonAuthorizing {
		t.Fatalf("unchanged feedback was promoted: %#v", feedback)
	}
}
