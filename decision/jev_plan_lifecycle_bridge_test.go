package decision

import (
	"testing"
	"time"
)

func TestBindExecutionEnvelopeJEVExecutionPlanToLifecycle(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	plan := DeriveExecutionEnvelopeJEVExecutionPlanFromSource(ExecutionEnvelopeJEVExecutionPlanSourceInput{
		PlanID:              "triage-plan",
		ContractID:          "gooo://gooo-jev/triage-plan",
		TaskSource:          "task: classify support ticket",
		WorkspaceSource:     "workspace: read-only ticket cache",
		GatewaySource:       "gateway: approved decision provider",
		ModelSource:         "model: jev-latest",
		SuspendResumeSource: "lifecycle: checkpoint and receipt",
		NonAuthorizing:      true,
	})
	grant, boundary, identity, decisionReceipt := lifecycleInputs(t, now)
	suspension, resumption, reverse := lifecycleReplayInputs(t, now)
	replay, err := NewExecutionLifecycleReplay(suspension, resumption, reverse, now.Add(3*time.Minute))
	if err != nil {
		t.Fatalf("NewExecutionLifecycleReplay() error = %v", err)
	}
	authorization, err := NewExecutionLifecycleAuthorization(replay, resumption, grant, boundary, identity, decisionReceipt, now.Add(4*time.Minute))
	if err != nil {
		t.Fatalf("NewExecutionLifecycleAuthorization() error = %v", err)
	}
	binding := BindExecutionEnvelopeJEVExecutionPlanToLifecycle(ExecutionEnvelopeJEVExecutionPlanLifecycleInput{
		PlanBinding: plan, Replay: replay, Authorization: authorization, NonAuthorizing: true,
	})
	if binding.Status != "bound" || binding.MissingStage != "" {
		t.Fatalf("binding = %#v, want bound", binding)
	}
	if binding.PlanDigest == "" || binding.ReplayDigest == "" || binding.AuthorizationDigest == "" || binding.BindingDigest == "" {
		t.Fatal("bound lifecycle bridge must retain all identity digests")
	}
	changedPlan := plan
	changedPlan.PlanDigest = "changed-plan-digest"
	changed := BindExecutionEnvelopeJEVExecutionPlanToLifecycle(ExecutionEnvelopeJEVExecutionPlanLifecycleInput{
		PlanBinding: changedPlan, Replay: replay, Authorization: authorization, NonAuthorizing: true,
	})
	if changed.Status != "bound" || changed.BindingDigest == binding.BindingDigest {
		t.Fatal("changed plan digest must change the lifecycle bridge digest")
	}
	unknownReplay, err := NewExecutionLifecycleReplay(ExecutionLifecycleReceipt{}, ExecutionLifecycleReceipt{}, ReverseObservation{}, now)
	if err != nil {
		t.Fatalf("unknown replay construction error = %v", err)
	}
	unknown := BindExecutionEnvelopeJEVExecutionPlanToLifecycle(ExecutionEnvelopeJEVExecutionPlanLifecycleInput{
		PlanBinding: plan, Replay: unknownReplay, Authorization: authorization, NonAuthorizing: true,
	})
	if unknown.Status != "UNKNOWN" || unknown.MissingStage != "suspension" {
		t.Fatalf("unknown replay bridge = %#v, want UNKNOWN at suspension", unknown)
	}
	unauthorized := BindExecutionEnvelopeJEVExecutionPlanToLifecycle(ExecutionEnvelopeJEVExecutionPlanLifecycleInput{
		PlanBinding: plan, Replay: replay, Authorization: authorization, NonAuthorizing: false,
	})
	if unauthorized.Status != "UNKNOWN" || unauthorized.NonAuthorizing || unauthorized.MissingStage != "authorization-boundary" {
		t.Fatalf("unauthorized lifecycle bridge = %#v, want non-authorizing UNKNOWN", unauthorized)
	}
}
