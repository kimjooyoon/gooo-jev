package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerCandidateReviewReceiptIsNonAuthorizing(t *testing.T) {
	signal := DecisionConfidenceEvidenceLedgerCandidateReviewSignal{
		ReviewOutcomeDigest: "signal-digest",
		Status:              "ready-for-candidate-review",
		NonAuthorizing:      true,
	}
	outcome := DecisionConfidenceEvidenceLedgerExternalReviewOutcome{
		HandoffDigest:  "review-digest",
		Status:         "accepted-for-analysis",
		NonAuthorizing: true,
	}
	receipt, err := BuildDecisionConfidenceEvidenceLedgerCandidateReviewReceipt(signal, outcome)
	if err != nil {
		t.Fatalf("build review receipt: %v", err)
	}
	if receipt.Status != "ready-for-candidate-review" || receipt.ReceiptDigest == "" || !receipt.NonAuthorizing {
		t.Fatalf("unexpected review receipt: %#v", receipt)
	}

	outcome.Status = "unknown"
	receipt, err = BuildDecisionConfidenceEvidenceLedgerCandidateReviewReceipt(signal, outcome)
	if err != nil {
		t.Fatalf("build held receipt: %v", err)
	}
	if receipt.Status != "hold" || !receipt.NonAuthorizing {
		t.Fatalf("unknown review escaped hold: %#v", receipt)
	}
}
