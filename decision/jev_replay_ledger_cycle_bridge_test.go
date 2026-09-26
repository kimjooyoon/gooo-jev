package decision

import "testing"

func TestObserveJEVReplayLedgerCycle(t *testing.T) {
	observation := ObserveJEVReplayLedgerCycle(ExecutionEnvelopeJEVReplayLedgerCycleInput{
		Checkpoint:       BindJEVFullProvenanceReplayLedgerCheckpoint(ExecutionEnvelopeJEVFullProvenanceReplayLedgerInput{Provenance: fullProvenanceReplayLedgerFixture(), Ledger: replayLedgerCheckpointFixture(), NonAuthorizing: true}),
		ChangePlanDigest: "change-plan-digest",
		Feedback:         feedbackCycleAggregation("stable-for-review"),
		NonAuthorizing:   true,
	})
	if observation.Status != "stable-for-review" || observation.CycleStatus != "stable-for-review" ||
		observation.CycleEvidenceDigest == "" || observation.LedgerEvidenceDigest == "" {
		t.Fatalf("observation = %#v, want stable-for-review with ledger evidence", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("observation should validate: %v", err)
	}
}

func TestObserveJEVReplayLedgerCycleRejectsCheckpointTampering(t *testing.T) {
	checkpoint := BindJEVFullProvenanceReplayLedgerCheckpoint(ExecutionEnvelopeJEVFullProvenanceReplayLedgerInput{
		Provenance: fullProvenanceReplayLedgerFixture(),
		Ledger:     replayLedgerCheckpointFixture(),
		NonAuthorizing: true,
	})
	checkpoint.LedgerEvidenceDigest = "tampered"
	observation := ObserveJEVReplayLedgerCycle(ExecutionEnvelopeJEVReplayLedgerCycleInput{
		Checkpoint:       checkpoint,
		ChangePlanDigest: "change-plan-digest",
		Feedback:         feedbackCycleAggregation("stable-for-review"),
		NonAuthorizing:   true,
	})
	if observation.Status != "UNKNOWN" || observation.MissingStage != "replay-ledger-checkpoint" {
		t.Fatalf("observation = %#v, want checkpoint UNKNOWN", observation)
	}
}

func TestObserveJEVReplayLedgerCyclePreservesCycleMissingStage(t *testing.T) {
	observation := ObserveJEVReplayLedgerCycle(ExecutionEnvelopeJEVReplayLedgerCycleInput{
		Checkpoint:       BindJEVFullProvenanceReplayLedgerCheckpoint(ExecutionEnvelopeJEVFullProvenanceReplayLedgerInput{Provenance: fullProvenanceReplayLedgerFixture(), Ledger: replayLedgerCheckpointFixture(), NonAuthorizing: true}),
		ChangePlanDigest: "change-plan-digest",
		Feedback:         feedbackCycleAggregation("UNKNOWN"),
		NonAuthorizing:   true,
	})
	if observation.Status != "UNKNOWN" || observation.MissingStage != "feedback-aggregation" {
		t.Fatalf("observation = %#v, want feedback-aggregation UNKNOWN", observation)
	}
}

func TestObserveJEVReplayLedgerCycleRejectsAuthorization(t *testing.T) {
	observation := ObserveJEVReplayLedgerCycle(ExecutionEnvelopeJEVReplayLedgerCycleInput{
		Checkpoint:       BindJEVFullProvenanceReplayLedgerCheckpoint(ExecutionEnvelopeJEVFullProvenanceReplayLedgerInput{Provenance: fullProvenanceReplayLedgerFixture(), Ledger: replayLedgerCheckpointFixture(), NonAuthorizing: true}),
		ChangePlanDigest: "change-plan-digest",
		Feedback:         feedbackCycleAggregation("stable-for-review"),
		NonAuthorizing:   false,
	})
	if observation.Status != "UNKNOWN" || observation.MissingStage != "authorization-boundary" || observation.NonAuthorizing {
		t.Fatalf("observation = %#v, want authorization UNKNOWN", observation)
	}
}