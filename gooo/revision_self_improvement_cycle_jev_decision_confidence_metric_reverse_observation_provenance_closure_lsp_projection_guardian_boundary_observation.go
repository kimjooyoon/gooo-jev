package gooo

import "fmt"

const (
	jevGuardianBoundaryUnknownSignal = "jev-guardian-boundary-unknown"
	jevGuardianBoundaryObservedSignal = "jev-guardian-boundary-observed"
	jevGuardianBoundaryDecisionDefer = "DEFER"
)

type RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjectionGuardianBoundaryObservation struct {
	Status                          string
	MissingStage                    string
	DecisionState                   string
	WorkloadIdentityEvidenceDigest  string
	CapabilityBoundaryEvidenceDigest string
	EvidencePrefixDigest            string
	GuardianSignal                  string
	ObservationDigest               string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

func ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjectionGuardianBoundaryObservation(
	input RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjection,
) RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjectionGuardianBoundaryObservation {
	output := RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjectionGuardianBoundaryObservation{
		Status:         "UNKNOWN",
		MissingStage:   input.MissingStage,
		DecisionState:  jevGuardianBoundaryDecisionDefer,
		GuardianSignal: jevGuardianBoundaryUnknownSignal,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "closure_lsp_projection"
	}
	if err := input.Validate(); err != nil {
		output.ObservationDigest = revisionSelfImprovementCycleJEVGuardianBoundaryObservationDigest(output)
		return output
	}

	output.EvidencePrefixDigest = input.EvidencePrefixDigest
	if input.Status != "BOUND" {
		output.ObservationDigest = revisionSelfImprovementCycleJEVGuardianBoundaryObservationDigest(output)
		return output
	}

	output.Status = "BOUND"
	output.MissingStage = ""
	output.WorkloadIdentityEvidenceDigest = digestString("jev-workload-identity-evidence|" + input.EvidencePrefixDigest)
	output.CapabilityBoundaryEvidenceDigest = digestString("jev-capability-boundary-evidence|" + input.MetricDigest)
	output.GuardianSignal = jevGuardianBoundaryObservedSignal
	output.ObservationDigest = revisionSelfImprovementCycleJEVGuardianBoundaryObservationDigest(output)
	return output
}

func (value RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjectionGuardianBoundaryObservation) Validate() error {
	if value.Status != "UNKNOWN" && value.Status != "BOUND" {
		return fmt.Errorf("invalid status %q", value.Status)
	}
	if value.DecisionState != jevGuardianBoundaryDecisionDefer {
		return fmt.Errorf("guardian boundary observation must remain DEFER")
	}
	if !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("guardian boundary observation must remain non-executing and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" {
			return fmt.Errorf("unknown guardian observation must preserve a missing stage")
		}
		if value.GuardianSignal != jevGuardianBoundaryUnknownSignal {
			return fmt.Errorf("invalid unknown guardian signal %q", value.GuardianSignal)
		}
	} else {
		if value.MissingStage != "" {
			return fmt.Errorf("bound guardian observation cannot have a missing stage")
		}
		if !validJEVDecisionConfidenceClosureDigest(value.WorkloadIdentityEvidenceDigest) {
			return fmt.Errorf("invalid workload identity evidence digest")
		}
		if !validJEVDecisionConfidenceClosureDigest(value.CapabilityBoundaryEvidenceDigest) {
			return fmt.Errorf("invalid capability boundary evidence digest")
		}
		if !validJEVDecisionConfidenceClosureDigest(value.EvidencePrefixDigest) {
			return fmt.Errorf("invalid evidence prefix digest")
		}
		if value.GuardianSignal != jevGuardianBoundaryObservedSignal {
			return fmt.Errorf("invalid bound guardian signal %q", value.GuardianSignal)
		}
	}
	expected := revisionSelfImprovementCycleJEVGuardianBoundaryObservationDigest(value)
	if value.ObservationDigest != expected {
		return fmt.Errorf("observation digest mismatch")
	}
	return nil
}

func revisionSelfImprovementCycleJEVGuardianBoundaryObservationDigest(
	value RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjectionGuardianBoundaryObservation,
) string {
	return digestString(fmt.Sprintf(
		"jev-guardian-boundary-observation|%s|%s|%s|%s|%s|%s|%s|%t|%t",
		value.Status,
		value.MissingStage,
		value.DecisionState,
		value.WorkloadIdentityEvidenceDigest,
		value.CapabilityBoundaryEvidenceDigest,
		value.EvidencePrefixDigest,
		value.GuardianSignal,
		value.NonExecuting,
		value.NonAuthorizing,
	))
}