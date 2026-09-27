package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation struct {
	Status                           string
	MissingStage                     string
	MetricObservationDigest          string
	MetricName                       string
	MetricDigest                     string
	ReconstructedMetricDigest        string
	DecisionObservationDigest        string
	LinkedStageCount                 int64
	RequiredStageCount               int64
	LinkedDigestCount                int64
	MetricSignal                     string
	ReverseSignal                    string
	ObservationDigest                string
	NonExecuting                     bool
	NonAuthorizing                   bool
}

func ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation(
	metric RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricObservation,
) (RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation, error) {
	result := RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation{
		Status:                  "UNKNOWN",
		MissingStage:            "revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation",
		MetricObservationDigest: metric.ObservationDigest,
		MetricName:              metric.MetricName,
		MetricDigest:            metric.MetricDigest,
		DecisionObservationDigest: metric.DecisionObservationDigest,
		LinkedStageCount:        metric.LinkedStageCount,
		RequiredStageCount:      metric.RequiredStageCount,
		LinkedDigestCount:       metric.LinkedDigestCount,
		MetricSignal:             metric.MetricSignal,
		NonExecuting:            true,
		NonAuthorizing:           true,
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation(result)
	}
	setDigest()

	if err := metric.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation-input"
		setDigest()
		return result, fmt.Errorf("decision observation metric is not valid: %w", err)
	}
	if metric.Status != "BOUND" || metric.MissingStage != "" {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation-input-status"
		setDigest()
		return result, fmt.Errorf("decision observation metric is not BOUND")
	}

	result.ReconstructedMetricDigest = digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(metric)
	if result.ReconstructedMetricDigest != metric.MetricDigest {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation-digest"
		setDigest()
		return result, fmt.Errorf("decision observation metric digest cannot be reconstructed")
	}
	result.ReverseSignal = "decision-observation-metric-reverse-complete"
	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation"
		result.ReverseSignal = "decision-observation-metric-reverse-incomplete"
		setDigest()
		return result, fmt.Errorf("decision observation metric reverse observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("decision observation metric reverse observation status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound decision observation metric reverse observation has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown decision observation metric reverse observation has no missing stage")
	}
	for name, digest := range map[string]string{
		"metric observation":   o.MetricObservationDigest,
		"metric":               o.MetricDigest,
		"reconstructed metric": o.ReconstructedMetricDigest,
		"decision observation": o.DecisionObservationDigest,
		"observation":          o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("decision observation metric reverse observation %s digest is invalid", name)
		}
	}
	if o.MetricName != "evidence-feedback-decision-observation-provenance-link-count" {
		return fmt.Errorf("decision observation metric reverse observation metric name is invalid")
	}
	if o.LinkedStageCount < 0 || o.RequiredStageCount < 0 || o.LinkedDigestCount < 0 {
		return fmt.Errorf("decision observation metric reverse observation counts are negative")
	}
	if o.MetricSignal != "decision-observation-complete" &&
		o.MetricSignal != "decision-observation-incomplete" {
		return fmt.Errorf("decision observation metric reverse observation metric signal is invalid")
	}
	if o.ReverseSignal != "decision-observation-metric-reverse-complete" &&
		o.ReverseSignal != "decision-observation-metric-reverse-incomplete" {
		return fmt.Errorf("decision observation metric reverse signal is invalid")
	}
	if o.Status == "BOUND" &&
		(o.LinkedStageCount != 2 ||
			o.RequiredStageCount != 2 ||
			o.LinkedDigestCount != 3 ||
			o.MetricSignal != "decision-observation-complete" ||
			o.ReverseSignal != "decision-observation-metric-reverse-complete" ||
			o.ReconstructedMetricDigest != o.MetricDigest) {
		return fmt.Errorf("bound decision observation metric reverse observation is incomplete")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("decision observation metric reverse observation must remain non-executing and non-authorizing")
	}
	if o.ObservationDigest != digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation(o) {
		return fmt.Errorf("decision observation metric reverse observation digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation(
	observation RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation,
) string {
	return digestString(strings.Join([]string{
		observation.Status,
		observation.MissingStage,
		observation.MetricObservationDigest,
		observation.MetricName,
		observation.MetricDigest,
		observation.ReconstructedMetricDigest,
		observation.DecisionObservationDigest,
		strconv.FormatInt(observation.LinkedStageCount, 10),
		strconv.FormatInt(observation.RequiredStageCount, 10),
		strconv.FormatInt(observation.LinkedDigestCount, 10),
		observation.MetricSignal,
		observation.ReverseSignal,
		strconv.FormatBool(observation.NonExecuting),
		strconv.FormatBool(observation.NonAuthorizing),
	}, "|"))
}
