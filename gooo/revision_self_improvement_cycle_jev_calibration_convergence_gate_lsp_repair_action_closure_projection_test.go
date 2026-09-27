package gooo

import "testing"

func TestProjectJEVCalibrationConvergenceGateLSPRepairActionClosureBoundIsInformational(t *testing.T) {
	projection := ProjectJEVCalibrationConvergenceGateLSPRepairActionClosure(JEVCalibrationConvergenceGateLSPRepairActionClosureMetric{
		Status:            JEVCalibrationConvergenceGateLSPRepairActionClosureBound,
		MetricDigest:      "metric-bound",
		MissingStageIndex: -1,
	})
	if projection.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureProjectionBound {
		t.Fatalf("status = %q, want BOUND", projection.Status)
	}
	if projection.IsActionable || projection.Kind != "info" {
		t.Fatalf("bound projection is actionable or not informational: %+v", projection)
	}
	if !projection.IsReadOnly || projection.CanExecute || projection.CanAuthorize || projection.Command != "" || len(projection.Edits) != 0 {
		t.Fatalf("bound projection crossed authority boundary: %+v", projection)
	}
}

func TestProjectJEVCalibrationConvergenceGateLSPRepairActionClosurePreservesUnknown(t *testing.T) {
	projection := ProjectJEVCalibrationConvergenceGateLSPRepairActionClosure(JEVCalibrationConvergenceGateLSPRepairActionClosureMetric{
		Status:            JEVCalibrationConvergenceGateLSPRepairActionClosureUnknown,
		MetricDigest:      "metric-unknown",
		MissingStage:      "reverse_observation",
		MissingStageIndex: 4,
	})
	if projection.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureProjectionUnknown {
		t.Fatalf("status = %q, want UNKNOWN", projection.Status)
	}
	if !projection.IsActionable || projection.MissingStage != "reverse_observation" || projection.MissingStageIndex != 4 {
		t.Fatalf("unknown provenance was not preserved: %+v", projection)
	}
	if projection.Code != "gooo.provenance.unknown" || projection.Title != "Inspect provenance closure" {
		t.Fatalf("unexpected unknown projection: %+v", projection)
	}
}

func TestProjectJEVCalibrationConvergenceGateLSPRepairActionClosureDefersProducerState(t *testing.T) {
	projection := ProjectJEVCalibrationConvergenceGateLSPRepairActionClosure(JEVCalibrationConvergenceGateLSPRepairActionClosureMetric{
		Status:       JEVCalibrationConvergenceGateLSPRepairActionClosureDeferred,
		MetricDigest: "metric-deferred",
	})
	if projection.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureProjectionDeferred || projection.Code != "gooo.provenance.deferred" {
		t.Fatalf("unexpected deferred projection: %+v", projection)
	}
	if !projection.IsActionable || !projection.IsReadOnly || projection.CanExecute || projection.CanAuthorize {
		t.Fatalf("deferred projection crossed authority boundary: %+v", projection)
	}
}