package decision

import "strings"

// ExecutionEnvelopeJEVReviewLedgerMetricBindingInput connects a verified
// review-ledger aggregation to a named metric source without claiming that the
// metric was improved or applied.
type ExecutionEnvelopeJEVReviewLedgerMetricBindingInput struct {
	Aggregation        ExecutionEnvelopeJEVImprovementReviewLedgerAggregation
	MetricName         string
	MetricValueDigest  string
	MetricSourceDigest string
	NonAuthorizing     bool
}

// ExecutionEnvelopeJEVReviewLedgerMetricBinding records the source and value
// identity of a metric derived from review evidence.
type ExecutionEnvelopeJEVReviewLedgerMetricBinding struct {
	Status                    string
	MetricName                string
	MetricValueDigest         string
	MetricSourceDigest        string
	InputLedgerDigest         string
	AggregationEvidenceDigest string
	BindingDigest             string
	MissingStage              string
	NonExecuting              bool
	NonAuthorizing            bool
}

// BindExecutionEnvelopeJEVReviewLedgerToMetric preserves the ledger evidence
// boundary while making a metric source explicit.
func BindExecutionEnvelopeJEVReviewLedgerToMetric(input ExecutionEnvelopeJEVReviewLedgerMetricBindingInput) ExecutionEnvelopeJEVReviewLedgerMetricBinding {
	output := ExecutionEnvelopeJEVReviewLedgerMetricBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Aggregation.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.Aggregation.NonExecuting {
		output.MissingStage = "execution-boundary"
		return output
	}
	if err := input.Aggregation.Validate(); err != nil {
		output.MissingStage = input.Aggregation.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "review-ledger-aggregation"
		}
		return output
	}
	if strings.TrimSpace(input.MetricName) == "" {
		output.MissingStage = "metric-name"
		return output
	}
	if strings.TrimSpace(input.MetricSourceDigest) == "" {
		output.MissingStage = "metric-source"
		return output
	}
	if strings.TrimSpace(input.MetricValueDigest) == "" {
		output.MissingStage = "metric-value"
		return output
	}
	bindingDigest, err := Digest(struct {
		MetricName                string
		MetricValueDigest         string
		MetricSourceDigest        string
		InputLedgerDigest         string
		AggregationEvidenceDigest string
	}{
		MetricName:                input.MetricName,
		MetricValueDigest:         input.MetricValueDigest,
		MetricSourceDigest:        input.MetricSourceDigest,
		InputLedgerDigest:         input.Aggregation.InputLedgerDigest,
		AggregationEvidenceDigest: input.Aggregation.EvidenceDigest,
	})
	if err != nil {
		output.MissingStage = "metric-binding-digest"
		return output
	}
	output.Status = "bound"
	output.MetricName = input.MetricName
	output.MetricValueDigest = input.MetricValueDigest
	output.MetricSourceDigest = input.MetricSourceDigest
	output.InputLedgerDigest = input.Aggregation.InputLedgerDigest
	output.AggregationEvidenceDigest = input.Aggregation.EvidenceDigest
	output.BindingDigest = bindingDigest
	return output
}
