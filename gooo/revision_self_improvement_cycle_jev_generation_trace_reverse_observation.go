package gooo

import "fmt"

const (
	jevGenerationTraceReverseMetricName = "jev-generation-trace-reverse-observation"
	jevGenerationTraceReverseUnknown = "jev-generation-trace-reverse-unknown"
	jevGenerationTraceReverseComplete = "jev-generation-trace-reverse-complete"
)

type RevisionSelfImprovementCycleJEVGenerationTraceReverseObservation struct {
	Status                  string
	MissingStage            string
	MetricName              string
	ContractName            string
	ContractDigest           string
	IRDigest                 string
	GeneratedArtifactDigest  string
	MetricDigest             string
	EvidencePrefixDigest     string
	ReconstructedTraceDigest string
	ReverseSignal            string
	ObservationDigest        string
	NonExecuting             bool
	NonAuthorizing           bool
}

func ObserveRevisionSelfImprovementCycleJEVGenerationTraceReverseObservation(
	input RevisionSelfImprovementCycleJEVGenerationTraceObservation,
) RevisionSelfImprovementCycleJEVGenerationTraceReverseObservation {
	output := RevisionSelfImprovementCycleJEVGenerationTraceReverseObservation{
		Status:         "UNKNOWN",
		MissingStage:   input.MissingStage,
		MetricName:     jevGenerationTraceReverseMetricName,
		ReverseSignal:  jevGenerationTraceReverseUnknown,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "generation_trace_observation"
	}
	if err := input.Validate(); err != nil {
		output.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceReverseObservationDigest(output)
		return output
	}
	if input.Status != "BOUND" {
		output.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceReverseObservationDigest(output)
		return output
	}

	output.Status = "BOUND"
	output.MissingStage = ""
	output.ContractName = input.ContractName
	output.ContractDigest = input.ContractDigest
	output.IRDigest = input.IRDigest
	output.GeneratedArtifactDigest = input.GeneratedArtifactDigest
	output.MetricDigest = input.MetricDigest
	output.EvidencePrefixDigest = input.EvidencePrefixDigest
	output.ReconstructedTraceDigest = revisionSelfImprovementCycleJEVGenerationTraceReconstructedDigest(output)
	output.ReverseSignal = jevGenerationTraceReverseComplete
	output.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceReverseObservationDigest(output)
	return output
}

func (value RevisionSelfImprovementCycleJEVGenerationTraceReverseObservation) Validate() error {
	if value.Status != "UNKNOWN" && value.Status != "BOUND" {
		return fmt.Errorf("invalid status %q", value.Status)
	}
	if value.MetricName != jevGenerationTraceReverseMetricName {
		return fmt.Errorf("invalid metric name %q", value.MetricName)
	}
	if !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("generation reverse observation must remain non-executing and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" || value.ReverseSignal != jevGenerationTraceReverseUnknown {
			return fmt.Errorf("unknown generation reverse observation must preserve its unresolved stage")
		}
	} else {
		if value.MissingStage != "" || value.ContractName != jevGenerationTraceContractName {
			return fmt.Errorf("bound generation reverse observation must identify its contract")
		}
		for name, digest := range map[string]string{
			"contract":           value.ContractDigest,
			"ir":                 value.IRDigest,
			"generated_artifact": value.GeneratedArtifactDigest,
			"metric":             value.MetricDigest,
			"evidence_prefix":    value.EvidencePrefixDigest,
			"reconstructed_trace": value.ReconstructedTraceDigest,
		} {
			if !validJEVDecisionConfidenceClosureDigest(digest) {
				return fmt.Errorf("invalid %s digest", name)
			}
		}
		if value.ReconstructedTraceDigest != revisionSelfImprovementCycleJEVGenerationTraceReconstructedDigest(value) {
			return fmt.Errorf("reconstructed trace digest mismatch")
		}
		if value.ReverseSignal != jevGenerationTraceReverseComplete {
			return fmt.Errorf("invalid bound generation reverse signal %q", value.ReverseSignal)
		}
	}
	expected := revisionSelfImprovementCycleJEVGenerationTraceReverseObservationDigest(value)
	if value.ObservationDigest != expected {
		return fmt.Errorf("observation digest mismatch")
	}
	return nil
}

func revisionSelfImprovementCycleJEVGenerationTraceReconstructedDigest(
	value RevisionSelfImprovementCycleJEVGenerationTraceReverseObservation,
) string {
	return digestString(fmt.Sprintf(
		"jev-generation-trace-reconstructed|%v|%v|%v|%v|%v|%v|%v",
		value.ContractDigest,
		value.IRDigest,
		value.GeneratedArtifactDigest,
		value.MetricDigest,
		value.EvidencePrefixDigest,
		value.ContractName,
		value.MetricName,
	))
}

func revisionSelfImprovementCycleJEVGenerationTraceReverseObservationDigest(
	value RevisionSelfImprovementCycleJEVGenerationTraceReverseObservation,
) string {
	return digestString(fmt.Sprintf(
		"jev-generation-trace-reverse-observation|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v|%v",
		value.Status,
		value.MissingStage,
		value.MetricName,
		value.ContractName,
		value.ContractDigest,
		value.IRDigest,
		value.GeneratedArtifactDigest,
		value.MetricDigest,
		value.EvidencePrefixDigest,
		value.ReconstructedTraceDigest,
		value.ReverseSignal,
		value.NonExecuting,
		value.NonAuthorizing,
	))
}