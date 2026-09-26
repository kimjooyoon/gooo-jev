package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerWriteSetObservationIsNonAuthorizing(t *testing.T) {
	signal := DecisionConfidenceEvidenceLedgerChangePlanReviewSignal{
		FeedbackDigest: "feedback-digest",
		Status:         "ready-for-external-change-plan-review",
		NonAuthorizing: true,
	}
	observation := ObserveDecisionConfidenceEvidenceLedgerWriteSet(signal, "write-set-digest")
	if observation.Status != "write-set-observed" || observation.WriteSetDigest == "" || !observation.NonAuthorizing {
		t.Fatalf("unexpected write-set observation: %#v", observation)
	}

	observation = ObserveDecisionConfidenceEvidenceLedgerWriteSet(signal, "")
	if observation.Status != "unknown" || !observation.NonAuthorizing {
		t.Fatalf("missing write-set escaped unknown: %#v", observation)
	}

	signal.Status = "hold"
	observation = ObserveDecisionConfidenceEvidenceLedgerWriteSet(signal, "write-set-digest")
	if observation.Status != "unknown" || !observation.NonAuthorizing {
		t.Fatalf("held signal escaped unknown: %#v", observation)
	}
}
