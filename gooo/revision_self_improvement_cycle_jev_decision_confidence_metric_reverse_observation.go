package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation
// records deterministic reverse observation for a JEV confidence metric.
type RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation struct {
	Status                    string
	MissingStage              string
	MetricName                string
	MetricDigest              string
	ReconstructedMetricDigest string
	DecisionObservationDigest string
	ReverseSignal             string
	ObservationDigest         string
	NonExecuting              bool
	NonAuthorizing            bool
}

func ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation(
	metric RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation,
) (RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation, error) {
	result := RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation{
		Status:                    "UNKNOWN",
		MissingStage:              "revision-self-improvement-cycle-jev-decision-confidence-metric-reverse-observation",
		MetricName:                metric.MetricName,
		MetricDigest:              metric.MetricDigest,
		DecisionObservationDigest: metric.DecisionObservationDigest,
		ReverseSignal:             "jev-decision-confidence-metric-reverse-unknown",
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	setDigest := func() {
		result.ObservationDigest =
			digestRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation(result)
	}
	setDigest()

	if err := metric.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-jev-decision-confidence-metric-reverse-observation-input"
		setDigest()
		return result, fmt.Errorf("jev decision confidence metric is not valid for reverse observation: %w", err)
	}
	if metric.Status != "BOUND" || metric.MissingStage != "" {
		result.MissingStage = "revision-self-improvement-cycle-jev-decision-confidence-metric-reverse-observation-input-status"
		setDigest()
		return result, fmt.Errorf("jev decision confidence metric is not BOUND")
	}

	result.ReconstructedMetricDigest =
		digestRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(metric)
	if result.ReconstructedMetricDigest != metric.MetricDigest {
		result.MissingStage = "revision-self-improvement-cycle-jev-decision-confidence-metric-reverse-observation-digest"
		setDigest()
		return result, fmt.Errorf("jev decision confidence metric digest cannot be reconstructed")
	}
	result.Status = "BOUND"
	result.MissingStage = ""
	result.ReverseSignal = "jev-decision-confidence-metric-reverse-complete"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-jev-decision-confidence-metric-reverse-observation"
		result.ReverseSignal = "jev-decision-confidence-metric-reverse-incomplete"
		setDigest()
		return result, fmt.Errorf("jev decision confidence metric reverse observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("jev decision confidence reverse observation status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound jev decision confidence reverse observation has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown jev decision confidence reverse observation has no missing stage")
	}
	if o.MetricName != "jev-typed-decision-confidence-milli" {
		return fmt.Errorf("jev decision confidence reverse observation metric name is invalid")
	}
	for name, digest := range map[string]string{
		"metric":               o.MetricDigest,
		"reconstructed metric": o.ReconstructedMetricDigest,
		"decision observation": o.DecisionObservationDigest,
		"observation":          o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("jev decision confidence reverse observation %s digest is invalid", name)
		}
	}
	if o.ReverseSignal != "jev-decision-confidence-metric-reverse-complete" &&
		o.ReverseSignal != "jev-decision-confidence-metric-reverse-incomplete" &&
		o.ReverseSignal != "jev-decision-confidence-metric-reverse-unknown" {
		return fmt.Errorf("jev decision confidence reverse observation signal is invalid")
	}
	if o.Status == "BOUND" &&
		(o.ReconstructedMetricDigest != o.MetricDigest ||
			o.ReverseSignal != "jev-decision-confidence-metric-reverse-complete") {
		return fmt.Errorf("bound jev decision confidence reverse observation is incomplete")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("jev decision confidence reverse observation must remain non-executing and non-authorizing")
	}
	if o.ObservationDigest !=
		digestRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation(o) {
		return fmt.Errorf("jev decision confidence reverse observation digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation(
	observation RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation,
) string {
	return digestString(strings.Join([]string{
		observation.Status,
		observation.MissingStage,
		observation.MetricName,
		observation.MetricDigest,
		observation.ReconstructedMetricDigest,
		observation.DecisionObservationDigest,
		observation.ReverseSignal,
		strconv.FormatBool(observation.NonExecuting),
		strconv.FormatBool(observation.NonAuthorizing),
	}, "|"))
}
