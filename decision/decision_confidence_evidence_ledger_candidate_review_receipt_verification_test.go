package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerCandidateReviewReceiptVerificationRejectsTampering(t *testing.T) {
	signal := DecisionConfidenceEvidenceLedgerCandidateReviewSignal{ReviewOutcomeDigest: "signal-digest", Status: "ready-for-candidate-review", NonAuthorizing: true}
	outcome := DecisionConfidenceEvidenceLedgerExternalReviewOutcome{HandoffDigest: "review-digest", Status: "accepted-for-analysis", NonAuthorizing: true}
	receipt, err := BuildDecisionConfidenceEvidenceLedgerCandidateReviewReceipt(signal, outcome)
	if err != nil {
		t.Fatalf("build receipt: %v", err)
	}
	verification := VerifyDecisionConfidenceEvidenceLedgerCandidateReviewReceipt(receipt)
	if verification.Status != "verified" || !verification.NonAuthorizing {
		t.Fatalf("receipt was not verified: %#v", verification)
	}

	receipt.ReviewOutcomeDigest = "tampered"
	verification = VerifyDecisionConfidenceEvidenceLedgerCandidateReviewReceipt(receipt)
	if verification.Status != "mismatch" || !verification.NonAuthorizing {
		t.Fatalf("tampered receipt was not rejected: %#v", verification)
	}

	receipt = DecisionConfidenceEvidenceLedgerCandidateReviewReceipt{ReceiptDigest: "digest", NonAuthorizing: false}
	verification = VerifyDecisionConfidenceEvidenceLedgerCandidateReviewReceipt(receipt)
	if verification.Status != "unknown" || !verification.NonAuthorizing {
		t.Fatalf("authorizing receipt escaped unknown: %#v", verification)
	}
}
