package gooo

import "fmt"

const (
	jevExecutionPlanCoverageMetricName = "jev-execution-plan-reverse-observation-coverage"
	jevExecutionPlanCoverageUnknown    = "jev-execution-plan-reverse-observation-coverage-unknown"
	jevExecutionPlanCoverageComplete   = "jev-execution-plan-reverse-observation-coverage-complete"
)

type RevisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetric struct {
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

func ObserveRevisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetric(
	input RevisionSelfImprovementCycleJEVExecutionPlanReverseObservation,
) RevisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetric {
	output := RevisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetric{
		Status:              "UNKNOWN",
		MissingStage:        input.MissingStage,
		MetricName:          jevExecutionPlanCoverageMetricName,
		RequiredDigestCount: 6,
		CoverageBand:        "unknown",
		MetricSignal:        jevExecutionPlanCoverageUnknown,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "execution_plan_reverse_observation"
	}
	if err := input.Validate(); err != nil {
		output.ObservationDigest = revisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricDigest(output)
		return output
	}
	if input.Status != "BOUND" {
		output.ObservationDigest = revisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricDigest(output)
		return output
	}

	output.Status = "BOUND"
	output.MissingStage = ""
	output.LinkedDigestCount = 6
	output.CoverageMilli = 1000
	output.CoverageBand = "high"
	output.MetricSignal = jevExecutionPlanCoverageComplete
	output.MetricDigest = digestString(fmt.Sprintf(
		"jev-execution-plan-reverse-observation-coverage|%s|%s|%s|%s|%s|%s",
		input.TaskDigest,
		input.WorkspaceDigest,
		input.GatewayPolicyDigest,
		input.ModelDigest,
		input.SuspendResumeDigest,
		input.EvidencePrefixDigest,
	))
	output.ReverseObservationDigest = input.ObservationDigest
	output.ObservationDigest = revisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricDigest(output)
	return output
}

func (value RevisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetric) Validate() error {
	if value.Status != "UNKNOWN" && value.Status != "BOUND" {
		return fmt.Errorf("invalid status %q", value.Status)
	}
	if value.MetricName != jevExecutionPlanCoverageMetricName {
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
		if value.MetricSignal != jevExecutionPlanCoverageUnknown {
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
		if value.MetricSignal != jevExecutionPlanCoverageComplete {
			return fmt.Errorf("invalid bound coverage signal %q", value.MetricSignal)
		}
	}
	expected := revisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricDigest(value)
	if value.ObservationDigest != expected {
		return fmt.Errorf("observation digest mismatch")
	}
	return nil
}

func revisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricDigest(
	value RevisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetric,
) string {
	return digestString(fmt.Sprintf(
		"jev-execution-plan-reverse-observation-coverage-metric|%s|%s|%s|%d|%d|%d|%s|%s|%s|%s|%t|%t",
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
