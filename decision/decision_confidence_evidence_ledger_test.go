package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerIsAppendOnlyAndNonAuthorizing(t *testing.T) {
	first, err := NewDecisionConfidenceEvidenceLedgerEntry(1, "decision", "source-a", "")
	if err != nil {
		t.Fatalf("create first entry: %v", err)
	}
	ledger, err := AppendDecisionConfidenceEvidenceLedgerEntry(nil, first)
	if err != nil {
		t.Fatalf("append first entry: %v", err)
	}
	second, err := NewDecisionConfidenceEvidenceLedgerEntry(2, "generation", "source-b", first.EntryDigest)
	if err != nil {
		t.Fatalf("create second entry: %v", err)
	}
	ledger, err = AppendDecisionConfidenceEvidenceLedgerEntry(ledger, second)
	if err != nil {
		t.Fatalf("append second entry: %v", err)
	}
	if len(ledger) != 2 || ledger[1].PreviousEntryDigest != ledger[0].EntryDigest {
		t.Fatalf("ledger chain was not preserved: %#v", ledger)
	}
	if ledger[0].EntryDigest == ledger[1].EntryDigest {
		t.Fatal("distinct entries unexpectedly share a digest")
	}
}

func TestDecisionConfidenceEvidenceLedgerRejectsTamperingAndGaps(t *testing.T) {
	first, err := NewDecisionConfidenceEvidenceLedgerEntry(1, "decision", "source-a", "")
	if err != nil {
		t.Fatalf("create first entry: %v", err)
	}
	tampered := first
	tampered.SourceDigest = "source-tampered"
	if err := VerifyDecisionConfidenceEvidenceLedgerEntry(tampered); err == nil {
		t.Fatal("tampered source digest was accepted")
	}
	gap, err := NewDecisionConfidenceEvidenceLedgerEntry(3, "generation", "source-c", first.EntryDigest)
	if err != nil {
		t.Fatalf("create gap entry: %v", err)
	}
	if _, err := AppendDecisionConfidenceEvidenceLedgerEntry([]DecisionConfidenceEvidenceLedgerEntry{first}, gap); err == nil {
		t.Fatal("sequence gap was accepted")
	}
	unauthorized := first
	unauthorized.NonAuthorizing = false
	if err := VerifyDecisionConfidenceEvidenceLedgerEntry(unauthorized); err == nil {
		t.Fatal("authorizing ledger entry was accepted")
	}
}
