package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptIsNonAuthorizing(t *testing.T) {
	handoff := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewHandoff{CounterexampleDigest: "counterexample", Status: "ready-for-external-review", NonAuthorizing: true}
	outcome := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewOutcome{HandoffDigest: "handoff", Status: "accepted-for-analysis", NonAuthorizing: true}
	receipt, err := BuildDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceipt(handoff, outcome)
	if err != nil {
		t.Fatalf("build integrity review receipt: %v", err)
	}
	if receipt.Status != "ready-for-integrity-analysis" || receipt.ReceiptDigest == "" || !receipt.NonAuthorizing {
		t.Fatalf("unexpected integrity receipt: %#v", receipt)
	}

	outcome.Status = "unknown"
	receipt, err = BuildDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceipt(handoff, outcome)
	if err != nil {
		t.Fatalf("build held integrity receipt: %v", err)
	}
	if receipt.Status != "hold" || !receipt.NonAuthorizing {
		t.Fatalf("unknown integrity review escaped hold: %#v", receipt)
	}
}
