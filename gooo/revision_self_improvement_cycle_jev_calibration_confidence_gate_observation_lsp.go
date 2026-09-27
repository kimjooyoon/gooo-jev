package gooo

// JEVCalibrationConfidenceGateObservationLSPStatus is the editor-facing
// severity for a confidence boundary. It is diagnostic only.
type JEVCalibrationConfidenceGateObservationLSPStatus string

const (
	JEVCalibrationConfidenceGateObservationLSPInformation JEVCalibrationConfidenceGateObservationLSPStatus = "Information"
	JEVCalibrationConfidenceGateObservationLSPHint        JEVCalibrationConfidenceGateObservationLSPStatus = "Hint"
	JEVCalibrationConfidenceGateObservationLSPError       JEVCalibrationConfidenceGateObservationLSPStatus = "Error"
)

// JEVCalibrationConfidenceGateObservationLSPInput mirrors the stable
// confidence-gate observation boundary without coupling the core to a runtime.
type JEVCalibrationConfidenceGateObservationLSPInput struct {
	Status               string
	Confidence           float64
	Threshold            float64
	EvidencePrefixDigest string
	SourceVersion        string
	MissingStageIndex    int
}

// JEVCalibrationConfidenceGateObservationLSP is a read-only projection for
// editor diagnostics. It never authorizes edits, execution, or improvement.
type JEVCalibrationConfidenceGateObservationLSP struct {
	Status               JEVCalibrationConfidenceGateObservationLSPStatus
	Code                 string
	Message              string
	Confidence           float64
	Threshold            float64
	EvidencePrefixDigest string
	SourceVersion        string
	MissingStageIndex    int
	IsReadOnly           bool
	CanEdit              bool
	CanExecute           bool
	CanAuthorize         bool
}

// ProjectJEVCalibrationConfidenceGateObservationLSP maps the explicit runtime
// boundary to diagnostics while preserving all supplied evidence fields.
func ProjectJEVCalibrationConfidenceGateObservationLSP(
	input JEVCalibrationConfidenceGateObservationLSPInput,
) JEVCalibrationConfidenceGateObservationLSP {
	projection := JEVCalibrationConfidenceGateObservationLSP{
		Status:               JEVCalibrationConfidenceGateObservationLSPError,
		Code:                 "jev.confidence_gate.unknown",
		Message:              "confidence-gate evidence is missing or invalid",
		Confidence:           input.Confidence,
		Threshold:            input.Threshold,
		EvidencePrefixDigest: input.EvidencePrefixDigest,
		SourceVersion:        input.SourceVersion,
		MissingStageIndex:    input.MissingStageIndex,
		IsReadOnly:           true,
		CanEdit:              false,
		CanExecute:           false,
		CanAuthorize:         false,
	}

	switch input.Status {
	case "AUTO":
		projection.Status = JEVCalibrationConfidenceGateObservationLSPInformation
		projection.Code = "jev.confidence_gate.auto"
		projection.Message = "confidence meets the configured routing threshold"
	case "REVIEW":
		projection.Status = JEVCalibrationConfidenceGateObservationLSPHint
		projection.Code = "jev.confidence_gate.review"
		projection.Message = "confidence is below the configured routing threshold; review is required"
	}

	return projection
}