package gooo

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

// JEVCalibrationConvergenceGateLSPRepairActionClosureStatus is the bounded
// state of a closure metric. It never implies that a repair action executed.
type JEVCalibrationConvergenceGateLSPRepairActionClosureStatus string

const (
	JEVCalibrationConvergenceGateLSPRepairActionClosureBound    JEVCalibrationConvergenceGateLSPRepairActionClosureStatus = "BOUND"
	JEVCalibrationConvergenceGateLSPRepairActionClosureUnknown  JEVCalibrationConvergenceGateLSPRepairActionClosureStatus = "UNKNOWN"
	JEVCalibrationConvergenceGateLSPRepairActionClosureDeferred JEVCalibrationConvergenceGateLSPRepairActionClosureStatus = "DEFERRED"
)

// JEVCalibrationConvergenceGateLSPRepairActionClosureInput is the evidence
// tuple that must remain stable from declaration through reverse observation.
type JEVCalibrationConvergenceGateLSPRepairActionClosureInput struct {
	SourceVersion             string
	ContractVersion           string
	Declaration               string
	IR                        string
	Generated                 string
	ActionMaterial            string
	ReverseObservation       string
	EvidencePrefix            string
	ActionStatus              string
	ReverseObservationStatus  string
	ActionDigest              string
	ReverseObservationDigest  string
	MissingStageIndex         int
}

// JEVCalibrationConvergenceGateLSPRepairActionClosureMetric is an evidence
// metric, not an execution result or an authorization decision.
type JEVCalibrationConvergenceGateLSPRepairActionClosureMetric struct {
	Status                    JEVCalibrationConvergenceGateLSPRepairActionClosureStatus
	SourceVersion             string
	ContractVersion           string
	DeclarationDigest         string
	IRDigest                  string
	GeneratedDigest           string
	ActionDigest              string
	ReverseObservationDigest  string
	EvidencePrefixDigest      string
	MissingStage               string
	MissingStageIndex         int
	Reason                    string
	MetricDigest              string
}

// ObserveJEVCalibrationConvergenceGateLSPRepairActionClosure calculates a
// bounded closure state. BOUND is emitted only when every stage and digest is
// present, exact, and explicitly marked BOUND by its producer.
func ObserveJEVCalibrationConvergenceGateLSPRepairActionClosure(in JEVCalibrationConvergenceGateLSPRepairActionClosureInput) JEVCalibrationConvergenceGateLSPRepairActionClosureMetric {
	metric := JEVCalibrationConvergenceGateLSPRepairActionClosureMetric{
		Status:                   JEVCalibrationConvergenceGateLSPRepairActionClosureUnknown,
		SourceVersion:            in.SourceVersion,
		ContractVersion:          in.ContractVersion,
		DeclarationDigest:        closureDigest(in.Declaration),
		IRDigest:                 closureDigest(in.IR),
		GeneratedDigest:          closureDigest(in.Generated),
		ActionDigest:             in.ActionDigest,
		ReverseObservationDigest: in.ReverseObservationDigest,
		EvidencePrefixDigest:     closureDigest(in.EvidencePrefix),
		MissingStageIndex:        in.MissingStageIndex,
	}

	if in.MissingStageIndex >= 0 {
		metric.MissingStage = "missing_stage_index"
		metric.Reason = "missing stage evidence"
		return finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureMetric(metric)
	}
	if missing := firstMissingJEVCalibrationConvergenceGateLSPRepairActionClosureStage(in); missing != "" {
		metric.MissingStage = missing
		metric.Reason = "missing closure evidence"
		return finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureMetric(metric)
	}
	if in.ActionStatus != string(JEVCalibrationConvergenceGateLSPRepairActionClosureBound) || in.ReverseObservationStatus != string(JEVCalibrationConvergenceGateLSPRepairActionClosureBound) {
		metric.Status = JEVCalibrationConvergenceGateLSPRepairActionClosureDeferred
		metric.Reason = "closure producer state is not BOUND"
		return finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureMetric(metric)
	}

	expectedActionDigest := closureDigest(in.SourceVersion, in.ContractVersion, in.Declaration, in.IR, in.Generated, in.ActionMaterial, in.EvidencePrefix, strconv.Itoa(in.MissingStageIndex))
	if in.ActionDigest != expectedActionDigest {
		metric.Reason = "action digest mismatch"
		return finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureMetric(metric)
	}
	expectedReverseDigest := closureDigest(in.ActionDigest, in.ReverseObservation)
	if in.ReverseObservationDigest != expectedReverseDigest {
		metric.Reason = "reverse observation digest mismatch"
		return finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureMetric(metric)
	}

	metric.Status = JEVCalibrationConvergenceGateLSPRepairActionClosureBound
	metric.Reason = "declaration, IR, generation, action, and reverse observation are digest-bound"
	return finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureMetric(metric)
}

func firstMissingJEVCalibrationConvergenceGateLSPRepairActionClosureStage(in JEVCalibrationConvergenceGateLSPRepairActionClosureInput) string {
	stages := []struct {
		name  string
		value string
	}{
		{"source_version", in.SourceVersion},
		{"contract_version", in.ContractVersion},
		{"declaration", in.Declaration},
		{"ir", in.IR},
		{"generated", in.Generated},
		{"action_material", in.ActionMaterial},
		{"reverse_observation", in.ReverseObservation},
		{"evidence_prefix", in.EvidencePrefix},
		{"action_status", in.ActionStatus},
		{"reverse_observation_status", in.ReverseObservationStatus},
		{"action_digest", in.ActionDigest},
		{"reverse_observation_digest", in.ReverseObservationDigest},
	}
	for _, stage := range stages {
		if stage.value == "" {
			return stage.name
		}
	}
	return ""
}

func finalizeJEVCalibrationConvergenceGateLSPRepairActionClosureMetric(metric JEVCalibrationConvergenceGateLSPRepairActionClosureMetric) JEVCalibrationConvergenceGateLSPRepairActionClosureMetric {
	metric.MetricDigest = closureDigest(
		string(metric.Status),
		metric.SourceVersion,
		metric.ContractVersion,
		metric.DeclarationDigest,
		metric.IRDigest,
		metric.GeneratedDigest,
		metric.ActionDigest,
		metric.ReverseObservationDigest,
		metric.EvidencePrefixDigest,
		metric.MissingStage,
		strconv.Itoa(metric.MissingStageIndex),
		metric.Reason,
	)
	return metric
}

func closureDigest(parts ...string) string {
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