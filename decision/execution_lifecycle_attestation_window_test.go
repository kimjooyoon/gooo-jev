package decision

import (
	"testing"
	"time"
)

func TestObserveExecutionLifecycleAttestationWindowAcceptsObservedTime(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	grant, boundary, identity, decisionReceipt := lifecycleInputs(t, now)
	suspension, resumption, reverse := lifecycleReplayInputs(t, now)
	replayObservation, err := NewExecutionLifecycleReplay(suspension, resumption, reverse, now.Add(3*time.Minute))
	if err != nil {
		t.Fatalf("NewExecutionLifecycleReplay() error = %v", err)
	}
	authorization, err := NewExecutionLifecycleAuthorization(
		replayObservation,
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
	window, err := ObserveExecutionLifecycleAttestationWindow(attestation, now, now.Add(10*time.Minute))
	if err != nil {
		t.Fatalf("ObserveExecutionLifecycleAttestationWindow() error = %v", err)
	}
	if window.Status != ExecutionLifecycleAttestationWithinWindow {
		t.Fatalf("window status = %q", window.Status)
	}
	if err := window.Validate(); err != nil {
		t.Fatalf("window Validate() error = %v", err)
	}
}

func TestObserveExecutionLifecycleAttestationWindowPreservesOutOfWindow(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	attestation := ExecutionLifecycleAttestation{
		Schema:              ExecutionLifecycleAttestationSchemaV1,
		AuthorizationDigest: "authorization-digest",
		ObservationDigest:   "reverse-observation-digest",
		Status:              ExecutionLifecycleAttested,
		ObservedAt:          now,
		AttestationDigest:   "attestation-digest",
	}
	window, err := ObserveExecutionLifecycleAttestationWindow(attestation, now.Add(time.Minute), now.Add(10*time.Minute))
	if err != nil {
		t.Fatalf("ObserveExecutionLifecycleAttestationWindow() error = %v", err)
	}
	if window.Status != ExecutionLifecycleAttestationWindowUnknown ||
		window.MissingStage != "attestation" {
		t.Fatalf("invalid attestation window = %#v", window)
	}
	if err := window.Validate(); err != nil {
		t.Fatalf("invalid attestation window Validate() error = %v", err)
	}
}

func TestObserveExecutionLifecycleAttestationWindowPreservesNestedUnknown(t *testing.T) {
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
	window, err := ObserveExecutionLifecycleAttestationWindow(attestation, now, now.Add(10*time.Minute))
	if err != nil {
		t.Fatalf("ObserveExecutionLifecycleAttestationWindow() error = %v", err)
	}
	if window.Status != ExecutionLifecycleAttestationWindowUnknown ||
		window.MissingStage != "attestation:authorization:replay" {
		t.Fatalf("nested unknown window = %#v", window)
	}
	if err := window.Validate(); err != nil {
		t.Fatalf("nested unknown window Validate() error = %v", err)
	}
}

func TestObserveExecutionLifecycleAttestationWindowRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	window, err := ObserveExecutionLifecycleAttestationWindow(
		ExecutionLifecycleAttestation{
			Schema:              ExecutionLifecycleAttestationSchemaV1,
			AuthorizationDigest: "authorization-digest",
			ObservationDigest:   "reverse-observation-digest",
			Status:              ExecutionLifecycleAttested,
			ObservedAt:          now,
			AttestationDigest:   "attestation-digest",
		},
		now.Add(time.Minute),
		now.Add(10*time.Minute),
	)
	if err != nil {
		t.Fatalf("ObserveExecutionLifecycleAttestationWindow() error = %v", err)
	}
	window.ObservationDigest = "tampered-observation"
	if err := window.Validate(); err == nil {
		t.Fatal("tampered execution lifecycle attestation window unexpectedly validated")
	}
}