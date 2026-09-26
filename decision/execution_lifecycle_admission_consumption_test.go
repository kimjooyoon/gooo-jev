package decision

import (
	"testing"
	"time"
)

func TestConsumeExecutionLifecycleAdmissionBindsUseAndLimit(t *testing.T) {
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
	consumption, err := ConsumeExecutionLifecycleAdmission(admission, "use-digest", 2, 1, now.Add(7*time.Minute))
	if err != nil {
		t.Fatalf("ConsumeExecutionLifecycleAdmission() error = %v", err)
	}
	if consumption.Status != ExecutionLifecycleAdmissionConsumed {
		t.Fatalf("consumption status = %q", consumption.Status)
	}
	if err := consumption.Validate(); err != nil {
		t.Fatalf("consumption Validate() error = %v", err)
	}
}

func TestConsumeExecutionLifecycleAdmissionPreservesUseLimitUnknown(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	admission := ExecutionLifecycleAdmission{
		Schema:            ExecutionLifecycleAdmissionSchemaV1,
		WindowDigest:      "window-digest",
		AttestationDigest: "attestation-digest",
		ObservationDigest: "reverse-observation-digest",
		Status:            ExecutionLifecycleAdmitted,
		AdmittedAt:        now,
		AdmissionDigest:   "admission-digest",
	}
	consumption, err := ConsumeExecutionLifecycleAdmission(admission, "use-digest", 1, 2, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("ConsumeExecutionLifecycleAdmission() error = %v", err)
	}
	if consumption.Status != ExecutionLifecycleAdmissionConsumptionUnknown ||
		consumption.MissingStage != "admission" {
		t.Fatalf("unknown consumption = %#v", consumption)
	}
	if err := consumption.Validate(); err != nil {
		t.Fatalf("unknown consumption Validate() error = %v", err)
	}
}

func TestConsumeExecutionLifecycleAdmissionRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	consumption := ExecutionLifecycleAdmissionConsumption{
		Schema:            ExecutionLifecycleAdmissionConsumptionSchemaV1,
		AdmissionDigest:   "admission-digest",
		ObservationDigest: "reverse-observation-digest",
		UseDigest:         "use-digest",
		MaxUses:           2,
		ConsumedUses:      1,
		Status:            ExecutionLifecycleAdmissionConsumed,
		ConsumedAt:        now,
		ConsumptionDigest: "consumption-digest",
	}
	consumption.UseDigest = "tampered-use"
	if err := consumption.Validate(); err == nil {
		t.Fatal("tampered execution lifecycle admission consumption unexpectedly validated")
	}
}
