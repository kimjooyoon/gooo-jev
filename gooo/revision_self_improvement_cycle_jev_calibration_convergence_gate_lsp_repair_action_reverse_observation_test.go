package gooo

import "testing"

func TestObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionBindsActionDigest(t *testing.T) {
	action, err := ProjectRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction(
		RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionInput{
			ProjectionObservationDigest: "projection",
			EvidencePrefixDigest:        "prefix",
			SourceObservationDigest:     "source",
			ReverseObservationDigest:    "reverse",
			GateDecision:                "jev-calibration-convergence-converged",
		},
	)
	if err != nil {
		t.Fatalf("unexpected action error: %v", err)
	}
	observation, err := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction(action)
	if err != nil {
		t.Fatalf("unexpected observation error: %v", err)
	}
	if observation.Status != "BOUND" || observation.ActionDigest != action.ActionDigest {
		t.Fatalf("observation = %#v, want bound action digest", observation)
	}
	if !observation.NonExecuting || !observation.NonAuthorizing {
		t.Fatal("reverse observation must remain non-executing and non-authorizing")
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionRejectsTampering(t *testing.T) {
	action, err := ProjectRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction(
		RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionInput{
			ProjectionObservationDigest: "projection",
			EvidencePrefixDigest:        "prefix",
			SourceObservationDigest:     "source",
			ReverseObservationDigest:    "reverse",
			GateDecision:                "jev-calibration-convergence-regressed",
		},
	)
	if err != nil {
		t.Fatalf("unexpected action error: %v", err)
	}
	action.Reason = "tampered"
	if _, err := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction(action); err == nil {
		t.Fatal("expected tampered action to be rejected")
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionPreservesUnknownStage(t *testing.T) {
	action, err := ProjectRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction(
		RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionInput{
			MissingStageIndex: 3,
		},
	)
	if err != nil {
		t.Fatalf("unexpected action error: %v", err)
	}
	observation, err := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction(action)
	if err != nil {
		t.Fatalf("unexpected observation error: %v", err)
	}
	if observation.Status != "UNKNOWN" || observation.MissingStageIndex != 3 {
		t.Fatalf("observation = %#v, want unknown stage 3", observation)
	}
}