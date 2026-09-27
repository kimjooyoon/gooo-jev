package gooo

import "fmt"

const (
	jevExecutionPlanReverseMetricName = "jev-execution-plan-reverse-observation"
	jevExecutionPlanReverseUnknown = "jev-execution-plan-reverse-unknown"
	jevExecutionPlanReverseComplete = "jev-execution-plan-reverse-complete"
)

type RevisionSelfImprovementCycleJEVExecutionPlanReverseObservation struct {
	Status               string
	MissingStage         string
	MetricName           string
	ContractName         string
	TaskDigest           string
	WorkspaceDigest      string
	GatewayPolicyDigest  string
	ModelDigest          string
	SuspendResumeDigest  string
	EvidencePrefixDigest string
	PlanSignal           string
	ObservationDigest    string
	NonExecuting         bool
	NonAuthorizing       bool
}

func ObserveRevisionSelfImprovementCycleJEVExecutionPlanReverseObservation(
	input RevisionSelfImprovementCycleJEVExecutionPlanObservation,
) RevisionSelfImprovementCycleJEVExecutionPlanReverseObservation {
	output := RevisionSelfImprovementCycleJEVExecutionPlanReverseObservation{
		Status:         "UNKNOWN",
		MissingStage:   input.MissingStage,
		MetricName:     jevExecutionPlanReverseMetricName,
		PlanSignal:     jevExecutionPlanReverseUnknown,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "execution_plan_observation"
	}
	if err := input.Validate(); err != nil {
		output.ObservationDigest = revisionSelfImprovementCycleJEVExecutionPlanReverseObservationDigest(output)
		return output
	}
	if input.Status != "BOUND" {
		output.ObservationDigest = revisionSelfImprovementCycleJEVExecutionPlanReverseObservationDigest(output)
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
	output.PlanSignal = jevExecutionPlanReverseComplete
	output.ObservationDigest = revisionSelfImprovementCycleJEVExecutionPlanReverseObservationDigest(output)
	return output
}

func (value RevisionSelfImprovementCycleJEVExecutionPlanReverseObservation) Validate() error {
	if value.Status != "UNKNOWN" && value.Status != "BOUND" {
		return fmt.Errorf("invalid status %q", value.Status)
	}
	if value.MetricName != jevExecutionPlanReverseMetricName {
		return fmt.Errorf("invalid metric name %q", value.MetricName)
	}
	if !value.NonExecuting || !value.NonAuthorizing {
		return fmt.Errorf("execution plan reverse observation must remain non-executing and non-authorizing")
	}
	if value.Status == "UNKNOWN" {
		if value.MissingStage == "" || value.PlanSignal != jevExecutionPlanReverseUnknown {
			return fmt.Errorf("unknown execution plan reverse observation must preserve its unresolved stage")
		}
	} else {
		if value.MissingStage != "" || value.ContractName != jevExecutionPlanContractName {
			return fmt.Errorf("bound execution plan reverse observation must identify its contract")
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
		if value.PlanSignal != jevExecutionPlanReverseComplete {
			return fmt.Errorf("invalid bound execution plan reverse signal %q", value.PlanSignal)
		}
	}
	if value.ObservationDigest != revisionSelfImprovementCycleJEVExecutionPlanReverseObservationDigest(value) {
		return fmt.Errorf("observation digest mismatch")
	}
	return nil
}

func revisionSelfImprovementCycleJEVExecutionPlanReverseObservationDigest(
	value RevisionSelfImprovementCycleJEVExecutionPlanReverseObservation,
) string {
	return digestString(fmt.Sprintf(
		"jev-execution-plan-reverse|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t",
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
