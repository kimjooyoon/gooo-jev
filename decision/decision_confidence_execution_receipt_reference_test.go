package decision

import "testing"

func TestReferenceDecisionConfidenceExecutionReceiptFailsClosed(t *testing.T) {
	_, err := ReferenceDecisionConfidenceExecutionReceipt(
		DecisionConfidenceChangePlanDispositionObservation{},
		ExecutionReceipt{},
	)
	if err == nil {
		t.Fatal("expected invalid disposition observation to be rejected")
	}
}

func TestDecisionConfidenceExecutionReceiptReferenceRejectsTampering(t *testing.T) {
	reference := DecisionConfidenceExecutionReceiptReference{
		OutcomeObservationDigest: "observation",
		ExecutionReceiptDigest:  "receipt",
		NonAuthorizing:          true,
		ReferenceDigest:         "tampered",
	}
	if err := reference.Validate(); err == nil {
		t.Fatal("expected tampered receipt reference to be rejected")
	}
}
