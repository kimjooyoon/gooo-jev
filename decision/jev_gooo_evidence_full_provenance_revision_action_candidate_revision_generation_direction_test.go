package decision

import "testing"

func confirmedGeneratedCandidateRevisionFeedbackForDirection(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationFeedbackBinding {
	t.Helper()
	metric := confirmedGeneratedCandidateRevisionMetricForFeedback(t)
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "generated-candidate-revision-direction-confirmed",
		NonAuthorizing:          true,
	})
}

func refutedGeneratedCandidateRevisionFeedbackForDirection(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationFeedbackBinding {
	t.Helper()
	metric := refutedGeneratedCandidateRevisionMetricForFeedback(t)
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "generated-candidate-revision-direction-refuted",
		NonAuthorizing:          true,
	})
}

func unknownGeneratedCandidateRevisionFeedbackForDirection(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationFeedbackBinding {
	t.Helper()
	metric := unknownGeneratedCandidateRevisionMetricForFeedback(t)
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "generated-candidate-revision-direction-unknown",
		NonAuthorizing:          true,
	})
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationDirectionReview(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationDirectionInput{
		Feedback:       confirmedGeneratedCandidateRevisionFeedbackForDirection(t),
		CandidateSource: "gooo://candidate/generated-revision-confirmed",
		NonAuthorizing: true,
	})
	if output.Status != "bound" || output.Directive != jevImprovementDirectiveExternalReview ||
		output.SourceCandidateID == "" || output.RevisionCandidateDigest == "" ||
		output.EvidencePrefixDigest == "" {
		t.Fatalf("output = %#v, want external review direction", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationDirectionRevision(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationDirectionInput{
		Feedback:       refutedGeneratedCandidateRevisionFeedbackForDirection(t),
		CandidateSource: "gooo://candidate/generated-revision-refuted",
		NonAuthorizing: true,
	})
	if output.Status != "bound" || output.Directive != jevImprovementDirectiveRevision ||
		output.DirectiveStatus != jevImprovementDirectiveRevision {
		t.Fatalf("output = %#v, want revision direction", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationDirectionHold(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationDirectionInput{
		Feedback:       unknownGeneratedCandidateRevisionFeedbackForDirection(t),
		CandidateSource: "gooo://candidate/generated-revision-unknown",
		NonAuthorizing: true,
	})
	if output.Status != "bound" || output.Directive != jevImprovementDirectiveEvidence {
		t.Fatalf("output = %#v, want evidence hold direction", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationDirectionRejectsTampering(t *testing.T) {
	feedback := confirmedGeneratedCandidateRevisionFeedbackForDirection(t)
	feedback.FeedbackDigest = "tampered"
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationDirectionInput{
		Feedback:       feedback,
		CandidateSource: "gooo://candidate/generated-revision-tampered",
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-feedback-validation" {
		t.Fatalf("output = %#v, want feedback validation UNKNOWN", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationDirectionRequiresSource(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationDirectionInput{
		Feedback:       confirmedGeneratedCandidateRevisionFeedbackForDirection(t),
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "candidate-source" {
		t.Fatalf("output = %#v, want candidate-source UNKNOWN", output)
	}
}
