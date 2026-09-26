package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateInput
// connects action feedback direction to the existing revision generator.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateInput struct {
	Feedback             ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackBinding
	Direction            ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionDirectionBinding
	RevisionSource       string
	RevisionChangeDigest string
	NonAuthorizing       bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateBinding
// preserves action, direction, and generated candidate provenance.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateBinding struct {
	Status                    string
	MissingStage              string
	MetricEvidenceDigest      string
	FeedbackDigest            string
	AggregationEvidenceDigest string
	DirectionEvidenceDigest   string
	CandidateStatus           string
	CandidateDigest           string
	CandidateEvidenceDigest   string
	RevisionSource            string
	BoundRevisionChangeDigest string
	EvidenceDigest            string
	NonExecuting              bool
	NonAuthorizing            bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
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
		return fmt.Errorf("incomplete Gooo revision action candidate binding")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo revision action candidate binding must be non-executing and non-authorizing")
	}
	expected, err := Digest(struct {
		MetricEvidenceDigest      string
		FeedbackDigest            string
		AggregationEvidenceDigest string
		DirectionEvidenceDigest   string
		CandidateStatus            string
		CandidateDigest            string
		CandidateEvidenceDigest    string
		RevisionSource             string
		BoundRevisionChangeDigest  string
	}{
		MetricEvidenceDigest:      b.MetricEvidenceDigest,
		FeedbackDigest:            b.FeedbackDigest,
		AggregationEvidenceDigest: b.AggregationEvidenceDigest,
		DirectionEvidenceDigest:   b.DirectionEvidenceDigest,
		CandidateStatus:            b.CandidateStatus,
		CandidateDigest:            b.CandidateDigest,
		CandidateEvidenceDigest:    b.CandidateEvidenceDigest,
		RevisionSource:             b.RevisionSource,
		BoundRevisionChangeDigest:  b.BoundRevisionChangeDigest,
	})
	if err != nil || b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo revision action candidate digest mismatch")
	}
	return nil
}

// GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate
// generates only from a revision directive and never applies the candidate.
func GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidate(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateBinding{
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
		return unknown("revision-action-feedback-validation")
	}
	if err := input.Direction.Validate(); err != nil {
		return unknown("revision-action-direction-validation")
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
		RevisionChangeDigest       string
		MetricEvidenceDigest       string
		FeedbackDigest             string
		AggregationEvidenceDigest  string
		DirectionEvidenceDigest    string
	}{
		RevisionChangeDigest:      input.RevisionChangeDigest,
		MetricEvidenceDigest:      input.Feedback.MetricEvidenceDigest,
		FeedbackDigest:            input.Feedback.FeedbackDigest,
		AggregationEvidenceDigest: input.Feedback.AggregationEvidenceDigest,
		DirectionEvidenceDigest:   input.Direction.DirectionEvidenceDigest,
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
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateBinding{
		Status:                    "bound",
		MetricEvidenceDigest:      input.Feedback.MetricEvidenceDigest,
		FeedbackDigest:            input.Feedback.FeedbackDigest,
		AggregationEvidenceDigest: input.Feedback.AggregationEvidenceDigest,
		DirectionEvidenceDigest:   input.Direction.DirectionEvidenceDigest,
		CandidateStatus:            candidate.Status,
		CandidateDigest:            candidate.CandidateDigest,
		CandidateEvidenceDigest:    candidate.EvidenceDigest,
		RevisionSource:            candidate.RevisionSource,
		BoundRevisionChangeDigest: boundChangeDigest,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	output.EvidenceDigest, _ = Digest(struct {
		MetricEvidenceDigest      string
		FeedbackDigest            string
		AggregationEvidenceDigest string
		DirectionEvidenceDigest   string
		CandidateStatus            string
		CandidateDigest            string
		CandidateEvidenceDigest   string
		RevisionSource             string
		BoundRevisionChangeDigest  string
	}{
		MetricEvidenceDigest:      output.MetricEvidenceDigest,
		FeedbackDigest:            output.FeedbackDigest,
		AggregationEvidenceDigest: output.AggregationEvidenceDigest,
		DirectionEvidenceDigest:   output.DirectionEvidenceDigest,
		CandidateStatus:            output.CandidateStatus,
		CandidateDigest:            output.CandidateDigest,
		CandidateEvidenceDigest:    output.CandidateEvidenceDigest,
		RevisionSource:             output.RevisionSource,
		BoundRevisionChangeDigest:  output.BoundRevisionChangeDigest,
	})
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-evidence")
	}
	return output
}