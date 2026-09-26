package decision

import "testing"

func TestDeriveExecutionEnvelopeJEVExecutionPlanFromSource(t *testing.T) {
	input := ExecutionEnvelopeJEVExecutionPlanSourceInput{
		PlanID:              "triage-plan",
		ContractID:          "gooo://gooo-jev/triage-plan",
		TaskSource:          "task: classify support ticket",
		WorkspaceSource:     "workspace: read-only ticket cache",
		GatewaySource:       "gateway: approved decision provider",
		ModelSource:         "model: jev-latest",
		SuspendResumeSource: "lifecycle: suspend after checkpoint; resume with receipt",
		NonAuthorizing:      true,
	}
	derived := DeriveExecutionEnvelopeJEVExecutionPlanFromSource(input)
	if derived.Status != "derived" || derived.MissingStage != "" {
		t.Fatalf("derived = %#v, want derived", derived)
	}
	if derived.PlanDigest == "" || derived.TaskDigest == "" || derived.WorkspaceDigest == "" ||
		derived.GatewayDigest == "" || derived.ModelDigest == "" || derived.SuspendResumeDigest == "" {
		t.Fatal("derived plan must retain every source digest")
	}
	changedLifecycle := input
	changedLifecycle.SuspendResumeSource += "; cancel on timeout"
	changed := DeriveExecutionEnvelopeJEVExecutionPlanFromSource(changedLifecycle)
	if changed.Status != "derived" || changed.PlanDigest == derived.PlanDigest || changed.SuspendResumeDigest == derived.SuspendResumeDigest {
		t.Fatal("changed suspend/resume source must change plan provenance")
	}
	missingGateway := input
	missingGateway.GatewaySource = ""
	gatewayUnknown := DeriveExecutionEnvelopeJEVExecutionPlanFromSource(missingGateway)
	if gatewayUnknown.Status != "UNKNOWN" || gatewayUnknown.MissingStage != "gateway-source" {
		t.Fatalf("missing gateway = %#v, want UNKNOWN at gateway-source", gatewayUnknown)
	}
	missingLifecycle := input
	missingLifecycle.SuspendResumeSource = ""
	lifecycleUnknown := DeriveExecutionEnvelopeJEVExecutionPlanFromSource(missingLifecycle)
	if lifecycleUnknown.Status != "UNKNOWN" || lifecycleUnknown.MissingStage != "suspend-resume-source" {
		t.Fatalf("missing lifecycle = %#v, want UNKNOWN at suspend-resume-source", lifecycleUnknown)
	}
	unauthorized := input
	unauthorized.NonAuthorizing = false
	unauthorizedPlan := DeriveExecutionEnvelopeJEVExecutionPlanFromSource(unauthorized)
	if unauthorizedPlan.Status != "UNKNOWN" || unauthorizedPlan.NonAuthorizing || unauthorizedPlan.MissingStage != "authorization-boundary" {
		t.Fatalf("unauthorized plan = %#v, want non-authorizing UNKNOWN", unauthorizedPlan)
	}
}
