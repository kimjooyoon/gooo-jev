package gooo

import "fmt"

const (
	jevGuardianBoundaryCoverageMetricName = "jev-guardian-boundary-evidence-coverage"
	jevGuardianBoundaryCoverageUnknown = "jev-guardian-boundary-evidence-coverage-unknown"
	jevGuardianBoundaryCoverageComplete = "jev-guardian-boundary-evidence-coverage-complete"
)

type RevisionSelfImprovementCycleJEVGuardianBoundaryEvidenceCoverageMetric struct {
	Status                   string
	MissingStage             string
	MetricName               string
	RequiredDigestCount      int
	LinkedDigestCount        int
	CoverageMilli             int
	CoverageBand              string
	MetricSignal              string
	MetricDigest              string
	GuardianObservationDigest string
	ObservationDigest         string
	NonExecuting              bool
	NonAuthorizing            bool
}

func ObserveRevisionSelfImprovementCycleJEVGuardianBoundaryEvidenceCoverageMetric(
	input RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureLSPProjectionGuardianBoundaryObservation,
) RevisionSelfImprovementCycleJEVGuardianBoundaryEvidenceCoverageMetric {
	output := RevisionSelfImprovementCycleJEVGuardianBoundaryEvidenceCoverageMetric{
		Status:              "UNKNOWN",
		MissingStage:        input.MissingStage,
		MetricName:          jevGuardianBoundaryCoverageMetricName,
		RequiredDigestCount: 3,
		CoverageBand:        "unknown",
		MetricSignal:        jevGuardianBoundaryCoverageUnknown,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "guardian_boundary_observation"
	}
	if err := input.Validate(); err != nil {
		output.ObservationDigest = revisionSelfImprovementCycleJEVGuardianBoundaryEvidenceCoverageMetricDigest(output)
		return output
	}
	if input.Status != "BOUND" {
		output.ObservationDigest = revisionSelfImprovementCycleJEVGuardianBoundaryEvidenceCoverageMetricDigest(output)
		return output
	}

	output.Status = "BOUND"
	output.MissingStage = ""
	output.LinkedDigestCount = 3
	output.CoverageMilli = 1000
	output.CoverageBand = "high"
	output.MetricSignal = jevGuardianBoundaryCoverageComplete
	output.MetricDigest = digestString(fmt.Sprintf(
		"jev-guardian-boundary-evidence-coverage|%s|%s|%s",
		input.WorkloadIdentityEvidenceDigest,
		input.CapabilityBoundaryEvidenceDigest,
		input.EvidencePrefixDigest,
	))
	output.GuardianObservationDigest = input.ObservationDigest
	output.ObservationDigest = revisionSelfImprovementCycleJEVGuardianBoundaryEvidenceCoverageMetricDigest(output)
	return output
}

func (value RevisionSelfImprovementCycleJEVGuardianBoundaryEvidenceCoverageMetric) Validate() error {
	if value.Status != "UNKNOWN" && value.Status != "BOUND" {
		return fmt.Errorf("invalid status %q", value.Status)
	}
	if value.MetricName != jevGuardianBoundaryCoverageMetricName {
		return fmt.Errorf("invalid metric name %q", value.MetricName)
	}
	if value.RequiredDigestCount != 3 || value.LinkedDigestCount < 0 || value.LinkedDigestCount > value.RequiredDigestCount {
		return fmt.Errorf("invalid digest coverage counts")
	}
	if !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("coverage metric must remain non-executing and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" {
			return fmt.Errorf("unknown coverage metric must preserve a missing stage")
		}
		if value.CoverageMilli != 0 || value.CoverageBand != "unknown" {
			return fmt.Errorf("unknown coverage metric must not claim coverage")
		}
		if value.MetricSignal != jevGuardianBoundaryCoverageUnknown {
			return fmt.Errorf("invalid unknown coverage signal %q", value.MetricSignal)
		}
	} else {
		if value.MissingStage != "" {
			return fmt.Errorf("bound coverage metric cannot have a missing stage")
		}
		if value.LinkedDigestCount != value.RequiredDigestCount || value.CoverageMilli != 1000 || value.CoverageBand != "high" {
			return fmt.Errorf("bound coverage metric must report complete coverage")
		}
		if !validJEVDecisionConfidenceClosureDigest(value.MetricDigest) ||
			!validJEVDecisionConfidenceClosureDigest(value.GuardianObservationDigest) {
			return fmt.Errorf("invalid coverage provenance digest")
		}
		if value.MetricSignal != jevGuardianBoundaryCoverageComplete {
			return fmt.Errorf("invalid bound coverage signal %q", value.MetricSignal)
		}
	}
	expected := revisionSelfImprovementCycleJEVGuardianBoundaryEvidenceCoverageMetricDigest(value)
	if value.ObservationDigest != expected {
		return fmt.Errorf("observation digest mismatch")
	}
	return nil
}

func revisionSelfImprovementCycleJEVGuardianBoundaryEvidenceCoverageMetricDigest(
	value RevisionSelfImprovementCycleJEVGuardianBoundaryEvidenceCoverageMetric,
) string {
	return digestString(fmt.Sprintf(
		"jev-guardian-boundary-evidence-coverage-metric|%s|%s|%s|%d|%d|%d|%s|%s|%s|%s|%t|%t",
		value.Status,
		value.MissingStage,
		value.MetricName,
		value.RequiredDigestCount,
		value.LinkedDigestCount,
		value.CoverageMilli,
		value.CoverageBand,
		value.MetricSignal,
		value.MetricDigest,
		value.GuardianObservationDigest,
		value.NonExecuting,
		value.NonAuthorizing,
	))
}