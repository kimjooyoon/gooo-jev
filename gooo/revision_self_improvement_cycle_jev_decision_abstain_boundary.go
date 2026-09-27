package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleJEVDecisionAbstainBoundaryPolicy struct {
	AcceptThresholdMilli     int64
	AbstainLowerMilli         int64
	AbstainUpperMilli         int64
	GenerationDigest          string
	ReverseObservationDigest  string
}

type RevisionSelfImprovementCycleJEVDecisionAbstainBoundaryObservation struct {
	Status                    string
	MissingStage              string
	MetricName                string
	QuestionID                string
	QuestionKind              string
	ConfidenceMilli            int64
	AcceptThresholdMilli      int64
	AbstainLowerMilli          int64
	AbstainUpperMilli          int64
	Disposition                string
	CalibrationStatus          string
	CalibrationRequired        bool
	DecisionDigest             string
	EvidenceDigest             string
	DecisionObservationDigest  string
	MetricDigest               string
	GenerationDigest           string
	ReverseObservationDigest   string
	PolicyDigest               string
	BoundarySignal             string
	ObservationDigest          string
	ReadOnly                  bool
	NonExecuting              bool
	NonAuthorizing             bool
}

func ObserveRevisionSelfImprovementCycleJEVDecisionAbstainBoundary(
	metric RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation,
	policy RevisionSelfImprovementCycleJEVDecisionAbstainBoundaryPolicy,
) RevisionSelfImprovementCycleJEVDecisionAbstainBoundaryObservation {
	result := RevisionSelfImprovementCycleJEVDecisionAbstainBoundaryObservation{
		Status: "UNKNOWN",
		MissingStage: "revision-self-improvement-cycle-jev-decision-abstain-boundary",
		MetricName: metric.MetricName, QuestionID: metric.QuestionID,
		QuestionKind: metric.QuestionKind, ConfidenceMilli: metric.ConfidenceMilli,
		AcceptThresholdMilli: policy.AcceptThresholdMilli,
		AbstainLowerMilli: policy.AbstainLowerMilli, AbstainUpperMilli: policy.AbstainUpperMilli,
		Disposition: "unknown", CalibrationStatus: "unverified", CalibrationRequired: true,
		DecisionDigest: metric.DecisionDigest, EvidenceDigest: metric.EvidenceDigest,
		DecisionObservationDigest: metric.DecisionObservationDigest, MetricDigest: metric.MetricDigest,
		GenerationDigest: policy.GenerationDigest,
		ReverseObservationDigest: policy.ReverseObservationDigest,
		PolicyDigest: digestRevisionSelfImprovementCycleJEVDecisionAbstainBoundaryPolicy(policy),
		BoundarySignal: "jev-decision-abstain-boundary-unknown",
		ReadOnly: true, NonExecuting: true, NonAuthorizing: true,
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCycleJEVDecisionAbstainBoundary(result)
	}
	setDigest()
	if err := metric.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-jev-decision-abstain-boundary-input-metric"
		setDigest()
		return result
	}
	if metric.Status != "BOUND" || metric.MissingStage != "" {
		result.MissingStage = "revision-self-improvement-cycle-jev-decision-abstain-boundary-input-metric-status"
		setDigest()
		return result
	}
	if metric.QuestionKind == "noul" {
		result.MissingStage = "revision-self-improvement-cycle-jev-decision-abstain-boundary-noul-confidence"
		setDigest()
		return result
	}
	if err := policy.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-jev-decision-abstain-boundary-policy"
		setDigest()
		return result
	}
	result.Status, result.MissingStage = "BOUND", ""
	switch {
	case metric.ConfidenceMilli >= policy.AcceptThresholdMilli:
		result.Disposition = "direct_filter_candidate"
	case metric.ConfidenceMilli >= policy.AbstainLowerMilli && metric.ConfidenceMilli <= policy.AbstainUpperMilli:
		result.Disposition = "abstain_for_calibration"
	default:
		result.Disposition = "escalate_to_full_validator"
	}
	result.BoundarySignal = "jev-decision-abstain-boundary-observed"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status, result.MissingStage = "UNKNOWN", "revision-self-improvement-cycle-jev-decision-abstain-boundary"
		result.Disposition, result.BoundarySignal = "unknown", "jev-decision-abstain-boundary-unknown"
		setDigest()
	}
	return result
}

func (p RevisionSelfImprovementCycleJEVDecisionAbstainBoundaryPolicy) Validate() error {
	if p.AcceptThresholdMilli < 0 || p.AcceptThresholdMilli > 1000 ||
		p.AbstainLowerMilli < 0 || p.AbstainLowerMilli > 1000 ||
		p.AbstainUpperMilli < 0 || p.AbstainUpperMilli > 1000 {
		return fmt.Errorf("boundary thresholds must be between 0 and 1000 milli")
	}
	if p.AbstainLowerMilli > p.AbstainUpperMilli || p.AbstainUpperMilli >= p.AcceptThresholdMilli {
		return fmt.Errorf("abstain interval must precede the accept threshold")
	}
	if !validDigest(p.GenerationDigest) || !validDigest(p.ReverseObservationDigest) {
		return fmt.Errorf("boundary policy lineage digests are invalid")
	}
	return nil
}

func (o RevisionSelfImprovementCycleJEVDecisionAbstainBoundaryObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("jev decision abstain boundary status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound boundary has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown boundary has no missing stage")
	}
	if o.Status == "UNKNOWN" {
		if o.Disposition != "unknown" || o.BoundarySignal != "jev-decision-abstain-boundary-unknown" ||
			o.CalibrationStatus != "unverified" || !o.CalibrationRequired {
			return fmt.Errorf("unknown boundary must preserve unresolved routing")
		}
	} else {
		if o.MetricName != "jev-typed-decision-confidence-milli" || o.QuestionKind == "noul" ||
			o.CalibrationStatus != "unverified" || !o.CalibrationRequired {
			return fmt.Errorf("bound boundary has invalid metric state")
		}
		if o.Disposition != "direct_filter_candidate" && o.Disposition != "abstain_for_calibration" &&
			o.Disposition != "escalate_to_full_validator" {
			return fmt.Errorf("bound boundary disposition is invalid")
		}
		policy := RevisionSelfImprovementCycleJEVDecisionAbstainBoundaryPolicy{
			AcceptThresholdMilli: o.AcceptThresholdMilli, AbstainLowerMilli: o.AbstainLowerMilli,
			AbstainUpperMilli: o.AbstainUpperMilli, GenerationDigest: o.GenerationDigest,
			ReverseObservationDigest: o.ReverseObservationDigest,
		}
		if err := policy.Validate(); err != nil {
			return fmt.Errorf("bound boundary policy is invalid: %w", err)
		}
		for name, value := range map[string]string{
			"decision": o.DecisionDigest, "evidence": o.EvidenceDigest,
			"observation": o.DecisionObservationDigest, "metric": o.MetricDigest,
			"generation": o.GenerationDigest, "reverse": o.ReverseObservationDigest,
			"policy": o.PolicyDigest,
		} {
			if !validDigest(value) {
				return fmt.Errorf("bound boundary %s digest is invalid", name)
			}
		}
		if o.BoundarySignal != "jev-decision-abstain-boundary-observed" {
			return fmt.Errorf("bound boundary signal is invalid")
		}
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("boundary must remain read-only, non-executing, and non-authorizing")
	}
	if o.ObservationDigest != digestRevisionSelfImprovementCycleJEVDecisionAbstainBoundary(o) {
		return fmt.Errorf("boundary digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleJEVDecisionAbstainBoundaryPolicy(
	p RevisionSelfImprovementCycleJEVDecisionAbstainBoundaryPolicy,
) string {
	return digestString(strings.Join([]string{
		strconv.FormatInt(p.AcceptThresholdMilli, 10),
		strconv.FormatInt(p.AbstainLowerMilli, 10),
		strconv.FormatInt(p.AbstainUpperMilli, 10),
		p.GenerationDigest, p.ReverseObservationDigest,
	}, "|"))
}

func digestRevisionSelfImprovementCycleJEVDecisionAbstainBoundary(
	o RevisionSelfImprovementCycleJEVDecisionAbstainBoundaryObservation,
) string {
	return digestString(strings.Join([]string{
		o.Status, o.MissingStage, o.MetricName, o.QuestionID, o.QuestionKind,
		strconv.FormatInt(o.ConfidenceMilli, 10),
		strconv.FormatInt(o.AcceptThresholdMilli, 10),
		strconv.FormatInt(o.AbstainLowerMilli, 10),
		strconv.FormatInt(o.AbstainUpperMilli, 10),
		o.Disposition, o.CalibrationStatus, strconv.FormatBool(o.CalibrationRequired),
		o.DecisionDigest, o.EvidenceDigest, o.DecisionObservationDigest, o.MetricDigest,
		o.GenerationDigest, o.ReverseObservationDigest, o.PolicyDigest, o.BoundarySignal,
		strconv.FormatBool(o.ReadOnly), strconv.FormatBool(o.NonExecuting),
		strconv.FormatBool(o.NonAuthorizing),
	}, "|"))
}