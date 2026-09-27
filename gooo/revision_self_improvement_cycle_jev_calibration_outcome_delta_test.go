package gooo

import "testing"

func jevCalibrationOutcomeDeltaInput() RevisionSelfImprovementCycleJEVCalibrationOutcomeDeltaInput {
	baselineInput := jevCalibrationOutcomeWindowInput()
	currentInput := jevCalibrationOutcomeWindowInput()
	currentInput.GenerationTraceDigest = digestString("generation-trace:revision-2")
	currentInput.Observations = []RevisionSelfImprovementCycleJEVCalibrationOutcomeObservation{
		jevCalibrationOutcomeWindowSample("route-a", "route-a", 900, false),
		jevCalibrationOutcomeWindowSample("route-a", "route-a", 700, false),
		jevCalibrationOutcomeWindowSample("route-a", "route-b", 600, true),
	}
	return RevisionSelfImprovementCycleJEVCalibrationOutcomeDeltaInput{
		Baseline: ObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeWindow(baselineInput),
		Current:  ObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeWindow(currentInput),
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeDelta(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeDelta(
		jevCalibrationOutcomeDeltaInput(),
	)
	if observation.Status != "BOUND" ||
		observation.AccuracyDeltaMilli != 500 ||
		observation.ConfidenceDeltaMilli != 100 ||
		observation.DeltaSignal != jevCalibrationOutcomeDeltaImproved {
		t.Fatalf("unexpected calibration outcome delta: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("calibration outcome delta did not validate: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeDeltaRejectsLineageMismatch(t *testing.T) {
	input := jevCalibrationOutcomeDeltaInput()
	currentInput := jevCalibrationOutcomeWindowInput()
	currentInput.SourceObservationDigest = digestString("source:other")
	input.Current = ObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeWindow(currentInput)
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeDelta(input)
	if observation.Status != "UNKNOWN" ||
		observation.MissingStage != "revision-self-improvement-cycle-jev-calibration-outcome-delta-source-lineage" {
		t.Fatalf("lineage mismatch must remain unknown: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("unknown lineage mismatch should validate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVCalibrationOutcomeDeltaRejectsTampering(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeDelta(
		jevCalibrationOutcomeDeltaInput(),
	)
	observation.ConfidenceDeltaMilli = 900
	if err := observation.Validate(); err == nil {
		t.Fatal("expected calibration outcome delta tampering to be rejected")
	}
}