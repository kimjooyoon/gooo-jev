package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackInput feeds
// materialization evaluation back into the existing replay feedback aggregate.
type ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackInput struct {
	Evaluation             ExecutionEnvelopeGoooEvidenceFullProvenanceMaterializationAdmissionEvaluationBinding
	CandidateDigest        string
	ReplayObservationDigest string
	NonAuthorizing         bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackBinding
// preserves the evaluation result and the newly derived feedback aggregate.
type ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackBinding struct {
	Status                  string
	MissingStage            string
	FeedbackStatus          string
	CandidateDigest         string
	ReplayObservationDigest string
	MetricDigest            string
	FeedbackEvidenceDigest  string
	AggregationStatus       string
	Total                   int
	Confirmed               int
	Refuted                 int
	Unknown                 int
	AggregationEvidenceDigest string
	EvidenceDigest          string
	NonExecuting            bool
	NonAuthorizing          bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
		b.FeedbackStatus == "" ||
		b.CandidateDigest == "" ||
		b.ReplayObservationDigest == "" ||
		b.MetricDigest == "" ||
		b.FeedbackEvidenceDigest == "" ||
		b.AggregationStatus == "" ||
		b.Total <= 0 ||
		b.Confirmed+b.Refuted+b.Unknown != b.Total ||
		b.AggregationEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo evidence evaluation feedback binding")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo evidence evaluation feedback binding must be non-executing and non-authorizing")
	}
	expected, err := Digest(struct {
		FeedbackStatus            string
		CandidateDigest           string
		ReplayObservationDigest   string
		MetricDigest              string
		FeedbackEvidenceDigest    string
		AggregationStatus         string
		Total                     int
		Confirmed                 int
		Refuted                   int
		Unknown                   int
		AggregationEvidenceDigest string
	}{
		FeedbackStatus:            b.FeedbackStatus,
		CandidateDigest:           b.CandidateDigest,
		ReplayObservationDigest:   b.ReplayObservationDigest,
		MetricDigest:              b.MetricDigest,
		FeedbackEvidenceDigest:    b.FeedbackEvidenceDigest,
		AggregationStatus:          b.AggregationStatus,
		Total:                      b.Total,
		Confirmed:                  b.Confirmed,
		Refuted:                    b.Refuted,
		Unknown:                    b.Unknown,
		AggregationEvidenceDigest: b.AggregationEvidenceDigest,
	})
	if err != nil || b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo evidence evaluation feedback digest mismatch")
	}
	return nil
}

// BindExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedback converts
// admission/evaluation outcomes into the existing feedback aggregation model.
func BindExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedback(input ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackInput) ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "evaluation-feedback"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackBinding{
			Status: "UNKNOWN", MissingStage: stage,
			NonExecuting: true, NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.Evaluation.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.Evaluation.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.Evaluation.Validate(); err != nil {
		return unknown("evaluation-validation")
	}
	if input.Evaluation.Status != "bound" {
		return unknown("evaluation")
	}
	if strings.TrimSpace(input.CandidateDigest) == "" {
		return unknown("candidate-digest")
	}
	if strings.TrimSpace(input.ReplayObservationDigest) == "" {
		return unknown("replay-observation")
	}
	feedbackStatus := jevReplayFeedbackUnknown
	switch input.Evaluation.AdmissionDecision {
	case "admit":
		feedbackStatus = jevReplayFeedbackConfirmed
	case "reject":
		feedbackStatus = jevReplayFeedbackRefuted
	case "hold":
		feedbackStatus = jevReplayFeedbackUnknown
	default:
		return unknown("admission-decision")
	}
	feedback := JEVImprovementReplayFeedback{
		Status:                  feedbackStatus,
		FeedbackKind:            feedbackStatus,
		CandidateDigest:         input.CandidateDigest,
		ReplayObservationDigest: input.ReplayObservationDigest,
		MetricDigest:            input.Evaluation.EvaluationBridgeDigest,
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
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceEvaluationFeedbackBinding{
		Status:                    "bound",
		FeedbackStatus:             feedbackStatus,
		CandidateDigest:            input.CandidateDigest,
		ReplayObservationDigest:   input.ReplayObservationDigest,
		MetricDigest:               feedback.MetricDigest,
		FeedbackEvidenceDigest:    feedback.EvidenceDigest,
		AggregationStatus:          aggregation.Status,
		Total:                     aggregation.Total,
		Confirmed:                 aggregation.Confirmed,
		Refuted:                   aggregation.Refuted,
		Unknown:                   aggregation.Unknown,
		AggregationEvidenceDigest: aggregation.EvidenceDigest,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	output.EvidenceDigest, _ = Digest(struct {
		FeedbackStatus            string
		CandidateDigest           string
		ReplayObservationDigest   string
		MetricDigest              string
		FeedbackEvidenceDigest    string
		AggregationStatus         string
		Total                     int
		Confirmed                 int
		Refuted                   int
		Unknown                   int
		AggregationEvidenceDigest string
	}{
		FeedbackStatus:            output.FeedbackStatus,
		CandidateDigest:           output.CandidateDigest,
		ReplayObservationDigest:   output.ReplayObservationDigest,
		MetricDigest:              output.MetricDigest,
		FeedbackEvidenceDigest:    output.FeedbackEvidenceDigest,
		AggregationStatus:          output.AggregationStatus,
		Total:                      output.Total,
		Confirmed:                  output.Confirmed,
		Refuted:                    output.Refuted,
		Unknown:                    output.Unknown,
		AggregationEvidenceDigest: output.AggregationEvidenceDigest,
	})
	if err := output.Validate(); err != nil {
		return unknown("evaluation-feedback-evidence")
	}
	return output
}
