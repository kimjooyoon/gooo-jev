package decision

import (
	"testing"
	"time"
)

func TestCapabilityGrantRejectsMismatchedWorkloadIdentity(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	boundary, err := NewCapabilityBoundary(
		"spiffe://example.org/ns/prod/sa/jev",
		"decision-review",
		"decision:review",
		"evidence-digest-1",
		now.Add(-time.Minute),
		now.Add(time.Hour),
	)
	if err != nil {
		t.Fatalf("NewCapabilityBoundary() error = %v", err)
	}
	identity, err := NewWorkloadIdentityObservation(
		"spiffe://example.org/ns/prod/sa/other",
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
	if _, err := NewCapabilityGrant(request, boundary, identity, "grant-evidence-1"); err == nil {
		t.Fatal("mismatched workload identity unexpectedly granted")
	}
}

