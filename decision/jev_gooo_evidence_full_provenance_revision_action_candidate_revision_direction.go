package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionInput
// adapts candidate revision feedback into the existing direction model.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionInput struct {
	Feedback       ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackBinding
	CandidateSource string
	NonAuthorizing bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionBinding
// preserves candidate revision lineage, aggregation, and direction evidence.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionBinding struct {
	Status                          string
	MissingStage                    string
	CandidateID                     string
	SourceCandidateID               string
	SourceRevisionCandidateDigest   string
	SourceCandidateEvidenceDigest   string
	RevisionCandidateDigest         string
	RevisionCandidateEvidenceDigest string
	MetricEvidenceDigest            string
	FeedbackDigest                  string
	AggregationEvidenceDigest       string
	DirectiveStatus                 string
	Directive                       string
	CandidateDigest                 string
	CandidateSource                 string
	InputEvidenceDigest             string
	DirectionEvidenceDigest         string
	EvidenceDigest                  string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
		b.CandidateID == "" ||
		b.SourceCandidateID == "" ||
		b.SourceRevisionCandidateDigest == "" ||
		b.SourceCandidateEvidenceDigest == "" ||
		b.RevisionCandidateDigest == "" ||
		b.RevisionCandidateEvidenceDigest == "" ||
		b.MetricEvidenceDigest == "" ||
		b.FeedbackDigest == "" ||
		b.AggregationEvidenceDigest == "" ||
		b.DirectiveStatus == "" ||
		b.Directive == "" ||
		b.CandidateDigest == "" ||
		b.CandidateSource == "" ||
		b.InputEvidenceDigest == "" ||
		b.DirectionEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo action-derived candidate revision direction binding")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo action-derived candidate revision direction must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirection(
		b.CandidateID,
		b.SourceCandidateID,
		b.SourceRevisionCandidateDigest,
		b.SourceCandidateEvidenceDigest,
		b.RevisionCandidateDigest,
		b.RevisionCandidateEvidenceDigest,
		b.MetricEvidenceDigest,
		b.FeedbackDigest,
		b.AggregationEvidenceDigest,
		b.DirectiveStatus,
		b.Directive,
		b.CandidateDigest,
		b.CandidateSource,
		b.InputEvidenceDigest,
		b.DirectionEvidenceDigest,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo action-derived candidate revision direction digest mismatch")
	}
	return nil
}

// DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirection
// derives a non-authorizing direction while preserving candidate revision lineage.
func DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirection(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-revision-direction"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionBinding{
			Status:         "UNKNOWN",
			MissingStage:   stage,
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.Feedback.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Feedback.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Feedback.Validate(); err != nil {
		return unknown("revision-action-candidate-revision-feedback-validation")
	}
	if strings.TrimSpace(input.CandidateSource) == "" {
		return unknown("candidate-source")
	}
	aggregation := JEVImprovementFeedbackAggregation{
		Status:              input.Feedback.AggregationStatus,
		Total:               input.Feedback.Total,
		Confirmed:           input.Feedback.Confirmed,
		Refuted:             input.Feedback.Refuted,
		Unknown:             input.Feedback.Unknown,
		InputEvidenceDigest: digestJEVImprovementFeedbackInputs([]string{input.Feedback.FeedbackDigest}),
		EvidenceDigest:       input.Feedback.AggregationEvidenceDigest,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	if err := aggregation.Validate(); err != nil {
		return unknown("feedback-aggregation")
	}
	direction := DeriveJEVImprovementDirectionDirective(JEVImprovementDirectionDirectiveInput{
		Aggregation:     aggregation,
		CandidateDigest: input.Feedback.CandidateDigest,
		CandidateSource: input.CandidateSource,
		NonAuthorizing: true,
	})
	if err := direction.Validate(); err != nil {
		stage := direction.MissingStage
		if strings.TrimSpace(stage) == "" {
			stage = "feedback-direction"
		}
		return unknown(stage)
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirectionBinding{
		Status:                          "bound",
		CandidateID:                     input.Feedback.CandidateID,
		SourceCandidateID:               input.Feedback.SourceCandidateID,
		SourceRevisionCandidateDigest:   input.Feedback.SourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   input.Feedback.SourceCandidateEvidenceDigest,
		RevisionCandidateDigest:         input.Feedback.RevisionCandidateDigest,
		RevisionCandidateEvidenceDigest: input.Feedback.RevisionCandidateEvidenceDigest,
		MetricEvidenceDigest:            input.Feedback.MetricEvidenceDigest,
		FeedbackDigest:                  input.Feedback.FeedbackDigest,
		AggregationEvidenceDigest:       input.Feedback.AggregationEvidenceDigest,
		DirectiveStatus:                 direction.Status,
		Directive:                       direction.Directive,
		CandidateDigest:                 direction.CandidateDigest,
		CandidateSource:                 direction.CandidateSource,
		InputEvidenceDigest:             direction.InputEvidenceDigest,
		DirectionEvidenceDigest:         direction.EvidenceDigest,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirection(
		output.CandidateID,
		output.SourceCandidateID,
		output.SourceRevisionCandidateDigest,
		output.SourceCandidateEvidenceDigest,
		output.RevisionCandidateDigest,
		output.RevisionCandidateEvidenceDigest,
		output.MetricEvidenceDigest,
		output.FeedbackDigest,
		output.AggregationEvidenceDigest,
		output.DirectiveStatus,
		output.Directive,
		output.CandidateDigest,
		output.CandidateSource,
		output.InputEvidenceDigest,
		output.DirectionEvidenceDigest,
	)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-revision-direction-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionDirection(
	candidateID,
	sourceCandidateID,
	sourceRevisionCandidateDigest,
	sourceCandidateEvidenceDigest,
	revisionCandidateDigest,
	revisionCandidateEvidenceDigest,
	metricEvidenceDigest,
	feedbackDigest,
	aggregationEvidenceDigest,
	directiveStatus,
	directive,
	candidateDigest,
	candidateSource,
	inputEvidenceDigest,
	directionEvidenceDigest string,
) string {
	digest, err := Digest(struct {
		CandidateID                     string
		SourceCandidateID               string
		SourceRevisionCandidateDigest   string
		SourceCandidateEvidenceDigest   string
		RevisionCandidateDigest         string
		RevisionCandidateEvidenceDigest string
		MetricEvidenceDigest            string
		FeedbackDigest                  string
		AggregationEvidenceDigest       string
		DirectiveStatus                 string
		Directive                       string
		CandidateDigest                 string
		CandidateSource                 string
		InputEvidenceDigest             string
		DirectionEvidenceDigest         string
	}{
		CandidateID:                     candidateID,
		SourceCandidateID:               sourceCandidateID,
		SourceRevisionCandidateDigest:   sourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   sourceCandidateEvidenceDigest,
		RevisionCandidateDigest:         revisionCandidateDigest,
		RevisionCandidateEvidenceDigest: revisionCandidateEvidenceDigest,
		MetricEvidenceDigest:            metricEvidenceDigest,
		FeedbackDigest:                  feedbackDigest,
		AggregationEvidenceDigest:       aggregationEvidenceDigest,
		DirectiveStatus:                 directiveStatus,
		Directive:                       directive,
		CandidateDigest:                 candidateDigest,
		CandidateSource:                 candidateSource,
		InputEvidenceDigest:             inputEvidenceDigest,
		DirectionEvidenceDigest:         directionEvidenceDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}