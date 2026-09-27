package gooo

import "testing"

func calibrationDecisionRouteInput() JEVCalibrationDecisionRouteInput {
	return JEVCalibrationDecisionRouteInput{
		SourceVersion:             "source-v1",
		ContractVersion:           "contract-v1",
		DeclarationDigest:         "decl-v1",
		IRDigest:                  "ir-v1",
		GeneratedDigest:            "generated-v1",
		ReverseObservationDigest:  "reverse-v1",
		CalibrationDigest:          "calibration-v1",
		OutcomeWindowDigest:        "outcome-v1",
		ChoiceSetDigest:            "choices-v1",
		CalibratedChoiceSetDigest: "choices-v1",
		OutcomeCount:              128,
		CalibratedProbability:     0.92,
		AcceptThreshold:           0.90,
		ReviewThreshold:           0.60,
	}
}

func TestProjectJEVCalibrationDecisionRouteProducesReadOnlyAcceptanceCandidate(t *testing.T) {
	route := ProjectJEVCalibrationDecisionRoute(calibrationDecisionRouteInput())
	if route.Status != JEVCalibrationDecisionRouteAcceptCandidate || route.EvidenceCoverage != 7 {
		t.Fatalf("unexpected route: %+v", route)
	}
	if route.DecisionDigest == "" || !route.IsReadOnly || route.CanExecute || route.CanAuthorize {
		t.Fatalf("route crossed authority boundary or lacks digest: %+v", route)
	}
}

func TestProjectJEVCalibrationDecisionRouteSeparatesReviewAndAbstain(t *testing.T) {
	tests := []struct {
		probability float64
		want        JEVCalibrationDecisionRouteStatus
	}{
		{probability: 0.60, want: JEVCalibrationDecisionRouteReviewCandidate},
		{probability: 0.20, want: JEVCalibrationDecisionRouteAbstainCandidate},
	}
	for _, test := range tests {
		input := calibrationDecisionRouteInput()
		input.CalibratedProbability = test.probability
		route := ProjectJEVCalibrationDecisionRoute(input)
		if route.Status != test.want {
			t.Fatalf("probability %v produced %q, want %q", test.probability, route.Status, test.want)
		}
	}
}

func TestProjectJEVCalibrationDecisionRoutePreservesFirstMissingStage(t *testing.T) {
	input := calibrationDecisionRouteInput()
	input.ReverseObservationDigest = ""
	route := ProjectJEVCalibrationDecisionRoute(input)
	if route.Status != JEVCalibrationDecisionRouteUnknown || route.MissingStage != "reverse_observation" || route.MissingStageIndex != 3 {
		t.Fatalf("unexpected missing-stage route: %+v", route)
	}
}

func TestProjectJEVCalibrationDecisionRouteRejectsChoiceSetDrift(t *testing.T) {
	input := calibrationDecisionRouteInput()
	input.CalibratedChoiceSetDigest = "choices-v2"
	route := ProjectJEVCalibrationDecisionRoute(input)
	if route.Status != JEVCalibrationDecisionRouteUnknown || route.MissingStage != "choice_set" {
		t.Fatalf("choice-set drift was not rejected: %+v", route)
	}
}

func TestProjectJEVCalibrationDecisionRouteDefersProducer(t *testing.T) {
	input := calibrationDecisionRouteInput()
	input.ProviderDeferred = true
	route := ProjectJEVCalibrationDecisionRoute(input)
	if route.Status != JEVCalibrationDecisionRouteDeferred || route.DecisionDigest == "" {
		t.Fatalf("deferred route was not preserved: %+v", route)
	}
}

func TestProjectJEVCalibrationDecisionRouteRequiresOutcomeEvidence(t *testing.T) {
	input := calibrationDecisionRouteInput()
	input.OutcomeCount = 0
	route := ProjectJEVCalibrationDecisionRoute(input)
	if route.Status != JEVCalibrationDecisionRouteUnknown || route.MissingStage != "outcome_window" {
		t.Fatalf("outcome evidence gap was not preserved: %+v", route)
	}
}