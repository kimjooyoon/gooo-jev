package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	jevCalibrationOutcomeWindowObserved = "jev-calibration-outcome-window-observed"
	jevCalibrationOutcomeWindowUnknown  = "jev-calibration-outcome-window-unknown"
)

type RevisionSelfImprovementCycleJEVCalibrationOutcomeWindowInput struct {
	Observations                  []RevisionSelfImprovementCycleJEVCalibrationOutcomeObservation
	SourceObservationDigest       string
	GenerationTraceDigest         string
	ReverseObservationCoverageDigest string
}

type RevisionSelfImprovementCycleJEVCalibrationOutcomeWindowObservation struct {
	Status                         string
	MissingStage                   string
	SampleCount                    int64
	CorrectCount                   int64
	IncorrectCount                 int64
	AbstainedCount                 int64
	EvaluatedCount                 int64
	AccuracyMilli                  int64
	MeanConfidenceMilli            int64
	OutcomeSamplesDigest            string
	SourceObservationDigest        string
	GenerationTraceDigest          string
	ReverseObservationCoverageDigest string
	WindowSignal                   string
	ObservationDigest              string
	ReadOnly                       bool
	NonExecuting                   bool
	NonAuthorizing                 bool
}

func ObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeWindow(
	input RevisionSelfImprovementCycleJEVCalibrationOutcomeWindowInput,
) RevisionSelfImprovementCycleJEVCalibrationOutcomeWindowObservation {
	result := RevisionSelfImprovementCycleJEVCalibrationOutcomeWindowObservation{
		Status:                           "UNKNOWN",
		MissingStage:                     "revision-self-improvement-cycle-jev-calibration-outcome-window",
		SourceObservationDigest:          input.SourceObservationDigest,
		GenerationTraceDigest:            input.GenerationTraceDigest,
		ReverseObservationCoverageDigest: input.ReverseObservationCoverageDigest,
		WindowSignal:                     jevCalibrationOutcomeWindowUnknown,
		ReadOnly:                         true,
		NonExecuting:                     true,
		NonAuthorizing:                   true,
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCycleJEVCalibrationOutcomeWindow(result)
	}
	setDigest()

	if len(input.Observations) == 0 {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-outcome-window-samples"
		setDigest()
		return result
	}
	if !validDigest(input.SourceObservationDigest) ||
		!validDigest(input.GenerationTraceDigest) ||
		!validDigest(input.ReverseObservationCoverageDigest) {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-outcome-window-lineage"
		setDigest()
		return result
	}

	sampleDigests := make([]string, 0, len(input.Observations))
	var confidenceSum int64
	for index, sample := range input.Observations {
		if err := sample.Validate(); err != nil ||
			sample.Status != "BOUND" ||
			sample.OutcomeSignal != jevCalibrationOutcomeObserved {
			result.MissingStage = fmt.Sprintf(
				"revision-self-improvement-cycle-jev-calibration-outcome-window-sample-%d",
				index,
			)
			setDigest()
			return result
		}
		sampleDigests = append(sampleDigests, sample.ObservationDigest)
		confidenceSum += sample.ConfidenceMilli
		switch sample.OutcomeClass {
		case jevCalibrationOutcomeCorrect:
			result.CorrectCount++
		case jevCalibrationOutcomeIncorrect:
			result.IncorrectCount++
		case jevCalibrationOutcomeAbstained:
			result.AbstainedCount++
		default:
			result.MissingStage = fmt.Sprintf(
				"revision-self-improvement-cycle-jev-calibration-outcome-window-sample-%d-class",
				index,
			)
			setDigest()
			return result
		}
	}
	result.SampleCount = int64(len(input.Observations))
	result.EvaluatedCount = result.CorrectCount + result.IncorrectCount
	if result.EvaluatedCount < 1 {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-outcome-window-evaluated-samples"
		setDigest()
		return result
	}
	result.OutcomeSamplesDigest = digestString(strings.Join(sampleDigests, "|"))
	result.AccuracyMilli = result.CorrectCount * 1000 / result.EvaluatedCount
	result.MeanConfidenceMilli = confidenceSum / result.SampleCount
	result.Status = "BOUND"
	result.MissingStage = ""
	result.WindowSignal = jevCalibrationOutcomeWindowObserved
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-outcome-window"
		result.WindowSignal = jevCalibrationOutcomeWindowUnknown
		setDigest()
	}
	return result
}

func (value RevisionSelfImprovementCycleJEVCalibrationOutcomeWindowObservation) Validate() error {
	if value.Status != "BOUND" && value.Status != "UNKNOWN" {
		return fmt.Errorf("calibration outcome window status is invalid")
	}
	if value.Status == "BOUND" && value.MissingStage != "" {
		return fmt.Errorf("bound calibration outcome window has a missing stage")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" ||
			value.WindowSignal != jevCalibrationOutcomeWindowUnknown {
			return fmt.Errorf("unknown calibration outcome window must preserve incomplete evidence")
		}
	} else {
		for name, digest := range map[string]string{
			"outcome_samples":    value.OutcomeSamplesDigest,
			"source":             value.SourceObservationDigest,
			"generation_trace":   value.GenerationTraceDigest,
			"reverse_coverage":   value.ReverseObservationCoverageDigest,
		} {
			if !validDigest(digest) {
				return fmt.Errorf("bound calibration outcome window %s digest is invalid", name)
			}
		}
		if value.SampleCount < 1 ||
			value.CorrectCount < 0 ||
			value.IncorrectCount < 0 ||
			value.AbstainedCount < 0 ||
			value.EvaluatedCount < 1 ||
			value.CorrectCount+value.IncorrectCount != value.EvaluatedCount ||
			value.EvaluatedCount+value.AbstainedCount != value.SampleCount ||
			value.AccuracyMilli < 0 ||
			value.AccuracyMilli > 1000 ||
			value.MeanConfidenceMilli < 0 ||
			value.MeanConfidenceMilli > 1000 ||
			value.WindowSignal != jevCalibrationOutcomeWindowObserved {
			return fmt.Errorf("bound calibration outcome window metrics are invalid")
		}
	}
	if !value.ReadOnly || !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("calibration outcome window must remain read-only, non-executing, and non-authorizing")
	}
	if value.ObservationDigest != digestRevisionSelfImprovementCycleJEVCalibrationOutcomeWindow(value) {
		return fmt.Errorf("calibration outcome window digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleJEVCalibrationOutcomeWindow(
	value RevisionSelfImprovementCycleJEVCalibrationOutcomeWindowObservation,
) string {
	return digestString(strings.Join([]string{
		value.Status,
		value.MissingStage,
		strconv.FormatInt(value.SampleCount, 10),
		strconv.FormatInt(value.CorrectCount, 10),
		strconv.FormatInt(value.IncorrectCount, 10),
		strconv.FormatInt(value.AbstainedCount, 10),
		strconv.FormatInt(value.EvaluatedCount, 10),
		strconv.FormatInt(value.AccuracyMilli, 10),
		strconv.FormatInt(value.MeanConfidenceMilli, 10),
		value.OutcomeSamplesDigest,
		value.SourceObservationDigest,
		value.GenerationTraceDigest,
		value.ReverseObservationCoverageDigest,
		value.WindowSignal,
		strconv.FormatBool(value.ReadOnly),
		strconv.FormatBool(value.NonExecuting),
		strconv.FormatBool(value.NonAuthorizing),
	}, "|"))
}
