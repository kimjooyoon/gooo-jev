package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReplayReviewSignalPreservesBoundaries(t *testing.T) {
	observation := DecisionConfidenceEvidenceLedgerReplayOutcomeObservation{
		DispositionDigest: "disposition-digest",
		Outcome:           "diverged",
		Status:            "diverged",
		NonAuthorizing:    true,
	}
	signal, err := DeriveDecisionConfidenceEvidenceLedgerReplayReviewSignal(observation)
	if err != nil {
		t.Fatalf("derive review signal: %v", err)
	}
	if signal.Status != "review-required" || signal.ObservationDigest == "" || !signal.NonAuthorizing {
		t.Fatalf("unexpected review signal: %#v", signal)
	}

	observation.Status = "reproduced"
	observation.Outcome = "reproduced"
	signal, err = DeriveDecisionConfidenceEvidenceLedgerReplayReviewSignal(observation)
	if err != nil {
		t.Fatalf("derive observation signal: %v", err)
	}
	if signal.Status != "observation-only" || !signal.NonAuthorizing {
		t.Fatalf("reproduced observation was promoted: %#v", signal)
	}
}
