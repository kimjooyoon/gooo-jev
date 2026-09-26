package decision

import (
	"testing"
	"time"
)

func TestReplayExecutionLifecycleAdmissionBindsObservation(t *testing.T) {
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
	replay, err := ReplayExecutionLifecycleAdmission(admission, "reverse-observation-digest", now.Add(7*time.Minute))
	if err != nil {
		t.Fatalf("ReplayExecutionLifecycleAdmission() error = %v", err)
	}
	if replay.Status != ExecutionLifecycleAdmissionReplayed {
		t.Fatalf("replay status = %q", replay.Status)
	}
	if err := replay.Validate(); err != nil {
		t.Fatalf("replay Validate() error = %v", err)
	}
}

func TestReplayExecutionLifecycleAdmissionPreservesNestedUnknown(t *testing.T) {
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
	replay, err := ReplayExecutionLifecycleAdmission(admission, "reverse-observation-digest", now.Add(3*time.Minute))
	if err != nil {
		t.Fatalf("ReplayExecutionLifecycleAdmission() error = %v", err)
	}
	if replay.Status != ExecutionLifecycleAdmissionReplayUnknown ||
		replay.MissingStage != "admission:window:attestation:authorization:replay" {
		t.Fatalf("unknown replay = %#v", replay)
	}
	if err := replay.Validate(); err != nil {
		t.Fatalf("unknown replay Validate() error = %v", err)
	}
}

func TestReplayExecutionLifecycleAdmissionRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	replay := ExecutionLifecycleAdmissionReplay{
		Schema:            ExecutionLifecycleAdmissionReplaySchemaV1,
		AdmissionDigest:   "admission-digest",
		WindowDigest:      "window-digest",
		AttestationDigest: "attestation-digest",
		ObservationDigest: "reverse-observation-digest",
		Status:            ExecutionLifecycleAdmissionReplayed,
		ObservedAt:        now,
		ReplayDigest:      "replay-digest",
	}
	replay.ObservationDigest = "tampered-observation"
	if err := replay.Validate(); err == nil {
		t.Fatal("tampered execution lifecycle admission replay unexpectedly validated")
	}
}
