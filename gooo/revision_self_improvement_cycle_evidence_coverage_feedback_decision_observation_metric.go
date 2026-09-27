package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricObservation struct {
	Status                    string
	MissingStage              string
	DecisionObservationDigest string
	MetricName                string
	LinkedStageCount          int64
	RequiredStageCount        int64
	LinkedDigestCount         int64
	MetricSignal              string
	MetricDigest              string
	ObservationDigest         string
	NonExecuting              bool
	NonAuthorizing            bool
}

func ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(
	decision RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionObservation,
) (RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricObservation, error) {
	result := RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricObservation{
		Status:                    "UNKNOWN",
		MissingStage:              "revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric",
		DecisionObservationDigest: decision.ObservationDigest,
		MetricName:                "evidence-feedback-decision-observation-provenance-link-count",
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	setDigest := func() {
		result.MetricDigest = digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(result)
		result.ObservationDigest = digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricObservation(result)
	}
	setDigest()

	if err := decision.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-input"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision observation is not valid: %w", err)
	}
	if decision.Status != "BOUND" || decision.MissingStage != "" {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-input-status"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision observation is not BOUND")
	}

	result.LinkedStageCount = 2
	result.RequiredStageCount = 2
	result.LinkedDigestCount = 3
	result.MetricSignal = "decision-observation-complete"
	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision observation metric is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("cycle evidence feedback decision observation metric status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound cycle evidence feedback decision observation metric has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown cycle evidence feedback decision observation metric has no missing stage")
	}
	for name, digest := range map[string]string{
		"decision observation": o.DecisionObservationDigest,
		"metric":               o.MetricDigest,
		"observation":          o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("cycle evidence feedback decision observation metric %s digest is invalid", name)
		}
	}
	if o.MetricName != "evidence-feedback-decision-observation-provenance-link-count" {
		return fmt.Errorf("cycle evidence feedback decision observation metric name is invalid")
	}
	if o.LinkedStageCount < 0 || o.RequiredStageCount < 0 || o.LinkedDigestCount < 0 {
		return fmt.Errorf("cycle evidence feedback decision observation metric counts are negative")
	}
	if o.MetricSignal != "decision-observation-complete" &&
		o.MetricSignal != "decision-observation-incomplete" {
		return fmt.Errorf("cycle evidence feedback decision observation metric signal is invalid")
	}
	if o.Status == "BOUND" &&
		(o.LinkedStageCount != 2 ||
			o.RequiredStageCount != 2 ||
			o.LinkedDigestCount != 3 ||
			o.MetricSignal != "decision-observation-complete") {
		return fmt.Errorf("bound cycle evidence feedback decision observation metric is incomplete")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("cycle evidence feedback decision observation metric must remain non-executing and non-authorizing")
	}
	if o.MetricDigest != digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(o) {
		return fmt.Errorf("cycle evidence feedback decision observation metric digest does not match its fields")
	}
	if o.ObservationDigest != digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricObservation(o) {
		return fmt.Errorf("cycle evidence feedback decision observation metric observation digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(
	metric RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricObservation,
) string {
	return digestString(strings.Join([]string{
		metric.Status,
		metric.MissingStage,
		metric.DecisionObservationDigest,
		metric.MetricName,
		strconv.FormatInt(metric.LinkedStageCount, 10),
		strconv.FormatInt(metric.RequiredStageCount, 10),
		strconv.FormatInt(metric.LinkedDigestCount, 10),
		metric.MetricSignal,
		strconv.FormatBool(metric.NonExecuting),
		strconv.FormatBool(metric.NonAuthorizing),
	}, "|"))
}

func digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricObservation(
	metric RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricObservation,
) string {
	return digestString(strings.Join([]string{
		metric.DecisionObservationDigest,
		metric.MetricDigest,
		metric.MetricSignal,
		strconv.FormatInt(metric.LinkedStageCount, 10),
		strconv.FormatInt(metric.RequiredStageCount, 10),
		strconv.FormatInt(metric.LinkedDigestCount, 10),
	}, "|"))
}
