package gooo

import "fmt"

const (
	jevExecutionPlanCoverageLSPUnknown  = "jev-execution-plan-coverage-lsp-unknown"
	jevExecutionPlanCoverageLSPComplete = "jev-execution-plan-coverage-lsp-complete"
)

type RevisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricLSPProjection struct {
	Status                   string
	MissingStage             string
	MissingStageIndex        int
	MetricName               string
	CoverageMilli            int
	CoverageBand             string
	MetricDigest             string
	EvidencePrefixDigest     string
	ReverseObservationDigest string
	CoverageSignal           string
	ObservationDigest        string
	NonExecuting             bool
	NonAuthorizing           bool
}

func ObserveRevisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricLSPProjection(
	input RevisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetric,
) RevisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricLSPProjection {
	output := RevisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricLSPProjection{
		Status:            "UNKNOWN",
		MissingStage:      input.MissingStage,
		MissingStageIndex: -1,
		MetricName:        jevExecutionPlanCoverageMetricName,
		CoverageBand:      "unknown",
		CoverageSignal:    jevExecutionPlanCoverageLSPUnknown,
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "execution_plan_reverse_observation_coverage_metric"
	}
	if err := input.Validate(); err != nil {
		output.ObservationDigest = revisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricLSPProjectionDigest(output)
		return output
	}

	output.EvidencePrefixDigest = input.ObservationDigest
	if input.Status != "BOUND" {
		output.ObservationDigest = revisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricLSPProjectionDigest(output)
		return output
	}

	output.Status = "BOUND"
	output.MissingStage = ""
	output.MissingStageIndex = 0
	output.CoverageMilli = input.CoverageMilli
	output.CoverageBand = input.CoverageBand
	output.MetricDigest = input.MetricDigest
	output.ReverseObservationDigest = input.ReverseObservationDigest
	output.CoverageSignal = jevExecutionPlanCoverageLSPComplete
	output.ObservationDigest = revisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricLSPProjectionDigest(output)
	return output
}

func (value RevisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricLSPProjection) Validate() error {
	if value.Status != "UNKNOWN" && value.Status != "BOUND" {
		return fmt.Errorf("invalid status %q", value.Status)
	}
	if value.MetricName != jevExecutionPlanCoverageMetricName {
		return fmt.Errorf("invalid metric name %q", value.MetricName)
	}
	if value.MissingStageIndex < -1 {
		return fmt.Errorf("invalid missing stage index %d", value.MissingStageIndex)
	}
	if !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("LSP coverage projection must remain non-executing and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" || value.MissingStageIndex != -1 {
			return fmt.Errorf("unknown LSP coverage projection must preserve an unresolved stage")
		}
		if value.CoverageMilli != 0 || value.CoverageBand != "unknown" {
			return fmt.Errorf("unknown LSP coverage projection must not claim coverage")
		}
		if value.CoverageSignal != jevExecutionPlanCoverageLSPUnknown {
			return fmt.Errorf("invalid unknown LSP coverage signal %q", value.CoverageSignal)
		}
		if value.EvidencePrefixDigest != "" && !validJEVDecisionConfidenceClosureDigest(value.EvidencePrefixDigest) {
			return fmt.Errorf("invalid unknown evidence prefix digest")
		}
	} else {
		if value.MissingStage != "" || value.MissingStageIndex != 0 {
			return fmt.Errorf("bound LSP coverage projection must close the stage")
		}
		if value.CoverageMilli != 1000 || value.CoverageBand != "high" {
			return fmt.Errorf("bound LSP coverage projection must preserve complete coverage")
		}
		if !validJEVDecisionConfidenceClosureDigest(value.MetricDigest) {
			return fmt.Errorf("invalid metric digest")
		}
		if !validJEVDecisionConfidenceClosureDigest(value.EvidencePrefixDigest) {
			return fmt.Errorf("invalid evidence prefix digest")
		}
		if !validJEVDecisionConfidenceClosureDigest(value.ReverseObservationDigest) {
			return fmt.Errorf("invalid reverse observation digest")
		}
		if value.CoverageSignal != jevExecutionPlanCoverageLSPComplete {
			return fmt.Errorf("invalid bound LSP coverage signal %q", value.CoverageSignal)
		}
	}
	expected := revisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricLSPProjectionDigest(value)
	if value.ObservationDigest != expected {
		return fmt.Errorf("observation digest mismatch")
	}
	return nil
}

func revisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricLSPProjectionDigest(
	value RevisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricLSPProjection,
) string {
	return digestString(fmt.Sprintf(
		"jev-execution-plan-coverage-lsp|%s|%s|%d|%s|%d|%s|%s|%s|%s|%s|%t|%t",
		value.Status,
		value.MissingStage,
		value.MissingStageIndex,
		value.MetricName,
		value.CoverageMilli,
		value.CoverageBand,
		value.MetricDigest,
		value.EvidencePrefixDigest,
		value.ReverseObservationDigest,
		value.CoverageSignal,
		value.NonExecuting,
		value.NonAuthorizing,
	))
}
