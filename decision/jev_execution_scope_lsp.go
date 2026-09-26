package decision

import (
	"strings"
	"time"
)

// ExecutionScopeLSPInput adapts a validated execution scope to an editor
// projection without granting the scope or executing a task.
type ExecutionScopeLSPInput struct {
	Scope          ExecutionScope
	NonAuthorizing bool
}

// ExecutionScopeLSPProjection preserves task, provider, network, and
// filesystem policy digests for security-aware editor consumers.
type ExecutionScopeLSPProjection struct {
	Status                  string
	Code                    string
	Severity                string
	Message                 string
	Schema                  string
	TaskDigest              string
	WorkspaceDigest         string
	GatewayDigest           string
	ModelDigest             string
	NetworkPolicyDigest     string
	FilesystemPolicyDigest  string
	ScopeDigest             string
	BoundAt                 string
	EvidenceDigest          string
	MissingStage             string
	NonExecuting            bool
	NonAuthorizing          bool
}

func digestExecutionScopeLSPProjection(projection ExecutionScopeLSPProjection) (string, error) {
	return Digest(struct {
		Status                  string
		Code                    string
		Severity                string
		Message                 string
		Schema                  string
		TaskDigest              string
		WorkspaceDigest         string
		GatewayDigest           string
		ModelDigest             string
		NetworkPolicyDigest     string
		FilesystemPolicyDigest  string
		ScopeDigest             string
		BoundAt                 string
		MissingStage             string
	}{
		Status:                 projection.Status,
		Code:                   projection.Code,
		Severity:                projection.Severity,
		Message:                projection.Message,
		Schema:                 projection.Schema,
		TaskDigest:             projection.TaskDigest,
		WorkspaceDigest:        projection.WorkspaceDigest,
		GatewayDigest:          projection.GatewayDigest,
		ModelDigest:            projection.ModelDigest,
		NetworkPolicyDigest:    projection.NetworkPolicyDigest,
		FilesystemPolicyDigest: projection.FilesystemPolicyDigest,
		ScopeDigest:            projection.ScopeDigest,
		BoundAt:                projection.BoundAt,
		MissingStage:           projection.MissingStage,
	})
}

// ProjectExecutionScopeLSP publishes a validated scope boundary and keeps
// incomplete or invalid network/filesystem policy evidence UNKNOWN.
func ProjectExecutionScopeLSP(input ExecutionScopeLSPInput) ExecutionScopeLSPProjection {
	output := ExecutionScopeLSPProjection{
		Status: "UNKNOWN", Code: "JEV_EXECUTION_SCOPE_UNKNOWN", Severity: "warning",
		NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		output.Message = "JEV execution scope is UNKNOWN: missing authorization boundary"
		return finalizeExecutionScopeLSP(output)
	}
	if input.Scope.Status != ExecutionScopeBound {
		output.MissingStage = input.Scope.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "execution-scope"
		}
		output.Message = "JEV execution scope is UNKNOWN: missing " + output.MissingStage
		return finalizeExecutionScopeLSP(output)
	}
	if err := input.Scope.Validate(); err != nil {
		output.MissingStage = "execution-scope-integrity"
		output.Message = "JEV execution scope is UNKNOWN: invalid execution-scope-integrity"
		return finalizeExecutionScopeLSP(output)
	}
	output.Status = "bound"
	output.Code = "JEV_EXECUTION_SCOPE_BOUND"
	output.Severity = "info"
	output.Message = "JEV execution scope is bound for review; no capability is granted by this projection"
	output.Schema = input.Scope.Schema
	output.TaskDigest = input.Scope.TaskDigest
	output.WorkspaceDigest = input.Scope.WorkspaceDigest
	output.GatewayDigest = input.Scope.GatewayDigest
	output.ModelDigest = input.Scope.ModelDigest
	output.NetworkPolicyDigest = input.Scope.NetworkPolicyDigest
	output.FilesystemPolicyDigest = input.Scope.FilesystemPolicyDigest
	output.ScopeDigest = input.Scope.ScopeDigest
	output.BoundAt = input.Scope.BoundAt.UTC().Format(time.RFC3339Nano)
	return finalizeExecutionScopeLSP(output)
}

func finalizeExecutionScopeLSP(output ExecutionScopeLSPProjection) ExecutionScopeLSPProjection {
	digest, err := digestExecutionScopeLSPProjection(output)
	if err != nil {
		output.Status = "UNKNOWN"
		output.Code = "JEV_EXECUTION_SCOPE_UNKNOWN"
		output.Severity = "warning"
		output.MissingStage = "lsp-evidence-digest"
		output.Message = "JEV execution scope is UNKNOWN: missing lsp-evidence-digest"
		output.EvidenceDigest = ""
		return output
	}
	output.EvidenceDigest = digest
	return output
}
