package decision

import "testing"

func confirmedGeneratedExtendedCandidateRevisionFeedbackForDirection(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedbackBinding {
	t.Helper()
	metric := confirmedGeneratedExtendedCandidateRevisionMetricForFeedback(t)
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "generated-extended-candidate-direction-confirmed",
		NonAuthorizing:          true,
	})
}

func refutedGeneratedExtendedCandidateRevisionFeedbackForDirection(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedbackBinding {
	t.Helper()
	metric := refutedGeneratedExtendedCandidateRevisionMetricForFeedback(t)
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "generated-extended-candidate-direction-refuted",
		NonAuthorizing:          true,
	})
}

func unknownGeneratedExtendedCandidateRevisionFeedbackForDirection(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedbackBinding {
	t.Helper()
	metric := unknownGeneratedExtendedCandidateRevisionMetricForFeedback(t)
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "generated-extended-candidate-direction-unknown",
		NonAuthorizing:          true,
	})
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirectionReview(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirectionInput{
		Feedback:       confirmedGeneratedExtendedCandidateRevisionFeedbackForDirection(t),
		CandidateSource: "gooo://candidate/generated-extended-direction-confirmed",
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

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirectionRevision(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirectionInput{
		Feedback:       refutedGeneratedExtendedCandidateRevisionFeedbackForDirection(t),
		CandidateSource: "gooo://candidate/generated-extended-direction-refuted",
		NonAuthorizing: true,
	})
	if output.Status != "bound" || output.Directive != jevImprovementDirectiveRevision ||
		output.DirectiveStatus != jevImprovementDirectiveRevision {
		t.Fatalf("output = %#v, want revision direction", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirectionHold(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirectionInput{
		Feedback:       unknownGeneratedExtendedCandidateRevisionFeedbackForDirection(t),
		CandidateSource: "gooo://candidate/generated-extended-direction-unknown",
		NonAuthorizing: true,
	})
	if output.Status != "bound" || output.Directive != jevImprovementDirectiveEvidence {
		t.Fatalf("output = %#v, want evidence hold direction", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirectionRejectsTampering(t *testing.T) {
	feedback := confirmedGeneratedExtendedCandidateRevisionFeedbackForDirection(t)
	feedback.FeedbackDigest = "tampered"
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirectionInput{
		Feedback:       feedback,
		CandidateSource: "gooo://candidate/generated-extended-direction-tampered",
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-feedback-validation" {
		t.Fatalf("output = %#v, want feedback validation UNKNOWN", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirectionRequiresSource(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirectionInput{
		Feedback:       confirmedGeneratedExtendedCandidateRevisionFeedbackForDirection(t),
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "candidate-source" {
		t.Fatalf("output = %#v, want candidate-source UNKNOWN", output)
	}
}
