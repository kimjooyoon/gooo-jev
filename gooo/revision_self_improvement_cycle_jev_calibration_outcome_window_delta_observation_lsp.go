package gooo

type JEVCalibrationOutcomeWindowDeltaObservationLSPStatus string

const (
	JEVCalibrationOutcomeWindowDeltaLSPBound    JEVCalibrationOutcomeWindowDeltaObservationLSPStatus = "BOUND"
	JEVCalibrationOutcomeWindowDeltaLSPDeferred JEVCalibrationOutcomeWindowDeltaObservationLSPStatus = "DEFERRED"
	JEVCalibrationOutcomeWindowDeltaLSPUnknown  JEVCalibrationOutcomeWindowDeltaObservationLSPStatus = "UNKNOWN"
)

// JEVCalibrationOutcomeWindowDeltaObservationLSPProjection is an
// inspection-only view of a signed metric. It cannot edit, execute, or authorize.
type JEVCalibrationOutcomeWindowDeltaObservationLSPProjection struct {
	Status       JEVCalibrationOutcomeWindowDeltaObservationLSPStatus
	Code         string
	Severity     string
	Title        string
	TargetStage  string
	Delta        int
	RecordDigest string
	Reason       string
	IsReadOnly   bool
	CanExecute   bool
	CanAuthorize bool
	Edits        []string
	Command      string
}

func ProjectJEVCalibrationOutcomeWindowDeltaObservationLSP(
	observation JEVCalibrationOutcomeWindowDeltaObservation,
) JEVCalibrationOutcomeWindowDeltaObservationLSPProjection {
	projection := JEVCalibrationOutcomeWindowDeltaObservationLSPProjection{
		Status:       JEVCalibrationOutcomeWindowDeltaLSPUnknown,
		Code:         "jev.outcome_window_delta.unknown",
		Severity:     "Error",
		Title:        "Signed outcome-window delta is unresolved",
		TargetStage:  observation.TargetStage,
		Delta:        observation.Delta,
		RecordDigest: observation.RecordDigest,
		Reason:       observation.Reason,
		IsReadOnly:   true,
		CanExecute:   false,
		CanAuthorize: false,
	}

	if observation.RecordDigest == "" {
		projection.TargetStage = "record_digest"
		return projection
	}

	switch observation.Status {
	case JEVCalibrationOutcomeWindowDeltaBound:
		projection.Status = JEVCalibrationOutcomeWindowDeltaLSPBound
		projection.Code = "jev.outcome_window_delta.bound"
		projection.Severity = "Information"
		projection.Title = "Signed outcome-window delta is recorded"
	case JEVCalibrationOutcomeWindowDeltaDeferred:
		projection.Status = JEVCalibrationOutcomeWindowDeltaLSPDeferred
		projection.Code = "jev.outcome_window_delta.deferred"
		projection.Severity = "Hint"
		projection.Title = "Signed outcome-window delta is deferred"
	}
	return projection
}
