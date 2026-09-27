package gooo

import (
	"fmt"
	"math"
)

const (
	ExecutionEnvelopeTypedDecisionCalibrationWindowLSPBound   = "BOUND"
	ExecutionEnvelopeTypedDecisionCalibrationWindowLSPUnknown = "UNKNOWN"
	ExecutionEnvelopeTypedDecisionCalibrationWindowLSPError   = "ERROR"

	executionEnvelopeTypedDecisionCalibrationWindowLSPBoundCode    = "gooo.typed_decision_calibration_window.bound"
	executionEnvelopeTypedDecisionCalibrationWindowLSPUnknownCode  = "gooo.typed_decision_calibration_window.unknown"
	executionEnvelopeTypedDecisionCalibrationWindowLSPIntegrityCode = "gooo.typed_decision_calibration_window.integrity"
	executionEnvelopeTypedDecisionCalibrationWindowLSPBoundaryCode  = "gooo.typed_decision_calibration_window.capability_boundary"
)

// ExecutionEnvelopeTypedDecisionCalibrationWindowLSPInput mirrors an
// already validated runtime window without executing or re-evaluating it.
type ExecutionEnvelopeTypedDecisionCalibrationWindowLSPInput struct {
	Status                    string
	SourceVersion             string
	ContractVersion           string
	ObservationCount          int
	KnownObservationCount     int
	WithinToleranceCount      int
	MeanAbsoluteError         float64
	EvidenceCoverage          float64
	MinimumWindow             int
	FirstMismatch             string
	MissingStage              string
	ObservationEvidenceDigest string
	NonExecuting              bool
	NonAuthorizing            bool
}

// ExecutionEnvelopeTypedDecisionCalibrationWindowLSPProjection exposes a
// read-only calibration-window summary to an editor.
type ExecutionEnvelopeTypedDecisionCalibrationWindowLSPProjection struct {
	Status                    string
	Code                      string
	Severity                  string
	Message                   string
	SourceVersion             string
	ContractVersion           string
	ObservationCount          int
	KnownObservationCount     int
	WithinToleranceCount      int
	MeanAbsoluteError         float64
	EvidenceCoverage          float64
	MinimumWindow             int
	FirstMismatch             string
	MissingStage              string
	ObservationEvidenceDigest string
	EvidenceDigest            string
	NonExecuting              bool
	NonAuthorizing            bool
}

// ProjectExecutionEnvelopeTypedDecisionCalibrationWindowLSP preserves the
// runtime window's evidence state and never claims improvement or authority.
func ProjectExecutionEnvelopeTypedDecisionCalibrationWindowLSP(
	input ExecutionEnvelopeTypedDecisionCalibrationWindowLSPInput,
) ExecutionEnvelopeTypedDecisionCalibrationWindowLSPProjection {
	output := ExecutionEnvelopeTypedDecisionCalibrationWindowLSPProjection{
		Status:                    ExecutionEnvelopeTypedDecisionCalibrationWindowLSPUnknown,
		Code:                      executionEnvelopeTypedDecisionCalibrationWindowLSPUnknownCode,
		Severity:                  "warning",
		Message:                   "typed decision calibration window is unresolved",
		SourceVersion:             input.SourceVersion,
		ContractVersion:           input.ContractVersion,
		ObservationCount:          input.ObservationCount,
		KnownObservationCount:     input.KnownObservationCount,
		WithinToleranceCount:      input.WithinToleranceCount,
		MeanAbsoluteError:         input.MeanAbsoluteError,
		EvidenceCoverage:          input.EvidenceCoverage,
		MinimumWindow:             input.MinimumWindow,
		FirstMismatch:             input.FirstMismatch,
		MissingStage:              input.MissingStage,
		ObservationEvidenceDigest: input.ObservationEvidenceDigest,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "calibration-window"
	}
	if output.FirstMismatch == "" {
		output.FirstMismatch = "calibration-window"
	}

	switch {
	case !input.NonExecuting || !input.NonAuthorizing:
		output.Status = ExecutionEnvelopeTypedDecisionCalibrationWindowLSPError
		output.Code = executionEnvelopeTypedDecisionCalibrationWindowLSPBoundaryCode
		output.Severity = "error"
		output.Message = "typed decision calibration window crossed a capability boundary"
		output.FirstMismatch = "capability-boundary"
		output.MissingStage = "capability-boundary"
	case input.Status == ExecutionEnvelopeTypedDecisionCalibrationWindowLSPUnknown:
		output.Status = ExecutionEnvelopeTypedDecisionCalibrationWindowLSPUnknown
		output.Code = executionEnvelopeTypedDecisionCalibrationWindowLSPUnknownCode
		output.Severity = "warning"
		output.Message = "typed decision calibration window is UNKNOWN: evidence is incomplete"
	case input.Status != ExecutionEnvelopeTypedDecisionCalibrationWindowLSPBound:
		output.Status = ExecutionEnvelopeTypedDecisionCalibrationWindowLSPError
		output.Code = executionEnvelopeTypedDecisionCalibrationWindowLSPIntegrityCode
		output.Severity = "error"
		output.Message = "typed decision calibration window status is invalid"
		output.FirstMismatch = "status"
		output.MissingStage = "calibration-window"
	case calibrationWindowLSPInputError(input) != "":
		output.Status = ExecutionEnvelopeTypedDecisionCalibrationWindowLSPError
		output.Code = executionEnvelopeTypedDecisionCalibrationWindowLSPIntegrityCode
		output.Severity = "error"
		output.Message = "typed decision calibration window failed integrity validation"
		output.FirstMismatch = calibrationWindowLSPInputError(input)
		output.MissingStage = "calibration-window"
	default:
		output.Status = ExecutionEnvelopeTypedDecisionCalibrationWindowLSPBound
		output.Code = executionEnvelopeTypedDecisionCalibrationWindowLSPBoundCode
		output.Severity = "info"
		output.Message = "typed decision calibration window is available for inspection"
		output.FirstMismatch = ""
		output.MissingStage = ""
	}
	output.EvidenceDigest = executionEnvelopeTypedDecisionCalibrationWindowLSPDigest(output)
	return output
}

func calibrationWindowLSPInputError(
	input ExecutionEnvelopeTypedDecisionCalibrationWindowLSPInput,
) string {
	switch {
	case input.SourceVersion == "":
		return "source-version"
	case input.ContractVersion == "":
		return "contract-version"
	case !validExternalTypedDecisionCalibrationDigest(input.ObservationEvidenceDigest):
		return "observation-evidence-digest"
	case input.ObservationCount <= 0:
		return "observation-count"
	case input.KnownObservationCount != input.ObservationCount:
		return "known-observation-count"
	case input.WithinToleranceCount < 0 || input.WithinToleranceCount > input.KnownObservationCount:
		return "within-tolerance-count"
	case input.MinimumWindow <= 0 || input.ObservationCount < input.MinimumWindow:
		return "minimum-window"
	case !finiteCalibrationWindowLSPUnit(input.MeanAbsoluteError):
		return "mean-absolute-error"
	case !finiteCalibrationWindowLSPUnit(input.EvidenceCoverage):
		return "evidence-coverage"
	case input.EvidenceCoverage != 1:
		return "evidence-coverage"
	case input.FirstMismatch != "" || input.MissingStage != "":
		return "observation-completeness"
	default:
		return ""
	}
}

func finiteCalibrationWindowLSPUnit(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value >= 0 && value <= 1
}

func executionEnvelopeTypedDecisionCalibrationWindowLSPDigest(
	projection ExecutionEnvelopeTypedDecisionCalibrationWindowLSPProjection,
) string {
	return digestString(fmt.Sprintf(
		"gooo-typed-decision-calibration-window-lsp|%s|%s|%s|%s|%s|%s|%d|%d|%d|%0.9f|%0.9f|%d|%s|%s|%s|%t|%t",
		projection.Status,
		projection.Code,
		projection.Severity,
		projection.Message,
		projection.SourceVersion,
		projection.ContractVersion,
		projection.ObservationCount,
		projection.KnownObservationCount,
		projection.WithinToleranceCount,
		projection.MeanAbsoluteError,
		projection.EvidenceCoverage,
		projection.MinimumWindow,
		projection.FirstMismatch,
		projection.MissingStage,
		projection.ObservationEvidenceDigest,
		projection.NonExecuting,
		projection.NonAuthorizing,
	))
}

func (projection ExecutionEnvelopeTypedDecisionCalibrationWindowLSPProjection) Validate() error {
	switch projection.Status {
	case ExecutionEnvelopeTypedDecisionCalibrationWindowLSPBound,
		ExecutionEnvelopeTypedDecisionCalibrationWindowLSPUnknown,
		ExecutionEnvelopeTypedDecisionCalibrationWindowLSPError:
	default:
		return fmt.Errorf("invalid typed decision calibration window LSP status %q", projection.Status)
	}
	if !projection.NonExecuting || !projection.NonAuthorizing {
		return fmt.Errorf("typed decision calibration window LSP crossed a capability boundary")
	}
	if projection.EvidenceDigest != executionEnvelopeTypedDecisionCalibrationWindowLSPDigest(projection) {
		return fmt.Errorf("typed decision calibration window LSP evidence digest mismatch")
	}
	if projection.Status == ExecutionEnvelopeTypedDecisionCalibrationWindowLSPBound {
		if calibrationWindowLSPInputError(ExecutionEnvelopeTypedDecisionCalibrationWindowLSPInput{
			Status:                    ExecutionEnvelopeTypedDecisionCalibrationWindowLSPBound,
			SourceVersion:             projection.SourceVersion,
			ContractVersion:           projection.ContractVersion,
			ObservationCount:          projection.ObservationCount,
			KnownObservationCount:     projection.KnownObservationCount,
			WithinToleranceCount:      projection.WithinToleranceCount,
			MeanAbsoluteError:         projection.MeanAbsoluteError,
			EvidenceCoverage:          projection.EvidenceCoverage,
			MinimumWindow:             projection.MinimumWindow,
			FirstMismatch:             projection.FirstMismatch,
			MissingStage:              projection.MissingStage,
			ObservationEvidenceDigest: projection.ObservationEvidenceDigest,
		}) != "" {
			return fmt.Errorf("bound typed decision calibration window is incomplete")
		}
	}
	if projection.Status == ExecutionEnvelopeTypedDecisionCalibrationWindowLSPUnknown &&
		(projection.FirstMismatch == "" || projection.MissingStage == "") {
		return fmt.Errorf("unknown typed decision calibration window lost its first mismatch")
	}
	return nil
}
