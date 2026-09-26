package decision

import (
	"testing"
	"time"
)

func lifecycleInputs(t *testing.T, now time.Time) (CapabilityGrant, CapabilityBoundary, WorkloadIdentityObservation, Receipt) {
	t.Helper()
	boundary, err := NewCapabilityBoundary(
		"spiffe://example.org/ns/prod/sa/jev",
		"decision-review",
		"decision:review",
		"boundary-evidence-1",
		now.Add(-time.Minute),
		now.Add(time.Hour),
	)
	if err != nil {
		t.Fatalf("NewCapabilityBoundary() error = %v", err)
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
		Subject:    boundary.Subject,
		Audience:   boundary.Audience,
		Capability: boundary.Capability,
	}
	grant, err := NewCapabilityGrant(request, boundary, identity, "grant-evidence-1")
	if err != nil {
		t.Fatalf("NewCapabilityGrant() error = %v", err)
	}
	result, err := Observe(
		Spec{
			ID:             "review",
			Question:       "Should this candidate be reviewed?",
			Kind:           KindChoice,
			AllowedChoices: []string{"accept", "review"},
			PolicyDigest:   "policy-digest-1",
		},
		State{Digest: "state-digest-1"},
		Result{
			SpecID:          "review",
			Kind:            KindChoice,
			Value:           Value{Choice: "review"},
			EvidenceDigest:  "decision-evidence-1",
			Provider:        "provider-neutral",
			Status:          StatusObserved,
			ObservedAt:      now,
		},
	)
	if err != nil {
		t.Fatalf("Observe() error = %v", err)
	}
	return grant, boundary, identity, result
}

func TestExecutionLifecycleSuspendsAndResumesWithBinding(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	grant, boundary, identity, decisionReceipt := lifecycleInputs(t, now)
	suspension, err := NewExecutionSuspension(
		grant,
		boundary,
		identity,
		decisionReceipt,
		"checkpoint-digest-1",
		"state-digest-1",
		"suspension-evidence-1",
		now,
	)
	if err != nil {
		t.Fatalf("NewExecutionSuspension() error = %v", err)
	}
	if suspension.Status != ExecutionLifecycleSuspended {
		t.Fatalf("suspension status = %q", suspension.Status)
	}
	if err := suspension.Validate(); err != nil {
		t.Fatalf("suspension Validate() error = %v", err)
	}
	resumption, err := NewExecutionResumption(
		suspension,
		grant,
		boundary,
		identity,
		decisionReceipt,
		"state-digest-2",
		"resumption-evidence-1",
		now.Add(time.Minute),
	)
	if err != nil {
		t.Fatalf("NewExecutionResumption() error = %v", err)
	}
	if resumption.Status != ExecutionLifecycleResumed {
		t.Fatalf("resumption status = %q", resumption.Status)
	}
	if resumption.ParentLifecycleDigest != suspension.LifecycleDigest {
		t.Fatalf("parent lifecycle digest = %q, want %q", resumption.ParentLifecycleDigest, suspension.LifecycleDigest)
	}
	if err := resumption.Validate(); err != nil {
		t.Fatalf("resumption Validate() error = %v", err)
	}
}

func TestExecutionLifecycleFailsClosedWithUnknownStage(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	unknown, err := NewExecutionSuspension(
		CapabilityGrant{},
		CapabilityBoundary{},
		WorkloadIdentityObservation{},
		Receipt{},
		"",
		"",
		"",
		now,
	)
	if err != nil {
		t.Fatalf("NewExecutionSuspension() error = %v", err)
	}
	if unknown.Status != ExecutionLifecycleUnknown || unknown.UnknownReason != "missing-or-invalid-grant" {
		t.Fatalf("unknown lifecycle receipt = %#v", unknown)
	}
	if err := unknown.Validate(); err != nil {
		t.Fatalf("unknown Validate() error = %v", err)
	}
}

func TestExecutionLifecycleRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	grant, boundary, identity, decisionReceipt := lifecycleInputs(t, now)
	suspension, err := NewExecutionSuspension(
		grant,
		boundary,
		identity,
		decisionReceipt,
		"checkpoint-digest-1",
		"state-digest-1",
		"suspension-evidence-1",
		now,
	)
	if err != nil {
		t.Fatalf("NewExecutionSuspension() error = %v", err)
	}
	suspension.CheckpointDigest = "tampered-checkpoint"
	if err := suspension.Validate(); err == nil {
		t.Fatal("tampered suspension unexpectedly validated")
	}
}
