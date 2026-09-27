package gooo

import "testing"

func TestProjectJEVCalibrationConfidenceGateObservationLSPAuto(t *testing.T) {
	projection := ProjectJEVCalibrationConfidenceGateObservationLSP(
		JEVCalibrationConfidenceGateObservationLSPInput{
			Status:               "AUTO",
			Confidence:           0.96,
			Threshold:            0.90,
			EvidencePrefixDigest: "prefix-digest",
			SourceVersion:        "jev-1.13.0",
			MissingStageIndex:    -1,
		},
	)

	if projection.Status != JEVCalibrationConfidenceGateObservationLSPInformation {
		t.Fatalf("status = %q, want Information", projection.Status)
	}
	if projection.Code != "jev.confidence_gate.auto" {
		t.Fatalf("code = %q, want jev.confidence_gate.auto", projection.Code)
	}
	if projection.EvidencePrefixDigest != "prefix-digest" || projection.SourceVersion != "jev-1.13.0" {
		t.Fatal("projection did not preserve evidence identity")
	}
	if !projection.IsReadOnly || projection.CanEdit || projection.CanExecute || projection.CanAuthorize {
		t.Fatal("LSP projection must remain read-only and non-authorizing")
	}
}

func TestProjectJEVCalibrationConfidenceGateObservationLSPReview(t *testing.T) {
	projection := ProjectJEVCalibrationConfidenceGateObservationLSP(
		JEVCalibrationConfidenceGateObservationLSPInput{
			Status:               "REVIEW",
			Confidence:           0.72,
			Threshold:            0.90,
			EvidencePrefixDigest: "prefix-digest",
			SourceVersion:        "jev-1.13.0",
			MissingStageIndex:    -1,
		},
	)

	if projection.Status != JEVCalibrationConfidenceGateObservationLSPHint {
		t.Fatalf("status = %q, want Hint", projection.Status)
	}
	if projection.Code != "jev.confidence_gate.review" {
		t.Fatalf("code = %q, want jev.confidence_gate.review", projection.Code)
	}
}

func TestProjectJEVCalibrationConfidenceGateObservationLSPUnknown(t *testing.T) {
	projection := ProjectJEVCalibrationConfidenceGateObservationLSP(
		JEVCalibrationConfidenceGateObservationLSPInput{
			Status:            "UNKNOWN",
			MissingStageIndex: 3,
		},
	)

	if projection.Status != JEVCalibrationConfidenceGateObservationLSPError {
		t.Fatalf("status = %q, want Error", projection.Status)
	}
	if projection.MissingStageIndex != 3 {
		t.Fatalf("missing stage = %d, want 3", projection.MissingStageIndex)
	}
}