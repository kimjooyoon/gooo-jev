package gooo

import (
	"encoding/hex"
	"fmt"
	"math"
	"strings"
)

const (
	ExecutionEnvelopeTypedDecisionCalibrationLSPBound  = "BOUND"
	ExecutionEnvelopeTypedDecisionCalibrationLSPUnknown = "UNKNOWN"
	ExecutionEnvelopeTypedDecisionCalibrationLSPError   = "ERROR"

	executionEnvelopeTypedDecisionCalibrationLSPBoundCode  = "gooo.typed_decision_calibration.bound"
	executionEnvelopeTypedDecisionCalibrationLSPUnknownCode = "gooo.typed_decision_calibration.unknown"
	executionEnvelopeTypedDecisionCalibrationLSPIntegrityCode = "gooo.typed_decision_calibration.integrity"
	executionEnvelopeTypedDecisionCalibrationLSPBoundaryCode = "gooo.typed_decision_calibration.capability_boundary"
)

// ExecutionEnvelopeTypedDecisionCalibrationLSPInput mirrors the runtime
// calibration observation across a repository boundary.
type ExecutionEnvelopeTypedDecisionCalibrationLSPInput struct {
	Status                    string
	SourceVersion             string
	ContractVersion           string
	ModelRevision             string
	QuestionID                string
	SignalEvidenceDigest      string
	CalibrationEvidenceDigest string
	OutcomeDigest             string
	SelectedProbability       float64
	Confidence                float64
	OutcomeKnown              bool
	ObservedOutcome           bool
	AbsoluteError             float64
	Tolerance                 float64
	WithinTolerance           bool
	WindowSize                int
	FirstMismatch             string
	MissingStage              string
	NonExecuting              bool
	NonAuthorizing            bool
}

// ExecutionEnvelopeTypedDecisionCalibrationLSPProjection exposes calibration
// evidence to an editor without executing or authorizing the decision.
type ExecutionEnvelopeTypedDecisionCalibrationLSPProjection struct {
	Status                    string
	Code                      string
	Severity                  string
	Message                   string
	SourceVersion             string
	ContractVersion           string
	ModelRevision             string
	QuestionID                string
	SignalEvidenceDigest      string
	CalibrationEvidenceDigest string
	OutcomeDigest             string
	SelectedProbability       float64
	Confidence                float64
	OutcomeKnown              bool
	ObservedOutcome           bool
	AbsoluteError             float64
	Tolerance                 float64
	WithinTolerance           bool
	WindowSize                int
	FirstMismatch              string
	MissingStage              string
	EvidenceDigest             string
	NonExecuting              bool
	NonAuthorizing             bool
}

// ProjectExecutionEnvelopeTypedDecisionCalibrationLSP preserves runtime
// calibration status and first-mismatch evidence as LSP-visible data.
func ProjectExecutionEnvelopeTypedDecisionCalibrationLSP(
	input ExecutionEnvelopeTypedDecisionCalibrationLSPInput,
) ExecutionEnvelopeTypedDecisionCalibrationLSPProjection {
	output := ExecutionEnvelopeTypedDecisionCalibrationLSPProjection{
		Status:                    ExecutionEnvelopeTypedDecisionCalibrationLSPUnknown,
		Code:                      executionEnvelopeTypedDecisionCalibrationLSPUnknownCode,
		Severity:                  "warning",
		Message:                   "typed decision calibration is unresolved",
		SourceVersion:             input.SourceVersion,
		ContractVersion:           input.ContractVersion,
		ModelRevision:             input.ModelRevision,
		QuestionID:                input.QuestionID,
		SignalEvidenceDigest:      input.SignalEvidenceDigest,
		CalibrationEvidenceDigest: input.CalibrationEvidenceDigest,
		OutcomeDigest:             input.OutcomeDigest,
		SelectedProbability:       input.SelectedProbability,
		Confidence:                input.Confidence,
		OutcomeKnown:              input.OutcomeKnown,
		ObservedOutcome:           input.ObservedOutcome,
		AbsoluteError:             input.AbsoluteError,
		Tolerance:                 input.Tolerance,
		WithinTolerance:           input.WithinTolerance,
		WindowSize:                input.WindowSize,
		FirstMismatch:             input.FirstMismatch,
		MissingStage:              input.MissingStage,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "typed_decision_calibration"
	}
	if output.FirstMismatch == "" {
		output.FirstMismatch = "typed_decision_calibration"
	}

	switch {
	case !input.NonExecuting || !input.NonAuthorizing:
		output.Status = ExecutionEnvelopeTypedDecisionCalibrationLSPError
		output.Code = executionEnvelopeTypedDecisionCalibrationLSPBoundaryCode
		output.Severity = "error"
		output.Message = "typed decision calibration crossed a capability boundary"
		output.FirstMismatch = "capability-boundary"
		output.MissingStage = "capability-boundary"
	case input.Status == ExecutionEnvelopeTypedDecisionCalibrationLSPUnknown:
		output.Status = ExecutionEnvelopeTypedDecisionCalibrationLSPUnknown
		output.Code = executionEnvelopeTypedDecisionCalibrationLSPUnknownCode
		output.Severity = "warning"
		output.Message = "typed decision calibration is UNKNOWN: evidence is incomplete"
	case input.Status != ExecutionEnvelopeTypedDecisionCalibrationLSPBound:
		output.Status = ExecutionEnvelopeTypedDecisionCalibrationLSPError
		output.Code = executionEnvelopeTypedDecisionCalibrationLSPIntegrityCode
		output.Severity = "error"
		output.Message = "typed decision calibration status is invalid"
		output.FirstMismatch = "status"
		output.MissingStage = "typed_decision_calibration"
	case calibrationInputError(input) != "":
		output.Status = ExecutionEnvelopeTypedDecisionCalibrationLSPError
		output.Code = executionEnvelopeTypedDecisionCalibrationLSPIntegrityCode
		output.Severity = "error"
		output.Message = "typed decision calibration failed integrity validation"
		output.FirstMismatch = calibrationInputError(input)
		output.MissingStage = "typed_decision_calibration"
	default:
		output.Status = ExecutionEnvelopeTypedDecisionCalibrationLSPBound
		output.Code = executionEnvelopeTypedDecisionCalibrationLSPBoundCode
		output.Severity = "info"
		output.Message = "typed decision calibration is available for inspection"
		output.FirstMismatch = ""
		output.MissingStage = ""
	}
	output.EvidenceDigest = executionEnvelopeTypedDecisionCalibrationLSPDigest(output)
	return output
}

func calibrationInputError(input ExecutionEnvelopeTypedDecisionCalibrationLSPInput) string {
	switch {
	case strings.TrimSpace(input.SourceVersion) == "":
		return "source-version"
	case strings.TrimSpace(input.ContractVersion) == "":
		return "contract-version"
	case strings.TrimSpace(input.ModelRevision) == "":
		return "model-revision"
	case strings.TrimSpace(input.QuestionID) == "":
		return "question-id"
	case !validExternalTypedDecisionCalibrationDigest(input.SignalEvidenceDigest):
		return "signal-evidence-digest"
	case !validExternalTypedDecisionCalibrationDigest(input.CalibrationEvidenceDigest):
		return "calibration-evidence-digest"
	case !validExternalTypedDecisionCalibrationDigest(input.OutcomeDigest):
		return "outcome-digest"
	case !finiteCalibrationLSPUnit(input.SelectedProbability):
		return "selected-probability"
	case !finiteCalibrationLSPUnit(input.Confidence):
		return "confidence"
	case !input.OutcomeKnown:
		return "outcome"
	case input.WindowSize <= 0:
		return "sample-window"
	case !finiteCalibrationLSPUnit(input.Tolerance):
		return "tolerance"
	case !finiteCalibrationLSPUnit(input.AbsoluteError):
		return "absolute-error"
	default:
		target := 0.0
		if input.ObservedOutcome {
			target = 1
		}
		if math.Abs(math.Abs(input.SelectedProbability-target)-input.AbsoluteError) > 1e-9 {
			return "absolute-error"
		}
		if input.WithinTolerance != (input.AbsoluteError <= input.Tolerance) {
			return "tolerance"
		}
		return ""
	}
}

func finiteCalibrationLSPUnit(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value >= 0 && value <= 1
}

func validExternalTypedDecisionCalibrationDigest(value string) bool {
	const prefix = "sha256:"
	if !strings.HasPrefix(value, prefix) || len(value) != len(prefix)+64 {
		return false
	}
	_, err := hex.DecodeString(value[len(prefix):])
	return err == nil
}

func executionEnvelopeTypedDecisionCalibrationLSPDigest(
	projection ExecutionEnvelopeTypedDecisionCalibrationLSPProjection,
) string {
	return digestString(fmt.Sprintf(
		"gooo-typed-decision-calibration-lsp|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%0.9f|%0.9f|%t|%t|%0.9f|%0.9f|%t|%d|%s|%s|%t|%t",
		projection.Status,
		projection.Code,
		projection.Severity,
		projection.Message,
		projection.SourceVersion,
		projection.ContractVersion,
		projection.ModelRevision,
		projection.QuestionID,
		projection.SignalEvidenceDigest,
		projection.CalibrationEvidenceDigest,
		projection.SelectedProbability,
		projection.Confidence,
		projection.OutcomeKnown,
		projection.ObservedOutcome,
		projection.AbsoluteError,
		projection.Tolerance,
		projection.WithinTolerance,
		projection.WindowSize,
		projection.FirstMismatch,
		projection.MissingStage,
		projection.NonExecuting,
		projection.NonAuthorizing,
	))
}

func (projection ExecutionEnvelopeTypedDecisionCalibrationLSPProjection) Validate() error {
	switch projection.Status {
	case ExecutionEnvelopeTypedDecisionCalibrationLSPBound,
		ExecutionEnvelopeTypedDecisionCalibrationLSPUnknown,
		ExecutionEnvelopeTypedDecisionCalibrationLSPError:
	default:
		return fmt.Errorf("invalid typed decision calibration LSP status %q", projection.Status)
	}
	if !projection.NonExecuting || !projection.NonAuthorizing {
		return fmt.Errorf("typed decision calibration LSP crossed a capability boundary")
	}
	if projection.EvidenceDigest != executionEnvelopeTypedDecisionCalibrationLSPDigest(projection) {
		return fmt.Errorf("typed decision calibration LSP evidence digest mismatch")
	}
	switch projection.Status {
	case ExecutionEnvelopeTypedDecisionCalibrationLSPBound:
		if projection.Code != executionEnvelopeTypedDecisionCalibrationLSPBoundCode ||
			projection.Severity != "info" ||
			projection.FirstMismatch != "" ||
			projection.MissingStage != "" ||
			calibrationInputError(ExecutionEnvelopeTypedDecisionCalibrationLSPInput{
				Status:                    ExecutionEnvelopeTypedDecisionCalibrationLSPBound,
				SourceVersion:             projection.SourceVersion,
				ContractVersion:           projection.ContractVersion,
				ModelRevision:             projection.ModelRevision,
				QuestionID:                projection.QuestionID,
				SignalEvidenceDigest:      projection.SignalEvidenceDigest,
				CalibrationEvidenceDigest: projection.CalibrationEvidenceDigest,
				OutcomeDigest:             projection.OutcomeDigest,
				SelectedProbability:       projection.SelectedProbability,
				Confidence:                projection.Confidence,
				OutcomeKnown:              projection.OutcomeKnown,
				ObservedOutcome:           projection.ObservedOutcome,
				AbsoluteError:             projection.AbsoluteError,
				Tolerance:                 projection.Tolerance,
				WithinTolerance:           projection.WithinTolerance,
				WindowSize:                projection.WindowSize,
				NonExecuting:              projection.NonExecuting,
				NonAuthorizing:            projection.NonAuthorizing,
			}) != "" {
			return fmt.Errorf("bound typed decision calibration LSP is incomplete")
		}
	case ExecutionEnvelopeTypedDecisionCalibrationLSPUnknown:
		if projection.Code != executionEnvelopeTypedDecisionCalibrationLSPUnknownCode ||
			projection.MissingStage == "" ||
			projection.FirstMismatch == "" {
			return fmt.Errorf("unknown typed decision calibration LSP lost its first mismatch")
		}
	case ExecutionEnvelopeTypedDecisionCalibrationLSPError:
		if projection.MissingStage == "" ||
			(projection.Code != executionEnvelopeTypedDecisionCalibrationLSPIntegrityCode &&
				projection.Code != executionEnvelopeTypedDecisionCalibrationLSPBoundaryCode) {
			return fmt.Errorf("error typed decision calibration LSP lost its failure stage")
		}
	}
	return nil
}