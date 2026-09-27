package gooo

import "fmt"

const (
	jevExecutionPlanCoverageRepairPlanName     = "jev-execution-plan-coverage-repair-plan"
	jevExecutionPlanCoverageRepairPlanUnknown  = "jev-execution-plan-coverage-repair-plan-unknown"
	jevExecutionPlanCoverageRepairPlanComplete = "jev-execution-plan-coverage-repair-plan-complete"
)

type RevisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlan struct {
	Status               string
	MissingStage         string
	PlanName             string
	TargetStage          string
	RepairAction         string
	EvidencePrefixDigest string
	ProjectionDigest     string
	PlanSignal           string
	ObservationDigest    string
	NonExecuting         bool
	NonAuthorizing       bool
}

func ObserveRevisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlan(
	input RevisionSelfImprovementCycleJEVExecutionPlanReverseObservationCoverageMetricLSPProjection,
) RevisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlan {
	output := RevisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlan{
		Status:               "UNKNOWN",
		MissingStage:         input.MissingStage,
		PlanName:             jevExecutionPlanCoverageRepairPlanName,
		RepairAction:         "resolve_missing_provenance",
		PlanSignal:           jevExecutionPlanCoverageRepairPlanUnknown,
		NonExecuting:         true,
		NonAuthorizing:       true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "execution_plan_coverage_lsp_projection"
	}
	output.TargetStage = output.MissingStage
	if err := input.Validate(); err != nil {
		output.ObservationDigest = revisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlanDigest(output)
		return output
	}

	output.EvidencePrefixDigest = input.EvidencePrefixDigest
	output.ProjectionDigest = input.ObservationDigest
	if input.Status != "BOUND" {
		output.ObservationDigest = revisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlanDigest(output)
		return output
	}

	output.Status = "BOUND"
	output.MissingStage = ""
	output.TargetStage = ""
	output.RepairAction = "none"
	output.PlanSignal = jevExecutionPlanCoverageRepairPlanComplete
	output.ObservationDigest = revisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlanDigest(output)
	return output
}

func (value RevisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlan) Validate() error {
	if value.Status != "UNKNOWN" && value.Status != "BOUND" {
		return fmt.Errorf("invalid status %q", value.Status)
	}
	if value.PlanName != jevExecutionPlanCoverageRepairPlanName {
		return fmt.Errorf("invalid repair plan name %q", value.PlanName)
	}
	if !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("repair plan must remain non-executing and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" || value.TargetStage != value.MissingStage {
			return fmt.Errorf("unknown repair plan must preserve its target stage")
		}
		if value.RepairAction != "resolve_missing_provenance" {
			return fmt.Errorf("invalid unknown repair action %q", value.RepairAction)
		}
		if value.PlanSignal != jevExecutionPlanCoverageRepairPlanUnknown {
			return fmt.Errorf("invalid unknown repair plan signal %q", value.PlanSignal)
		}
		if value.EvidencePrefixDigest != "" && !validJEVDecisionConfidenceClosureDigest(value.EvidencePrefixDigest) {
			return fmt.Errorf("invalid unknown evidence prefix digest")
		}
		if value.ProjectionDigest != "" && !validJEVDecisionConfidenceClosureDigest(value.ProjectionDigest) {
			return fmt.Errorf("invalid unknown projection digest")
		}
	} else {
		if value.MissingStage != "" || value.TargetStage != "" {
			return fmt.Errorf("bound repair plan cannot retain a target stage")
		}
		if value.RepairAction != "none" {
			return fmt.Errorf("bound repair plan must not request an edit")
		}
		if !validJEVDecisionConfidenceClosureDigest(value.EvidencePrefixDigest) ||
			!validJEVDecisionConfidenceClosureDigest(value.ProjectionDigest) {
			return fmt.Errorf("invalid bound repair plan provenance digest")
		}
		if value.PlanSignal != jevExecutionPlanCoverageRepairPlanComplete {
			return fmt.Errorf("invalid bound repair plan signal %q", value.PlanSignal)
		}
	}
	expected := revisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlanDigest(value)
	if value.ObservationDigest != expected {
		return fmt.Errorf("observation digest mismatch")
	}
	return nil
}

func revisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlanDigest(
	value RevisionSelfImprovementCycleJEVExecutionPlanCoverageLSPRepairPlan,
) string {
	return digestString(fmt.Sprintf(
		"jev-execution-plan-coverage-lsp-repair-plan|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t",
		value.Status,
		value.MissingStage,
		value.PlanName,
		value.TargetStage,
		value.RepairAction,
		value.EvidencePrefixDigest,
		value.ProjectionDigest,
		value.PlanSignal,
		value.NonExecuting,
		value.NonAuthorizing,
	))
}
