package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPInput
// projects a closed review disposition into LSP metadata without adding execution authority.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPInput struct {
	Review           ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewBinding
	ProjectionSource string
	NonAuthorizing   bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPBinding
// preserves review, confidence, fallback, and evidence-prefix provenance for editor clients.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPBinding struct {
	Status                    string
	MissingStage              string
	ProjectionStatus          string
	CandidateStatus           string
	CandidateDigest           string
	CandidateSource           string
	RevisionSource            string
	ReviewChoice              string
	ReviewDisposition         string
	Severity                  string
	FallbackStage             string
	ConfidenceEvidenceDigest  string
	ReviewEvidenceDigest      string
	EvidencePrefixDigest      string
	ProjectionSource          string
	EvidenceDigest            string
	NonExecuting              bool
	NonAuthorizing            bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
		b.ProjectionStatus != "projected" ||
		b.CandidateStatus != jevImprovementRevisionCandidateReady ||
		b.CandidateDigest == "" ||
		b.CandidateSource == "" ||
		b.RevisionSource == "" ||
		b.ReviewChoice == "" ||
		b.ReviewDisposition == "" ||
		b.Severity == "" ||
		b.FallbackStage == "" ||
		b.ConfidenceEvidenceDigest == "" ||
		b.ReviewEvidenceDigest == "" ||
		b.EvidencePrefixDigest == "" ||
		b.ProjectionSource == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended lineage candidate revision review LSP binding")
	}
	expectedDisposition, expectedFallback := reviewDispositionForGoooExtendedLineageCandidateRevision(b.ReviewChoice)
	if expectedDisposition == "" ||
		b.ReviewDisposition != expectedDisposition ||
		b.FallbackStage != expectedFallback {
		return fmt.Errorf("invalid Gooo extended lineage candidate revision review LSP disposition")
	}
	expectedSeverity := severityForGoooExtendedLineageCandidateRevisionReview(b.ReviewChoice)
	if b.Severity != expectedSeverity {
		return fmt.Errorf("invalid Gooo extended lineage candidate revision review LSP severity")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage candidate revision review LSP must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSP(b)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended lineage candidate revision review LSP digest mismatch")
	}
	return nil
}

// ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSP
// projects review state for LSP consumers while keeping deterministic policy authoritative.
func ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSP(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-generation-extended-lineage-candidate-revision-review-lsp"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPBinding{
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
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-revision-review-validation")
	}
	if strings.TrimSpace(input.ProjectionSource) == "" {
		return unknown("projection-source")
	}
	severity := severityForGoooExtendedLineageCandidateRevisionReview(input.Review.ReviewChoice)
	if severity == "" {
		return unknown("review-severity")
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPBinding{
		Status:                   "bound",
		ProjectionStatus:         "projected",
		CandidateStatus:          input.Review.CandidateStatus,
		CandidateDigest:          input.Review.CandidateDigest,
		CandidateSource:          input.Review.CandidateSource,
		RevisionSource:            input.Review.RevisionSource,
		ReviewChoice:              input.Review.ReviewChoice,
		ReviewDisposition:         input.Review.ReviewDisposition,
		Severity:                 severity,
		FallbackStage:             input.Review.FallbackStage,
		ConfidenceEvidenceDigest: input.Review.ConfidenceEvidenceDigest,
		ReviewEvidenceDigest:      input.Review.ReviewEvidenceDigest,
		EvidencePrefixDigest:      input.Review.EvidencePrefixDigest,
		ProjectionSource:          input.ProjectionSource,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSP(output)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-revision-review-lsp-evidence")
	}
	return output
}

func severityForGoooExtendedLineageCandidateRevisionReview(choice string) string {
	switch strings.TrimSpace(choice) {
	case "support":
		return "information"
	case "reject":
		return "warning"
	case "abstain":
		return "hint"
	default:
		return ""
	}
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSP(b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewLSPBinding) string {
	digest, err := Digest(struct {
		Status                   string
		ProjectionStatus         string
		CandidateStatus          string
		CandidateDigest          string
		CandidateSource          string
		RevisionSource           string
		ReviewChoice             string
		ReviewDisposition        string
		Severity                 string
		FallbackStage            string
		ConfidenceEvidenceDigest string
		ReviewEvidenceDigest     string
		EvidencePrefixDigest     string
		ProjectionSource         string
		NonExecuting             bool
		NonAuthorizing           bool
	}{
		Status:                   b.Status,
		ProjectionStatus:         b.ProjectionStatus,
		CandidateStatus:           b.CandidateStatus,
		CandidateDigest:           b.CandidateDigest,
		CandidateSource:           b.CandidateSource,
		RevisionSource:            b.RevisionSource,
		ReviewChoice:              b.ReviewChoice,
		ReviewDisposition:         b.ReviewDisposition,
		Severity:                 b.Severity,
		FallbackStage:             b.FallbackStage,
		ConfidenceEvidenceDigest: b.ConfidenceEvidenceDigest,
		ReviewEvidenceDigest:     b.ReviewEvidenceDigest,
		EvidencePrefixDigest:      b.EvidencePrefixDigest,
		ProjectionSource:          b.ProjectionSource,
		NonExecuting:              b.NonExecuting,
		NonAuthorizing:            b.NonAuthorizing,
	})
	if err != nil {
		return ""
	}
	return digest
}