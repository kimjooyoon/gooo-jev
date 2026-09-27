package gooo

import "testing"

func TestProjectJEVCalibrationDecisionRoutingPlanObservationLSPProjectsBound(t *testing.T) {
	observation := ObserveJEVCalibrationDecisionRoutingPlan(JEVCalibrationDecisionRoutingPlanObservationInput{
		SourceVersion:             "source-v1",
		ContractVersion:           "contract-v1",
		DeclarationDigest:         "sha256:declaration-v1",
		IRDigest:                  "sha256:ir-v1",
		GeneratedDigest:            "sha256:generated-v1",
		ReverseObservationStatus:  "BOUND",
		ReverseObservationDigest:  "sha256:reverse-v1",
		ProposalStatus:            "PROPOSED",
		PlanDigest:                "sha256:plan-v1",
		ExpectedPlanDigest:        "sha256:plan-v1",
		OutcomeWindowDigest:       "sha256:window-v1",
		OutcomeCount:              3,
	})
	projection := ProjectJEVCalibrationDecisionRoutingPlanObservationLSP(observation)
	if projection.Status != JEVCalibrationDecisionRoutingPlanObservationLSPBound || projection.Code != "jev.plan_observation.bound" || projection.Severity != "Information" {
		t.Fatalf("unexpected bound projection: %+v", projection)
	}
	if !projection.IsReadOnly || projection.CanExecute || projection.CanAuthorize || len(projection.Edits) != 0 || projection.Command != "" {
		t.Fatalf("projection crossed authority boundary: %+v", projection)
	}
}

func TestProjectJEVCalibrationDecisionRoutingPlanObservationLSPProjectsDeferred(t *testing.T) {
	observation := ObserveJEVCalibrationDecisionRoutingPlan(JEVCalibrationDecisionRoutingPlanObservationInput{
		SourceVersion:             "source-v1",
		ContractVersion:           "contract-v1",
		DeclarationDigest:         "sha256:declaration-v1",
		IRDigest:                  "sha256:ir-v1",
		GeneratedDigest:            "sha256:generated-v1",
		ReverseObservationStatus:  "DEFERRED",
		ReverseObservationDigest:  "sha256:reverse-v1",
		ProposalStatus:            "PROPOSED",
		PlanDigest:                "sha256:plan-v1",
		ExpectedPlanDigest:        "sha256:plan-v1",
		OutcomeWindowDigest:       "sha256:window-v1",
		OutcomeCount:              3,
	})
	projection := ProjectJEVCalibrationDecisionRoutingPlanObservationLSP(observation)
	if projection.Status != JEVCalibrationDecisionRoutingPlanObservationLSPDeferred || projection.Code != "jev.plan_observation.deferred" || projection.Severity != "Hint" {
		t.Fatalf("unexpected deferred projection: %+v", projection)
	}
}

func TestProjectJEVCalibrationDecisionRoutingPlanObservationLSPPreservesMissingDigest(t *testing.T) {
	projection := ProjectJEVCalibrationDecisionRoutingPlanObservationLSP(JEVCalibrationDecisionRoutingPlanObservation{
		Status:      JEVCalibrationDecisionRoutingPlanObservationUnknown,
		TargetStage: "reverse_observation",
		Reason:      "reverse observation evidence is missing",
		IsReadOnly:  true,
	})
	if projection.Code != "jev.plan_observation.unknown" || projection.TargetStage != "observation_digest" || projection.Severity != "Error" {
		t.Fatalf("missing observation digest was not preserved: %+v", projection)
	}
}