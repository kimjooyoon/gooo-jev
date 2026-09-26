package decision

import "testing"

func extendedLineageCandidateDirectionForGeneration(t *testing.T, feedback ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackBinding) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageDirectionBinding {
	t.Helper()
	return DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageDirectionInput{
		Feedback:       feedback,
		CandidateSource: "gooo://candidate/extended-lineage-generation",
		NonAuthorizing: true,
	})
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidate(t *testing.T) {
	feedback := refutedExtendedLineageCandidateRevisionFeedbackForDirection(t)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateInput{
		Feedback:             feedback,
		Direction:            extendedLineageCandidateDirectionForGeneration(t, feedback),
		RevisionSource:       "gooo://revision/extended-lineage-next",
		RevisionChangeDigest: "extended-lineage-next-change",
		NonAuthorizing:       true,
	})
	if output.Status != "bound" || output.CandidateStatus != jevImprovementRevisionCandidateReady ||
		output.SourceCandidateID == "" || output.CandidateDigest == "" ||
		output.VerificationEvidenceDigest == "" || output.EvidencePrefixDigest == "" {
		t.Fatalf("output = %#v, want bound extended lineage candidate", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRejectsReview(t *testing.T) {
	feedback := confirmedExtendedLineageCandidateRevisionFeedbackForDirection(t)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateInput{
		Feedback:             feedback,
		Direction:            extendedLineageCandidateDirectionForGeneration(t, feedback),
		RevisionSource:       "gooo://revision/extended-lineage-next",
		RevisionChangeDigest: "extended-lineage-next-change",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-directive" {
		t.Fatalf("output = %#v, want revision-directive UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRejectsHold(t *testing.T) {
	feedback := unknownExtendedLineageCandidateRevisionFeedbackForDirection(t)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateInput{
		Feedback:             feedback,
		Direction:            extendedLineageCandidateDirectionForGeneration(t, feedback),
		RevisionSource:       "gooo://revision/extended-lineage-next",
		RevisionChangeDigest: "extended-lineage-next-change",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-directive" {
		t.Fatalf("output = %#v, want revision-directive UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRequiresChange(t *testing.T) {
	feedback := refutedExtendedLineageCandidateRevisionFeedbackForDirection(t)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateInput{
		Feedback:       feedback,
		Direction:      extendedLineageCandidateDirectionForGeneration(t, feedback),
		RevisionSource: "gooo://revision/extended-lineage-next",
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-change" {
		t.Fatalf("output = %#v, want revision-change UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRejectsDirectionTampering(t *testing.T) {
	feedback := refutedExtendedLineageCandidateRevisionFeedbackForDirection(t)
	direction := extendedLineageCandidateDirectionForGeneration(t, feedback)
	direction.DirectionEvidenceDigest = "tampered"
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateInput{
		Feedback:             feedback,
		Direction:            direction,
		RevisionSource:       "gooo://revision/extended-lineage-tampered",
		RevisionChangeDigest: "extended-lineage-tampered-change",
		NonAuthorizing:       true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-lineage-direction-validation" {
		t.Fatalf("output = %#v, want direction-validation UNKNOWN", output)
	}
}
