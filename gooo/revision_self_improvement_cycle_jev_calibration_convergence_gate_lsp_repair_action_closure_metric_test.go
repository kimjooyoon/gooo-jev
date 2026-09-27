package gooo

import "testing"

func boundJEVCalibrationConvergenceGateLSPRepairActionClosureInput() JEVCalibrationConvergenceGateLSPRepairActionClosureInput {
	input := JEVCalibrationConvergenceGateLSPRepairActionClosureInput{
		SourceVersion:            "source-v1",
		ContractVersion:          "contract-v1",
		Declaration:              "decl:closure",
		IR:                       "ir:closure",
		Generated:                "generated:closure",
		ActionMaterial:           "action:repair",
		ReverseObservation:       "reverse:closure",
		EvidencePrefix:            "prefix:digest",
		ActionStatus:              string(JEVCalibrationConvergenceGateLSPRepairActionClosureBound),
		ReverseObservationStatus:  string(JEVCalibrationConvergenceGateLSPRepairActionClosureBound),
		MissingStageIndex:         -1,
	}
	input.ActionDigest = closureDigest(input.SourceVersion, input.ContractVersion, input.Declaration, input.IR, input.Generated, input.ActionMaterial, input.EvidencePrefix, "-1")
	input.ReverseObservationDigest = closureDigest(input.ActionDigest, input.ReverseObservation)
	return input
}

func TestObserveJEVCalibrationConvergenceGateLSPRepairActionClosureBindsAllStages(t *testing.T) {
	metric := ObserveJEVCalibrationConvergenceGateLSPRepairActionClosure(boundJEVCalibrationConvergenceGateLSPRepairActionClosureInput())
	if metric.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureBound {
		t.Fatalf("status = %q, want BOUND", metric.Status)
	}
	if metric.MissingStage != "" || metric.MissingStageIndex != -1 {
		t.Fatalf("missing stage metadata = %q/%d", metric.MissingStage, metric.MissingStageIndex)
	}
	if metric.MetricDigest == "" {
		t.Fatal("metric digest is empty")
	}
}

func TestObserveJEVCalibrationConvergenceGateLSPRepairActionClosurePreservesUnknownStage(t *testing.T) {
	input := boundJEVCalibrationConvergenceGateLSPRepairActionClosureInput()
	input.MissingStageIndex = 3
	metric := ObserveJEVCalibrationConvergenceGateLSPRepairActionClosure(input)
	if metric.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureUnknown {
		t.Fatalf("status = %q, want UNKNOWN", metric.Status)
	}
	if metric.MissingStage != "missing_stage_index" || metric.MissingStageIndex != 3 {
		t.Fatalf("missing stage metadata = %q/%d", metric.MissingStage, metric.MissingStageIndex)
	}
}

func TestObserveJEVCalibrationConvergenceGateLSPRepairActionClosureRejectsTampering(t *testing.T) {
	input := boundJEVCalibrationConvergenceGateLSPRepairActionClosureInput()
	input.ReverseObservation = "reverse:tampered"
	metric := ObserveJEVCalibrationConvergenceGateLSPRepairActionClosure(input)
	if metric.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureUnknown {
		t.Fatalf("status = %q, want UNKNOWN", metric.Status)
	}
	if metric.Reason != "reverse observation digest mismatch" {
		t.Fatalf("reason = %q", metric.Reason)
	}
}

func TestObserveJEVCalibrationConvergenceGateLSPRepairActionClosureDefersUnboundProducer(t *testing.T) {
	input := boundJEVCalibrationConvergenceGateLSPRepairActionClosureInput()
	input.ActionStatus = string(JEVCalibrationConvergenceGateLSPRepairActionClosureUnknown)
	metric := ObserveJEVCalibrationConvergenceGateLSPRepairActionClosure(input)
	if metric.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureDeferred {
		t.Fatalf("status = %q, want DEFERRED", metric.Status)
	}
}