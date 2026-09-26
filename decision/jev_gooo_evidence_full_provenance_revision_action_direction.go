package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionInput
// adapts action feedback aggregation into the existing direction model.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionInput struct {
	Feedback       ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackBinding
	CandidateSource string
	NonAuthorizing bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionBinding
// preserves action metric, feedback, aggregation, and direction evidence.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionBinding struct {
	Status                    string
	MissingStage              string
	MetricEvidenceDigest      string
	FeedbackDigest            string
	AggregationEvidenceDigest string
	DirectiveStatus           string
	Directive                 string
	CandidateDigest           string
	CandidateSource           string
	InputEvidenceDigest       string
	DirectionEvidenceDigest   string
	EvidenceDigest            string
	NonExecuting              bool
	NonAuthorizing            bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
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
		return fmt.Errorf("incomplete Gooo revision action direction binding")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo revision action direction binding must be non-executing and non-authorizing")
	}
	expected, err := Digest(struct {
		MetricEvidenceDigest      string
		FeedbackDigest            string
		AggregationEvidenceDigest string
		DirectiveStatus           string
		Directive                 string
		CandidateDigest           string
		CandidateSource           string
		InputEvidenceDigest       string
		DirectionEvidenceDigest   string
	}{
		MetricEvidenceDigest:      b.MetricEvidenceDigest,
		FeedbackDigest:            b.FeedbackDigest,
		AggregationEvidenceDigest: b.AggregationEvidenceDigest,
		DirectiveStatus:            b.DirectiveStatus,
		Directive:                 b.Directive,
		CandidateDigest:            b.CandidateDigest,
		CandidateSource:            b.CandidateSource,
		InputEvidenceDigest:        b.InputEvidenceDigest,
		DirectionEvidenceDigest:    b.DirectionEvidenceDigest,
	})
	if err != nil || b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo revision action direction digest mismatch")
	}
	return nil
}

// DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirection
// never authorizes the chosen direction and preserves missing evidence.
func DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirection(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-direction"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionBinding{
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
		return unknown("revision-action-feedback-validation")
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
		InputEvidenceDigest: digestJEVImprovementFeedbackInputs([]string{input.Feedback.FeedbackEvidenceDigest}),
		EvidenceDigest:      input.Feedback.AggregationEvidenceDigest,
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
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionBinding{
		Status:                    "bound",
		MetricEvidenceDigest:      input.Feedback.MetricEvidenceDigest,
		FeedbackDigest:            input.Feedback.FeedbackDigest,
		AggregationEvidenceDigest: input.Feedback.AggregationEvidenceDigest,
		DirectiveStatus:           direction.Status,
		Directive:                 direction.Directive,
		CandidateDigest:           direction.CandidateDigest,
		CandidateSource:           direction.CandidateSource,
		InputEvidenceDigest:       direction.InputEvidenceDigest,
		DirectionEvidenceDigest:   direction.EvidenceDigest,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	output.EvidenceDigest, _ = Digest(struct {
		MetricEvidenceDigest      string
		FeedbackDigest            string
		AggregationEvidenceDigest string
		DirectiveStatus           string
		Directive                 string
		CandidateDigest           string
		CandidateSource           string
		InputEvidenceDigest       string
		DirectionEvidenceDigest   string
	}{
		MetricEvidenceDigest:      output.MetricEvidenceDigest,
		FeedbackDigest:            output.FeedbackDigest,
		AggregationEvidenceDigest: output.AggregationEvidenceDigest,
		DirectiveStatus:            output.DirectiveStatus,
		Directive:                 output.Directive,
		CandidateDigest:           output.CandidateDigest,
		CandidateSource:           output.CandidateSource,
		InputEvidenceDigest:       output.InputEvidenceDigest,
		DirectionEvidenceDigest:   output.DirectionEvidenceDigest,
	})
	if err := output.Validate(); err != nil {
		return unknown("revision-action-direction-evidence")
	}
	return output
}