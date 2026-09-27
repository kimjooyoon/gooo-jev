package gooo

import "testing"

func TestProjectJEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPProjectsBound(t *testing.T) {
	observation := ObserveJEVCalibrationDecisionRoutingPlanOutcomeWindow(JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationInput{
		SourceVersion:                          "source-v1",
		ContractVersion:                        "contract-v1",
		PlanObservationDigest:                  "sha256:plan-v1",
		ExpectedPlanObservationDigest:          "sha256:plan-v1",
		OutcomeWindowObservationDigest:         "sha256:window-v1",
		ExpectedOutcomeWindowObservationDigest: "sha256:window-v1",
	})
	projection := ProjectJEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSP(observation)
	if projection.Status != JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPBound || projection.Code != "jev.plan_outcome_window.bound" || projection.Severity != "Information" {
		t.Fatalf("unexpected bound projection: %+v", projection)
	}
	if !projection.IsReadOnly || projection.CanExecute || projection.CanAuthorize || len(projection.Edits) != 0 || projection.Command != "" {
		t.Fatalf("projection crossed authority boundary: %+v", projection)
	}
}

func TestProjectJEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPProjectsDeferred(t *testing.T) {
	observation := ObserveJEVCalibrationDecisionRoutingPlanOutcomeWindow(JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationInput{
		SourceVersion:                          "source-v1",
		ContractVersion:                        "contract-v1",
		PlanObservationDigest:                  "sha256:plan-v1",
		ExpectedPlanObservationDigest:          "sha256:plan-v1",
		OutcomeWindowObservationDigest:         "sha256:window-v1",
		ExpectedOutcomeWindowObservationDigest: "sha256:window-v1",
		ProducerDeferred:                       true,
	})
	projection := ProjectJEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSP(observation)
	if projection.Status != JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPDeferred || projection.Code != "jev.plan_outcome_window.deferred" || projection.Severity != "Hint" {
		t.Fatalf("unexpected deferred projection: %+v", projection)
	}
}

func TestProjectJEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPPreservesMissingDigest(t *testing.T) {
	projection := ProjectJEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSP(JEVCalibrationDecisionRoutingPlanOutcomeWindowObservation{
		Status:      JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationUnknown,
		TargetStage: "outcome_window_observation_digest",
		Reason:      "outcome-window observation digest is missing",
	})
	if projection.Code != "jev.plan_outcome_window.unknown" || projection.TargetStage != "chain_digest" || projection.Severity != "Error" {
		t.Fatalf("missing chain digest was not preserved: %+v", projection)
	}
}
