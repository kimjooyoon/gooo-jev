package decision

import "testing"

func fullProvenanceReplayLedgerFixture() ExecutionEnvelopeFullProvenanceSourceBinding {
	return ExecutionEnvelopeFullProvenanceSourceBinding{
		Status:                         "complete",
		DeclarationDigest:              "declaration-digest",
		IRDigest:                       "ir-digest",
		GenerationDigest:               "generation-digest",
		BindingDigest:                  "binding-digest",
		ReverseObservationDigest:       "reverse-observation-digest",
		MetricDigest:                   "metric-a",
		EvidenceDigest:                 "provenance-evidence-digest",
		CompletenessDigest:             "completeness-digest",
		NonExecuting:                   true,
		NonAuthorizing:                 true,
	}
}

func replayLedgerCheckpointFixture() JEVImprovementReplayFeedbackLedger {
	return AppendJEVImprovementReplayFeedbackLedger(ExecutionEnvelopeJEVReplayFeedbackLedgerInput{
		Feedback:       replayFeedbackLedgerFixture(jevReplayFeedbackConfirmed, "metric-a"),
		MetricDigest:   "metric-a",
		NonAuthorizing: true,
	})
}

func TestBindJEVFullProvenanceReplayLedgerCheckpoint(t *testing.T) {
	checkpoint := BindJEVFullProvenanceReplayLedgerCheckpoint(ExecutionEnvelopeJEVFullProvenanceReplayLedgerInput{
		Provenance:     fullProvenanceReplayLedgerFixture(),
		Ledger:         replayLedgerCheckpointFixture(),
		NonAuthorizing: true,
	})
	if checkpoint.Status != jevFullProvenanceReplayLedgerCheckpointComplete ||
		checkpoint.MissingStage != "" {
		t.Fatalf("checkpoint = %#v, want complete", checkpoint)
	}
	if err := checkpoint.Validate(); err != nil {
		t.Fatalf("checkpoint should validate: %v", err)
	}
}

func TestBindJEVFullProvenanceReplayLedgerCheckpointRejectsMissingStage(t *testing.T) {
	provenance := fullProvenanceReplayLedgerFixture()
	provenance.IRDigest = ""
	checkpoint := BindJEVFullProvenanceReplayLedgerCheckpoint(ExecutionEnvelopeJEVFullProvenanceReplayLedgerInput{
		Provenance:     provenance,
		Ledger:         replayLedgerCheckpointFixture(),
		NonAuthorizing: true,
	})
	if checkpoint.Status != jevFullProvenanceReplayLedgerCheckpointUnknown ||
		checkpoint.MissingStage != "ir" {
		t.Fatalf("checkpoint = %#v, want UNKNOWN at ir", checkpoint)
	}
}

func TestBindJEVFullProvenanceReplayLedgerCheckpointRejectsLedgerTampering(t *testing.T) {
	ledger := replayLedgerCheckpointFixture()
	ledger.EvidenceDigest = "tampered"
	checkpoint := BindJEVFullProvenanceReplayLedgerCheckpoint(ExecutionEnvelopeJEVFullProvenanceReplayLedgerInput{
		Provenance:     fullProvenanceReplayLedgerFixture(),
		Ledger:         ledger,
		NonAuthorizing: true,
	})
	if checkpoint.Status != jevFullProvenanceReplayLedgerCheckpointUnknown ||
		checkpoint.MissingStage != "replay-feedback-ledger" {
		t.Fatalf("checkpoint = %#v, want UNKNOWN at replay-feedback-ledger", checkpoint)
	}
}

func TestBindJEVFullProvenanceReplayLedgerCheckpointRejectsAuthorization(t *testing.T) {
	checkpoint := BindJEVFullProvenanceReplayLedgerCheckpoint(ExecutionEnvelopeJEVFullProvenanceReplayLedgerInput{
		Provenance:     fullProvenanceReplayLedgerFixture(),
		Ledger:         replayLedgerCheckpointFixture(),
		NonAuthorizing: false,
	})
	if checkpoint.Status != jevFullProvenanceReplayLedgerCheckpointUnknown ||
		checkpoint.MissingStage != "authorization-boundary" ||
		checkpoint.NonAuthorizing {
		t.Fatalf("checkpoint = %#v, want authorization UNKNOWN", checkpoint)
	}
}