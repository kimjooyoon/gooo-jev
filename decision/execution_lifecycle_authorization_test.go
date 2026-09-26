package decision

import (
	"testing"
	"time"
)

func TestExecutionLifecycleAuthorizationBindsReplayAndIdentity(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	grant, boundary, identity, decisionReceipt := lifecycleInputs(t, now)
	suspension, resumption, reverse := lifecycleReplayInputs(t, now)
	replay, err := NewExecutionLifecycleReplay(suspension, resumption, reverse, now.Add(3*time.Minute))
	if err != nil {
		t.Fatalf("NewExecutionLifecycleReplay() error = %v", err)
	}
	authorization, err := NewExecutionLifecycleAuthorization(
		replay,
		resumption,
		grant,
		boundary,
		identity,
		decisionReceipt,
		now.Add(4*time.Minute),
	)
	if err != nil {
		t.Fatalf("NewExecutionLifecycleAuthorization() error = %v", err)
	}
	if authorization.Status != ExecutionLifecycleAuthorized {
		t.Fatalf("authorization status = %q", authorization.Status)
	}
	if err := authorization.Validate(); err != nil {
		t.Fatalf("authorization Validate() error = %v", err)
	}
}

func TestExecutionLifecycleAuthorizationPreservesUnknownStage(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	unknown, err := NewExecutionLifecycleAuthorization(
		ExecutionLifecycleReplayObservation{},
		ExecutionLifecycleReceipt{},
		CapabilityGrant{},
		CapabilityBoundary{},
		WorkloadIdentityObservation{},
		Receipt{},
		now,
	)
	if err != nil {
		t.Fatalf("NewExecutionLifecycleAuthorization() error = %v", err)
	}
	if unknown.Status != ExecutionLifecycleAuthorizationUnknown || unknown.MissingStage != "replay" {
		t.Fatalf("unknown authorization = %#v", unknown)
	}
	if err := unknown.Validate(); err != nil {
		t.Fatalf("unknown authorization Validate() error = %v", err)
	}
}

func TestExecutionLifecycleAuthorizationRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	grant, boundary, identity, decisionReceipt := lifecycleInputs(t, now)
	suspension, resumption, reverse := lifecycleReplayInputs(t, now)
	replay, err := NewExecutionLifecycleReplay(suspension, resumption, reverse, now.Add(3*time.Minute))
	if err != nil {
		t.Fatalf("NewExecutionLifecycleReplay() error = %v", err)
	}
	authorization, err := NewExecutionLifecycleAuthorization(
		replay,
		resumption,
		grant,
		boundary,
		identity,
		decisionReceipt,
		now.Add(4*time.Minute),
	)
	if err != nil {
		t.Fatalf("NewExecutionLifecycleAuthorization() error = %v", err)
	}
	authorization.DecisionReceiptDigest = "tampered-decision-receipt"
	if err := authorization.Validate(); err == nil {
		t.Fatal("tampered execution lifecycle authorization unexpectedly validated")
	}
}
