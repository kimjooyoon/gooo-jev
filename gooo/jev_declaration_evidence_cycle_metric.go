package gooo

import "fmt"

const (
	ExecutionEnvelopeDeclarationEvidenceCycleMetricBound    = "BOUND"
	ExecutionEnvelopeDeclarationEvidenceCycleMetricDeferred = "DEFERRED"
	ExecutionEnvelopeDeclarationEvidenceCycleMetricUnknown  = "UNKNOWN"
	ExecutionEnvelopeDeclarationEvidenceCycleMetricError    = "ERROR"

	executionEnvelopeDeclarationEvidenceCycleMetricStageTotal = 3
)

// ExecutionEnvelopeDeclarationEvidenceCycleMetric is a read-only coverage
// observation over the declaration, IR/generation, and reverse stages.
type ExecutionEnvelopeDeclarationEvidenceCycleMetric struct {
	Status            string
	StageCount        int
	StageTotal        int
	Coverage          float64
	FirstMissingStage string
	Reason            string
	EvidenceDigest    string
	MetricDigest      string
	NonExecuting      bool
	NonAuthorizing    bool
}

// MeasureExecutionEnvelopeDeclarationEvidenceCycleMetric converts an already
// projected cycle into a coverage metric without re-executing any stage.
func MeasureExecutionEnvelopeDeclarationEvidenceCycleMetric(
	projection ExecutionEnvelopeDeclarationEvidenceCycleProjection,
) ExecutionEnvelopeDeclarationEvidenceCycleMetric {
	metric := ExecutionEnvelopeDeclarationEvidenceCycleMetric{
		Status:            ExecutionEnvelopeDeclarationEvidenceCycleMetricUnknown,
		StageTotal:        executionEnvelopeDeclarationEvidenceCycleMetricStageTotal,
		FirstMissingStage: projection.MissingStage,
		Reason:            projection.FirstMismatch,
		EvidenceDigest:    projection.EvidenceDigest,
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	if metric.FirstMissingStage == "" {
		metric.FirstMissingStage = "declaration_source"
	}
	if metric.Reason == "" {
		metric.Reason = "evidence cycle remains unresolved"
	}
	if err := projection.Validate(); err != nil {
		metric.Status = ExecutionEnvelopeDeclarationEvidenceCycleMetricError
		metric.StageCount = 0
		metric.Coverage = 0
		metric.FirstMissingStage = "projection_integrity"
		metric.Reason = err.Error()
		return finalizeExecutionEnvelopeDeclarationEvidenceCycleMetric(metric)
	}

	switch projection.Status {
	case ExecutionEnvelopeDeclarationEvidenceCycleBound:
		metric.Status = ExecutionEnvelopeDeclarationEvidenceCycleMetricBound
		metric.StageCount = executionEnvelopeDeclarationEvidenceCycleMetricStageTotal
		metric.Coverage = 1
		metric.FirstMissingStage = ""
		metric.Reason = "declaration, IR generation, and reverse observation are available"
	case ExecutionEnvelopeDeclarationEvidenceCycleDeferred:
		metric.Status = ExecutionEnvelopeDeclarationEvidenceCycleMetricDeferred
		metric.StageCount = 2
		metric.Coverage = 2.0 / 3.0
		metric.FirstMissingStage = "reverse_observation"
		metric.Reason = "reverse observation is deferred"
	case ExecutionEnvelopeDeclarationEvidenceCycleUnknown:
		metric.Status = ExecutionEnvelopeDeclarationEvidenceCycleMetricUnknown
		metric.StageCount = executionEnvelopeDeclarationEvidenceCycleMetricStageCount(projection.MissingStage)
		metric.Coverage = float64(metric.StageCount) / float64(metric.StageTotal)
	case ExecutionEnvelopeDeclarationEvidenceCycleError:
		metric.Status = ExecutionEnvelopeDeclarationEvidenceCycleMetricError
		metric.StageCount = 0
		metric.Coverage = 0
	default:
		metric.Status = ExecutionEnvelopeDeclarationEvidenceCycleMetricError
		metric.StageCount = 0
		metric.Coverage = 0
		metric.FirstMissingStage = "projection_status"
		metric.Reason = "declaration evidence cycle status is not recognized"
	}
	return finalizeExecutionEnvelopeDeclarationEvidenceCycleMetric(metric)
}

func executionEnvelopeDeclarationEvidenceCycleMetricStageCount(missingStage string) int {
	switch missingStage {
	case "declaration_ir_generation":
		return 1
	case "reverse_observation":
		return 2
	default:
		return 0
	}
}

func finalizeExecutionEnvelopeDeclarationEvidenceCycleMetric(
	metric ExecutionEnvelopeDeclarationEvidenceCycleMetric,
) ExecutionEnvelopeDeclarationEvidenceCycleMetric {
	metric.MetricDigest = digestString(fmt.Sprintf(
		"gooo-declaration-evidence-cycle-metric|%s|%d|%d|%.17g|%s|%s|%s|%t|%t",
		metric.Status,
		metric.StageCount,
		metric.StageTotal,
		metric.Coverage,
		metric.FirstMissingStage,
		metric.Reason,
		metric.EvidenceDigest,
		metric.NonExecuting,
		metric.NonAuthorizing,
	))
	return metric
}

// Validate checks metric integrity and the non-executing, non-authorizing
// boundary. It does not interpret coverage as correctness or improvement.
func (metric ExecutionEnvelopeDeclarationEvidenceCycleMetric) Validate() error {
	if metric.Status != ExecutionEnvelopeDeclarationEvidenceCycleMetricBound &&
		metric.Status != ExecutionEnvelopeDeclarationEvidenceCycleMetricDeferred &&
		metric.Status != ExecutionEnvelopeDeclarationEvidenceCycleMetricUnknown &&
		metric.Status != ExecutionEnvelopeDeclarationEvidenceCycleMetricError {
		return fmt.Errorf("invalid declaration evidence cycle metric status %q", metric.Status)
	}
	if metric.StageTotal != executionEnvelopeDeclarationEvidenceCycleMetricStageTotal ||
		metric.StageCount < 0 || metric.StageCount > metric.StageTotal ||
		metric.Coverage < 0 || metric.Coverage > 1 {
		return fmt.Errorf("declaration evidence cycle metric coverage is invalid")
	}
	if !metric.NonExecuting || !metric.NonAuthorizing {
		return fmt.Errorf("declaration evidence cycle metric crossed a capability boundary")
	}
	if metric.MetricDigest != digestString(fmt.Sprintf(
		"gooo-declaration-evidence-cycle-metric|%s|%d|%d|%.17g|%s|%s|%s|%t|%t",
		metric.Status,
		metric.StageCount,
		metric.StageTotal,
		metric.Coverage,
		metric.FirstMissingStage,
		metric.Reason,
		metric.EvidenceDigest,
		metric.NonExecuting,
		metric.NonAuthorizing,
	)) {
		return fmt.Errorf("declaration evidence cycle metric digest mismatch")
	}
	if metric.Status == ExecutionEnvelopeDeclarationEvidenceCycleMetricBound &&
		(metric.StageCount != metric.StageTotal || metric.FirstMissingStage != "" || metric.Coverage != 1) {
		return fmt.Errorf("bound declaration evidence cycle metric is incomplete")
	}
	if metric.Status == ExecutionEnvelopeDeclarationEvidenceCycleMetricUnknown && metric.FirstMissingStage == "" {
		return fmt.Errorf("unknown declaration evidence cycle metric lost its first missing stage")
	}
	return nil
}