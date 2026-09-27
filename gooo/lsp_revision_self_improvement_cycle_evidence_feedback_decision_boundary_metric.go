package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric struct {
	Status             string
	MissingStage       string
	BoundaryDigest     string
	MetricName         string
	LinkedStageCount   int64
	RequiredStageCount int64
	LinkedDigestCount  int64
	MetricSignal       string
	MetricDigest       string
	ProjectionDigest   string
	ReadOnly           bool
	NonExecuting       bool
	NonAuthorizing     bool
}

func ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric(
	metric RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricObservation,
) (LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric, error) {
	result := LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric{
		Status:             "UNKNOWN",
		MissingStage:       "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric",
		BoundaryDigest:     metric.BoundaryDigest,
		MetricName:         metric.MetricName,
		LinkedStageCount:   metric.LinkedStageCount,
		RequiredStageCount: metric.RequiredStageCount,
		LinkedDigestCount:  metric.LinkedDigestCount,
		MetricSignal:       metric.MetricSignal,
		MetricDigest:       metric.MetricDigest,
		ReadOnly:           true,
		NonExecuting:       true,
		NonAuthorizing:     true,
	}
	setDigest := func() {
		result.ProjectionDigest = digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric(result)
	}
	setDigest()

	if err := metric.Validate(); err != nil {
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-input"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision boundary metric is not valid: %w", err)
	}
	if metric.Status != "BOUND" || metric.MissingStage != "" {
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-input-status"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision boundary metric is not BOUND")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric"
		setDigest()
		return result, fmt.Errorf("lsp evidence feedback decision boundary metric projection is not valid: %w", err)
	}
	return result, nil
}

func (o LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("lsp evidence feedback decision boundary metric status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound lsp evidence feedback decision boundary metric has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown lsp evidence feedback decision boundary metric has no missing stage")
	}
	for name, digest := range map[string]string{
		"boundary":   o.BoundaryDigest,
		"metric":     o.MetricDigest,
		"projection": o.ProjectionDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("lsp evidence feedback decision boundary metric %s digest is invalid", name)
		}
	}
	if o.MetricName != "evidence-feedback-decision-boundary-provenance-link-count" {
		return fmt.Errorf("lsp evidence feedback decision boundary metric name is invalid")
	}
	if o.LinkedStageCount < 0 || o.RequiredStageCount < 0 || o.LinkedDigestCount < 0 {
		return fmt.Errorf("lsp evidence feedback decision boundary metric counts are negative")
	}
	if o.MetricSignal != "boundary-complete" && o.MetricSignal != "boundary-incomplete" {
		return fmt.Errorf("lsp evidence feedback decision boundary metric signal is invalid")
	}
	if o.Status == "BOUND" &&
		(o.LinkedStageCount != 4 ||
			o.RequiredStageCount != 4 ||
			o.LinkedDigestCount != 4 ||
			o.MetricSignal != "boundary-complete") {
		return fmt.Errorf("bound lsp evidence feedback decision boundary metric is incomplete")
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("lsp evidence feedback decision boundary metric must remain read-only and non-executing")
	}
	if o.ProjectionDigest != digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric(o) {
		return fmt.Errorf("lsp evidence feedback decision boundary metric projection digest does not match its fields")
	}
	return nil
}

func digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric(
	metric LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric,
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
		strconv.FormatBool(metric.ReadOnly),
		strconv.FormatBool(metric.NonExecuting),
		strconv.FormatBool(metric.NonAuthorizing),
	}, "|"))
}
