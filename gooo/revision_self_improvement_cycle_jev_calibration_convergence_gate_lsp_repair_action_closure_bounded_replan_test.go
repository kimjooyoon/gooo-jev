package gooo

import "testing"

func replanTrajectory(status JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryStatus) JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryObservation {
	return JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryObservation{
		SourceVersion:        "source-v1",
		ContractVersion:      "contract-v1",
		Status:               status,
		DeltaDigest:          "delta-v1",
	}
}

func replanMetric(status JEVCalibrationConvergenceGateLSPRepairActionClosureStatus, missing string) JEVCalibrationConvergenceGateLSPRepairActionClosureMetric {
	return JEVCalibrationConvergenceGateLSPRepairActionClosureMetric{
		SourceVersion:     "source-v1",
		ContractVersion:   "contract-v1",
		Status:            status,
		MissingStage:      missing,
		MissingStageIndex: -1,
		MetricDigest:      "metric-v1",
	}
}

func TestProposeJEVCalibrationConvergenceGateLSPRepairActionClosureReplanTargetsFirstUnresolvedStage(t *testing.T) {
	candidate := ProposeJEVCalibrationConvergenceGateLSPRepairActionClosureReplan(
		replanTrajectory(JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryRegressed),
		replanMetric(JEVCalibrationConvergenceGateLSPRepairActionClosureUnknown, "reverse_observation"),
	)
	if candidate.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateProposed || candidate.TargetStage != "reverse_observation" {
		t.Fatalf("unexpected replan candidate: %+v", candidate)
	}
	if candidate.PlanDigest == "" || !candidate.IsReadOnly || candidate.CanExecute || candidate.CanAuthorize {
		t.Fatalf("replan crossed authority boundary or lacks digest: %+v", candidate)
	}
}

func TestProposeJEVCalibrationConvergenceGateLSPRepairActionClosureReplanUsesClosureDigestWhenStageIsNotNamed(t *testing.T) {
	candidate := ProposeJEVCalibrationConvergenceGateLSPRepairActionClosureReplan(
		replanTrajectory(JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryUnknown),
		replanMetric(JEVCalibrationConvergenceGateLSPRepairActionClosureBound, ""),
	)
	if candidate.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateProposed || candidate.TargetStage != "closure_digest" {
		t.Fatalf("unexpected fallback target: %+v", candidate)
	}
}

func TestProposeJEVCalibrationConvergenceGateLSPRepairActionClosureReplanDoesNotProposeForImprovement(t *testing.T) {
	candidate := ProposeJEVCalibrationConvergenceGateLSPRepairActionClosureReplan(
		replanTrajectory(JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryImproved),
		replanMetric(JEVCalibrationConvergenceGateLSPRepairActionClosureBound, ""),
	)
	if candidate.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateNoAction || candidate.TargetStage != "" {
		t.Fatalf("improvement produced a replan: %+v", candidate)
	}
}

func TestProposeJEVCalibrationConvergenceGateLSPRepairActionClosureReplanDefersProducer(t *testing.T) {
	candidate := ProposeJEVCalibrationConvergenceGateLSPRepairActionClosureReplan(
		replanTrajectory(JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryDeferred),
		replanMetric(JEVCalibrationConvergenceGateLSPRepairActionClosureDeferred, "stage"),
	)
	if candidate.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateDeferred {
		t.Fatalf("status = %q, want DEFERRED", candidate.Status)
	}
}

func TestProposeJEVCalibrationConvergenceGateLSPRepairActionClosureReplanRejectsIdentityChange(t *testing.T) {
	trajectory := replanTrajectory(JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryRegressed)
	metric := replanMetric(JEVCalibrationConvergenceGateLSPRepairActionClosureUnknown, "stage")
	trajectory.SourceVersion = "source-v2"
	candidate := ProposeJEVCalibrationConvergenceGateLSPRepairActionClosureReplan(trajectory, metric)
	if candidate.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateUnknown || candidate.Reason != "source or contract identity changed" {
		t.Fatalf("identity change was not rejected: %+v", candidate)
	}
}