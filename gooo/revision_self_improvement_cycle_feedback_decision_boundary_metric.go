package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleFeedbackDecisionBoundaryMetricObservation struct {
	Status             string
	MissingStage       string
	BoundaryDigest     string
	MetricName         string
	LinkedStageCount   int64
	RequiredStageCount int64
	LinkedDigestCount  int64
	MetricSignal       string
	MetricDigest       string
	ObservationDigest  string
	NonExecuting       bool
	NonAuthorizing     bool
}

func ObserveRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetric(
	boundary RevisionSelfImprovementCycleFeedbackDecisionBoundaryObservation,
) (RevisionSelfImprovementCycleFeedbackDecisionBoundaryMetricObservation, error) {
	result := RevisionSelfImprovementCycleFeedbackDecisionBoundaryMetricObservation{
		Status:            "UNKNOWN",
		MissingStage:      "revision-self-improvement-cycle-feedback-decision-boundary-metric",
		BoundaryDigest:    boundary.ObservationDigest,
		MetricName:        "provenance-link-count",
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	setDigest := func() {
		result.MetricDigest = digestRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetric(result)
		result.ObservationDigest = digestRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetricObservation(result)
	}
	setDigest()

	if err := boundary.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-feedback-decision-boundary-metric-boundary"
		setDigest()
		return result, fmt.Errorf("cycle feedback decision boundary is not valid: %w", err)
	}
	if boundary.Status != "BOUND" || boundary.MissingStage != "" {
		result.MissingStage = "revision-self-improvement-cycle-feedback-decision-boundary-metric-status"
		setDigest()
		return result, fmt.Errorf("cycle feedback decision boundary is not BOUND")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.LinkedStageCount = 3
	result.RequiredStageCount = 3
	result.LinkedDigestCount = 3
	result.MetricSignal = "boundary-complete"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-feedback-decision-boundary-metric"
		setDigest()
		return result, fmt.Errorf("cycle feedback decision boundary metric is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementCycleFeedbackDecisionBoundaryMetricObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("cycle feedback decision boundary metric status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound cycle feedback decision boundary metric has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown cycle feedback decision boundary metric has no missing stage")
	}
	if !validDigest(o.BoundaryDigest) || !validDigest(o.MetricDigest) || !validDigest(o.ObservationDigest) {
		return fmt.Errorf("cycle feedback decision boundary metric digest is invalid")
	}
	if o.MetricName != "provenance-link-count" {
		return fmt.Errorf("cycle feedback decision boundary metric name is invalid")
	}
	if o.LinkedStageCount < 0 || o.RequiredStageCount < 0 || o.LinkedDigestCount < 0 {
		return fmt.Errorf("cycle feedback decision boundary metric counts are negative")
	}
	if o.MetricSignal != "boundary-complete" && o.MetricSignal != "boundary-incomplete" {
		return fmt.Errorf("cycle feedback decision boundary metric signal is invalid")
	}
	if o.Status == "BOUND" && (o.LinkedStageCount != o.RequiredStageCount || o.LinkedStageCount != 3 || o.LinkedDigestCount != 3 || o.MetricSignal != "boundary-complete") {
		return fmt.Errorf("bound cycle feedback decision boundary metric is incomplete")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("cycle feedback decision boundary metric must remain non-executing and non-authorizing")
	}
	if o.MetricDigest != digestRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetric(o) {
		return fmt.Errorf("cycle feedback decision boundary metric digest does not match its fields")
	}
	if o.ObservationDigest != digestRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetricObservation(o) {
		return fmt.Errorf("cycle feedback decision boundary metric observation digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetric(
	metric RevisionSelfImprovementCycleFeedbackDecisionBoundaryMetricObservation,
) string {
	return digestString(strings.Join([]string{
		metric.BoundaryDigest,
		metric.MetricName,
		strconv.FormatInt(metric.LinkedStageCount, 10),
		strconv.FormatInt(metric.RequiredStageCount, 10),
		strconv.FormatInt(metric.LinkedDigestCount, 10),
		metric.MetricSignal,
	}, "|"))
}

func digestRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetricObservation(
	metric RevisionSelfImprovementCycleFeedbackDecisionBoundaryMetricObservation,
) string {
	return digestString(strings.Join([]string{
		metric.Status,
		metric.MissingStage,
		metric.BoundaryDigest,
		metric.MetricName,
		strconv.FormatInt(metric.LinkedStageCount, 10),
		strconv.FormatInt(metric.RequiredStageCount, 10),
		strconv.FormatInt(metric.LinkedDigestCount, 10),
		metric.MetricSignal,
		metric.MetricDigest,
		strconv.FormatBool(metric.NonExecuting),
		strconv.FormatBool(metric.NonAuthorizing),
	}, "|"))
}
