package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionObservation struct {
	Status              string
	MissingStage        string
	FeedbackDigest      string
	MetricSignal        string
	FeedbackSignal      string
	FeedbackReason      string
	DecisionSignal      string
	DecisionReason      string
	RequiresObservation bool
	RequiresReplan      bool
	DecisionAligned     bool
	DecisionDigest      string
	ObservationDigest   string
	NonExecuting        bool
	NonAuthorizing      bool
}

func ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecision(
	feedback RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackObservation,
) (RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionObservation, error) {
	result := RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionObservation{
		Status:              "UNKNOWN",
		MissingStage:        "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback-decision",
		FeedbackDigest:      feedback.FeedbackDigest,
		MetricSignal:        feedback.MetricSignal,
		FeedbackSignal:      feedback.FeedbackSignal,
		FeedbackReason:      feedback.FeedbackReason,
		RequiresObservation: feedback.RequiresObservation,
		RequiresReplan:      feedback.RequiresReplan,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecision(result)
	}
	setDigest()

	if err := feedback.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback-decision-input"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision boundary metric feedback is not valid: %w", err)
	}
	if feedback.Status != "BOUND" || feedback.MissingStage != "" {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback-decision-input-status"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision boundary metric feedback is not BOUND")
	}

	switch feedback.FeedbackSignal {
	case "observe":
		result.DecisionSignal = "observe-next-cycle"
		result.DecisionReason = "complete boundary feedback permits next observation"
	case "replan":
		result.DecisionSignal = "replan-next-cycle"
		result.DecisionReason = "incomplete boundary feedback requires replan before next observation"
	default:
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback-decision-signal"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision boundary metric feedback signal is not recognized")
	}

	result.DecisionAligned = true
	result.DecisionDigest = digestString(strings.Join([]string{
		result.FeedbackDigest,
		result.DecisionSignal,
		result.DecisionReason,
	}, "|"))
	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback-decision"
		result.DecisionAligned = false
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("cycle evidence feedback decision status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound cycle evidence feedback decision has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown cycle evidence feedback decision has no missing stage")
	}
	for name, digest := range map[string]string{
		"feedback":    o.FeedbackDigest,
		"decision":    o.DecisionDigest,
		"observation": o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("cycle evidence feedback decision %s digest is invalid", name)
		}
	}
	if o.MetricSignal != "boundary-complete" && o.MetricSignal != "boundary-incomplete" {
		return fmt.Errorf("cycle evidence feedback decision metric signal is invalid")
	}
	if o.FeedbackSignal != "observe" && o.FeedbackSignal != "replan" {
		return fmt.Errorf("cycle evidence feedback decision feedback signal is invalid")
	}
	if o.FeedbackReason == "" || !o.RequiresObservation {
		return fmt.Errorf("cycle evidence feedback decision feedback evidence is incomplete")
	}
	if o.Status == "BOUND" {
		switch o.FeedbackSignal {
		case "observe":
			if o.DecisionSignal != "observe-next-cycle" ||
				o.RequiresReplan ||
				!o.DecisionAligned {
				return fmt.Errorf("observe feedback does not match its decision observation")
			}
		case "replan":
			if o.DecisionSignal != "replan-next-cycle" ||
				!o.RequiresReplan ||
				!o.DecisionAligned {
				return fmt.Errorf("replan feedback does not match its decision observation")
			}
		}
		if o.DecisionDigest != digestString(strings.Join([]string{
			o.FeedbackDigest,
			o.DecisionSignal,
			o.DecisionReason,
		}, "|")) {
			return fmt.Errorf("cycle evidence feedback decision digest does not match its fields")
		}
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("cycle evidence feedback decision must remain non-executing and non-authorizing")
	}
	if o.ObservationDigest != digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecision(o) {
		return fmt.Errorf("cycle evidence feedback decision observation digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecision(
	decision RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionObservation,
) string {
	return digestString(strings.Join([]string{
		decision.Status,
		decision.MissingStage,
		decision.FeedbackDigest,
		decision.MetricSignal,
		decision.FeedbackSignal,
		decision.FeedbackReason,
		decision.DecisionSignal,
		decision.DecisionReason,
		strconv.FormatBool(decision.RequiresObservation),
		strconv.FormatBool(decision.RequiresReplan),
		strconv.FormatBool(decision.DecisionAligned),
		strconv.FormatBool(decision.NonExecuting),
		strconv.FormatBool(decision.NonAuthorizing),
	}, "|"))
}
