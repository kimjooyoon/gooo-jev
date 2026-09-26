package decision

import (
	"testing"
	"time"
)

func TestReplayExecutionLifecycleAdmissionConsumptionPreservesBoundedUse(t *testing.T) {
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
	}
	if err := consumption.assignDigest(); err != nil {
		t.Fatalf("consumption.assignDigest() error = %v", err)
	}
	replay, err := ReplayExecutionLifecycleAdmissionConsumption(consumption, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("ReplayExecutionLifecycleAdmissionConsumption() error = %v", err)
	}
	if replay.Status != ExecutionLifecycleAdmissionConsumptionReplayed ||
		replay.MaxUses != 2 || replay.ConsumedUses != 1 {
		t.Fatalf("replay = %#v", replay)
	}
	if err := replay.Validate(); err != nil {
		t.Fatalf("replay Validate() error = %v", err)
	}
}

func TestReplayExecutionLifecycleAdmissionConsumptionPreservesUnknown(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	consumption := ExecutionLifecycleAdmissionConsumption{
		Schema:            ExecutionLifecycleAdmissionConsumptionSchemaV1,
		AdmissionDigest:   "admission-digest",
		ObservationDigest: "reverse-observation-digest",
		UseDigest:         "use-digest",
		MaxUses:           1,
		ConsumedUses:      2,
		Status:            ExecutionLifecycleAdmissionConsumed,
		ConsumedAt:        now,
		ConsumptionDigest: "tampered-consumption-digest",
	}
	replay, err := ReplayExecutionLifecycleAdmissionConsumption(consumption, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("ReplayExecutionLifecycleAdmissionConsumption() error = %v", err)
	}
	if replay.Status != ExecutionLifecycleAdmissionConsumptionReplayUnknown ||
		replay.MissingStage != "consumption" {
		t.Fatalf("unknown replay = %#v", replay)
	}
	if err := replay.Validate(); err != nil {
		t.Fatalf("unknown replay Validate() error = %v", err)
	}
}

func TestReplayExecutionLifecycleAdmissionConsumptionRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	replay := ExecutionLifecycleAdmissionConsumptionReplay{
		Schema:            ExecutionLifecycleAdmissionConsumptionReplaySchemaV1,
		ConsumptionDigest: "consumption-digest",
		AdmissionDigest:   "admission-digest",
		ObservationDigest: "reverse-observation-digest",
		UseDigest:         "use-digest",
		MaxUses:           2,
		ConsumedUses:      1,
		Status:            ExecutionLifecycleAdmissionConsumptionReplayed,
		ObservedAt:        now,
	}
	if err := replay.assignDigest(); err != nil {
		t.Fatalf("replay.assignDigest() error = %v", err)
	}
	replay.UseDigest = "tampered-use"
	if err := replay.Validate(); err == nil {
		t.Fatal("tampered execution lifecycle admission consumption replay unexpectedly validated")
	}
}
