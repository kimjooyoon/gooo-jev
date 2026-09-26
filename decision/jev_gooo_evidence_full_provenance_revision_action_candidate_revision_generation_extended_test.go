package decision

import "testing"

func generatedCandidateRevisionDirectionForExtendedGeneration(t *testing.T, feedback ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationFeedbackBinding) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationDirectionBinding {
	t.Helper()
	return DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationDirectionInput{
		Feedback:       feedback,
		CandidateSource: "gooo://candidate/generated-revision-generation",
		NonAuthorizing: true,
	})
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(t *testing.T) {
	feedback := refutedGeneratedCandidateRevisionFeedbackForDirection(t)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedInput{
		Feedback:             feedback,
		Direction:            generatedCandidateRevisionDirectionForExtendedGeneration(t, feedback),
		RevisionSource:       "gooo://revision/generated-candidate-next",
		RevisionChangeDigest: "generated-candidate-next-revision-change",
		NonAuthorizing:       true,
	})
	if output.Status != "bound" || output.CandidateStatus != jevImprovementRevisionCandidateReady ||
		output.SourceCandidateID == "" || output.CandidateDigest == "" ||
		output.VerificationEvidenceDigest == "" || output.EvidencePrefixDigest == "" {
		t.Fatalf("output = %#v, want bound generated candidate revision", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedRejectsReview(t *testing.T) {
	feedback := confirmedGeneratedCandidateRevisionFeedbackForDirection(t)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedInput{
		Feedback:             feedback,
		Direction:            generatedCandidateRevisionDirectionForExtendedGeneration(t, feedback),
		RevisionSource:       "gooo://revision/generated-candidate-next",
		RevisionChangeDigest: "generated-candidate-next-revision-change",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-directive" {
		t.Fatalf("output = %#v, want revision-directive UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedRejectsHold(t *testing.T) {
	feedback := unknownGeneratedCandidateRevisionFeedbackForDirection(t)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedInput{
		Feedback:             feedback,
		Direction:            generatedCandidateRevisionDirectionForExtendedGeneration(t, feedback),
		RevisionSource:       "gooo://revision/generated-candidate-next",
		RevisionChangeDigest: "generated-candidate-next-revision-change",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-directive" {
		t.Fatalf("output = %#v, want revision-directive UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedRequiresChange(t *testing.T) {
	feedback := refutedGeneratedCandidateRevisionFeedbackForDirection(t)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedInput{
		Feedback:       feedback,
		Direction:      generatedCandidateRevisionDirectionForExtendedGeneration(t, feedback),
		RevisionSource: "gooo://revision/generated-candidate-next",
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-change" {
		t.Fatalf("output = %#v, want revision-change UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedRejectsDirectionTampering(t *testing.T) {
	feedback := refutedGeneratedCandidateRevisionFeedbackForDirection(t)
	direction := generatedCandidateRevisionDirectionForExtendedGeneration(t, feedback)
	direction.DirectionEvidenceDigest = "tampered"
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtended(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/generated-candidate-tampered",
		RevisionChangeDigest: "generated-candidate-tampered-change",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-direction-validation" {
		t.Fatalf("output = %#v, want direction-validation UNKNOWN", output)
	}
}
