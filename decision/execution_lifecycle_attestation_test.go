package decision

import (
	"testing"
	"time"
)

func TestExecutionLifecycleAttestationBindsObservation(t *testing.T) {
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
	attestation, err := NewExecutionLifecycleAttestation(authorization, "reverse-observation-digest", now.Add(5*time.Minute))
	if err != nil {
		t.Fatalf("NewExecutionLifecycleAttestation() error = %v", err)
	}
	if attestation.Status != ExecutionLifecycleAttested {
		t.Fatalf("attestation status = %q", attestation.Status)
	}
	if err := attestation.Validate(); err != nil {
		t.Fatalf("attestation Validate() error = %v", err)
	}
}

func TestExecutionLifecycleAttestationPreservesAuthorizationUnknown(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	authorization, err := NewExecutionLifecycleAuthorization(
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
	attestation, err := NewExecutionLifecycleAttestation(authorization, "reverse-observation-digest", now.Add(time.Minute))
	if err != nil {
		t.Fatalf("NewExecutionLifecycleAttestation() error = %v", err)
	}
	if attestation.Status != ExecutionLifecycleAttestationUnknown ||
		attestation.MissingStage != "authorization:replay" {
		t.Fatalf("unknown attestation = %#v", attestation)
	}
	if err := attestation.Validate(); err != nil {
		t.Fatalf("unknown attestation Validate() error = %v", err)
	}
}

func TestExecutionLifecycleAttestationPreservesMissingObservation(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	authorization, err := NewExecutionLifecycleAuthorization(
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
	authorization.Status = ExecutionLifecycleAuthorized
	authorization.MissingStage = ""
	attestation, err := NewExecutionLifecycleAttestation(authorization, "", now.Add(time.Minute))
	if err != nil {
		t.Fatalf("NewExecutionLifecycleAttestation() error = %v", err)
	}
	if attestation.Status != ExecutionLifecycleAttestationUnknown ||
		attestation.MissingStage != "authorization" {
		t.Fatalf("invalid authorization attestation = %#v", attestation)
	}
	if err := attestation.Validate(); err != nil {
		t.Fatalf("invalid authorization attestation Validate() error = %v", err)
	}
}

func TestExecutionLifecycleAttestationRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	attestation, err := NewExecutionLifecycleAttestation(
		ExecutionLifecycleAuthorization{
			Schema:              ExecutionLifecycleAuthorizationSchemaV1,
			AuthorizationDigest: "authorization-digest",
			Status:              ExecutionLifecycleAuthorizationUnknown,
			MissingStage:        "replay",
			ObservedAt:          now,
		},
		"reverse-observation-digest",
		now.Add(time.Minute),
	)
	if err != nil {
		t.Fatalf("NewExecutionLifecycleAttestation() error = %v", err)
	}
	attestation.ObservationDigest = "tampered-observation"
	if err := attestation.Validate(); err == nil {
		t.Fatal("tampered execution lifecycle attestation unexpectedly validated")
	}
}