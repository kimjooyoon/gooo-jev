package decision

import (
	"testing"
	"time"
)

func TestAppendVerifiedBuildsAndValidatesReceipt(t *testing.T) {
	spec := Spec{
		ID:             "verified-review",
		Question:       "Is this observation complete?",
		Kind:           KindChoice,
		AllowedChoices: []string{"yes", "no"},
		PolicyDigest:   "policy-verified",
	}
	result := Result{
		SpecID:         spec.ID,
		Kind:           spec.Kind,
		Value:          Value{Choice: "yes"},
		EvidenceDigest: "evidence-verified",
		Provider:       "rule",
		Status:         StatusObserved,
		ObservedAt:     time.Unix(70, 0).UTC(),
	}
	ledger, err := (Ledger{}).AppendVerified(spec, State{Digest: "state-verified"}, result)
	if err != nil {
		t.Fatalf("AppendVerified() error = %v", err)
	}
	if len(ledger.Entries) != 1 {
		t.Fatalf("entry count = %d, want 1", len(ledger.Entries))
	}
	if err := ledger.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}
