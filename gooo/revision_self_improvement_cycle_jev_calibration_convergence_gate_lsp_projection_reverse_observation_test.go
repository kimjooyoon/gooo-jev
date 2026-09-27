package gooo

import "testing"

func TestObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionReverseObservation(
	t *testing.T,
) {
	gate := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGate(
		testJEVCalibrationConvergenceGateInput(
			jevCalibrationOutcomeDeltaImproved,
			"BOUND",
			jevCalibrationConvergenceGateChoiceStable,
			jevCalibrationConvergenceGateChoiceCompare,
		),
	)
	projection := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjection(gate)
	reverse := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionReverseObservation(projection)
	if reverse.Status != "BOUND" ||
		reverse.MissingStageIndex != 0 ||
		reverse.ReverseSignal != jevCalibrationConvergenceGateLSPReverseComplete {
		t.Fatalf("expected bounded LSP reverse observation, got %#v", reverse)
	}
	if reverse.ProjectionObservationDigest != projection.ObservationDigest {
		t.Fatalf("projection digest = %q, want %q", reverse.ProjectionObservationDigest, projection.ObservationDigest)
	}
	if err := reverse.Validate(); err != nil {
		t.Fatalf("expected valid LSP reverse observation: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionReverseObservationPreservesUnknown(
	t *testing.T,
) {
	gate := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGate(
		testJEVCalibrationConvergenceGateInput(
			jevCalibrationOutcomeDeltaImproved,
			"UNKNOWN",
			jevCalibrationConvergenceGateChoiceUnknown,
			jevCalibrationConvergenceGateChoiceDefer,
		),
	)
	projection := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjection(gate)
	reverse := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionReverseObservation(projection)
	if reverse.Status != "UNKNOWN" ||
		reverse.MissingStageIndex != -1 ||
		reverse.ReverseSignal != jevCalibrationConvergenceGateLSPReverseUnknown {
		t.Fatalf("expected unknown LSP reverse observation, got %#v", reverse)
	}
	if err := reverse.Validate(); err != nil {
		t.Fatalf("expected valid unknown reverse observation: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionReverseObservationRejectsTampering(
	t *testing.T,
) {
	gate := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGate(
		testJEVCalibrationConvergenceGateInput(
			jevCalibrationOutcomeDeltaImproved,
			"BOUND",
			jevCalibrationConvergenceGateChoiceStable,
			jevCalibrationConvergenceGateChoiceCompare,
		),
	)
	projection := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjection(gate)
	reverse := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionReverseObservation(projection)
	reverse.ProjectionObservationDigest = digestString("tampered")
	if err := reverse.Validate(); err == nil {
		t.Fatal("expected projection digest tampering to be rejected")
	}
}