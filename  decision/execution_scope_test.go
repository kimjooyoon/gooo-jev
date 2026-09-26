package decision

import (
	"testing"
	"time"
)

func TestBindExecutionScopeBindsRuntimeBoundaries(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	scope, err := BindExecutionScope(
		"task-digest",
		"workspace-digest",
		"gateway-digest",
		"model-digest",
		"network-policy-digest",
		"filesystem-policy-digest",
		now,
	)
	if err != nil {
		t.Fatalf("BindExecutionScope() error = %v", err)
	}
	if scope.Status != ExecutionScopeBound || scope.MissingStage != "" {
		t.Fatalf("scope = %#v", scope)
	}
	if err := scope.Validate(); err != nil {
		t.Fatalf("scope Validate() error = %v", err)
	}
}

func TestBindExecutionScopePreservesUnknownStage(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	scope, err := BindExecutionScope(
		"task-digest",
		"workspace-digest",
		"",
		"model-digest",
		"network-policy-digest",
		"filesystem-policy-digest",
		now,
	)
	if err != nil {
		t.Fatalf("BindExecutionScope() error = %v", err)
	}
	if scope.Status != ExecutionScopeUnknown || scope.MissingStage != "gateway" {
		t.Fatalf("unknown scope = %#v", scope)
	}
	if err := scope.Validate(); err != nil {
		t.Fatalf("unknown scope Validate() error = %v", err)
	}
}

func TestExecutionScopeRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	scope, err := BindExecutionScope(
		"task-digest",
		"workspace-digest",
		"gateway-digest",
		"model-digest",
		"network-policy-digest",
		"filesystem-policy-digest",
		now,
	)
	if err != nil {
		t.Fatalf("BindExecutionScope() error = %v", err)
	}
	scope.GatewayDigest = "tampered-gateway"
	if err := scope.Validate(); err == nil {
		t.Fatal("tampered execution scope unexpectedly validated")
	}
}
