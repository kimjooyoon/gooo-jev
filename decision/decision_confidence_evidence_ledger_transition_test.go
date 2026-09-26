package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerTransitionObservesContinuity(t *testing.T) {
	first, err := NewDecisionConfidenceEvidenceLedgerEntry(1, "decision", "source-a", "")
	if err != nil {
		t.Fatalf("create first entry: %v", err)
	}
	second, err := NewDecisionConfidenceEvidenceLedgerEntry(2, "generation", "source-b", first.EntryDigest)
	if err != nil {
		t.Fatalf("create second entry: %v", err)
	}
	observed := ObserveDecisionConfidenceEvidenceLedgerTransition(first, second)
	if observed.Status != "extended" || observed.SequenceDelta != 1 || !observed.NonAuthorizing {
		t.Fatalf("unexpected extended transition: %#v", observed)
	}

	broken := second
	broken.PreviousEntryDigest = "other-predecessor"
	observed = ObserveDecisionConfidenceEvidenceLedgerTransition(first, broken)
	if observed.Status != "unknown" || observed.SequenceDelta != 1 || !observed.NonAuthorizing {
		t.Fatalf("unexpected unknown transition: %#v", observed)
	}
}
