package gooo

type JEVCalibrationDecisionRoutingPlanObservationLSPStatus string

const (
	JEVCalibrationDecisionRoutingPlanObservationLSPBound    JEVCalibrationDecisionRoutingPlanObservationLSPStatus = "BOUND"
	JEVCalibrationDecisionRoutingPlanObservationLSPDeferred JEVCalibrationDecisionRoutingPlanObservationLSPStatus = "DEFERRED"
	JEVCalibrationDecisionRoutingPlanObservationLSPUnknown  JEVCalibrationDecisionRoutingPlanObservationLSPStatus = "UNKNOWN"
)

// JEVCalibrationDecisionRoutingPlanObservationLSPProjection is an
// inspection-only diagnostic view. It cannot edit, execute, or authorize.
type JEVCalibrationDecisionRoutingPlanObservationLSPProjection struct {
	Status            JEVCalibrationDecisionRoutingPlanObservationLSPStatus
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

func ProjectJEVCalibrationDecisionRoutingPlanObservationLSP(
	observation JEVCalibrationDecisionRoutingPlanObservation,
) JEVCalibrationDecisionRoutingPlanObservationLSPProjection {
	projection := JEVCalibrationDecisionRoutingPlanObservationLSPProjection{
		Status:            JEVCalibrationDecisionRoutingPlanObservationLSPUnknown,
		Code:              "jev.plan_observation.unknown",
		Severity:          "Error",
		Title:             "Plan provenance evidence is unresolved",
		TargetStage:       observation.TargetStage,
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
	case JEVCalibrationDecisionRoutingPlanObservationBound:
		projection.Status = JEVCalibrationDecisionRoutingPlanObservationLSPBound
		projection.Code = "jev.plan_observation.bound"
		projection.Severity = "Information"
		projection.Title = "Plan provenance evidence is bound"
	case JEVCalibrationDecisionRoutingPlanObservationDeferred:
		projection.Status = JEVCalibrationDecisionRoutingPlanObservationLSPDeferred
		projection.Code = "jev.plan_observation.deferred"
		projection.Severity = "Hint"
		projection.Title = "Plan provenance evidence is deferred"
	}
	return projection
}