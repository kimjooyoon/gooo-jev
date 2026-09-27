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
	CoverageMilli     int
	FirstMissingStage string
	MissingStageIndex int
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
		MissingStageIndex: declarationEvidenceCycleMetricMissingStageIndex(projection.MissingStage),
		Reason:            projection.FirstMismatch,
		EvidenceDigest:    projection.EvidenceDigest,
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	if metric.FirstMissingStage == "" {
		metric.FirstMissingStage = "declaration_source"
		metric.MissingStageIndex = declarationEvidenceCycleMetricMissingStageIndex(metric.FirstMissingStage)
	}
	if metric.Reason == "" {
		metric.Reason = "evidence cycle remains unresolved"
	}
	if err := projection.Validate(); err != nil {
		metric.Status = ExecutionEnvelopeDeclarationEvidenceCycleMetricError
		metric.StageCount = 0
		metric.Coverage = 0
		metric.CoverageMilli = 0
		metric.FirstMissingStage = "projection_integrity"
		metric.MissingStageIndex = -1
		metric.Reason = err.Error()
		return finalizeExecutionEnvelopeDeclarationEvidenceCycleMetric(metric)
	}

	switch projection.Status {
	case ExecutionEnvelopeDeclarationEvidenceCycleBound:
		metric.Status = ExecutionEnvelopeDeclarationEvidenceCycleMetricBound
		metric.StageCount = executionEnvelopeDeclarationEvidenceCycleMetricStageTotal
		metric.Coverage = 1
		metric.CoverageMilli = 1000
		metric.FirstMissingStage = ""
		metric.MissingStageIndex = -1
		metric.Reason = "declaration, IR generation, and reverse observation are available"
	case ExecutionEnvelopeDeclarationEvidenceCycleDeferred:
		metric.Status = ExecutionEnvelopeDeclarationEvidenceCycleMetricDeferred
		metric.StageCount = 2
		metric.Coverage = 2.0 / 3.0
		metric.CoverageMilli = 666
		metric.FirstMissingStage = "reverse_observation"
		metric.MissingStageIndex = declarationEvidenceCycleMetricMissingStageIndex(metric.FirstMissingStage)
		metric.Reason = "reverse observation is deferred"
	case ExecutionEnvelopeDeclarationEvidenceCycleUnknown:
		metric.Status = ExecutionEnvelopeDeclarationEvidenceCycleMetricUnknown
		metric.StageCount = executionEnvelopeDeclarationEvidenceCycleMetricStageCount(projection.MissingStage)
		metric.Coverage = float64(metric.StageCount) / float64(metric.StageTotal)
		metric.CoverageMilli = metric.StageCount * 1000 / metric.StageTotal
		metric.MissingStageIndex = declarationEvidenceCycleMetricMissingStageIndex(metric.FirstMissingStage)
	case ExecutionEnvelopeDeclarationEvidenceCycleError:
		metric.Status = ExecutionEnvelopeDeclarationEvidenceCycleMetricError
		metric.StageCount = 0
		metric.Coverage = 0
		metric.CoverageMilli = 0
		metric.MissingStageIndex = -1
	default:
		metric.Status = ExecutionEnvelopeDeclarationEvidenceCycleMetricError
		metric.StageCount = 0
		metric.Coverage = 0
		metric.CoverageMilli = 0
		metric.FirstMissingStage = "projection_status"
		metric.MissingStageIndex = -1
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

func declarationEvidenceCycleMetricMissingStageIndex(missingStage string) int {
	switch missingStage {
	case "declaration_source":
		return 0
	case "declaration_ir_generation":
		return 1
	case "reverse_observation":
		return 2
	default:
		return -1
	}
}

func finalizeExecutionEnvelopeDeclarationEvidenceCycleMetric(
	metric ExecutionEnvelopeDeclarationEvidenceCycleMetric,
) ExecutionEnvelopeDeclarationEvidenceCycleMetric {
	metric.MetricDigest = digestString(fmt.Sprintf(
		"gooo-declaration-evidence-cycle-metric|%s|%d|%d|%.17g|%d|%s|%d|%s|%s|%t|%t",
		metric.Status,
		metric.StageCount,
		metric.StageTotal,
		metric.Coverage,
		metric.CoverageMilli,
		metric.FirstMissingStage,
		metric.MissingStageIndex,
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
		metric.Coverage < 0 || metric.Coverage > 1 ||
		metric.CoverageMilli < 0 || metric.CoverageMilli > 1000 ||
		metric.CoverageMilli != metric.StageCount*1000/metric.StageTotal {
		return fmt.Errorf("declaration evidence cycle metric coverage is invalid")
	}
	if metric.FirstMissingStage == "" && metric.MissingStageIndex != -1 {
		return fmt.Errorf("bound declaration evidence cycle metric has a missing stage index")
	}
	if metric.FirstMissingStage != "" && metric.Status != ExecutionEnvelopeDeclarationEvidenceCycleMetricError && metric.MissingStageIndex < 0 {
		return fmt.Errorf("declaration evidence cycle metric lost its missing stage index")
	}
	if !metric.NonExecuting || !metric.NonAuthorizing {
		return fmt.Errorf("declaration evidence cycle metric crossed a capability boundary")
	}
	if metric.MetricDigest != digestString(fmt.Sprintf(
		"gooo-declaration-evidence-cycle-metric|%s|%d|%d|%.17g|%d|%s|%d|%s|%s|%t|%t",
		metric.Status,
		metric.StageCount,
		metric.StageTotal,
		metric.Coverage,
		metric.CoverageMilli,
		metric.FirstMissingStage,
		metric.MissingStageIndex,
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
