package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	jevCalibrationOutcomeObserved = "jev-calibration-outcome-observed"
	jevCalibrationOutcomeUnknown  = "jev-calibration-outcome-unknown"
	jevCalibrationOutcomeCorrect  = "correct"
	jevCalibrationOutcomeIncorrect = "incorrect"
	jevCalibrationOutcomeAbstained = "abstained"
	jevCalibrationOutcomeUnknownClass = "unknown"
)

type RevisionSelfImprovementCycleJEVCalibrationOutcomeInput struct {
	DecisionDigest          string
	DeclarationDigest       string
	IRDigest                string
	GenerationDigest        string
	ReverseObservationDigest string
	PredictedLabel          string
	ObservedLabel           string
	ConfidenceMilli         int64
	Abstained               bool
	OutcomeStatus           string
}

type RevisionSelfImprovementCycleJEVCalibrationOutcomeObservation struct {
	Status                   string
	MissingStage             string
	DecisionDigest           string
	DeclarationDigest        string
	IRDigest                 string
	GenerationDigest         string
	ReverseObservationDigest string
	PredictedLabel           string
	ObservedLabel            string
	ConfidenceMilli          int64
	Abstained                bool
	OutcomeStatus             string
	OutcomeClass             string
	OutcomeSignal            string
	ObservationDigest        string
	ReadOnly                 bool
	NonExecuting             bool
	NonAuthorizing           bool
}

func ObserveRevisionSelfImprovementCycleJEVCalibrationOutcome(
	input RevisionSelfImprovementCycleJEVCalibrationOutcomeInput,
) RevisionSelfImprovementCycleJEVCalibrationOutcomeObservation {
	result := RevisionSelfImprovementCycleJEVCalibrationOutcomeObservation{
		Status:                    "UNKNOWN",
		MissingStage:              "revision-self-improvement-cycle-jev-calibration-outcome",
		DecisionDigest:             input.DecisionDigest,
		DeclarationDigest:          input.DeclarationDigest,
		IRDigest:                   input.IRDigest,
		GenerationDigest:           input.GenerationDigest,
		ReverseObservationDigest:  input.ReverseObservationDigest,
		PredictedLabel:             input.PredictedLabel,
		ObservedLabel:              input.ObservedLabel,
		ConfidenceMilli:            input.ConfidenceMilli,
		Abstained:                  input.Abstained,
		OutcomeStatus:              input.OutcomeStatus,
		OutcomeClass:               jevCalibrationOutcomeUnknownClass,
		OutcomeSignal:              jevCalibrationOutcomeUnknown,
		ReadOnly:                   true,
		NonExecuting:               true,
		NonAuthorizing:             true,
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCycleJEVCalibrationOutcome(result)
	}
	setDigest()

	if !validDigest(input.DecisionDigest) {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-outcome-decision"
		setDigest()
		return result
	}
	if !validDigest(input.DeclarationDigest) ||
		!validDigest(input.IRDigest) ||
		!validDigest(input.GenerationDigest) ||
		!validDigest(input.ReverseObservationDigest) {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-outcome-lineage"
		setDigest()
		return result
	}
	if strings.TrimSpace(input.PredictedLabel) == "" ||
		strings.TrimSpace(input.ObservedLabel) == "" {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-outcome-labels"
		setDigest()
		return result
	}
	if input.ConfidenceMilli < 0 || input.ConfidenceMilli > 1000 {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-outcome-confidence"
		setDigest()
		return result
	}
	if input.OutcomeStatus != "observed" {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-outcome-observation"
		setDigest()
		return result
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	if input.Abstained {
		result.OutcomeClass = jevCalibrationOutcomeAbstained
	} else if input.PredictedLabel == input.ObservedLabel {
		result.OutcomeClass = jevCalibrationOutcomeCorrect
	} else {
		result.OutcomeClass = jevCalibrationOutcomeIncorrect
	}
	result.OutcomeSignal = jevCalibrationOutcomeObserved
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-outcome"
		result.OutcomeClass = jevCalibrationOutcomeUnknownClass
		result.OutcomeSignal = jevCalibrationOutcomeUnknown
		setDigest()
	}
	return result
}

func (value RevisionSelfImprovementCycleJEVCalibrationOutcomeObservation) Validate() error {
	if value.Status != "BOUND" && value.Status != "UNKNOWN" {
		return fmt.Errorf("calibration outcome status is invalid")
	}
	if value.Status == "BOUND" && value.MissingStage != "" {
		return fmt.Errorf("bound calibration outcome has a missing stage")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" ||
			value.OutcomeClass != jevCalibrationOutcomeUnknownClass ||
			value.OutcomeSignal != jevCalibrationOutcomeUnknown {
			return fmt.Errorf("unknown calibration outcome must preserve incomplete evidence")
		}
	} else {
		for name, digest := range map[string]string{
			"decision":  value.DecisionDigest,
			"declaration": value.DeclarationDigest,
			"ir":         value.IRDigest,
			"generation": value.GenerationDigest,
			"reverse":   value.ReverseObservationDigest,
		} {
			if !validDigest(digest) {
				return fmt.Errorf("bound calibration outcome %s digest is invalid", name)
			}
		}
		if strings.TrimSpace(value.PredictedLabel) == "" ||
			strings.TrimSpace(value.ObservedLabel) == "" ||
			value.ConfidenceMilli < 0 ||
			value.ConfidenceMilli > 1000 ||
			value.OutcomeStatus != "observed" ||
			value.OutcomeSignal != jevCalibrationOutcomeObserved {
			return fmt.Errorf("bound calibration outcome evidence is incomplete")
		}
		switch value.OutcomeClass {
		case jevCalibrationOutcomeCorrect:
			if value.Abstained || value.PredictedLabel != value.ObservedLabel {
				return fmt.Errorf("correct calibration outcome is inconsistent")
			}
		case jevCalibrationOutcomeIncorrect:
			if value.Abstained || value.PredictedLabel == value.ObservedLabel {
				return fmt.Errorf("incorrect calibration outcome is inconsistent")
			}
		case jevCalibrationOutcomeAbstained:
			if !value.Abstained {
				return fmt.Errorf("abstained calibration outcome is inconsistent")
			}
		default:
			return fmt.Errorf("bound calibration outcome class is invalid")
		}
	}
	if !value.ReadOnly || !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("calibration outcome must remain read-only, non-executing, and non-authorizing")
	}
	if value.ObservationDigest != digestRevisionSelfImprovementCycleJEVCalibrationOutcome(value) {
		return fmt.Errorf("calibration outcome digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleJEVCalibrationOutcome(
	value RevisionSelfImprovementCycleJEVCalibrationOutcomeObservation,
) string {
	return digestString(strings.Join([]string{
		value.Status,
		value.MissingStage,
		value.DecisionDigest,
		value.DeclarationDigest,
		value.IRDigest,
		value.GenerationDigest,
		value.ReverseObservationDigest,
		value.PredictedLabel,
		value.ObservedLabel,
		strconv.FormatInt(value.ConfidenceMilli, 10),
		strconv.FormatBool(value.Abstained),
		value.OutcomeStatus,
		value.OutcomeClass,
		value.OutcomeSignal,
		strconv.FormatBool(value.ReadOnly),
		strconv.FormatBool(value.NonExecuting),
		strconv.FormatBool(value.NonAuthorizing),
	}, "|"))
}
