package decision

import (
	"strings"
	"testing"
	"time"
)

func TestLedgerAppendsAndRejectsTampering(t *testing.T) {
	digest := "evidence-ledger"
	spec := Spec{
		ID:           "ledger-review",
		Question:     "Should this observation be retained?",
		Kind:         KindChoice,
		AllowedChoices: []string{"yes", "no"},
		PolicyDigest: "policy-ledger",
	}
	result := Result{
		SpecID:         spec.ID,
		Kind:           spec.Kind,
		Value:          Value{Choice: "yes"},
		EvidenceDigest: digest,
		Provider:       "rule",
		Status:         StatusObserved,
		ObservedAt:     time.Unix(60, 0).UTC(),
	}
	receipt, err := Observe(spec, State{Digest: "state-ledger"}, result)
	if err != nil {
		t.Fatalf("Observe() error = %v", err)
	}
	ledger, err := (Ledger{}).Append(receipt)
	if err != nil {
		t.Fatalf("first Append() error = %v", err)
	}
	ledger, err = ledger.Append(receipt)
	if err != nil {
		t.Fatalf("second Append() error = %v", err)
	}
	if err := ledger.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if ledger.Digest() == "" || !strings.Contains(ledger.Digest(), "") {
		t.Fatal("ledger digest must be present")
	}
	ledger.Entries[1].PreviousDigest = "tampered"
	if err := ledger.Validate(); err == nil {
		t.Fatal("Validate() accepted a tampered chain link")
	}
}
