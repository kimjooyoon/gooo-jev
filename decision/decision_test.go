package decision

import (
	"testing"
	"time"
)

func TestObserveProducesNonAuthorizingReceipt(t *testing.T) {
	score := 0.8
	spec := Spec{
		ID:           "risk-review",
		Question:     "Is the proposed change low risk?",
		Kind:         KindScore,
		Threshold:    &score,
		PolicyDigest: "policy-1",
	}
	result := Result{
		SpecID:         spec.ID,
		Kind:           spec.Kind,
		Value:          Value{Score: &score},
		Confidence:     &score,
		EvidenceDigest: "evidence-1",
		Provider:       "rule",
		Status:         StatusObserved,
		ObservedAt:     time.Unix(10, 0).UTC(),
	}
	receipt, err := Observe(spec, State{Digest: "state-1"}, result)
	if err != nil {
		t.Fatalf("Observe() error = %v", err)
	}
	if err := receipt.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if !receipt.NonAuthorizing {
		t.Fatal("receipt must not authorize execution")
	}
}

func TestObservePreservesUnknown(t *testing.T) {
	spec := Spec{
		ID:           "freshness",
		Question:     "Is the evidence fresh enough?",
		Kind:         KindNoul,
		PolicyDigest: "policy-2",
	}
	result := Result{
		SpecID:        spec.ID,
		Kind:          spec.Kind,
		Provider:      "jev",
		Status:        StatusUnknown,
		EvidenceDigest: "evidence-2",
		ObservedAt:    time.Unix(20, 0).UTC(),
	}
	receipt, err := Observe(spec, State{Digest: "state-2"}, result)
	if err != nil {
		t.Fatalf("Observe() error = %v", err)
	}
	if receipt.Status != StatusUnknown {
		t.Fatalf("status = %q, want %q", receipt.Status, StatusUnknown)
	}
}

func TestDigestChangesWithState(t *testing.T) {
	first, err := Digest(State{Digest: "state-a"})
	if err != nil {
		t.Fatalf("Digest() error = %v", err)
	}
	second, err := Digest(State{Digest: "state-b"})
	if err != nil {
		t.Fatalf("Digest() error = %v", err)
	}
	if first == second {
		t.Fatal("state digests must differ")
	}
}
