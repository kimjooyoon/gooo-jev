package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetricObservation struct {
	Status                    string
	MissingStage              string
	MetricObservationDigest   string
	MetricDigest              string
	ReconstructedMetricDigest string
	DecisionObservationDigest string
	ReverseObservationDigest  string
	LinkedDigestCount         int64
	RequiredDigestCount       int64
	MetricSignal              string
	ClosureDigest              string
	ProjectionDigest          string
	ReadOnly                  bool
	NonExecuting              bool
	NonAuthorizing            bool
}

func ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric(
	metric RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric,
) (LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetricObservation, error) {
	result := LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetricObservation{
		Status:                    "UNKNOWN",
		MissingStage:              "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation-provenance-closure-metric",
		MetricObservationDigest:   metric.MetricObservationDigest,
		MetricDigest:              metric.MetricDigest,
		ReconstructedMetricDigest: metric.ReconstructedMetricDigest,
		DecisionObservationDigest: metric.DecisionObservationDigest,
		ReverseObservationDigest:  metric.ReverseObservationDigest,
		LinkedDigestCount:         metric.LinkedDigestCount,
		RequiredDigestCount:       metric.RequiredDigestCount,
		MetricSignal:              metric.MetricSignal,
		ClosureDigest:              metric.ClosureDigest,
		ReadOnly:                  true,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	setDigest := func() {
		result.ProjectionDigest = digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric(result)
	}
	setDigest()

	if err := metric.Validate(); err != nil {
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation-provenance-closure-metric-input"
		setDigest()
		return result, fmt.Errorf("reverse observation provenance closure metric is not valid: %w", err)
	}
	if metric.Status != "BOUND" || metric.MissingStage != "" {
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation-provenance-closure-metric-input-status"
		setDigest()
		return result, fmt.Errorf("reverse observation provenance closure metric is not BOUND")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation-provenance-closure-metric"
		setDigest()
		return result, fmt.Errorf("lsp reverse observation provenance closure metric projection is not valid: %w", err)
	}
	return result, nil
}

func (o LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetricObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("lsp reverse observation provenance closure metric status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound lsp reverse observation provenance closure metric has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown lsp reverse observation provenance closure metric has no missing stage")
	}
	for name, digest := range map[string]string{
		"metric observation":   o.MetricObservationDigest,
		"metric":               o.MetricDigest,
		"reconstructed metric": o.ReconstructedMetricDigest,
		"decision observation": o.DecisionObservationDigest,
		"reverse observation":  o.ReverseObservationDigest,
		"closure":              o.ClosureDigest,
		"projection":           o.ProjectionDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("lsp reverse observation provenance closure metric %s digest is invalid", name)
		}
	}
	if o.LinkedDigestCount < 0 || o.RequiredDigestCount < 0 {
		return fmt.Errorf("lsp reverse observation provenance closure metric counts are negative")
	}
	if o.MetricSignal != "reverse-observation-provenance-complete" &&
		o.MetricSignal != "reverse-observation-provenance-incomplete" {
		return fmt.Errorf("lsp reverse observation provenance closure metric signal is invalid")
	}
	if o.Status == "BOUND" &&
		(o.LinkedDigestCount != 5 ||
			o.RequiredDigestCount != 5 ||
			o.MetricSignal != "reverse-observation-provenance-complete" ||
			o.ReconstructedMetricDigest != o.MetricDigest) {
		return fmt.Errorf("bound lsp reverse observation provenance closure metric is incomplete")
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("lsp reverse observation provenance closure metric must remain read-only, non-executing, and non-authorizing")
	}
	if o.ProjectionDigest != digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric(o) {
		return fmt.Errorf("lsp reverse observation provenance closure metric projection digest does not match its fields")
	}
	return nil
}

func digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric(
	metric LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetricObservation,
) string {
	return digestString(strings.Join([]string{
		metric.Status,
		metric.MissingStage,
		metric.MetricObservationDigest,
		metric.MetricDigest,
		metric.ReconstructedMetricDigest,
		metric.DecisionObservationDigest,
		metric.ReverseObservationDigest,
		strconv.FormatInt(metric.LinkedDigestCount, 10),
		strconv.FormatInt(metric.RequiredDigestCount, 10),
		metric.MetricSignal,
		metric.ClosureDigest,
		strconv.FormatBool(metric.ReadOnly),
		strconv.FormatBool(metric.NonExecuting),
		strconv.FormatBool(metric.NonAuthorizing),
	}, "|"))
}
