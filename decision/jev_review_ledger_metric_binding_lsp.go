package decision

import "strings"

// ExecutionEnvelopeJEVReviewLedgerMetricBindingLSPInput projects metric
// provenance into deterministic language-tooling evidence.
type ExecutionEnvelopeJEVReviewLedgerMetricBindingLSPInput struct {
	Binding        ExecutionEnvelopeJEVReviewLedgerMetricBinding
	NonAuthorizing bool
}

// ExecutionEnvelopeJEVReviewLedgerMetricBindingLSPProjection exposes metric
// source and value identity without claiming metric improvement.
type ExecutionEnvelopeJEVReviewLedgerMetricBindingLSPProjection struct {
	Status            string
	Code              string
	Severity          string
	Message           string
	MetricName        string
	MetricSourceDigest string
	MetricValueDigest  string
	BindingDigest     string
	EvidenceDigest    string
	MissingStage      string
	NonExecuting      bool
	NonAuthorizing    bool
}

func digestExecutionEnvelopeJEVReviewLedgerMetricBindingLSPProjection(projection ExecutionEnvelopeJEVReviewLedgerMetricBindingLSPProjection) (string, error) {
	return Digest(struct {
		Status            string
		Code              string
		Severity          string
		Message           string
		MetricName        string
		MetricSourceDigest string
		MetricValueDigest  string
		BindingDigest     string
		MissingStage      string
	}{
		Status:            projection.Status,
		Code:              projection.Code,
		Severity:          projection.Severity,
		Message:           projection.Message,
		MetricName:        projection.MetricName,
		MetricSourceDigest: projection.MetricSourceDigest,
		MetricValueDigest:  projection.MetricValueDigest,
		BindingDigest:     projection.BindingDigest,
		MissingStage:      projection.MissingStage,
	})
}

// ProjectExecutionEnvelopeJEVReviewLedgerMetricBindingLSP projects bound
// metric provenance or its first UNKNOWN stage for LSP consumers.
func ProjectExecutionEnvelopeJEVReviewLedgerMetricBindingLSP(input ExecutionEnvelopeJEVReviewLedgerMetricBindingLSPInput) ExecutionEnvelopeJEVReviewLedgerMetricBindingLSPProjection {
	output := ExecutionEnvelopeJEVReviewLedgerMetricBindingLSPProjection{
		Status: "UNKNOWN", Code: "JEV_METRIC_PROVENANCE_UNKNOWN", Severity: "warning",
		NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Binding.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		output.Message = "JEV metric provenance is UNKNOWN: missing authorization boundary"
		return finalizeExecutionEnvelopeJEVReviewLedgerMetricBindingLSP(output)
	}
	if !input.Binding.NonExecuting {
		output.MissingStage = "execution-boundary"
		output.Message = "JEV metric provenance is UNKNOWN: missing execution boundary"
		return finalizeExecutionEnvelopeJEVReviewLedgerMetricBindingLSP(output)
	}
	if input.Binding.Status != "bound" || strings.TrimSpace(input.Binding.BindingDigest) == "" {
		output.MissingStage = input.Binding.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "metric-binding"
		}
		output.Message = "JEV metric provenance is UNKNOWN: missing " + output.MissingStage
		return finalizeExecutionEnvelopeJEVReviewLedgerMetricBindingLSP(output)
	}
	if strings.TrimSpace(input.Binding.MetricName) == "" {
		output.MissingStage = "metric-name"
		output.Message = "JEV metric provenance is UNKNOWN: missing metric-name"
		return finalizeExecutionEnvelopeJEVReviewLedgerMetricBindingLSP(output)
	}
	if strings.TrimSpace(input.Binding.MetricSourceDigest) == "" {
		output.MissingStage = "metric-source"
		output.Message = "JEV metric provenance is UNKNOWN: missing metric-source"
		return finalizeExecutionEnvelopeJEVReviewLedgerMetricBindingLSP(output)
	}
	if strings.TrimSpace(input.Binding.MetricValueDigest) == "" {
		output.MissingStage = "metric-value"
		output.Message = "JEV metric provenance is UNKNOWN: missing metric-value"
		return finalizeExecutionEnvelopeJEVReviewLedgerMetricBindingLSP(output)
	}
	if strings.TrimSpace(input.Binding.AggregationEvidenceDigest) == "" {
		output.MissingStage = "ledger-aggregation-evidence"
		output.Message = "JEV metric provenance is UNKNOWN: missing ledger-aggregation-evidence"
		return finalizeExecutionEnvelopeJEVReviewLedgerMetricBindingLSP(output)
	}
	output.Status = "bound"
	output.Code = "JEV_METRIC_PROVENANCE_BOUND"
	output.Severity = "info"
	output.MetricName = input.Binding.MetricName
	output.MetricSourceDigest = input.Binding.MetricSourceDigest
	output.MetricValueDigest = input.Binding.MetricValueDigest
	output.BindingDigest = input.Binding.BindingDigest
	output.Message = "JEV metric source, value, and ledger evidence are bound"
	return finalizeExecutionEnvelopeJEVReviewLedgerMetricBindingLSP(output)
}

func finalizeExecutionEnvelopeJEVReviewLedgerMetricBindingLSP(output ExecutionEnvelopeJEVReviewLedgerMetricBindingLSPProjection) ExecutionEnvelopeJEVReviewLedgerMetricBindingLSPProjection {
	digest, err := digestExecutionEnvelopeJEVReviewLedgerMetricBindingLSPProjection(output)
	if err != nil {
		output.Status = "UNKNOWN"
		output.Code = "JEV_METRIC_PROVENANCE_UNKNOWN"
		output.Severity = "warning"
		output.MissingStage = "lsp-evidence-digest"
		output.Message = "JEV metric provenance is UNKNOWN: missing lsp-evidence-digest"
		output.EvidenceDigest = ""
		return output
	}
	output.EvidenceDigest = digest
	return output
}
