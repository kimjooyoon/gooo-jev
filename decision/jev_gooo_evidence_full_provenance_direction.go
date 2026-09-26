package decision

import "strings"

// ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionInput adapts the
// evidence feedback aggregate into the existing direction decision model.
type ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionInput struct {
	Feedback       ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackBinding
	CandidateSource string
	NonAuthorizing bool
}

// DeriveExecutionEnvelopeGoooEvidenceFullProvenanceDirection converts the
// aggregate status back into review, revision, or evidence direction.
func DeriveExecutionEnvelopeGoooEvidenceFullProvenanceDirection(input ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionInput) JEVImprovementDirectionDirective {
	unknown := func(stage string) JEVImprovementDirectionDirective {
		if strings.TrimSpace(stage) == "" {
			stage = "evidence-direction"
		}
		return JEVImprovementDirectionDirective{
			Status: "UNKNOWN", MissingStage: stage,
			NonExecuting: true, NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.Feedback.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Feedback.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Feedback.Validate(); err != nil {
		return unknown("evaluation-feedback-validation")
	}
	if strings.TrimSpace(input.CandidateSource) == "" {
		return unknown("candidate-source")
	}
	aggregationInputEvidence := digestJEVImprovementFeedbackInputs(
		[]string{input.Feedback.FeedbackEvidenceDigest},
	)
	aggregation := JEVImprovementFeedbackAggregation{
		Status:              input.Feedback.AggregationStatus,
		Total:               input.Feedback.Total,
		Confirmed:           input.Feedback.Confirmed,
		Refuted:             input.Feedback.Refuted,
		Unknown:             input.Feedback.Unknown,
		InputEvidenceDigest: aggregationInputEvidence,
		EvidenceDigest:      input.Feedback.AggregationEvidenceDigest,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	if err := aggregation.Validate(); err != nil {
		return unknown("feedback-aggregation")
	}
	return DeriveJEVImprovementDirectionDirective(JEVImprovementDirectionDirectiveInput{
		Aggregation:     aggregation,
		CandidateDigest: input.Feedback.CandidateDigest,
		CandidateSource: input.CandidateSource,
		NonAuthorizing: true,
	})
}
