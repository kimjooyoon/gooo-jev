package gooo

import "fmt"

const (
	jevGenerationTraceContractName = "revision_self_improvement_cycle_jev_guardian_boundary_evidence_coverage_metric"
	jevGenerationTraceMetricName = "jev-generation-trace-closure"
	jevGenerationTraceUnknownSignal = "jev-generation-trace-unknown"
	jevGenerationTraceCompleteSignal = "jev-generation-trace-complete"
)

type RevisionSelfImprovementCycleJEVGenerationTraceEvidence struct {
	Status                  string
	MissingStage            string
	ContractName            string
	ContractDigest           string
	IRDigest                 string
	GeneratedArtifactDigest  string
	ReverseObservationDigest string
	MetricDigest             string
	EvidencePrefixDigest     string
	ObservationDigest        string
	NonExecuting             bool
	NonAuthorizing           bool
}

type RevisionSelfImprovementCycleJEVGenerationTraceObservation struct {
	Status                  string
	MissingStage            string
	MetricName              string
	ContractName             string
	ContractDigest           string
	IRDigest                 string
	GeneratedArtifactDigest  string
	ReverseObservationDigest string
	MetricDigest             string
	EvidencePrefixDigest     string
	TraceSignal              string
	ObservationDigest        string
	NonExecuting             bool
	NonAuthorizing           bool
}

func ObserveRevisionSelfImprovementCycleJEVGenerationTraceObservation(
	input RevisionSelfImprovementCycleJEVGenerationTraceEvidence,
) RevisionSelfImprovementCycleJEVGenerationTraceObservation {
	output := RevisionSelfImprovementCycleJEVGenerationTraceObservation{
		Status:         "UNKNOWN",
		MissingStage:   input.MissingStage,
		MetricName:     jevGenerationTraceMetricName,
		TraceSignal:    jevGenerationTraceUnknownSignal,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "generation_trace_evidence"
	}
	if err := input.Validate(); err != nil {
		output.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceObservationDigest(output)
		return output
	}
	if input.Status != "BOUND" {
		output.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceObservationDigest(output)
		return output
	}

	output.Status = "BOUND"
	output.MissingStage = ""
	output.ContractName = input.ContractName
	output.ContractDigest = input.ContractDigest
	output.IRDigest = input.IRDigest
	output.GeneratedArtifactDigest = input.GeneratedArtifactDigest
	output.ReverseObservationDigest = input.ReverseObservationDigest
	output.MetricDigest = input.MetricDigest
	output.EvidencePrefixDigest = input.EvidencePrefixDigest
	output.TraceSignal = jevGenerationTraceCompleteSignal
	output.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceObservationDigest(output)
	return output
}

func (value RevisionSelfImprovementCycleJEVGenerationTraceEvidence) Validate() error {
	if value.Status != "UNKNOWN" && value.Status != "BOUND" {
		return fmt.Errorf("invalid status %q", value.Status)
	}
	if !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("generation trace evidence must remain non-executing and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" {
			return fmt.Errorf("unknown generation trace evidence must preserve a missing stage")
		}
		return nil
	}
	if value.MissingStage != "" || value.ContractName != jevGenerationTraceContractName {
		return fmt.Errorf("bound generation trace evidence must identify its contract")
	}
	for name, digest := range map[string]string{
		"contract":            value.ContractDigest,
		"ir":                  value.IRDigest,
		"generated_artifact":  value.GeneratedArtifactDigest,
		"reverse_observation": value.ReverseObservationDigest,
		"metric":              value.MetricDigest,
		"evidence_prefix":     value.EvidencePrefixDigest,
	} {
		if !validJEVDecisionConfidenceClosureDigest(digest) {
			return fmt.Errorf("invalid %s digest", name)
		}
	}
	expected := revisionSelfImprovementCycleJEVGenerationTraceEvidenceDigest(value)
	if value.ObservationDigest != expected {
		return fmt.Errorf("generation trace evidence digest mismatch")
	}
	return nil
}

func (value RevisionSelfImprovementCycleJEVGenerationTraceObservation) Validate() error {
	if value.Status != "UNKNOWN" && value.Status != "BOUND" {
		return fmt.Errorf("invalid status %q", value.Status)
	}
	if value.MetricName != jevGenerationTraceMetricName {
		return fmt.Errorf("invalid metric name %q", value.MetricName)
	}
	if !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("generation trace observation must remain non-executing and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" || value.TraceSignal != jevGenerationTraceUnknownSignal {
			return fmt.Errorf("unknown generation trace observation must preserve its unresolved stage")
		}
	} else {
		if value.MissingStage != "" || value.ContractName != jevGenerationTraceContractName {
			return fmt.Errorf("bound generation trace observation must identify its contract")
		}
		for name, digest := range map[string]string{
			"contract":            value.ContractDigest,
			"ir":                  value.IRDigest,
			"generated_artifact":  value.GeneratedArtifactDigest,
			"reverse_observation": value.ReverseObservationDigest,
			"metric":              value.MetricDigest,
			"evidence_prefix":     value.EvidencePrefixDigest,
		} {
			if !validJEVDecisionConfidenceClosureDigest(digest) {
				return fmt.Errorf("invalid %s digest", name)
			}
		}
		if value.TraceSignal != jevGenerationTraceCompleteSignal {
			return fmt.Errorf("invalid bound generation trace signal %q", value.TraceSignal)
		}
	}
	expected := revisionSelfImprovementCycleJEVGenerationTraceObservationDigest(value)
	if value.ObservationDigest != expected {
		return fmt.Errorf("observation digest mismatch")
	}
	return nil
}

func revisionSelfImprovementCycleJEVGenerationTraceEvidenceDigest(
	value RevisionSelfImprovementCycleJEVGenerationTraceEvidence,
) string {
	return digestString(fmt.Sprintf(
		"jev-generation-trace-evidence|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v",
		value.Status,
		value.MissingStage,
		value.ContractName,
		value.ContractDigest,
		value.IRDigest,
		value.GeneratedArtifactDigest,
		value.ReverseObservationDigest,
		value.MetricDigest,
		value.EvidencePrefixDigest,
		value.NonExecuting,
		value.NonAuthorizing,
	))
}

func revisionSelfImprovementCycleJEVGenerationTraceObservationDigest(
	value RevisionSelfImprovementCycleJEVGenerationTraceObservation,
) string {
	return digestString(fmt.Sprintf(
		"jev-generation-trace-observation|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v",
		value.Status,
		value.MissingStage,
		value.MetricName,
		value.ContractName,
		value.ContractDigest,
		value.IRDigest,
		value.GeneratedArtifactDigest,
		value.ReverseObservationDigest,
		value.MetricDigest,
		value.EvidencePrefixDigest,
		value.TraceSignal,
		value.NonExecuting,
		value.NonAuthorizing,
	))
}