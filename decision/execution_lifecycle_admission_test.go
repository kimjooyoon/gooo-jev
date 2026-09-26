package decision

import (
	"testing"
	"time"
)

func TestAdmitExecutionLifecycleAttestationWindowAcceptsCompleteWindow(t *testing.T) {
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
	admission, err := AdmitExecutionLifecycleAttestationWindow(window, now.Add(6*time.Minute))
	if err != nil {
		t.Fatalf("AdmitExecutionLifecycleAttestationWindow() error = %v", err)
	}
	if admission.Status != ExecutionLifecycleAdmitted {
		t.Fatalf("admission status = %q", admission.Status)
	}
	if err := admission.Validate(); err != nil {
		t.Fatalf("admission Validate() error = %v", err)
	}
}

func TestAdmitExecutionLifecycleAttestationWindowPreservesUnknownStage(t *testing.T) {
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
	admission, err := AdmitExecutionLifecycleAttestationWindow(window, now.Add(2*time.Minute))
	if err != nil {
		t.Fatalf("AdmitExecutionLifecycleAttestationWindow() error = %v", err)
	}
	if admission.Status != ExecutionLifecycleAdmissionUnknown ||
		admission.MissingStage != "window:attestation:authorization:replay" {
		t.Fatalf("unknown admission = %#v", admission)
	}
	if err := admission.Validate(); err != nil {
		t.Fatalf("unknown admission Validate() error = %v", err)
	}
}

func TestAdmitExecutionLifecycleAttestationWindowRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	window := ExecutionLifecycleAttestationWindow{
		Schema:            ExecutionLifecycleAttestationWindowSchemaV1,
		AttestationDigest: "attestation-digest",
		ObservationDigest: "reverse-observation-digest",
		NotBefore:         now,
		NotAfter:          now.Add(10 * time.Minute),
		Status:            ExecutionLifecycleAttestationWithinWindow,
		ObservedAt:        now.Add(time.Minute),
		WindowDigest:      "window-digest",
	}
	admission, err := AdmitExecutionLifecycleAttestationWindow(window, now.Add(2*time.Minute))
	if err != nil {
		t.Fatalf("AdmitExecutionLifecycleAttestationWindow() error = %v", err)
	}
	admission.ObservationDigest = "tampered-observation"
	if err := admission.Validate(); err == nil {
		t.Fatal("tampered execution lifecycle admission unexpectedly validated")
	}
}
