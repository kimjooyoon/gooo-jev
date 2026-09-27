package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

const jevCalibrationWindowMaxSamples int64 = 1000000

type RevisionSelfImprovementCycleJEVCalibrationWindowInput struct {
	MetricName              string
	QuestionKind            string
	SampleCount             int64
	CorrectCount            int64
	ConfidenceSumMilli      int64
	CalibrationDataDigest   string
	GenerationDigest        string
	ReverseObservationDigest string
	SourceObservationDigest string
}

type RevisionSelfImprovementCycleJEVCalibrationWindowObservation struct {
	Status                    string
	MissingStage              string
	MetricName                string
	QuestionKind              string
	SampleCount               int64
	CorrectCount               int64
	ConfidenceSumMilli         int64
	AccuracyMilli              int64
	MeanConfidenceMilli        int64
	CalibrationErrorMilli      int64
	CalibrationStatus           string
	CalibrationDataDigest       string
	GenerationDigest            string
	ReverseObservationDigest    string
	SourceObservationDigest    string
	CalibrationSignal           string
	ObservationDigest           string
	ReadOnly                    bool
	NonExecuting                bool
	NonAuthorizing              bool
}

func ObserveRevisionSelfImprovementCycleJEVCalibrationWindow(
	input RevisionSelfImprovementCycleJEVCalibrationWindowInput,
) RevisionSelfImprovementCycleJEVCalibrationWindowObservation {
	result := RevisionSelfImprovementCycleJEVCalibrationWindowObservation{
		Status: "UNKNOWN",
		MissingStage: "revision-self-improvement-cycle-jev-calibration-window",
		MetricName: input.MetricName, QuestionKind: input.QuestionKind,
		SampleCount: input.SampleCount, CorrectCount: input.CorrectCount,
		ConfidenceSumMilli: input.ConfidenceSumMilli,
		CalibrationStatus: "unknown",
		CalibrationDataDigest: input.CalibrationDataDigest,
		GenerationDigest: input.GenerationDigest,
		ReverseObservationDigest: input.ReverseObservationDigest,
		SourceObservationDigest: input.SourceObservationDigest,
		CalibrationSignal: "jev-calibration-window-unknown",
		ReadOnly: true, NonExecuting: true, NonAuthorizing: true,
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCycleJEVCalibrationWindow(result)
	}
	setDigest()
	if err := input.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-window-input"
		setDigest()
		return result
	}

	result.Status, result.MissingStage = "BOUND", ""
	result.AccuracyMilli = (input.CorrectCount * 1000) / input.SampleCount
	result.MeanConfidenceMilli = input.ConfidenceSumMilli / input.SampleCount
	result.CalibrationErrorMilli = result.MeanConfidenceMilli - result.AccuracyMilli
	if result.CalibrationErrorMilli < 0 {
		result.CalibrationErrorMilli = -result.CalibrationErrorMilli
	}
	result.CalibrationStatus = "observed"
	result.CalibrationSignal = "jev-calibration-window-observed"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status, result.MissingStage = "UNKNOWN", "revision-self-improvement-cycle-jev-calibration-window"
		result.CalibrationStatus = "unknown"
		result.CalibrationSignal = "jev-calibration-window-unknown"
		result.AccuracyMilli, result.MeanConfidenceMilli, result.CalibrationErrorMilli = 0, 0, 0
		setDigest()
	}
	return result
}

func (i RevisionSelfImprovementCycleJEVCalibrationWindowInput) Validate() error {
	if i.MetricName != "jev-typed-decision-confidence-milli" {
		return fmt.Errorf("calibration window metric name is invalid")
	}
	if i.QuestionKind != "choice" && i.QuestionKind != "score" {
		return fmt.Errorf("calibration window question kind is invalid")
	}
	if i.SampleCount < 1 || i.SampleCount > jevCalibrationWindowMaxSamples {
		return fmt.Errorf("calibration window sample count is out of range")
	}
	if i.CorrectCount < 0 || i.CorrectCount > i.SampleCount {
		return fmt.Errorf("calibration window correct count is out of range")
	}
	if i.ConfidenceSumMilli < 0 || i.ConfidenceSumMilli > i.SampleCount*1000 {
		return fmt.Errorf("calibration window confidence sum is out of range")
	}
	for name, value := range map[string]string{
		"calibration": i.CalibrationDataDigest,
		"generation": i.GenerationDigest,
		"reverse": i.ReverseObservationDigest,
		"source": i.SourceObservationDigest,
	} {
		if !validDigest(value) {
			return fmt.Errorf("calibration window %s digest is invalid", name)
		}
	}
	return nil
}

func (o RevisionSelfImprovementCycleJEVCalibrationWindowObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("calibration window status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound calibration window has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown calibration window has no missing stage")
	}
	if o.Status == "UNKNOWN" {
		if o.CalibrationStatus != "unknown" ||
			o.CalibrationSignal != "jev-calibration-window-unknown" {
			return fmt.Errorf("unknown calibration window must preserve missing outcome evidence")
		}
	} else {
		input := RevisionSelfImprovementCycleJEVCalibrationWindowInput{
			MetricName: o.MetricName, QuestionKind: o.QuestionKind,
			SampleCount: o.SampleCount, CorrectCount: o.CorrectCount,
			ConfidenceSumMilli: o.ConfidenceSumMilli,
			CalibrationDataDigest: o.CalibrationDataDigest,
			GenerationDigest: o.GenerationDigest,
			ReverseObservationDigest: o.ReverseObservationDigest,
			SourceObservationDigest: o.SourceObservationDigest,
		}
		if err := input.Validate(); err != nil {
			return fmt.Errorf("bound calibration window input is invalid: %w", err)
		}
		wantAccuracy := (o.CorrectCount * 1000) / o.SampleCount
		wantMean := o.ConfidenceSumMilli / o.SampleCount
		wantError := wantMean - wantAccuracy
		if wantError < 0 {
			wantError = -wantError
		}
		if o.AccuracyMilli != wantAccuracy || o.MeanConfidenceMilli != wantMean ||
			o.CalibrationErrorMilli != wantError ||
			o.CalibrationStatus != "observed" ||
			o.CalibrationSignal != "jev-calibration-window-observed" {
			return fmt.Errorf("bound calibration window metrics are inconsistent")
		}
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("calibration window must remain read-only, non-executing, and non-authorizing")
	}
	if o.ObservationDigest != digestRevisionSelfImprovementCycleJEVCalibrationWindow(o) {
		return fmt.Errorf("calibration window digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleJEVCalibrationWindow(
	o RevisionSelfImprovementCycleJEVCalibrationWindowObservation,
) string {
	return digestString(strings.Join([]string{
		o.Status, o.MissingStage, o.MetricName, o.QuestionKind,
		strconv.FormatInt(o.SampleCount, 10),
		strconv.FormatInt(o.CorrectCount, 10),
		strconv.FormatInt(o.ConfidenceSumMilli, 10),
		strconv.FormatInt(o.AccuracyMilli, 10),
		strconv.FormatInt(o.MeanConfidenceMilli, 10),
		strconv.FormatInt(o.CalibrationErrorMilli, 10),
		o.CalibrationStatus, o.CalibrationDataDigest, o.GenerationDigest,
		o.ReverseObservationDigest, o.SourceObservationDigest, o.CalibrationSignal,
		strconv.FormatBool(o.ReadOnly), strconv.FormatBool(o.NonExecuting),
		strconv.FormatBool(o.NonAuthorizing),
	}, "|"))
}