package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackObservation struct {
	Status              string
	MissingStage        string
	BoundaryMetricDigest string
	MetricName          string
	MetricSignal        string
	FeedbackSignal      string
	FeedbackReason      string
	RequiresObservation bool
	RequiresReview      bool
	RequiresReplan      bool
	RequiresMeasurement bool
	FeedbackDigest      string
	NonExecuting        bool
	NonAuthorizing      bool
}

func ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback(
	metric RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricObservation,
) (RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackObservation, error) {
	result := RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackObservation{
		Status:               "UNKNOWN",
		MissingStage:         "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback",
		BoundaryMetricDigest: metric.ObservationDigest,
		MetricName:            metric.MetricName,
		MetricSignal:          metric.MetricSignal,
		RequiresObservation:  true,
		NonExecuting:         true,
		NonAuthorizing:       true,
	}
	setDigest := func() {
		result.FeedbackDigest = digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback(result)
	}
	setDigest()

	if err := metric.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback-metric"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision boundary metric is not valid: %w", err)
	}
	if metric.Status != "BOUND" || metric.MissingStage != "" {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback-metric-status"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision boundary metric is not BOUND")
	}

	switch metric.MetricSignal {
	case "boundary-complete":
		result.FeedbackSignal = "observe"
		result.FeedbackReason = "complete boundary requires follow-up observation"
	case "boundary-incomplete":
		result.FeedbackSignal = "replan"
		result.FeedbackReason = "incomplete boundary requires replan"
		result.RequiresReplan = true
	default:
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback-signal"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision boundary metric signal is not recognized")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision boundary metric feedback is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("cycle evidence feedback decision boundary metric feedback status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound cycle evidence feedback decision boundary metric feedback has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown cycle evidence feedback decision boundary metric feedback has no missing stage")
	}
	for name, digest := range map[string]string{
		"boundary metric": o.BoundaryMetricDigest,
		"feedback":        o.FeedbackDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("cycle evidence feedback decision boundary metric feedback %s digest is invalid", name)
		}
	}
	if o.MetricName != "evidence-feedback-decision-boundary-provenance-link-count" {
		return fmt.Errorf("cycle evidence feedback decision boundary metric feedback metric name is invalid")
	}
	if o.MetricSignal != "boundary-complete" && o.MetricSignal != "boundary-incomplete" {
		return fmt.Errorf("cycle evidence feedback decision boundary metric feedback metric signal is invalid")
	}
	if o.FeedbackSignal != "observe" && o.FeedbackSignal != "replan" {
		return fmt.Errorf("cycle evidence feedback decision boundary metric feedback signal is invalid")
	}
	if o.FeedbackReason == "" {
		return fmt.Errorf("cycle evidence feedback decision boundary metric feedback reason is empty")
	}
	if !o.RequiresObservation {
		return fmt.Errorf("cycle evidence feedback decision boundary metric feedback must require observation")
	}
	if o.Status == "BOUND" {
		switch o.MetricSignal {
		case "boundary-complete":
			if o.FeedbackSignal != "observe" ||
				o.RequiresReview ||
				o.RequiresReplan ||
				o.RequiresMeasurement {
				return fmt.Errorf("complete boundary feedback does not match its signal")
			}
		case "boundary-incomplete":
			if o.FeedbackSignal != "replan" ||
				!o.RequiresReplan ||
				o.RequiresReview ||
				o.RequiresMeasurement {
				return fmt.Errorf("incomplete boundary feedback does not match its signal")
			}
		}
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("cycle evidence feedback decision boundary metric feedback must remain non-executing and non-authorizing")
	}
	if o.FeedbackDigest != digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback(o) {
		return fmt.Errorf("cycle evidence feedback decision boundary metric feedback digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback(
	feedback RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackObservation,
) string {
	return digestString(strings.Join([]string{
		feedback.Status,
		feedback.MissingStage,
		feedback.BoundaryMetricDigest,
		feedback.MetricName,
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
