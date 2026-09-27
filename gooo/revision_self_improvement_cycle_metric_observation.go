package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleMetricInput struct {
	CycleDigest             string
	MetricName              string
	SourceDigest            string
	CandidateSourceDigest   string
	GeneratedIRDigest       string
	ReverseObservationDigest string
	Direction               string
	BaselineValue           int64
	CandidateValue          int64
}

type RevisionSelfImprovementCycleMetricObservation struct {
	Status                   string
	MissingStage             string
	CycleDigest              string
	MetricName               string
	SourceDigest             string
	CandidateSourceDigest    string
	GeneratedIRDigest        string
	ReverseObservationDigest string
	Direction                string
	BaselineValue            int64
	CandidateValue           int64
	MetricDelta              int64
	MetricSignal             string
	MetricEvidenceDigest     string
	ObservationDigest        string
	NonExecuting             bool
	NonAuthorizing           bool
}

func ObserveRevisionSelfImprovementCycleMetric(
	cycle RevisionSelfImprovementCycleObservation,
	input RevisionSelfImprovementCycleMetricInput,
) (RevisionSelfImprovementCycleMetricObservation, error) {
	result := RevisionSelfImprovementCycleMetricObservation{
		Status:                   "UNKNOWN",
		MissingStage:             "revision-self-improvement-cycle-metric",
		CycleDigest:              input.CycleDigest,
		MetricName:               input.MetricName,
		SourceDigest:             input.SourceDigest,
		CandidateSourceDigest:    input.CandidateSourceDigest,
		GeneratedIRDigest:        input.GeneratedIRDigest,
		ReverseObservationDigest: input.ReverseObservationDigest,
		Direction:                input.Direction,
		BaselineValue:            input.BaselineValue,
		CandidateValue:           input.CandidateValue,
		MetricDelta:               input.CandidateValue - input.BaselineValue,
		NonExecuting:             true,
		NonAuthorizing:           true,
	}
	setDigest := func() {
		result.MetricEvidenceDigest = digestRevisionSelfImprovementCycleMetricEvidence(result)
		result.ObservationDigest = digestRevisionSelfImprovementCycleMetric(result)
	}
	setDigest()

	if err := cycle.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-metric-cycle"
		setDigest()
		return result, fmt.Errorf("self-improvement cycle is not valid: %w", err)
	}
	if input.CycleDigest != cycle.ObservationDigest {
		result.MissingStage = "revision-self-improvement-cycle-metric-cycle-link"
		setDigest()
		return result, fmt.Errorf("metric input is not linked to the cycle observation")
	}
	if input.SourceDigest != cycle.SourceDigest ||
		input.CandidateSourceDigest != cycle.NextSourceDigest ||
		input.GeneratedIRDigest != cycle.NextIRDigest ||
		input.ReverseObservationDigest != cycle.ReverseObservationDigest {
		result.MissingStage = "revision-self-improvement-cycle-metric-evidence-link"
		setDigest()
		return result, fmt.Errorf("metric input evidence is not linked to the cycle")
	}
	if input.MetricName == "" {
		result.MissingStage = "revision-self-improvement-cycle-metric-name"
		setDigest()
		return result, fmt.Errorf("metric name is empty")
	}
	if input.Direction != "higher-is-better" && input.Direction != "lower-is-better" {
		result.MissingStage = "revision-self-improvement-cycle-metric-direction"
		setDigest()
		return result, fmt.Errorf("metric direction is invalid")
	}
	if !validDigest(input.CycleDigest) ||
		!validDigest(input.SourceDigest) ||
		!validDigest(input.CandidateSourceDigest) ||
		!validDigest(input.GeneratedIRDigest) ||
		!validDigest(input.ReverseObservationDigest) {
		result.MissingStage = "revision-self-improvement-cycle-metric-evidence"
		setDigest()
		return result, fmt.Errorf("metric evidence digest is invalid")
	}

	switch {
	case input.CandidateValue == input.BaselineValue:
		result.MetricSignal = "metric-stable"
	case input.Direction == "higher-is-better" && input.CandidateValue > input.BaselineValue:
		result.MetricSignal = "metric-improved"
	case input.Direction == "lower-is-better" && input.CandidateValue < input.BaselineValue:
		result.MetricSignal = "metric-improved"
	default:
		result.MetricSignal = "metric-regressed"
	}
	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-metric"
		setDigest()
		return result, fmt.Errorf("self-improvement cycle metric is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementCycleMetricObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("cycle metric status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound cycle metric has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown cycle metric has no missing stage")
	}
	for name, digest := range map[string]string{
		"cycle":              o.CycleDigest,
		"source":             o.SourceDigest,
		"candidate source":   o.CandidateSourceDigest,
		"generated ir":       o.GeneratedIRDigest,
		"reverse observation": o.ReverseObservationDigest,
		"metric evidence":    o.MetricEvidenceDigest,
		"observation":        o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("cycle metric %s digest is invalid", name)
		}
	}
	if o.MetricName == "" {
		return fmt.Errorf("cycle metric name is empty")
	}
	if o.Direction != "higher-is-better" && o.Direction != "lower-is-better" {
		return fmt.Errorf("cycle metric direction is invalid")
	}
	if o.MetricDelta != o.CandidateValue-o.BaselineValue {
		return fmt.Errorf("cycle metric delta does not match values")
	}
	if o.MetricSignal != "metric-improved" &&
		o.MetricSignal != "metric-regressed" &&
		o.MetricSignal != "metric-stable" {
		return fmt.Errorf("cycle metric signal is invalid")
	}
	if o.CandidateValue == o.BaselineValue && o.MetricSignal != "metric-stable" {
		return fmt.Errorf("equal cycle metric values must be stable")
	}
	if o.CandidateValue != o.BaselineValue &&
		o.MetricSignal == "metric-stable" {
		return fmt.Errorf("different cycle metric values cannot be stable")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("cycle metric must remain non-executing and non-authorizing")
	}
	if o.MetricEvidenceDigest != digestRevisionSelfImprovementCycleMetricEvidence(o) {
		return fmt.Errorf("cycle metric evidence digest does not match its fields")
	}
	if o.ObservationDigest != digestRevisionSelfImprovementCycleMetric(o) {
		return fmt.Errorf("cycle metric observation digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleMetricEvidence(
	o RevisionSelfImprovementCycleMetricObservation,
) string {
	return digestString(strings.Join([]string{
		o.CycleDigest,
		o.MetricName,
		o.SourceDigest,
		o.CandidateSourceDigest,
		o.GeneratedIRDigest,
		o.ReverseObservationDigest,
		o.Direction,
		strconv.FormatInt(o.BaselineValue, 10),
		strconv.FormatInt(o.CandidateValue, 10),
	}, "|"))
}

func digestRevisionSelfImprovementCycleMetric(
	o RevisionSelfImprovementCycleMetricObservation,
) string {
	return digestString(strings.Join([]string{
		o.Status,
		o.MissingStage,
		o.CycleDigest,
		o.MetricName,
		o.SourceDigest,
		o.CandidateSourceDigest,
		o.GeneratedIRDigest,
		o.ReverseObservationDigest,
		o.Direction,
		strconv.FormatInt(o.BaselineValue, 10),
		strconv.FormatInt(o.CandidateValue, 10),
		strconv.FormatInt(o.MetricDelta, 10),
		o.MetricSignal,
		o.MetricEvidenceDigest,
		strconv.FormatBool(o.NonExecuting),
		strconv.FormatBool(o.NonAuthorizing),
	}, "|"))
}