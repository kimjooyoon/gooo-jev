package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleEvidenceCoverageFeedbackObservation struct {
	Status                   string
	MissingStage             string
	BridgeDigest             string
	CoverageMetricDigest     string
	FeedbackDigest           string
	SourceDigest             string
	CandidateSourceDigest    string
	GeneratedIRDigest        string
	ReverseObservationDigest string
	CoverageSignal           string
	FeedbackSignal           string
	FeedbackReason           string
	RequiresObservation      bool
	RequiresReview           bool
	RequiresReplan           bool
	RequiresMeasurement      bool
	FeedbackAligned          bool
	LinkSignal               string
	ObservationDigest        string
	NonExecuting             bool
	NonAuthorizing           bool
}

func ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedback(
	bridge RevisionSelfImprovementCycleFeedbackBridgeObservation,
	coverage RevisionSelfImprovementCycleEvidenceCoverageMetricObservation,
	feedback RevisionSelfImprovementCycleMetricFeedback,
) (RevisionSelfImprovementCycleEvidenceCoverageFeedbackObservation, error) {
	result := RevisionSelfImprovementCycleEvidenceCoverageFeedbackObservation{
		Status:                   "UNKNOWN",
		MissingStage:             "revision-self-improvement-cycle-evidence-coverage-feedback",
		BridgeDigest:             bridge.ObservationDigest,
		CoverageMetricDigest:     coverage.ObservationDigest,
		FeedbackDigest:           feedback.FeedbackDigest,
		SourceDigest:             coverage.SourceDigest,
		CandidateSourceDigest:    coverage.CandidateSourceDigest,
		GeneratedIRDigest:        coverage.GeneratedIRDigest,
		ReverseObservationDigest: coverage.ReverseObservationDigest,
		CoverageSignal:            coverage.CoverageSignal,
		FeedbackSignal:            feedback.FeedbackSignal,
		FeedbackReason:            feedback.FeedbackReason,
		RequiresObservation:      feedback.RequiresObservation,
		RequiresReview:           feedback.RequiresReview,
		RequiresReplan:           feedback.RequiresReplan,
		RequiresMeasurement:      feedback.RequiresMeasurement,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCycleEvidenceCoverageFeedback(result)
	}
	setDigest()

	if err := bridge.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-bridge"
		setDigest()
		return result, fmt.Errorf("cycle feedback bridge is not valid: %w", err)
	}
	if err := coverage.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-coverage"
		setDigest()
		return result, fmt.Errorf("cycle evidence coverage metric is not valid: %w", err)
	}
	if err := feedback.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-feedback"
		setDigest()
		return result, fmt.Errorf("cycle metric feedback is not valid: %w", err)
	}
	if bridge.Status != "BOUND" || coverage.Status != "BOUND" || feedback.Status != "BOUND" {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-status"
		setDigest()
		return result, fmt.Errorf("cycle evidence coverage feedback inputs are not all BOUND")
	}
	if coverage.BridgeDigest != bridge.ObservationDigest {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-coverage-bridge-link"
		setDigest()
		return result, fmt.Errorf("cycle evidence coverage metric is not linked to the feedback bridge")
	}
	if coverage.SourceDigest != bridge.SourceDigest ||
		coverage.CandidateSourceDigest != bridge.CandidateSourceDigest ||
		coverage.GeneratedIRDigest != bridge.GeneratedIRDigest ||
		coverage.ReverseObservationDigest != bridge.ReverseObservationDigest {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-evidence-link"
		setDigest()
		return result, fmt.Errorf("cycle evidence coverage metric lost bridge evidence links")
	}
	if feedback.CycleMetricDigest != bridge.MetricDigest {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-feedback-metric-link"
		setDigest()
		return result, fmt.Errorf("cycle feedback is not linked to the feedback bridge metric")
	}
	if feedback.CycleDigest != bridge.CycleDigest {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-feedback-cycle-link"
		setDigest()
		return result, fmt.Errorf("cycle feedback is not linked to the feedback bridge cycle")
	}
	if feedback.FeedbackDigest != bridge.FeedbackDigest {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-feedback-link"
		setDigest()
		return result, fmt.Errorf("cycle feedback digest is not linked to the feedback bridge")
	}
	if feedback.FeedbackSignal != bridge.FeedbackSignal ||
		feedback.FeedbackReason != bridge.FeedbackReason ||
		feedback.RequiresObservation != bridge.RequiresObservation ||
		feedback.RequiresReview != bridge.RequiresReview ||
		feedback.RequiresReplan != bridge.RequiresReplan ||
		feedback.RequiresMeasurement != bridge.RequiresMeasurement {
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback-signal-link"
		setDigest()
		return result, fmt.Errorf("cycle feedback semantics are not linked to the feedback bridge")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.FeedbackAligned = true
	result.LinkSignal = "cycle-evidence-feedback-linked"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-evidence-coverage-feedback"
		result.FeedbackAligned = false
		result.LinkSignal = ""
		setDigest()
		return result, fmt.Errorf("cycle evidence coverage feedback is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementCycleEvidenceCoverageFeedbackObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("cycle evidence coverage feedback status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound cycle evidence coverage feedback has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown cycle evidence coverage feedback has no missing stage")
	}
	for name, digest := range map[string]string{
		"bridge":              o.BridgeDigest,
		"coverage metric":    o.CoverageMetricDigest,
		"feedback":           o.FeedbackDigest,
		"source":              o.SourceDigest,
		"candidate source":   o.CandidateSourceDigest,
		"generated ir":       o.GeneratedIRDigest,
		"reverse observation": o.ReverseObservationDigest,
		"observation":        o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("cycle evidence coverage feedback %s digest is invalid", name)
		}
	}
	if o.CoverageSignal != "evidence-complete" && o.CoverageSignal != "evidence-incomplete" {
		return fmt.Errorf("cycle evidence coverage feedback coverage signal is invalid")
	}
	if o.FeedbackSignal != "observe" &&
		o.FeedbackSignal != "replan" &&
		o.FeedbackSignal != "review" {
		return fmt.Errorf("cycle evidence coverage feedback signal is invalid")
	}
	if o.FeedbackReason == "" {
		return fmt.Errorf("cycle evidence coverage feedback reason is empty")
	}
	if o.Status == "BOUND" {
		if o.CoverageSignal != "evidence-complete" {
			return fmt.Errorf("bound cycle evidence coverage feedback is incomplete")
		}
		if !o.FeedbackAligned || o.LinkSignal != "cycle-evidence-feedback-linked" {
			return fmt.Errorf("bound cycle evidence coverage feedback is not linked")
		}
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("cycle evidence coverage feedback must remain non-executing and non-authorizing")
	}
	if o.ObservationDigest != digestRevisionSelfImprovementCycleEvidenceCoverageFeedback(o) {
		return fmt.Errorf("cycle evidence coverage feedback digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleEvidenceCoverageFeedback(
	observation RevisionSelfImprovementCycleEvidenceCoverageFeedbackObservation,
) string {
	return digestString(strings.Join([]string{
		observation.Status,
		observation.MissingStage,
		observation.BridgeDigest,
		observation.CoverageMetricDigest,
		observation.FeedbackDigest,
		observation.SourceDigest,
		observation.CandidateSourceDigest,
		observation.GeneratedIRDigest,
		observation.ReverseObservationDigest,
		observation.CoverageSignal,
		observation.FeedbackSignal,
		observation.FeedbackReason,
		strconv.FormatBool(observation.RequiresObservation),
		strconv.FormatBool(observation.RequiresReview),
		strconv.FormatBool(observation.RequiresReplan),
		strconv.FormatBool(observation.RequiresMeasurement),
		strconv.FormatBool(observation.FeedbackAligned),
		observation.LinkSignal,
		strconv.FormatBool(observation.NonExecuting),
		strconv.FormatBool(observation.NonAuthorizing),
	}, "|"))
}
