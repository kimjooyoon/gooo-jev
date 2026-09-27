package gooo

import "fmt"

const (
	jevCalibrationConvergenceGateLSPMetricName = "jev-calibration-convergence-gate-lsp"
	jevCalibrationConvergenceGateLSPUnknown    = "jev-calibration-convergence-gate-lsp-unknown"
	jevCalibrationConvergenceGateLSPConverged  = "jev-calibration-convergence-gate-lsp-converged"
	jevCalibrationConvergenceGateLSPRegressed  = "jev-calibration-convergence-gate-lsp-regressed"
	jevCalibrationConvergenceGateLSPDefer      = "jev-calibration-convergence-gate-lsp-defer"
)

type RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjection struct {
	Status                           string
	MissingStage                     string
	MissingStageIndex                int
	MetricName                       string
	GateDecision                     string
	OutcomeDeltaSignal               string
	ChoiceSetSensitivitySignal       string
	ChoiceSetDecisionSignal          string
	GateObservationDigest            string
	EvidencePrefixDigest             string
	SourceObservationDigest          string
	ReverseObservationCoverageDigest string
	ProjectionSignal                 string
	ObservationDigest                string
	NonExecuting                     bool
	NonAuthorizing                   bool
}

func ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjection(
	input RevisionSelfImprovementCycleJEVCalibrationConvergenceGateObservation,
) RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjection {
	output := RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjection{
		Status:            "UNKNOWN",
		MissingStage:      input.MissingStage,
		MissingStageIndex: -1,
		MetricName:        jevCalibrationConvergenceGateLSPMetricName,
		GateDecision:      jevCalibrationConvergenceGateUnknown,
		ProjectionSignal:  jevCalibrationConvergenceGateLSPUnknown,
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "jev_calibration_convergence_gate"
	}
	setDigest := func() {
		output.ObservationDigest = revisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionDigest(output)
	}
	setDigest()

	if err := input.Validate(); err != nil {
		return output
	}

	output.GateObservationDigest = input.ObservationDigest
	output.EvidencePrefixDigest = input.ObservationDigest
	output.OutcomeDeltaSignal = input.OutcomeDeltaSignal
	output.ChoiceSetSensitivitySignal = input.ChoiceSetSensitivitySignal
	output.ChoiceSetDecisionSignal = input.ChoiceSetDecisionSignal
	output.SourceObservationDigest = input.SourceObservationDigest
	output.ReverseObservationCoverageDigest = input.ReverseObservationCoverageDigest
	output.GateDecision = input.GateDecision
	if input.Status != "BOUND" {
		setDigest()
		return output
	}

	output.Status = "BOUND"
	output.MissingStage = ""
	output.MissingStageIndex = 0
	switch input.GateDecision {
	case jevCalibrationConvergenceGateConverged:
		output.ProjectionSignal = jevCalibrationConvergenceGateLSPConverged
	case jevCalibrationConvergenceGateRegressed:
		output.ProjectionSignal = jevCalibrationConvergenceGateLSPRegressed
	case jevCalibrationConvergenceGateDefer:
		output.ProjectionSignal = jevCalibrationConvergenceGateLSPDefer
	default:
		output.Status = "UNKNOWN"
		output.MissingStage = "jev_calibration_convergence_gate_lsp_decision"
		output.MissingStageIndex = -1
		output.GateDecision = jevCalibrationConvergenceGateUnknown
		output.ProjectionSignal = jevCalibrationConvergenceGateLSPUnknown
	}
	setDigest()
	if err := output.Validate(); err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "jev_calibration_convergence_gate_lsp_projection"
		output.MissingStageIndex = -1
		output.GateDecision = jevCalibrationConvergenceGateUnknown
		output.ProjectionSignal = jevCalibrationConvergenceGateLSPUnknown
		setDigest()
	}
	return output
}

func (value RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjection) Validate() error {
	if value.Status != "UNKNOWN" && value.Status != "BOUND" {
		return fmt.Errorf("invalid status %q", value.Status)
	}
	if value.MetricName != jevCalibrationConvergenceGateLSPMetricName {
		return fmt.Errorf("invalid metric name %q", value.MetricName)
	}
	if value.MissingStageIndex < -1 {
		return fmt.Errorf("invalid missing stage index %d", value.MissingStageIndex)
	}
	if !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("LSP convergence projection must remain non-executing and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" ||
			value.MissingStageIndex != -1 ||
			value.GateDecision != jevCalibrationConvergenceGateUnknown ||
			value.ProjectionSignal != jevCalibrationConvergenceGateLSPUnknown {
			return fmt.Errorf("unknown LSP convergence projection must preserve an unresolved stage")
		}
		if value.GateObservationDigest != "" && !validDigest(value.GateObservationDigest) {
			return fmt.Errorf("invalid unknown gate observation digest")
		}
		if value.EvidencePrefixDigest != "" && !validDigest(value.EvidencePrefixDigest) {
			return fmt.Errorf("invalid unknown evidence prefix digest")
		}
	} else {
		if value.MissingStage != "" || value.MissingStageIndex != 0 {
			return fmt.Errorf("bound LSP convergence projection must close the stage")
		}
		if !validDigest(value.GateObservationDigest) ||
			!validDigest(value.EvidencePrefixDigest) ||
			value.GateObservationDigest != value.EvidencePrefixDigest ||
			!validDigest(value.SourceObservationDigest) ||
			!validDigest(value.ReverseObservationCoverageDigest) {
			return fmt.Errorf("invalid bound LSP convergence provenance")
		}
		switch value.GateDecision {
		case jevCalibrationConvergenceGateConverged:
			if value.ProjectionSignal != jevCalibrationConvergenceGateLSPConverged {
				return fmt.Errorf("converged gate has invalid LSP signal")
			}
		case jevCalibrationConvergenceGateRegressed:
			if value.ProjectionSignal != jevCalibrationConvergenceGateLSPRegressed {
				return fmt.Errorf("regressed gate has invalid LSP signal")
			}
		case jevCalibrationConvergenceGateDefer:
			if value.ProjectionSignal != jevCalibrationConvergenceGateLSPDefer {
				return fmt.Errorf("deferred gate has invalid LSP signal")
			}
		default:
			return fmt.Errorf("invalid bound LSP convergence decision")
		}
	}
	if value.ObservationDigest != revisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionDigest(value) {
		return fmt.Errorf("LSP convergence projection digest mismatch")
	}
	return nil
}

func revisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionDigest(
	value RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjection,
) string {
	return digestString(fmt.Sprintf(
		"jev-calibration-convergence-gate-lsp|%s|%s|%d|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t",
		value.Status,
		value.MissingStage,
		value.MissingStageIndex,
		value.MetricName,
		value.GateDecision,
		value.OutcomeDeltaSignal,
		value.ChoiceSetSensitivitySignal,
		value.ChoiceSetDecisionSignal,
		value.GateObservationDigest,
		value.EvidencePrefixDigest,
		value.SourceObservationDigest,
		value.ReverseObservationCoverageDigest,
		value.ProjectionSignal,
		value.NonExecuting,
		value.NonAuthorizing,
	))
}