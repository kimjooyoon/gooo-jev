package decision

import "testing"

func extendedLineageDirectionForGeneration(t *testing.T, feedback ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedbackBinding) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirectionBinding {
	t.Helper()
	return DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirectionInput{
		Feedback:       feedback,
		CandidateSource: "gooo://candidate/generated-extended-lineage",
		NonAuthorizing: true,
	})
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(t *testing.T) {
	feedback := refutedGeneratedExtendedCandidateRevisionFeedbackForDirection(t)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageInput{
		Feedback:             feedback,
		Direction:            extendedLineageDirectionForGeneration(t, feedback),
		RevisionSource:       "gooo://revision/generated-extended-lineage-next",
		RevisionChangeDigest: "generated-extended-lineage-next-change",
		NonAuthorizing:       true,
	})
	if output.Status != "bound" || output.CandidateStatus != jevImprovementRevisionCandidateReady ||
		output.SourceCandidateID == "" || output.CandidateDigest == "" ||
		output.VerificationEvidenceDigest == "" || output.EvidencePrefixDigest == "" {
		t.Fatalf("output = %#v, want bound extended lineage candidate revision", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageRejectsReview(t *testing.T) {
	feedback := confirmedGeneratedExtendedCandidateRevisionFeedbackForDirection(t)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageInput{
		Feedback:             feedback,
		Direction:            extendedLineageDirectionForGeneration(t, feedback),
		RevisionSource:       "gooo://revision/generated-extended-lineage-next",
		RevisionChangeDigest: "generated-extended-lineage-next-change",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-directive" {
		t.Fatalf("output = %#v, want revision-directive UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageRejectsHold(t *testing.T) {
	feedback := unknownGeneratedExtendedCandidateRevisionFeedbackForDirection(t)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageInput{
		Feedback:             feedback,
		Direction:            extendedLineageDirectionForGeneration(t, feedback),
		RevisionSource:       "gooo://revision/generated-extended-lineage-next",
		RevisionChangeDigest: "generated-extended-lineage-next-change",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-directive" {
		t.Fatalf("output = %#v, want revision-directive UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageRequiresChange(t *testing.T) {
	feedback := refutedGeneratedExtendedCandidateRevisionFeedbackForDirection(t)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageInput{
		Feedback:       feedback,
		Direction:      extendedLineageDirectionForGeneration(t, feedback),
		RevisionSource: "gooo://revision/generated-extended-lineage-next",
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-change" {
		t.Fatalf("output = %#v, want revision-change UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageRejectsDirectionTampering(t *testing.T) {
	feedback := refutedGeneratedExtendedCandidateRevisionFeedbackForDirection(t)
	direction := extendedLineageDirectionForGeneration(t, feedback)
	direction.DirectionEvidenceDigest = "tampered"
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/generated-extended-lineage-tampered",
		RevisionChangeDigest: "generated-extended-lineage-tampered-change",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-direction-validation" {
		t.Fatalf("output = %#v, want direction-validation UNKNOWN", output)
	}
}
