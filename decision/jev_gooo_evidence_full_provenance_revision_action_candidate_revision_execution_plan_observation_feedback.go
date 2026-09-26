package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackInput
// binds explicit external feedback to an observed execution-plan metric.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackInput struct {
	Metric            ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricBinding
	FeedbackChoice    string
	FeedbackEvidenceDigest string
	NonAuthorizing    bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackBinding
// preserves feedback, aggregation counts, metric, and reverse-observation provenance.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackBinding struct {
	Status                    string
	MissingStage              string
	FeedbackStatus            string
	ObservationStatus         string
	LifecycleDisposition      string
	CandidateStatus           string
	CandidateDigest           string
	PermissionState           string
	MetricName                string
	MetricValue               int
	MetricUnit                string
	MetricEvidenceDigest      string
	FeedbackChoice             string
	FeedbackEvidenceDigest     string
	AggregationStatus          string
	Total                     int
	Confirmed                 int
	Refuted                   int
	Unknown                   int
	AggregationEvidenceDigest string
	EvidencePrefixDigest      string
	ReverseObservationDigest  string
	EvidenceDigest            string
	NonExecuting              bool
	NonAuthorizing            bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackBinding) Validate() error {
	if b.Status != "bound" ||
		b.MissingStage != "" ||
		b.FeedbackStatus == "" ||
		b.ObservationStatus == "" ||
		b.LifecycleDisposition == "" ||
		b.CandidateStatus != jevImprovementRevisionCandidateReady ||
		b.CandidateDigest == "" ||
		b.PermissionState != "not-authorized" ||
		b.MetricName == "" ||
		b.MetricValue < 0 ||
		b.MetricUnit == "" ||
		b.MetricEvidenceDigest == "" ||
		b.FeedbackChoice == "" ||
		b.FeedbackEvidenceDigest == "" ||
		b.AggregationStatus != "aggregated" ||
		b.Total != 1 ||
		b.Confirmed+b.Refuted+b.Unknown != b.Total ||
		b.AggregationEvidenceDigest == "" ||
		b.EvidencePrefixDigest == "" ||
		b.ReverseObservationDigest == "" ||
		b.EvidenceDigest == "" {
		return fmt.Errorf("incomplete Gooo extended lineage execution plan observation feedback binding")
	}
	expectedStatus, expectedConfirmed, expectedRefuted, expectedUnknown := feedbackDispositionForGoooExtendedLineageExecutionPlan(b.FeedbackChoice)
	if expectedStatus == "" ||
		b.FeedbackStatus != expectedStatus ||
		b.Confirmed != expectedConfirmed ||
		b.Refuted != expectedRefuted ||
		b.Unknown != expectedUnknown {
		return fmt.Errorf("invalid Gooo extended lineage execution plan feedback disposition")
	}
	expectedAggregationDigest := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackAggregation(
		b.CandidateDigest,
		b.MetricEvidenceDigest,
		b.FeedbackChoice,
		b.FeedbackEvidenceDigest,
		b.ReverseObservationDigest,
	)
	if b.AggregationEvidenceDigest != expectedAggregationDigest {
		return fmt.Errorf("Gooo extended lineage execution plan aggregation digest mismatch")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("Gooo extended lineage execution plan observation feedback must be non-executing and non-authorizing")
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedback(b)
	if b.EvidenceDigest != expected {
		return fmt.Errorf("Gooo extended lineage execution plan observation feedback digest mismatch")
	}
	return nil
}

// BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedback
// closes the metric-to-feedback loop without inferring feedback from the metric.
func BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedback(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-candidate-generation-extended-lineage-candidate-execution-plan-observation-feedback"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackBinding{
			Status:         "UNKNOWN",
			MissingStage:   stage,
			PermissionState: "not-authorized",
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
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-execution-plan-observation-metric-validation")
	}
	status, confirmed, refuted, unknownCount := feedbackDispositionForGoooExtendedLineageExecutionPlan(input.FeedbackChoice)
	if status == "" {
		return unknown("feedback-choice")
	}
	if strings.TrimSpace(input.FeedbackEvidenceDigest) == "" {
		return unknown("feedback-evidence")
	}
	aggregationDigest := digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackAggregation(
		input.Metric.CandidateDigest,
		input.Metric.MetricEvidenceDigest,
		input.FeedbackChoice,
		input.FeedbackEvidenceDigest,
		input.Metric.ReverseObservationDigest,
	)
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackBinding{
		Status:                    "bound",
		FeedbackStatus:            status,
		ObservationStatus:         input.Metric.ObservationStatus,
		LifecycleDisposition:      input.Metric.LifecycleDisposition,
		CandidateStatus:           input.Metric.CandidateStatus,
		CandidateDigest:           input.Metric.CandidateDigest,
		PermissionState:           "not-authorized",
		MetricName:                input.Metric.MetricName,
		MetricValue:               input.Metric.MetricValue,
		MetricUnit:                input.Metric.MetricUnit,
		MetricEvidenceDigest:      input.Metric.MetricEvidenceDigest,
		FeedbackChoice:             strings.TrimSpace(input.FeedbackChoice),
		FeedbackEvidenceDigest:     input.FeedbackEvidenceDigest,
		AggregationStatus:          "aggregated",
		Total:                     1,
		Confirmed:                 confirmed,
		Refuted:                   refuted,
		Unknown:                   unknownCount,
		AggregationEvidenceDigest: aggregationDigest,
		EvidencePrefixDigest:      input.Metric.EvidencePrefixDigest,
		ReverseObservationDigest:  input.Metric.ReverseObservationDigest,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedback(output)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-candidate-generation-extended-lineage-candidate-execution-plan-observation-feedback-evidence")
	}
	return output
}

func feedbackDispositionForGoooExtendedLineageExecutionPlan(choice string) (string, int, int, int) {
	switch strings.TrimSpace(choice) {
	case "confirmed":
		return "confirmed", 1, 0, 0
	case "refuted":
		return "refuted", 0, 1, 0
	case "unknown":
		return "unknown", 0, 0, 1
	default:
		return "", 0, 0, 0
	}
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackAggregation(candidateDigest, metricEvidenceDigest, feedbackChoice, feedbackEvidenceDigest, reverseObservationDigest string) string {
	digest, err := Digest(struct {
		CandidateDigest          string
		MetricEvidenceDigest     string
		FeedbackChoice           string
		FeedbackEvidenceDigest   string
		ReverseObservationDigest string
	}{
		CandidateDigest:          candidateDigest,
		MetricEvidenceDigest:     metricEvidenceDigest,
		FeedbackChoice:           feedbackChoice,
		FeedbackEvidenceDigest:   feedbackEvidenceDigest,
		ReverseObservationDigest: reverseObservationDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedback(b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackBinding) string {
	digest, err := Digest(struct {
		Status                    string
		FeedbackStatus            string
		ObservationStatus         string
		LifecycleDisposition      string
		CandidateStatus           string
		CandidateDigest           string
		PermissionState           string
		MetricName                string
		MetricValue               int
		MetricUnit                string
		MetricEvidenceDigest      string
		FeedbackChoice             string
		FeedbackEvidenceDigest     string
		AggregationStatus          string
		Total                     int
		Confirmed                 int
		Refuted                   int
		Unknown                   int
		AggregationEvidenceDigest string
		EvidencePrefixDigest      string
		ReverseObservationDigest  string
		NonExecuting              bool
		NonAuthorizing            bool
	}{
		Status:                    b.Status,
		FeedbackStatus:            b.FeedbackStatus,
		ObservationStatus:         b.ObservationStatus,
		LifecycleDisposition:      b.LifecycleDisposition,
		CandidateStatus:           b.CandidateStatus,
		CandidateDigest:           b.CandidateDigest,
		PermissionState:            b.PermissionState,
		MetricName:                b.MetricName,
		MetricValue:               b.MetricValue,
		MetricUnit:                b.MetricUnit,
		MetricEvidenceDigest:      b.MetricEvidenceDigest,
		FeedbackChoice:             b.FeedbackChoice,
		FeedbackEvidenceDigest:    b.FeedbackEvidenceDigest,
		AggregationStatus:          b.AggregationStatus,
		Total:                     b.Total,
		Confirmed:                 b.Confirmed,
		Refuted:                   b.Refuted,
		Unknown:                   b.Unknown,
		AggregationEvidenceDigest: b.AggregationEvidenceDigest,
		EvidencePrefixDigest:      b.EvidencePrefixDigest,
		ReverseObservationDigest:  b.ReverseObservationDigest,
		NonExecuting:              b.NonExecuting,
		NonAuthorizing:            b.NonAuthorizing,
	})
	if err != nil {
		return ""
	}
	return digest
}