package gooo

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type JEVCalibrationDecisionRoutingPlanObservationStatus string

const (
	JEVCalibrationDecisionRoutingPlanObservationBound    JEVCalibrationDecisionRoutingPlanObservationStatus = "BOUND"
	JEVCalibrationDecisionRoutingPlanObservationDeferred JEVCalibrationDecisionRoutingPlanObservationStatus = "DEFERRED"
	JEVCalibrationDecisionRoutingPlanObservationUnknown  JEVCalibrationDecisionRoutingPlanObservationStatus = "UNKNOWN"
)

// JEVCalibrationDecisionRoutingPlanObservationInput is the evidence chain
// required to observe a plan application boundary without authorizing it.
type JEVCalibrationDecisionRoutingPlanObservationInput struct {
	SourceVersion             string
	ContractVersion           string
	DeclarationDigest         string
	IRDigest                  string
	GeneratedDigest           string
	ReverseObservationStatus  string
	ReverseObservationDigest  string
	ProposalStatus            string
	PlanDigest                string
	ExpectedPlanDigest        string
	OutcomeWindowDigest       string
	OutcomeCount              int
	ProducerDeferred          bool
}

// JEVCalibrationDecisionRoutingPlanObservation is a read-only provenance
// result. It describes evidence alignment, not execution or authorization.
type JEVCalibrationDecisionRoutingPlanObservation struct {
	Status                   JEVCalibrationDecisionRoutingPlanObservationStatus
	SourceVersion            string
	ContractVersion          string
	DeclarationDigest        string
	IRDigest                 string
	GeneratedDigest          string
	ReverseObservationDigest string
	PlanDigest               string
	OutcomeWindowDigest      string
	OutcomeCount             int
	TargetStage              string
	Reason                   string
	ObservationDigest        string
	IsReadOnly               bool
	CanExecute               bool
	CanAuthorize             bool
	Edits                    []string
	Command                  string
}

func ObserveJEVCalibrationDecisionRoutingPlan(
	input JEVCalibrationDecisionRoutingPlanObservationInput,
) JEVCalibrationDecisionRoutingPlanObservation {
	observation := JEVCalibrationDecisionRoutingPlanObservation{
		Status:                   JEVCalibrationDecisionRoutingPlanObservationUnknown,
		SourceVersion:            input.SourceVersion,
		ContractVersion:          input.ContractVersion,
		DeclarationDigest:        input.DeclarationDigest,
		IRDigest:                 input.IRDigest,
		GeneratedDigest:          input.GeneratedDigest,
		ReverseObservationDigest: input.ReverseObservationDigest,
		PlanDigest:               input.PlanDigest,
		OutcomeWindowDigest:      input.OutcomeWindowDigest,
		OutcomeCount:             input.OutcomeCount,
		TargetStage:              "plan_application_observation",
		IsReadOnly:               true,
		CanExecute:               false,
		CanAuthorize:             false,
	}

	switch {
	case input.SourceVersion == "" || input.ContractVersion == "":
		observation.TargetStage = "identity"
		observation.Reason = "source or contract identity is missing"
	case input.DeclarationDigest == "":
		observation.TargetStage = "declaration"
		observation.Reason = "declaration evidence is missing"
	case input.IRDigest == "":
		observation.TargetStage = "ir"
		observation.Reason = "IR evidence is missing"
	case input.GeneratedDigest == "":
		observation.TargetStage = "generated"
		observation.Reason = "generated evidence is missing"
	case input.ReverseObservationDigest == "":
		observation.TargetStage = "reverse_observation"
		observation.Reason = "reverse observation evidence is missing"
	case input.PlanDigest == "":
		observation.TargetStage = "plan_digest"
		observation.Reason = "plan digest is missing"
	case input.ExpectedPlanDigest == "":
		observation.TargetStage = "expected_plan_digest"
		observation.Reason = "expected plan digest is missing"
	case input.OutcomeWindowDigest == "":
		observation.TargetStage = "outcome_window"
		observation.Reason = "outcome window evidence is missing"
	case input.ProducerDeferred:
		observation.Status = JEVCalibrationDecisionRoutingPlanObservationDeferred
		observation.Reason = "plan observation producer is deferred"
	case input.ReverseObservationStatus == "DEFERRED":
		observation.Status = JEVCalibrationDecisionRoutingPlanObservationDeferred
		observation.Reason = "reverse observation is deferred"
	case input.ProposalStatus != "PROPOSED":
		observation.TargetStage = "proposal_status"
		observation.Reason = "plan proposal status is not PROPOSED"
	case input.PlanDigest != input.ExpectedPlanDigest:
		observation.TargetStage = "plan_digest_match"
		observation.Reason = "plan digest does not match expected digest"
	case input.ReverseObservationStatus == "UNKNOWN":
		observation.TargetStage = "reverse_observation_status"
		observation.Reason = "reverse observation remains unresolved"
	case input.OutcomeCount <= 0:
		observation.TargetStage = "outcome_count"
		observation.Reason = "outcome window has no positive observations"
	case input.ReverseObservationStatus == "BOUND":
		observation.Status = JEVCalibrationDecisionRoutingPlanObservationBound
		observation.Reason = "declaration, IR, generated, reverse, plan, and outcome evidence are aligned"
	default:
		observation.TargetStage = "reverse_observation_status"
		observation.Reason = "reverse observation status is not recognized"
	}

	observation.ObservationDigest = digestJEVCalibrationDecisionRoutingPlanObservation(
		string(observation.Status),
		observation.SourceVersion,
		observation.ContractVersion,
		observation.DeclarationDigest,
		observation.IRDigest,
		observation.GeneratedDigest,
		observation.ReverseObservationDigest,
		observation.PlanDigest,
		input.ExpectedPlanDigest,
		observation.OutcomeWindowDigest,
		fmt.Sprintf("%d", observation.OutcomeCount),
		observation.TargetStage,
		observation.Reason,
	)
	return observation
}

func digestJEVCalibrationDecisionRoutingPlanObservation(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(sum[:]))
}