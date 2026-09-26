package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const ExecutionLifecycleSchemaV1 = "gooo/jev-execution-lifecycle/v1"

type ExecutionLifecycleStatus string

const (
	ExecutionLifecycleSuspended ExecutionLifecycleStatus = "suspended"
	ExecutionLifecycleResumed  ExecutionLifecycleStatus = "resumed"
	ExecutionLifecycleUnknown  ExecutionLifecycleStatus = "unknown"
)

type ExecutionLifecycleReceipt struct {
	Schema                    string
	GrantDigest               string
	CapabilityBoundaryDigest  string
	WorkloadIdentityDigest    string
	DecisionReceiptDigest     string
	CheckpointDigest          string
	StateDigest               string
	ParentLifecycleDigest     string
	EvidenceDigest            string
	Status                    ExecutionLifecycleStatus
	UnknownReason             string
	ObservedAt                time.Time
	LifecycleDigest           string
}

func NewExecutionSuspension(
	grant CapabilityGrant,
	boundary CapabilityBoundary,
	identity WorkloadIdentityObservation,
	decisionReceipt Receipt,
	checkpointDigest string,
	stateDigest string,
	evidenceDigest string,
	observedAt time.Time,
) (ExecutionLifecycleReceipt, error) {
	if observedAt.IsZero() {
		return ExecutionLifecycleReceipt{}, errors.New("execution lifecycle observation time is required")
	}
	receipt := ExecutionLifecycleReceipt{
		Schema:                   ExecutionLifecycleSchemaV1,
		GrantDigest:              grant.GrantDigest,
		CapabilityBoundaryDigest: boundary.BoundaryDigest,
		WorkloadIdentityDigest:   identity.ObservationDigest,
		DecisionReceiptDigest:    decisionReceipt.DecisionDigest,
		CheckpointDigest:         checkpointDigest,
		StateDigest:              stateDigest,
		EvidenceDigest:           evidenceDigest,
		Status:                   ExecutionLifecycleSuspended,
		ObservedAt:               observedAt.UTC(),
	}
	if reason := validateExecutionLifecycleInputs(grant, boundary, identity, decisionReceipt, checkpointDigest, stateDigest, evidenceDigest, observedAt); reason != "" {
		receipt.Status = ExecutionLifecycleUnknown
		receipt.UnknownReason = reason
	}
	if err := receipt.assignDigest(); err != nil {
		return ExecutionLifecycleReceipt{}, fmt.Errorf("digest execution suspension: %w", err)
	}
	if err := receipt.Validate(); err != nil {
		return ExecutionLifecycleReceipt{}, err
	}
	return receipt, nil
}

func NewExecutionResumption(
	suspension ExecutionLifecycleReceipt,
	grant CapabilityGrant,
	boundary CapabilityBoundary,
	identity WorkloadIdentityObservation,
	decisionReceipt Receipt,
	stateDigest string,
	evidenceDigest string,
	observedAt time.Time,
) (ExecutionLifecycleReceipt, error) {
	if observedAt.IsZero() {
		return ExecutionLifecycleReceipt{}, errors.New("execution lifecycle observation time is required")
	}
	receipt := ExecutionLifecycleReceipt{
		Schema:                   ExecutionLifecycleSchemaV1,
		GrantDigest:              grant.GrantDigest,
		CapabilityBoundaryDigest: boundary.BoundaryDigest,
		WorkloadIdentityDigest:   identity.ObservationDigest,
		DecisionReceiptDigest:    decisionReceipt.DecisionDigest,
		CheckpointDigest:         suspension.CheckpointDigest,
		StateDigest:              stateDigest,
		ParentLifecycleDigest:    suspension.LifecycleDigest,
		EvidenceDigest:           evidenceDigest,
		Status:                   ExecutionLifecycleResumed,
		ObservedAt:               observedAt.UTC(),
	}
	switch {
	case suspension.Validate() != nil:
		receipt.Status = ExecutionLifecycleUnknown
		receipt.UnknownReason = "missing-or-invalid-suspension"
	case suspension.Status != ExecutionLifecycleSuspended:
		receipt.Status = ExecutionLifecycleUnknown
		receipt.UnknownReason = "suspension-not-active"
	case suspension.GrantDigest != grant.GrantDigest ||
		suspension.CapabilityBoundaryDigest != boundary.BoundaryDigest ||
		suspension.WorkloadIdentityDigest != identity.ObservationDigest ||
		suspension.DecisionReceiptDigest != decisionReceipt.DecisionDigest:
		receipt.Status = ExecutionLifecycleUnknown
		receipt.UnknownReason = "suspension-binding-mismatch"
	default:
		if reason := validateExecutionLifecycleInputs(grant, boundary, identity, decisionReceipt, suspension.CheckpointDigest, stateDigest, evidenceDigest, observedAt); reason != "" {
			receipt.Status = ExecutionLifecycleUnknown
			receipt.UnknownReason = reason
		}
	}
	if err := receipt.assignDigest(); err != nil {
		return ExecutionLifecycleReceipt{}, fmt.Errorf("digest execution resumption: %w", err)
	}
	if err := receipt.Validate(); err != nil {
		return ExecutionLifecycleReceipt{}, err
	}
	return receipt, nil
}

func validateExecutionLifecycleInputs(
	grant CapabilityGrant,
	boundary CapabilityBoundary,
	identity WorkloadIdentityObservation,
	decisionReceipt Receipt,
	checkpointDigest string,
	stateDigest string,
	evidenceDigest string,
	observedAt time.Time,
) string {
	if err := grant.Validate(); err != nil {
		return "missing-or-invalid-grant"
	}
	if err := boundary.Validate(); err != nil || boundary.StatusAt(observedAt) != CapabilityBoundaryActive {
		return "inactive-or-invalid-capability-boundary"
	}
	if err := identity.Validate(); err != nil || identity.Status != WorkloadIdentityObserved {
		return "missing-or-invalid-workload-identity"
	}
	if grant.CapabilityBoundaryDigest != boundary.BoundaryDigest || grant.WorkloadIdentityDigest != identity.ObservationDigest {
		return "grant-binding-mismatch"
	}
	if err := decisionReceipt.Validate(); err != nil {
		return "missing-or-invalid-decision-receipt"
	}
	if strings.TrimSpace(checkpointDigest) == "" {
		return "missing-checkpoint"
	}
	if strings.TrimSpace(stateDigest) == "" {
		return "missing-state"
	}
	if strings.TrimSpace(evidenceDigest) == "" {
		return "missing-evidence"
	}
	return ""
}

func (receipt *ExecutionLifecycleReceipt) assignDigest() error {
	digest, err := receipt.computeDigest()
	if err != nil {
		return err
	}
	receipt.LifecycleDigest = digest
	return nil
}

func (receipt ExecutionLifecycleReceipt) Validate() error {
	if err := receipt.validateShape(); err != nil {
		return err
	}
	expected, err := receipt.computeDigest()
	if err != nil {
		return fmt.Errorf("digest execution lifecycle receipt: %w", err)
	}
	if receipt.LifecycleDigest != expected {
		return errors.New("execution lifecycle receipt digest mismatch")
	}
	return nil
}

func (receipt ExecutionLifecycleReceipt) validateShape() error {
	if receipt.Schema != ExecutionLifecycleSchemaV1 {
		return errors.New("unsupported execution lifecycle schema")
	}
	if receipt.ObservedAt.IsZero() {
		return errors.New("execution lifecycle observation time is required")
	}
	if strings.TrimSpace(receipt.LifecycleDigest) == "" {
		return errors.New("execution lifecycle digest is required")
	}
	switch receipt.Status {
	case ExecutionLifecycleSuspended:
		if strings.TrimSpace(receipt.GrantDigest) == "" ||
			strings.TrimSpace(receipt.CapabilityBoundaryDigest) == "" ||
			strings.TrimSpace(receipt.WorkloadIdentityDigest) == "" ||
			strings.TrimSpace(receipt.DecisionReceiptDigest) == "" ||
			strings.TrimSpace(receipt.CheckpointDigest) == "" ||
			strings.TrimSpace(receipt.StateDigest) == "" ||
			strings.TrimSpace(receipt.EvidenceDigest) == "" ||
			strings.TrimSpace(receipt.ParentLifecycleDigest) != "" {
			return errors.New("suspended execution lifecycle receipt is incomplete")
		}
	case ExecutionLifecycleResumed:
		if strings.TrimSpace(receipt.GrantDigest) == "" ||
			strings.TrimSpace(receipt.CapabilityBoundaryDigest) == "" ||
			strings.TrimSpace(receipt.WorkloadIdentityDigest) == "" ||
			strings.TrimSpace(receipt.DecisionReceiptDigest) == "" ||
			strings.TrimSpace(receipt.CheckpointDigest) == "" ||
			strings.TrimSpace(receipt.StateDigest) == "" ||
			strings.TrimSpace(receipt.ParentLifecycleDigest) == "" ||
			strings.TrimSpace(receipt.EvidenceDigest) == "" {
			return errors.New("resumed execution lifecycle receipt is incomplete")
		}
	case ExecutionLifecycleUnknown:
		if strings.TrimSpace(receipt.UnknownReason) == "" {
			return errors.New("unknown execution lifecycle receipt requires an explicit reason")
		}
	default:
		return fmt.Errorf("unsupported execution lifecycle status %q", receipt.Status)
	}
	return nil
}

func (receipt ExecutionLifecycleReceipt) computeDigest() (string, error) {
	return Digest(struct {
		Schema                   string
		GrantDigest              string
		CapabilityBoundaryDigest string
		WorkloadIdentityDigest   string
		DecisionReceiptDigest    string
		CheckpointDigest         string
		StateDigest              string
		ParentLifecycleDigest    string
		EvidenceDigest            string
		Status                   ExecutionLifecycleStatus
		UnknownReason             string
		ObservedAt               time.Time
	}{
		Schema:                   receipt.Schema,
		GrantDigest:              receipt.GrantDigest,
		CapabilityBoundaryDigest: receipt.CapabilityBoundaryDigest,
		WorkloadIdentityDigest:   receipt.WorkloadIdentityDigest,
		DecisionReceiptDigest:    receipt.DecisionReceiptDigest,
		CheckpointDigest:         receipt.CheckpointDigest,
		StateDigest:              receipt.StateDigest,
		ParentLifecycleDigest:    receipt.ParentLifecycleDigest,
		EvidenceDigest:           receipt.EvidenceDigest,
		Status:                   receipt.Status,
		UnknownReason:            receipt.UnknownReason,
		ObservedAt:               receipt.ObservedAt.UTC(),
	})
}
