package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageInput
// connects extended candidate direction to the next revision candidate.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageInput struct {
	Feedback             ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedbackBinding
	Direction            ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedDirectionBinding
	RevisionSource       string
	RevisionChangeDigest string
	NonAuthorizing       bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageBinding
// preserves extended candidate lineage and next-candidate provenance.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageBinding struct {
	Status                          string
	MissingStage                    string
	CandidateID                     string
	SourceCandidateID               string
	SourceRevisionCandidateDigest   string
	SourceCandidateEvidenceDigest   string
	RevisionCandidateDigest         string
	RevisionCandidateEvidenceDigest string
	GuardEvidenceDigest             string
	VerificationEvidenceDigest      string
	ReverseEvidenceDigest           string
	EvidencePrefixDigest            string
	MetricEvidenceDigest            string
	FeedbackDigest                  string
	AggregationEvidenceDigest       string
	DirectionEvidenceDigest         string
	CandidateStatus                 string
	CandidateDigest                 string
	CandidateEvidenceDigest         string
	RevisionSource                  string
	BoundRevisionChangeDigest       string
	EvidenceDigest                  string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
		b.CandidateID == "" ||
		b.SourceCandidateID == "" ||
		b.SourceRevisionCandidateDigest == "" ||
		b.SourceCandidateEvidenceDigest == "" ||
		b.RevisionCandidateDigest == "" ||
		b.RevisionCandidateEvidenceDigest == "" ||
		b.GuardEvidenceDigest == "" ||
		b.VerificationEvidenceDigest == "" ||
		b.ReverseEvidenceDigest == "" ||
		b.EvidencePrefixDigest == "" ||
		b.MetricEvidenceDigest == "" ||
		b.FeedbackDigest == "" ||
		b.AggregationEvidenceDigest == "" ||
		b.DirectionEvidenceDigest == "" ||
		b.CandidateStatus != jevImprovementRevisionCandidateReady ||
		b.CandidateDigest == "" ||
		b.CandidateEvidenceDigest == "" ||
		b.RevisionSource == "" ||
		b.BoundRevisionChangeDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended lineage candidate revision generation binding")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage candidate revision generation must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(
		b.CandidateID,
		b.SourceCandidateID,
		b.SourceRevisionCandidateDigest,
		b.SourceCandidateEvidenceDigest,
		b.RevisionCandidateDigest,
		b.RevisionCandidateEvidenceDigest,
		b.GuardEvidenceDigest,
		b.VerificationEvidenceDigest,
		b.ReverseEvidenceDigest,
		b.EvidencePrefixDigest,
		b.MetricEvidenceDigest,
		b.FeedbackDigest,
		b.AggregationEvidenceDigest,
		b.DirectionEvidenceDigest,
		b.CandidateStatus,
		b.CandidateDigest,
		b.CandidateEvidenceDigest,
		b.RevisionSource,
		b.BoundRevisionChangeDigest,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended lineage candidate revision generation digest mismatch")
	}
	return nil
}

// GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage
// generates only from an extended direction and never applies the candidate.
func GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-generation-extended-lineage"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageBinding{
			Status:         "UNKNOWN",
			MissingStage:   stage,
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.Feedback.NonAuthorizing || !input.Direction.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Feedback.NonExecuting || !input.Direction.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Feedback.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-feedback-validation")
	}
	if err := input.Direction.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-direction-validation")
	}
	if input.Direction.CandidateDigest != input.Feedback.CandidateDigest {
		return unknown("direction-feedback-candidate")
	}
	if input.Direction.InputEvidenceDigest != input.Feedback.AggregationEvidenceDigest {
		return unknown("direction-feedback-evidence")
	}
	if strings.TrimSpace(input.RevisionSource) == "" {
		return unknown("revision-source")
	}
	if strings.TrimSpace(input.RevisionChangeDigest) == "" {
		return unknown("revision-change")
	}
	boundChangeDigest, err := Digest(struct {
		RevisionChangeDigest      string
		VerificationEvidenceDigest string
		EvidencePrefixDigest      string
		MetricEvidenceDigest      string
		FeedbackDigest            string
		AggregationEvidenceDigest string
		DirectionEvidenceDigest  string
	}{
		RevisionChangeDigest:       input.RevisionChangeDigest,
		VerificationEvidenceDigest: input.Feedback.VerificationEvidenceDigest,
		EvidencePrefixDigest:       input.Feedback.EvidencePrefixDigest,
		MetricEvidenceDigest:       input.Feedback.MetricEvidenceDigest,
		FeedbackDigest:             input.Feedback.FeedbackDigest,
		AggregationEvidenceDigest: input.Feedback.AggregationEvidenceDigest,
		DirectionEvidenceDigest:    input.Direction.DirectionEvidenceDigest,
	})
	if err != nil {
		return unknown("revision-change-evidence")
	}
	direction := JEVImprovementDirectionDirective{
		Status:              input.Direction.DirectiveStatus,
		Directive:           input.Direction.Directive,
		CandidateDigest:     input.Direction.CandidateDigest,
		CandidateSource:     input.Direction.CandidateSource,
		InputEvidenceDigest: input.Direction.InputEvidenceDigest,
		EvidenceDigest:      input.Direction.DirectionEvidenceDigest,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	candidate := GenerateJEVImprovementRevisionCandidate(JEVImprovementRevisionCandidateInput{
		Directive:            direction,
		RevisionSource:       input.RevisionSource,
		RevisionChangeDigest: boundChangeDigest,
		NonAuthorizing:       true,
	})
	if err := candidate.Validate(); err != nil {
		stage := candidate.MissingStage
		if strings.TrimSpace(stage) == "" {
			stage = "revision-candidate"
		}
		return unknown(stage)
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageBinding{
		Status:                          "bound",
		CandidateID:                     input.Feedback.CandidateID,
		SourceCandidateID:               input.Feedback.SourceCandidateID,
		SourceRevisionCandidateDigest:   input.Feedback.SourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   input.Feedback.SourceCandidateEvidenceDigest,
		RevisionCandidateDigest:         input.Feedback.RevisionCandidateDigest,
		RevisionCandidateEvidenceDigest: input.Feedback.RevisionCandidateEvidenceDigest,
		GuardEvidenceDigest:             input.Feedback.GuardEvidenceDigest,
		VerificationEvidenceDigest:      input.Feedback.VerificationEvidenceDigest,
		ReverseEvidenceDigest:           input.Feedback.ReverseEvidenceDigest,
		EvidencePrefixDigest:            input.Feedback.EvidencePrefixDigest,
		MetricEvidenceDigest:            input.Feedback.MetricEvidenceDigest,
		FeedbackDigest:                  input.Feedback.FeedbackDigest,
		AggregationEvidenceDigest:       input.Feedback.AggregationEvidenceDigest,
		DirectionEvidenceDigest:         input.Direction.DirectionEvidenceDigest,
		CandidateStatus:                 candidate.Status,
		CandidateDigest:                 candidate.CandidateDigest,
		CandidateEvidenceDigest:         candidate.EvidenceDigest,
		RevisionSource:                  candidate.RevisionSource,
		BoundRevisionChangeDigest:       boundChangeDigest,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(
		output.CandidateID,
		output.SourceCandidateID,
		output.SourceRevisionCandidateDigest,
		output.SourceCandidateEvidenceDigest,
		output.RevisionCandidateDigest,
		output.RevisionCandidateEvidenceDigest,
		output.GuardEvidenceDigest,
		output.VerificationEvidenceDigest,
		output.ReverseEvidenceDigest,
		output.EvidencePrefixDigest,
		output.MetricEvidenceDigest,
		output.FeedbackDigest,
		output.AggregationEvidenceDigest,
		output.DirectionEvidenceDigest,
		output.CandidateStatus,
		output.CandidateDigest,
		output.CandidateEvidenceDigest,
		output.RevisionSource,
		output.BoundRevisionChangeDigest,
	)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineage(
	candidateID,
	sourceCandidateID,
	sourceRevisionCandidateDigest,
	sourceCandidateEvidenceDigest,
	revisionCandidateDigest,
	revisionCandidateEvidenceDigest,
	guardEvidenceDigest,
	verificationEvidenceDigest,
	reverseEvidenceDigest,
	evidencePrefixDigest,
	metricEvidenceDigest,
	feedbackDigest,
	aggregationEvidenceDigest,
	directionEvidenceDigest,
	candidateStatus,
	candidateDigest,
	candidateEvidenceDigest,
	revisionSource,
	boundRevisionChangeDigest string,
) string {
	digest, err := Digest(struct {
		CandidateID                     string
		SourceCandidateID               string
		SourceRevisionCandidateDigest   string
		SourceCandidateEvidenceDigest   string
		RevisionCandidateDigest         string
		RevisionCandidateEvidenceDigest string
		GuardEvidenceDigest             string
		VerificationEvidenceDigest      string
		ReverseEvidenceDigest           string
		EvidencePrefixDigest            string
		MetricEvidenceDigest            string
		FeedbackDigest                  string
		AggregationEvidenceDigest       string
		DirectionEvidenceDigest         string
		CandidateStatus                 string
		CandidateDigest                 string
		CandidateEvidenceDigest         string
		RevisionSource                  string
		BoundRevisionChangeDigest       string
	}{
		CandidateID:                     candidateID,
		SourceCandidateID:               sourceCandidateID,
		SourceRevisionCandidateDigest:   sourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   sourceCandidateEvidenceDigest,
		RevisionCandidateDigest:         revisionCandidateDigest,
		RevisionCandidateEvidenceDigest: revisionCandidateEvidenceDigest,
		GuardEvidenceDigest:             guardEvidenceDigest,
		VerificationEvidenceDigest:      verificationEvidenceDigest,
		ReverseEvidenceDigest:           reverseEvidenceDigest,
		EvidencePrefixDigest:            evidencePrefixDigest,
		MetricEvidenceDigest:            metricEvidenceDigest,
		FeedbackDigest:                  feedbackDigest,
		AggregationEvidenceDigest:       aggregationEvidenceDigest,
		DirectionEvidenceDigest:         directionEvidenceDigest,
		CandidateStatus:                 candidateStatus,
		CandidateDigest:                 candidateDigest,
		CandidateEvidenceDigest:         candidateEvidenceDigest,
		RevisionSource:                  revisionSource,
		BoundRevisionChangeDigest:       boundRevisionChangeDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}
