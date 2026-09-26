package decision

import (
	"testing"
	"time"
)

func TestReceiptVerifyRejectsTampering(t *testing.T) {
	spec := Spec{
		ID:             "deployment-review",
		Question:       "Is deployment evidence complete?",
		Kind:           KindChoice,
		AllowedChoices: []string{"yes", "no"},
		PolicyDigest:   "policy-3",
	}
	result := Result{
		SpecID:         spec.ID,
		Kind:           spec.Kind,
		Value:          Value{Choice: "yes"},
		EvidenceDigest: "evidence-3",
		Provider:       "offline-fixture",
		Status:         StatusObserved,
		ObservedAt:     time.Unix(50, 0).UTC(),
	}
	state := State{Digest: "state-5"}
	receipt, err := Observe(spec, state, result)
	if err != nil {
		t.Fatalf("Observe() error = %v", err)
	}
	if err := receipt.Verify(spec, state, result); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	receipt.ResultDigest = "tampered"
	if err := receipt.Verify(spec, state, result); err == nil {
		t.Fatal("Verify() accepted a tampered result digest")
	}
}
