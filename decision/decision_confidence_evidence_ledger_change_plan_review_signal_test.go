package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerChangePlanReviewSignalIsNonExecuting(t *testing.T) {
	feedback := DecisionConfidenceEvidenceLedgerFeedbackSignal{
		ComparisonDigest: "comparison-digest",
		Delta:            "declined",
		Status:           "review-required",
		NonAuthorizing:   true,
	}
	signal, err := DeriveDecisionConfidenceEvidenceLedgerChangePlanReviewSignal(feedback)
	if err != nil {
		t.Fatalf("derive change-plan signal: %v", err)
	}
	if signal.Status != "ready-for-external-change-plan-review" || signal.FeedbackDigest == "" || !signal.NonAuthorizing {
		t.Fatalf("unexpected change-plan signal: %#v", signal)
	}

	feedback.Status = "unknown"
	signal, err = DeriveDecisionConfidenceEvidenceLedgerChangePlanReviewSignal(feedback)
	if err != nil {
		t.Fatalf("derive held signal: %v", err)
	}
	if signal.Status != "hold" || !signal.NonAuthorizing {
		t.Fatalf("unknown feedback escaped hold: %#v", signal)
	}
}
