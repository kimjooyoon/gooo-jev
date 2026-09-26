package decision

import "testing"

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionReview(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionInput{
		Feedback:       BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackInput{Metric: clearRevisionActionMetric(t), CandidateDigest: "candidate-review", ReplayObservationDigest: "replay-review", NonAuthorizing: true}),
		CandidateSource: "gooo://candidate/action-review",
		NonAuthorizing: true,
	})
	if output.Status != "bound" || output.Directive != jevImprovementDirectiveExternalReview {
		t.Fatalf("output = %#v, want external review direction", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionRevision(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionInput{
		Feedback:       BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackInput{Metric: counterexampleRevisionActionMetric(t), CandidateDigest: "candidate-revision", ReplayObservationDigest: "replay-revision", NonAuthorizing: true}),
		CandidateSource: "gooo://candidate/action-revision",
		NonAuthorizing: true,
	})
	if output.Status != "bound" || output.Directive != jevImprovementDirectiveRevision ||
		output.DirectiveStatus != jevImprovementDirectiveRevision {
		t.Fatalf("output = %#v, want revision direction", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionHold(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionInput{
		Feedback:       BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackInput{Metric: unknownRevisionActionMetric(t), CandidateDigest: "candidate-hold", ReplayObservationDigest: "replay-hold", NonAuthorizing: true}),
		CandidateSource: "gooo://candidate/action-hold",
		NonAuthorizing: true,
	})
	if output.Status != "bound" || output.Directive != jevImprovementDirectiveEvidence {
		t.Fatalf("output = %#v, want evidence hold direction", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionRejectsTampering(t *testing.T) {
	feedback := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackInput{
		Metric:                  clearRevisionActionMetric(t),
		CandidateDigest:         "candidate-tampered",
		ReplayObservationDigest: "replay-tampered",
		NonAuthorizing:          true,
	})
	feedback.FeedbackDigest = "tampered"
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionInput{
		Feedback:       feedback,
		CandidateSource: "gooo://candidate/action-tampered",
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-feedback-validation" {
		t.Fatalf("output = %#v, want feedback validation UNKNOWN", output)
	}
}