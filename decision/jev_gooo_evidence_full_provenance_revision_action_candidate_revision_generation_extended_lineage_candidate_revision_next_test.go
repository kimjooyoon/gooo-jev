package decision

import "testing"

func refutedExtendedLineageCandidateRevisionDirectionForNext(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionBinding {
	t.Helper()
	return DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionInput{
		Feedback:       refutedExtendedLineageCandidateRevisionFeedbackForRevisionDirection(t),
		CandidateSource: "gooo://candidate/extended-lineage-candidate-revision-next",
		NonAuthorizing: true,
	})
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextRevision(t *testing.T) {
	direction := refutedExtendedLineageCandidateRevisionDirectionForNext(t)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNext(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextInput{
		Direction:                       direction,
		PreviousCandidateDigest:         direction.DirectionCandidateDigest,
		PreviousCandidateEvidenceDigest: direction.CandidateEvidenceDigest,
		CandidateSource:                 "gooo://candidate/extended-lineage-candidate-revision-next/generated",
		RevisionSource:                  "gooo://revision/extended-lineage-candidate-revision-next",
		RevisionChangeDigest:            "change-next-revision-candidate",
		NonAuthorizing:                  true,
	})
	if output.Status != "bound" || output.CandidateStatus != jevImprovementRevisionCandidateReady ||
		output.ParentCandidateDigest != output.PreviousCandidateDigest ||
		output.CandidateDigest == "" || output.EvidenceDigest == "" {
		t.Fatalf("output = %#v, want bound next revision candidate", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextRejectsDirectionTampering(t *testing.T) {
	direction := refutedExtendedLineageCandidateRevisionDirectionForNext(t)
	direction.DirectionEvidenceDigest = "tampered"
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNext(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextInput{
		Direction:                       direction,
		PreviousCandidateDigest:         direction.DirectionCandidateDigest,
		PreviousCandidateEvidenceDigest: direction.CandidateEvidenceDigest,
		CandidateSource:                 "gooo://candidate/extended-lineage-candidate-revision-next/tampered",
		RevisionSource:                  "gooo://revision/extended-lineage-candidate-revision-next/tampered",
		RevisionChangeDigest:            "change-next-revision-candidate-tampered",
		NonAuthorizing:                  true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-lineage-candidate-revision-direction-validation" {
		t.Fatalf("output = %#v, want direction validation UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextRejectsParentMismatch(t *testing.T) {
	direction := refutedExtendedLineageCandidateRevisionDirectionForNext(t)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNext(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextInput{
		Direction:                       direction,
		PreviousCandidateDigest:         "wrong-parent",
		PreviousCandidateEvidenceDigest: direction.CandidateEvidenceDigest,
		CandidateSource:                 "gooo://candidate/extended-lineage-candidate-revision-next/mismatch",
		RevisionSource:                  "gooo://revision/extended-lineage-candidate-revision-next/mismatch",
		RevisionChangeDigest:            "change-next-revision-candidate-mismatch",
		NonAuthorizing:                  true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "previous-candidate" {
		t.Fatalf("output = %#v, want previous-candidate UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextRequiresRevisionChange(t *testing.T) {
	direction := refutedExtendedLineageCandidateRevisionDirectionForNext(t)
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNext(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextInput{
		Direction:                       direction,
		PreviousCandidateDigest:         direction.DirectionCandidateDigest,
		PreviousCandidateEvidenceDigest: direction.CandidateEvidenceDigest,
		CandidateSource:                 "gooo://candidate/extended-lineage-candidate-revision-next/missing-change",
		RevisionSource:                  "gooo://revision/extended-lineage-candidate-revision-next/missing-change",
		NonAuthorizing:                  true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-change" {
		t.Fatalf("output = %#v, want revision-change UNKNOWN", output)
	}
}

func TestGenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextRejectsNonRevisionDirection(t *testing.T) {
	direction := DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirection(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionDirectionInput{
		Feedback:       confirmedExtendedLineageCandidateRevisionFeedbackForRevisionDirection(t),
		CandidateSource: "gooo://candidate/extended-lineage-candidate-revision-next/external-review",
		NonAuthorizing: true,
	})
	output := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNext(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextInput{
		Direction:                       direction,
		PreviousCandidateDigest:         direction.DirectionCandidateDigest,
		PreviousCandidateEvidenceDigest: direction.CandidateEvidenceDigest,
		CandidateSource:                 "gooo://candidate/extended-lineage-candidate-revision-next/external-review/generated",
		RevisionSource:                  "gooo://revision/extended-lineage-candidate-revision-next/external-review",
		RevisionChangeDigest:            "change-next-revision-candidate-review",
		NonAuthorizing:                  true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-directive" {
		t.Fatalf("output = %#v, want revision-directive UNKNOWN", output)
	}
}