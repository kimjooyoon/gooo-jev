package gooo

import "fmt"

const (
	jevGenerationTraceDecisionConfidenceClosureMetricName = "jev-generation-trace-decision-confidence-closure"
	jevGenerationTraceDecisionConfidenceClosureUnknown = "jev-generation-trace-decision-confidence-closure-unknown"
	jevGenerationTraceDecisionConfidenceClosureComplete = "jev-generation-trace-decision-confidence-closure-complete"
)

type RevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetricInput struct {
	GenerationTraceMetric   RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetric
	ConfidenceClosureMetric RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetric
}

type RevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetric struct {
	Status                                  string
	MissingStage                            string
	MetricName                              string
	RequiredMetricCount                     int
	LinkedMetricCount                       int
	GenerationTraceMetricDigest             string
	GenerationTraceReverseObservationDigest string
	ConfidenceClosureMetricDigest            string
	ConfidenceReverseObservationDigest       string
	ClosureSignal                           string
	ObservationDigest                       string
	NonExecuting                            bool
	NonAuthorizing                           bool
}

// ObserveRevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetric
// links generation evidence coverage with decision-confidence reverse provenance
// without turning either metric into an improvement or authorization claim.
func ObserveRevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetric(
	input RevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetricInput,
) RevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetric {
	output := RevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetric{
		Status:              "UNKNOWN",
		MissingStage:        "generation_trace_coverage_metric",
		MetricName:          jevGenerationTraceDecisionConfidenceClosureMetricName,
		RequiredMetricCount:  2,
		ClosureSignal:       jevGenerationTraceDecisionConfidenceClosureUnknown,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}

	if !input.GenerationTraceMetric.NonExecuting ||
		!input.GenerationTraceMetric.NonAuthorizing ||
		!input.ConfidenceClosureMetric.NonExecuting ||
		!input.ConfidenceClosureMetric.NonAuthorizing {
		output.MissingStage = "capability_boundary"
		output.NonAuthorizing = false
		output.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetricDigest(output)
		return output
	}
	if err := input.GenerationTraceMetric.Validate(); err != nil {
		output.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetricDigest(output)
		return output
	}
	if err := input.ConfidenceClosureMetric.Validate(); err != nil {
		output.MissingStage = "decision_confidence_closure_metric"
		output.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetricDigest(output)
		return output
	}

	output.GenerationTraceMetricDigest = input.GenerationTraceMetric.MetricDigest
	output.GenerationTraceReverseObservationDigest = input.GenerationTraceMetric.ReverseObservationDigest
	output.ConfidenceClosureMetricDigest = input.ConfidenceClosureMetric.MetricDigest
	output.ConfidenceReverseObservationDigest = input.ConfidenceClosureMetric.ReverseObservationDigest

	if input.GenerationTraceMetric.Status != "BOUND" {
		output.MissingStage = input.GenerationTraceMetric.MissingStage
		if output.MissingStage == "" {
			output.MissingStage = "generation_trace_coverage_metric"
		}
		output.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetricDigest(output)
		return output
	}
	if input.ConfidenceClosureMetric.Status != "BOUND" {
		output.MissingStage = input.ConfidenceClosureMetric.MissingStage
		if output.MissingStage == "" {
			output.MissingStage = "decision_confidence_closure_metric"
		}
		output.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetricDigest(output)
		return output
	}

	output.Status = "BOUND"
	output.MissingStage = ""
	output.LinkedMetricCount = output.RequiredMetricCount
	output.ClosureSignal = jevGenerationTraceDecisionConfidenceClosureComplete
	output.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetricDigest(output)
	return output
}

func (value RevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetric) Validate() error {
	if value.Status != "UNKNOWN" && value.Status != "BOUND" {
		return fmt.Errorf("invalid generation trace decision confidence closure status %q", value.Status)
	}
	if value.MetricName != jevGenerationTraceDecisionConfidenceClosureMetricName {
		return fmt.Errorf("invalid generation trace decision confidence closure metric name %q", value.MetricName)
	}
	if value.RequiredMetricCount != 2 || value.LinkedMetricCount < 0 || value.LinkedMetricCount > value.RequiredMetricCount {
		return fmt.Errorf("invalid generation trace decision confidence metric counts")
	}
	if !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("generation trace decision confidence closure must remain non-executing and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" ||
			value.LinkedMetricCount != 0 ||
			value.ClosureSignal != jevGenerationTraceDecisionConfidenceClosureUnknown {
			return fmt.Errorf("unknown generation trace decision confidence closure is incomplete")
		}
	} else {
		if value.MissingStage != "" ||
			value.LinkedMetricCount != value.RequiredMetricCount ||
			value.ClosureSignal != jevGenerationTraceDecisionConfidenceClosureComplete {
			return fmt.Errorf("bound generation trace decision confidence closure is incomplete")
		}
		for name, digest := range map[string]string{
			"generation trace metric":             value.GenerationTraceMetricDigest,
			"generation trace reverse observation": value.GenerationTraceReverseObservationDigest,
			"confidence closure metric":           value.ConfidenceClosureMetricDigest,
			"confidence reverse observation":      value.ConfidenceReverseObservationDigest,
		} {
			if !validJEVDecisionConfidenceClosureDigest(digest) {
				return fmt.Errorf("invalid %s digest", name)
			}
		}
	}
	expected := revisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetricDigest(value)
	if value.ObservationDigest != expected {
		return fmt.Errorf("generation trace decision confidence closure digest mismatch")
	}
	return nil
}

func revisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetricDigest(
	value RevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetric,
) string {
	return digestString(fmt.Sprintf(
		"jev-generation-trace-decision-confidence-closure|%s|%s|%s|%d|%d|%s|%s|%s|%s|%s|%t|%t",
		value.Status,
		value.MissingStage,
		value.MetricName,
		value.RequiredMetricCount,
		value.LinkedMetricCount,
		value.GenerationTraceMetricDigest,
		value.GenerationTraceReverseObservationDigest,
		value.ConfidenceClosureMetricDigest,
		value.ConfidenceReverseObservationDigest,
		value.ClosureSignal,
		value.NonExecuting,
		value.NonAuthorizing,
	))
}