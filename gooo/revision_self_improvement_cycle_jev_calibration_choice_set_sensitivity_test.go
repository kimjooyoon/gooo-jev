package gooo

import "testing"

func jevCalibrationChoiceSetSensitivityInput() RevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivityInput {
	baselineInput := jevCalibrationOutcomeWindowInput()
	currentInput := jevCalibrationOutcomeWindowInput()
	currentInput.GenerationTraceDigest = digestString("generation-trace:revision-2")
	currentInput.Observations = []RevisionSelfImprovementCycleJEVCalibrationOutcomeObservation{
		jevCalibrationOutcomeWindowSample("route-a", "route-a", 900, false),
		jevCalibrationOutcomeWindowSample("route-a", "route-a", 700, false),
		jevCalibrationOutcomeWindowSample("route-a", "route-b", 600, true),
	}
	return RevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivityInput{
		Baseline:                  ObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeWindow(baselineInput),
		Current:                   ObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeWindow(currentInput),
		BaselineChoiceSetDigest:   digestString("choices:v1"),
		CurrentChoiceSetDigest:    digestString("choices:v1"),
		BaselineChoiceOrderDigest: digestString("order:v1"),
		CurrentChoiceOrderDigest:  digestString("order:v1"),
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivity(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivity(
		jevCalibrationChoiceSetSensitivityInput(),
	)
	if observation.Status != "BOUND" ||
		observation.AccuracyDeltaMilli != 500 ||
		observation.SensitivitySignal != jevCalibrationChoiceSetSensitivityStable ||
		observation.DecisionSignal != jevCalibrationChoiceSetDecisionCompare {
		t.Fatalf("unexpected stable choice-set sensitivity: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("stable choice-set sensitivity did not validate: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivityDefersOnChoiceSetChange(t *testing.T) {
	input := jevCalibrationChoiceSetSensitivityInput()
	input.CurrentChoiceSetDigest = digestString("choices:v2")
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivity(input)
	if observation.Status != "BOUND" ||
		observation.SensitivitySignal != jevCalibrationChoiceSetSensitivityChanged ||
		observation.DecisionSignal != jevCalibrationChoiceSetDecisionDefer {
		t.Fatalf("choice-set change must defer comparison: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("choice-set change should validate as deferred evidence: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivityPreservesUnknown(t *testing.T) {
	input := jevCalibrationChoiceSetSensitivityInput()
	input.Current.SourceObservationDigest = digestString("source:other")
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivity(input)
	if observation.Status != "UNKNOWN" ||
		observation.MissingStage != "revision-self-improvement-cycle-jev-calibration-choice-set-sensitivity-current" {
		t.Fatalf("tampered current window must remain unknown: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("unknown choice-set sensitivity should validate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivityRejectsTampering(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationChoiceSetSensitivity(
		jevCalibrationChoiceSetSensitivityInput(),
	)
	observation.DecisionSignal = jevCalibrationChoiceSetDecisionDefer
	if err := observation.Validate(); err == nil {
		t.Fatal("expected stable choice-set sensitivity tampering to be rejected")
	}
}