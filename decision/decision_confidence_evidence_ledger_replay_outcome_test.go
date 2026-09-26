package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReplayOutcomeIsObservational(t *testing.T) {
	disposition := DecisionConfidenceEvidenceLedgerReplayDisposition{
		ObservationDigest: "observation-digest",
		Status:            "ready-for-external-replay",
		NonAuthorizing:    true,
	}
	observation, err := ObserveDecisionConfidenceEvidenceLedgerReplayOutcome(disposition, "reproduced")
	if err != nil {
		t.Fatalf("observe replay outcome: %v", err)
	}
	if observation.Status != "reproduced" || observation.Outcome != "reproduced" || !observation.NonAuthorizing {
		t.Fatalf("unexpected replay outcome: %#v", observation)
	}

	observation, err = ObserveDecisionConfidenceEvidenceLedgerReplayOutcome(disposition, "improved")
	if err != nil {
		t.Fatalf("observe invalid outcome: %v", err)
	}
	if observation.Status != "unknown" || !observation.NonAuthorizing {
		t.Fatalf("invalid outcome was promoted: %#v", observation)
	}

	disposition.Status = "hold"
	observation, err = ObserveDecisionConfidenceEvidenceLedgerReplayOutcome(disposition, "reproduced")
	if err != nil {
		t.Fatalf("observe held outcome: %v", err)
	}
	if observation.Status != "unknown" || !observation.NonAuthorizing {
		t.Fatalf("held replay escaped unknown: %#v", observation)
	}
}
