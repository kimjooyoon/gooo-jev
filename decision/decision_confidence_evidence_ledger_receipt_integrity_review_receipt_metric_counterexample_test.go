package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricCounterexamplePreservesBoundary(t *testing.T) {
	feedback := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricFeedback{
		ComparisonDigest: "comparison-digest",
		Delta:            "changed",
		Status:           "review-required",
		NonAuthorizing:   true,
	}
	counterexample, err := ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricCounterexample(feedback)
	if err != nil {
		t.Fatalf("observe changed counterexample: %v", err)
	}
	if counterexample.Status != "counterexample" || counterexample.FeedbackDigest == "" || !counterexample.NonAuthorizing {
		t.Fatalf("unexpected changed counterexample: %#v", counterexample)
	}

	feedback.Status = "observation-only"
	counterexample, err = ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricCounterexample(feedback)
	if err != nil {
		t.Fatalf("observe unchanged counterexample: %v", err)
	}
	if counterexample.Status != "no-counterexample" || !counterexample.NonAuthorizing {
		t.Fatalf("unexpected unchanged counterexample: %#v", counterexample)
	}

	feedback.Status = "hold"
	counterexample, err = ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricCounterexample(feedback)
	if err != nil {
		t.Fatalf("observe held counterexample: %v", err)
	}
	if counterexample.Status != "unknown" || !counterexample.NonAuthorizing {
		t.Fatalf("held feedback escaped unknown: %#v", counterexample)
	}
}
