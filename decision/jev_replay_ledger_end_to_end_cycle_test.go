package decision

import "testing"

func replayLedgerEndToEndInput(ledger JEVImprovementReplayFeedbackLedger) ExecutionEnvelopeJEVReplayLedgerEndToEndCycleInput {
	return ExecutionEnvelopeJEVReplayLedgerEndToEndCycleInput{
		Checkpoint:       BindJEVFullProvenanceReplayLedgerCheckpoint(ExecutionEnvelopeJEVFullProvenanceReplayLedgerInput{Provenance: fullProvenanceReplayLedgerFixture(), Ledger: ledger, NonAuthorizing: true}),
		Ledger:           ledger,
		ChangePlanDigest: "change-plan-digest",
		NonAuthorizing:   true,
	}
}

func TestObserveJEVReplayLedgerCycleFromLedger(t *testing.T) {
	stableLedger := replayLedgerCheckpointFixture()
	stable := ObserveJEVReplayLedgerCycleFromLedger(replayLedgerEndToEndInput(stableLedger))
	if stable.Status != "stable-for-review" || stable.LedgerEvidenceDigest != stableLedger.EvidenceDigest {
		t.Fatalf("stable = %#v, want stable cycle", stable)
	}

	revisionLedger := AppendJEVImprovementReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerInput{
		Previous:       stableLedger,
		Feedback:       replayFeedbackLedgerFixture(jevReplayFeedbackRefuted, "metric-b"),
		MetricDigest:   "metric-b",
		NonAuthorizing: true,
	})
	revision := ObserveJEVReplayLedgerCycleFromLedger(replayLedgerEndToEndInput(revisionLedger))
	if revision.Status != "needs-revision" || revision.CycleStatus != "needs-revision" {
		t.Fatalf("revision = %#v, want needs-revision cycle", revision)
	}
}

func TestObserveJEVReplayLedgerCycleFromLedgerPreservesHold(t *testing.T) {
	holdLedger := AppendJEVImprovementReplayFeedbackLedger(ExecutionEnvelopeJEVImprovementReplayFeedbackLedgerInput{
		Feedback:       replayFeedbackLedgerFixture(jevReplayFeedbackUnknown, "metric-unknown"),
		MetricDigest:   "metric-unknown",
		NonAuthorizing: true,
	})
	hold := ObserveJEVReplayLedgerCycleFromLedger(replayLedgerEndToEndInput(holdLedger))
	if hold.Status != "hold" || hold.CycleStatus != "hold" {
		t.Fatalf("hold = %#v, want hold cycle", hold)
	}
}

func TestObserveJEVReplayLedgerCycleFromLedgerRejectsLedgerTampering(t *testing.T) {
	ledger := replayLedgerCheckpointFixture()
	ledger.EvidenceDigest = "tampered"
	observation := ObserveJEVReplayLedgerCycleFromLedger(replayLedgerEndToEndInput(ledger))
	if observation.Status != "UNKNOWN" || observation.MissingStage != "replay-feedback-ledger" {
		t.Fatalf("observation = %#v, want ledger UNKNOWN", observation)
	}
}

func TestObserveJEVReplayLedgerCycleFromLedgerRejectsAuthorization(t *testing.T) {
	input := replayLedgerEndToEndInput(replayLedgerCheckpointFixture())
	input.NonAuthorizing = false
	observation := ObserveJEVReplayLedgerCycleFromLedger(input)
	if observation.Status != "UNKNOWN" || observation.MissingStage != "authorization-boundary" || observation.NonAuthorizing {
		t.Fatalf("observation = %#v, want authorization UNKNOWN", observation)
	}
}