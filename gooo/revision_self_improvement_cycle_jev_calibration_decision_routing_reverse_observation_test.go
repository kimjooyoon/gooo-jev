package gooo

import "testing"

func reverseDecisionRouteExpected() JEVCalibrationDecisionRoute {
	return ProjectJEVCalibrationDecisionRoute(JEVCalibrationDecisionRouteInput{
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
		ReviewThreshold:            0.60,
	})
}

func TestObserveJEVCalibrationDecisionRouteReverseObservationBindsExactEvidence(t *testing.T) {
	expected := reverseDecisionRouteExpected()
	observation := ObserveJEVCalibrationDecisionRouteReverseObservation(
		expected,
		"source-v1",
		"contract-v1",
		string(expected.Status),
		expected.DecisionDigest,
	)
	if observation.Status != JEVCalibrationDecisionRouteReverseBound || observation.ObservationDigest == "" {
		t.Fatalf("exact route was not bound: %+v", observation)
	}
	if !observation.IsReadOnly || observation.CanExecute || observation.CanAuthorize {
		t.Fatalf("reverse observation crossed authority boundary: %+v", observation)
	}
}

func TestObserveJEVCalibrationDecisionRouteReverseObservationRejectsDigestDrift(t *testing.T) {
	expected := reverseDecisionRouteExpected()
	observation := ObserveJEVCalibrationDecisionRouteReverseObservation(
		expected,
		"source-v1",
		"contract-v1",
		string(expected.Status),
		"sha256:tampered",
	)
	if observation.Status != JEVCalibrationDecisionRouteReverseUnknown || observation.Reason != "observed route digest changed" {
		t.Fatalf("digest drift was not preserved: %+v", observation)
	}
}

func TestObserveJEVCalibrationDecisionRouteReverseObservationRejectsIdentityDrift(t *testing.T) {
	expected := reverseDecisionRouteExpected()
	observation := ObserveJEVCalibrationDecisionRouteReverseObservation(
		expected,
		"source-v2",
		"contract-v1",
		string(expected.Status),
		expected.DecisionDigest,
	)
	if observation.Status != JEVCalibrationDecisionRouteReverseUnknown || observation.Reason != "observed source or contract identity changed" {
		t.Fatalf("identity drift was not preserved: %+v", observation)
	}
}

func TestObserveJEVCalibrationDecisionRouteReverseObservationPreservesDeferred(t *testing.T) {
	input := JEVCalibrationDecisionRouteInput{
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
	if observation.Status != JEVCalibrationDecisionRouteReverseDeferred {
		t.Fatalf("deferred route was not preserved: %+v", observation)
	}
}

func TestObserveJEVCalibrationDecisionRouteReverseObservationDoesNotCallUnknownBound(t *testing.T) {
	input := JEVCalibrationDecisionRouteInput{
		SourceVersion:             "source-v1",
		ContractVersion:           "contract-v1",
		DeclarationDigest:         "decl-v1",
		IRDigest:                  "ir-v1",
		GeneratedDigest:            "generated-v1",
		ReverseObservationDigest:  "",
		CalibrationDigest:          "calibration-v1",
		OutcomeWindowDigest:        "outcome-v1",
		ChoiceSetDigest:            "choices-v1",
		CalibratedChoiceSetDigest: "choices-v1",
		OutcomeCount:              128,
		CalibratedProbability:     0.92,
		AcceptThreshold:           0.90,
		ReviewThreshold:           0.60,
	}
	expected := ProjectJEVCalibrationDecisionRoute(input)
	observation := ObserveJEVCalibrationDecisionRouteReverseObservation(
		expected,
		"source-v1",
		"contract-v1",
		string(expected.Status),
		expected.DecisionDigest,
	)
	if observation.Status != JEVCalibrationDecisionRouteReverseUnknown {
		t.Fatalf("unresolved route was treated as bound: %+v", observation)
	}
}