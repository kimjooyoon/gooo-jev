package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleFeedbackDecisionBoundaryObservation struct {
	Status            string
	MissingStage      string
	CycleDigest       string
	BridgeDigest      string
	DecisionLinkDigest string
	DecisionSignal    string
	BoundarySignal    string
	BoundaryAligned   bool
	ObservationDigest string
	NonExecuting      bool
	NonAuthorizing    bool
}

func ObserveRevisionSelfImprovementCycleFeedbackDecisionBoundary(
	cycle RevisionSelfImprovementCycleObservation,
	bridge RevisionSelfImprovementCycleFeedbackBridgeObservation,
	link RevisionSelfImprovementCycleFeedbackDecisionLinkObservation,
) (RevisionSelfImprovementCycleFeedbackDecisionBoundaryObservation, error) {
	result := RevisionSelfImprovementCycleFeedbackDecisionBoundaryObservation{
		Status:             "UNKNOWN",
		MissingStage:       "revision-self-improvement-cycle-feedback-decision-boundary",
		CycleDigest:        cycle.ObservationDigest,
		BridgeDigest:       bridge.ObservationDigest,
		DecisionLinkDigest: link.ObservationDigest,
		DecisionSignal:     link.DecisionSignal,
		NonExecuting:       true,
		NonAuthorizing:     true,
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCycleFeedbackDecisionBoundary(result)
	}
	setDigest()

	if err := cycle.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-feedback-decision-boundary-cycle"
		setDigest()
		return result, fmt.Errorf("cycle is not valid: %w", err)
	}
	if err := bridge.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-feedback-decision-boundary-bridge"
		setDigest()
		return result, fmt.Errorf("cycle feedback bridge is not valid: %w", err)
	}
	if err := link.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-feedback-decision-boundary-link"
		setDigest()
		return result, fmt.Errorf("cycle feedback decision link is not valid: %w", err)
	}
	if bridge.CycleDigest != cycle.ObservationDigest {
		result.MissingStage = "revision-self-improvement-cycle-feedback-decision-boundary-cycle-link"
		setDigest()
		return result, fmt.Errorf("cycle feedback bridge is not linked to the cycle")
	}
	if link.BridgeDigest != bridge.ObservationDigest {
		result.MissingStage = "revision-self-improvement-cycle-feedback-decision-boundary-bridge-link"
		setDigest()
		return result, fmt.Errorf("decision link is not linked to the cycle feedback bridge")
	}
	if cycle.Status != "BOUND" || bridge.Status != "BOUND" || link.Status != "BOUND" {
		result.MissingStage = "revision-self-improvement-cycle-feedback-decision-boundary-status"
		setDigest()
		return result, fmt.Errorf("cycle feedback decision boundary inputs are not all BOUND")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.BoundarySignal = "cycle-feedback-decision-boundary-linked"
	result.BoundaryAligned = true
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-feedback-decision-boundary"
		result.BoundaryAligned = false
		setDigest()
		return result, fmt.Errorf("cycle feedback decision boundary is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementCycleFeedbackDecisionBoundaryObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("cycle feedback decision boundary status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound cycle feedback decision boundary has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown cycle feedback decision boundary has no missing stage")
	}
	for name, digest := range map[string]string{
		"cycle":         o.CycleDigest,
		"bridge":        o.BridgeDigest,
		"decision link": o.DecisionLinkDigest,
		"observation":   o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("cycle feedback decision boundary %s digest is invalid", name)
		}
	}
	if o.DecisionSignal != revisionSelfImprovementDecisionReceiptAllowNextIteration &&
		o.DecisionSignal != revisionSelfImprovementDecisionReceiptRequireReplan {
		return fmt.Errorf("cycle feedback decision boundary decision signal is invalid")
	}
	if o.BoundarySignal != "cycle-feedback-decision-boundary-linked" || !o.BoundaryAligned {
		return fmt.Errorf("cycle feedback decision boundary is not aligned")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("cycle feedback decision boundary must remain non-executing and non-authorizing")
	}
	if o.ObservationDigest != digestRevisionSelfImprovementCycleFeedbackDecisionBoundary(o) {
		return fmt.Errorf("cycle feedback decision boundary digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleFeedbackDecisionBoundary(
	boundary RevisionSelfImprovementCycleFeedbackDecisionBoundaryObservation,
) string {
	return digestString(strings.Join([]string{
		boundary.Status,
		boundary.MissingStage,
		boundary.CycleDigest,
		boundary.BridgeDigest,
		boundary.DecisionLinkDigest,
		boundary.DecisionSignal,
		boundary.BoundarySignal,
		strconv.FormatBool(boundary.BoundaryAligned),
		strconv.FormatBool(boundary.NonExecuting),
		strconv.FormatBool(boundary.NonAuthorizing),
	}, "|"))
}
