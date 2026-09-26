package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const ExecutionScopeSchemaV1 = "gooo/jev-execution-scope/v1"

type ExecutionScopeStatus string

const (
	ExecutionScopeBound   ExecutionScopeStatus = "bound"
	ExecutionScopeUnknown ExecutionScopeStatus = "unknown"
)

type ExecutionScope struct {
	Schema                  string
	TaskDigest              string
	WorkspaceDigest         string
	GatewayDigest            string
	ModelDigest              string
	NetworkPolicyDigest      string
	FilesystemPolicyDigest  string
	Status                  ExecutionScopeStatus
	MissingStage             string
	BoundAt                  time.Time
	ScopeDigest              string
}

func BindExecutionScope(
	taskDigest string,
	workspaceDigest string,
	gatewayDigest string,
	modelDigest string,
	networkPolicyDigest string,
	filesystemPolicyDigest string,
	boundAt time.Time,
) (ExecutionScope, error) {
	if boundAt.IsZero() {
		return ExecutionScope{}, errors.New("execution scope binding time is required")
	}
	scope := ExecutionScope{
		Schema:                  ExecutionScopeSchemaV1,
		TaskDigest:              strings.TrimSpace(taskDigest),
		WorkspaceDigest:         strings.TrimSpace(workspaceDigest),
		GatewayDigest:           strings.TrimSpace(gatewayDigest),
		ModelDigest:             strings.TrimSpace(modelDigest),
		NetworkPolicyDigest:     strings.TrimSpace(networkPolicyDigest),
		FilesystemPolicyDigest:  strings.TrimSpace(filesystemPolicyDigest),
		Status:                  ExecutionScopeBound,
		BoundAt:                 boundAt.UTC(),
	}
	for _, input := range []struct {
		value string
		stage string
	}{
		{value: scope.TaskDigest, stage: "task"},
		{value: scope.WorkspaceDigest, stage: "workspace"},
		{value: scope.GatewayDigest, stage: "gateway"},
		{value: scope.ModelDigest, stage: "model"},
		{value: scope.NetworkPolicyDigest, stage: "network-policy"},
		{value: scope.FilesystemPolicyDigest, stage: "filesystem-policy"},
	} {
		if input.value == "" {
			scope.Status = ExecutionScopeUnknown
			scope.MissingStage = input.stage
			break
		}
	}
	if err := scope.assignDigest(); err != nil {
		return ExecutionScope{}, fmt.Errorf("digest execution scope: %w", err)
	}
	if err := scope.Validate(); err != nil {
		return ExecutionScope{}, err
	}
	return scope, nil
}

func (scope *ExecutionScope) assignDigest() error {
	digest, err := scope.computeDigest()
	if err != nil {
		return err
	}
	scope.ScopeDigest = digest
	return nil
}

func (scope ExecutionScope) Validate() error {
	if err := scope.validateShape(); err != nil {
		return err
	}
	expected, err := scope.computeDigest()
	if err != nil {
		return fmt.Errorf("digest execution scope: %w", err)
	}
	if scope.ScopeDigest != expected {
		return errors.New("execution scope digest mismatch")
	}
	return nil
}

func (scope ExecutionScope) validateShape() error {
	if scope.Schema != ExecutionScopeSchemaV1 {
		return errors.New("unsupported execution scope schema")
	}
	if scope.BoundAt.IsZero() {
		return errors.New("execution scope binding time is required")
	}
	if strings.TrimSpace(scope.ScopeDigest) == "" {
		return errors.New("execution scope digest is required")
	}
	switch scope.Status {
	case ExecutionScopeBound:
		if strings.TrimSpace(scope.TaskDigest) == "" ||
			strings.TrimSpace(scope.WorkspaceDigest) == "" ||
			strings.TrimSpace(scope.GatewayDigest) == "" ||
			strings.TrimSpace(scope.ModelDigest) == "" ||
			strings.TrimSpace(scope.NetworkPolicyDigest) == "" ||
			strings.TrimSpace(scope.FilesystemPolicyDigest) == "" ||
			strings.TrimSpace(scope.MissingStage) != "" {
			return errors.New("bound execution scope is incomplete")
		}
	case ExecutionScopeUnknown:
		if strings.TrimSpace(scope.MissingStage) == "" {
			return errors.New("unknown execution scope requires a missing stage")
		}
	default:
		return fmt.Errorf("unsupported execution scope status %q", scope.Status)
	}
	return nil
}

func (scope ExecutionScope) computeDigest() (string, error) {
	return Digest(struct {
		Schema                 string
		TaskDigest             string
		WorkspaceDigest        string
		GatewayDigest           string
		ModelDigest             string
		NetworkPolicyDigest     string
		FilesystemPolicyDigest string
		Status                 ExecutionScopeStatus
		MissingStage            string
		BoundAt                 time.Time
	}{
		Schema:                  scope.Schema,
		TaskDigest:              scope.TaskDigest,
		WorkspaceDigest:         scope.WorkspaceDigest,
		GatewayDigest:            scope.GatewayDigest,
		ModelDigest:             scope.ModelDigest,
		NetworkPolicyDigest:     scope.NetworkPolicyDigest,
		FilesystemPolicyDigest:  scope.FilesystemPolicyDigest,
		Status:                  scope.Status,
		MissingStage:             scope.MissingStage,
		BoundAt:                 scope.BoundAt.UTC(),
	})
}
