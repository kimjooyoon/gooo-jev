package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleJEVCalibrationClosureInput struct {
	CalibrationWindow                RevisionSelfImprovementCycleJEVCalibrationWindowObservation
	GenerationTraceDigest            string
	ReverseObservationCoverageDigest string
}

type RevisionSelfImprovementCycleJEVCalibrationClosureObservation struct {
	Status                         string
	MissingStage                   string
	MetricName                     string
	QuestionKind                   string
	SampleCount                    int64
	CorrectCount                   int64
	AccuracyMilli                  int64
	MeanConfidenceMilli            int64
	CalibrationErrorMilli          int64
	CalibrationStatus              string
	CalibrationDataDigest          string
	SourceObservationDigest        string
	GenerationDigest               string
	ReverseObservationDigest       string
	GenerationTraceDigest           string
	ReverseObservationCoverageDigest string
	ClosureSignal                  string
	ObservationDigest              string
	ReadOnly                       bool
	NonExecuting                   bool
	NonAuthorizing                 bool
}

func ObserveRevisionSelfImprovementCycleJEVCalibrationClosure(
	input RevisionSelfImprovementCycleJEVCalibrationClosureInput,
) RevisionSelfImprovementCycleJEVCalibrationClosureObservation {
	window := input.CalibrationWindow
	result := RevisionSelfImprovementCycleJEVCalibrationClosureObservation{
		Status: "UNKNOWN",
		MissingStage: "revision-self-improvement-cycle-jev-calibration-closure",
		MetricName: window.MetricName, QuestionKind: window.QuestionKind,
		SampleCount: window.SampleCount, CorrectCount: window.CorrectCount,
		AccuracyMilli: window.AccuracyMilli, MeanConfidenceMilli: window.MeanConfidenceMilli,
		CalibrationErrorMilli: window.CalibrationErrorMilli,
		CalibrationStatus: window.CalibrationStatus,
		CalibrationDataDigest: window.CalibrationDataDigest,
		SourceObservationDigest: window.SourceObservationDigest,
		GenerationDigest: window.GenerationDigest,
		ReverseObservationDigest: window.ReverseObservationDigest,
		GenerationTraceDigest: input.GenerationTraceDigest,
		ReverseObservationCoverageDigest: input.ReverseObservationCoverageDigest,
		ClosureSignal: "jev-calibration-closure-unknown",
		ReadOnly: true, NonExecuting: true, NonAuthorizing: true,
	}
	setDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCycleJEVCalibrationClosure(result)
	}
	setDigest()
	if err := window.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-closure-input-window"
		setDigest()
		return result
	}
	if window.Status != "BOUND" || window.CalibrationStatus != "observed" {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-closure-input-window-status"
		setDigest()
		return result
	}
	if !validDigest(input.GenerationTraceDigest) ||
		!validDigest(input.ReverseObservationCoverageDigest) {
		result.MissingStage = "revision-self-improvement-cycle-jev-calibration-closure-lineage"
		setDigest()
		return result
	}
	result.Status, result.MissingStage = "BOUND", ""
	result.ClosureSignal = "jev-calibration-closure-observed"
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status, result.MissingStage = "UNKNOWN", "revision-self-improvement-cycle-jev-calibration-closure"
		result.ClosureSignal = "jev-calibration-closure-unknown"
		setDigest()
	}
	return result
}

func (o RevisionSelfImprovementCycleJEVCalibrationClosureObservation) Validate() error {
	if o.Status != "BOUND" && o.Status != "UNKNOWN" {
		return fmt.Errorf("calibration closure status is invalid")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound calibration closure has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown calibration closure has no missing stage")
	}
	if o.Status == "UNKNOWN" {
		if o.ClosureSignal != "jev-calibration-closure-unknown" {
			return fmt.Errorf("unknown calibration closure must preserve incomplete lineage")
		}
	} else {
		if o.MetricName != "jev-typed-decision-confidence-milli" ||
			(o.QuestionKind != "choice" && o.QuestionKind != "score") ||
			o.CalibrationStatus != "observed" ||
			o.ClosureSignal != "jev-calibration-closure-observed" {
			return fmt.Errorf("bound calibration closure has invalid state")
		}
		if o.SampleCount < 1 || o.CorrectCount < 0 || o.CorrectCount > o.SampleCount ||
			o.AccuracyMilli < 0 || o.AccuracyMilli > 1000 ||
			o.MeanConfidenceMilli < 0 || o.MeanConfidenceMilli > 1000 ||
			o.CalibrationErrorMilli < 0 || o.CalibrationErrorMilli > 1000 {
			return fmt.Errorf("bound calibration closure metrics are out of range")
		}
		for name, value := range map[string]string{
			"calibration":          o.CalibrationDataDigest,
			"source":               o.SourceObservationDigest,
			"generation":           o.GenerationDigest,
			"reverse":              o.ReverseObservationDigest,
			"generation_trace":     o.GenerationTraceDigest,
			"reverse_coverage":     o.ReverseObservationCoverageDigest,
		} {
			if !validDigest(value) {
				return fmt.Errorf("bound calibration closure %s digest is invalid", name)
			}
		}
	}
	if !o.ReadOnly || !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("calibration closure must remain read-only, non-executing, and non-authorizing")
	}
	if o.ObservationDigest != digestRevisionSelfImprovementCycleJEVCalibrationClosure(o) {
		return fmt.Errorf("calibration closure digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleJEVCalibrationClosure(
	o RevisionSelfImprovementCycleJEVCalibrationClosureObservation,
) string {
	return digestString(strings.Join([]string{
		o.Status, o.MissingStage, o.MetricName, o.QuestionKind,
		strconv.FormatInt(o.SampleCount, 10),
		strconv.FormatInt(o.CorrectCount, 10),
		strconv.FormatInt(o.AccuracyMilli, 10),
		strconv.FormatInt(o.MeanConfidenceMilli, 10),
		strconv.FormatInt(o.CalibrationErrorMilli, 10),
		o.CalibrationStatus, o.CalibrationDataDigest, o.SourceObservationDigest,
		o.GenerationDigest, o.ReverseObservationDigest, o.GenerationTraceDigest,
		o.ReverseObservationCoverageDigest, o.ClosureSignal,
		strconv.FormatBool(o.ReadOnly), strconv.FormatBool(o.NonExecuting),
		strconv.FormatBool(o.NonAuthorizing),
	}, "|"))
}