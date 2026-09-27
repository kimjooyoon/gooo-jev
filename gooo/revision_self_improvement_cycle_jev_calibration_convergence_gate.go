package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	jevCalibrationConvergenceGateConverged = "jev-calibration-convergence-converged"
	jevCalibrationConvergenceGateRegressed = "jev-calibration-convergence-regressed"
	jevCalibrationConvergenceGateDefer = "jev-calibration-convergence-defer"
	jevCalibrationConvergenceGateUnknown = "jev-calibration-convergence-unknown"

	jevCalibrationConvergenceGateChoiceStable = "jev-calibration-choice-set-sensitivity-stable"
	jevCalibrationConvergenceGateChoiceChanged = "jev-calibration-choice-set-sensitivity-choice-set-changed"
	jevCalibrationConvergenceGateChoiceOrderChanged = "jev-calibration-choice-set-sensitivity-order-changed"
	jevCalibrationConvergenceGateChoiceBothChanged = "jev-calibration-choice-set-sensitivity-choice-set-and-order-changed"
	jevCalibrationConvergenceGateChoiceUnknown = "jev-calibration-choice-set-sensitivity-unknown"
	jevCalibrationConvergenceGateChoiceCompare = "jev-calibration-choice-set-decision-compare"
	jevCalibrationConvergenceGateChoiceDefer = "jev-calibration-choice-set-decision-defer"
)

type RevisionSelfImprovementCycleJEVCalibrationConvergenceGateInput struct {
	OutcomeDelta                          RevisionSelfImprovementCycleJEVCalibrationOutcomeDeltaObservation
	ChoiceSetSensitivityStatus            string
	ChoiceSetSensitivitySignal            string
	ChoiceSetDecisionSignal               string
	ChoiceSetSensitivityObservationDigest string
	Coverage                               RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetric
}

type RevisionSelfImprovementCycleJEVCalibrationConvergenceGateObservation struct {
	Status                                string
	MissingStage                          string
	OutcomeDeltaSignal                    string
	ChoiceSetSensitivitySignal            string
	ChoiceSetDecisionSignal               string
	OutcomeDeltaObservationDigest         string
	ChoiceSetSensitivityObservationDigest string
	SourceObservationDigest               string
	ReverseObservationCoverageDigest      string
	CoverageMetricObservationDigest       string
	CoverageMetricSignal                  string
	GateDecision                          string
	ObservationDigest                     string
	ReadOnly                              bool
	NonExecuting                          bool
	NonAuthorizing                        bool
}

func ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGate(
	input RevisionSelfImprovementCycleJEVCalibrationConvergenceGateInput,
) RevisionSelfImprovementCycleJEVCalibrationConvergenceGateObservation {
	result := RevisionSelfImprovementCycleJEVCalibrationConvergenceGateObservation{
		Status:                     "UNKNOWN",
		MissingStage:               "revision-self-improvement-cycle-jev-calibration-convergence-gate",
		ChoiceSetSensitivitySignal: jevCalibrationConvergenceGateChoiceUnknown,
		ChoiceSetDecisionSignal:    jevCalibrationConvergenceGateChoiceDefer,
		CoverageMetricSignal:       jevGenerationTraceCoverageUnknown,
		GateDecision:               jevCalibrationConvergenceGateUnknown,
		ReadOnly:                   true,
		NonExecuting:               true,
		NonAuthorizing:             true,
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCycleJEVCalibrationConvergenceGate(result)
	}
	setDigest()

	if err := input.OutcomeDelta.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-convergence-gate-outcome-delta"
		setDigest()
		return result
	}
	if input.OutcomeDelta.Status != "BOUND" {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-convergence-gate-outcome-delta-status"
		setDigest()
		return result
	}
	if err := input.Coverage.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-convergence-gate-coverage"
		setDigest()
		return result
	}
	if input.Coverage.Status != "BOUND" {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-convergence-gate-coverage-status"
		setDigest()
		return result
	}
	if input.OutcomeDelta.ReverseObservationCoverageDigest != input.Coverage.ObservationDigest {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-convergence-gate-reverse-coverage-lineage"
		setDigest()
		return result
	}
	if input.ChoiceSetSensitivityStatus != "BOUND" && input.ChoiceSetSensitivityStatus != "UNKNOWN" {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-convergence-gate-choice-set-sensitivity-status"
		setDigest()
		return result
	}
	if !validDigest(input.ChoiceSetSensitivityObservationDigest) {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-convergence-gate-choice-set-sensitivity-digest"
		setDigest()
		return result
	}
	if input.ChoiceSetSensitivityStatus == "UNKNOWN" {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-convergence-gate-choice-set-sensitivity"
		setDigest()
		return result
	}
	if !validJEVCalibrationConvergenceGateChoiceEvidence(
		input.ChoiceSetSensitivitySignal,
		input.ChoiceSetDecisionSignal,
	) {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-convergence-gate-choice-set-sensitivity-evidence"
		setDigest()
		return result
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.OutcomeDeltaSignal = input.OutcomeDelta.DeltaSignal
	result.ChoiceSetSensitivitySignal = input.ChoiceSetSensitivitySignal
	result.ChoiceSetDecisionSignal = input.ChoiceSetDecisionSignal
	result.OutcomeDeltaObservationDigest = input.OutcomeDelta.ObservationDigest
	result.ChoiceSetSensitivityObservationDigest = input.ChoiceSetSensitivityObservationDigest
	result.SourceObservationDigest = input.OutcomeDelta.SourceObservationDigest
	result.ReverseObservationCoverageDigest = input.OutcomeDelta.ReverseObservationCoverageDigest
	result.CoverageMetricObservationDigest = input.Coverage.ObservationDigest
	result.CoverageMetricSignal = input.Coverage.MetricSignal
	switch {
	case input.ChoiceSetDecisionSignal == jevCalibrationConvergenceGateChoiceDefer:
		result.GateDecision = jevCalibrationConvergenceGateDefer
	case input.OutcomeDelta.DeltaSignal == jevCalibrationOutcomeDeltaImproved:
		result.GateDecision = jevCalibrationConvergenceGateConverged
	case input.OutcomeDelta.DeltaSignal == jevCalibrationOutcomeDeltaRegressed:
		result.GateDecision = jevCalibrationConvergenceGateRegressed
	case input.OutcomeDelta.DeltaSignal == jevCalibrationOutcomeDeltaFlat:
		result.GateDecision = jevCalibrationConvergenceGateDefer
	default:
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-convergence-gate-outcome-delta-evidence"
		result.GateDecision = jevCalibrationConvergenceGateUnknown
	}
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-convergence-gate"
		result.GateDecision = jevCalibrationConvergenceGateUnknown
		setDigest()
	}
	return result
}

func validJEVCalibrationConvergenceGateChoiceEvidence(signal, decision string) bool {
	switch signal {
	case jevCalibrationConvergenceGateChoiceStable:
		return decision == jevCalibrationConvergenceGateChoiceCompare
	case jevCalibrationConvergenceGateChoiceChanged,
		jevCalibrationConvergenceGateChoiceOrderChanged,
		jevCalibrationConvergenceGateChoiceBothChanged:
		return decision == jevCalibrationConvergenceGateChoiceDefer
	default:
		return false
	}
}

func (value RevisionSelfImprovementCycleJEVCalibrationConvergenceGateObservation) Validate() error {
	if value.Status != "BOUND" && value.Status != "UNKNOWN" {
		return fmt.Errorf("calibration convergence gate status is invalid")
	}
	if value.Status == "BOUND" && value.MissingStage != "" {
		return fmt.Errorf("bound calibration convergence gate has a missing stage")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" || value.GateDecision != jevCalibrationConvergenceGateUnknown {
			return fmt.Errorf("unknown calibration convergence gate must preserve incomplete evidence")
		}
	} else {
		if !validDigest(value.OutcomeDeltaObservationDigest) ||
			!validDigest(value.ChoiceSetSensitivityObservationDigest) ||
			!validDigest(value.SourceObservationDigest) ||
			!validDigest(value.ReverseObservationCoverageDigest) ||
			!validDigest(value.CoverageMetricObservationDigest) ||
			value.CoverageMetricObservationDigest != value.ReverseObservationCoverageDigest ||
			value.CoverageMetricSignal != jevGenerationTraceCoverageComplete ||
			!validJEVCalibrationConvergenceGateChoiceEvidence(
				value.ChoiceSetSensitivitySignal,
				value.ChoiceSetDecisionSignal,
			) {
			return fmt.Errorf("bound calibration convergence gate evidence is invalid")
		}
		switch value.OutcomeDeltaSignal {
		case jevCalibrationOutcomeDeltaImproved:
			if value.ChoiceSetDecisionSignal == jevCalibrationConvergenceGateChoiceCompare &&
				value.GateDecision != jevCalibrationConvergenceGateConverged {
				return fmt.Errorf("improved comparable calibration must converge")
			}
		case jevCalibrationOutcomeDeltaRegressed:
			if value.ChoiceSetDecisionSignal == jevCalibrationConvergenceGateChoiceCompare &&
				value.GateDecision != jevCalibrationConvergenceGateRegressed {
				return fmt.Errorf("regressed comparable calibration must regress")
			}
		case jevCalibrationOutcomeDeltaFlat:
			if value.GateDecision != jevCalibrationConvergenceGateDefer {
				return fmt.Errorf("flat calibration must defer")
			}
		default:
			return fmt.Errorf("bound calibration convergence gate signal is invalid")
		}
		if value.ChoiceSetDecisionSignal == jevCalibrationConvergenceGateChoiceDefer &&
			value.GateDecision != jevCalibrationConvergenceGateDefer {
			return fmt.Errorf("changed choice evidence must defer")
		}
		switch value.GateDecision {
		case jevCalibrationConvergenceGateConverged,
			jevCalibrationConvergenceGateRegressed,
			jevCalibrationConvergenceGateDefer:
		default:
			return fmt.Errorf("bound calibration convergence gate decision is invalid")
		}
	}
	if !value.ReadOnly || !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("calibration convergence gate must remain read-only, non-executing, and non-authorizing")
	}
	if value.ObservationDigest != digestRevisionSelfImprovementCycleJEVCalibrationConvergenceGate(value) {
		return fmt.Errorf("calibration convergence gate digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleJEVCalibrationConvergenceGate(
	value RevisionSelfImprovementCycleJEVCalibrationConvergenceGateObservation,
) string {
	return digestString(strings.Join([]string{
		value.Status,
		value.MissingStage,
		value.OutcomeDeltaSignal,
		value.ChoiceSetSensitivitySignal,
		value.ChoiceSetDecisionSignal,
		value.OutcomeDeltaObservationDigest,
		value.ChoiceSetSensitivityObservationDigest,
		value.SourceObservationDigest,
		value.ReverseObservationCoverageDigest,
		value.CoverageMetricObservationDigest,
		value.CoverageMetricSignal,
		value.GateDecision,
		strconv.FormatBool(value.ReadOnly),
		strconv.FormatBool(value.NonExecuting),
		strconv.FormatBool(value.NonAuthorizing),
	}, "|"))
}