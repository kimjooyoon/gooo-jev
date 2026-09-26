package decision

import "testing"

func replayLedgerCycleLSPObservation(status string) JEVReplayLedgerCycleObservation {
	return ObserveJEVReplayLedgerCycle(ExecutionEnvelopeJEVReplayLedgerCycleInput{
		Checkpoint:       BindJEVFullProvenanceReplayLedgerCheckpoint(ExecutionEnvelopeJEVFullProvenanceReplayLedgerInput{Provenance: fullProvenanceReplayLedgerFixture(), Ledger: replayLedgerCheckpointFixture(), NonAuthorizing: true}),
		ChangePlanDigest: "change-plan-digest",
		Feedback:         feedbackCycleAggregation(status),
		NonAuthorizing:   true,
	})
}

func TestProjectJEVReplayLedgerCycleToLSP(t *testing.T) {
	ready := ProjectJEVReplayLedgerCycleToLSP(ExecutionEnvelopeJEVReplayLedgerCycleLSPInput{
		Observation:    replayLedgerCycleLSPObservation("stable-for-review"),
		NonAuthorizing: true,
	})
	if ready.Status != "ready" || !ready.Publishable || ready.Severity != "info" ||
		ready.EvidenceDigest == "" || ready.EvidencePrefixDigest == "" ||
		ready.LedgerEvidenceDigest == "" {
		t.Fatalf("ready = %#v, want publishable stable projection", ready)
	}
	if ready.Code != "jev.replay-ledger-cycle.stable-for-review" {
		t.Fatalf("ready code = %q", ready.Code)
	}

	revision := ProjectJEVReplayLedgerCycleToLSP(ExecutionEnvelopeJEVReplayLedgerCycleLSPInput{
		Observation:    replayLedgerCycleLSPObservation("needs-revision"),
		NonAuthorizing: true,
	})
	if revision.Status != "ready" || revision.Severity != "warning" ||
		revision.Code != "jev.replay-ledger-cycle.needs-revision" {
		t.Fatalf("revision = %#v, want warning projection", revision)
	}
}

func TestProjectJEVReplayLedgerCycleToLSPPreservesUnknown(t *testing.T) {
	unknown := ProjectJEVReplayLedgerCycleToLSP(ExecutionEnvelopeJEVReplayLedgerCycleLSPInput{
		Observation:    replayLedgerCycleLSPObservation("UNKNOWN"),
		NonAuthorizing: true,
	})
	if unknown.Status != "diagnostic" || !unknown.Publishable ||
		unknown.MissingStage != "feedback-aggregation" ||
		unknown.MissingStageIndex != 5 || unknown.EvidencePrefixDigest == "" {
		t.Fatalf("unknown = %#v, want feedback diagnostic", unknown)
	}
}

func TestProjectJEVReplayLedgerCycleToLSPRejectsTampering(t *testing.T) {
	observation := replayLedgerCycleLSPObservation("stable-for-review")
	observation.EvidenceDigest = "tampered"
	projected := ProjectJEVReplayLedgerCycleToLSP(ExecutionEnvelopeJEVReplayLedgerCycleLSPInput{
		Observation:    observation,
		NonAuthorizing: true,
	})
	if projected.Status != "UNKNOWN" || projected.Publishable ||
		projected.Code != "jev-replay-ledger-cycle-integrity" {
		t.Fatalf("tampered = %#v, want integrity UNKNOWN", projected)
	}
}

func TestProjectJEVReplayLedgerCycleToLSPRejectsAuthorization(t *testing.T) {
	projected := ProjectJEVReplayLedgerCycleToLSP(ExecutionEnvelopeJEVReplayLedgerCycleLSPInput{
		Observation:    replayLedgerCycleLSPObservation("stable-for-review"),
		NonAuthorizing: false,
	})
	if projected.Status != "UNKNOWN" || projected.NonAuthorizing ||
		projected.Code != "authorization-boundary" {
		t.Fatalf("unauthorized = %#v, want authorization UNKNOWN", projected)
	}
}