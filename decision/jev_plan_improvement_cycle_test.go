package decision

import "testing"

func TestBindExecutionEnvelopeJEVPlanToImprovementCycle(t *testing.T) {
	cycleDisposition := cycleObservationDisposition()
	cycleLSP := ProjectDecisionConfidenceChangePlanFeedbackDispositionLSP(cycleDisposition)
	cycle := ObserveJEVImprovementCycle(JEVImprovementCycleObservationInput{
		DeclarationDigest:        "declaration-digest",
		IRDigest:                 "ir-digest",
		GenerationDigest:         "generation-digest",
		ReverseObservationDigest: "reverse-observation-digest",
		MetricDigest:             "metric-digest",
		Disposition:              cycleDisposition,
		LSPDiagnostic:             cycleLSP,
		NonAuthorizing:           true,
	})
	if err := cycle.Validate(); err != nil {
		t.Fatalf("cycle Validate() error = %v", err)
	}
	planLifecycle := ExecutionEnvelopeJEVExecutionPlanLifecycleBinding{
		Status: "bound", PlanID: "triage-plan", PlanDigest: "plan-digest",
		BindingDigest: "lifecycle-binding-digest",
		NonExecuting: true, NonAuthorizing: true,
	}
	binding := BindExecutionEnvelopeJEVPlanToImprovementCycle(ExecutionEnvelopeJEVPlanImprovementCycleInput{
		PlanLifecycle: planLifecycle, Cycle: cycle, NonAuthorizing: true,
	})
	if binding.Status != "bound" || binding.MissingStage != "" {
		t.Fatalf("binding = %#v, want bound", binding)
	}
	if binding.BindingDigest == "" || binding.PlanDigest == "" || binding.CycleEvidenceDigest == "" {
		t.Fatal("bound cycle bridge must retain plan and cycle evidence")
	}
	changedPlan := planLifecycle
	changedPlan.BindingDigest = "changed-lifecycle-binding"
	changed := BindExecutionEnvelopeJEVPlanToImprovementCycle(ExecutionEnvelopeJEVPlanImprovementCycleInput{
		PlanLifecycle: changedPlan, Cycle: cycle, NonAuthorizing: true,
	})
	if changed.Status != "bound" || changed.BindingDigest == binding.BindingDigest {
		t.Fatal("changed lifecycle identity must change cycle binding digest")
	}
	tamperedCycle := cycle
	tamperedCycle.EvidenceDigest = "tampered-cycle-evidence"
	cycleUnknown := BindExecutionEnvelopeJEVPlanToImprovementCycle(ExecutionEnvelopeJEVPlanImprovementCycleInput{
		PlanLifecycle: planLifecycle, Cycle: tamperedCycle, NonAuthorizing: true,
	})
	if cycleUnknown.Status != "UNKNOWN" || cycleUnknown.MissingStage != "improvement-cycle" {
		t.Fatalf("tampered cycle = %#v, want UNKNOWN at improvement-cycle", cycleUnknown)
	}
	planUnknown := planLifecycle
	planUnknown.Status = "UNKNOWN"
	planUnknown.MissingStage = "suspend-resume-source"
	planResult := BindExecutionEnvelopeJEVPlanToImprovementCycle(ExecutionEnvelopeJEVPlanImprovementCycleInput{
		PlanLifecycle: planUnknown, Cycle: cycle, NonAuthorizing: true,
	})
	if planResult.Status != "UNKNOWN" || planResult.MissingStage != "suspend-resume-source" {
		t.Fatalf("unknown plan = %#v, want preserved missing stage", planResult)
	}
	unauthorized := BindExecutionEnvelopeJEVPlanToImprovementCycle(ExecutionEnvelopeJEVPlanImprovementCycleInput{
		PlanLifecycle: planLifecycle, Cycle: cycle, NonAuthorizing: false,
	})
	if unauthorized.Status != "UNKNOWN" || unauthorized.NonAuthorizing || unauthorized.MissingStage != "authorization-boundary" {
		t.Fatalf("unauthorized cycle = %#v, want non-authorizing UNKNOWN", unauthorized)
	}
}
