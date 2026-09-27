package gooo

import "fmt"

const (
	jevGenerationTraceCoverageMetricName = "jev-generation-trace-reverse-observation-coverage"
	jevGenerationTraceCoverageUnknown = "jev-generation-trace-reverse-observation-coverage-unknown"
	jevGenerationTraceCoverageComplete = "jev-generation-trace-reverse-observation-coverage-complete"
)

type RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetric struct {
	Status                  string
	MissingStage            string
	MetricName              string
	RequiredDigestCount     int
	LinkedDigestCount       int
	CoverageMilli            int
	CoverageBand             string
	MetricSignal             string
	MetricDigest             string
	ReverseObservationDigest string
	ObservationDigest        string
	NonExecuting             bool
	NonAuthorizing           bool
}

func ObserveRevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetric(
	input RevisionSelfImprovementCycleJEVGenerationTraceReverseObservation,
) RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetric {
	output := RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetric{
		Status:              "UNKNOWN",
		MissingStage:        input.MissingStage,
		MetricName:          jevGenerationTraceCoverageMetricName,
		RequiredDigestCount: 6,
		CoverageBand:        "unknown",
		MetricSignal:        jevGenerationTraceCoverageUnknown,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "generation_trace_reverse_observation"
	}
	if err := input.Validate(); err != nil {
		output.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricDigest(output)
		return output
	}
	if input.Status != "BOUND" {
		output.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricDigest(output)
		return output
	}

	output.Status = "BOUND"
	output.MissingStage = ""
	output.LinkedDigestCount = 6
	output.CoverageMilli = 1000
	output.CoverageBand = "high"
	output.MetricSignal = jevGenerationTraceCoverageComplete
	output.MetricDigest = digestString(fmt.Sprintf(
		"jev-generation-trace-reverse-observation-coverage|%s|%s|%s|%s|%s|%s",
		input.ContractDigest,
		input.IRDigest,
		input.GeneratedArtifactDigest,
		input.MetricDigest,
		input.EvidencePrefixDigest,
		input.ReconstructedTraceDigest,
	))
	output.ReverseObservationDigest = input.ObservationDigest
	output.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricDigest(output)
	return output
}

func (value RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetric) Validate() error {
	if value.Status != "UNKNOWN" && value.Status != "BOUND" {
		return fmt.Errorf("invalid status %q", value.Status)
	}
	if value.MetricName != jevGenerationTraceCoverageMetricName {
		return fmt.Errorf("invalid metric name %q", value.MetricName)
	}
	if value.RequiredDigestCount != 6 || value.LinkedDigestCount < 0 || value.LinkedDigestCount > value.RequiredDigestCount {
		return fmt.Errorf("invalid digest coverage counts")
	}
	if !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("coverage metric must remain non-executing and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" {
			return fmt.Errorf("unknown coverage metric must preserve a missing stage")
		}
		if value.CoverageMilli != 0 || value.CoverageBand != "unknown" {
			return fmt.Errorf("unknown coverage metric must not claim coverage")
		}
		if value.MetricSignal != jevGenerationTraceCoverageUnknown {
			return fmt.Errorf("invalid unknown coverage signal %q", value.MetricSignal)
		}
	} else {
		if value.MissingStage != "" {
			return fmt.Errorf("bound coverage metric cannot have a missing stage")
		}
		if value.LinkedDigestCount != value.RequiredDigestCount || value.CoverageMilli != 1000 || value.CoverageBand != "high" {
			return fmt.Errorf("bound coverage metric must report complete coverage")
		}
		if !validJEVDecisionConfidenceClosureDigest(value.MetricDigest) ||
			!validJEVDecisionConfidenceClosureDigest(value.ReverseObservationDigest) {
			return fmt.Errorf("invalid coverage provenance digest")
		}
		if value.MetricSignal != jevGenerationTraceCoverageComplete {
			return fmt.Errorf("invalid bound coverage signal %q", value.MetricSignal)
		}
	}
	expected := revisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricDigest(value)
	if value.ObservationDigest != expected {
		return fmt.Errorf("observation digest mismatch")
	}
	return nil
}

func revisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricDigest(
	value RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetric,
) string {
	return digestString(fmt.Sprintf(
		"jev-generation-trace-reverse-observation-coverage-metric|%s|%s|%s|%d|%d|%d|%s|%s|%s|%s|%t|%t",
		value.Status,
		value.MissingStage,
		value.MetricName,
		value.RequiredDigestCount,
		value.LinkedDigestCount,
		value.CoverageMilli,
		value.CoverageBand,
		value.MetricSignal,
		value.MetricDigest,
		value.ReverseObservationDigest,
		value.NonExecuting,
		value.NonAuthorizing,
	))
}
