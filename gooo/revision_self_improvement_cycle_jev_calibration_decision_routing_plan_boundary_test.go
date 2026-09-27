package gooo

import "testing"

func planBoundaryInput() JEVCalibrationDecisionRoutePlanBoundaryInput {
	return JEVCalibrationDecisionRoutePlanBoundaryInput{
		SourceVersion:            "source-v1",
		ContractVersion:          "contract-v1",
		RouteStatus:              JEVCalibrationDecisionRouteAcceptCandidate,
		RouteDecisionDigest:      "sha256:route-v1",
		ReverseObservationStatus: JEVCalibrationDecisionRouteReverseBound,
		ObservationDigest:        "sha256:observation-v1",
		PlanDigest:               "sha256:plan-v1",
	}
}

func TestProjectJEVCalibrationDecisionRoutePlanBoundaryProposesOnlyBoundEvidence(t *testing.T) {
	boundary := ProjectJEVCalibrationDecisionRoutePlanBoundary(planBoundaryInput())
	if boundary.Status != JEVCalibrationDecisionRoutePlanProposed || boundary.BoundaryDigest == "" {
		t.Fatalf("unexpected plan boundary: %+v", boundary)
	}
	if !boundary.IsReadOnly || boundary.CanExecute || boundary.CanAuthorize {
		t.Fatalf("plan boundary crossed authority: %+v", boundary)
	}
}

func TestProjectJEVCalibrationDecisionRoutePlanBoundaryRejectsUnresolvedObservation(t *testing.T) {
	input := planBoundaryInput()
	input.ReverseObservationStatus = JEVCalibrationDecisionRouteReverseUnknown
	boundary := ProjectJEVCalibrationDecisionRoutePlanBoundary(input)
	if boundary.Status != JEVCalibrationDecisionRoutePlanUnknown || boundary.TargetStage != "reverse_observation" {
		t.Fatalf("unresolved observation was not preserved: %+v", boundary)
	}
}

func TestProjectJEVCalibrationDecisionRoutePlanBoundaryPreservesDeferred(t *testing.T) {
	input := planBoundaryInput()
	input.ReverseObservationStatus = JEVCalibrationDecisionRouteReverseDeferred
	boundary := ProjectJEVCalibrationDecisionRoutePlanBoundary(input)
	if boundary.Status != JEVCalibrationDecisionRoutePlanDeferred {
		t.Fatalf("deferred observation was not preserved: %+v", boundary)
	}
}

func TestProjectJEVCalibrationDecisionRoutePlanBoundaryRequiresAllDigests(t *testing.T) {
	input := planBoundaryInput()
	input.PlanDigest = ""
	boundary := ProjectJEVCalibrationDecisionRoutePlanBoundary(input)
	if boundary.Status != JEVCalibrationDecisionRoutePlanUnknown || boundary.TargetStage != "plan_digest" {
		t.Fatalf("missing plan digest was not preserved: %+v", boundary)
	}
}