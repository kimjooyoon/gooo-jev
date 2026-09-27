package gooo

import "testing"

func closureTrajectoryMetric(status JEVCalibrationConvergenceGateLSPRepairActionClosureStatus, digest string) JEVCalibrationConvergenceGateLSPRepairActionClosureMetric {
	return JEVCalibrationConvergenceGateLSPRepairActionClosureMetric{
		Status:          status,
		SourceVersion:   "source-v1",
		ContractVersion: "contract-v1",
		MetricDigest:    digest,
	}
}

func TestObserveJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryImprovesFromUnknown(t *testing.T) {
	observation := ObserveJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectory(
		closureTrajectoryMetric(JEVCalibrationConvergenceGateLSPRepairActionClosureUnknown, "metric-unknown"),
		closureTrajectoryMetric(JEVCalibrationConvergenceGateLSPRepairActionClosureBound, "metric-bound"),
	)
	if observation.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryImproved {
		t.Fatalf("status = %q, want IMPROVED", observation.Status)
	}
}

func TestObserveJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryRecordsRegression(t *testing.T) {
	observation := ObserveJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectory(
		closureTrajectoryMetric(JEVCalibrationConvergenceGateLSPRepairActionClosureBound, "metric-bound"),
		closureTrajectoryMetric(JEVCalibrationConvergenceGateLSPRepairActionClosureUnknown, "metric-unknown"),
	)
	if observation.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryRegressed {
		t.Fatalf("status = %q, want REGRESSED", observation.Status)
	}
}

func TestObserveJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryKeepsUnchangedBound(t *testing.T) {
	metric := closureTrajectoryMetric(JEVCalibrationConvergenceGateLSPRepairActionClosureBound, "metric-bound")
	observation := ObserveJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectory(metric, metric)
	if observation.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryUnchanged {
		t.Fatalf("status = %q, want UNCHANGED", observation.Status)
	}
}

func TestObserveJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryDoesNotInferChangedBoundQuality(t *testing.T) {
	observation := ObserveJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectory(
		closureTrajectoryMetric(JEVCalibrationConvergenceGateLSPRepairActionClosureBound, "metric-before"),
		closureTrajectoryMetric(JEVCalibrationConvergenceGateLSPRepairActionClosureBound, "metric-after"),
	)
	if observation.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryUnknown {
		t.Fatalf("status = %q, want UNKNOWN", observation.Status)
	}
}

func TestObserveJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryRejectsIdentityChange(t *testing.T) {
	previous := closureTrajectoryMetric(JEVCalibrationConvergenceGateLSPRepairActionClosureBound, "metric-before")
	current := closureTrajectoryMetric(JEVCalibrationConvergenceGateLSPRepairActionClosureBound, "metric-before")
	current.SourceVersion = "source-v2"
	observation := ObserveJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectory(previous, current)
	if observation.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryUnknown || observation.Reason != "source or contract identity changed" {
		t.Fatalf("identity change was not preserved as UNKNOWN: %+v", observation)
	}
}