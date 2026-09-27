package gooo

import "testing"

func jevCalibrationOutcomeWindowSample(
	predictedLabel string,
	observedLabel string,
	confidenceMilli int64,
	abstained bool,
) RevisionSelfImprovementCycleJEVCalibrationOutcomeObservation {
	input := jevCalibrationOutcomeInput()
	input.PredictedLabel = predictedLabel
	input.ObservedLabel = observedLabel
	input.ConfidenceMilli = confidenceMilli
	input.Abstained = abstained
	return ObserveRevisionSelfImprovementCycleJEVCalibrationOutcome(input)
}

func jevCalibrationOutcomeWindowInput() RevisionSelfImprovementCycleJEVCalibrationOutcomeWindowInput {
	return RevisionSelfImprovementCycleJEVCalibrationOutcomeWindowInput{
		Observations: []RevisionSelfImprovementCycleJEVCalibrationOutcomeObservation{
			jevCalibrationOutcomeWindowSample("route-a", "route-a", 800, false),
			jevCalibrationOutcomeWindowSample("route-a", "route-b", 600, false),
			jevCalibrationOutcomeWindowSample("route-a", "route-b", 500, true),
		},
		SourceObservationDigest:          digestString("source:window"),
		GenerationTraceDigest:            digestString("generation-trace:window"),
		ReverseObservationCoverageDigest: digestString("reverse-coverage:window"),
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeWindow(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeWindow(
		jevCalibrationOutcomeWindowInput(),
	)
	if observation.Status != "BOUND" ||
		observation.SampleCount != 3 ||
		observation.CorrectCount != 1 ||
		observation.IncorrectCount != 1 ||
		observation.AbstainedCount != 1 ||
		observation.EvaluatedCount != 2 ||
		observation.AccuracyMilli != 500 ||
		observation.MeanConfidenceMilli != 633 {
		t.Fatalf("unexpected calibration outcome window: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("calibration outcome window did not validate: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeWindowPreservesUnknown(t *testing.T) {
	input := jevCalibrationOutcomeWindowInput()
	input.Observations = nil
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeWindow(input)
	if observation.Status != "UNKNOWN" ||
		observation.MissingStage != "revision-self-improvement-cycle-jev-calibration-outcome-window-samples" ||
		observation.WindowSignal != jevCalibrationOutcomeWindowUnknown {
		t.Fatalf("empty outcome window must remain unknown: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("unknown outcome window should validate: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeWindowRejectsUnknownSample(t *testing.T) {
	input := jevCalibrationOutcomeWindowInput()
	input.Observations[1].Status = "UNKNOWN"
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeWindow(input)
	if observation.Status != "UNKNOWN" ||
		observation.MissingStage != "revision-self-improvement-cycle-jev-calibration-outcome-window-sample-1" {
		t.Fatalf("unknown sample must not enter a bound window: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("unknown sample window should validate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVCalibrationOutcomeWindowRejectsTampering(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeWindow(
		jevCalibrationOutcomeWindowInput(),
	)
	observation.AccuracyMilli = 1000
	if err := observation.Validate(); err == nil {
		t.Fatal("expected calibration outcome window tampering to be rejected")
	}
}
