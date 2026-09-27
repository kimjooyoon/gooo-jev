package gooo

type JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPStatus string

const (
	JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPBound    JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPStatus = "BOUND"
	JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPDeferred JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPStatus = "DEFERRED"
	JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPUnknown  JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPStatus = "UNKNOWN"
)

// JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPProjection is
// an inspection-only diagnostic view. It cannot edit, execute, or authorize.
type JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPProjection struct {
	Status       JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPStatus
	Code         string
	Severity     string
	Title        string
	TargetStage  string
	ChainDigest  string
	Reason       string
	IsReadOnly   bool
	CanExecute   bool
	CanAuthorize bool
	Edits        []string
	Command      string
}

func ProjectJEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSP(
	observation JEVCalibrationDecisionRoutingPlanOutcomeWindowObservation,
) JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPProjection {
	projection := JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPProjection{
		Status:       JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPUnknown,
		Code:         "jev.plan_outcome_window.unknown",
		Severity:     "Error",
		Title:        "Plan and outcome-window provenance is unresolved",
		TargetStage:  observation.TargetStage,
		ChainDigest:  observation.ChainDigest,
		Reason:       observation.Reason,
		IsReadOnly:   true,
		CanExecute:   false,
		CanAuthorize: false,
	}

	if observation.ChainDigest == "" {
		projection.TargetStage = "chain_digest"
		return projection
	}

	switch observation.Status {
	case JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationBound:
		projection.Status = JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPBound
		projection.Code = "jev.plan_outcome_window.bound"
		projection.Severity = "Information"
		projection.Title = "Plan and outcome-window provenance is bound"
	case JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationDeferred:
		projection.Status = JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationLSPDeferred
		projection.Code = "jev.plan_outcome_window.deferred"
		projection.Severity = "Hint"
		projection.Title = "Plan and outcome-window provenance is deferred"
	}
	return projection
}
