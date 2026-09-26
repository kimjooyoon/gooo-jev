package decision

import "testing"

func confirmedActionDerivedCandidateFeedback(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackBinding {
	t.Helper()
	metric := confirmedActionDerivedCandidateMetric(t)
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-direction-confirmed",
		NonAuthorizing:          true,
	})
}

func refutedActionDerivedCandidateFeedback(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackBinding {
	t.Helper()
	metric := refutedActionDerivedCandidateMetric(t)
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-direction-refuted",
		NonAuthorizing:          true,
	})
}

func unknownActionDerivedCandidateFeedback(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackBinding {
	t.Helper()
	metric := unknownActionDerivedCandidateMetric(t)
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-direction-unknown",
		NonAuthorizing:          true,
	})
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionReview(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionInput{
		Feedback:       confirmedActionDerivedCandidateFeedback(t),
		CandidateSource: "gooo://candidate/action-confirmed",
		NonAuthorizing: true,
	})
	if output.Status != "bound" || output.Directive != jevImprovementDirectiveExternalReview ||
		output.CandidateID == "" || output.RevisionCandidateDigest == "" {
		t.Fatalf("output = %#v, want external review direction", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionRevision(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionInput{
		Feedback:       refutedActionDerivedCandidateFeedback(t),
		CandidateSource: "gooo://candidate/action-refuted",
		NonAuthorizing: true,
	})
	if output.Status != "bound" || output.Directive != jevImprovementDirectiveRevision ||
		output.DirectiveStatus != jevImprovementDirectiveRevision {
		t.Fatalf("output = %#v, want revision direction", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionHold(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionInput{
		Feedback:       unknownActionDerivedCandidateFeedback(t),
		CandidateSource: "gooo://candidate/action-unknown",
		NonAuthorizing: true,
	})
	if output.Status != "bound" || output.Directive != jevImprovementDirectiveEvidence {
		t.Fatalf("output = %#v, want evidence hold direction", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionRejectsTampering(t *testing.T) {
	feedback := confirmedActionDerivedCandidateFeedback(t)
	feedback.FeedbackDigest = "tampered"
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionInput{
		Feedback:       feedback,
		CandidateSource: "gooo://candidate/action-tampered",
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-feedback-validation" {
		t.Fatalf("output = %#v, want candidate feedback validation UNKNOWN", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionRequiresSource(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionInput{
		Feedback:       confirmedActionDerivedCandidateFeedback(t),
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "candidate-source" {
		t.Fatalf("output = %#v, want candidate-source UNKNOWN", output)
	}
}