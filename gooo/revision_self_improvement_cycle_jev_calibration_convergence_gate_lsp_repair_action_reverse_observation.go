package gooo

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionReverseObservation struct {
	Status                       string
	MissingStage                 string
	MissingStageIndex            int
	ActionDigest                 string
	GateDecision                string
	SuggestedAction              string
	ProjectionObservationDigest string
	EvidencePrefixDigest        string
	SourceObservationDigest     string
	ReverseObservationDigest    string
	NonExecuting                bool
	NonAuthorizing              bool
}

func ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction(
	action RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction,
) (RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionReverseObservation, error) {
	observation := RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionReverseObservation{
		Status:           "UNKNOWN",
		MissingStage:     "jev-calibration-convergence-gate-lsp-repair-action-reverse-observation",
		MissingStageIndex: action.MissingStageIndex,
		ActionDigest:     action.ActionDigest,
		NonExecuting:     true,
		NonAuthorizing:   true,
	}
	setDigest := func() {
		observation.ReverseObservationDigest = digestRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionReverseObservation(observation)
	}
	if err := action.Validate(); err != nil {
		setDigest()
		return observation, err
	}
	observation.Status = action.Status
	observation.MissingStage = action.MissingStage
	observation.GateDecision = action.GateDecision
	observation.SuggestedAction = action.SuggestedAction
	observation.ProjectionObservationDigest = action.ProjectionObservationDigest
	observation.EvidencePrefixDigest = action.EvidencePrefixDigest
	observation.SourceObservationDigest = action.SourceObservationDigest
	setDigest()
	if err := observation.Validate(); err != nil {
		observation.Status = "UNKNOWN"
		observation.MissingStage = "jev-calibration-convergence-gate-lsp-repair-action-reverse-observation"
		setDigest()
		return observation, err
	}
	return observation, nil
}

func (o RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionReverseObservation) Validate() error {
	if o.Status == "" {
		return errors.New("reverse observation status is empty")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return errors.New("bound reverse observation has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return errors.New("unknown reverse observation has no missing stage")
	}
	if o.MissingStageIndex < 0 {
		return errors.New("reverse observation missing stage index is negative")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return errors.New("reverse observation must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionReverseObservation(o) != o.ReverseObservationDigest {
		return errors.New("reverse observation digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionReverseObservation(
	observation RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionReverseObservation,
) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{
		observation.Status,
		observation.MissingStage,
		strconv.Itoa(observation.MissingStageIndex),
		observation.ActionDigest,
		observation.GateDecision,
		observation.SuggestedAction,
		observation.ProjectionObservationDigest,
		observation.EvidencePrefixDigest,
		observation.SourceObservationDigest,
		strconv.FormatBool(observation.NonExecuting),
		strconv.FormatBool(observation.NonAuthorizing),
	}, ":")))
	return hex.EncodeToString(sum[:])
}