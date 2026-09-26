package decision

import "testing"

func goooRevisionCandidateInput(status string) ExecutionEnvelopeGoooRevisionCandidateInput {
	feedback := feedbackCycleAggregation(status)
	directive := DeriveJEVImprovementDirectionDirective(JEVImprovementDirectionDirectiveInput{
		Aggregation:     feedback,
		CandidateDigest: "candidate-digest",
		CandidateSource: "gooo://candidate/one",
		NonAuthorizing:  true,
	})
	return ExecutionEnvelopeGoooRevisionCandidateInput{
		Provenance:           fullProvenanceCycleBinding(),
		ChangePlanDigest:     "change-plan-digest",
		Feedback:             feedback,
		Directive:            directive,
		RevisionSource:       "gooo://revision/one",
		RevisionChangeDigest: "revision-change-digest",
		NonAuthorizing:       true,
	}
}

func TestGenerateExecutionEnvelopeJEVRevisionCandidateFromGoooFullProvenance(t *testing.T) {
	for _, status := range []string{"stable-for-review", "needs-revision"} {
		t.Run(status, func(t *testing.T) {
			got := GenerateExecutionEnvelopeJEVRevisionCandidateFromGoooFullProvenance(goooRevisionCandidateInput(status))
			if got.Status != "bound" || got.CandidateDigest == "" || got.CycleEvidenceDigest == "" || got.BindingDigest == "" || got.MissingStage != "" || !got.NonExecuting || !got.NonAuthorizing {
				t.Fatalf("got %+v", got)
			}
		})
	}
}

func TestGenerateExecutionEnvelopeJEVRevisionCandidateFromGoooFullProvenanceHoldsUnknownFeedback(t *testing.T) {
	got := GenerateExecutionEnvelopeJEVRevisionCandidateFromGoooFullProvenance(goooRevisionCandidateInput("hold"))
	if got.Status != "UNKNOWN" || got.MissingStage != "feedback-hold" || got.BindingDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestGenerateExecutionEnvelopeJEVRevisionCandidateFromGoooFullProvenanceRejectsMixedDirective(t *testing.T) {
	input := goooRevisionCandidateInput("stable-for-review")
	input.Directive.InputEvidenceDigest = "other-feedback"
	got := GenerateExecutionEnvelopeJEVRevisionCandidateFromGoooFullProvenance(input)
	if got.Status != "UNKNOWN" || got.MissingStage != "feedback-direction-binding" || got.BindingDigest != "" {
		t.Fatalf("got %+v", got)
	}
}
