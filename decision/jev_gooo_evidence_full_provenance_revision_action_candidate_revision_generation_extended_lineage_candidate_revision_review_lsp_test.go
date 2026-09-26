package decision

import "testing"

func supportRevisionCandidateReviewForLSP(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewBinding {
	t.Helper()
	return ClassifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReview(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewInput{
		Candidate:                nextRevisionCandidateForReview(t),
		ReviewChoice:             "support",
		ConfidenceEvidenceDigest: "confidence-band-lsp-support",
		ReviewEvidenceDigest:     "review-evidence-lsp-support",
		NonAuthorizing:           true,
	})
}

func abstainRevisionCandidateReviewForLSP(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewBinding {
	t.Helper()
	return ClassifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReview(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewInput{
		Candidate:                nextRevisionCandidateForReview(t),
		ReviewChoice:             "abstain",
		ConfidenceEvidenceDigest: "confidence-band-lsp-abstain",
		ReviewEvidenceDigest:     "review-evidence-lsp-abstain",
		NonAuthorizing:           true,
	})
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPSupport(t *testing.T) {
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPInput{
		Review:           supportRevisionCandidateReviewForLSP(t),
		ProjectionSource: "gooo://lsp/extended-lineage-candidate-revision-review-support",
		NonAuthorizing:   true,
	})
	if output.Status != "bound" || output.ProjectionStatus != "projected" ||
		output.Severity != "information" || output.EvidencePrefixDigest == "" {
		t.Fatalf("output = %#v, want informational LSP projection", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPAbstain(t *testing.T) {
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPInput{
		Review:           abstainRevisionCandidateReviewForLSP(t),
		ProjectionSource: "gooo://lsp/extended-lineage-candidate-revision-review-abstain",
		NonAuthorizing:   true,
	})
	if output.Status != "bound" || output.Severity != "hint" ||
		output.FallbackStage != "human-or-stronger-model-review" {
		t.Fatalf("output = %#v, want abstention fallback projection", output)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPRequiresSource(t *testing.T) {
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPInput{
		Review:         supportRevisionCandidateReviewForLSP(t),
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "projection-source" {
		t.Fatalf("output = %#v, want projection-source UNKNOWN", output)
	}
}

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPRejectsReviewTampering(t *testing.T) {
	review := supportRevisionCandidateReviewForLSP(t)
	review.ReviewEvidenceDigest = "tampered"
	output := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPInput{
		Review:           review,
		ProjectionSource: "gooo://lsp/extended-lineage-candidate-revision-review-tampered",
		NonAuthorizing:   true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-lineage-candidate-revision-review-validation" {
		t.Fatalf("output = %#v, want review validation UNKNOWN", output)
	}
}