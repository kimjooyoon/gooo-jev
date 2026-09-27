package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric struct {
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
	NonExecuting              bool
	NonAuthorizing            bool
}

func ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric(
	observation RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation,
) (RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric, error) {
	result := RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric{
		Status:                    "UNKNOWN",
		MissingStage:              "revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation-provenance-closure-metric",
		MetricObservationDigest:   observation.MetricObservationDigest,
		MetricDigest:              observation.MetricDigest,
		ReconstructedMetricDigest: observation.ReconstructedMetricDigest,
		DecisionObservationDigest: observation.DecisionObservationDigest,
		ReverseObservationDigest:  observation.ObservationDigest,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	setDigest := func() {
		result.ClosureDigest = digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric(result)
	}
	setDigest()

	if err := observation.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation-provenance-closure-metric-input"
		setDigest()
		return result, fmt.Errorf("decision observation metric reverse observation is not valid: %w", err)
	}
	if observation.Status != "BOUND" || observation.MissingStage != "" {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation-provenance-closure-metric-input-status"
		setDigest()
		return result, fmt.Errorf("decision observation metric reverse observation is not BOUND")
	}

	result.LinkedDigestCount = 5
	result.RequiredDigestCount = 5
	result.MetricSignal = "reverse-observation-provenance-complete"
	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation-provenance-closure-metric"
		result.MetricSignal = "reverse-observation-provenance-incomplete"
		setDigest()
		return result, fmt.Errorf("decision observation metric reverse observation provenance closure metric is not valid: %w", err)
	}
	return result, nil
}

func (m RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric) Validate() error {
	if m.Status != "BOUND" && m.Status != "UNKNOWN" {
		return fmt.Errorf("reverse observation provenance closure metric status is invalid")
	}
	if m.Status == "BOUND" && m.MissingStage != "" {
		return fmt.Errorf("bound reverse observation provenance closure metric has a missing stage")
	}
	if m.Status == "UNKNOWN" && m.MissingStage == "" {
		return fmt.Errorf("unknown reverse observation provenance closure metric has no missing stage")
	}
	for name, digest := range map[string]string{
		"metric observation":   m.MetricObservationDigest,
		"metric":               m.MetricDigest,
		"reconstructed metric": m.ReconstructedMetricDigest,
		"decision observation": m.DecisionObservationDigest,
		"reverse observation":  m.ReverseObservationDigest,
		"closure":              m.ClosureDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("reverse observation provenance closure metric %s digest is invalid", name)
		}
	}
	if m.LinkedDigestCount < 0 || m.RequiredDigestCount < 0 {
		return fmt.Errorf("reverse observation provenance closure metric counts are negative")
	}
	if m.MetricSignal != "reverse-observation-provenance-complete" &&
		m.MetricSignal != "reverse-observation-provenance-incomplete" {
		return fmt.Errorf("reverse observation provenance closure metric signal is invalid")
	}
	if m.Status == "BOUND" &&
		(m.LinkedDigestCount != 5 ||
			m.RequiredDigestCount != 5 ||
			m.MetricSignal != "reverse-observation-provenance-complete" ||
			m.ReconstructedMetricDigest != m.MetricDigest) {
		return fmt.Errorf("bound reverse observation provenance closure metric is incomplete")
	}
	if !m.NonExecuting || !m.NonAuthorizing {
		return fmt.Errorf("reverse observation provenance closure metric must remain non-executing and non-authorizing")
	}
	if m.ClosureDigest != digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric(m) {
		return fmt.Errorf("reverse observation provenance closure metric digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric(
	metric RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric,
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
		strconv.FormatBool(metric.NonExecuting),
		strconv.FormatBool(metric.NonAuthorizing),
	}, "|"))
}
