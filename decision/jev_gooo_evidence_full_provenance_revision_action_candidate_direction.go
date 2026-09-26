package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionInput
// adapts candidate feedback into the existing direction model.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionInput struct {
	Feedback       ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackBinding
	CandidateSource string
	NonAuthorizing bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionBinding
// preserves candidate lineage, feedback aggregation, and direction evidence.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionBinding struct {
	Status                     string
	MissingStage               string
	CandidateID                string
	RevisionCandidateDigest    string
	CandidateEvidenceDigest    string
	MetricEvidenceDigest       string
	FeedbackDigest             string
	AggregationEvidenceDigest  string
	DirectiveStatus            string
	Directive                  string
	CandidateDigest            string
	CandidateSource            string
	InputEvidenceDigest        string
	DirectionEvidenceDigest    string
	EvidenceDigest             string
	NonExecuting               bool
	NonAuthorizing             bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
		b.CandidateID == "" ||
		b.RevisionCandidateDigest == "" ||
		b.CandidateEvidenceDigest == "" ||
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
		return fmt.Errorf("incomplete Gooo action-derived candidate direction binding")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo action-derived candidate direction must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateDirection(
		b.CandidateID,
		b.RevisionCandidateDigest,
		b.CandidateEvidenceDigest,
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
		return fmt.Errorf("Gooo action-derived candidate direction digest mismatch")
	}
	return nil
}

// DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirection
// derives a non-authorizing direction while preserving candidate lineage.
func DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirection(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-direction"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionBinding{
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
		return unknown("revision-action-candidate-feedback-validation")
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
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateDirectionBinding{
		Status:                    "bound",
		CandidateID:               input.Feedback.CandidateID,
		RevisionCandidateDigest:   input.Feedback.RevisionCandidateDigest,
		CandidateEvidenceDigest:   input.Feedback.CandidateEvidenceDigest,
		MetricEvidenceDigest:      input.Feedback.MetricEvidenceDigest,
		FeedbackDigest:            input.Feedback.FeedbackDigest,
		AggregationEvidenceDigest: input.Feedback.AggregationEvidenceDigest,
		DirectiveStatus:            direction.Status,
		Directive:                 direction.Directive,
		CandidateDigest:            direction.CandidateDigest,
		CandidateSource:            direction.CandidateSource,
		InputEvidenceDigest:        direction.InputEvidenceDigest,
		DirectionEvidenceDigest:   direction.EvidenceDigest,
		NonExecuting:               true,
		NonAuthorizing:             true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateDirection(
		output.CandidateID,
		output.RevisionCandidateDigest,
		output.CandidateEvidenceDigest,
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
		return unknown("revision-action-candidate-direction-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateDirection(
	candidateID,
	revisionCandidateDigest,
	candidateEvidenceDigest,
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
		CandidateID               string
		RevisionCandidateDigest  string
		CandidateEvidenceDigest  string
		MetricEvidenceDigest     string
		FeedbackDigest            string
		AggregationEvidenceDigest string
		DirectiveStatus           string
		Directive                 string
		CandidateDigest           string
		CandidateSource           string
		InputEvidenceDigest       string
		DirectionEvidenceDigest   string
	}{
		CandidateID:               candidateID,
		RevisionCandidateDigest:  revisionCandidateDigest,
		CandidateEvidenceDigest:   candidateEvidenceDigest,
		MetricEvidenceDigest:      metricEvidenceDigest,
		FeedbackDigest:            feedbackDigest,
		AggregationEvidenceDigest: aggregationEvidenceDigest,
		DirectiveStatus:            directiveStatus,
		Directive:                 directive,
		CandidateDigest:            candidateDigest,
		CandidateSource:            candidateSource,
		InputEvidenceDigest:        inputEvidenceDigest,
		DirectionEvidenceDigest:    directionEvidenceDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}