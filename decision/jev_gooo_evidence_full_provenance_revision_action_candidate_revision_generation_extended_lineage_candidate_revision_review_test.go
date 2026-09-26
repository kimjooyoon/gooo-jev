package decision

import "testing"

func nextRevisionCandidateForReview(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextBinding {
	t.Helper()
	direction := refutedExtendedLineageCandidateRevisionDirectionForNext(t)
	return GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNext(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextInput{
		Direction:                       direction,
		PreviousCandidateDigest:         direction.DirectionCandidateDigest,
		PreviousCandidateEvidenceDigest: direction.CandidateEvidenceDigest,
		CandidateSource:                 "gooo://candidate/extended-lineage-candidate-revision-review/generated",
		RevisionSource:                  "gooo://revision/extended-lineage-candidate-revision-review",
		RevisionChangeDigest:            "change-next-revision-candidate-review",
		NonAuthorizing:                  true,
	})
}

func TestClassifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewSupport(t *testing.T) {
	candidate := nextRevisionCandidateForReview(t)
	output := ClassifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReview(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewInput{
		Candidate:                candidate,
		ReviewChoice:             "support",
		ConfidenceEvidenceDigest: "confidence-band-review-support",
		ReviewEvidenceDigest:     "review-evidence-support",
		NonAuthorizing:           true,
	})
	if output.Status != "bound" || output.ReviewDisposition != "review-supported" ||
		output.FallbackStage != "deterministic-policy-remains-authoritative" {
		t.Fatalf("output = %#v, want non-authorizing support disposition", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestClassifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionReviewAbstain(t *testing.T) {
	candidate := nextRevisionCandidateForReview(t)
	output := ClassifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReview(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewInput{
		Candidate:                candidate,
		ReviewChoice:             "abstain",
		ConfidenceEvidenceDigest: "confidence-band-review-abstain",
		ReviewEvidenceDigest:     "review-evidence-abstain",
		NonAuthorizing:           true,
	})
	if output.Status != "bound" || output.ReviewDisposition != "review-abstained" ||
		output.FallbackStage != "human-or-stronger-model-review" {
		t.Fatalf("output = %#v, want abstain fallback", output)
	}
}

func TestClassifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionReviewRejectsChoice(t *testing.T) {
	output := ClassifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReview(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewInput{
		Candidate:                nextRevisionCandidateForReview(t),
		ReviewChoice:             "maybe",
		ConfidenceEvidenceDigest: "confidence-band-review-invalid",
		ReviewEvidenceDigest:     "review-evidence-invalid",
		NonAuthorizing:           true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "review-choice" {
		t.Fatalf("output = %#v, want review-choice UNKNOWN", output)
	}
}

func TestClassifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionReviewRequiresEvidence(t *testing.T) {
	output := ClassifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionReview(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewInput{
		Candidate:                nextRevisionCandidateForReview(t),
		ReviewChoice:             "reject",
		ReviewEvidenceDigest:     "review-evidence-reject",
		NonAuthorizing:           true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "confidence-evidence" {
		t.Fatalf("output = %#v, want confidence-evidence UNKNOWN", output)
	}
}

func TestClassifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionReviewRejectsCandidateTampering(t *testing.T) {
	candidate := nextRevisionCandidateForReview(t)
	candidate.CandidateEvidenceDigest = "tampered"
	output := ClassifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionReview(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionReviewInput{
		Candidate:                candidate,
		ReviewChoice:             "reject",
		ConfidenceEvidenceDigest: "confidence-band-review-tampered",
		ReviewEvidenceDigest:     "review-evidence-tampered",
		NonAuthorizing:           true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-lineage-candidate-revision-next-validation" {
		t.Fatalf("output = %#v, want candidate validation UNKNOWN", output)
	}
}