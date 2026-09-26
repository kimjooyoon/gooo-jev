package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerExternalApplyReviewOutcomeDoesNotAuthorize(t *testing.T) {
	disposition := DecisionConfidenceEvidenceLedgerWriteSetDisposition{
		IntegrityDigest: "integrity-digest",
		Status:         "ready-for-external-apply-review",
		NonAuthorizing: true,
	}
	outcome, err := ObserveDecisionConfidenceEvidenceLedgerExternalApplyReviewOutcome(disposition, "accepted-for-apply-review")
	if err != nil {
		t.Fatalf("observe accepted apply review: %v", err)
	}
	if outcome.Status != "accepted-for-apply-review" || outcome.DispositionDigest == "" || !outcome.NonAuthorizing {
		t.Fatalf("unexpected accepted review: %#v", outcome)
	}

	outcome, err = ObserveDecisionConfidenceEvidenceLedgerExternalApplyReviewOutcome(disposition, "approved")
	if err != nil {
		t.Fatalf("observe invalid apply review: %v", err)
	}
	if outcome.Status != "unknown" || !outcome.NonAuthorizing {
		t.Fatalf("invalid decision was promoted: %#v", outcome)
	}

	disposition.Status = "abort"
	outcome, err = ObserveDecisionConfidenceEvidenceLedgerExternalApplyReviewOutcome(disposition, "accepted-for-apply-review")
	if err != nil {
		t.Fatalf("observe aborted apply review: %v", err)
	}
	if outcome.Status != "unknown" || !outcome.NonAuthorizing {
		t.Fatalf("aborted disposition escaped unknown: %#v", outcome)
	}
}
