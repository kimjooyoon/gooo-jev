package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReceiptIntegrityCounterexamplePreservesChange(t *testing.T) {
	feedback := DecisionConfidenceEvidenceLedgerReceiptIntegrityFeedback{
		ComparisonDigest: "comparison-digest",
		Delta:            "changed",
		Status:           "review-required",
		NonAuthorizing:   true,
	}
	counterexample, err := ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityCounterexample(feedback)
	if err != nil {
		t.Fatalf("observe changed counterexample: %v", err)
	}
	if counterexample.Status != "counterexample" || counterexample.FeedbackDigest == "" || !counterexample.NonAuthorizing {
		t.Fatalf("unexpected changed counterexample: %#v", counterexample)
	}

	feedback.Status = "unknown"
	counterexample, err = ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityCounterexample(feedback)
	if err != nil {
		t.Fatalf("observe held counterexample: %v", err)
	}
	if counterexample.Status != "unknown" || !counterexample.NonAuthorizing {
		t.Fatalf("unknown feedback escaped hold: %#v", counterexample)
	}
}
