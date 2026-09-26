package decision

import "testing"

func makeImprovementReplayEvidenceBindingInput(t *testing.T, status ImprovementReplayStatus) ImprovementReplayEvidenceBindingInput {
	t.Helper()
	cycle, err := BuildImprovementCycle(
		"ledger-before-replay-binding",
		"ledger-after-replay-binding",
		"candidate://replay-binding",
		"gooo://jev/source/replay-binding",
		ImprovementObserved,
	)
	if err != nil {
		t.Fatalf("BuildImprovementCycle() error = %v", err)
	}
	recomputed := "ledger-after-replay-binding"
	if status == ImprovementDiverged {
		recomputed = "ledger-diverged"
	}
	if status == ImprovementReplayUnknown {
		recomputed = ""
	}
	replay, err := ReplayImprovementCycle(cycle, recomputed, "replay-evidence")
	if err != nil {
		t.Fatalf("ReplayImprovementCycle() error = %v", err)
	}
	return ImprovementReplayEvidenceBindingInput{
		CycleBindingStatus:                 "recorded",
		CycleDigest:                        cycle.CycleDigest,
		CycleEvidenceDigest:                "cycle-evidence",
		CandidateDigest:                    "candidate-digest",
		AdmissionDigest:                    "admission-digest",
		ReplayCycleDigest:                  replay.CycleDigest,
		ReplayRecordedResultLedgerDigest:   replay.RecordedResultLedgerDigest,
		ReplayRecomputedResultLedgerDigest: replay.RecomputedResultLedgerDigest,
		ReplayEvidenceDigest:               replay.EvidenceDigest,
		ReplayStatus:                       replay.Status,
		ReplayDigest:                       replay.ReplayDigest,
		NonAuthorizing:                     true,
	}
}

func TestBindImprovementReplayEvidence(t *testing.T) {
	output := BindImprovementReplayEvidence(makeImprovementReplayEvidenceBindingInput(t, ImprovementReplayed))
	if output.Status != "replayed" || output.ReplayStatus != ImprovementReplayed ||
		output.EvidenceDigest == "" || !output.NonExecuting ||
		!output.NonAuthorizing {
		t.Fatalf("unexpected replayed binding: %#v", output)
	}

	output = BindImprovementReplayEvidence(makeImprovementReplayEvidenceBindingInput(t, ImprovementDiverged))
	if output.Status != "diverged" || output.ReplayStatus != ImprovementDiverged {
		t.Fatalf("unexpected diverged binding: %#v", output)
	}

	output = BindImprovementReplayEvidence(makeImprovementReplayEvidenceInput(t, ImprovementReplayUnknown))
	if output.Status != "review" || output.MissingStage != "replay-unknown" {
		t.Fatalf("unexpected unknown replay binding: %#v", output)
	}

	input := makeImprovementReplayEvidenceBindingInput(t, ImprovementReplayed)
	input.CycleBindingStatus = "review"
	output = BindImprovementReplayEvidence(input)
	if output.Status != "review" || output.MissingStage != "cycle-outcome" {
		t.Fatalf("unexpected cycle review binding: %#v", output)
	}

	input = makeImprovementReplayEvidenceBindingInput(t, ImprovementReplayed)
	input.ReplayDigest = "tampered"
	output = BindImprovementReplayEvidence(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "improvement-replay" {
		t.Fatalf("unexpected tampered replay: %#v", output)
	}

	input = makeImprovementReplayEvidenceBindingInput(t, ImprovementReplayed)
	input.ReplayCycleDigest = "other-cycle"
	output = BindImprovementReplayEvidence(input)
	if output.Status != "UNKNOWN" || output.MissingStage != "cycle-replay-binding" {
		t.Fatalf("unexpected cycle mismatch: %#v", output)
	}
}