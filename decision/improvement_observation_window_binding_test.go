package decision

import (
	"testing"
	"time"
)

func makeImprovementObservationWindow(t *testing.T, status ExecutionLifecycleAttestationWindowStatus) ExecutionLifecycleAttestationWindow {
	t.Helper()
	now := time.Unix(1_800_000_000, 0).UTC()
	window := ExecutionLifecycleAttestationWindow{
		Schema:            ExecutionLifecycleAttestationWindowSchemaV1,
		AttestationDigest: "attestation-digest",
		ObservationDigest: "observation-digest",
		NotBefore:         now.Add(-time.Minute),
		NotAfter:          now.Add(time.Hour),
		Status:            status,
		ObservedAt:        now,
	}
	if status == ExecutionLifecycleAttestationWindowUnknown {
		window.MissingStage = "attestation"
	}
	digest, err := window.computeDigest()
	if err != nil {
		t.Fatalf("computeDigest() error = %v", err)
	}
	window.WindowDigest = digest
	if err := window.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	return window
}

func TestBindImprovementObservationWindow(t *testing.T) {
	output := BindImprovementObservationWindow(ImprovementObservationWindowBindingInput{
		ReviewGateStatus: "accepted-for-observation",
		CandidateDigest:  "candidate-digest",
		ReviewDigest:     "review-digest",
		Window:           makeImprovementObservationWindow(t, ExecutionLifecycleAttestationWithinWindow),
		NonAuthorizing:   true,
	})
	if output.Status != "observation-admitted" || !output.ObservationOnly ||
		output.AdmissionDigest == "" || output.ExecutionGranted ||
		!output.NonAuthorizing {
		t.Fatalf("unexpected observation admission: %#v", output)
	}

	output = BindImprovementObservationWindow(ImprovementObservationWindowBindingInput{
		ReviewGateStatus: "rejected",
		CandidateDigest:  "candidate-digest",
		ReviewDigest:     "review-digest",
		Window:           makeImprovementObservationWindow(t, ExecutionLifecycleAttestationWithinWindow),
		NonAuthorizing:   true,
	})
	if output.Status != "rejected" || output.ObservationOnly || output.ExecutionGranted {
		t.Fatalf("unexpected rejection: %#v", output)
	}

	output = BindImprovementObservationWindow(ImprovementObservationWindowBindingInput{
		ReviewGateStatus: "review",
		CandidateDigest:  "candidate-digest",
		ReviewDigest:     "review-digest",
		Window:           makeImprovementObservationWindow(t, ExecutionLifecycleAttestationWithinWindow),
		NonAuthorizing:   true,
	})
	if output.Status != "review" || output.MissingStage != "review-gate" {
		t.Fatalf("unexpected review hold: %#v", output)
	}

	output = BindImprovementObservationWindow(ImprovementObservationWindowBindingInput{
		ReviewGateStatus: "accepted-for-observation",
		CandidateDigest:  "candidate-digest",
		ReviewDigest:     "review-digest",
		Window:           makeImprovementObservationWindow(t, ExecutionLifecycleAttestationWindowUnknown),
		NonAuthorizing:   true,
	})
	if output.Status != "hold" || output.MissingStage != "attestation-window:attestation" {
		t.Fatalf("unexpected window hold: %#v", output)
	}

	input := ImprovementObservationWindowBindingInput{
		ReviewGateStatus: "accepted-for-observation",
		CandidateDigest:  "candidate-digest",
		ReviewDigest:     "review-digest",
		Window:           makeImprovementObservationWindow(t, ExecutionLifecycleAttestationWithinWindow),
		NonAuthorizing:   true,
	}
	input.Window.WindowDigest = "tampered"
	output = BindImprovementObservationWindow(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "attestation-window" {
		t.Fatalf("unexpected tampered window: %#v", output)
	}

	input.Window = makeImprovementObservationWindow(t, ExecutionLifecycleAttestationWithinWindow)
	input.ExecutionGranted = true
	output = BindImprovementObservationWindow(input)
	if output.Status != "UNKNOWN" || output.NonAuthorizing ||
		output.MissingStage != "authorization-boundary" {
		t.Fatalf("unexpected execution claim: %#v", output)
	}
}