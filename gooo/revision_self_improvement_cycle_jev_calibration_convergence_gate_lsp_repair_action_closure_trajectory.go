package gooo

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

// JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryStatus records
// an evidence-backed state transition without inventing a quality score.
type JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryStatus string

const (
	JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryImproved  JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryStatus = "IMPROVED"
	JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryRegressed  JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryStatus = "REGRESSED"
	JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryUnchanged  JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryStatus = "UNCHANGED"
	JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryUnknown    JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryStatus = "UNKNOWN"
	JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryDeferred   JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryStatus = "DEFERRED"
)

// JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryObservation
// compares only statuses and exact identity. It does not infer semantic quality
// from a changed digest.
type JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryObservation struct {
	SourceVersion          string
	ContractVersion        string
	PreviousStatus         JEVCalibrationConvergenceGateLSPRepairActionClosureStatus
	CurrentStatus          JEVCalibrationConvergenceGateLSPRepairActionClosureStatus
	Status                 JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryStatus
	PreviousMetricDigest   string
	CurrentMetricDigest    string
	Reason                 string
	DeltaDigest            string
}

// ObserveJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectory
// preserves improvement and regression transitions while refusing to call an
// incomparable digest change an improvement.
func ObserveJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectory(previous, current JEVCalibrationConvergenceGateLSPRepairActionClosureMetric) JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryObservation {
	observation := JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryObservation{
		SourceVersion:        current.SourceVersion,
		ContractVersion:      current.ContractVersion,
		PreviousStatus:       previous.Status,
		CurrentStatus:        current.Status,
		Status:               JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryUnknown,
		PreviousMetricDigest: previous.MetricDigest,
		CurrentMetricDigest:  current.MetricDigest,
	}

	if previous.SourceVersion != current.SourceVersion || previous.ContractVersion != current.ContractVersion {
		observation.Reason = "source or contract identity changed"
		return finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryObservation(observation)
	}
	if previous.MetricDigest == "" || current.MetricDigest == "" {
		observation.Reason = "metric digest is missing"
		return finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryObservation(observation)
	}
	if previous.Status == JEVCalibrationConvergenceGateLSPRepairActionClosureDeferred || current.Status == JEVCalibrationConvergenceGateLSPRepairActionClosureDeferred {
		observation.Status = JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryDeferred
		observation.Reason = "closure producer state is deferred"
		return finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryObservation(observation)
	}
	if current.Status == JEVCalibrationConvergenceGateLSPRepairActionClosureBound && previous.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureBound {
		observation.Status = JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryImproved
		observation.Reason = "closure became BOUND from a non-BOUND state"
		return finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryObservation(observation)
	}
	if previous.Status == JEVCalibrationConvergenceGateLSPRepairActionClosureBound && current.Status != JEVCalibrationConvergenceGateLSPRepairActionClosureBound {
		observation.Status = JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryRegressed
		observation.Reason = "closure left BOUND for a non-BOUND state"
		return finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryObservation(observation)
	}
	if previous.Status == JEVCalibrationConvergenceGateLSPRepairActionClosureBound && current.Status == JEVCalibrationConvergenceGateLSPRepairActionClosureBound && previous.MetricDigest == current.MetricDigest {
		observation.Status = JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryUnchanged
		observation.Reason = "BOUND closure digest is unchanged"
		return finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryObservation(observation)
	}

	observation.Reason = "bound metric changed without a comparable quality score"
	return finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryObservation(observation)
}

func finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryObservation(observation JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryObservation) JEVCalibrationConvergenceGateLSPRepairActionClosureTrajectoryObservation {
	observation.DeltaDigest = closureTrajectoryDigest(
		observation.SourceVersion,
		observation.ContractVersion,
		string(observation.PreviousStatus),
		string(observation.CurrentStatus),
		string(observation.Status),
		observation.PreviousMetricDigest,
		observation.CurrentMetricDigest,
		observation.Reason,
	)
	return observation
}

func closureTrajectoryDigest(parts ...string) string {
	var encoded strings.Builder
	for _, part := range parts {
		encoded.WriteString(strconv.Itoa(len(part)))
		encoded.WriteByte(':')
		encoded.WriteString(part)
		encoded.WriteByte('|')
	}
	sum := sha256.Sum256([]byte(encoded.String()))
	return hex.EncodeToString(sum[:])
}