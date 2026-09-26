package decision

import "strings"

// ExecutionEnvelopeJEVExecutionPlanLifecycleInput connects a derived JEV
// execution plan to already observed lifecycle evidence without authorizing it.
type ExecutionEnvelopeJEVExecutionPlanLifecycleInput struct {
	PlanBinding    ExecutionEnvelopeJEVExecutionPlanSourceBinding
	Replay         ExecutionLifecycleReplayObservation
	Authorization  ExecutionLifecycleAuthorization
	NonAuthorizing bool
}

// ExecutionEnvelopeJEVExecutionPlanLifecycleBinding records the immutable
// identity bridge across plan, replay, and authorization evidence.
type ExecutionEnvelopeJEVExecutionPlanLifecycleBinding struct {
	Status              string
	PlanID              string
	ContractID          string
	PlanDigest          string
	ReplayDigest        string
	AuthorizationDigest string
	BindingDigest       string
	MissingStage        string
	NonExecuting        bool
	NonAuthorizing      bool
}

// BindExecutionEnvelopeJEVExecutionPlanToLifecycle verifies existing plan,
// replay, and authorization evidence without granting new capabilities.
func BindExecutionEnvelopeJEVExecutionPlanToLifecycle(input ExecutionEnvelopeJEVExecutionPlanLifecycleInput) ExecutionEnvelopeJEVExecutionPlanLifecycleBinding {
	output := ExecutionEnvelopeJEVExecutionPlanLifecycleBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.PlanBinding.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if input.PlanBinding.Status != "derived" || strings.TrimSpace(input.PlanBinding.PlanDigest) == "" {
		output.MissingStage = input.PlanBinding.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "plan-source"
		}
		return output
	}
	if err := input.Replay.Validate(); err != nil {
		output.MissingStage = "replay"
		return output
	}
	if input.Replay.Status != ExecutionLifecycleReplayObserved {
		output.MissingStage = "replay-status"
		return output
	}
	if err := input.Authorization.Validate(); err != nil {
		output.MissingStage = "authorization"
		return output
	}
	if input.Authorization.Status != ExecutionLifecycleAuthorized {
		output.MissingStage = "authorization-status"
		return output
	}
	if input.Authorization.ReplayDigest != input.Replay.ReplayDigest {
		output.MissingStage = "lifecycle-binding"
		return output
	}
	bindingDigest, err := Digest(struct {
		PlanID              string
		ContractID          string
		PlanDigest          string
		ReplayDigest        string
		AuthorizationDigest string
	}{
		PlanID:              input.PlanBinding.PlanID,
		ContractID:          input.PlanBinding.ContractID,
		PlanDigest:          input.PlanBinding.PlanDigest,
		ReplayDigest:        input.Replay.ReplayDigest,
		AuthorizationDigest: input.Authorization.AuthorizationDigest,
	})
	if err != nil {
		output.MissingStage = "lifecycle-binding-digest"
		return output
	}
	output.Status = "bound"
	output.PlanID = input.PlanBinding.PlanID
	output.ContractID = input.PlanBinding.ContractID
	output.PlanDigest = input.PlanBinding.PlanDigest
	output.ReplayDigest = input.Replay.ReplayDigest
	output.AuthorizationDigest = input.Authorization.AuthorizationDigest
	output.BindingDigest = bindingDigest
	return output
}
