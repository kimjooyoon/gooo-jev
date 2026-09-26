package decision

import "testing"

func candidateFeedbackForRevisionGeneration(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedbackBinding {
	t.Helper()
	metric := refutedExtendedLineageCandidateMetricForFeedback(t)
	return BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "extended-lineage-candidate-revision-generation",
		NonAuthorizing:          true,
	})
}

func candidateDirectionForRevisionGeneration(t *testing.T, feedback ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedbackBinding) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateDirectionBinding {
	t.Helper()
	return DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateDirectionInput{
		Feedback:       feedback,
		CandidateSource: "gooo://candidate/extended-lineage-candidate-revision-generation",
		NonAuthorizing: true,
	})
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(t *testing.T) {
	feedback := candidateFeedbackForRevisionGeneration(t)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionInput{
		Feedback:             feedback,
		Direction:            candidateDirectionForRevisionGeneration(t, feedback),
		RevisionSource:       "gooo://revision/extended-lineage-candidate-revision",
		RevisionChangeDigest: "extended-lineage-candidate-revision-change",
		NonAuthorizing:       true,
	})
	if output.Status != "bound" || output.CandidateStatus != jevImprovementRevisionCandidateReady ||
		output.ParentCandidateDigest == "" || output.CandidateDigest == "" ||
		output.CandidateEvidenceDigest == "" || output.DirectionEvidenceDigest == "" {
		t.Fatalf("output = %#v, want bound extended lineage candidate revision", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionRejectsReview(t *testing.T) {
	metric := confirmedExtendedLineageCandidateMetricForFeedback(t)
	feedback := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "extended-lineage-candidate-revision-review",
		NonAuthorizing:          true,
	})
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionInput{
		Feedback:             feedback,
		Direction:            candidateDirectionForRevisionGeneration(t, feedback),
		RevisionSource:       "gooo://revision/extended-lineage-candidate-revision",
		RevisionChangeDigest: "extended-lineage-candidate-revision-change",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-directive" {
		t.Fatalf("output = %#v, want revision-directive UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionRejectsHold(t *testing.T) {
	feedback := unknownExtendedLineageCandidateFeedbackForDirection(t)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionInput{
		Feedback:             feedback,
		Direction:            candidateDirectionForRevisionGeneration(t, feedback),
		RevisionSource:       "gooo://revision/extended-lineage-candidate-revision",
		RevisionChangeDigest: "extended-lineage-candidate-revision-change",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-directive" {
		t.Fatalf("output = %#v, want revision-directive UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionRequiresChange(t *testing.T) {
	feedback := candidateFeedbackForRevisionGeneration(t)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionInput{
		Feedback:       feedback,
		Direction:      candidateDirectionForRevisionGeneration(t, feedback),
		RevisionSource: "gooo://revision/extended-lineage-candidate-revision",
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-change" {
		t.Fatalf("output = %#v, want revision-change UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionRejectsDirectionTampering(t *testing.T) {
	feedback := candidateFeedbackForRevisionGeneration(t)
	direction := candidateDirectionForRevisionGeneration(t, feedback)
	direction.DirectionEvidenceDigest = "tampered"
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/extended-lineage-candidate-revision-tampered",
		RevisionChangeDigest: "extended-lineage-candidate-revision-tampered-change",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-lineage-candidate-direction-validation" {
		t.Fatalf("output = %#v, want direction-validation UNKNOWN", output)
	}
}