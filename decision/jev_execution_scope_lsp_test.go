package decision

import (
	"testing"
	"time"
)

func TestProjectExecutionScopeLSP(t *testing.T) {
	scope, err := BindExecutionScope(
		"task-digest",
		"workspace-digest",
		"gateway-digest",
		"model-digest",
		"network-policy-digest",
		"filesystem-policy-digest",
		time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC),
	)
	if err != nil {
		t.Fatal(err)
	}
	got := ProjectExecutionScopeLSP(ExecutionScopeLSPInput{
		Scope:          scope,
		NonAuthorizing: true,
	})
	if got.Status != "bound" || got.Code != "JEV_EXECUTION_SCOPE_BOUND" || got.ScopeDigest != scope.ScopeDigest || got.NetworkPolicyDigest != "network-policy-digest" || got.FilesystemPolicyDigest != "filesystem-policy-digest" || got.EvidenceDigest == "" || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("got %+v", got)
	}
}

func TestProjectExecutionScopeLSPPreservesUnknownStage(t *testing.T) {
	got := ProjectExecutionScopeLSP(ExecutionScopeLSPInput{
		Scope: ExecutionScope{
			Status:      ExecutionScopeUnknown,
			MissingStage: "filesystem-policy",
		},
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "filesystem-policy" || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestProjectExecutionScopeLSPRejectsAuthorization(t *testing.T) {
	got := ProjectExecutionScopeLSP(ExecutionScopeLSPInput{})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}
