package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptVerificationRejectsTampering(t *testing.T) {
	handoff := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewHandoff{CounterexampleDigest: "counterexample", Status: "ready-for-external-review", NonAuthorizing: true}
	outcome := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewOutcome{HandoffDigest: "handoff", Status: "accepted-for-analysis", NonAuthorizing: true}
	receipt, err := BuildDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceipt(handoff, outcome)
	if err != nil {
		t.Fatalf("build receipt: %v", err)
	}
	verification := VerifyDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceipt(receipt)
	if verification.Status != "verified" || !verification.NonAuthorizing {
		t.Fatalf("receipt was not verified: %#v", verification)
	}

	receipt.OutcomeDigest = "tampered"
	verification = VerifyDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceipt(receipt)
	if verification.Status != "mismatch" || !verification.NonAuthorizing {
		t.Fatalf("tampered receipt was not rejected: %#v", verification)
	}

	receipt = DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceipt{ReceiptDigest: "digest", NonAuthorizing: false}
	verification = VerifyDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceipt(receipt)
	if verification.Status != "unknown" || !verification.NonAuthorizing {
		t.Fatalf("authorizing receipt escaped unknown: %#v", verification)
	}
}
