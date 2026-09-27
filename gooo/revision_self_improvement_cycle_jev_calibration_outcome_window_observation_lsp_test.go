package gooo

import "testing"

func TestProjectJEVCalibrationOutcomeWindowObservationLSPProjectsBound(t *testing.T) {
	observation := ObserveJEVCalibrationOutcomeWindow(JEVCalibrationOutcomeWindowObservationInput{
		SourceVersion:           "source-v1",
		ContractVersion:         "contract-v1",
		ChoiceSetDigest:         "sha256:choices-v1",
		ExpectedChoiceSetDigest: "sha256:choices-v1",
		EvidencePrefixDigest:    "sha256:prefix-v1",
		ObservationDigest:       "sha256:observation-v1",
		OutcomeCount:            4,
	})
	projection := ProjectJEVCalibrationOutcomeWindowObservationLSP(observation)
	if projection.Status != JEVCalibrationOutcomeWindowObservationLSPBound || projection.Code != "jev.outcome_window.bound" || projection.Severity != "Information" {
		t.Fatalf("unexpected bound projection: %+v", projection)
	}
	if !projection.IsReadOnly || projection.CanExecute || projection.CanAuthorize || len(projection.Edits) != 0 || projection.Command != "" {
		t.Fatalf("projection crossed authority boundary: %+v", projection)
	}
}

func TestProjectJEVCalibrationOutcomeWindowObservationLSPProjectsDeferred(t *testing.T) {
	observation := ObserveJEVCalibrationOutcomeWindow(JEVCalibrationOutcomeWindowObservationInput{
		SourceVersion:           "source-v1",
		ContractVersion:         "contract-v1",
		ChoiceSetDigest:         "sha256:choices-v1",
		ExpectedChoiceSetDigest: "sha256:choices-v1",
		EvidencePrefixDigest:    "sha256:prefix-v1",
		ObservationDigest:       "sha256:observation-v1",
		ProducerDeferred:        true,
	})
	projection := ProjectJEVCalibrationOutcomeWindowObservationLSP(observation)
	if projection.Status != JEVCalibrationOutcomeWindowObservationLSPDeferred || projection.Code != "jev.outcome_window.deferred" || projection.Severity != "Hint" {
		t.Fatalf("unexpected deferred projection: %+v", projection)
	}
}

func TestProjectJEVCalibrationOutcomeWindowObservationLSPPreservesMissingDigest(t *testing.T) {
	projection := ProjectJEVCalibrationOutcomeWindowObservationLSP(JEVCalibrationOutcomeWindowObservation{
		Status:      JEVCalibrationOutcomeWindowObservationUnknown,
		TargetStage: "outcome_window",
		Reason:      "outcome observation digest is missing",
		IsReadOnly:  true,
	})
	if projection.Code != "jev.outcome_window.unknown" || projection.TargetStage != "record_digest" || projection.Severity != "Error" {
		t.Fatalf("missing record digest was not preserved: %+v", projection)
	}
}
