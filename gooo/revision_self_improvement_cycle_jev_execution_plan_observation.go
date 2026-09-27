package gooo

import "fmt"

const (
	jevExecutionPlanContractName = "revision_self_improvement_cycle_jev_execution_plan_provenance"
	jevExecutionPlanMetricName = "jev-execution-plan-provenance"
	jevExecutionPlanUnknown = "jev-execution-plan-unknown"
	jevExecutionPlanComplete = "jev-execution-plan-complete"
)

type RevisionSelfImprovementCycleJEVExecutionPlanEvidence struct {
	Status               string
	MissingStage         string
	ContractName         string
	TaskDigest           string
	WorkspaceDigest      string
	GatewayPolicyDigest  string
	ModelDigest          string
	SuspendResumeDigest  string
	EvidencePrefixDigest string
	ObservationDigest    string
	NonExecuting         bool
	NonAuthorizing       bool
}

type RevisionSelfImprovementCycleJEVExecutionPlanObservation struct {
	Status               string
	MissingStage         string
	MetricName           string
	ContractName         string
	TaskDigest           string
	WorkspaceDigest      string
	GatewayPolicyDigest  string
	ModelDigest           string
	SuspendResumeDigest   string
	EvidencePrefixDigest  string
	PlanSignal           string
	ObservationDigest     string
	NonExecuting          bool
	NonAuthorizing        bool
}

func ObserveRevisionSelfImprovementCycleJEVExecutionPlanObservation(
	input RevisionSelfImprovementCycleJEVExecutionPlanEvidence,
) RevisionSelfImprovementCycleJEVExecutionPlanObservation {
	output := RevisionSelfImprovementCycleJEVExecutionPlanObservation{
		Status:         "UNKNOWN",
		MissingStage:   input.MissingStage,
		MetricName:     jevExecutionPlanMetricName,
		PlanSignal:     jevExecutionPlanUnknown,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "execution_plan_evidence"
	}
	if err := input.Validate(); err != nil {
		output.ObservationDigest = revisionSelfImprovementCycleJEVExecutionPlanObservationDigest(output)
		return output
	}
	if input.Status != "BOUND" {
		output.ObservationDigest = revisionSelfImprovementCycleJEVExecutionPlanObservationDigest(output)
		return output
	}

	output.Status = "BOUND"
	output.MissingStage = ""
	output.ContractName = input.ContractName
	output.TaskDigest = input.TaskDigest
	output.WorkspaceDigest = input.WorkspaceDigest
	output.GatewayPolicyDigest = input.GatewayPolicyDigest
	output.ModelDigest = input.ModelDigest
	output.SuspendResumeDigest = input.SuspendResumeDigest
	output.EvidencePrefixDigest = input.EvidencePrefixDigest
	output.PlanSignal = jevExecutionPlanComplete
	output.ObservationDigest = revisionSelfImprovementCycleJEVExecutionPlanObservationDigest(output)
	return output
}

func (value RevisionSelfImprovementCycleJEVExecutionPlanEvidence) Validate() error {
	if value.Status != "UNKNOWN" && value.Status != "BOUND" {
		return fmt.Errorf("invalid status %q", value.Status)
	}
	if !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("execution plan evidence must remain non-executing and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" {
			return fmt.Errorf("unknown execution plan evidence must preserve a missing stage")
		}
		return nil
	}
	if value.MissingStage != "" || value.ContractName != jevExecutionPlanContractName {
		return fmt.Errorf("bound execution plan evidence must identify its contract")
	}
	for name, digest := range map[string]string{
		"task":            value.TaskDigest,
		"workspace":       value.WorkspaceDigest,
		"gateway_policy":  value.GatewayPolicyDigest,
		"model":           value.ModelDigest,
		"suspend_resume":  value.SuspendResumeDigest,
		"evidence_prefix": value.EvidencePrefixDigest,
	} {
		if !validJEVDecisionConfidenceClosureDigest(digest) {
			return fmt.Errorf("invalid %s digest", name)
		}
	}
	if value.ObservationDigest != revisionSelfImprovementCycleJEVExecutionPlanEvidenceDigest(value) {
		return fmt.Errorf("execution plan evidence digest mismatch")
	}
	return nil
}

func (value RevisionSelfImprovementCycleJEVExecutionPlanObservation) Validate() error {
	if value.Status != "UNKNOWN" && value.Status != "BOUND" {
		return fmt.Errorf("invalid status %q", value.Status)
	}
	if value.MetricName != jevExecutionPlanMetricName {
		return fmt.Errorf("invalid metric name %q", value.MetricName)
	}
	if !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("execution plan observation must remain non-executing and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" || value.PlanSignal != jevExecutionPlanUnknown {
			return fmt.Errorf("unknown execution plan observation must preserve its unresolved stage")
		}
	} else {
		if value.MissingStage != "" || value.ContractName != jevExecutionPlanContractName {
			return fmt.Errorf("bound execution plan observation must identify its contract")
		}
		for name, digest := range map[string]string{
			"task":            value.TaskDigest,
			"workspace":       value.WorkspaceDigest,
			"gateway_policy":  value.GatewayPolicyDigest,
			"model":           value.ModelDigest,
			"suspend_resume":  value.SuspendResumeDigest,
			"evidence_prefix": value.EvidencePrefixDigest,
		} {
			if !validJEVDecisionConfidenceClosureDigest(digest) {
				return fmt.Errorf("invalid %s digest", name)
			}
		}
		if value.PlanSignal != jevExecutionPlanComplete {
			return fmt.Errorf("invalid bound execution plan signal %q", value.PlanSignal)
		}
	}
	if value.ObservationDigest != revisionSelfImprovementCycleJEVExecutionPlanObservationDigest(value) {
		return fmt.Errorf("observation digest mismatch")
	}
	return nil
}

func revisionSelfImprovementCycleJEVExecutionPlanEvidenceDigest(
	value RevisionSelfImprovementCycleJEVExecutionPlanEvidence,
) string {
	return digestString(fmt.Sprintf(
		"jev-execution-plan-evidence|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t",
		value.Status,
		value.MissingStage,
		value.ContractName,
		value.TaskDigest,
		value.WorkspaceDigest,
		value.GatewayPolicyDigest,
		value.ModelDigest,
		value.SuspendResumeDigest,
		value.EvidencePrefixDigest,
		value.NonExecuting,
		value.NonAuthorizing,
	))
}

func revisionSelfImprovementCycleJEVExecutionPlanObservationDigest(
	value RevisionSelfImprovementCycleJEVExecutionPlanObservation,
) string {
	return digestString(fmt.Sprintf(
		"jev-execution-plan-observation|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t",
		value.Status,
		value.MissingStage,
		value.MetricName,
		value.ContractName,
		value.TaskDigest,
		value.WorkspaceDigest,
		value.GatewayPolicyDigest,
		value.ModelDigest,
		value.SuspendResumeDigest,
		value.EvidencePrefixDigest,
		value.PlanSignal,
		value.NonExecuting,
		value.NonAuthorizing,
	))
}
