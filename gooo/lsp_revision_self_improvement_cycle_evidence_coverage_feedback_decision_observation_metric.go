package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricObservation struct {
	Status                    string
	MissingStage              string
	DecisionObservationDigest string
	MetricName                string
	LinkedStageCount          int64
	RequiredStageCount        int64
	LinkedDigestCount         int64
	MetricSignal              string
	MetricDigest              string
	ProjectionDigest          string
	ReadOnly                  bool
	NonExecuting              bool
	NonAuthorizing            bool
}

func ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(
	metric RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricObservation,
) (LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricObservation, error) {
	result := LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricObservation{
		Status:                    "UNKNOWN",
		MissingStage:              "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric",
		DecisionObservationDigest: metric.DecisionObservationDigest,
		MetricName:                metric.MetricName,
		LinkedStageCount:          metric.LinkedStageCount,
		RequiredStageCount:        metric.RequiredStageCount,
		LinkedDigestCount:         metric.LinkedDigestCount,
		MetricSignal:              metric.MetricSignal,
		MetricDigest:              metric.MetricDigest,
		ReadOnly:                  true,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	setDigest := func() {
		result.ProjectionDigest = digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(result)
	}
	setDigest()

	if err := metric.Validate(); err != nil {
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-input"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision observation metric is not valid: %w", err)
	}
	if metric.Status != "BOUND" || metric.MissingStage != "" {
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-input-status"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision observation metric is not BOUND")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric"
		setDigest()
		return result, fmt.Errorf("lsp cycle evidence feedback decision observation metric projection is not valid: %w", err)
	}
	return result, nil
}

func (o LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("lsp cycle evidence feedback decision observation metric status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound lsp cycle evidence feedback decision observation metric has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown lsp cycle evidence feedback decision observation metric has no missing stage")
	}
	for name, digest := range map[string]string{
		"decision observation": o.DecisionObservationDigest,
		"metric":               o.MetricDigest,
		"projection":           o.ProjectionDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("lsp cycle evidence feedback decision observation metric %s digest is invalid", name)
		}
	}
	if o.MetricName != "evidence-feedback-decision-observation-provenance-link-count" {
		return fmt.Errorf("lsp cycle evidence feedback decision observation metric name is invalid")
	}
	if o.LinkedStageCount < 0 || o.RequiredStageCount < 0 || o.LinkedDigestCount < 0 {
		return fmt.Errorf("lsp cycle evidence feedback decision observation metric counts are negative")
	}
	if o.MetricSignal != "decision-observation-complete" &&
		o.MetricSignal != "decision-observation-incomplete" {
		return fmt.Errorf("lsp cycle evidence feedback decision observation metric signal is invalid")
	}
	if o.Status == "BOUND" &&
		(o.LinkedStageCount != 2 ||
			o.RequiredStageCount != 2 ||
			o.LinkedDigestCount != 3 ||
			o.MetricSignal != "decision-observation-complete") {
		return fmt.Errorf("bound lsp cycle evidence feedback decision observation metric is incomplete")
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("lsp cycle evidence feedback decision observation metric must remain read-only, non-executing, and non-authorizing")
	}
	if o.ProjectionDigest != digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(o) {
		return fmt.Errorf("lsp cycle evidence feedback decision observation metric projection digest does not match its fields")
	}
	return nil
}

func digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(
	metric LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricObservation,
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
		metric.MetricDigest,
		strconv.FormatBool(metric.ReadOnly),
		strconv.FormatBool(metric.NonExecuting),
		strconv.FormatBool(metric.NonAuthorizing),
	}, "|"))
}
