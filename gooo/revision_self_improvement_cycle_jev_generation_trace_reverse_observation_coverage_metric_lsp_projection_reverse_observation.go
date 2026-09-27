package gooo

import "fmt"

const (
	jevGenerationTraceCoverageLSPReverseMetricName = "jev-generation-trace-coverage-lsp-reverse-observation"
	jevGenerationTraceCoverageLSPReverseUnknown = "jev-generation-trace-coverage-lsp-reverse-unknown"
	jevGenerationTraceCoverageLSPReverseComplete = "jev-generation-trace-coverage-lsp-reverse-complete"
)

type RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjectionReverseObservation struct {
	Status                    string
	MissingStage              string
	MissingStageIndex         int
	MetricName                string
	CoverageMilli             int
	CoverageBand              string
	MetricDigest              string
	EvidencePrefixDigest      string
	ProjectionObservationDigest string
	ReverseSignal             string
	ObservationDigest         string
	NonExecuting              bool
	NonAuthorizing            bool
}

func ObserveRevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjectionReverseObservation(
	input RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjection,
) RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjectionReverseObservation {
	output := RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjectionReverseObservation{
		Status:            "UNKNOWN",
		MissingStage:      input.MissingStage,
		MissingStageIndex: -1,
		MetricName:        jevGenerationTraceCoverageLSPReverseMetricName,
		CoverageBand:      "unknown",
		ReverseSignal:     jevGenerationTraceCoverageLSPReverseUnknown,
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "generation_trace_coverage_lsp_projection"
	}
	if err := input.Validate(); err != nil {
		output.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjectionReverseObservationDigest(output)
		return output
	}
	if input.Status != "BOUND" {
		output.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjectionReverseObservationDigest(output)
		return output
	}

	output.Status = "BOUND"
	output.MissingStage = ""
	output.MissingStageIndex = 0
	output.CoverageMilli = input.CoverageMilli
	output.CoverageBand = input.CoverageBand
	output.MetricDigest = input.MetricDigest
	output.EvidencePrefixDigest = input.EvidencePrefixDigest
	output.ProjectionObservationDigest = input.ObservationDigest
	output.ReverseSignal = jevGenerationTraceCoverageLSPReverseComplete
	output.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjectionReverseObservationDigest(output)
	return output
}

func (value RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjectionReverseObservation) Validate() error {
	if value.Status != "UNKNOWN" && value.Status != "BOUND" {
		return fmt.Errorf("invalid status %q", value.Status)
	}
	if value.MetricName != jevGenerationTraceCoverageLSPReverseMetricName {
		return fmt.Errorf("invalid metric name %q", value.MetricName)
	}
	if value.MissingStageIndex < -1 {
		return fmt.Errorf("invalid missing stage index %d", value.MissingStageIndex)
	}
	if !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("LSP reverse observation must remain non-executing and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" || value.MissingStageIndex != -1 {
			return fmt.Errorf("unknown LSP reverse observation must preserve an unresolved stage")
		}
		if value.CoverageMilli != 0 || value.CoverageBand != "unknown" {
			return fmt.Errorf("unknown LSP reverse observation must not claim coverage")
		}
		if value.ReverseSignal != jevGenerationTraceCoverageLSPReverseUnknown {
			return fmt.Errorf("invalid unknown LSP reverse signal %q", value.ReverseSignal)
		}
	} else {
		if value.MissingStage != "" || value.MissingStageIndex != 0 {
			return fmt.Errorf("bound LSP reverse observation must close the stage")
		}
		if value.CoverageMilli != 1000 || value.CoverageBand != "high" {
			return fmt.Errorf("bound LSP reverse observation must preserve complete coverage")
		}
		if !validJEVDecisionConfidenceClosureDigest(value.MetricDigest) ||
			!validJEVDecisionConfidenceClosureDigest(value.EvidencePrefixDigest) ||
			!validJEVDecisionConfidenceClosureDigest(value.ProjectionObservationDigest) {
			return fmt.Errorf("invalid LSP reverse provenance digest")
		}
		if value.ReverseSignal != jevGenerationTraceCoverageLSPReverseComplete {
			return fmt.Errorf("invalid bound LSP reverse signal %q", value.ReverseSignal)
		}
	}
	expected := revisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjectionReverseObservationDigest(value)
	if value.ObservationDigest != expected {
		return fmt.Errorf("observation digest mismatch")
	}
	return nil
}

func revisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjectionReverseObservationDigest(
	value RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricLSPProjectionReverseObservation,
) string {
	return digestString(fmt.Sprintf(
		"jev-generation-trace-coverage-lsp-reverse|%s|%s|%d|%s|%d|%s|%s|%s|%s|%s|%t|%t",
		value.Status,
		value.MissingStage,
		value.MissingStageIndex,
		value.MetricName,
		value.CoverageMilli,
		value.CoverageBand,
		value.MetricDigest,
		value.EvidencePrefixDigest,
		value.ProjectionObservationDigest,
		value.ReverseSignal,
		value.NonExecuting,
		value.NonAuthorizing,
	))
}
