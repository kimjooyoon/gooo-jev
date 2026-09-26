package decision

import "testing"

func cycleCandidateObservation(status string) JEVImprovementCycleObservation {
	observation := JEVImprovementCycleObservation{
		Status: status, DeclarationDigest: "declaration-digest", IRDigest: "ir-digest", GenerationDigest: "generation-digest", ReverseObservationDigest: "reverse-observation-digest", MetricDigest: "metric-digest", ChangePlanDigest: "change-plan-digest", DispositionStatus: status, LSPCode: "jev.change-plan." + status, NonExecuting: true, NonAuthorizing: true,
	}
	observation.EvidenceDigest = digestJEVImprovementCycleObservation(observation.DeclarationDigest, observation.IRDigest, observation.GenerationDigest, observation.ReverseObservationDigest, observation.MetricDigest, observation.ChangePlanDigest, observation.DispositionStatus, observation.LSPCode)
	return observation
}

func TestGenerateExecutionEnvelopeJEVRevisionCandidateFromCycle(t *testing.T) {
	for _, status := range []string{"stable-for-review", "needs-revision"} {
		t.Run(status, func(t *testing.T) {
			got := GenerateExecutionEnvelopeJEVRevisionCandidateFromCycle(ExecutionEnvelopeJEVCycleRevisionCandidateGenerationInput{
				Cycle:                cycleCandidateObservation(status),
				Directive:            revisionCandidateDirective(),
				RevisionSource:       "gooo://revision/source/one",
				RevisionChangeDigest: "revision-change-digest",
				NonAuthorizing:       true,
			})
			if got.Status != "bound" || got.CycleStatus != status || got.CandidateStatus == "" || got.CandidateDigest == "" || got.CandidateEvidenceDigest == "" || got.BindingDigest == "" {
				t.Fatalf("got %+v", got)
			}
			if !got.NonExecuting || !got.NonAuthorizing || got.BoundRevisionChangeDigest == "" {
				t.Fatalf("missing boundary %+v", got)
			}
		})
	}
}

func TestGenerateExecutionEnvelopeJEVRevisionCandidateFromCycleBlocksHold(t *testing.T) {
	got := GenerateExecutionEnvelopeJEVRevisionCandidateFromCycle(ExecutionEnvelopeJEVCycleRevisionCandidateGenerationInput{
		Cycle:                cycleCandidateObservation("hold"),
		Directive:            revisionCandidateDirective(),
		RevisionSource:       "gooo://revision/source/one",
		RevisionChangeDigest: "revision-change-digest",
		NonAuthorizing:       true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "feedback-hold" || got.BindingDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestGenerateExecutionEnvelopeJEVRevisionCandidateFromCyclePreservesTampering(t *testing.T) {
	cycle := cycleCandidateObservation("stable-for-review")
	cycle.EvidenceDigest = "tampered"
	got := GenerateExecutionEnvelopeJEVRevisionCandidateFromCycle(ExecutionEnvelopeJEVCycleRevisionCandidateGenerationInput{
		Cycle:                cycle,
		Directive:            revisionCandidateDirective(),
		RevisionSource:       "gooo://revision/source/one",
		RevisionChangeDigest: "revision-change-digest",
		NonAuthorizing:       true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "improvement-cycle-observation" || got.BindingDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestGenerateExecutionEnvelopeJEVRevisionCandidateFromCycleRejectsAuthorization(t *testing.T) {
	got := GenerateExecutionEnvelopeJEVRevisionCandidateFromCycle(ExecutionEnvelopeJEVCycleRevisionCandidateGenerationInput{NonAuthorizing: false})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.BindingDigest != "" {
		t.Fatalf("got %+v", got)
	}
}
