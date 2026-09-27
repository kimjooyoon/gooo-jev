package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleFeedbackBridgeObservation struct {
	Status                   string
	MissingStage             string
	CycleDigest              string
	MetricDigest             string
	FeedbackDigest           string
	SourceDigest             string
	CandidateSourceDigest    string
	GeneratedIRDigest        string
	ReverseObservationDigest string
	MetricSignal             string
	FeedbackSignal           string
	FeedbackReason           string
	BridgeSignal             string
	RequiresObservation      bool
	RequiresReview           bool
	RequiresReplan           bool
	RequiresMeasurement      bool
	SignalsAligned            bool
	ObservationDigest        string
	NonExecuting              bool
	NonAuthorizing           bool
}

func ObserveRevisionSelfImprovementCycleFeedbackBridge(
	cycle RevisionSelfImprovementCycleObservation,
	metric RevisionSelfImprovementCycleMetricObservation,
	feedback RevisionSelfImprovementCycleMetricFeedback,
) (RevisionSelfImprovementCycleFeedbackBridgeObservation, error) {
	result := RevisionSelfImprovementCycleFeedbackBridgeObservation{
		Status:                   "UNKNOWN",
		MissingStage:             "revision-self-improvement-cycle-feedback-bridge",
		CycleDigest:              cycle.ObservationDigest,
		MetricDigest:             metric.ObservationDigest,
		FeedbackDigest:           feedback.FeedbackDigest,
		SourceDigest:             metric.SourceDigest,
		CandidateSourceDigest:    metric.CandidateSourceDigest,
		GeneratedIRDigest:        metric.GeneratedIRDigest,
		ReverseObservationDigest: metric.ReverseObservationDigest,
		MetricSignal:             metric.MetricSignal,
		FeedbackSignal:           feedback.FeedbackSignal,
		FeedbackReason:           feedback.FeedbackReason,
		RequiresObservation:      feedback.RequiresObservation,
		RequiresReview:           feedback.RequiresReview,
		RequiresReplan:           feedback.RequiresReplan,
		RequiresMeasurement:      feedback.RequiresMeasurement,
		NonExecuting:             true,
		NonAuthorizing:           true,
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCycleFeedbackBridge(result)
	}
	setDigest()

	if err := cycle.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-feedback-bridge-cycle"
		setDigest()
		return result, fmt.Errorf("cycle is not valid: %w", err)
	}
	if err := metric.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-feedback-bridge-metric"
		setDigest()
		return result, fmt.Errorf("cycle metric is not valid: %w", err)
	}
	if err := feedback.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-feedback-bridge-feedback"
		setDigest()
		return result, fmt.Errorf("cycle metric feedback is not valid: %w", err)
	}
	if metric.CycleDigest != cycle.ObservationDigest {
		result.MissingStage = "revision-self-improvement-cycle-feedback-bridge-metric-link"
		setDigest()
		return result, fmt.Errorf("cycle metric is not linked to the cycle")
	}
	if feedback.CycleMetricDigest != metric.ObservationDigest {
		result.MissingStage = "revision-self-improvement-cycle-feedback-bridge-feedback-metric-link"
		setDigest()
		return result, fmt.Errorf("cycle feedback is not linked to the cycle metric")
	}
	if feedback.CycleDigest != cycle.ObservationDigest {
		result.MissingStage = "revision-self-improvement-cycle-feedback-bridge-feedback-cycle-link"
		setDigest()
		return result, fmt.Errorf("cycle feedback is not linked to the cycle")
	}
	if cycle.Status != "BOUND" || metric.Status != "BOUND" || feedback.Status != "BOUND" {
		result.MissingStage = "revision-self-improvement-cycle-feedback-bridge-status"
		setDigest()
		return result, fmt.Errorf("cycle feedback bridge inputs are not all BOUND")
	}

	result.SignalsAligned = true
	result.BridgeSignal = "cycle-feedback-linked"
	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-feedback-bridge"
		result.SignalsAligned = false
		setDigest()
		return result, fmt.Errorf("cycle feedback bridge is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementCycleFeedbackBridgeObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("cycle feedback bridge status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound cycle feedback bridge has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown cycle feedback bridge has no missing stage")
	}
	for name, digest := range map[string]string{
		"cycle":              o.CycleDigest,
		"metric":             o.MetricDigest,
		"feedback":           o.FeedbackDigest,
		"source":             o.SourceDigest,
		"candidate source":   o.CandidateSourceDigest,
		"generated ir":       o.GeneratedIRDigest,
		"reverse observation": o.ReverseObservationDigest,
		"observation":        o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("cycle feedback bridge %s digest is invalid", name)
		}
	}
	if o.MetricSignal != "metric-improved" &&
		o.MetricSignal != "metric-regressed" &&
		o.MetricSignal != "metric-stable" {
		return fmt.Errorf("cycle feedback bridge metric signal is invalid")
	}
	if o.FeedbackSignal != "observe" &&
		o.FeedbackSignal != "replan" &&
		o.FeedbackSignal != "review" {
		return fmt.Errorf("cycle feedback bridge feedback signal is invalid")
	}
	if o.FeedbackReason == "" {
		return fmt.Errorf("cycle feedback bridge feedback reason is empty")
	}
	if o.BridgeSignal != "cycle-feedback-linked" {
		return fmt.Errorf("cycle feedback bridge signal is invalid")
	}
	if !o.SignalsAligned {
		return fmt.Errorf("bound cycle feedback bridge must be aligned")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("cycle feedback bridge must remain non-executing and non-authorizing")
	}
	if o.ObservationDigest != digestRevisionSelfImprovementCycleFeedbackBridge(o) {
		return fmt.Errorf("cycle feedback bridge digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleFeedbackBridge(
	observation RevisionSelfImprovementCycleFeedbackBridgeObservation,
) string {
	return digestString(strings.Join([]string{
		observation.Status,
		observation.MissingStage,
		observation.CycleDigest,
		observation.MetricDigest,
		observation.FeedbackDigest,
		observation.SourceDigest,
		observation.CandidateSourceDigest,
		observation.GeneratedIRDigest,
		observation.ReverseObservationDigest,
		observation.MetricSignal,
		observation.FeedbackSignal,
		observation.FeedbackReason,
		observation.BridgeSignal,
		strconv.FormatBool(observation.RequiresObservation),
		strconv.FormatBool(observation.RequiresReview),
		strconv.FormatBool(observation.RequiresReplan),
		strconv.FormatBool(observation.RequiresMeasurement),
		strconv.FormatBool(observation.SignalsAligned),
		strconv.FormatBool(observation.NonExecuting),
		strconv.FormatBool(observation.NonAuthorizing),
	}, "|"))
}