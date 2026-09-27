package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryObservation struct {
	Status                   string
	MissingStage             string
	CycleDigest              string
	BridgeDigest             string
	EvidenceFeedbackDigest   string
	DecisionBoundaryDigest   string
	CoverageSignal           string
	FeedbackSignal           string
	FeedbackReason           string
	DecisionSignal           string
	BoundarySignal           string
	EvidenceFeedbackAligned  bool
	BoundaryAligned          bool
	ObservationDigest        string
	NonExecuting             bool
	NonAuthorizing           bool
}

func ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(
	cycle RevisionSelfImprovementCycleObservation,
	evidenceFeedback RevisionSelfImprovementCycleEvidenceCoverageFeedbackObservation,
	boundary RevisionSelfImprovementCycleFeedbackDecisionBoundaryObservation,
) (RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryObservation, error) {
	result := RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryObservation{
		Status:                 "UNKNOWN",
		MissingStage:           "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary",
		CycleDigest:            cycle.ObservationDigest,
		BridgeDigest:           evidenceFeedback.BridgeDigest,
		EvidenceFeedbackDigest: evidenceFeedback.ObservationDigest,
		DecisionBoundaryDigest: boundary.ObservationDigest,
		CoverageSignal:         evidenceFeedback.CoverageSignal,
		FeedbackSignal:         evidenceFeedback.FeedbackSignal,
		FeedbackReason:         evidenceFeedback.FeedbackReason,
		DecisionSignal:         boundary.DecisionSignal,
		NonExecuting:           true,
		NonAuthorizing:         true,
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(result)
	}
	setDigest()

	if err := cycle.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-cycle"
		setDigest()
		return result, fmt.Errorf("cycle is not valid: %w", err)
	}
	if err := evidenceFeedback.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-evidence-feedback"
		setDigest()
		return result, fmt.Errorf("cycle evidence coverage feedback is not valid: %w", err)
	}
	if err := boundary.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-boundary"
		setDigest()
		return result, fmt.Errorf("cycle feedback decision boundary is not valid: %w", err)
	}
	if boundary.CycleDigest != cycle.ObservationDigest {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-cycle-link"
		setDigest()
		return result, fmt.Errorf("decision boundary is not linked to the cycle")
	}
	if evidenceFeedback.BridgeDigest != boundary.BridgeDigest {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-bridge-link"
		setDigest()
		return result, fmt.Errorf("evidence feedback is not linked to the decision boundary bridge")
	}
	if evidenceFeedback.Status != "BOUND" ||
		evidenceFeedback.CoverageSignal != "evidence-complete" ||
		!evidenceFeedback.FeedbackAligned {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-evidence-status"
		setDigest()
		return result, fmt.Errorf("evidence feedback is not complete and aligned")
	}
	if cycle.Status != "BOUND" || evidenceFeedback.Status != "BOUND" || boundary.Status != "BOUND" {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-status"
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision boundary inputs are not all BOUND")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.EvidenceFeedbackAligned = true
	result.BoundaryAligned = true
	result.BoundarySignal = "cycle-evidence-feedback-decision-boundary-linked"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary"
		result.EvidenceFeedbackAligned = false
		result.BoundaryAligned = false
		result.BoundarySignal = ""
		setDigest()
		return result, fmt.Errorf("cycle evidence feedback decision boundary is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("cycle evidence feedback decision boundary status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound cycle evidence feedback decision boundary has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown cycle evidence feedback decision boundary has no missing stage")
	}
	for name, digest := range map[string]string{
		"cycle":             o.CycleDigest,
		"bridge":            o.BridgeDigest,
		"evidence feedback": o.EvidenceFeedbackDigest,
		"decision boundary": o.DecisionBoundaryDigest,
		"observation":       o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("cycle evidence feedback decision boundary %s digest is invalid", name)
		}
	}
	if o.CoverageSignal != "evidence-complete" && o.CoverageSignal != "evidence-incomplete" {
		return fmt.Errorf("cycle evidence feedback decision boundary coverage signal is invalid")
	}
	if o.FeedbackSignal != "observe" &&
		o.FeedbackSignal != "replan" &&
		o.FeedbackSignal != "review" {
		return fmt.Errorf("cycle evidence feedback decision boundary feedback signal is invalid")
	}
	if o.FeedbackReason == "" {
		return fmt.Errorf("cycle evidence feedback decision boundary feedback reason is empty")
	}
	if o.DecisionSignal != revisionSelfImprovementDecisionReceiptAllowNextIteration &&
		o.DecisionSignal != revisionSelfImprovementDecisionReceiptRequireReplan {
		return fmt.Errorf("cycle evidence feedback decision boundary decision signal is invalid")
	}
	if o.Status == "BOUND" &&
		(o.CoverageSignal != "evidence-complete" ||
			!o.EvidenceFeedbackAligned ||
			!o.BoundaryAligned ||
			o.BoundarySignal != "cycle-evidence-feedback-decision-boundary-linked") {
		return fmt.Errorf("bound cycle evidence feedback decision boundary is not aligned")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("cycle evidence feedback decision boundary must remain non-executing and non-authorizing")
	}
	if o.ObservationDigest != digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(o) {
		return fmt.Errorf("cycle evidence feedback decision boundary digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(
	boundary RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryObservation,
) string {
	return digestString(strings.Join([]string{
		boundary.Status,
		boundary.MissingStage,
		boundary.CycleDigest,
		boundary.BridgeDigest,
		boundary.EvidenceFeedbackDigest,
		boundary.DecisionBoundaryDigest,
		boundary.CoverageSignal,
		boundary.FeedbackSignal,
		boundary.FeedbackReason,
		boundary.DecisionSignal,
		boundary.BoundarySignal,
		strconv.FormatBool(boundary.EvidenceFeedbackAligned),
		strconv.FormatBool(boundary.BoundaryAligned),
		strconv.FormatBool(boundary.NonExecuting),
		strconv.FormatBool(boundary.NonAuthorizing),
	}, "|"))
}
