package gooo

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type JEVCalibrationDecisionRouteStatus string

const (
	JEVCalibrationDecisionRouteAcceptCandidate JEVCalibrationDecisionRouteStatus = "ACCEPT_CANDIDATE"
	JEVCalibrationDecisionRouteReviewCandidate JEVCalibrationDecisionRouteStatus = "REVIEW_CANDIDATE"
	JEVCalibrationDecisionRouteAbstainCandidate JEVCalibrationDecisionRouteStatus = "ABSTAIN_CANDIDATE"
	JEVCalibrationDecisionRouteUnknown         JEVCalibrationDecisionRouteStatus = "UNKNOWN"
	JEVCalibrationDecisionRouteDeferred        JEVCalibrationDecisionRouteStatus = "DEFERRED"
)

type JEVCalibrationDecisionRouteInput struct {
	SourceVersion                string
	ContractVersion              string
	DeclarationDigest            string
	IRDigest                     string
	GeneratedDigest              string
	ReverseObservationDigest    string
	CalibrationDigest            string
	OutcomeWindowDigest          string
	ChoiceSetDigest              string
	CalibratedChoiceSetDigest    string
	OutcomeCount                 int
	CalibratedProbability        float64
	AcceptThreshold              float64
	ReviewThreshold              float64
	ProviderDeferred             bool
}

type JEVCalibrationDecisionRoute struct {
	Status                JEVCalibrationDecisionRouteStatus
	SourceVersion         string
	ContractVersion       string
	MissingStage          string
	MissingStageIndex     int
	EvidenceCoverage      int
	CalibratedProbability float64
	AcceptThreshold       float64
	ReviewThreshold       float64
	Reason                string
	DecisionDigest        string
	IsReadOnly            bool
	CanExecute            bool
	CanAuthorize          bool
}

func ProjectJEVCalibrationDecisionRoute(input JEVCalibrationDecisionRouteInput) JEVCalibrationDecisionRoute {
	route := JEVCalibrationDecisionRoute{
		Status:                JEVCalibrationDecisionRouteUnknown,
		SourceVersion:         input.SourceVersion,
		ContractVersion:       input.ContractVersion,
		MissingStageIndex:     -1,
		CalibratedProbability: input.CalibratedProbability,
		AcceptThreshold:       input.AcceptThreshold,
		ReviewThreshold:       input.ReviewThreshold,
		IsReadOnly:             true,
		CanExecute:             false,
		CanAuthorize:           false,
	}
	if input.SourceVersion == "" || input.ContractVersion == "" {
		route.MissingStage = "identity"
		route.Reason = "source or contract identity is missing"
		return finalizeJEVCalibrationDecisionRoute(route)
	}
	if input.ProviderDeferred {
		route.Status = JEVCalibrationDecisionRouteDeferred
		route.Reason = "decision producer is deferred"
		return finalizeJEVCalibrationDecisionRoute(route)
	}

	stages := []struct {
		name  string
		value string
	}{
		{name: "declaration", value: input.DeclarationDigest},
		{name: "ir", value: input.IRDigest},
		{name: "generated", value: input.GeneratedDigest},
		{name: "reverse_observation", value: input.ReverseObservationDigest},
		{name: "calibration", value: input.CalibrationDigest},
		{name: "outcome_window", value: input.OutcomeWindowDigest},
		{name: "choice_set", value: input.ChoiceSetDigest},
	}
	for index, stage := range stages {
		if stage.value == "" {
			route.MissingStage = stage.name
			route.MissingStageIndex = index
			route.EvidenceCoverage = index
			route.Reason = "first required evidence stage is missing"
			return finalizeJEVCalibrationDecisionRoute(route)
		}
	}
	route.EvidenceCoverage = len(stages)
	if input.CalibratedChoiceSetDigest != input.ChoiceSetDigest {
		route.MissingStage = "choice_set"
		route.MissingStageIndex = len(stages) - 1
		route.Reason = "calibration evidence was produced for a different choice set"
		return finalizeJEVCalibrationDecisionRoute(route)
	}
	if input.OutcomeCount <= 0 {
		route.MissingStage = "outcome_window"
		route.MissingStageIndex = 5
		route.Reason = "outcome-backed calibration evidence is absent or empty"
		return finalizeJEVCalibrationDecisionRoute(route)
	}
	if math.IsNaN(input.CalibratedProbability) || math.IsInf(input.CalibratedProbability, 0) ||
		input.CalibratedProbability < 0 || input.CalibratedProbability > 1 {
		route.Reason = "calibrated probability is outside [0,1]"
		return finalizeJEVCalibrationDecisionRoute(route)
	}
	if input.ReviewThreshold < 0 || input.ReviewThreshold > 1 ||
		input.AcceptThreshold < 0 || input.AcceptThreshold > 1 ||
		input.ReviewThreshold > input.AcceptThreshold {
		route.Reason = "routing thresholds are invalid"
		return finalizeJEVCalibrationDecisionRoute(route)
	}

	switch {
	case input.CalibratedProbability >= input.AcceptThreshold:
		route.Status = JEVCalibrationDecisionRouteAcceptCandidate
		route.Reason = "calibrated evidence supports an acceptance candidate"
	case input.CalibratedProbability >= input.ReviewThreshold:
		route.Status = JEVCalibrationDecisionRouteReviewCandidate
		route.Reason = "calibrated evidence requires review before any acceptance"
	default:
		route.Status = JEVCalibrationDecisionRouteAbstainCandidate
		route.Reason = "calibrated evidence is below the review threshold"
	}
	return finalizeJEVCalibrationDecisionRoute(route)
}

func finalizeJEVCalibrationDecisionRoute(route JEVCalibrationDecisionRoute) JEVCalibrationDecisionRoute {
	route.DecisionDigest = jevCalibrationDecisionRouteDigest(
		string(route.Status),
		route.SourceVersion,
		route.ContractVersion,
		route.MissingStage,
		strconv.Itoa(route.MissingStageIndex),
		strconv.Itoa(route.EvidenceCoverage),
		strconv.FormatFloat(route.CalibratedProbability, 'g', -1, 64),
		strconv.FormatFloat(route.AcceptThreshold, 'g', -1, 64),
		strconv.FormatFloat(route.ReviewThreshold, 'g', -1, 64),
		route.Reason,
	)
	return route
}

func jevCalibrationDecisionRouteDigest(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return fmt.Sprintf("sha256:%s", hex.EncodeToString(sum[:]))
}