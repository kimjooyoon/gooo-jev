package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleEvidenceCoverageMetricObservation struct {
	Status                   string
	MissingStage             string
	BridgeDigest             string
	MetricName               string
	SourceDigest             string
	CandidateSourceDigest    string
	GeneratedIRDigest        string
	ReverseObservationDigest string
	LinkedEvidenceCount      int64
	RequiredEvidenceCount    int64
	CoverageSignal            string
	MetricDigest              string
	ObservationDigest         string
	NonExecuting              bool
	NonAuthorizing            bool
}

func ObserveRevisionSelfImprovementCycleEvidenceCoverageMetric(
	bridge RevisionSelfImprovementCycleFeedbackBridgeObservation,
) (RevisionSelfImprovementCycleEvidenceCoverageMetricObservation, error) {
	result := RevisionSelfImprovementCycleEvidenceCoverageMetricObservation{
		Status:                   "UNKNOWN",
		MissingStage:             "revision-self-improvement-cycle-evidence-coverage-metric",
		BridgeDigest:             bridge.ObservationDigest,
		MetricName:               "source-ir-generation-reverse-observation-coverage",
		SourceDigest:             bridge.SourceDigest,
		CandidateSourceDigest:    bridge.CandidateSourceDigest,
		GeneratedIRDigest:        bridge.GeneratedIRDigest,
		ReverseObservationDigest: bridge.ReverseObservationDigest,
		NonExecuting:             true,
		NonAuthorizing:           true,
	}
	setDigest := func() {
		result.MetricDigest = digestRevisionSelfImprovementCycleEvidenceCoverageMetric(result)
		result.ObservationDigest = digestRevisionSelfImprovementCycleEvidenceCoverageMetricObservation(result)
	}
	setDigest()

	if err := bridge.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-metric-bridge"
		setDigest()
		return result, fmt.Errorf("cycle feedback bridge is not valid: %w", err)
	}
	if bridge.Status != "BOUND" || bridge.MissingStage != "" {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-metric-status"
		setDigest()
		return result, fmt.Errorf("cycle feedback bridge is not BOUND")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.LinkedEvidenceCount = 4
	result.RequiredEvidenceCount = 4
	result.CoverageSignal = "evidence-complete"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-metric"
		setDigest()
		return result, fmt.Errorf("cycle evidence coverage metric is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementCycleEvidenceCoverageMetricObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("cycle evidence coverage metric status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound cycle evidence coverage metric has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown cycle evidence coverage metric has no missing stage")
	}
	for name, digest := range map[string]string{
		"bridge":              o.BridgeDigest,
		"source":              o.SourceDigest,
		"candidate source":    o.CandidateSourceDigest,
		"generated ir":        o.GeneratedIRDigest,
		"reverse observation": o.ReverseObservationDigest,
		"metric":              o.MetricDigest,
		"observation":         o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("cycle evidence coverage metric %s digest is invalid", name)
		}
	}
	if o.MetricName != "source-ir-generation-reverse-observation-coverage" {
		return fmt.Errorf("cycle evidence coverage metric name is invalid")
	}
	if o.LinkedEvidenceCount < 0 || o.RequiredEvidenceCount < 0 {
		return fmt.Errorf("cycle evidence coverage metric counts are negative")
	}
	if o.CoverageSignal != "evidence-complete" && o.CoverageSignal != "evidence-incomplete" {
		return fmt.Errorf("cycle evidence coverage metric signal is invalid")
	}
	if o.Status == "BOUND" && (o.LinkedEvidenceCount != 4 || o.RequiredEvidenceCount != 4 || o.CoverageSignal != "evidence-complete") {
		return fmt.Errorf("bound cycle evidence coverage metric is incomplete")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("cycle evidence coverage metric must remain non-executing and non-authorizing")
	}
	if o.MetricDigest != digestRevisionSelfImprovementCycleEvidenceCoverageMetric(o) {
		return fmt.Errorf("cycle evidence coverage metric digest does not match its fields")
	}
	if o.ObservationDigest != digestRevisionSelfImprovementCycleEvidenceCoverageMetricObservation(o) {
		return fmt.Errorf("cycle evidence coverage metric observation digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleEvidenceCoverageMetric(
	metric RevisionSelfImprovementCycleEvidenceCoverageMetricObservation,
) string {
	return digestString(strings.Join([]string{
		metric.BridgeDigest,
		metric.MetricName,
		metric.SourceDigest,
		metric.CandidateSourceDigest,
		metric.GeneratedIRDigest,
		metric.ReverseObservationDigest,
		strconv.FormatInt(metric.LinkedEvidenceCount, 10),
		strconv.FormatInt(metric.RequiredEvidenceCount, 10),
		metric.CoverageSignal,
	}, "|"))
}

func digestRevisionSelfImprovementCycleEvidenceCoverageMetricObservation(
	metric RevisionSelfImprovementCycleEvidenceCoverageMetricObservation,
) string {
	return digestString(strings.Join([]string{
		metric.Status,
		metric.MissingStage,
		metric.BridgeDigest,
		metric.MetricName,
		metric.SourceDigest,
		metric.CandidateSourceDigest,
		metric.GeneratedIRDigest,
		metric.ReverseObservationDigest,
		strconv.FormatInt(metric.LinkedEvidenceCount, 10),
		strconv.FormatInt(metric.RequiredEvidenceCount, 10),
		metric.CoverageSignal,
		metric.MetricDigest,
		strconv.FormatBool(metric.NonExecuting),
		strconv.FormatBool(metric.NonAuthorizing),
	}, "|"))
}
