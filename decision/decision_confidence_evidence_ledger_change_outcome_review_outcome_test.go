package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerChangeOutcomeReviewOutcomeDoesNotAuthorize(t *testing.T) {
	handoff := DecisionConfidenceEvidenceLedgerChangeOutcomeReviewHandoff{
		CounterexampleDigest: "counterexample-digest",
		Status:              "ready-for-external-review",
		NonAuthorizing:      true,
	}
	outcome, err := ObserveDecisionConfidenceEvidenceLedgerChangeOutcomeReviewOutcome(handoff, "accepted-for-analysis")
	if err != nil {
		t.Fatalf("observe accepted review: %v", err)
	}
	if outcome.Status != "accepted-for-analysis" || outcome.HandoffDigest == "" || !outcome.NonAuthorizing {
		t.Fatalf("unexpected accepted review: %#v", outcome)
	}

	outcome, err = ObserveDecisionConfidenceEvidenceLedgerChangeOutcomeReviewOutcome(handoff, "approved")
	if err != nil {
		t.Fatalf("observe invalid review: %v", err)
	}
	if outcome.Status != "unknown" || !outcome.NonAuthorizing {
		t.Fatalf("invalid review was promoted: %#v", outcome)
	}

	handoff.Status = "hold"
	outcome, err = ObserveDecisionConfidenceEvidenceLedgerChangeOutcomeReviewOutcome(handoff, "rejected")
	if err != nil {
		t.Fatalf("observe held review: %v", err)
	}
	if outcome.Status != "unknown" || !outcome.NonAuthorizing {
		t.Fatalf("held review escaped unknown: %#v", outcome)
	}
}
