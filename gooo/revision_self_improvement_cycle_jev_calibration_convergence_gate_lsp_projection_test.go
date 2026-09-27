package gooo

import "testing"

func TestObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjection(t *testing.T) {
	gate := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGate(
		testJEVCalibrationConvergenceGateInput(
			jevCalibrationOutcomeDeltaImproved,
			"BOUND",
			jevCalibrationConvergenceGateChoiceStable,
			jevCalibrationConvergenceGateChoiceCompare,
		),
	)
	projection := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjection(gate)
	if projection.Status != "BOUND" ||
		projection.MissingStageIndex != 0 ||
		projection.ProjectionSignal != jevCalibrationConvergenceGateLSPConverged {
		t.Fatalf("expected converged LSP projection, got %#v", projection)
	}
	if projection.EvidencePrefixDigest != gate.ObservationDigest {
		t.Fatalf("evidence prefix = %q, want gate digest %q", projection.EvidencePrefixDigest, gate.ObservationDigest)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("expected valid LSP projection: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionPreservesUnknown(
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
	if projection.Status != "UNKNOWN" ||
		projection.MissingStageIndex != -1 ||
		projection.ProjectionSignal != jevCalibrationConvergenceGateLSPUnknown {
		t.Fatalf("expected unknown LSP projection, got %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("expected valid unknown LSP projection: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionRejectsTampering(
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
	projection.EvidencePrefixDigest = digestString("tampered")
	if err := projection.Validate(); err == nil {
		t.Fatal("expected evidence prefix tampering to be rejected")
	}
}