package decision

import "testing"

func supportRevisionCandidateReviewLSPForConfidence(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPBinding {
	t.Helper()
	return ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPInput{
		Review:           supportRevisionCandidateReviewForLSP(t),
		ProjectionSource: "gooo://lsp/extended-lineage-candidate-revision-confidence-support",
		NonAuthorizing:   true,
	})
}

func abstainRevisionCandidateReviewLSPForConfidence(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPBinding {
	t.Helper()
	return ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPInput{
		Review:           abstainRevisionCandidateReviewForLSP(t),
		ProjectionSource: "gooo://lsp/extended-lineage-candidate-revision-confidence-abstain",
		NonAuthorizing:   true,
	})
}

func TestCalibrateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceLow(t *testing.T) {
	output := CalibrateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidence(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceInput{
		Review:                    supportRevisionCandidateReviewLSPForConfidence(t),
		ConfidenceBand:            "low",
		CalibrationEvidenceDigest: "calibration-evidence-low",
		NonAuthorizing:            true,
	})
	if output.Status != "bound" || output.AutomationMode != "shadow-only" ||
		output.FallbackStage != "collect-more-evidence" {
		t.Fatalf("output = %#v, want shadow-only calibration", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestCalibrateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceHigh(t *testing.T) {
	output := CalibrateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidence(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceInput{
		Review:                    abstainRevisionCandidateReviewLSPForConfidence(t),
		ConfidenceBand:            "high",
		CalibrationEvidenceDigest: "calibration-evidence-high",
		NonAuthorizing:            true,
	})
	if output.Status != "bound" || output.AutomationMode != "deterministic-policy-required" ||
		output.FallbackStage != "human-or-deterministic-policy" {
		t.Fatalf("output = %#v, want deterministic policy fallback", output)
	}
}

func TestCalibrateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceRejectsUnknownBand(t *testing.T) {
	output := CalibrateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidence(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceInput{
		Review:                    supportRevisionCandidateReviewLSPForConfidence(t),
		ConfidenceBand:            "certain",
		CalibrationEvidenceDigest: "calibration-evidence-invalid-band",
		NonAuthorizing:            true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "confidence-band" {
		t.Fatalf("output = %#v, want confidence-band UNKNOWN", output)
	}
}

func TestCalibrateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceRequiresEvidence(t *testing.T) {
	output := CalibrateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidence(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceInput{
		Review:         supportRevisionCandidateReviewLSPForConfidence(t),
		ConfidenceBand: "medium",
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "calibration-evidence" {
		t.Fatalf("output = %#v, want calibration-evidence UNKNOWN", output)
	}
}

func TestCalibrateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceRejectsReviewTampering(t *testing.T) {
	review := supportRevisionCandidateReviewLSPForConfidence(t)
	review.ReviewEvidenceDigest = "tampered"
	output := CalibrateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidence(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceInput{
		Review:                    review,
		ConfidenceBand:            "medium",
		CalibrationEvidenceDigest: "calibration-evidence-tampered",
		NonAuthorizing:            true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-lineage-candidate-revision-review-lsp-validation" {
		t.Fatalf("output = %#v, want review LSP validation UNKNOWN", output)
	}
}