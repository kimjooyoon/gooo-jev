package gooo

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type JEVCalibrationDecisionRouteReverseStatus string

const (
	JEVCalibrationDecisionRouteReverseBound    JEVCalibrationDecisionRouteReverseStatus = "BOUND"
	JEVCalibrationDecisionRouteReverseUnknown  JEVCalibrationDecisionRouteReverseStatus = "UNKNOWN"
	JEVCalibrationDecisionRouteReverseDeferred JEVCalibrationDecisionRouteReverseStatus = "DEFERRED"
)

// JEVCalibrationDecisionRouteReverseObservation records whether an observed
// route exactly matches the evidence-bound route.
type JEVCalibrationDecisionRouteReverseObservation struct {
	Status                  JEVCalibrationDecisionRouteReverseStatus
	SourceVersion           string
	ContractVersion         string
	ExpectedStatus          JEVCalibrationDecisionRouteStatus
	ObservedStatus          string
	ExpectedDecisionDigest  string
	ObservedDecisionDigest  string
	Reason                  string
	ObservationDigest       string
	IsReadOnly              bool
	CanExecute              bool
	CanAuthorize            bool
}

// ObserveJEVCalibrationDecisionRouteReverseObservation compares generated or
// projected output with the route that produced it.
func ObserveJEVCalibrationDecisionRouteReverseObservation(
	expected JEVCalibrationDecisionRoute,
	observedSourceVersion string,
	observedContractVersion string,
	observedStatus string,
	observedDecisionDigest string,
) JEVCalibrationDecisionRouteReverseObservation {
	observation := JEVCalibrationDecisionRouteReverseObservation{
		Status:                 JEVCalibrationDecisionRouteReverseUnknown,
		SourceVersion:          expected.SourceVersion,
		ContractVersion:       expected.ContractVersion,
		ExpectedStatus:         expected.Status,
		ObservedStatus:         observedStatus,
		ExpectedDecisionDigest: expected.DecisionDigest,
		ObservedDecisionDigest: observedDecisionDigest,
		IsReadOnly:             true,
		CanExecute:             false,
		CanAuthorize:           false,
	}
	switch {
	case expected.SourceVersion == "" || expected.ContractVersion == "":
		observation.Reason = "expected route identity is missing"
	case observedSourceVersion == "" || observedContractVersion == "":
		observation.Reason = "observed route identity is missing"
	case observedSourceVersion != expected.SourceVersion || observedContractVersion != expected.ContractVersion:
		observation.Reason = "observed source or contract identity changed"
	case expected.DecisionDigest == "" || observedDecisionDigest == "":
		observation.Reason = "route decision evidence is missing"
	case observedStatus != string(expected.Status):
		observation.Reason = "observed route status changed"
	case observedDecisionDigest != expected.DecisionDigest:
		observation.Reason = "observed route digest changed"
	case expected.Status == JEVCalibrationDecisionRouteDeferred:
		observation.Status = JEVCalibrationDecisionRouteReverseDeferred
		observation.Reason = "route producer remains deferred"
	case expected.Status == JEVCalibrationDecisionRouteUnknown:
		observation.Reason = "expected route remains unresolved"
	default:
		observation.Status = JEVCalibrationDecisionRouteReverseBound
		observation.Reason = "observed route matches source, contract, status, and digest"
	}
	observation.ObservationDigest = jevCalibrationDecisionRouteReverseDigest(
		string(observation.Status),
		observation.SourceVersion,
		observation.ContractVersion,
		string(observation.ExpectedStatus),
		observation.ObservedStatus,
		observation.ExpectedDecisionDigest,
		observation.ObservedDecisionDigest,
		observation.Reason,
	)
	return observation
}

func jevCalibrationDecisionRouteReverseDigest(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(sum[:]))
}