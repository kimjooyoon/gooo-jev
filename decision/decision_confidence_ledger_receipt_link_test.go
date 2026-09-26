package decision

import "testing"

func TestDecisionConfidenceLedgerReceiptLinkUsesDigestOnly(t *testing.T) {
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
	if link.Status != "linked" || link.ExecutionReceiptDigest == "" || !link.NonAuthorizing {
		t.Fatalf("unexpected receipt link: %#v", link)
	}

	unknown := transition
	unknown.Status = "unknown"
	link, err = LinkDecisionConfidenceEvidenceLedgerTransitionToExecutionReceipt(unknown, ExecutionReceipt{})
	if err != nil {
		t.Fatalf("observe unknown link: %v", err)
	}
	if link.Status != "unknown" || !link.NonAuthorizing {
		t.Fatalf("unknown transition was promoted: %#v", link)
	}
}
