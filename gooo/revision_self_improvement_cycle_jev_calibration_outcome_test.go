package gooo

import "testing"

func jevCalibrationOutcomeInput() RevisionSelfImprovementCycleJEVCalibrationOutcomeInput {
	return RevisionSelfImprovementCycleJEVCalibrationOutcomeInput{
		DecisionDigest:           digestString("decision:route"),
		DeclarationDigest:        digestString("declaration:route"),
		IRDigest:                 digestString("ir:route"),
		GenerationDigest:         digestString("generation:route"),
		ReverseObservationDigest: digestString("reverse:route"),
		PredictedLabel:           "route-a",
		ObservedLabel:            "route-a",
		ConfidenceMilli:          820,
		OutcomeStatus:            "observed",
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeCorrect(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationOutcome(jevCalibrationOutcomeInput())
	if observation.Status != "BOUND" ||
		observation.OutcomeClass != jevCalibrationOutcomeCorrect ||
		observation.OutcomeSignal != jevCalibrationOutcomeObserved {
		t.Fatalf("unexpected correct calibration outcome: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("correct calibration outcome did not validate: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeIncorrect(t *testing.T) {
	input := jevCalibrationOutcomeInput()
	input.ObservedLabel = "route-b"
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationOutcome(input)
	if observation.Status != "BOUND" || observation.OutcomeClass != jevCalibrationOutcomeIncorrect {
		t.Fatalf("unexpected incorrect calibration outcome: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("incorrect calibration outcome did not validate: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationOutcomePreservesUnknown(t *testing.T) {
	input := jevCalibrationOutcomeInput()
	input.ReverseObservationDigest = ""
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationOutcome(input)
	if observation.Status != "UNKNOWN" ||
		observation.MissingStage != "revision-self-improvement-cycle-jev-calibration-outcome-lineage" ||
		observation.OutcomeClass != jevCalibrationOutcomeUnknownClass {
		t.Fatalf("incomplete calibration outcome must remain unknown: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("unknown calibration outcome should validate: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationOutcomeAbstained(t *testing.T) {
	input := jevCalibrationOutcomeInput()
	input.Abstained = true
	input.ObservedLabel = "route-b"
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationOutcome(input)
	if observation.Status != "BOUND" || observation.OutcomeClass != jevCalibrationOutcomeAbstained {
		t.Fatalf("unexpected abstained calibration outcome: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("abstained calibration outcome did not validate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVCalibrationOutcomeRejectsTampering(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationOutcome(jevCalibrationOutcomeInput())
	observation.OutcomeClass = jevCalibrationOutcomeIncorrect
	if err := observation.Validate(); err == nil {
		t.Fatal("expected calibration outcome tampering to be rejected")
	}
}
