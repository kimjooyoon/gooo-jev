package gooo

import (
	"encoding/hex"
	"fmt"
)

const (
	jevDecisionConfidenceMetricName = "jev-decision-confidence-metric"
	jevDecisionConfidenceReverseClosureUnknown = "jev-confidence-reverse-provenance-closure-unknown"
	jevDecisionConfidenceReverseClosureComplete = "jev-confidence-reverse-provenance-closure-complete"
)

type RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetric struct {
	Status                   string
	MissingStage             string
	MetricName               string
	MetricDigest             string
	ReverseObservationDigest string
	ClosureSignal            string
	ObservationDigest        string
	NonExecuting             bool
	NonAuthorizing           bool
}

func ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetric(
	input RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation,
) RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetric {
	output := RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetric{
		Status:         "UNKNOWN",
		MissingStage:   input.MissingStage,
		MetricName:     jevDecisionConfidenceMetricName,
		ClosureSignal:  jevDecisionConfidenceReverseClosureUnknown,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "reverse_observation"
	}
	if err := input.Validate(); err != nil {
		output.ObservationDigest = revisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetricDigest(output)
		return output
	}
	if input.Status != "BOUND" {
		if input.MissingStage != "" {
			output.MissingStage = input.MissingStage
		}
		output.ObservationDigest = revisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetricDigest(output)
		return output
	}

	output.Status = "BOUND"
	output.MissingStage = ""
	output.MetricDigest = input.MetricDigest
	output.ReverseObservationDigest = input.ObservationDigest
	output.ClosureSignal = jevDecisionConfidenceReverseClosureComplete
	output.ObservationDigest = revisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetricDigest(output)
	return output
}

func (value RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetric) Validate() error {
	if value.Status != "UNKNOWN" && value.Status != "BOUND" {
		return fmt.Errorf("invalid status %q", value.Status)
	}
	if value.MetricName != "" && value.MetricName != jevDecisionConfidenceMetricName {
		return fmt.Errorf("invalid metric name %q", value.MetricName)
	}
	if !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("provenance closure metric must remain non-executing and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" {
			return fmt.Errorf("unknown closure metric must preserve a missing stage")
		}
		if value.ClosureSignal != jevDecisionConfidenceReverseClosureUnknown {
			return fmt.Errorf("invalid unknown closure signal %q", value.ClosureSignal)
		}
	} else {
		if value.MissingStage != "" {
			return fmt.Errorf("bound closure metric cannot have a missing stage")
		}
		if value.MetricName != jevDecisionConfidenceMetricName {
			return fmt.Errorf("bound closure metric must identify the JEV confidence metric")
		}
		if !validJEVDecisionConfidenceClosureDigest(value.MetricDigest) {
			return fmt.Errorf("invalid metric digest")
		}
		if !validJEVDecisionConfidenceClosureDigest(value.ReverseObservationDigest) {
			return fmt.Errorf("invalid reverse observation digest")
		}
		if value.ClosureSignal != jevDecisionConfidenceReverseClosureComplete {
			return fmt.Errorf("invalid bound closure signal %q", value.ClosureSignal)
		}
	}
	expected := revisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetricDigest(value)
	if value.ObservationDigest != expected {
		return fmt.Errorf("observation digest mismatch")
	}
	return nil
}

func revisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetricDigest(
	value RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetric,
) string {
	return digestString(fmt.Sprintf(
		"jev-confidence-reverse-provenance-closure|%s|%s|%s|%s|%s|%s|%t|%t",
		value.Status,
		value.MissingStage,
		value.MetricName,
		value.MetricDigest,
		value.ReverseObservationDigest,
		value.ClosureSignal,
		value.NonExecuting,
		value.NonAuthorizing,
	))
}

func validJEVDecisionConfidenceClosureDigest(value string) bool {
	if len(value) != 64 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}