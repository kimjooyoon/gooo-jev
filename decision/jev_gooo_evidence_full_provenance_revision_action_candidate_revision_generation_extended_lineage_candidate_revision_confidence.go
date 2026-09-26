package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceInput
// calibrates review confidence into a bounded automation mode without making it permission.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceInput struct {
	Review                 ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPBinding
	ConfidenceBand         string
	CalibrationEvidenceDigest string
	NonAuthorizing         bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceBinding
// records confidence bands, fallback, and an explicitly non-authorizing automation mode.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceBinding struct {
	Status                    string
	MissingStage              string
	CalibrationStatus         string
	CandidateStatus           string
	CandidateDigest           string
	CandidateSource           string
	RevisionSource            string
	ReviewDisposition          string
	ConfidenceBand             string
	AutomationMode             string
	CalibrationEvidenceDigest  string
	FallbackStage              string
	EvidencePrefixDigest       string
	ReviewEvidenceDigest       string
	ProjectionSource           string
	EvidenceDigest             string
	NonExecuting               bool
	NonAuthorizing             bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
		b.CalibrationStatus != "calibrated" ||
		b.CandidateStatus != jevImprovementRevisionCandidateReady ||
		b.CandidateDigest == "" ||
		b.CandidateSource == "" ||
		b.RevisionSource == "" ||
		b.ReviewDisposition == "" ||
		b.ConfidenceBand == "" ||
		b.AutomationMode == "" ||
		b.CalibrationEvidenceDigest == "" ||
		b.FallbackStage == "" ||
		b.EvidencePrefixDigest == "" ||
		b.ReviewEvidenceDigest == "" ||
		b.ProjectionSource == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended lineage candidate revision confidence binding")
	}
	expectedMode, expectedFallback := automationForGoooExtendedLineageCandidateRevisionConfidence(b.ConfidenceBand)
	if expectedMode == "" ||
		b.AutomationMode != expectedMode ||
		b.FallbackStage != expectedFallback {
		return fmt.Errorf("invalid Gooo extended lineage candidate revision confidence mode")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage candidate revision confidence must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidence(b)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended lineage candidate revision confidence digest mismatch")
	}
	return nil
}

// CalibrateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidence
// treats confidence as an input to a bounded mode, never as an authorization proof.
func CalibrateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidence(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-generation-extended-lineage-candidate-revision-confidence"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceBinding{
			Status:         "UNKNOWN",
			MissingStage:   stage,
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.Review.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Review.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Review.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-revision-review-lsp-validation")
	}
	mode, fallback := automationForGoooExtendedLineageCandidateRevisionConfidence(input.ConfidenceBand)
	if mode == "" {
		return unknown("confidence-band")
	}
	if strings.TrimSpace(input.CalibrationEvidenceDigest) == "" {
		return unknown("calibration-evidence")
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceBinding{
		Status:                   "bound",
		CalibrationStatus:        "calibrated",
		CandidateStatus:          input.Review.CandidateStatus,
		CandidateDigest:          input.Review.CandidateDigest,
		CandidateSource:          input.Review.CandidateSource,
		RevisionSource:            input.Review.RevisionSource,
		ReviewDisposition:         input.Review.ReviewDisposition,
		ConfidenceBand:            strings.TrimSpace(input.ConfidenceBand),
		AutomationMode:            mode,
		CalibrationEvidenceDigest: input.CalibrationEvidenceDigest,
		FallbackStage:             fallback,
		EvidencePrefixDigest:      input.Review.EvidencePrefixDigest,
		ReviewEvidenceDigest:      input.Review.ReviewEvidenceDigest,
		ProjectionSource:          input.Review.ProjectionSource,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidence(output)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-revision-confidence-evidence")
	}
	return output
}

func automationForGoooExtendedLineageCandidateRevisionConfidence(band string) (string, string) {
	switch strings.TrimSpace(band) {
	case "low":
		return "shadow-only", "collect-more-evidence"
	case "medium":
		return "reversible-review-only", "deterministic-reversible-policy"
	case "high":
		return "deterministic-policy-required", "human-or-deterministic-policy"
	default:
		return "", ""
	}
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidence(b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionConfidenceBinding) string {
	digest, err := Digest(struct {
		Status                    string
		CalibrationStatus         string
		CandidateStatus           string
		CandidateDigest           string
		CandidateSource           string
		RevisionSource            string
		ReviewDisposition          string
		ConfidenceBand             string
		AutomationMode             string
		CalibrationEvidenceDigest  string
		FallbackStage              string
		EvidencePrefixDigest       string
		ReviewEvidenceDigest       string
		ProjectionSource           string
		NonExecuting               bool
		NonAuthorizing             bool
	}{
		Status:                   b.Status,
		CalibrationStatus:        b.CalibrationStatus,
		CandidateStatus:          b.CandidateStatus,
		CandidateDigest:          b.CandidateDigest,
		CandidateSource:          b.CandidateSource,
		RevisionSource:            b.RevisionSource,
		ReviewDisposition:         b.ReviewDisposition,
		ConfidenceBand:            b.ConfidenceBand,
		AutomationMode:            b.AutomationMode,
		CalibrationEvidenceDigest: b.CalibrationEvidenceDigest,
		FallbackStage:             b.FallbackStage,
		EvidencePrefixDigest:      b.EvidencePrefixDigest,
		ReviewEvidenceDigest:      b.ReviewEvidenceDigest,
		ProjectionSource:          b.ProjectionSource,
		NonExecuting:              b.NonExecuting,
		NonAuthorizing:            b.NonAuthorizing,
	})
	if err != nil {
		return ""
	}
	return digest
}