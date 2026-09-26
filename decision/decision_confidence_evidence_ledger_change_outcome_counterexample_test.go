package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerChangeOutcomeCounterexamplePreservesRollback(t *testing.T) {
	feedback := DecisionConfidenceEvidenceLedgerChangeOutcomeFeedback{
		PlanDigest:     "plan-digest",
		Delta:          "rolled-back",
		Status:         "review-required",
		NonAuthorizing: true,
	}
	counterexample, err := ObserveDecisionConfidenceEvidenceLedgerChangeOutcomeCounterexample(feedback)
	if err != nil {
		t.Fatalf("observe rollback counterexample: %v", err)
	}
	if counterexample.Status != "counterexample" || counterexample.FeedbackDigest == "" || !counterexample.NonAuthorizing {
		t.Fatalf("unexpected rollback counterexample: %#v", counterexample)
	}

	feedback.Status = "unknown"
	counterexample, err = ObserveDecisionConfidenceEvidenceLedgerChangeOutcomeCounterexample(feedback)
	if err != nil {
		t.Fatalf("observe held counterexample: %v", err)
	}
	if counterexample.Status != "unknown" || !counterexample.NonAuthorizing {
		t.Fatalf("unknown feedback escaped hold: %#v", counterexample)
	}
}
