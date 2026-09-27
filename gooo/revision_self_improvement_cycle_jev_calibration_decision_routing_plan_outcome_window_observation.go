package gooo

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationStatus string

const (
	JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationBound    JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationStatus = "BOUND"
	JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationDeferred JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationStatus = "DEFERRED"
	JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationUnknown  JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationStatus = "UNKNOWN"
)

// JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationInput links plan
// provenance to an outcome window without mutating metrics or granting authority.
type JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationInput struct {
	SourceVersion                        string
	ContractVersion                      string
	PlanObservationDigest                string
	ExpectedPlanObservationDigest        string
	OutcomeWindowObservationDigest       string
	ExpectedOutcomeWindowObservationDigest string
	ProducerDeferred                     bool
}

// JEVCalibrationDecisionRoutingPlanOutcomeWindowObservation is a read-only
// provenance chain result, not an execution or authorization result.
type JEVCalibrationDecisionRoutingPlanOutcomeWindowObservation struct {
	Status                          JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationStatus
	SourceVersion                   string
	ContractVersion                 string
	PlanObservationDigest           string
	OutcomeWindowObservationDigest  string
	TargetStage                     string
	Reason                          string
	ChainDigest                     string
	IsReadOnly                      bool
	CanExecute                      bool
	CanAuthorize                    bool
	Edits                           []string
	Command                         string
}

func ObserveJEVCalibrationDecisionRoutingPlanOutcomeWindow(
	input JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationInput,
) JEVCalibrationDecisionRoutingPlanOutcomeWindowObservation {
	observation := JEVCalibrationDecisionRoutingPlanOutcomeWindowObservation{
		Status:                         JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationUnknown,
		SourceVersion:                  input.SourceVersion,
		ContractVersion:                input.ContractVersion,
		PlanObservationDigest:          input.PlanObservationDigest,
		OutcomeWindowObservationDigest: input.OutcomeWindowObservationDigest,
		TargetStage:                    "plan_outcome_window",
		IsReadOnly:                     true,
		CanExecute:                     false,
		CanAuthorize:                   false,
	}

	switch {
	case input.SourceVersion == "" || input.ContractVersion == "":
		observation.TargetStage = "identity"
		observation.Reason = "source or contract identity is missing"
	case input.PlanObservationDigest == "":
		observation.TargetStage = "plan_observation_digest"
		observation.Reason = "plan observation digest is missing"
	case input.ExpectedPlanObservationDigest == "":
		observation.TargetStage = "expected_plan_observation_digest"
		observation.Reason = "expected plan observation digest is missing"
	case input.OutcomeWindowObservationDigest == "":
		observation.TargetStage = "outcome_window_observation_digest"
		observation.Reason = "outcome-window observation digest is missing"
	case input.ExpectedOutcomeWindowObservationDigest == "":
		observation.TargetStage = "expected_outcome_window_observation_digest"
		observation.Reason = "expected outcome-window observation digest is missing"
	case input.ProducerDeferred:
		observation.Status = JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationDeferred
		observation.Reason = "plan-outcome producer is deferred"
	case input.PlanObservationDigest != input.ExpectedPlanObservationDigest:
		observation.TargetStage = "plan_observation_digest_match"
		observation.Reason = "plan observation digest does not match expected digest"
	case input.OutcomeWindowObservationDigest != input.ExpectedOutcomeWindowObservationDigest:
		observation.TargetStage = "outcome_window_observation_digest_match"
		observation.Reason = "outcome-window observation digest does not match expected digest"
	default:
		observation.Status = JEVCalibrationDecisionRoutingPlanOutcomeWindowObservationBound
		observation.Reason = "plan observation and outcome-window observation are aligned"
	}

	observation.ChainDigest = digestJEVCalibrationDecisionRoutingPlanOutcomeWindow(
		string(observation.Status),
		observation.SourceVersion,
		observation.ContractVersion,
		observation.PlanObservationDigest,
		input.ExpectedPlanObservationDigest,
		observation.OutcomeWindowObservationDigest,
		input.ExpectedOutcomeWindowObservationDigest,
		observation.TargetStage,
		observation.Reason,
	)
	return observation
}

func digestJEVCalibrationDecisionRoutingPlanOutcomeWindow(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(sum[:]))
}
