package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation records a
// deterministic confidence metric for one typed JEV decision signal.
type RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation struct {
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
	ReadOnly                  bool
	NonExecuting              bool
	NonAuthorizing            bool
}

func ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(
	signal RevisionSelfImprovementCycleJEVTypedDecisionSignalObservation,
) (RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation, error) {
	result := RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation{
		Status:                    "UNKNOWN",
		MissingStage:              "revision-self-improvement-cycle-jev-decision-confidence-metric",
		MetricName:                "jev-typed-decision-confidence-milli",
		QuestionID:                signal.QuestionID,
		QuestionKind:              signal.QuestionKind,
		ConfidenceMilli:           signal.ConfidenceMilli,
		ConfidenceBand:            "unknown",
		DecisionDigest:            signal.DecisionDigest,
		EvidenceDigest:            signal.EvidenceDigest,
		DecisionObservationDigest: signal.DecisionObservationDigest,
		RequiredDigestCount:       3,
		MetricSignal:              "jev-decision-confidence-unknown",
		ReadOnly:                  true,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	setDigest := func() {
		result.MetricDigest = digestRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(result)
	}
	setDigest()

	if err := signal.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-jev-decision-confidence-metric-input"
		setDigest()
		return result, fmt.Errorf("jev typed decision signal is not valid for confidence metric: %w", err)
	}
	if signal.Status != "BOUND" || signal.MissingStage != "" {
		result.MissingStage = "revision-self-improvement-cycle-jev-decision-confidence-metric-input-status"
		setDigest()
		return result, fmt.Errorf("jev typed decision signal is not BOUND")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.ConfidenceBand = jevDecisionConfidenceBand(signal.QuestionKind, signal.ConfidenceMilli)
	result.LinkedDigestCount = 3
	result.MetricSignal = "jev-decision-confidence-observed"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-jev-decision-confidence-metric"
		result.ConfidenceBand = "unknown"
		result.LinkedDigestCount = 0
		result.MetricSignal = "jev-decision-confidence-unknown"
		setDigest()
		return result, fmt.Errorf("jev decision confidence metric is not valid: %w", err)
	}
	return result, nil
}

func jevDecisionConfidenceBand(questionKind string, confidenceMilli int64) string {
	if questionKind == "noul" {
		return "none"
	}
	if confidenceMilli >= 800 {
		return "high"
	}
	if confidenceMilli >= 500 {
		return "medium"
	}
	return "low"
}

func (o RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("jev decision confidence metric status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound jev decision confidence metric has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown jev decision confidence metric has no missing stage")
	}
	if o.MetricName != "jev-typed-decision-confidence-milli" {
		return fmt.Errorf("jev decision confidence metric name is invalid")
	}
	if err := (RevisionSelfImprovementCycleJEVTypedDecisionSignalObservationInput{
		QuestionID:       o.QuestionID,
		QuestionKind:     o.QuestionKind,
		DecisionDigest:   o.DecisionDigest,
		EvidenceDigest:   o.EvidenceDigest,
		ConfidenceMilli:  o.ConfidenceMilli,
	}).Validate(); err != nil {
		return fmt.Errorf("jev decision confidence metric input is invalid: %w", err)
	}
	if !validDigest(o.DecisionObservationDigest) || !validDigest(o.MetricDigest) {
		return fmt.Errorf("jev decision confidence metric digest is invalid")
	}
	if o.ConfidenceBand != "high" && o.ConfidenceBand != "medium" &&
		o.ConfidenceBand != "low" && o.ConfidenceBand != "none" &&
		o.ConfidenceBand != "unknown" {
		return fmt.Errorf("jev decision confidence band is invalid")
	}
	if o.MetricSignal != "jev-decision-confidence-observed" &&
		o.MetricSignal != "jev-decision-confidence-unknown" {
		return fmt.Errorf("jev decision confidence metric signal is invalid")
	}
	if o.Status == "BOUND" &&
		(o.LinkedDigestCount != 3 ||
			o.RequiredDigestCount != 3 ||
			o.MetricSignal != "jev-decision-confidence-observed" ||
			o.ConfidenceBand != jevDecisionConfidenceBand(o.QuestionKind, o.ConfidenceMilli)) {
		return fmt.Errorf("bound jev decision confidence metric is incomplete")
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("jev decision confidence metric must remain read-only, non-executing, and non-authorizing")
	}
	if o.MetricDigest != digestRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(o) {
		return fmt.Errorf("jev decision confidence metric digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(
	metric RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation,
) string {
	return digestString(strings.Join([]string{
		metric.Status,
		metric.MissingStage,
		metric.MetricName,
		metric.QuestionID,
		metric.QuestionKind,
		strconv.FormatInt(metric.ConfidenceMilli, 10),
		metric.ConfidenceBand,
		metric.DecisionDigest,
		metric.EvidenceDigest,
		metric.DecisionObservationDigest,
		strconv.FormatInt(metric.LinkedDigestCount, 10),
		strconv.FormatInt(metric.RequiredDigestCount, 10),
		metric.MetricSignal,
		strconv.FormatBool(metric.ReadOnly),
		strconv.FormatBool(metric.NonExecuting),
		strconv.FormatBool(metric.NonAuthorizing),
	}, "|"))
}
