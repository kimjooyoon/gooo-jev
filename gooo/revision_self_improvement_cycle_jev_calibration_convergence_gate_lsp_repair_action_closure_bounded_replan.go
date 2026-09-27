package gooo

// JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateStatus
// describes a generated proposal, never an execution result.
type JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateStatus string

const (
	JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateProposed  JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateStatus = "PROPOSED"
	JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateNoAction  JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateStatus = "NO_ACTION"
	JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateUnknown   JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateStatus = "UNKNOWN"
	JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateDeferred  JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateStatus = "DEFERRED"
)

// JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidate is a
// bounded next-step proposal derived from reverse-observed provenance.
type JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidate struct {
	Status             JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateStatus
	SourceVersion      string
	ContractVersion    string
	TrajectoryStatus   JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryStatus
	MetricStatus       JEVCalibrationConvergenceGateLSPRepairActionClosureStatus
	TargetStage        string
	MissingStageIndex  int
	BasisDeltaDigest   string
	MetricDigest       string
	Reason             string
	PlanDigest         string
	IsReadOnly         bool
	CanExecute         bool
	CanAuthorize       bool
}

// ProposeJEVCalibrationConvergenceGateLSPRepairActionClosureReplan creates a
// proposal only for a regressed or unresolved trajectory. It never emits a
// command, edit, execution token, or authorization decision.
func ProposeJEVCalibrationConvergenceGateLSPRepairActionClosureReplan(trajectory JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryObservation, metric JEVCalibrationConvergenceGateLSPRepairActionClosureMetric) JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidate {
	candidate := JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidate{
		Status:            JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateUnknown,
		SourceVersion:     metric.SourceVersion,
		ContractVersion:   metric.ContractVersion,
		TrajectoryStatus:  trajectory.Status,
		MetricStatus:      metric.Status,
		MissingStageIndex: metric.MissingStageIndex,
		BasisDeltaDigest:  trajectory.DeltaDigest,
		MetricDigest:      metric.MetricDigest,
		IsReadOnly:        true,
		CanExecute:        false,
		CanAuthorize:      false,
	}

	if trajectory.SourceVersion != metric.SourceVersion || trajectory.ContractVersion != metric.ContractVersion {
		candidate.Reason = "source or contract identity changed"
		return finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidate(candidate)
	}
	if trajectory.DeltaDigest == "" || metric.MetricDigest == "" {
		candidate.Reason = "replan basis evidence is missing"
		return finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidate(candidate)
	}

	switch trajectory.Status {
	case JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryImproved, JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryUnchanged:
		candidate.Status = JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateNoAction
		candidate.Reason = "trajectory does not require a replan"
	case JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryDeferred:
		candidate.Status = JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateDeferred
		candidate.Reason = "trajectory producer is deferred"
	case JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryRegressed, JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryUnknown:
		candidate.Status = JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidateProposed
		candidate.TargetStage = metric.MissingStage
		if candidate.TargetStage == "" {
			candidate.TargetStage = "closure_digest"
		}
		candidate.Reason = "generate a bounded proposal for the first unresolved closure stage"
	default:
		candidate.Reason = "trajectory status is not recognized"
	}
	return finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidate(candidate)
}

func finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidate(candidate JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidate) JEVCalibrationConvergenceGateLSPRepairActionClosureReplanCandidate {
	candidate.PlanDigest = closureTrajectoryDigest(
		string(candidate.Status),
		candidate.SourceVersion,
		candidate.ContractVersion,
		string(candidate.TrajectoryStatus),
		string(candidate.MetricStatus),
		candidate.TargetStage,
		candidate.BasisDeltaDigest,
		candidate.MetricDigest,
		candidate.Reason,
	)
	return candidate
}