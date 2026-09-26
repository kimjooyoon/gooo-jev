package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewInput
// classifies a generated revision candidate into a closed, non-authorizing review disposition.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewInput struct {
	Candidate                  ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionNextBinding
	ReviewChoice               string
	ConfidenceEvidenceDigest   string
	ReviewEvidenceDigest       string
	NonAuthorizing             bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewBinding
// records a Jev-like choice without treating confidence as permission or safety proof.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewBinding struct {
	Status                       string
	MissingStage                 string
	CandidateStatus              string
	ParentCandidateDigest        string
	CandidateDigest              string
	CandidateEvidenceDigest      string
	CandidateSource              string
	RevisionSource               string
	BoundRevisionChangeDigest    string
	ReviewChoice                 string
	ReviewDisposition            string
	ConfidenceEvidenceDigest     string
	ReviewEvidenceDigest         string
	FallbackStage                string
	DirectionEvidenceDigest      string
	ReverseEvidenceDigest        string
	EvidencePrefixDigest         string
	MetricEvidenceDigest         string
	FeedbackDigest               string
	AggregationEvidenceDigest    string
	EvidenceDigest               string
	NonExecuting                 bool
	NonAuthorizing               bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
		b.CandidateStatus != jevImprovementRevisionCandidateReady ||
		b.ParentCandidateDigest == "" ||
		b.CandidateDigest == "" ||
		b.CandidateEvidenceDigest == "" ||
		b.CandidateSource == "" ||
		b.RevisionSource == "" ||
		b.BoundRevisionChangeDigest == "" ||
		b.ReviewChoice == "" ||
		b.ReviewDisposition == "" ||
		b.ConfidenceEvidenceDigest == "" ||
		b.ReviewEvidenceDigest == "" ||
		b.FallbackStage == "" ||
		b.DirectionEvidenceDigest == "" ||
		b.ReverseEvidenceDigest == "" ||
		b.EvidencePrefixDigest == "" ||
		b.MetricEvidenceDigest == "" ||
		b.FeedbackDigest == "" ||
		b.AggregationEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended lineage candidate revision review binding")
	}
	expectedDisposition, expectedFallback := reviewDispositionForGoooExtendedLineageCandidateRevision(b.ReviewChoice)
	if expectedDisposition == "" ||
		b.ReviewDisposition != expectedDisposition ||
		b.FallbackStage != expectedFallback {
		return fmt.Errorf("invalid Gooo extended lineage candidate revision review disposition")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage candidate revision review must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReview(b)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended lineage candidate revision review digest mismatch")
	}
	return nil
}

// ClassifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReview
// preserves a closed review choice and leaves deterministic execution policy authoritative.
func ClassifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReview(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-generation-extended-lineage-candidate-revision-review"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewBinding{
			Status:         "UNKNOWN",
			MissingStage:   stage,
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.Candidate.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Candidate.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Candidate.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-revision-next-validation")
	}
	disposition, fallback := reviewDispositionForGoooExtendedLineageCandidateRevision(input.ReviewChoice)
	if disposition == "" {
		return unknown("review-choice")
	}
	if strings.TrimSpace(input.ConfidenceEvidenceDigest) == "" {
		return unknown("confidence-evidence")
	}
	if strings.TrimSpace(input.ReviewEvidenceDigest) == "" {
		return unknown("review-evidence")
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewBinding{
		Status:                    "bound",
		CandidateStatus:           input.Candidate.CandidateStatus,
		ParentCandidateDigest:     input.Candidate.ParentCandidateDigest,
		CandidateDigest:           input.Candidate.CandidateDigest,
		CandidateEvidenceDigest:   input.Candidate.CandidateEvidenceDigest,
		CandidateSource:           input.Candidate.CandidateSource,
		RevisionSource:             input.Candidate.RevisionSource,
		BoundRevisionChangeDigest: input.Candidate.BoundRevisionChangeDigest,
		ReviewChoice:              input.ReviewChoice,
		ReviewDisposition:         disposition,
		ConfidenceEvidenceDigest:  input.ConfidenceEvidenceDigest,
		ReviewEvidenceDigest:      input.ReviewEvidenceDigest,
		FallbackStage:             fallback,
		DirectionEvidenceDigest:   input.Candidate.DirectionEvidenceDigest,
		ReverseEvidenceDigest:     input.Candidate.ReverseEvidenceDigest,
		EvidencePrefixDigest:      input.Candidate.EvidencePrefixDigest,
		MetricEvidenceDigest:      input.Candidate.MetricEvidenceDigest,
		FeedbackDigest:            input.Candidate.FeedbackDigest,
		AggregationEvidenceDigest: input.Candidate.AggregationEvidenceDigest,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReview(output)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-revision-review-evidence")
	}
	return output
}

func reviewDispositionForGoooExtendedLineageCandidateRevision(choice string) (string, string) {
	switch strings.TrimSpace(choice) {
	case "support":
		return "review-supported", "deterministic-policy-remains-authoritative"
	case "reject":
		return "review-rejected", "deterministic-policy-remains-authoritative"
	case "abstain":
		return "review-abstained", "human-or-stronger-model-review"
	default:
		return "", ""
	}
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReview(b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionReviewBinding) string {
	digest, err := Digest(struct {
		Status                    string
		CandidateStatus           string
		ParentCandidateDigest     string
		CandidateDigest           string
		CandidateEvidenceDigest   string
		CandidateSource            string
		RevisionSource             string
		BoundRevisionChangeDigest string
		ReviewChoice              string
		ReviewDisposition         string
		ConfidenceEvidenceDigest  string
		ReviewEvidenceDigest      string
		FallbackStage             string
		DirectionEvidenceDigest   string
		ReverseEvidenceDigest     string
		EvidencePrefixDigest      string
		MetricEvidenceDigest      string
		FeedbackDigest            string
		AggregationEvidenceDigest string
		NonExecuting              bool
		NonAuthorizing            bool
	}{
		Status:                    b.Status,
		CandidateStatus:           b.CandidateStatus,
		ParentCandidateDigest:     b.ParentCandidateDigest,
		CandidateDigest:           b.CandidateDigest,
		CandidateEvidenceDigest:   b.CandidateEvidenceDigest,
		CandidateSource:           b.CandidateSource,
		RevisionSource:             b.RevisionSource,
		BoundRevisionChangeDigest: b.BoundRevisionChangeDigest,
		ReviewChoice:              b.ReviewChoice,
		ReviewDisposition:         b.ReviewDisposition,
		ConfidenceEvidenceDigest:  b.ConfidenceEvidenceDigest,
		ReviewEvidenceDigest:      b.ReviewEvidenceDigest,
		FallbackStage:             b.FallbackStage,
		DirectionEvidenceDigest:   b.DirectionEvidenceDigest,
		ReverseEvidenceDigest:     b.ReverseEvidenceDigest,
		EvidencePrefixDigest:      b.EvidencePrefixDigest,
		MetricEvidenceDigest:      b.MetricEvidenceDigest,
		FeedbackDigest:            b.FeedbackDigest,
		AggregationEvidenceDigest: b.AggregationEvidenceDigest,
		NonExecuting:              b.NonExecuting,
		NonAuthorizing:            b.NonAuthorizing,
	})
	if err != nil {
		return ""
	}
	return digest
}