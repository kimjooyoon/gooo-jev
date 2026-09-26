package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackInput
// adapts an evidence-linked action metric into the replay feedback ledger.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackInput struct {
	Metric                 ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricBinding
	CandidateDigest        string
	ReplayObservationDigest string
	NonAuthorizing         bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackBinding
// preserves action, metric, feedback, and aggregation evidence.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackBinding struct {
	Status                    string
	MissingStage              string
	MetricEvidenceDigest      string
	FeedbackStatus            string
	CandidateDigest           string
	ReplayObservationDigest   string
	FeedbackEvidenceDigest    string
	FeedbackDigest            string
	AggregationStatus         string
	Total                     int
	Confirmed                 int
	Refuted                   int
	Unknown                   int
	AggregationEvidenceDigest string
	EvidenceDigest            string
	NonExecuting              bool
	NonAuthorizing            bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
		b.MetricEvidenceDigest == "" ||
		b.FeedbackStatus == "" ||
		b.CandidateDigest == "" ||
		b.ReplayObservationDigest == "" ||
		b.FeedbackEvidenceDigest == "" ||
		b.FeedbackDigest == "" ||
		b.AggregationStatus == "" ||
		b.Total <= 0 ||
		b.Confirmed+b.Refuted+b.Unknown != b.Total ||
		b.AggregationEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo revision action feedback binding")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo revision action feedback binding must be non-executing and non-authorizing")
	}
	expected, err := Digest(struct {
		MetricEvidenceDigest      string
		FeedbackStatus            string
		CandidateDigest           string
		ReplayObservationDigest   string
		FeedbackEvidenceDigest    string
		FeedbackDigest            string
		AggregationStatus         string
		Total                     int
		Confirmed                 int
		Refuted                   int
		Unknown                   int
		AggregationEvidenceDigest string
	}{
		MetricEvidenceDigest:      b.MetricEvidenceDigest,
		FeedbackStatus:            b.FeedbackStatus,
		CandidateDigest:           b.CandidateDigest,
		ReplayObservationDigest:   b.ReplayObservationDigest,
		FeedbackEvidenceDigest:    b.FeedbackEvidenceDigest,
		FeedbackDigest:            b.FeedbackDigest,
		AggregationStatus:         b.AggregationStatus,
		Total:                     b.Total,
		Confirmed:                 b.Confirmed,
		Refuted:                   b.Refuted,
		Unknown:                   b.Unknown,
		AggregationEvidenceDigest: b.AggregationEvidenceDigest,
	})
	if err != nil || b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo revision action feedback digest mismatch")
	}
	return nil
}

// BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedback
// keeps review and unknown dispositions conservative before direction derivation.
func BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedback(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-feedback"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackBinding{
			Status:         "UNKNOWN",
			MissingStage:   stage,
			NonExecuting:   true,
			NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.Metric.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Metric.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Metric.Validate(); err != nil {
		return unknown("revision-action-metric-validation")
	}
	if strings.TrimSpace(input.CandidateDigest) == "" {
		return unknown("candidate-digest")
	}
	if strings.TrimSpace(input.ReplayObservationDigest) == "" {
		return unknown("replay-observation")
	}
	feedbackStatus := jevReplayFeedbackUnknown
	switch {
	case input.Metric.CounterexampleCount > 0:
		feedbackStatus = jevReplayFeedbackRefuted
	case input.Metric.UnknownCount > 0 || input.Metric.ReviewCount > 0:
		feedbackStatus = jevReplayFeedbackUnknown
	case input.Metric.VerifiedCount > 0:
		feedbackStatus = jevReplayFeedbackConfirmed
	default:
		return unknown("feedback-disposition")
	}
	feedback := JEVImprovementReplayFeedback{
		Status:                  feedbackStatus,
		FeedbackKind:            feedbackStatus,
		CandidateDigest:         input.CandidateDigest,
		ReplayObservationDigest: input.ReplayObservationDigest,
		MetricDigest:            input.Metric.MetricEvidenceDigest,
		NonExecuting:            true,
		NonAuthorizing:          true,
	}
	feedback.FeedbackEvidenceDigest = digestJEVImprovementReplayFeedbackEvidence(
		feedback.ReplayObservationDigest,
		feedback.MetricDigest,
		feedback.FeedbackKind,
	)
	feedback.EvidenceDigest = digestJEVImprovementReplayFeedback(
		feedback.Status,
		feedback.FeedbackKind,
		feedback.CandidateDigest,
		feedback.ReplayObservationDigest,
		feedback.MetricDigest,
		feedback.FeedbackEvidenceDigest,
	)
	if err := feedback.Validate(); err != nil {
		return unknown("feedback-evidence")
	}
	aggregation := AggregateJEVImprovementFeedback(JEVImprovementFeedbackAggregationInput{
		Feedback:       []JEVImprovementReplayFeedback{feedback},
		NonAuthorizing: true,
	})
	if err := aggregation.Validate(); err != nil {
		return unknown(aggregation.MissingStage)
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackBinding{
		Status:                    "bound",
		MetricEvidenceDigest:      input.Metric.MetricEvidenceDigest,
		FeedbackStatus:            feedbackStatus,
		CandidateDigest:           input.CandidateDigest,
		ReplayObservationDigest:   input.ReplayObservationDigest,
		FeedbackEvidenceDigest:    feedback.FeedbackEvidenceDigest,
		FeedbackDigest:            feedback.EvidenceDigest,
		AggregationStatus:         aggregation.Status,
		Total:                     aggregation.Total,
		Confirmed:                 aggregation.Confirmed,
		Refuted:                   aggregation.Refuted,
		Unknown:                   aggregation.Unknown,
		AggregationEvidenceDigest: aggregation.EvidenceDigest,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	output.EvidenceDigest, _ = Digest(struct {
		MetricEvidenceDigest      string
		FeedbackStatus            string
		CandidateDigest           string
		ReplayObservationDigest   string
		FeedbackEvidenceDigest    string
		FeedbackDigest            string
		AggregationStatus         string
		Total                     int
		Confirmed                 int
		Refuted                   int
		Unknown                   int
		AggregationEvidenceDigest string
	}{
		MetricEvidenceDigest:      output.MetricEvidenceDigest,
		FeedbackStatus:            output.FeedbackStatus,
		CandidateDigest:           output.CandidateDigest,
		ReplayObservationDigest:   output.ReplayObservationDigest,
		FeedbackEvidenceDigest:    output.FeedbackEvidenceDigest,
		FeedbackDigest:            output.FeedbackDigest,
		AggregationStatus:         output.AggregationStatus,
		Total:                     output.Total,
		Confirmed:                 output.Confirmed,
		Refuted:                   output.Refuted,
		Unknown:                   output.Unknown,
		AggregationEvidenceDigest: output.AggregationEvidenceDigest,
	})
	if err := output.Validate(); err != nil {
		return unknown("revision-action-feedback-evidence")
	}
	return output
}