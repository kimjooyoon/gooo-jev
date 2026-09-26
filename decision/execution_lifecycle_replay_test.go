package decision

import (
	"testing"
	"time"
)

func lifecycleReplayInputs(t *testing.T, now time.Time) (ExecutionLifecycleReceipt, ExecutionLifecycleReceipt, ReverseObservation) {
	t.Helper()
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
	execution, err := NewExecutionReceipt(
		grant,
		boundary,
		identity,
		decisionReceipt,
		"result-digest-1",
		ExecutionCompleted,
		now.Add(2*time.Minute),
	)
	if err != nil {
		t.Fatalf("NewExecutionReceipt() error = %v", err)
	}
	reverse, err := NewReverseObservation(
		execution,
		"output-digest-1",
		"verifier-digest-1",
		"",
		ProvenanceObserved,
		now.Add(2*time.Minute),
	)
	if err != nil {
		t.Fatalf("NewReverseObservation() error = %v", err)
	}
	return suspension, resumption, reverse
}

func TestExecutionLifecycleReplayBindsReverseObservation(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	suspension, resumption, reverse := lifecycleReplayInputs(t, now)
	observation, err := NewExecutionLifecycleReplay(suspension, resumption, reverse, now.Add(3*time.Minute))
	if err != nil {
		t.Fatalf("NewExecutionLifecycleReplay() error = %v", err)
	}
	if observation.Status != ExecutionLifecycleReplayObserved {
		t.Fatalf("replay status = %q", observation.Status)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("replay Validate() error = %v", err)
	}
}

func TestExecutionLifecycleReplayPreservesUnknownStage(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	unknown, err := NewExecutionLifecycleReplay(
		ExecutionLifecycleReceipt{},
		ExecutionLifecycleReceipt{},
		ReverseObservation{},
		now,
	)
	if err != nil {
		t.Fatalf("NewExecutionLifecycleReplay() error = %v", err)
	}
	if unknown.Status != ExecutionLifecycleReplayUnknown || unknown.MissingStage != "suspension" {
		t.Fatalf("unknown replay = %#v", unknown)
	}
	if err := unknown.Validate(); err != nil {
		t.Fatalf("unknown replay Validate() error = %v", err)
	}
}

func TestExecutionLifecycleReplayRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	suspension, resumption, reverse := lifecycleReplayInputs(t, now)
	observation, err := NewExecutionLifecycleReplay(suspension, resumption, reverse, now.Add(3*time.Minute))
	if err != nil {
		t.Fatalf("NewExecutionLifecycleReplay() error = %v", err)
	}
	observation.ReverseObservationDigest = "tampered-reverse-digest"
	if err := observation.Validate(); err == nil {
		t.Fatal("tampered lifecycle replay unexpectedly validated")
	}
}
