package gooo

import "testing"

func TestProjectRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionBindsConvergence(t *testing.T) {
	action, err := ProjectRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction(
		RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionInput{
			ProjectionObservationDigest: "projection",
			EvidencePrefixDigest:        "prefix",
			SourceObservationDigest:     "source",
			ReverseObservationDigest:    "reverse",
			GateDecision:                "jev-calibration-convergence-converged",
			ProjectionSignal:             "converged",
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action.Status != "BOUND" || action.SuggestedAction != "observe-convergence" {
		t.Fatalf("action = %#v, want bound convergence action", action)
	}
	if !action.NonExecuting || !action.NonAuthorizing {
		t.Fatal("repair action must remain non-executing and non-authorizing")
	}
}

func TestProjectRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionPreservesUnknown(t *testing.T) {
	action, err := ProjectRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction(
		RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionInput{
			EvidencePrefixDigest: "prefix",
			MissingStageIndex:    4,
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action.Status != "UNKNOWN" || action.MissingStageIndex != 4 {
		t.Fatalf("action = %#v, want unknown stage 4", action)
	}
}

func TestRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionRejectsTampering(t *testing.T) {
	action, err := ProjectRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction(
		RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionInput{
			ProjectionObservationDigest: "projection",
			EvidencePrefixDigest:        "prefix",
			SourceObservationDigest:     "source",
			ReverseObservationDigest:    "reverse",
			GateDecision:                "jev-calibration-convergence-defer",
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	action.Reason = "tampered"
	if err := action.Validate(); err == nil {
		t.Fatal("expected tampered action to be rejected")
	}
}