package decision

import "strings"

// ExecutionEnvelopeJEVSelfImprovementEvidenceSummaryInput joins the complete
// observed improvement path without treating the summary as an authorization.
type ExecutionEnvelopeJEVSelfImprovementEvidenceSummaryInput struct {
	CycleBinding     ExecutionEnvelopeJEVPlanImprovementCycleBinding
	LedgerAggregation ExecutionEnvelopeJEVImprovementReviewLedgerAggregation
	MetricBinding    ExecutionEnvelopeJEVReviewLedgerMetricBinding
	NonAuthorizing   bool
}

// ExecutionEnvelopeJEVSelfImprovementEvidenceSummary is the stable identity of
// cycle, review ledger, and metric provenance for downstream tooling.
type ExecutionEnvelopeJEVSelfImprovementEvidenceSummary struct {
	Status                 string
	CycleStatus            string
	CycleEvidenceDigest    string
	LedgerStatus           string
	LedgerEvidenceDigest   string
	MetricName             string
	MetricSourceDigest     string
	MetricValueDigest      string
	MetricBindingDigest    string
	SummaryDigest          string
	MissingStage           string
	NonExecuting           bool
	NonAuthorizing         bool
}

// BindExecutionEnvelopeJEVSelfImprovementEvidenceSummary binds all observed
// stages into one digest while preserving the first incomplete stage.
func BindExecutionEnvelopeJEVSelfImprovementEvidenceSummary(input ExecutionEnvelopeJEVSelfImprovementEvidenceSummaryInput) ExecutionEnvelopeJEVSelfImprovementEvidenceSummary {
	output := ExecutionEnvelopeJEVSelfImprovementEvidenceSummary{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.CycleBinding.NonAuthorizing || !input.LedgerAggregation.NonAuthorizing || !input.MetricBinding.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return finalizeExecutionEnvelopeJEVSelfImprovementEvidenceSummary(output)
	}
	if !input.CycleBinding.NonExecuting || !input.LedgerAggregation.NonExecuting || !input.MetricBinding.NonExecuting {
		output.MissingStage = "execution-boundary"
		return finalizeExecutionEnvelopeJEVSelfImprovementEvidenceSummary(output)
	}
	if input.CycleBinding.Status != "bound" || strings.TrimSpace(input.CycleBinding.BindingDigest) == "" || strings.TrimSpace(input.CycleBinding.CycleEvidenceDigest) == "" {
		output.MissingStage = input.CycleBinding.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "improvement-cycle"
		}
		return finalizeExecutionEnvelopeJEVSelfImprovementEvidenceSummary(output)
	}
	if err := input.LedgerAggregation.Validate(); err != nil {
		output.MissingStage = input.LedgerAggregation.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "review-ledger-aggregation"
		}
		return finalizeExecutionEnvelopeJEVSelfImprovementEvidenceSummary(output)
	}
	if input.MetricBinding.Status != "bound" || strings.TrimSpace(input.MetricBinding.BindingDigest) == "" {
		output.MissingStage = input.MetricBinding.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "metric-binding"
		}
		return finalizeExecutionEnvelopeJEVSelfImprovementEvidenceSummary(output)
	}
	if strings.TrimSpace(input.MetricBinding.MetricSourceDigest) == "" || strings.TrimSpace(input.MetricBinding.MetricValueDigest) == "" || strings.TrimSpace(input.MetricBinding.AggregationEvidenceDigest) == "" {
		output.MissingStage = "metric-evidence"
		return finalizeExecutionEnvelopeJEVSelfImprovementEvidenceSummary(output)
	}
	output.Status = "bound"
	output.CycleStatus = input.CycleBinding.CycleStatus
	output.CycleEvidenceDigest = input.CycleBinding.CycleEvidenceDigest
	output.LedgerStatus = input.LedgerAggregation.Status
	output.LedgerEvidenceDigest = input.LedgerAggregation.EvidenceDigest
	output.MetricName = input.MetricBinding.MetricName
	output.MetricSourceDigest = input.MetricBinding.MetricSourceDigest
	output.MetricValueDigest = input.MetricBinding.MetricValueDigest
	output.MetricBindingDigest = input.MetricBinding.BindingDigest
	return finalizeExecutionEnvelopeJEVSelfImprovementEvidenceSummary(output)
}

func digestExecutionEnvelopeJEVSelfImprovementEvidenceSummary(summary ExecutionEnvelopeJEVSelfImprovementEvidenceSummary) (string, error) {
	return Digest(struct {
		Status                 string
		CycleStatus            string
		CycleEvidenceDigest    string
		LedgerStatus           string
		LedgerEvidenceDigest   string
		MetricName             string
		MetricSourceDigest     string
		MetricValueDigest      string
		MetricBindingDigest    string
		MissingStage           string
	}{
		Status:                 summary.Status,
		CycleStatus:            summary.CycleStatus,
		CycleEvidenceDigest:    summary.CycleEvidenceDigest,
		LedgerStatus:           summary.LedgerStatus,
		LedgerEvidenceDigest:   summary.LedgerEvidenceDigest,
		MetricName:             summary.MetricName,
		MetricSourceDigest:     summary.MetricSourceDigest,
		MetricValueDigest:      summary.MetricValueDigest,
		MetricBindingDigest:    summary.MetricBindingDigest,
		MissingStage:           summary.MissingStage,
	})
}

func finalizeExecutionEnvelopeJEVSelfImprovementEvidenceSummary(output ExecutionEnvelopeJEVSelfImprovementEvidenceSummary) ExecutionEnvelopeJEVSelfImprovementEvidenceSummary {
	digest, err := digestExecutionEnvelopeJEVSelfImprovementEvidenceSummary(output)
	if err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "self-improvement-summary-digest"
		output.SummaryDigest = ""
		return output
	}
	output.SummaryDigest = digest
	return output
}
