package decision

import "strings"

// ExecutionEnvelopeJEVPlanImprovementCycleInput connects the execution-plan
// lifecycle identity to a validated self-improvement cycle observation.
type ExecutionEnvelopeJEVPlanImprovementCycleInput struct {
	PlanLifecycle  ExecutionEnvelopeJEVExecutionPlanLifecycleBinding
	Cycle          JEVImprovementCycleObservation
	NonAuthorizing bool
}

// ExecutionEnvelopeJEVPlanImprovementCycleBinding records the plan and cycle
// identities without claiming that a revision was executed or authorized.
type ExecutionEnvelopeJEVPlanImprovementCycleBinding struct {
	Status                 string
	PlanID                 string
	PlanDigest             string
	LifecycleBindingDigest string
	CycleStatus            string
	CycleEvidenceDigest    string
	BindingDigest          string
	MissingStage           string
	NonExecuting           bool
	NonAuthorizing         bool
}

// BindExecutionEnvelopeJEVPlanToImprovementCycle admits only validated plan
// lifecycle and self-improvement cycle evidence into one immutable identity.
func BindExecutionEnvelopeJEVPlanToImprovementCycle(input ExecutionEnvelopeJEVPlanImprovementCycleInput) ExecutionEnvelopeJEVPlanImprovementCycleBinding {
	output := ExecutionEnvelopeJEVPlanImprovementCycleBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.PlanLifecycle.NonAuthorizing || !input.Cycle.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if input.PlanLifecycle.Status != "bound" || strings.TrimSpace(input.PlanLifecycle.BindingDigest) == "" {
		output.MissingStage = input.PlanLifecycle.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "plan-lifecycle"
		}
		return output
	}
	if err := input.Cycle.Validate(); err != nil {
		output.MissingStage = "improvement-cycle"
		return output
	}
	if strings.TrimSpace(input.Cycle.Status) == "" || strings.TrimSpace(input.Cycle.EvidenceDigest) == "" {
		output.MissingStage = "improvement-cycle"
		return output
	}
	bindingDigest, err := Digest(struct {
		PlanID                 string
		PlanDigest             string
		LifecycleBindingDigest string
		CycleStatus            string
		CycleEvidenceDigest    string
	}{
		PlanID:                 input.PlanLifecycle.PlanID,
		PlanDigest:             input.PlanLifecycle.PlanDigest,
		LifecycleBindingDigest: input.PlanLifecycle.BindingDigest,
		CycleStatus:            input.Cycle.Status,
		CycleEvidenceDigest:    input.Cycle.EvidenceDigest,
	})
	if err != nil {
		output.MissingStage = "improvement-cycle-binding-digest"
		return output
	}
	output.Status = "bound"
	output.PlanID = input.PlanLifecycle.PlanID
	output.PlanDigest = input.PlanLifecycle.PlanDigest
	output.LifecycleBindingDigest = input.PlanLifecycle.BindingDigest
	output.CycleStatus = input.Cycle.Status
	output.CycleEvidenceDigest = input.Cycle.EvidenceDigest
	output.BindingDigest = bindingDigest
	return output
}
