package gooo

import "testing"

func calibrationWindowInput() RevisionSelfImprovementCycleJEVCalibrationWindowInput {
	return RevisionSelfImprovementCycleJEVCalibrationWindowInput{
		MetricName: "jev-typed-decision-confidence-milli", QuestionKind: "choice",
		SampleCount: 10, CorrectCount: 8, ConfidenceSumMilli: 8200,
		CalibrationDataDigest: digestString("calibration:route"),
		GenerationDigest: digestString("generation:route"),
		ReverseObservationDigest: digestString("reverse:route"),
		SourceObservationDigest: digestString("source:route"),
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationWindow(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationWindow(calibrationWindowInput())
	if observation.Status != "BOUND" || observation.CalibrationStatus != "observed" ||
		observation.AccuracyMilli != 800 || observation.MeanConfidenceMilli != 820 ||
		observation.CalibrationErrorMilli != 20 {
		t.Fatalf("unexpected calibration window: %#v", observation)
	}
	if !observation.NonExecuting || !observation.NonAuthorizing {
		t.Fatalf("calibration window crossed an authority boundary: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("calibration window did not validate: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationWindowPreservesUnknown(t *testing.T) {
	input := calibrationWindowInput()
	input.SampleCount = 0
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationWindow(input)
	if observation.Status != "UNKNOWN" || observation.MissingStage == "" ||
		observation.CalibrationStatus != "unknown" {
		t.Fatalf("invalid calibration window must remain unknown: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("unknown calibration window should validate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVCalibrationWindowRejectsTampering(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationWindow(calibrationWindowInput())
	observation.CalibrationErrorMilli = 0
	if err := observation.Validate(); err == nil {
		t.Fatal("expected calibration error tampering to be rejected")
	}
}

func TestRevisionSelfImprovementCycleJEVCalibrationWindowRejectsMissingLineage(t *testing.T) {
	input := calibrationWindowInput()
	input.ReverseObservationDigest = ""
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationWindow(input)
	if observation.Status != "UNKNOWN" || observation.MissingStage == "" {
		t.Fatalf("missing reverse observation must remain unknown: %#v", observation)
	}
}