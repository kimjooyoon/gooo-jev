package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation struct {
	Status                    string
	MissingStage              string
	MetricObservationDigest   string
	MetricName                string
	MetricDigest              string
	ReconstructedMetricDigest string
	DecisionObservationDigest string
	LinkedStageCount          int64
	RequiredStageCount        int64
	LinkedDigestCount         int64
	MetricSignal              string
	ReverseSignal             string
	ObservationDigest         string
	ProjectionDigest          string
	ReadOnly                  bool
	NonExecuting              bool
	NonAuthorizing            bool
}

func ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation(
	observation RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation,
) (LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation, error) {
	result := LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation{
		Status:                    "UNKNOWN",
		MissingStage:              "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation",
		MetricObservationDigest:   observation.MetricObservationDigest,
		MetricName:                observation.MetricName,
		MetricDigest:              observation.MetricDigest,
		ReconstructedMetricDigest: observation.ReconstructedMetricDigest,
		DecisionObservationDigest: observation.DecisionObservationDigest,
		LinkedStageCount:          observation.LinkedStageCount,
		RequiredStageCount:        observation.RequiredStageCount,
		LinkedDigestCount:         observation.LinkedDigestCount,
		MetricSignal:              observation.MetricSignal,
		ReverseSignal:             observation.ReverseSignal,
		ObservationDigest:         observation.ObservationDigest,
		ReadOnly:                  true,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	setDigest := func() {
		result.ProjectionDigest = digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation(result)
	}
	setDigest()

	if err := observation.Validate(); err != nil {
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation-input"
		setDigest()
		return result, fmt.Errorf("decision observation metric reverse observation is not valid: %w", err)
	}
	if observation.Status != "BOUND" || observation.MissingStage != "" {
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation-input-status"
		setDigest()
		return result, fmt.Errorf("decision observation metric reverse observation is not BOUND")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation"
		setDigest()
		return result, fmt.Errorf("lsp decision observation metric reverse observation projection is not valid: %w", err)
	}
	return result, nil
}

func (o LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("lsp decision observation metric reverse observation status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound lsp decision observation metric reverse observation has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown lsp decision observation metric reverse observation has no missing stage")
	}
	for name, digest := range map[string]string{
		"metric observation":   o.MetricObservationDigest,
		"metric":              o.MetricDigest,
		"reconstructed metric": o.ReconstructedMetricDigest,
		"decision observation": o.DecisionObservationDigest,
		"observation":         o.ObservationDigest,
		"projection":          o.ProjectionDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("lsp decision observation metric reverse observation %s digest is invalid", name)
		}
	}
	if o.MetricName != "evidence-feedback-decision-observation-provenance-link-count" {
		return fmt.Errorf("lsp decision observation metric reverse observation metric name is invalid")
	}
	if o.LinkedStageCount < 0 || o.RequiredStageCount < 0 || o.LinkedDigestCount < 0 {
		return fmt.Errorf("lsp decision observation metric reverse observation counts are negative")
	}
	if o.MetricSignal != "decision-observation-complete" &&
		o.MetricSignal != "decision-observation-incomplete" {
		return fmt.Errorf("lsp decision observation metric reverse observation metric signal is invalid")
	}
	if o.ReverseSignal != "decision-observation-metric-reverse-complete" &&
		o.ReverseSignal != "decision-observation-metric-reverse-incomplete" {
		return fmt.Errorf("lsp decision observation metric reverse observation reverse signal is invalid")
	}
	if o.Status == "BOUND" &&
		(o.LinkedStageCount != 2 ||
			o.RequiredStageCount != 2 ||
			o.LinkedDigestCount != 3 ||
			o.MetricSignal != "decision-observation-complete" ||
			o.ReverseSignal != "decision-observation-metric-reverse-complete" ||
			o.ReconstructedMetricDigest != o.MetricDigest) {
		return fmt.Errorf("bound lsp decision observation metric reverse observation is incomplete")
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("lsp decision observation metric reverse observation must remain read-only, non-executing, and non-authorizing")
	}
	if o.ProjectionDigest != digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation(o) {
		return fmt.Errorf("lsp decision observation metric reverse observation projection digest does not match its fields")
	}
	return nil
}

func digestLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation(
	observation LSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation,
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
		observation.ObservationDigest,
		strconv.FormatBool(observation.ReadOnly),
		strconv.FormatBool(observation.NonExecuting),
		strconv.FormatBool(observation.NonAuthorizing),
	}, "|"))
}
