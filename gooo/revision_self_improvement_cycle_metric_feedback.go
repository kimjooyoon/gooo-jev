package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleMetricFeedback struct {
	Status            string
	MissingStage      string
	CycleMetricDigest string
	CycleDigest       string
	MetricName        string
	Direction         string
	BaselineValue     int64
	CandidateValue    int64
	MetricDelta       int64
	MetricSignal      string
	FeedbackSignal    string
	FeedbackReason    string
	RequiresObservation bool
	RequiresReview      bool
	RequiresReplan      bool
	RequiresMeasurement bool
	FeedbackDigest     string
	NonExecuting       bool
	NonAuthorizing     bool
}

func ObserveRevisionSelfImprovementCycleMetricFeedback(
	metric RevisionSelfImprovementCycleMetricObservation,
) (RevisionSelfImprovementCycleMetricFeedback, error) {
	feedback := RevisionSelfImprovementCycleMetricFeedback{
		Status:              "UNKNOWN",
		MissingStage:        "revision-self-improvement-cycle-metric-feedback",
		CycleMetricDigest:   metric.ObservationDigest,
		CycleDigest:         metric.CycleDigest,
		MetricName:          metric.MetricName,
		Direction:           metric.Direction,
		BaselineValue:       metric.BaselineValue,
		CandidateValue:      metric.CandidateValue,
		MetricDelta:         metric.MetricDelta,
		MetricSignal:        metric.MetricSignal,
		RequiresObservation: true,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	setDigest := func() {
		feedback.FeedbackDigest = digestRevisionSelfImprovementCycleMetricFeedback(feedback)
	}
	setDigest()

	if err := metric.Validate(); err != nil {
		feedback.MissingStage = "revision-self-improvement-cycle-metric-feedback-metric"
		setDigest()
		return feedback, fmt.Errorf("cycle metric is not valid: %w", err)
	}
	if metric.Status != "BOUND" || metric.MissingStage != "" {
		feedback.MissingStage = "revision-self-improvement-cycle-metric-feedback-metric-status"
		setDigest()
		return feedback, fmt.Errorf("cycle metric is not BOUND")
	}

	switch metric.MetricSignal {
	case "metric-improved":
		feedback.FeedbackSignal = "observe"
		feedback.FeedbackReason = "metric-improved-requires-follow-up-observation"
	case "metric-regressed":
		feedback.FeedbackSignal = "replan"
		feedback.FeedbackReason = "metric-regressed-requires-replan"
		feedback.RequiresReplan = true
	case "metric-stable":
		feedback.FeedbackSignal = "review"
		feedback.FeedbackReason = "metric-stable-requires-review"
		feedback.RequiresReview = true
		feedback.RequiresMeasurement = true
	default:
		feedback.MissingStage = "revision-self-improvement-cycle-metric-feedback-signal"
		setDigest()
		return feedback, fmt.Errorf("cycle metric signal is not recognized")
	}

	feedback.Status = "BOUND"
	feedback.MissingStage = ""
	setDigest()
	if err := feedback.Validate(); err != nil {
		feedback.Status = "UNKNOWN"
		feedback.MissingStage = "revision-self-improvement-cycle-metric-feedback"
		setDigest()
		return feedback, fmt.Errorf("cycle metric feedback is not valid: %w", err)
	}
	return feedback, nil
}

func (f RevisionSelfImprovementCycleMetricFeedback) Validate() error {
	if f.Status != "BOUND" && f.Status != "UNKNOWN" {
		return fmt.Errorf("cycle metric feedback status is invalid")
	}
	if f.Status == "BOUND" && f.MissingStage != "" {
		return fmt.Errorf("bound cycle metric feedback has a missing stage")
	}
	if f.Status == "UNKNOWN" && f.MissingStage == "" {
		return fmt.Errorf("unknown cycle metric feedback has no missing stage")
	}
	for name, digest := range map[string]string{
		"cycle metric": f.CycleMetricDigest,
		"cycle":        f.CycleDigest,
		"feedback":     f.FeedbackDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("cycle metric feedback %s digest is invalid", name)
		}
	}
	if f.MetricName == "" {
		return fmt.Errorf("cycle metric feedback metric name is empty")
	}
	if f.Direction != "higher-is-better" && f.Direction != "lower-is-better" {
		return fmt.Errorf("cycle metric feedback direction is invalid")
	}
	if f.MetricDelta != f.CandidateValue-f.BaselineValue {
		return fmt.Errorf("cycle metric feedback delta is invalid")
	}
	if f.MetricSignal != "metric-improved" &&
		f.MetricSignal != "metric-regressed" &&
		f.MetricSignal != "metric-stable" {
		return fmt.Errorf("cycle metric feedback metric signal is invalid")
	}
	if f.FeedbackSignal != "observe" &&
		f.FeedbackSignal != "replan" &&
		f.FeedbackSignal != "review" {
		return fmt.Errorf("cycle metric feedback signal is invalid")
	}
	if !f.RequiresObservation {
		return fmt.Errorf("cycle metric feedback must require observation")
	}
	if f.Status == "BOUND" {
		switch f.MetricSignal {
		case "metric-improved":
			if f.FeedbackSignal != "observe" || f.RequiresReview || f.RequiresReplan || f.RequiresMeasurement {
				return fmt.Errorf("improved metric feedback does not match its signal")
			}
		case "metric-regressed":
			if f.FeedbackSignal != "replan" || !f.RequiresReplan || f.RequiresReview || f.RequiresMeasurement {
				return fmt.Errorf("regressed metric feedback does not match its signal")
			}
		case "metric-stable":
			if f.FeedbackSignal != "review" || f.RequiresReplan || !f.RequiresReview || !f.RequiresMeasurement {
				return fmt.Errorf("stable metric feedback does not match its signal")
			}
		}
	}
	if !f.NonExecuting || !f.NonAuthorizing {
		return fmt.Errorf("cycle metric feedback must remain non-executing and non-authorizing")
	}
	if f.FeedbackDigest != digestRevisionSelfImprovementCycleMetricFeedback(f) {
		return fmt.Errorf("cycle metric feedback digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleMetricFeedback(
	feedback RevisionSelfImprovementCycleMetricFeedback,
) string {
	return digestString(strings.Join([]string{
		feedback.Status,
		feedback.MissingStage,
		feedback.CycleMetricDigest,
		feedback.CycleDigest,
		feedback.MetricName,
		feedback.Direction,
		strconv.FormatInt(feedback.BaselineValue, 10),
		strconv.FormatInt(feedback.CandidateValue, 10),
		strconv.FormatInt(feedback.MetricDelta, 10),
		feedback.MetricSignal,
		feedback.FeedbackSignal,
		feedback.FeedbackReason,
		strconv.FormatBool(feedback.RequiresObservation),
		strconv.FormatBool(feedback.RequiresReview),
		strconv.FormatBool(feedback.RequiresReplan),
		strconv.FormatBool(feedback.RequiresMeasurement),
		strconv.FormatBool(feedback.NonExecuting),
		strconv.FormatBool(feedback.NonAuthorizing),
	}, "|"))
}