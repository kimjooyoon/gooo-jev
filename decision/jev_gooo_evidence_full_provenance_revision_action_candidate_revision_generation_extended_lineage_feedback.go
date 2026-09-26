package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackInput
// adapts the extended lineage candidate metric into replay feedback.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackInput struct {
	Metric                  ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageMetricBinding
	CandidateDigest         string
	ReplayObservationDigest string
	NonAuthorizing          bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackBinding
// preserves extended lineage, feedback, and aggregation evidence.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackBinding struct {
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
	FeedbackStatus                  string
	CandidateDigest                 string
	ReplayObservationDigest         string
	FeedbackEvidenceDigest          string
	FeedbackDigest                  string
	AggregationStatus               string
	Total                           int
	Confirmed                       int
	Refuted                         int
	Unknown                         int
	AggregationEvidenceDigest       string
	EvidenceDigest                  string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackBinding) Validate() error {
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
		return fmt.Errorf("incomplete Gooo extended lineage candidate feedback binding")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage candidate feedback must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedback(
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
		return fmt.Errorf("Gooo extended lineage candidate feedback digest mismatch")
	}
	return nil
}

// BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedback
// keeps review and unknown dispositions conservative before direction derivation.
func BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedback(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-generation-extended-lineage-feedback"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackBinding{
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
		return unknown("revision-action-candidate-generation-extended-lineage-metric-validation")
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
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackBinding{
		Status:                          "bound",
		CandidateID:                     input.Metric.CandidateID,
		SourceCandidateID:               input.Metric.SourceCandidateID,
		SourceRevisionCandidateDigest:   input.Metric.SourceRevisionCandidateDigest,
		SourceCandidateEvidenceDigest:   input.Metric.SourceCandidateEvidenceDigest,
		RevisionCandidateDigest:         input.Metric.RevisionCandidateDigest,
		RevisionCandidateEvidenceDigest: input.Metric.RevisionCandidateEvidenceDigest,
		GuardEvidenceDigest:             input.Metric.GuardEvidenceDigest,
		VerificationEvidenceDigest:      input.Metric.VerificationEvidenceDigest,
		ReverseEvidenceDigest:           input.Metric.ReverseEvidenceDigest,
		EvidencePrefixDigest:            input.Metric.EvidencePrefixDigest,
		MetricEvidenceDigest:            input.Metric.MetricEvidenceDigest,
		FeedbackStatus:                  feedbackStatus,
		CandidateDigest:                 input.CandidateDigest,
		ReplayObservationDigest:         input.ReplayObservationDigest,
		FeedbackEvidenceDigest:          feedback.FeedbackEvidenceDigest,
		FeedbackDigest:                  feedback.EvidenceDigest,
		AggregationStatus:               aggregation.Status,
		Total:                           aggregation.Total,
		Confirmed:                       aggregation.Confirmed,
		Refuted:                         aggregation.Refuted,
		Unknown:                         aggregation.Unknown,
		AggregationEvidenceDigest:       aggregation.EvidenceDigest,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedback(
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
		return unknown("revision-action-candidate-generation-extended-lineage-feedback-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedback(
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
		FeedbackStatus                  string
		CandidateDigest                 string
		ReplayObservationDigest         string
		FeedbackEvidenceDigest          string
		FeedbackDigest                  string
		AggregationStatus               string
		Total                           int
		Confirmed                       int
		Refuted                         int
		Unknown                         int
		AggregationEvidenceDigest       string
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
		FeedbackStatus:                  feedbackStatus,
		CandidateDigest:                 candidateDigest,
		ReplayObservationDigest:         replayObservationDigest,
		FeedbackEvidenceDigest:          feedbackEvidenceDigest,
		FeedbackDigest:                  feedbackDigest,
		AggregationStatus:               aggregationStatus,
		Total:                           total,
		Confirmed:                       confirmed,
		Refuted:                         refuted,
		Unknown:                         unknown,
		AggregationEvidenceDigest:       aggregationEvidenceDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}
