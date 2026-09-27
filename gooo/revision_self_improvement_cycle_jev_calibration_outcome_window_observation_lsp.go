package gooo

type JEVCalibrationOutcomeWindowObservationLSPStatus string

const (
	JEVCalibrationOutcomeWindowObservationLSPBound    JEVCalibrationOutcomeWindowObservationLSPStatus = "BOUND"
	JEVCalibrationOutcomeWindowObservationLSPDeferred JEVCalibrationOutcomeWindowObservationLSPStatus = "DEFERRED"
	JEVCalibrationOutcomeWindowObservationLSPUnknown  JEVCalibrationOutcomeWindowObservationLSPStatus = "UNKNOWN"
)

// JEVCalibrationOutcomeWindowObservationLSPProjection is an
// inspection-only diagnostic view. It cannot edit, execute, or authorize.
type JEVCalibrationOutcomeWindowObservationLSPProjection struct {
	Status       JEVCalibrationOutcomeWindowObservationLSPStatus
	Code         string
	Severity     string
	Title        string
	TargetStage  string
	RecordDigest string
	Reason       string
	IsReadOnly   bool
	CanExecute   bool
	CanAuthorize bool
	Edits        []string
	Command      string
}

func ProjectJEVCalibrationOutcomeWindowObservationLSP(
	observation JEVCalibrationOutcomeWindowObservation,
) JEVCalibrationOutcomeWindowObservationLSPProjection {
	projection := JEVCalibrationOutcomeWindowObservationLSPProjection{
		Status:       JEVCalibrationOutcomeWindowObservationLSPUnknown,
		Code:         "jev.outcome_window.unknown",
		Severity:     "Error",
		Title:        "Outcome-window evidence is unresolved",
		TargetStage:  observation.TargetStage,
		RecordDigest: observation.RecordDigest,
		Reason:        observation.Reason,
		IsReadOnly:   true,
		CanExecute:   false,
		CanAuthorize: false,
	}

	if observation.RecordDigest == "" {
		projection.TargetStage = "record_digest"
		return projection
	}

	switch observation.Status {
	case JEVCalibrationOutcomeWindowObservationBound:
		projection.Status = JEVCalibrationOutcomeWindowObservationLSPBound
		projection.Code = "jev.outcome_window.bound"
		projection.Severity = "Information"
		projection.Title = "Outcome-window evidence is bound"
	case JEVCalibrationOutcomeWindowObservationDeferred:
		projection.Status = JEVCalibrationOutcomeWindowObservationLSPDeferred
		projection.Code = "jev.outcome_window.deferred"
		projection.Severity = "Hint"
		projection.Title = "Outcome-window evidence is deferred"
	}
	return projection
}
