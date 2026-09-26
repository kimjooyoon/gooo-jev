package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextInput
// converts observed-plan feedback into the next self-improvement action.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextInput struct {
	Feedback            ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackBinding
	CandidateSource     string
	RevisionSource      string
	RevisionChangeDigest string
	NonAuthorizing      bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextBinding
// preserves feedback, aggregation, direction, and optional next-candidate provenance.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextBinding struct {
	Status                       string
	MissingStage                 string
	NextAction                   string
	ParentCandidateDigest        string
	CandidateSource              string
	FeedbackStatus               string
	FeedbackEvidenceDigest       string
	FeedbackBindingEvidenceDigest string
	MetricEvidenceDigest         string
	ReverseObservationDigest     string
	EvidencePrefixDigest         string
	AggregationStatus            string
	AggregationEvidenceDigest    string
	DirectionStatus              string
	Directive                    string
	DirectionEvidenceDigest      string
	RevisionSource               string
	BoundRevisionChangeDigest    string
	RevisionCandidateStatus      string
	RevisionCandidateDigest      string
	RevisionCandidateEvidenceDigest string
	EvidenceDigest               string
	NonExecuting                 bool
	NonAuthorizing               bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
		b.NextAction == "" ||
		b.ParentCandidateDigest == "" ||
		b.CandidateSource == "" ||
		b.FeedbackStatus == "" ||
		b.FeedbackEvidenceDigest == "" ||
		b.FeedbackBindingEvidenceDigest == "" ||
		b.MetricEvidenceDigest == "" ||
		b.ReverseObservationDigest == "" ||
		b.EvidencePrefixDigest == "" ||
		b.AggregationStatus == "" ||
		b.AggregationEvidenceDigest == "" ||
		b.DirectionStatus == "" ||
		b.Directive == "" ||
		b.DirectionEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended lineage feedback next binding")
	}
	expectedAction, expectedAggregation := nextActionForGoooExtendedLineageFeedback(b.FeedbackStatus)
	if expectedAction == "" ||
		b.NextAction != expectedAction ||
		b.AggregationStatus != expectedAggregation {
		return fmt.Errorf("invalid Gooo extended lineage feedback next action")
	}
	expectedDirection := directionForGoooExtendedLineageFeedback(b.FeedbackStatus)
	if b.Directive != expectedDirection || b.DirectionStatus != expectedDirection {
		return fmt.Errorf("invalid Gooo extended lineage feedback direction")
	}
	if b.NextAction == "generate-revision-candidate" {
		if b.RevisionSource == "" ||
			b.BoundRevisionChangeDigest == "" ||
			b.RevisionCandidateStatus != jevImprovementRevisionCandidateReady ||
			b.RevisionCandidateDigest == "" ||
			b.RevisionCandidateEvidenceDigest == "" {
			return fmt.Errorf("incomplete Gooo extended lineage next revision candidate")
		}
		candidate := JEVImprovementRevisionCandidate{
			Status:                  b.RevisionCandidateStatus,
			ParentCandidateDigest:   b.ParentCandidateDigest,
			RevisionSource:          b.RevisionSource,
			RevisionChangeDigest:    b.BoundRevisionChangeDigest,
			DirectiveEvidenceDigest: b.DirectionEvidenceDigest,
			CandidateDigest:         b.RevisionCandidateDigest,
			EvidenceDigest:          b.RevisionCandidateEvidenceDigest,
			NonExecuting:            true,
			NonAuthorizing:          true,
		}
		if err := candidate.Validate(); err != nil {
			return fmt.Errorf("invalid Gooo extended lineage next revision candidate: %w", err)
		}
	} else if b.RevisionCandidateStatus != "" ||
		b.RevisionCandidateDigest != "" ||
		b.RevisionCandidateEvidenceDigest != "" ||
		b.RevisionSource != "" ||
		b.BoundRevisionChangeDigest != "" {
		return fmt.Errorf("unexpected revision candidate for non-revision feedback")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage feedback next binding must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNext(b)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended lineage feedback next digest mismatch")
	}
	return nil
}

// DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNext
// turns explicit feedback into review, hold, or revision generation without executing changes.
func DeriveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNext(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-generation-extended-lineage-candidate-execution-plan-feedback-next"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextBinding{
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
		return unknown("execution-plan-observation-feedback-validation")
	}
	if strings.TrimSpace(input.CandidateSource) == "" {
		return unknown("candidate-source")
	}
	nextAction, aggregationStatus := nextActionForGoooExtendedLineageFeedback(input.Feedback.FeedbackStatus)
	if nextAction == "" {
		return unknown("feedback-direction")
	}
	inputEvidenceDigest := digestJEVImprovementFeedbackInputs([]string{
		input.Feedback.EvidenceDigest,
		input.Feedback.MetricEvidenceDigest,
		input.Feedback.ReverseObservationDigest,
		input.Feedback.AggregationEvidenceDigest,
	})
	aggregation := JEVImprovementFeedbackAggregation{
		Status:              aggregationStatus,
		Total:               1,
		InputEvidenceDigest: inputEvidenceDigest,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	switch input.Feedback.FeedbackStatus {
	case "confirmed":
		aggregation.Confirmed = 1
	case "refuted":
		aggregation.Refuted = 1
	case "unknown":
		aggregation.Unknown = 1
	default:
		return unknown("feedback-direction")
	}
	aggregation.EvidenceDigest = digestJEVImprovementFeedbackAggregation(
		aggregation.Status,
		aggregation.Total,
		aggregation.Confirmed,
		aggregation.Refuted,
		aggregation.Unknown,
		aggregation.InputEvidenceDigest,
	)
	if err := aggregation.Validate(); err != nil {
		return unknown("feedback-aggregation")
	}
	direction := DeriveJEVImprovementDirectionDirective(JEVImprovementDirectionDirectiveInput{
		Aggregation:     aggregation,
		CandidateDigest: input.Feedback.CandidateDigest,
		CandidateSource: input.CandidateSource,
		NonAuthorizing:  true,
	})
	if err := direction.Validate(); err != nil {
		return unknown("feedback-direction")
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextBinding{
		Status:                       "bound",
		NextAction:                   nextAction,
		ParentCandidateDigest:        input.Feedback.CandidateDigest,
		CandidateSource:              input.CandidateSource,
		FeedbackStatus:               input.Feedback.FeedbackStatus,
		FeedbackEvidenceDigest:       input.Feedback.FeedbackEvidenceDigest,
		FeedbackBindingEvidenceDigest: input.Feedback.EvidenceDigest,
		MetricEvidenceDigest:         input.Feedback.MetricEvidenceDigest,
		ReverseObservationDigest:     input.Feedback.ReverseObservationDigest,
		EvidencePrefixDigest:         input.Feedback.EvidencePrefixDigest,
		AggregationStatus:            aggregation.Status,
		AggregationEvidenceDigest:    aggregation.EvidenceDigest,
		DirectionStatus:              direction.Status,
		Directive:                    direction.Directive,
		DirectionEvidenceDigest:      direction.EvidenceDigest,
		NonExecuting:                 true,
		NonAuthorizing:               true,
	}
	if nextAction == "generate-revision-candidate" {
		if strings.TrimSpace(input.RevisionSource) == "" {
			return unknown("revision-source")
		}
		if strings.TrimSpace(input.RevisionChangeDigest) == "" {
			return unknown("revision-change")
		}
		boundChangeDigest := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextChange(
			input.RevisionChangeDigest,
			input.Feedback.CandidateDigest,
			input.Feedback.EvidenceDigest,
			input.Feedback.MetricEvidenceDigest,
			input.Feedback.ReverseObservationDigest,
			input.Feedback.EvidencePrefixDigest,
			aggregation.EvidenceDigest,
			direction.EvidenceDigest,
		)
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
		output.RevisionSource = candidate.RevisionSource
		output.BoundRevisionChangeDigest = candidate.RevisionChangeDigest
		output.RevisionCandidateStatus = candidate.Status
		output.RevisionCandidateDigest = candidate.CandidateDigest
		output.RevisionCandidateEvidenceDigest = candidate.EvidenceDigest
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNext(output)
	if err := output.Validate(); err != nil {
		return unknown("feedback-next-evidence")
	}
	return output
}

func nextActionForGoooExtendedLineageFeedback(feedbackStatus string) (string, string) {
	switch strings.TrimSpace(feedbackStatus) {
	case "confirmed":
		return "external-review-required", jevImprovementFeedbackStableForReview
	case "refuted":
		return "generate-revision-candidate", jevImprovementFeedbackNeedsRevision
	case "unknown":
		return "evidence-hold", jevImprovementFeedbackHold
	default:
		return "", ""
	}
}

func directionForGoooExtendedLineageFeedback(feedbackStatus string) string {
	switch strings.TrimSpace(feedbackStatus) {
	case "confirmed":
		return jevImprovementDirectiveExternalReview
	case "refuted":
		return jevImprovementDirectiveRevision
	case "unknown":
		return jevImprovementDirectiveEvidence
	default:
		return ""
	}
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextChange(
	revisionChangeDigest,
	parentCandidateDigest,
	feedbackBindingEvidenceDigest,
	metricEvidenceDigest,
	reverseObservationDigest,
	evidencePrefixDigest,
	aggregationEvidenceDigest,
	directionEvidenceDigest string,
) string {
	digest, err := Digest(struct {
		RevisionChangeDigest      string
		ParentCandidateDigest     string
		FeedbackBindingEvidenceDigest string
		MetricEvidenceDigest      string
		ReverseObservationDigest  string
		EvidencePrefixDigest      string
		AggregationEvidenceDigest string
		DirectionEvidenceDigest   string
	}{
		RevisionChangeDigest:         revisionChangeDigest,
		ParentCandidateDigest:        parentCandidateDigest,
		FeedbackBindingEvidenceDigest: feedbackBindingEvidenceDigest,
		MetricEvidenceDigest:         metricEvidenceDigest,
		ReverseObservationDigest:     reverseObservationDigest,
		EvidencePrefixDigest:         evidencePrefixDigest,
		AggregationEvidenceDigest:   aggregationEvidenceDigest,
		DirectionEvidenceDigest:      directionEvidenceDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNext(b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanFeedbackNextBinding) string {
	digest, err := Digest(struct {
		Status                       string
		NextAction                   string
		ParentCandidateDigest        string
		CandidateSource              string
		FeedbackStatus               string
		FeedbackEvidenceDigest       string
		FeedbackBindingEvidenceDigest string
		MetricEvidenceDigest         string
		ReverseObservationDigest     string
		EvidencePrefixDigest         string
		AggregationStatus            string
		AggregationEvidenceDigest    string
		DirectionStatus              string
		Directive                    string
		DirectionEvidenceDigest      string
		RevisionSource               string
		BoundRevisionChangeDigest    string
		RevisionCandidateStatus      string
		RevisionCandidateDigest      string
		RevisionCandidateEvidenceDigest string
		NonExecuting                 bool
		NonAuthorizing               bool
	}{
		Status:                       b.Status,
		NextAction:                   b.NextAction,
		ParentCandidateDigest:        b.ParentCandidateDigest,
		CandidateSource:              b.CandidateSource,
		FeedbackStatus:               b.FeedbackStatus,
		FeedbackEvidenceDigest:       b.FeedbackEvidenceDigest,
		FeedbackBindingEvidenceDigest: b.FeedbackBindingEvidenceDigest,
		MetricEvidenceDigest:         b.MetricEvidenceDigest,
		ReverseObservationDigest:     b.ReverseObservationDigest,
		EvidencePrefixDigest:         b.EvidencePrefixDigest,
		AggregationStatus:            b.AggregationStatus,
		AggregationEvidenceDigest:    b.AggregationEvidenceDigest,
		DirectionStatus:              b.DirectionStatus,
		Directive:                    b.Directive,
		DirectionEvidenceDigest:      b.DirectionEvidenceDigest,
		RevisionSource:               b.RevisionSource,
		BoundRevisionChangeDigest:    b.BoundRevisionChangeDigest,
		RevisionCandidateStatus:      b.RevisionCandidateStatus,
		RevisionCandidateDigest:      b.RevisionCandidateDigest,
		RevisionCandidateEvidenceDigest: b.RevisionCandidateEvidenceDigest,
		NonExecuting:                 b.NonExecuting,
		NonAuthorizing:               b.NonAuthorizing,
	})
	if err != nil {
		return ""
	}
	return digest
}