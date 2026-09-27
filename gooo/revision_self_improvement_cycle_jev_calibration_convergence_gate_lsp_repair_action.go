package gooo

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
)

type RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionInput struct {
	ProjectionObservationDigest string
	EvidencePrefixDigest        string
	SourceObservationDigest     string
	ReverseObservationDigest    string
	GateDecision                string
	ProjectionSignal            string
	MissingStage                string
	MissingStageIndex           int
}

type RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction struct {
	Status                       string
	MissingStage                 string
	MissingStageIndex            int
	GateDecision                string
	ProjectionSignal             string
	SuggestedAction              string
	ActionKind                   string
	Title                        string
	Reason                       string
	ProjectionObservationDigest string
	EvidencePrefixDigest        string
	SourceObservationDigest     string
	ReverseObservationDigest    string
	ActionDigest                 string
	NonExecuting                bool
	NonAuthorizing               bool
}

func ProjectRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction(
	input RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairActionInput,
) (RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction, error) {
	action := RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction{
		Status:                       "UNKNOWN",
		MissingStage:                 input.MissingStage,
		MissingStageIndex:            input.MissingStageIndex,
		GateDecision:                input.GateDecision,
		ProjectionSignal:             input.ProjectionSignal,
		ProjectionObservationDigest: input.ProjectionObservationDigest,
		EvidencePrefixDigest:         input.EvidencePrefixDigest,
		SourceObservationDigest:      input.SourceObservationDigest,
		ReverseObservationDigest:     input.ReverseObservationDigest,
		ActionKind:                   "provenance-inspection",
		NonExecuting:                 true,
		NonAuthorizing:               true,
	}
	setDigest := func() {
		action.ActionDigest = digestRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction(action)
	}
	if action.MissingStageIndex < 0 {
		setDigest()
		return action, errors.New("missing stage index must not be negative")
	}
	if action.MissingStage == "" {
		action.MissingStage = "jev-calibration-convergence-gate-lsp-repair-action"
	}
	if action.ProjectionObservationDigest == "" {
		setDigest()
		return action, nil
	}
	if action.EvidencePrefixDigest == "" {
		action.MissingStage = "jev-calibration-convergence-gate-lsp-projection-evidence-prefix"
		setDigest()
		return action, nil
	}
	if action.ReverseObservationDigest == "" {
		action.MissingStage = "jev-calibration-convergence-gate-lsp-reverse-observation"
		setDigest()
		return action, nil
	}
	switch action.GateDecision {
	case "jev-calibration-convergence-converged":
		action.Status = "BOUND"
		action.SuggestedAction = "observe-convergence"
		action.Title = "Observe comparable convergence"
		action.Reason = "calibration improved on a stable choice set"
	case "jev-calibration-convergence-regressed":
		action.Status = "BOUND"
		action.SuggestedAction = "review-regression"
		action.Title = "Review calibration regression"
		action.Reason = "calibration regressed on a stable choice set"
	case "jev-calibration-convergence-defer":
		action.Status = "BOUND"
		action.SuggestedAction = "wait-for-comparable-evidence"
		action.Title = "Wait for comparable evidence"
		action.Reason = "choice-set or evidence lineage changed"
	default:
		action.MissingStage = "jev-calibration-convergence-gate-decision"
	}
	setDigest()
	if err := action.Validate(); err != nil {
		action.Status = "UNKNOWN"
		action.MissingStage = "jev-calibration-convergence-gate-lsp-repair-action"
		setDigest()
		return action, err
	}
	return action, nil
}

func (a RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction) Validate() error {
	if a.Status == "" {
		return errors.New("repair action status is empty")
	}
	if a.Status == "BOUND" && a.MissingStage != "" {
		return errors.New("bound repair action has a missing stage")
	}
	if a.Status == "UNKNOWN" && a.MissingStage == "" {
		return errors.New("unknown repair action has no missing stage")
	}
	if a.MissingStageIndex < 0 {
		return errors.New("repair action missing stage index is negative")
	}
	if !a.NonExecuting || !a.NonAuthorizing {
		return errors.New("repair action must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction(a) != a.ActionDigest {
		return errors.New("repair action digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction(
	action RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPRepairAction,
) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{
		action.Status,
		action.MissingStage,
		strconv.Itoa(action.MissingStageIndex),
		action.GateDecision,
		action.ProjectionSignal,
		action.SuggestedAction,
		action.ActionKind,
		action.Title,
		action.Reason,
		action.ProjectionObservationDigest,
		action.EvidencePrefixDigest,
		action.SourceObservationDigest,
		action.ReverseObservationDigest,
		strconv.FormatBool(action.NonExecuting),
		strconv.FormatBool(action.NonAuthorizing),
	}, ":")))
	return hex.EncodeToString(sum[:])
}