package gooo

import "testing"

func closureWindow() RevisionSelfImprovementCycleJEVCalibrationWindowObservation {
	return ObserveRevisionSelfImprovementCycleJEVCalibrationWindow(
		RevisionSelfImprovementCycleJEVCalibrationWindowInput{
			MetricName: "jev-typed-decision-confidence-milli", QuestionKind: "choice",
			SampleCount: 10, CorrectCount: 8, ConfidenceSumMilli: 8200,
			CalibrationDataDigest: digestString("calibration:route"),
			GenerationDigest: digestString("generation:route"),
			ReverseObservationDigest: digestString("reverse:route"),
			SourceObservationDigest: digestString("source:route"),
		},
	)
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationClosure(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationClosure(
		RevisionSelfImprovementCycleJEVCalibrationClosureInput{
			CalibrationWindow: closureWindow(),
			GenerationTraceDigest: digestString("generation-trace:route"),
			ReverseObservationCoverageDigest: digestString("reverse-coverage:route"),
		},
	)
	if observation.Status != "BOUND" || observation.ClosureSignal != "jev-calibration-closure-observed" ||
		observation.CalibrationErrorMilli != 20 {
		t.Fatalf("unexpected calibration closure: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("calibration closure did not validate: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationClosurePreservesUnknown(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationClosure(
		RevisionSelfImprovementCycleJEVCalibrationClosureInput{},
	)
	if observation.Status != "UNKNOWN" || observation.MissingStage == "" ||
		observation.ClosureSignal != "jev-calibration-closure-unknown" {
		t.Fatalf("incomplete closure must remain unknown: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("unknown closure should validate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVCalibrationClosureRejectsTampering(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationClosure(
		RevisionSelfImprovementCycleJEVCalibrationClosureInput{
			CalibrationWindow: closureWindow(),
			GenerationTraceDigest: digestString("generation-trace:route"),
			ReverseObservationCoverageDigest: digestString("reverse-coverage:route"),
		},
	)
	observation.ClosureSignal = "execute"
	if err := observation.Validate(); err == nil {
		t.Fatal("expected authority-crossing closure signal to be rejected")
	}
}