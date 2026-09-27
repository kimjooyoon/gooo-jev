package gooo

import (
	"strings"
	"testing"
)

func typedDecisionCalibrationWindowLSPBoundInput() ExecutionEnvelopeTypedDecisionCalibrationWindowLSPInput {
	return ExecutionEnvelopeTypedDecisionCalibrationWindowLSPInput{
		Status:                    ExecutionEnvelopeTypedDecisionCalibrationWindowLSPBound,
		SourceVersion:             "gooo-jev-runtime",
		ContractVersion:           "jev-calibration-window/v1",
		ObservationCount:          4,
		KnownObservationCount:     4,
		WithinToleranceCount:      3,
		MeanAbsoluteError:         0.25,
		EvidenceCoverage:          1,
		MinimumWindow:             4,
		ObservationEvidenceDigest: "sha256:" + strings.Repeat("4", 64),
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
}

func TestProjectTypedDecisionCalibrationWindowLSPBound(t *testing.T) {
	projection := ProjectExecutionEnvelopeTypedDecisionCalibrationWindowLSP(typedDecisionCalibrationWindowLSPBoundInput())
	if projection.Status != ExecutionEnvelopeTypedDecisionCalibrationWindowLSPBound ||
		projection.FirstMismatch != "" ||
		projection.EvidenceCoverage != 1 {
		t.Fatalf("projection = %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestProjectTypedDecisionCalibrationWindowLSPPreservesUnknown(t *testing.T) {
	projection := ProjectExecutionEnvelopeTypedDecisionCalibrationWindowLSP(
		ExecutionEnvelopeTypedDecisionCalibrationWindowLSPInput{
			Status:         ExecutionEnvelopeTypedDecisionCalibrationWindowLSPUnknown,
			FirstMismatch:  "minimum-window",
			MissingStage:   "calibration-window",
			NonExecuting:   true,
			NonAuthorizing: true,
		},
	)
	if projection.Status != ExecutionEnvelopeTypedDecisionCalibrationWindowLSPUnknown ||
		projection.FirstMismatch != "minimum-window" ||
		projection.MissingStage != "calibration-window" {
		t.Fatalf("projection = %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestProjectTypedDecisionCalibrationWindowLSPRejectsTamperedDigest(t *testing.T) {
	input := typedDecisionCalibrationWindowLSPBoundInput()
	input.ObservationEvidenceDigest = "sha256:tampered"
	projection := ProjectExecutionEnvelopeTypedDecisionCalibrationWindowLSP(input)
	if projection.Status != ExecutionEnvelopeTypedDecisionCalibrationWindowLSPError ||
		projection.FirstMismatch != "observation-evidence-digest" {
		t.Fatalf("projection = %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestProjectTypedDecisionCalibrationWindowLSPRejectsCapabilityBoundary(t *testing.T) {
	input := typedDecisionCalibrationWindowLSPBoundInput()
	input.NonAuthorizing = false
	projection := ProjectExecutionEnvelopeTypedDecisionCalibrationWindowLSP(input)
	if projection.Status != ExecutionEnvelopeTypedDecisionCalibrationWindowLSPError ||
		projection.FirstMismatch != "capability-boundary" {
		t.Fatalf("projection = %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}
