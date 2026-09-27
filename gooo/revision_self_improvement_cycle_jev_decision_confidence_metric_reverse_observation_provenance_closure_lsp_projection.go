package gooo

import "fmt"

const (
	jevDecisionConfidenceClosureLSPUnknown = "jev-confidence-closure-lsp-unknown"
	jevDecisionConfidenceClosureLSPComplete = "jev-confidence-closure-lsp-complete"
)

type RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjection struct {
	Status                   string
	MissingStage             string
	MissingStageIndex        int
	MetricName               string
	MetricDigest             string
	EvidencePrefixDigest     string
	ReverseObservationDigest string
	ClosureSignal            string
	ObservationDigest        string
	NonExecuting             bool
	NonAuthorizing           bool
}

func ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjection(
	input RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetric,
) RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjection {
	output := RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjection{
		Status:             "UNKNOWN",
		MissingStage:       input.MissingStage,
		MissingStageIndex:  -1,
		MetricName:         jevDecisionConfidenceMetricName,
		ClosureSignal:      jevDecisionConfidenceClosureLSPUnknown,
		NonExecuting:       true,
		NonAuthorizing:     true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "reverse_observation_provenance_closure"
	}
	if err := input.Validate(); err != nil {
		output.ObservationDigest = revisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjectionDigest(output)
		return output
	}

	output.EvidencePrefixDigest = input.ObservationDigest
	if input.Status != "BOUND" {
		output.ObservationDigest = revisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjectionDigest(output)
		return output
	}

	output.Status = "BOUND"
	output.MissingStage = ""
	output.MissingStageIndex = 0
	output.MetricDigest = input.MetricDigest
	output.ReverseObservationDigest = input.ReverseObservationDigest
	output.ClosureSignal = jevDecisionConfidenceClosureLSPComplete
	output.ObservationDigest = revisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjectionDigest(output)
	return output
}

func (value RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjection) Validate() error {
	if value.Status != "UNKNOWN" && value.Status != "BOUND" {
		return fmt.Errorf("invalid status %q", value.Status)
	}
	if value.MetricName != jevDecisionConfidenceMetricName {
		return fmt.Errorf("invalid metric name %q", value.MetricName)
	}
	if value.MissingStageIndex < -1 {
		return fmt.Errorf("invalid missing stage index %d", value.MissingStageIndex)
	}
	if !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("LSP projection must remain non-executing and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" || value.MissingStageIndex != -1 {
			return fmt.Errorf("unknown LSP projection must preserve an unresolved stage")
		}
		if value.ClosureSignal != jevDecisionConfidenceClosureLSPUnknown {
			return fmt.Errorf("invalid unknown LSP closure signal %q", value.ClosureSignal)
		}
		if value.EvidencePrefixDigest != "" && !validJEVDecisionConfidenceClosureDigest(value.EvidencePrefixDigest) {
			return fmt.Errorf("invalid unknown evidence prefix digest")
		}
	} else {
		if value.MissingStage != "" || value.MissingStageIndex != 0 {
			return fmt.Errorf("bound LSP projection must close the stage")
		}
		if !validJEVDecisionConfidenceClosureDigest(value.MetricDigest) {
			return fmt.Errorf("invalid metric digest")
		}
		if !validJEVDecisionConfidenceClosureDigest(value.EvidencePrefixDigest) {
			return fmt.Errorf("invalid evidence prefix digest")
		}
		if !validJEVDecisionConfidenceClosureDigest(value.ReverseObservationDigest) {
			return fmt.Errorf("invalid reverse observation digest")
		}
		if value.ClosureSignal != jevDecisionConfidenceClosureLSPComplete {
			return fmt.Errorf("invalid bound LSP closure signal %q", value.ClosureSignal)
		}
	}
	expected := revisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjectionDigest(value)
	if value.ObservationDigest != expected {
		return fmt.Errorf("observation digest mismatch")
	}
	return nil
}

func revisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjectionDigest(
	value RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjection,
) string {
	return digestString(fmt.Sprintf(
		"jev-confidence-closure-lsp|%s|%s|%d|%s|%s|%s|%s|%s|%t|%t",
		value.Status,
		value.MissingStage,
		value.MissingStageIndex,
		value.MetricName,
		value.MetricDigest,
		value.EvidencePrefixDigest,
		value.ReverseObservationDigest,
		value.ClosureSignal,
		value.NonExecuting,
		value.NonAuthorizing,
	))
}