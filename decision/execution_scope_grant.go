package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const ExecutionScopeGrantSchemaV1 = "gooo/jev-execution-scope-grant/v1"

type ExecutionScopeGrantStatus string

const (
	ExecutionScopeGranted       ExecutionScopeGrantStatus = "granted"
	ExecutionScopeGrantUnknown  ExecutionScopeGrantStatus = "unknown"
)

type ExecutionScopeGrant struct {
	Schema                 string
	ScopeDigest            string
	CapabilityGrantDigest  string
	DecisionDigest         string
	Status                 ExecutionScopeGrantStatus
	MissingStage           string
	GrantedAt              time.Time
	GrantDigest            string
}

func GrantExecutionScope(
	scope ExecutionScope,
	grant CapabilityGrant,
	decisionReceipt Receipt,
	grantedAt time.Time,
) (ExecutionScopeGrant, error) {
	if grantedAt.IsZero() {
		return ExecutionScopeGrant{}, errors.New("execution scope grant time is required")
	}
	binding := ExecutionScopeGrant{
		Schema:                ExecutionScopeGrantSchemaV1,
		ScopeDigest:           scope.ScopeDigest,
		CapabilityGrantDigest: grant.GrantDigest,
		DecisionDigest:        decisionReceipt.DecisionDigest,
		Status:                ExecutionScopeGranted,
		GrantedAt:             grantedAt.UTC(),
	}
	switch {
	case scope.Validate() != nil:
		binding.Status = ExecutionScopeGrantUnknown
		binding.MissingStage = "scope"
	case scope.Status != ExecutionScopeBound:
		binding.Status = ExecutionScopeGrantUnknown
		binding.MissingStage = "scope-status"
		if strings.TrimSpace(scope.MissingStage) != "" {
			binding.MissingStage = "scope:" + scope.MissingStage
		}
	case grant.Validate() != nil:
		binding.Status = ExecutionScopeGrantUnknown
		binding.MissingStage = "capability-grant"
	case decisionReceipt.Validate() != nil:
		binding.Status = ExecutionScopeGrantUnknown
		binding.MissingStage = "decision-receipt"
	}
	if err := binding.assignDigest(); err != nil {
		return ExecutionScopeGrant{}, fmt.Errorf("digest execution scope grant: %w", err)
	}
	if err := binding.Validate(); err != nil {
		return ExecutionScopeGrant{}, err
	}
	return binding, nil
}

func (binding *ExecutionScopeGrant) assignDigest() error {
	digest, err := binding.computeDigest()
	if err != nil {
		return err
	}
	binding.GrantDigest = digest
	return nil
}

func (binding ExecutionScopeGrant) Validate() error {
	if err := binding.validateShape(); err != nil {
		return err
	}
	expected, err := binding.computeDigest()
	if err != nil {
		return fmt.Errorf("digest execution scope grant: %w", err)
	}
	if binding.GrantDigest != expected {
		return errors.New("execution scope grant digest mismatch")
	}
	return nil
}

func (binding ExecutionScopeGrant) validateShape() error {
	if binding.Schema != ExecutionScopeGrantSchemaV1 {
		return errors.New("unsupported execution scope grant schema")
	}
	if binding.GrantedAt.IsZero() {
		return errors.New("execution scope grant time is required")
	}
	if strings.TrimSpace(binding.GrantDigest) == "" {
		return errors.New("execution scope grant digest is required")
	}
	switch binding.Status {
	case ExecutionScopeGranted:
		if strings.TrimSpace(binding.ScopeDigest) == "" ||
			strings.TrimSpace(binding.CapabilityGrantDigest) == "" ||
			strings.TrimSpace(binding.DecisionDigest) == "" ||
			strings.TrimSpace(binding.MissingStage) != "" {
			return errors.New("granted execution scope is incomplete")
		}
	case ExecutionScopeGrantUnknown:
		if strings.TrimSpace(binding.MissingStage) == "" {
			return errors.New("unknown execution scope grant requires a missing stage")
		}
	default:
		return fmt.Errorf("unsupported execution scope grant status %q", binding.Status)
	}
	return nil
}

func (binding ExecutionScopeGrant) computeDigest() (string, error) {
	return Digest(struct {
		Schema                string
		ScopeDigest           string
		CapabilityGrantDigest string
		DecisionDigest        string
		Status                ExecutionScopeGrantStatus
		MissingStage          string
		GrantedAt             time.Time
	}{
		Schema:                binding.Schema,
		ScopeDigest:           binding.ScopeDigest,
		CapabilityGrantDigest: binding.CapabilityGrantDigest,
		DecisionDigest:        binding.DecisionDigest,
		Status:                binding.Status,
		MissingStage:          binding.MissingStage,
		GrantedAt:             binding.GrantedAt.UTC(),
	})
}
