package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReplayObservationIsNonExecuting(t *testing.T) {
	first, err := NewDecisionConfidenceEvidenceLedgerEntry(1, "decision", "source-a", "")
	if err != nil {
		t.Fatalf("create first entry: %v", err)
	}
	second, err := NewDecisionConfidenceEvidenceLedgerEntry(2, "generation", "source-b", first.EntryDigest)
	if err != nil {
		t.Fatalf("create second entry: %v", err)
	}
	transition := ObserveDecisionConfidenceEvidenceLedgerTransition(first, second)
	link, err := LinkDecisionConfidenceEvidenceLedgerTransitionToExecutionReceipt(transition, ExecutionReceipt{})
	if err != nil {
		t.Fatalf("link receipt: %v", err)
	}
	observation, err := ObserveDecisionConfidenceEvidenceLedgerReplay(transition, link)
	if err != nil {
		t.Fatalf("observe replay: %v", err)
	}
	if observation.Status != "replayable" || observation.TransitionDigest == "" || !observation.NonAuthorizing {
		t.Fatalf("unexpected replay observation: %#v", observation)
	}

	link.Status = "unknown"
	observation, err = ObserveDecisionConfidenceEvidenceLedgerReplay(transition, link)
	if err != nil {
		t.Fatalf("observe unknown replay: %v", err)
	}
	if observation.Status != "unknown" || !observation.NonAuthorizing {
		t.Fatalf("invalid link was promoted: %#v", observation)
	}
}
