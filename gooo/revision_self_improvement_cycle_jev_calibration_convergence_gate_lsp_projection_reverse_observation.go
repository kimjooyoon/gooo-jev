package gooo

import "fmt"

const (
	jevCalibrationConvergenceGateLSPReverseMetricName = "jev-calibration-convergence-gate-lsp-reverse-observation"
	jevCalibrationConvergenceGateLSPReverseUnknown    = "jev-calibration-convergence-gate-lsp-reverse-unknown"
	jevCalibrationConvergenceGateLSPReverseComplete   = "jev-calibration-convergence-gate-lsp-reverse-complete"
)

type RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionReverseObservation struct {
	Status                           string
	MissingStage                     string
	MissingStageIndex                int
	MetricName                       string
	GateDecision                     string
	ProjectionSignal                 string
	ProjectionObservationDigest      string
	EvidencePrefixDigest             string
	SourceObservationDigest          string
	ReverseObservationCoverageDigest string
	ReverseSignal                    string
	ObservationDigest                string
	NonExecuting                     bool
	NonAuthorizing                   bool
}

func ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionReverseObservation(
	input RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjection,
) RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionReverseObservation {
	output := RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionReverseObservation{
		Status:            "UNKNOWN",
		MissingStage:      input.MissingStage,
		MissingStageIndex: -1,
		MetricName:        jevCalibrationConvergenceGateLSPReverseMetricName,
		GateDecision:      jevCalibrationConvergenceGateUnknown,
		ReverseSignal:     jevCalibrationConvergenceGateLSPReverseUnknown,
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "jev_calibration_convergence_gate_lsp_projection"
	}
	setDigest := func() {
		output.ObservationDigest = revisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionReverseObservationDigest(output)
	}
	setDigest()

	if err := input.Validate(); err != nil {
		return output
	}

	output.ProjectionObservationDigest = input.ObservationDigest
	output.EvidencePrefixDigest = input.EvidencePrefixDigest
	output.SourceObservationDigest = input.SourceObservationDigest
	output.ReverseObservationCoverageDigest = input.ReverseObservationCoverageDigest
	output.GateDecision = input.GateDecision
	output.ProjectionSignal = input.ProjectionSignal
	if input.Status != "BOUND" {
		setDigest()
		return output
	}

	output.Status = "BOUND"
	output.MissingStage = ""
	output.MissingStageIndex = 0
	output.ReverseSignal = jevCalibrationConvergenceGateLSPReverseComplete
	setDigest()
	if err := output.Validate(); err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "jev_calibration_convergence_gate_lsp_projection_reverse_observation"
		output.MissingStageIndex = -1
		output.GateDecision = jevCalibrationConvergenceGateUnknown
		output.ReverseSignal = jevCalibrationConvergenceGateLSPReverseUnknown
		setDigest()
	}
	return output
}

func (value RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionReverseObservation) Validate() error {
	if value.Status != "UNKNOWN" && value.Status != "BOUND" {
		return fmt.Errorf("invalid status %q", value.Status)
	}
	if value.MetricName != jevCalibrationConvergenceGateLSPReverseMetricName {
		return fmt.Errorf("invalid metric name %q", value.MetricName)
	}
	if value.MissingStageIndex < -1 {
		return fmt.Errorf("invalid missing stage index %d", value.MissingStageIndex)
	}
	if !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("LSP convergence reverse observation must remain non-executing and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" ||
			value.MissingStageIndex != -1 ||
			value.GateDecision != jevCalibrationConvergenceGateUnknown ||
			value.ReverseSignal != jevCalibrationConvergenceGateLSPReverseUnknown {
			return fmt.Errorf("unknown LSP convergence reverse observation must preserve an unresolved stage")
		}
		if value.ProjectionObservationDigest != "" && !validDigest(value.ProjectionObservationDigest) {
			return fmt.Errorf("invalid unknown projection observation digest")
		}
		if value.EvidencePrefixDigest != "" && !validDigest(value.EvidencePrefixDigest) {
			return fmt.Errorf("invalid unknown evidence prefix digest")
		}
	} else {
		if value.MissingStage != "" || value.MissingStageIndex != 0 {
			return fmt.Errorf("bound LSP convergence reverse observation must close the stage")
		}
		if !validDigest(value.ProjectionObservationDigest) ||
			!validDigest(value.EvidencePrefixDigest) ||
			!validDigest(value.SourceObservationDigest) ||
			!validDigest(value.ReverseObservationCoverageDigest) ||
			value.ReverseSignal != jevCalibrationConvergenceGateLSPReverseComplete {
			return fmt.Errorf("invalid bound LSP convergence reverse provenance")
		}
		switch value.GateDecision {
		case jevCalibrationConvergenceGateConverged:
			if value.ProjectionSignal != jevCalibrationConvergenceGateLSPConverged {
				return fmt.Errorf("converged reverse observation has invalid projection signal")
			}
		case jevCalibrationConvergenceGateRegressed:
			if value.ProjectionSignal != jevCalibrationConvergenceGateLSPRegressed {
				return fmt.Errorf("regressed reverse observation has invalid projection signal")
			}
		case jevCalibrationConvergenceGateDefer:
			if value.ProjectionSignal != jevCalibrationConvergenceGateLSPDefer {
				return fmt.Errorf("deferred reverse observation has invalid projection signal")
			}
		default:
			return fmt.Errorf("invalid bound LSP convergence reverse decision")
		}
	}
	if value.ObservationDigest != revisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionReverseObservationDigest(value) {
		return fmt.Errorf("LSP convergence reverse observation digest mismatch")
	}
	return nil
}

func revisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionReverseObservationDigest(
	value RevisionSelfImprovementCycleJEVCalibrationConvergenceGateLSPProjectionReverseObservation,
) string {
	return digestString(fmt.Sprintf(
		"jev-calibration-convergence-gate-lsp-reverse|%s|%s|%d|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t",
		value.Status,
		value.MissingStage,
		value.MissingStageIndex,
		value.MetricName,
		value.GateDecision,
		value.ProjectionSignal,
		value.ProjectionObservationDigest,
		value.EvidencePrefixDigest,
		value.SourceObservationDigest,
		value.ReverseObservationCoverageDigest,
		value.ReverseSignal,
		value.NonExecuting,
		value.NonAuthorizing,
	))
}