package decision

import (
	"testing"
	"time"
)

func TestExecutionReceiptRejectsGrantBindingMismatch(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	firstBoundary, err := NewCapabilityBoundary(
		"spiffe://example.org/ns/prod/sa/jev",
		"decision-review",
		"decision:review",
		"evidence-digest-1",
		now.Add(-time.Minute),
		now.Add(time.Hour),
	)
	if err != nil {
		t.Fatalf("NewCapabilityBoundary(first) error = %v", err)
	}
	secondBoundary, err := NewCapabilityBoundary(
		"spiffe://example.org/ns/prod/sa/jev",
		"decision-review",
		"decision:review",
		"evidence-digest-2",
		now.Add(-time.Minute),
		now.Add(time.Hour),
	)
	if err != nil {
		t.Fatalf("NewCapabilityBoundary(second) error = %v", err)
	}
	identity, err := NewWorkloadIdentityObservation(
		"spiffe://example.org/ns/prod/sa/jev",
		"identity-evidence-1",
		WorkloadIdentityObserved,
		now,
	)
	if err != nil {
		t.Fatalf("NewWorkloadIdentityObservation() error = %v", err)
	}
	request := CapabilityRequest{
		Subject:    "spiffe://example.org/ns/prod/sa/jev",
		Audience:   "decision-review",
		Capability: "decision:review",
	}
	grant, err := NewCapabilityGrant(request, firstBoundary, identity, "grant-evidence-1")
	if err != nil {
		t.Fatalf("NewCapabilityGrant() error = %v", err)
	}
	receipt, err := NewExecutionReceipt(
		grant,
		secondBoundary,
		identity,
		Receipt{},
		"result-digest-1",
		ExecutionCompleted,
		now,
	)
	if err != nil {
		t.Fatalf("NewExecutionReceipt() error = %v", err)
	}
	if receipt.Status != ExecutionUnknown || receipt.UnknownReason != "grant-binding-mismatch" {
		t.Fatalf("binding mismatch = %#v", receipt)
	}
}

