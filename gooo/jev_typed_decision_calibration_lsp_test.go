package gooo

import (
	"strings"
	"testing"
)

func typedDecisionCalibrationLSPBoundInput() ExecutionEnvelopeTypedDecisionCalibrationLSPInput {
	return ExecutionEnvelopeTypedDecisionCalibrationLSPInput{
		Status:                    ExecutionEnvelopeTypedDecisionCalibrationLSPBound,
		SourceVersion:             "gooo-jev-runtime",
		ContractVersion:           "jev-typed-decision/v1",
		ModelRevision:             "jev-latest",
		QuestionID:                "route",
		SignalEvidenceDigest:      "sha256:" + strings.Repeat("1", 64),
		CalibrationEvidenceDigest: "sha256:" + strings.Repeat("2", 64),
		OutcomeDigest:             "sha256:" + strings.Repeat("3", 64),
		SelectedProbability:       0.8,
		Confidence:                0.9,
		ConfidenceMethod:          "calibrated",
		OutcomeKnown:              true,
		ObservedOutcome:           false,
		AbsoluteError:             0.8,
		Tolerance:                 0.2,
		WithinTolerance:           false,
		WindowSize:                10,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
}

func TestProjectTypedDecisionCalibrationLSPBound(t *testing.T) {
	projection := ProjectExecutionEnvelopeTypedDecisionCalibrationLSP(typedDecisionCalibrationLSPBoundInput())
	if projection.Status != ExecutionEnvelopeTypedDecisionCalibrationLSPBound ||
		projection.MissingStage != "" ||
		projection.AbsoluteError != 0.8 {
		t.Fatalf("projection = %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestProjectTypedDecisionCalibrationLSPPreservesUnknown(t *testing.T) {
	projection := ProjectExecutionEnvelopeTypedDecisionCalibrationLSP(
		ExecutionEnvelopeTypedDecisionCalibrationLSPInput{
			Status:             ExecutionEnvelopeTypedDecisionCalibrationLSPUnknown,
			FirstMismatch:      "outcome",
			MissingStage:       "reverse_observation",
			NonExecuting:       true,
			NonAuthorizing:     true,
		},
	)
	if projection.Status != ExecutionEnvelopeTypedDecisionCalibrationLSPUnknown ||
		projection.FirstMismatch != "outcome" ||
		projection.MissingStage != "reverse_observation" {
		t.Fatalf("projection = %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestProjectTypedDecisionCalibrationLSPRejectsTampering(t *testing.T) {
	projection := ProjectExecutionEnvelopeTypedDecisionCalibrationLSP(typedDecisionCalibrationLSPBoundInput())
	projection.CalibrationEvidenceDigest = "sha256:tampered"
	if err := projection.Validate(); err == nil {
		t.Fatal("expected tampered calibration evidence to fail validation")
	}
}

func TestProjectTypedDecisionCalibrationLSPRejectsCapabilityBoundary(t *testing.T) {
	input := typedDecisionCalibrationLSPBoundInput()
	input.NonAuthorizing = false
	projection := ProjectExecutionEnvelopeTypedDecisionCalibrationLSP(input)
	if projection.Status != ExecutionEnvelopeTypedDecisionCalibrationLSPError ||
		projection.Code != executionEnvelopeTypedDecisionCalibrationLSPBoundaryCode ||
		projection.FirstMismatch != "capability-boundary" {
		t.Fatalf("projection = %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestProjectTypedDecisionCalibrationLSPRejectsUnknownConfidenceMethod(t *testing.T) {
	input := typedDecisionCalibrationLSPBoundInput()
	input.ConfidenceMethod = "untrusted-method"
	projection := ProjectExecutionEnvelopeTypedDecisionCalibrationLSP(input)
	if projection.Status != ExecutionEnvelopeTypedDecisionCalibrationLSPError ||
		projection.FirstMismatch != "confidence-method" {
		t.Fatalf("projection = %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}
