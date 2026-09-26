package decision

import (
    "testing"
    "time"
)

func TestCapabilityBoundaryLifecycleAndTamperDetection(t *testing.T) {
    issuedAt := time.Unix(1_800_000_000, 0).UTC()
    expiresAt := issuedAt.Add(time.Hour)
    boundary, err := NewCapabilityBoundary(
        "spiffe://example.org/ns/prod/sa/jev",
        "decision-review",
        "decision:review",
        "evidence-digest-1",
        issuedAt,
        expiresAt,
    )
    if err != nil {
        t.Fatalf("NewCapabilityBoundary() error = %v", err)
    }
    if err := boundary.Validate(); err != nil {
        t.Fatalf("Validate() error = %v", err)
    }
    if got := boundary.StatusAt(issuedAt.Add(-time.Nanosecond)); got != CapabilityBoundaryNotYet {
        t.Fatalf("future status = %q", got)
    }
    if got := boundary.StatusAt(issuedAt.Add(time.Minute)); got != CapabilityBoundaryActive {
        t.Fatalf("active status = %q", got)
    }
    if got := boundary.StatusAt(expiresAt); got != CapabilityBoundaryExpired {
        t.Fatalf("expired status = %q", got)
    }
    tampered := boundary
    tampered.Capability = "decision:admin"
    if err := tampered.Validate(); err == nil {
        t.Fatal("tampered capability unexpectedly validated")
    }
    tampered.BoundaryDigest = ""
    if got := tampered.StatusAt(issuedAt.Add(time.Minute)); got != CapabilityBoundaryUnknown {
        t.Fatalf("invalid status = %q", got)
    }
}