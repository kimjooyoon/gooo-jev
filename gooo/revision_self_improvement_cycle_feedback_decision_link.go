package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleFeedbackDecisionLinkObservation struct {
	Status             string
	MissingStage       string
	BridgeDigest       string
	DecisionDigest     string
	QuestionDigest     string
	FeedbackSignal     string
	DecisionSignal     string
	DecisionReason     string
	ReceiptSignal      string
	DecisionAligned    bool
	LinkSignal         string
	ObservationDigest  string
	NonExecuting       bool
	NonAuthorizing     bool
}

func ObserveRevisionSelfImprovementCycleFeedbackDecisionLink(
	bridge RevisionSelfImprovementCycleFeedbackBridgeObservation,
	decision RevisionSelfImprovementDecisionReceiptObservation,
	questionDigest string,
) (RevisionSelfImprovementCycleFeedbackDecisionLinkObservation, error) {
	result := RevisionSelfImprovementCycleFeedbackDecisionLinkObservation{
		Status:            "UNKNOWN",
		MissingStage:      "revision-self-improvement-cycle-feedback-decision-link",
		BridgeDigest:      bridge.ObservationDigest,
		DecisionDigest:    decision.ObservationDigest,
		QuestionDigest:    questionDigest,
		FeedbackSignal:    bridge.FeedbackSignal,
		DecisionSignal:    decision.DecisionSignal,
		DecisionReason:    decision.DecisionReason,
		ReceiptSignal:     decision.ReceiptSignal,
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCycleFeedbackDecisionLink(result)
	}
	setDigest()

	if err := bridge.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-feedback-decision-link-bridge"
		setDigest()
		return result, fmt.Errorf("cycle feedback bridge is not valid: %w", err)
	}
	if err := decision.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-feedback-decision-link-decision"
		setDigest()
		return result, fmt.Errorf("decision receipt is not valid: %w", err)
	}
	if !validDigest(questionDigest) {
		result.MissingStage = "revision-self-improvement-cycle-feedback-decision-link-question"
		setDigest()
		return result, fmt.Errorf("decision link question digest is invalid")
	}
	if questionDigest != digestString(strings.Join([]string{"cycle-feedback", bridge.FeedbackDigest}, "|")) {
		result.MissingStage = "revision-self-improvement-cycle-feedback-decision-link-question-bridge"
		setDigest()
		return result, fmt.Errorf("decision link question is not derived from cycle feedback")
	}
	if decision.QuestionDigest != questionDigest {
		result.MissingStage = "revision-self-improvement-cycle-feedback-decision-link-question-receipt"
		setDigest()
		return result, fmt.Errorf("decision receipt question is not linked to cycle feedback")
	}
	expectedSignal := revisionSelfImprovementDecisionReceiptAllowNextIteration
	if bridge.FeedbackSignal == "replan" {
		expectedSignal = revisionSelfImprovementDecisionReceiptRequireReplan
	}
	if decision.DecisionSignal != expectedSignal {
		result.MissingStage = "revision-self-improvement-cycle-feedback-decision-link-signal"
		setDigest()
		return result, fmt.Errorf("decision receipt signal is not aligned with cycle feedback")
	}
	if bridge.Status != "BOUND" || decision.Status != revisionSelfImprovementDecisionReceiptBoundStatus {
		result.MissingStage = "revision-self-improvement-cycle-feedback-decision-link-status"
		setDigest()
		return result, fmt.Errorf("cycle feedback decision link inputs are not all BOUND")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.DecisionAligned = true
	result.LinkSignal = "cycle-feedback-decision-linked"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-feedback-decision-link"
		result.DecisionAligned = false
		setDigest()
		return result, fmt.Errorf("cycle feedback decision link is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementCycleFeedbackDecisionLinkObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("cycle feedback decision link status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound cycle feedback decision link has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown cycle feedback decision link has no missing stage")
	}
	for name, digest := range map[string]string{
		"bridge":      o.BridgeDigest,
		"decision":    o.DecisionDigest,
		"question":    o.QuestionDigest,
		"observation": o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("cycle feedback decision link %s digest is invalid", name)
		}
	}
	if o.FeedbackSignal != "observe" && o.FeedbackSignal != "replan" && o.FeedbackSignal != "review" {
		return fmt.Errorf("cycle feedback decision link feedback signal is invalid")
	}
	if o.DecisionSignal != revisionSelfImprovementDecisionReceiptAllowNextIteration &&
		o.DecisionSignal != revisionSelfImprovementDecisionReceiptRequireReplan {
		return fmt.Errorf("cycle feedback decision link decision signal is invalid")
	}
	if o.DecisionReason == "" || o.ReceiptSignal != revisionSelfImprovementDecisionReceiptAcceptedSignal {
		return fmt.Errorf("cycle feedback decision link receipt evidence is incomplete")
	}
	if o.LinkSignal != "cycle-feedback-decision-linked" || !o.DecisionAligned {
		return fmt.Errorf("cycle feedback decision link is not aligned")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("cycle feedback decision link must remain non-executing and non-authorizing")
	}
	if o.ObservationDigest != digestRevisionSelfImprovementCycleFeedbackDecisionLink(o) {
		return fmt.Errorf("cycle feedback decision link digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleFeedbackDecisionLink(
	observation RevisionSelfImprovementCycleFeedbackDecisionLinkObservation,
) string {
	return digestString(strings.Join([]string{
		observation.Status,
		observation.MissingStage,
		observation.BridgeDigest,
		observation.DecisionDigest,
		observation.QuestionDigest,
		observation.FeedbackSignal,
		observation.DecisionSignal,
		observation.DecisionReason,
		observation.ReceiptSignal,
		strconv.FormatBool(observation.DecisionAligned),
		observation.LinkSignal,
		strconv.FormatBool(observation.NonExecuting),
		strconv.FormatBool(observation.NonAuthorizing),
	}, "|"))
}
