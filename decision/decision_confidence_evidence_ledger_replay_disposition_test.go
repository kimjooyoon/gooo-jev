package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReplayDispositionIsNonExecuting(t *testing.T) {
	replayable := DecisionConfidenceEvidenceLedgerReplayObservation{
		TransitionDigest:       "transition-digest",
		ExecutionReceiptDigest: "receipt-digest",
		Status:                 "replayable",
		NonAuthorizing:         true,
	}
	disposition, err := DeriveDecisionConfidenceEvidenceLedgerReplayDisposition(replayable)
	if err != nil {
		t.Fatalf("derive replay disposition: %v", err)
	}
	if disposition.Status != "ready-for-external-replay" || disposition.ObservationDigest == "" || !disposition.NonAuthorizing {
		t.Fatalf("unexpected replay disposition: %#v", disposition)
	}

	unknown := replayable
	unknown.Status = "unknown"
	disposition, err = DeriveDecisionConfidenceEvidenceLedgerReplayDisposition(unknown)
	if err != nil {
		t.Fatalf("derive hold disposition: %v", err)
	}
	if disposition.Status != "hold" || !disposition.NonAuthorizing {
		t.Fatalf("unknown observation escaped hold: %#v", disposition)
	}
}
