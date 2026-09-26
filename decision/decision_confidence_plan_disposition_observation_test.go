package decision

import "testing"

func TestObserveDecisionConfidenceChangePlanDispositionFailsClosed(t *testing.T) {
	_, err := ObserveDecisionConfidenceChangePlanDisposition(
		DecisionConfidenceChangePlanDisposition{},
		DecisionConfidenceChangePlanObservedApplied,
		"evidence",
	)
	if err == nil {
		t.Fatal("expected invalid disposition to be rejected")
	}
}

func TestObserveDecisionConfidenceChangePlanDispositionRejectsAppliedAbort(t *testing.T) {
	disposition := DecisionConfidenceChangePlanDisposition{
		ChangePlanDigest:   "plan",
		VerificationDigest: "verification",
		Status:             DecisionConfidenceChangePlanAbort,
		NonExecuting:       true,
	}
	digest, err := Digest(disposition)
	if err != nil {
		t.Fatal(err)
	}
	disposition.DispositionDigest = digest
	if _, err := ObserveDecisionConfidenceChangePlanDisposition(disposition, DecisionConfidenceChangePlanObservedApplied, "evidence"); err == nil {
		t.Fatal("expected aborted disposition to reject applied observation")
	}
}

func TestDecisionConfidenceChangePlanDispositionObservationRejectsTampering(t *testing.T) {
	observation := DecisionConfidenceChangePlanDispositionObservation{
		DispositionDigest:  "disposition",
		Status:             DecisionConfidenceChangePlanObservedRolledBack,
		EvidenceDigest:     "evidence",
		NonAuthorizing:     true,
		ObservationDigest:  "tampered",
	}
	if err := observation.Validate(); err == nil {
		t.Fatal("expected tampered disposition observation to be rejected")
	}
}
