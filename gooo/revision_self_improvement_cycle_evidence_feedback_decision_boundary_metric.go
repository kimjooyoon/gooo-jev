package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricObservation struct {
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

func ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric(
	boundary RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryObservation,
) (RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricObservation, error) {
	result := RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricObservation{
		Status:         "UNKNOWN",
		MissingStage:   "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric",
		BoundaryDigest: boundary.ObservationDigest,
		MetricName:     "evidence-feedback-decision-boundary-provenance-link-count",
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	setDigest := func() {
		result.MetricDigest = digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric(result)
		result.ObservationDigest = digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricObservation(result)
	}
	setDigest()

	if err := boundary.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-boundary"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision boundary is not valid: %w", err)
	}
	if boundary.Status != "BOUND" || boundary.MissingStage != "" {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-status"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision boundary is not BOUND")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.LinkedStageCount = 4
	result.RequiredStageCount = 4
	result.LinkedDigestCount = 4
	result.MetricSignal = "boundary-complete"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision boundary metric is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("cycle evidence feedback decision boundary metric status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound cycle evidence feedback decision boundary metric has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown cycle evidence feedback decision boundary metric has no missing stage")
	}
	if !validDigest(o.BoundaryDigest) || !validDigest(o.MetricDigest) || !validDigest(o.ObservationDigest) {
		return fmt.Errorf("cycle evidence feedback decision boundary metric digest is invalid")
	}
	if o.MetricName != "evidence-feedback-decision-boundary-provenance-link-count" {
		return fmt.Errorf("cycle evidence feedback decision boundary metric name is invalid")
	}
	if o.LinkedStageCount < 0 || o.RequiredStageCount < 0 || o.LinkedDigestCount < 0 {
		return fmt.Errorf("cycle evidence feedback decision boundary metric counts are negative")
	}
	if o.MetricSignal != "boundary-complete" && o.MetricSignal != "boundary-incomplete" {
		return fmt.Errorf("cycle evidence feedback decision boundary metric signal is invalid")
	}
	if o.Status == "BOUND" &&
		(o.LinkedStageCount != 4 ||
			o.RequiredStageCount != 4 ||
			o.LinkedDigestCount != 4 ||
			o.MetricSignal != "boundary-complete") {
		return fmt.Errorf("bound cycle evidence feedback decision boundary metric is incomplete")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("cycle evidence feedback decision boundary metric must remain non-executing and non-authorizing")
	}
	if o.MetricDigest != digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric(o) {
		return fmt.Errorf("cycle evidence feedback decision boundary metric digest does not match its fields")
	}
	if o.ObservationDigest != digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricObservation(o) {
		return fmt.Errorf("cycle evidence feedback decision boundary metric observation digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric(
	metric RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricObservation,
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

func digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricObservation(
	metric RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricObservation,
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
