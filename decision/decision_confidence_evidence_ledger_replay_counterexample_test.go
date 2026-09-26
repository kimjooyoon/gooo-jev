package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReplayCounterexamplePreservesReviewState(t *testing.T) {
	signal := DecisionConfidenceEvidenceLedgerReplayReviewSignal{
		ObservationDigest: "observation-digest",
		Outcome:           "rolled-back",
		Status:            "review-required",
		NonAuthorizing:    true,
	}
	counterexample, err := ObserveDecisionConfidenceEvidenceLedgerReplayCounterexample(signal)
	if err != nil {
		t.Fatalf("observe counterexample: %v", err)
	}
	if counterexample.Status != "counterexample" || counterexample.SignalDigest == "" || !counterexample.NonAuthorizing {
		t.Fatalf("unexpected counterexample: %#v", counterexample)
	}

	signal.Status = "observation-only"
	signal.Outcome = "reproduced"
	counterexample, err = ObserveDecisionConfidenceEvidenceLedgerReplayCounterexample(signal)
	if err != nil {
		t.Fatalf("observe non-counterexample: %v", err)
	}
	if counterexample.Status != "no-counterexample" || !counterexample.NonAuthorizing {
		t.Fatalf("observation-only signal was promoted: %#v", counterexample)
	}
}
