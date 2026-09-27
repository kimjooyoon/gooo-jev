package gooo

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type JEVCalibrationDecisionRoutePlanBoundaryStatus string

const (
	JEVCalibrationDecisionRoutePlanProposed  JEVCalibrationDecisionRoutePlanBoundaryStatus = "PROPOSED"
	JEVCalibrationDecisionRoutePlanDeferred  JEVCalibrationDecisionRoutePlanBoundaryStatus = "DEFERRED"
	JEVCalibrationDecisionRoutePlanUnknown   JEVCalibrationDecisionRoutePlanBoundaryStatus = "UNKNOWN"
)

type JEVCalibrationDecisionRoutePlanBoundaryInput struct {
	SourceVersion             string
	ContractVersion           string
	RouteStatus               JEVCalibrationDecisionRouteStatus
	RouteDecisionDigest       string
	ReverseObservationStatus  JEVCalibrationDecisionRouteReverseStatus
	ObservationDigest         string
	PlanDigest                string
}

type JEVCalibrationDecisionRoutePlanBoundary struct {
	Status                  JEVCalibrationDecisionRoutePlanBoundaryStatus
	SourceVersion           string
	ContractVersion         string
	RouteDecisionDigest     string
	ObservationDigest       string
	PlanDigest              string
	TargetStage             string
	Reason                  string
	BoundaryDigest          string
	IsReadOnly              bool
	CanExecute              bool
	CanAuthorize            bool
}

func ProjectJEVCalibrationDecisionRoutePlanBoundary(
	input JEVCalibrationDecisionRoutePlanBoundaryInput,
) JEVCalibrationDecisionRoutePlanBoundary {
	boundary := JEVCalibrationDecisionRoutePlanBoundary{
		Status:              JEVCalibrationDecisionRoutePlanUnknown,
		SourceVersion:       input.SourceVersion,
		ContractVersion:     input.ContractVersion,
		RouteDecisionDigest: input.RouteDecisionDigest,
		ObservationDigest:   input.ObservationDigest,
		PlanDigest:          input.PlanDigest,
		TargetStage:         "plan_application",
		IsReadOnly:          true,
		CanExecute:          false,
		CanAuthorize:        false,
	}
	switch {
	case input.SourceVersion == "" || input.ContractVersion == "":
		boundary.TargetStage = "identity"
		boundary.Reason = "source or contract identity is missing"
	case input.RouteDecisionDigest == "":
		boundary.TargetStage = "route_decision_digest"
		boundary.Reason = "route decision evidence is missing"
	case input.ObservationDigest == "":
		boundary.TargetStage = "reverse_observation_digest"
		boundary.Reason = "reverse observation evidence is missing"
	case input.PlanDigest == "":
		boundary.TargetStage = "plan_digest"
		boundary.Reason = "plan digest is missing"
	case input.ReverseObservationStatus == JEVCalibrationDecisionRouteReverseUnknown:
		boundary.TargetStage = "reverse_observation"
		boundary.Reason = "reverse observation remains unresolved"
	case input.ReverseObservationStatus == JEVCalibrationDecisionRouteReverseDeferred ||
		input.RouteStatus == JEVCalibrationDecisionRouteDeferred:
		boundary.Status = JEVCalibrationDecisionRoutePlanDeferred
		boundary.Reason = "route producer or reverse observation is deferred"
	case input.RouteStatus == JEVCalibrationDecisionRouteUnknown:
		boundary.TargetStage = "route_status"
		boundary.Reason = "route status remains unresolved"
	case input.ReverseObservationStatus != JEVCalibrationDecisionRouteReverseBound:
		boundary.TargetStage = "reverse_observation_status"
		boundary.Reason = "reverse observation status is not bound"
	case input.RouteStatus == JEVCalibrationDecisionRouteAcceptCandidate ||
		input.RouteStatus == JEVCalibrationDecisionRouteReviewCandidate ||
		input.RouteStatus == JEVCalibrationDecisionRouteAbstainCandidate:
		boundary.Status = JEVCalibrationDecisionRoutePlanProposed
		boundary.Reason = "form a plan proposal from bound route evidence"
	default:
		boundary.TargetStage = "route_status"
		boundary.Reason = "route status is not recognized"
	}
	boundary.BoundaryDigest = jevCalibrationDecisionRoutePlanDigest(
		string(boundary.Status),
		boundary.SourceVersion,
		boundary.ContractVersion,
		boundary.RouteDecisionDigest,
		boundary.ObservationDigest,
		boundary.PlanDigest,
		boundary.TargetStage,
		boundary.Reason,
	)
	return boundary
}

func jevCalibrationDecisionRoutePlanDigest(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(sum[:]))
}