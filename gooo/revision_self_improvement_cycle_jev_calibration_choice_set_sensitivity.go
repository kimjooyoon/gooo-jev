package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	jevCalibrationChoiceSetSensitivityStable = "jev-calibration-choice-set-sensitivity-stable"
	jevCalibrationChoiceSetSensitivityChanged = "jev-calibration-choice-set-sensitivity-choice-set-changed"
	jevCalibrationChoiceSetSensitivityOrderChanged = "jev-calibration-choice-set-sensitivity-order-changed"
	jevCalibrationChoiceSetSensitivityBothChanged = "jev-calibration-choice-set-sensitivity-choice-set-and-order-changed"
	jevCalibrationChoiceSetSensitivityUnknown = "jev-calibration-choice-set-sensitivity-unknown"
	jevCalibrationChoiceSetDecisionCompare = "jev-calibration-choice-set-decision-compare"
	jevCalibrationChoiceSetDecisionDefer = "jev-calibration-choice-set-decision-defer"
)

type RevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivityInput struct {
	Baseline                       RevisionSelfImprovementCycleJEVCalibrationOutcomeWindowObservation
	Current                        RevisionSelfImprovementCycleJEVCalibrationOutcomeWindowObservation
	BaselineChoiceSetDigest        string
	CurrentChoiceSetDigest         string
	BaselineChoiceOrderDigest      string
	CurrentChoiceOrderDigest       string
}

type RevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivityObservation struct {
	Status                            string
	MissingStage                      string
	BaselineAccuracyMilli             int64
	CurrentAccuracyMilli              int64
	AccuracyDeltaMilli                int64
	BaselineMeanConfidenceMilli       int64
	CurrentMeanConfidenceMilli        int64
	ConfidenceDeltaMilli              int64
	BaselineOutcomeSamplesDigest      string
	CurrentOutcomeSamplesDigest       string
	BaselineChoiceSetDigest           string
	CurrentChoiceSetDigest            string
	BaselineChoiceOrderDigest         string
	CurrentChoiceOrderDigest          string
	SourceObservationDigest           string
	BaselineGenerationTraceDigest     string
	CurrentGenerationTraceDigest      string
	ReverseObservationCoverageDigest  string
	SensitivitySignal                 string
	DecisionSignal                    string
	ObservationDigest                 string
	ReadOnly                          bool
	NonExecuting                      bool
	NonAuthorizing                    bool
}

func ObserveRevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivity(
	input RevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivityInput,
) RevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivityObservation {
	result := RevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivityObservation{
		Status:           "UNKNOWN",
		MissingStage:     "revision-self-improvement-cycle-jev-calibration-choice-set-sensitivity",
		SensitivitySignal: jevCalibrationChoiceSetSensitivityUnknown,
		DecisionSignal:   jevCalibrationChoiceSetDecisionDefer,
		ReadOnly:         true,
		NonExecuting:     true,
		NonAuthorizing:   true,
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivity(result)
	}
	setDigest()

	if err := input.Baseline.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-choice-set-sensitivity-baseline"
		setDigest()
		return result
	}
	if input.Baseline.Status != "BOUND" {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-choice-set-sensitivity-baseline-status"
		setDigest()
		return result
	}
	if err := input.Current.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-choice-set-sensitivity-current"
		setDigest()
		return result
	}
	if input.Current.Status != "BOUND" {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-choice-set-sensitivity-current-status"
		setDigest()
		return result
	}
	for name, digest := range map[string]string{
		"baseline_choice_set":   input.BaselineChoiceSetDigest,
		"current_choice_set":    input.CurrentChoiceSetDigest,
		"baseline_choice_order": input.BaselineChoiceOrderDigest,
		"current_choice_order":  input.CurrentChoiceOrderDigest,
	} {
		if !validDigest(digest) {
			result.MissingStage = "revision-self-improvement-cycle-jev-calibration-choice-set-sensitivity-" + name
			setDigest()
			return result
		}
	}
	if input.Baseline.SourceObservationDigest != input.Current.SourceObservationDigest {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-choice-set-sensitivity-source-lineage"
		setDigest()
		return result
	}
	if input.Baseline.ReverseObservationCoverageDigest != input.Current.ReverseObservationCoverageDigest {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-choice-set-sensitivity-reverse-coverage-lineage"
		setDigest()
		return result
	}

	result.BaselineAccuracyMilli = input.Baseline.AccuracyMilli
	result.CurrentAccuracyMilli = input.Current.AccuracyMilli
	result.AccuracyDeltaMilli = result.CurrentAccuracyMilli - result.BaselineAccuracyMilli
	result.BaselineMeanConfidenceMilli = input.Baseline.MeanConfidenceMilli
	result.CurrentMeanConfidenceMilli = input.Current.MeanConfidenceMilli
	result.ConfidenceDeltaMilli = result.CurrentMeanConfidenceMilli - result.BaselineMeanConfidenceMilli
	result.BaselineOutcomeSamplesDigest = input.Baseline.OutcomeSamplesDigest
	result.CurrentOutcomeSamplesDigest = input.Current.OutcomeSamplesDigest
	result.BaselineChoiceSetDigest = input.BaselineChoiceSetDigest
	result.CurrentChoiceSetDigest = input.CurrentChoiceSetDigest
	result.BaselineChoiceOrderDigest = input.BaselineChoiceOrderDigest
	result.CurrentChoiceOrderDigest = input.CurrentChoiceOrderDigest
	result.SourceObservationDigest = input.Current.SourceObservationDigest
	result.BaselineGenerationTraceDigest = input.Baseline.GenerationTraceDigest
	result.CurrentGenerationTraceDigest = input.Current.GenerationTraceDigest
	result.ReverseObservationCoverageDigest = input.Current.ReverseObservationCoverageDigest

	choiceSetChanged := result.BaselineChoiceSetDigest != result.CurrentChoiceSetDigest
	choiceOrderChanged := result.BaselineChoiceOrderDigest != result.CurrentChoiceOrderDigest
	switch {
	case choiceSetChanged && choiceOrderChanged:
		result.SensitivitySignal = jevCalibrationChoiceSetSensitivityBothChanged
	case choiceSetChanged:
		result.SensitivitySignal = jevCalibrationChoiceSetSensitivityChanged
	case choiceOrderChanged:
		result.SensitivitySignal = jevCalibrationChoiceSetSensitivityOrderChanged
	default:
		result.SensitivitySignal = jevCalibrationChoiceSetSensitivityStable
	}
	if result.SensitivitySignal == jevCalibrationChoiceSetSensitivityStable {
		result.DecisionSignal = jevCalibrationChoiceSetDecisionCompare
	}
	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-choice-set-sensitivity"
		result.SensitivitySignal = jevCalibrationChoiceSetSensitivityUnknown
		result.DecisionSignal = jevCalibrationChoiceSetDecisionDefer
		setDigest()
	}
	return result
}

func (value RevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivityObservation) Validate() error {
	if value.Status != "BOUND" && value.Status != "UNKNOWN" {
		return fmt.Errorf("calibration choice-set sensitivity status is invalid")
	}
	if value.Status == "BOUND" && value.MissingStage != "" {
		return fmt.Errorf("bound calibration choice-set sensitivity has a missing stage")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" ||
			value.SensitivitySignal != jevCalibrationChoiceSetSensitivityUnknown ||
			value.DecisionSignal != jevCalibrationChoiceSetDecisionDefer {
			return fmt.Errorf("unknown calibration choice-set sensitivity must preserve incomplete evidence")
		}
	} else {
		for name, digest := range map[string]string{
			"baseline_outcome_samples":     value.BaselineOutcomeSamplesDigest,
			"current_outcome_samples":      value.CurrentOutcomeSamplesDigest,
			"baseline_choice_set":          value.BaselineChoiceSetDigest,
			"current_choice_set":           value.CurrentChoiceSetDigest,
			"baseline_choice_order":        value.BaselineChoiceOrderDigest,
			"current_choice_order":         value.CurrentChoiceOrderDigest,
			"source":                       value.SourceObservationDigest,
			"baseline_generation_trace":    value.BaselineGenerationTraceDigest,
			"current_generation_trace":     value.CurrentGenerationTraceDigest,
			"reverse_coverage":             value.ReverseObservationCoverageDigest,
		} {
			if !validDigest(digest) {
				return fmt.Errorf("bound calibration choice-set sensitivity %s digest is invalid", name)
			}
		}
		if value.BaselineAccuracyMilli < 0 ||
			value.BaselineAccuracyMilli > 1000 ||
			value.CurrentAccuracyMilli < 0 ||
			value.CurrentAccuracyMilli > 1000 ||
			value.AccuracyDeltaMilli != value.CurrentAccuracyMilli-value.BaselineAccuracyMilli ||
			value.BaselineMeanConfidenceMilli < 0 ||
			value.BaselineMeanConfidenceMilli > 1000 ||
			value.CurrentMeanConfidenceMilli < 0 ||
			value.CurrentMeanConfidenceMilli > 1000 ||
			value.ConfidenceDeltaMilli != value.CurrentMeanConfidenceMilli-value.BaselineMeanConfidenceMilli {
			return fmt.Errorf("bound calibration choice-set sensitivity metrics are invalid")
		}
		switch value.SensitivitySignal {
		case jevCalibrationChoiceSetSensitivityStable:
			if value.DecisionSignal != jevCalibrationChoiceSetDecisionCompare {
				return fmt.Errorf("stable calibration choice-set sensitivity must be comparable")
			}
		case jevCalibrationChoiceSetSensitivityChanged,
			jevCalibrationChoiceSetSensitivityOrderChanged,
			jevCalibrationChoiceSetSensitivityBothChanged:
			if value.DecisionSignal != jevCalibrationChoiceSetDecisionDefer {
				return fmt.Errorf("changed calibration choice-set sensitivity must defer")
			}
		default:
			return fmt.Errorf("bound calibration choice-set sensitivity signal is invalid")
		}
	}
	if !value.ReadOnly || !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("calibration choice-set sensitivity must remain read-only, non-executing, and non-authorizing")
	}
	if value.ObservationDigest != digestRevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivity(value) {
		return fmt.Errorf("calibration choice-set sensitivity digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivity(
	value RevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivityObservation,
) string {
	return digestString(strings.Join([]string{
		value.Status,
		value.MissingStage,
		strconv.FormatInt(value.BaselineAccuracyMilli, 10),
		strconv.FormatInt(value.CurrentAccuracyMilli, 10),
		strconv.FormatInt(value.AccuracyDeltaMilli, 10),
		strconv.FormatInt(value.BaselineMeanConfidenceMilli, 10),
		strconv.FormatInt(value.CurrentMeanConfidenceMilli, 10),
		strconv.FormatInt(value.ConfidenceDeltaMilli, 10),
		value.BaselineOutcomeSamplesDigest,
		value.CurrentOutcomeSamplesDigest,
		value.BaselineChoiceSetDigest,
		value.CurrentChoiceSetDigest,
		value.BaselineChoiceOrderDigest,
		value.CurrentChoiceOrderDigest,
		value.SourceObservationDigest,
		value.BaselineGenerationTraceDigest,
		value.CurrentGenerationTraceDigest,
		value.ReverseObservationCoverageDigest,
		value.SensitivitySignal,
		value.DecisionSignal,
		strconv.FormatBool(value.ReadOnly),
		strconv.FormatBool(value.NonExecuting),
		strconv.FormatBool(value.NonAuthorizing),
	}, "|"))
}