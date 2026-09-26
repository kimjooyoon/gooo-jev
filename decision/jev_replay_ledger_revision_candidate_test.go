package decision

import "testing"

func replayLedgerRevisionCandidateCycle(status string) JEVReplayLedgerCycleObservation {
	return replayLedgerCycleLSPObservation(status)
}

func TestGenerateJEVReplayLedgerRevisionCandidateFromCycle(t *testing.T) {
	for _, status := range []string{"stable-for-review", "needs-revision"} {
		t.Run(status, func(t *testing.T) {
			binding := GenerateJEVReplayLedgerRevisionCandidateFromCycle(ExecutionEnvelopeJEVReplayLedgerRevisionCandidateInput{
				Cycle:                replayLedgerRevisionCandidateCycle(status),
				Directive:            revisionCandidateDirective(),
				RevisionSource:       "gooo://revision/source/replay-ledger",
				RevisionChangeDigest: "revision-change-digest",
				NonAuthorizing:       true,
			})
			if binding.Status != "bound" || binding.CycleStatus != status ||
				binding.CandidateDigest == "" || binding.CandidateEvidenceDigest == "" ||
				binding.BoundRevisionChangeDigest == "" || binding.BindingDigest == "" {
				t.Fatalf("binding = %#v, want bound candidate", binding)
			}
			if !binding.NonExecuting || !binding.NonAuthorizing {
				t.Fatalf("missing safety boundary: %#v", binding)
			}
		})
	}
}

func TestGenerateJEVReplayLedgerRevisionCandidateFromCycleBlocksHold(t *testing.T) {
	binding := GenerateJEVReplayLedgerRevisionCandidateFromCycle(ExecutionEnvelopeJEVReplayLedgerRevisionCandidateInput{
		Cycle:                replayLedgerRevisionCandidateCycle("hold"),
		Directive:            revisionCandidateDirective(),
		RevisionSource:       "gooo://revision/source/replay-ledger",
		RevisionChangeDigest: "revision-change-digest",
		NonAuthorizing:       true,
	})
	if binding.Status != "UNKNOWN" || binding.MissingStage != "feedback-hold" || binding.BindingDigest != "" {
		t.Fatalf("binding = %#v, want feedback-hold UNKNOWN", binding)
	}
}

func TestGenerateJEVReplayLedgerRevisionCandidateFromCycleRejectsTampering(t *testing.T) {
	cycle := replayLedgerRevisionCandidateCycle("stable-for-review")
	cycle.CycleEvidenceDigest = "tampered"
	binding := GenerateJEVReplayLedgerRevisionCandidateFromCycle(ExecutionEnvelopeJEVReplayLedgerRevisionCandidateInput{
		Cycle:                cycle,
		Directive:            revisionCandidateDirective(),
		RevisionSource:       "gooo://revision/source/replay-ledger",
		RevisionChangeDigest: "revision-change-digest",
		NonAuthorizing:       true,
	})
	if binding.Status != "UNKNOWN" || binding.MissingStage != "replay-ledger-cycle-observation" {
		t.Fatalf("binding = %#v, want cycle observation UNKNOWN", binding)
	}
}

func TestGenerateJEVReplayLedgerRevisionCandidateFromCycleRejectsAuthorization(t *testing.T) {
	binding := GenerateJEVReplayLedgerRevisionCandidateFromCycle(ExecutionEnvelopeJEVReplayLedgerRevisionCandidateInput{
		NonAuthorizing: false,
	})
	if binding.Status != "UNKNOWN" || binding.MissingStage != "authorization-boundary" || binding.NonAuthorizing {
		t.Fatalf("binding = %#v, want authorization UNKNOWN", binding)
	}
}