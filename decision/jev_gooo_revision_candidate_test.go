package decision

import "testing"

func goooRevisionDirective(feedback JEVImprovementFeedbackAggregation) JEVImprovementDirectionDirective {
	directive := JEVImprovementDirectionDirective{
		Status:              jevImprovementDirectiveRevision,
		Directive:            jevImprovementDirectiveRevision,
		CandidateDigest:     "candidate-digest",
		CandidateSource:     "gooo://candidate/one",
		InputEvidenceDigest: feedback.EvidenceDigest,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	directive.EvidenceDigest = digestJEVImprovementDirectionDirective(directive.Status, directive.Directive, directive.CandidateDigest, directive.CandidateSource, directive.InputEvidenceDigest)
	return directive
}

func goooRevisionCandidateInput(status string) ExecutionEnvelopeGoooRevisionCandidateInput {
	feedback := feedbackCycleAggregation(status)
	return ExecutionEnvelopeGoooRevisionCandidateInput{
		Provenance:           fullProvenanceCycleBinding(),
		ChangePlanDigest:     "change-plan-digest",
		Feedback:             feedback,
		Directive:            goooRevisionDirective(feedback),
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
	input.Directive.EvidenceDigest = digestJEVImprovementDirectionDirective(input.Directive.Status, input.Directive.Directive, input.Directive.CandidateDigest, input.Directive.CandidateSource, input.Directive.InputEvidenceDigest)
	got := GenerateExecutionEnvelopeJEVRevisionCandidateFromGoooFullProvenance(input)
	if got.Status != "UNKNOWN" || got.MissingStage != "feedback-direction-binding" || got.BindingDigest != "" {
		t.Fatalf("got %+v", got)
	}
}
