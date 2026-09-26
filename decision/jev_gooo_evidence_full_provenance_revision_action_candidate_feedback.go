package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackInput
// adapts an action-derived candidate metric into the replay feedback ledger.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackInput struct {
	Metric                  ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricBinding
	CandidateDigest         string
	ReplayObservationDigest string
	NonAuthorizing          bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackBinding
// preserves candidate lineage, feedback, and aggregation evidence.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackBinding struct {
	Status                     string
	MissingStage               string
	CandidateID                string
	RevisionCandidateDigest    string
	CandidateEvidenceDigest    string
	MetricEvidenceDigest       string
	FeedbackStatus             string
	CandidateDigest            string
	ReplayObservationDigest    string
	FeedbackEvidenceDigest     string
	FeedbackDigest             string
	AggregationStatus          string
	Total                      int
	Confirmed                  int
	Refuted                    int
	Unknown                    int
	AggregationEvidenceDigest  string
	EvidenceDigest             string
	NonExecuting               bool
	NonAuthorizing             bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
		b.CandidateID == "" ||
		b.RevisionCandidateDigest == "" ||
		b.CandidateEvidenceDigest == "" ||
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
		return fmt.Errorf("incomplete Gooo action-derived candidate feedback binding")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo action-derived candidate feedback must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateFeedback(
		b.CandidateID,
		b.RevisionCandidateDigest,
		b.CandidateEvidenceDigest,
		b.MetricEvidenceDigest,
		b.FeedbackStatus,
		b.CandidateDigest,
		b.ReplayObservationDigest,
		b.FeedbackEvidenceDigest,
		b.FeedbackDigest,
		b.AggregationStatus,
		b.Total,
		b.Confirmed,
		b.Refuted,
		b.Unknown,
		b.AggregationEvidenceDigest,
	)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo action-derived candidate feedback digest mismatch")
	}
	return nil
}

// BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedback
// keeps review and unknown dispositions conservative before direction derivation.
func BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedback(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-feedback"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackBinding{
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
		return unknown("revision-action-candidate-metric-validation")
	}
	if strings.TrimSpace(input.CandidateDigest) == "" {
		return unknown("candidate-digest")
	}
	if input.CandidateDigest != input.Metric.RevisionCandidateDigest {
		return unknown("candidate-digest-mismatch")
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
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackBinding{
		Status:                    "bound",
		CandidateID:               input.Metric.CandidateID,
		RevisionCandidateDigest:   input.Metric.RevisionCandidateDigest,
		CandidateEvidenceDigest:   input.Metric.CandidateEvidenceDigest,
		MetricEvidenceDigest:      input.Metric.MetricEvidenceDigest,
		FeedbackStatus:             feedbackStatus,
		CandidateDigest:            input.CandidateDigest,
		ReplayObservationDigest:   input.ReplayObservationDigest,
		FeedbackEvidenceDigest:    feedback.FeedbackEvidenceDigest,
		FeedbackDigest:            feedback.EvidenceDigest,
		AggregationStatus:          aggregation.Status,
		Total:                     aggregation.Total,
		Confirmed:                 aggregation.Confirmed,
		Refuted:                   aggregation.Refuted,
		Unknown:                   aggregation.Unknown,
		AggregationEvidenceDigest: aggregation.EvidenceDigest,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateFeedback(
		output.CandidateID,
		output.RevisionCandidateDigest,
		output.CandidateEvidenceDigest,
		output.MetricEvidenceDigest,
		output.FeedbackStatus,
		output.CandidateDigest,
		output.ReplayObservationDigest,
		output.FeedbackEvidenceDigest,
		output.FeedbackDigest,
		output.AggregationStatus,
		output.Total,
		output.Confirmed,
		output.Refuted,
		output.Unknown,
		output.AggregationEvidenceDigest,
	)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-feedback-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateFeedback(
	candidateID,
	revisionCandidateDigest,
	candidateEvidenceDigest,
	metricEvidenceDigest,
	feedbackStatus,
	candidateDigest,
	replayObservationDigest,
	feedbackEvidenceDigest,
	feedbackDigest,
	aggregationStatus string,
	total,
	confirmed,
	refuted,
	unknown int,
	aggregationEvidenceDigest string,
) string {
	digest, err := Digest(struct {
		CandidateID               string
		RevisionCandidateDigest  string
		CandidateEvidenceDigest  string
		MetricEvidenceDigest     string
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
		CandidateID:               candidateID,
		RevisionCandidateDigest:  revisionCandidateDigest,
		CandidateEvidenceDigest:  candidateEvidenceDigest,
		MetricEvidenceDigest:     metricEvidenceDigest,
		FeedbackStatus:            feedbackStatus,
		CandidateDigest:           candidateDigest,
		ReplayObservationDigest:   replayObservationDigest,
		FeedbackEvidenceDigest:    feedbackEvidenceDigest,
		FeedbackDigest:            feedbackDigest,
		AggregationStatus:         aggregationStatus,
		Total:                     total,
		Confirmed:                 confirmed,
		Refuted:                   refuted,
		Unknown:                   unknown,
		AggregationEvidenceDigest: aggregationEvidenceDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}