package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// LSPRevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation carries
// a JEV confidence metric into an editor-facing, read-only projection.
type LSPRevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation struct {
	Status                    string
	MissingStage              string
	MetricName                string
	QuestionID                string
	QuestionKind              string
	ConfidenceMilli           int64
	ConfidenceBand            string
	DecisionDigest            string
	EvidenceDigest            string
	DecisionObservationDigest string
	LinkedDigestCount         int64
	RequiredDigestCount       int64
	MetricSignal              string
	MetricDigest              string
	ProjectionSignal          string
	ProjectionDigest          string
	ReadOnly                  bool
	NonExecuting              bool
	NonAuthorizing            bool
}

func ObserveLSPRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(
	metric RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation,
) (LSPRevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation, error) {
	result := LSPRevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation{
		Status:                    "UNKNOWN",
		MissingStage:              "lsp-revision-self-improvement-cycle-jev-decision-confidence-metric",
		MetricName:                metric.MetricName,
		QuestionID:                metric.QuestionID,
		QuestionKind:              metric.QuestionKind,
		ConfidenceMilli:           metric.ConfidenceMilli,
		ConfidenceBand:            metric.ConfidenceBand,
		DecisionDigest:             metric.DecisionDigest,
		EvidenceDigest:             metric.EvidenceDigest,
		DecisionObservationDigest: metric.DecisionObservationDigest,
		LinkedDigestCount:          metric.LinkedDigestCount,
		RequiredDigestCount:       metric.RequiredDigestCount,
		MetricSignal:               metric.MetricSignal,
		MetricDigest:               metric.MetricDigest,
		ProjectionSignal:           "jev-decision-confidence-lsp-unknown",
		ReadOnly:                  true,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	setDigest := func() {
		result.ProjectionDigest = digestLSPRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(result)
	}
	setDigest()

	if err := metric.Validate(); err != nil {
		result.MissingStage = "lsp-revision-self-improvement-cycle-jev-decision-confidence-metric-input"
		setDigest()
		return result, fmt.Errorf("jev decision confidence metric is not valid for lsp projection: %w", err)
	}
	if metric.Status != "BOUND" || metric.MissingStage != "" {
		result.MissingStage = "lsp-revision-self-improvement-cycle-jev-decision-confidence-metric-input-status"
		setDigest()
		return result, fmt.Errorf("jev decision confidence metric is not BOUND")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.ProjectionSignal = "jev-decision-confidence-lsp-projected"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "lsp-revision-self-improvement-cycle-jev-decision-confidence-metric"
		result.ProjectionSignal = "jev-decision-confidence-lsp-unknown"
		setDigest()
		return result, fmt.Errorf("jev decision confidence lsp projection is not valid: %w", err)
	}
	return result, nil
}

func (o LSPRevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("jev decision confidence lsp status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound jev decision confidence lsp projection has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown jev decision confidence lsp projection has no missing stage")
	}
	source := RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation{
		Status:                    o.Status,
		MissingStage:              o.MissingStage,
		MetricName:                o.MetricName,
		QuestionID:                o.QuestionID,
		QuestionKind:              o.QuestionKind,
		ConfidenceMilli:           o.ConfidenceMilli,
		ConfidenceBand:            o.ConfidenceBand,
		DecisionDigest:            o.DecisionDigest,
		EvidenceDigest:            o.EvidenceDigest,
		DecisionObservationDigest: o.DecisionObservationDigest,
		LinkedDigestCount:         o.LinkedDigestCount,
		RequiredDigestCount:       o.RequiredDigestCount,
		MetricSignal:               o.MetricSignal,
		MetricDigest:               o.MetricDigest,
		ReadOnly:                  o.ReadOnly,
		NonExecuting:              o.NonExecuting,
		NonAuthorizing:            o.NonAuthorizing,
	}
	if err := source.Validate(); err != nil {
		return fmt.Errorf("jev decision confidence lsp source metric is invalid: %w", err)
	}
	if o.ProjectionSignal != "jev-decision-confidence-lsp-projected" &&
		o.ProjectionSignal != "jev-decision-confidence-lsp-unknown" {
		return fmt.Errorf("jev decision confidence lsp projection signal is invalid")
	}
	if o.Status == "BOUND" && o.ProjectionSignal != "jev-decision-confidence-lsp-projected" {
		return fmt.Errorf("bound jev decision confidence lsp projection has an unknown signal")
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("jev decision confidence lsp projection must remain read-only, non-executing, and non-authorizing")
	}
	if o.ProjectionDigest != digestLSPRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(o) {
		return fmt.Errorf("jev decision confidence lsp projection digest does not match its fields")
	}
	return nil
}

func digestLSPRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(
	observation LSPRevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation,
) string {
	return digestString(strings.Join([]string{
		observation.Status,
		observation.MissingStage,
		observation.MetricName,
		observation.QuestionID,
		observation.QuestionKind,
		strconv.FormatInt(observation.ConfidenceMilli, 10),
		observation.ConfidenceBand,
		observation.DecisionDigest,
		observation.EvidenceDigest,
		observation.DecisionObservationDigest,
		strconv.FormatInt(observation.LinkedDigestCount, 10),
		strconv.FormatInt(observation.RequiredDigestCount, 10),
		observation.MetricSignal,
		observation.MetricDigest,
		observation.ProjectionSignal,
		strconv.FormatBool(observation.ReadOnly),
		strconv.FormatBool(observation.NonExecuting),
		strconv.FormatBool(observation.NonAuthorizing),
	}, "|"))
}
