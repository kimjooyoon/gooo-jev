package decision

import "strings"

// ExecutionEnvelopeJEVExecutionPlanSourceInput describes the exact source
// boundaries of a declarative JEV execution plan.
type ExecutionEnvelopeJEVExecutionPlanSourceInput struct {
	PlanID              string
	ContractID          string
	TaskSource          string
	WorkspaceSource     string
	GatewaySource       string
	ModelSource         string
	SuspendResumeSource string
	NonAuthorizing      bool
}

// ExecutionEnvelopeJEVExecutionPlanSourceBinding records each source digest
// and a deterministic plan digest without authorizing execution.
type ExecutionEnvelopeJEVExecutionPlanSourceBinding struct {
	Status                  string
	PlanID                  string
	ContractID              string
	TaskDigest              string
	WorkspaceDigest         string
	GatewayDigest           string
	ModelDigest             string
	SuspendResumeDigest     string
	PlanDigest              string
	MissingStage            string
	NonExecuting            bool
	NonAuthorizing          bool
}

func digestJEVExecutionPlanSource(kind, source string) (string, error) {
	return Digest(struct {
		Kind   string
		Source string
	}{Kind: kind, Source: source})
}

// DeriveExecutionEnvelopeJEVExecutionPlanFromSource seals the plan's source
// boundaries while keeping missing or unauthorized evidence UNKNOWN.
func DeriveExecutionEnvelopeJEVExecutionPlanFromSource(input ExecutionEnvelopeJEVExecutionPlanSourceInput) ExecutionEnvelopeJEVExecutionPlanSourceBinding {
	output := ExecutionEnvelopeJEVExecutionPlanSourceBinding{
		Status: "UNKNOWN", PlanID: input.PlanID, ContractID: input.ContractID,
		NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if strings.TrimSpace(input.PlanID) == "" {
		output.MissingStage = "plan-id"
		return output
	}
	if strings.TrimSpace(input.ContractID) == "" {
		output.MissingStage = "contract-id"
		return output
	}
	stages := []struct {
		name  string
		value string
		set   func(string)
	}{
		{name: "task-source", value: input.TaskSource, set: func(value string) { output.TaskDigest = value }},
		{name: "workspace-source", value: input.WorkspaceSource, set: func(value string) { output.WorkspaceDigest = value }},
		{name: "gateway-source", value: input.GatewaySource, set: func(value string) { output.GatewayDigest = value }},
		{name: "model-source", value: input.ModelSource, set: func(value string) { output.ModelDigest = value }},
		{name: "suspend-resume-source", value: input.SuspendResumeSource, set: func(value string) { output.SuspendResumeDigest = value }},
	}
	for _, stage := range stages {
		if strings.TrimSpace(stage.value) == "" {
			output.MissingStage = stage.name
			return output
		}
		digest, err := digestJEVExecutionPlanSource(stage.name, stage.value)
		if err != nil {
			output.MissingStage = stage.name + "-digest"
			return output
		}
		stage.set(digest)
	}
	planDigest, err := Digest(struct {
		PlanID              string
		ContractID          string
		TaskDigest          string
		WorkspaceDigest     string
		GatewayDigest       string
		ModelDigest         string
		SuspendResumeDigest string
	}{
		PlanID:              input.PlanID,
		ContractID:          input.ContractID,
		TaskDigest:          output.TaskDigest,
		WorkspaceDigest:     output.WorkspaceDigest,
		GatewayDigest:       output.GatewayDigest,
		ModelDigest:         output.ModelDigest,
		SuspendResumeDigest: output.SuspendResumeDigest,
	})
	if err != nil {
		output.MissingStage = "plan-digest"
		return output
	}
	output.Status = "derived"
	output.PlanDigest = planDigest
	return output
}
