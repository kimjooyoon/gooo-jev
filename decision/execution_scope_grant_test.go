package decision

import (
	"testing"
	"time"
)

func TestGrantExecutionScopeBindsCapabilityAndDecision(t *testing.T) {
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
	grant := CapabilityGrant{
		RequestDigest:            "request-digest",
		CapabilityBoundaryDigest: "boundary-digest",
		WorkloadIdentityDigest:   "identity-digest",
		EvidenceDigest:           "grant-evidence-digest",
		Granted:                  true,
	}
	if err := func() error {
		digest, err := grant.computeDigest()
		grant.GrantDigest = digest
		return err
	}(); err != nil {
		t.Fatalf("grant digest error = %v", err)
	}
	receipt := Receipt{
		Schema:         SchemaV1,
		SpecDigest:     "spec-digest",
		StateDigest:    "state-digest",
		ResultDigest:   "result-digest",
		PolicyDigest:   "policy-digest",
		Provider:       "provider",
		Status:         StatusObserved,
		ObservedAt:     now,
		NonAuthorizing: true,
		DecisionDigest: "decision-digest",
	}
	binding, err := GrantExecutionScope(scope, grant, receipt, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("GrantExecutionScope() error = %v", err)
	}
	if binding.Status != ExecutionScopeGranted || binding.MissingStage != "" {
		t.Fatalf("binding = %#v", binding)
	}
	if err := binding.Validate(); err != nil {
		t.Fatalf("binding Validate() error = %v", err)
	}
}

func TestGrantExecutionScopePreservesUnknownStage(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	scope, err := BindExecutionScope("task-digest", "workspace-digest", "", "model-digest", "network-policy-digest", "filesystem-policy-digest", now)
	if err != nil {
		t.Fatalf("BindExecutionScope() error = %v", err)
	}
	binding, err := GrantExecutionScope(scope, CapabilityGrant{}, Receipt{}, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("GrantExecutionScope() error = %v", err)
	}
	if binding.Status != ExecutionScopeGrantUnknown || binding.MissingStage != "scope:gateway" {
		t.Fatalf("unknown binding = %#v", binding)
	}
	if err := binding.Validate(); err != nil {
		t.Fatalf("unknown binding Validate() error = %v", err)
	}
}

func TestExecutionScopeGrantRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	binding := ExecutionScopeGrant{
		Schema:                ExecutionScopeGrantSchemaV1,
		ScopeDigest:           "scope-digest",
		CapabilityGrantDigest: "grant-digest",
		DecisionDigest:        "decision-digest",
		Status:                ExecutionScopeGranted,
		GrantedAt:             now,
	}
	if err := binding.assignDigest(); err != nil {
		t.Fatalf("binding.assignDigest() error = %v", err)
	}
	binding.DecisionDigest = "tampered-decision"
	if err := binding.Validate(); err == nil {
		t.Fatal("tampered execution scope grant unexpectedly validated")
	}
}
