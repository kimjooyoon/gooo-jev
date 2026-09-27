package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	jevCalibrationOutcomeDeltaImproved = "jev-calibration-outcome-delta-improved"
	jevCalibrationOutcomeDeltaRegressed = "jev-calibration-outcome-delta-regressed"
	jevCalibrationOutcomeDeltaFlat = "jev-calibration-outcome-delta-flat"
	jevCalibrationOutcomeDeltaUnknown = "jev-calibration-outcome-delta-unknown"
)

type RevisionSelfImprovementCycleJEVCalibrationOutcomeDeltaInput struct {
	Baseline RevisionSelfImprovementCycleJEVCalibrationOutcomeWindowObservation
	Current  RevisionSelfImprovementCycleJEVCalibrationOutcomeWindowObservation
}

type RevisionSelfImprovementCycleJEVCalibrationOutcomeDeltaObservation struct {
	Status                         string
	MissingStage                   string
	BaselineAccuracyMilli          int64
	CurrentAccuracyMilli           int64
	AccuracyDeltaMilli             int64
	BaselineMeanConfidenceMilli    int64
	CurrentMeanConfidenceMilli     int64
	ConfidenceDeltaMilli           int64
	BaselineOutcomeSamplesDigest   string
	CurrentOutcomeSamplesDigest    string
	SourceObservationDigest        string
	BaselineGenerationTraceDigest  string
	CurrentGenerationTraceDigest   string
	ReverseObservationCoverageDigest string
	DeltaSignal                    string
	ObservationDigest              string
	ReadOnly                       bool
	NonExecuting                   bool
	NonAuthorizing                 bool
}

func ObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeDelta(
	input RevisionSelfImprovementCycleJEVCalibrationOutcomeDeltaInput,
) RevisionSelfImprovementCycleJEVCalibrationOutcomeDeltaObservation {
	result := RevisionSelfImprovementCycleJEVCalibrationOutcomeDeltaObservation{
		Status:           "UNKNOWN",
		MissingStage:     "revision-self-improvement-cycle-jev-calibration-outcome-delta",
		DeltaSignal:      jevCalibrationOutcomeDeltaUnknown,
		ReadOnly:         true,
		NonExecuting:     true,
		NonAuthorizing:   true,
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCycleJEVCalibrationOutcomeDelta(result)
	}
	setDigest()

	if err := input.Baseline.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-outcome-delta-baseline"
		setDigest()
		return result
	}
	if input.Baseline.Status != "BOUND" {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-outcome-delta-baseline-status"
		setDigest()
		return result
	}
	if err := input.Current.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-outcome-delta-current"
		setDigest()
		return result
	}
	if input.Current.Status != "BOUND" {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-outcome-delta-current-status"
		setDigest()
		return result
	}
	if input.Baseline.SourceObservationDigest != input.Current.SourceObservationDigest {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-outcome-delta-source-lineage"
		setDigest()
		return result
	}
	if input.Baseline.ReverseObservationCoverageDigest != input.Current.ReverseObservationCoverageDigest {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-outcome-delta-reverse-coverage-lineage"
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
	result.SourceObservationDigest = input.Current.SourceObservationDigest
	result.BaselineGenerationTraceDigest = input.Baseline.GenerationTraceDigest
	result.CurrentGenerationTraceDigest = input.Current.GenerationTraceDigest
	result.ReverseObservationCoverageDigest = input.Current.ReverseObservationCoverageDigest

	switch {
	case result.AccuracyDeltaMilli > 0 ||
		(result.AccuracyDeltaMilli == 0 && result.ConfidenceDeltaMilli > 0):
		result.DeltaSignal = jevCalibrationOutcomeDeltaImproved
	case result.AccuracyDeltaMilli < 0 ||
		(result.AccuracyDeltaMilli == 0 && result.ConfidenceDeltaMilli < 0):
		result.DeltaSignal = jevCalibrationOutcomeDeltaRegressed
	default:
		result.DeltaSignal = jevCalibrationOutcomeDeltaFlat
	}
	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-outcome-delta"
		result.DeltaSignal = jevCalibrationOutcomeDeltaUnknown
		setDigest()
	}
	return result
}

func (value RevisionSelfImprovementCycleJEVCalibrationOutcomeDeltaObservation) Validate() error {
	if value.Status != "BOUND" && value.Status != "UNKNOWN" {
		return fmt.Errorf("calibration outcome delta status is invalid")
	}
	if value.Status == "BOUND" && value.MissingStage != "" {
		return fmt.Errorf("bound calibration outcome delta has a missing stage")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" || value.DeltaSignal != jevCalibrationOutcomeDeltaUnknown {
			return fmt.Errorf("unknown calibration outcome delta must preserve incomplete evidence")
		}
	} else {
		for name, digest := range map[string]string{
			"baseline_outcome_samples":    value.BaselineOutcomeSamplesDigest,
			"current_outcome_samples":     value.CurrentOutcomeSamplesDigest,
			"source":                      value.SourceObservationDigest,
			"baseline_generation_trace":   value.BaselineGenerationTraceDigest,
			"current_generation_trace":    value.CurrentGenerationTraceDigest,
			"reverse_coverage":            value.ReverseObservationCoverageDigest,
		} {
			if !validDigest(digest) {
				return fmt.Errorf("bound calibration outcome delta %s digest is invalid", name)
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
			return fmt.Errorf("bound calibration outcome delta metrics are invalid")
		}
		switch value.DeltaSignal {
		case jevCalibrationOutcomeDeltaImproved,
			jevCalibrationOutcomeDeltaRegressed,
			jevCalibrationOutcomeDeltaFlat:
		default:
			return fmt.Errorf("bound calibration outcome delta signal is invalid")
		}
	}
	if !value.ReadOnly || !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("calibration outcome delta must remain read-only, non-executing, and non-authorizing")
	}
	if value.ObservationDigest != digestRevisionSelfImprovementCycleJEVCalibrationOutcomeDelta(value) {
		return fmt.Errorf("calibration outcome delta digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleJEVCalibrationOutcomeDelta(
	value RevisionSelfImprovementCycleJEVCalibrationOutcomeDeltaObservation,
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
		value.SourceObservationDigest,
		value.BaselineGenerationTraceDigest,
		value.CurrentGenerationTraceDigest,
		value.ReverseObservationCoverageDigest,
		value.DeltaSignal,
		strconv.FormatBool(value.ReadOnly),
		strconv.FormatBool(value.NonExecuting),
		strconv.FormatBool(value.NonAuthorizing),
	}, "|"))
}