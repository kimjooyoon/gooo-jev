package decision

import (
	"testing"
	"time"
)

func TestReplayExecutionLifecycleAttestationBindsSameObservation(t *testing.T) {
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
	replayed, err := ReplayExecutionLifecycleAttestation(attestation, "reverse-observation-digest", now.Add(6*time.Minute))
	if err != nil {
		t.Fatalf("ReplayExecutionLifecycleAttestation() error = %v", err)
	}
	if replayed.Status != ExecutionLifecycleAttestationReplayed {
		t.Fatalf("replay status = %q", replayed.Status)
	}
	if err := replayed.Validate(); err != nil {
		t.Fatalf("replay Validate() error = %v", err)
	}
}

func TestReplayExecutionLifecycleAttestationPreservesMismatch(t *testing.T) {
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
	replayed, err := ReplayExecutionLifecycleAttestation(attestation, "tampered-observation", now.Add(2*time.Minute))
	if err != nil {
		t.Fatalf("ReplayExecutionLifecycleAttestation() error = %v", err)
	}
	if replayed.Status != ExecutionLifecycleAttestationReplayUnknown ||
		replayed.MissingStage != "attestation" {
		t.Fatalf("unknown replay = %#v", replayed)
	}
	if err := replayed.Validate(); err != nil {
		t.Fatalf("unknown replay Validate() error = %v", err)
	}
}

func TestReplayExecutionLifecycleAttestationPreservesObservationMismatch(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	attestation := ExecutionLifecycleAttestation{
		Schema:              ExecutionLifecycleAttestationSchemaV1,
		AuthorizationDigest: "authorization-digest",
		ObservationDigest:   "reverse-observation-digest",
		Status:              ExecutionLifecycleAttested,
		ObservedAt:          now,
		AttestationDigest:   "attestation-digest",
	}
	replayed, err := ReplayExecutionLifecycleAttestation(attestation, "other-observation-digest", now.Add(time.Minute))
	if err != nil {
		t.Fatalf("ReplayExecutionLifecycleAttestation() error = %v", err)
	}
	if replayed.Status != ExecutionLifecycleAttestationReplayUnknown ||
		replayed.MissingStage != "attestation" {
		t.Fatalf("invalid attestation replay = %#v", replayed)
	}
	if err := replayed.Validate(); err != nil {
		t.Fatalf("invalid attestation replay Validate() error = %v", err)
	}
}

func TestReplayExecutionLifecycleAttestationRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	replayed, err := ReplayExecutionLifecycleAttestation(
		ExecutionLifecycleAttestation{
			Schema:              ExecutionLifecycleAttestationSchemaV1,
			AuthorizationDigest: "authorization-digest",
			ObservationDigest:   "reverse-observation-digest",
			Status:              ExecutionLifecycleAttested,
			ObservedAt:          now,
			AttestationDigest:   "attestation-digest",
		},
		"reverse-observation-digest",
		now.Add(time.Minute),
	)
	if err != nil {
		t.Fatalf("ReplayExecutionLifecycleAttestation() error = %v", err)
	}
	replayed.ObservationDigest = "tampered-observation"
	if err := replayed.Validate(); err == nil {
		t.Fatal("tampered execution lifecycle attestation replay unexpectedly validated")
	}
}