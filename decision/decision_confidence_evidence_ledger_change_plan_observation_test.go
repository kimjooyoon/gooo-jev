package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerChangePlanObservationJoinsEvidenceOnly(t *testing.T) {
	review := DecisionConfidenceEvidenceLedgerExternalApplyReviewOutcome{
		DispositionDigest: "disposition-digest",
		Decision:          "accepted-for-apply-review",
		Status:            "accepted-for-apply-review",
		NonAuthorizing:    true,
	}
	writeSet := DecisionConfidenceEvidenceLedgerWriteSetObservation{
		ReviewSignalDigest: "review-signal-digest",
		WriteSetDigest:     "write-set-digest",
		Status:             "write-set-observed",
		NonAuthorizing:     true,
	}
	observation, err := ObserveDecisionConfidenceEvidenceLedgerChangePlan(review, writeSet)
	if err != nil {
		t.Fatalf("observe change plan: %v", err)
	}
	if observation.Status != "plan-observed" || observation.PlanDigest == "" || !observation.NonAuthorizing {
		t.Fatalf("unexpected plan observation: %#v", observation)
	}

	writeSet.Status = "mismatch"
	observation, err = ObserveDecisionConfidenceEvidenceLedgerChangePlan(review, writeSet)
	if err != nil {
		t.Fatalf("observe rejected plan: %v", err)
	}
	if observation.Status != "rejected" || !observation.NonAuthorizing {
		t.Fatalf("mismatched write-set was not rejected: %#v", observation)
	}
}
