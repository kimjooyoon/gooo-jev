package gooo

import "testing"

func TestProjectJEVCalibrationDecisionRouteReverseLSPBindsExactObservation(t *testing.T) {
	expected := reverseDecisionRouteExpected()
	observation := ObserveJEVCalibrationDecisionRouteReverseObservation(
		expected,
		"source-v1",
		"contract-v1",
		string(expected.Status),
		expected.DecisionDigest,
	)
	projection := ProjectJEVCalibrationDecisionRouteReverseLSP(observation)
	if projection.Status != JEVCalibrationDecisionRouteReverseLSPBound || projection.Code != "jev.route.reverse.bound" || projection.Severity != "Information" {
		t.Fatalf("unexpected bound projection: %+v", projection)
	}
	if !projection.IsReadOnly || projection.CanExecute || projection.CanAuthorize || len(projection.Edits) != 0 || projection.Command != "" {
		t.Fatalf("reverse LSP projection crossed authority boundary: %+v", projection)
	}
}

func TestProjectJEVCalibrationDecisionRouteReverseLSPPreservesDeferred(t *testing.T) {
	input := JEVCalibrationDecisionRouteInput{
		SourceVersion:             "source-v1",
		ContractVersion:           "contract-v1",
		DeclarationDigest:         "decl-v1",
		IRDigest:                  "ir-v1",
		GeneratedDigest:           "generated-v1",
		ReverseObservationDigest:  "reverse-v1",
		CalibrationDigest:          "calibration-v1",
		OutcomeWindowDigest:        "outcome-v1",
		ChoiceSetDigest:            "choices-v1",
		CalibratedChoiceSetDigest: "choices-v1",
		OutcomeCount:              128,
		CalibratedProbability:     0.92,
		AcceptThreshold:           0.90,
		ReviewThreshold:           0.60,
		ProviderDeferred:           true,
	}
	expected := ProjectJEVCalibrationDecisionRoute(input)
	observation := ObserveJEVCalibrationDecisionRouteReverseObservation(
		expected,
		"source-v1",
		"contract-v1",
		string(expected.Status),
		expected.DecisionDigest,
	)
	projection := ProjectJEVCalibrationDecisionRouteReverseLSP(observation)
	if projection.Status != JEVCalibrationDecisionRouteReverseLSPDeferred || projection.Code != "jev.route.reverse.deferred" {
		t.Fatalf("deferred observation was not preserved: %+v", projection)
	}
}

func TestProjectJEVCalibrationDecisionRouteReverseLSPPreservesUnknown(t *testing.T) {
	expected := reverseDecisionRouteExpected()
	observation := ObserveJEVCalibrationDecisionRouteReverseObservation(
		expected,
		"source-v2",
		"contract-v1",
		string(expected.Status),
		expected.DecisionDigest,
	)
	projection := ProjectJEVCalibrationDecisionRouteReverseLSP(observation)
	if projection.Status != JEVCalibrationDecisionRouteReverseLSPUnknown || projection.TargetStage != "reverse_observation" {
		t.Fatalf("unknown observation was not preserved: %+v", projection)
	}
}