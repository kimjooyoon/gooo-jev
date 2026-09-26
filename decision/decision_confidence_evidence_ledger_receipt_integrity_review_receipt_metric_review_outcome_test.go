package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewOutcomeDoesNotAuthorize(t *testing.T) {
	handoff := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewHandoff{
		CounterexampleDigest: "counterexample-digest",
		Status:              "ready-for-external-review",
		NonAuthorizing:      true,
	}
	outcome, err := ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewOutcome(handoff, "accepted-for-analysis")
	if err != nil {
		t.Fatalf("observe accepted integrity review: %v", err)
	}
	if outcome.Status != "accepted-for-analysis" || outcome.HandoffDigest == "" || !outcome.NonAuthorizing {
		t.Fatalf("unexpected accepted review: %#v", outcome)
	}

	outcome, err = ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewOutcome(handoff, "approved")
	if err != nil {
		t.Fatalf("observe invalid integrity review: %v", err)
	}
	if outcome.Status != "unknown" || !outcome.NonAuthorizing {
		t.Fatalf("invalid decision was promoted: %#v", outcome)
	}

	handoff.Status = "hold"
	outcome, err = ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricReviewOutcome(handoff, "rejected")
	if err != nil {
		t.Fatalf("observe held integrity review: %v", err)
	}
	if outcome.Status != "unknown" || !outcome.NonAuthorizing {
		t.Fatalf("held review escaped unknown: %#v", outcome)
	}
}
