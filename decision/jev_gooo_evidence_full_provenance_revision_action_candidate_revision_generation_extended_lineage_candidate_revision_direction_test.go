package decision

import "testing"

func confirmedExtendedLineageCandidateRevisionFeedbackForRevisionDirection(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionFeedbackBinding {
	t.Helper()
	metric := confirmedExtendedLineageCandidateRevisionMetricForRevisionFeedback(t)
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "extended-lineage-candidate-revision-direction-confirmed",
		NonAuthorizing:          true,
	})
}

func refutedExtendedLineageCandidateRevisionFeedbackForRevisionDirection(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionFeedbackBinding {
	t.Helper()
	metric := refutedExtendedLineageCandidateRevisionMetricForRevisionFeedback(t)
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "extended-lineage-candidate-revision-direction-refuted",
		NonAuthorizing:          true,
	})
}

func unknownExtendedLineageCandidateRevisionFeedbackForRevisionDirection(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionFeedbackBinding {
	t.Helper()
	metric := unknownExtendedLineageCandidateRevisionMetricForRevisionFeedback(t)
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "extended-lineage-candidate-revision-direction-unknown",
		NonAuthorizing:          true,
	})
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionReview(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionInput{
		Feedback:       confirmedExtendedLineageCandidateRevisionFeedbackForRevisionDirection(t),
		CandidateSource: "gooo://candidate/extended-lineage-candidate-revision-direction-confirmed",
		NonAuthorizing: true,
	})
	if output.Status != "bound" || output.Directive != jevImprovementDirectiveExternalReview ||
		output.SourceCandidateID == "" || output.ParentCandidateDigest == "" ||
		output.CandidateDigest == "" || output.EvidencePrefixDigest == "" {
		t.Fatalf("output = %#v, want external review direction", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionRevision(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionInput{
		Feedback:       refutedExtendedLineageCandidateRevisionFeedbackForRevisionDirection(t),
		CandidateSource: "gooo://candidate/extended-lineage-candidate-revision-direction-refuted",
		NonAuthorizing: true,
	})
	if output.Status != "bound" || output.Directive != jevImprovementDirectiveRevision ||
		output.DirectiveStatus != jevImprovementDirectiveRevision {
		t.Fatalf("output = %#v, want revision direction", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionHold(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionInput{
		Feedback:       unknownExtendedLineageCandidateRevisionFeedbackForRevisionDirection(t),
		CandidateSource: "gooo://candidate/extended-lineage-candidate-revision-direction-unknown",
		NonAuthorizing: true,
	})
	if output.Status != "bound" || output.Directive != jevImprovementDirectiveEvidence {
		t.Fatalf("output = %#v, want evidence hold direction", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionRejectsTampering(t *testing.T) {
	feedback := confirmedExtendedLineageCandidateRevisionFeedbackForRevisionDirection(t)
	feedback.FeedbackDigest = "tampered"
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionInput{
		Feedback:       feedback,
		CandidateSource: "gooo://candidate/extended-lineage-candidate-revision-direction-tampered",
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-lineage-candidate-revision-feedback-validation" {
		t.Fatalf("output = %#v, want feedback validation UNKNOWN", output)
	}
}

func TestDeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionRequiresSource(t *testing.T) {
	output := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionInput{
		Feedback:       confirmedExtendedLineageCandidateRevisionFeedbackForRevisionDirection(t),
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "candidate-source" {
		t.Fatalf("output = %#v, want candidate-source UNKNOWN", output)
	}
}