package gooo

type JEVCalibrationDecisionRouteReverseLSPStatus string

const (
	JEVCalibrationDecisionRouteReverseLSPBound    JEVCalibrationDecisionRouteReverseLSPStatus = "BOUND"
	JEVCalibrationDecisionRouteReverseLSPDeferred JEVCalibrationDecisionRouteReverseLSPStatus = "DEFERRED"
	JEVCalibrationDecisionRouteReverseLSPUnknown  JEVCalibrationDecisionRouteReverseLSPStatus = "UNKNOWN"
)

// JEVCalibrationDecisionRouteReverseLSPProjection is an inspection-only
// diagnostic projection of reverse-observed route evidence.
type JEVCalibrationDecisionRouteReverseLSPProjection struct {
	Status            JEVCalibrationDecisionRouteReverseLSPStatus
	Code              string
	Severity          string
	Title             string
	TargetStage       string
	ObservationDigest string
	Reason            string
	IsReadOnly        bool
	CanExecute        bool
	CanAuthorize      bool
	Edits             []string
	Command           string
}

func ProjectJEVCalibrationDecisionRouteReverseLSP(
	observation JEVCalibrationDecisionRouteReverseObservation,
) JEVCalibrationDecisionRouteReverseLSPProjection {
	projection := JEVCalibrationDecisionRouteReverseLSPProjection{
		Status:            JEVCalibrationDecisionRouteReverseLSPUnknown,
		Code:              "jev.route.reverse.unknown",
		Severity:          "Error",
		Title:             "Reverse route evidence is unresolved",
		ObservationDigest: observation.ObservationDigest,
		Reason:            observation.Reason,
		IsReadOnly:        true,
		CanExecute:        false,
		CanAuthorize:      false,
	}
	if observation.ObservationDigest == "" {
		projection.TargetStage = "observation_digest"
		return projection
	}

	switch observation.Status {
	case JEVCalibrationDecisionRouteReverseBound:
		projection.Status = JEVCalibrationDecisionRouteReverseLSPBound
		projection.Code = "jev.route.reverse.bound"
		projection.Severity = "Information"
		projection.Title = "Reverse route evidence is bound"
	case JEVCalibrationDecisionRouteReverseDeferred:
		projection.Status = JEVCalibrationDecisionRouteReverseLSPDeferred
		projection.Code = "jev.route.reverse.deferred"
		projection.Severity = "Hint"
		projection.Title = "Reverse route evidence is deferred"
	default:
		projection.TargetStage = "reverse_observation"
	}
	return projection
}