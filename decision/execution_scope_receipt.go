package decision

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const ExecutionScopeReceiptSchemaV1 = "gooo/jev-execution-scope-receipt/v1"

type ExecutionScopeReceiptStatus string

const (
	ExecutionScopeExecutionCompleted ExecutionScopeReceiptStatus = "completed"
	ExecutionScopeReceiptUnknown      ExecutionScopeReceiptStatus = "unknown"
)

type ExecutionScopeReceipt struct {
	Schema                 string
	ScopeGrantDigest       string
	ExecutionReceiptDigest string
	Status                 ExecutionScopeReceiptStatus
	MissingStage           string
	ObservedAt             time.Time
	ReceiptDigest          string
}

func RecordScopedExecution(
	scopeGrant ExecutionScopeGrant,
	execution ExecutionReceipt,
	observedAt time.Time,
) (ExecutionScopeReceipt, error) {
	if observedAt.IsZero() {
		return ExecutionScopeReceipt{}, errors.New("execution scope receipt observation time is required")
	}
	receipt := ExecutionScopeReceipt{
		Schema:                 ExecutionScopeReceiptSchemaV1,
		ScopeGrantDigest:       scopeGrant.GrantDigest,
		ExecutionReceiptDigest: execution.ReceiptDigest,
		Status:                 ExecutionScopeExecutionCompleted,
		ObservedAt:             observedAt.UTC(),
	}
	switch {
	case scopeGrant.Validate() != nil:
		receipt.Status = ExecutionScopeReceiptUnknown
		receipt.MissingStage = "scope-grant"
	case scopeGrant.Status != ExecutionScopeGranted:
		receipt.Status = ExecutionScopeReceiptUnknown
		receipt.MissingStage = "scope-grant-status"
		if strings.TrimSpace(scopeGrant.MissingStage) != "" {
			receipt.MissingStage = "scope-grant:" + scopeGrant.MissingStage
		}
	case execution.Validate() != nil:
		receipt.Status = ExecutionScopeReceiptUnknown
		receipt.MissingStage = "execution"
	case execution.Status != ExecutionCompleted:
		receipt.Status = ExecutionScopeReceiptUnknown
		receipt.MissingStage = "terminal-execution"
	}
	if err := receipt.assignDigest(); err != nil {
		return ExecutionScopeReceipt{}, fmt.Errorf("digest execution scope receipt: %w", err)
	}
	if err := receipt.Validate(); err != nil {
		return ExecutionScopeReceipt{}, err
	}
	return receipt, nil
}

func (receipt *ExecutionScopeReceipt) assignDigest() error {
	digest, err := receipt.computeDigest()
	if err != nil {
		return err
	}
	receipt.ReceiptDigest = digest
	return nil
}

func (receipt ExecutionScopeReceipt) Validate() error {
	if err := receipt.validateShape(); err != nil {
		return err
	}
	expected, err := receipt.computeDigest()
	if err != nil {
		return fmt.Errorf("digest execution scope receipt: %w", err)
	}
	if receipt.ReceiptDigest != expected {
		return errors.New("execution scope receipt digest mismatch")
	}
	return nil
}

func (receipt ExecutionScopeReceipt) validateShape() error {
	if receipt.Schema != ExecutionScopeReceiptSchemaV1 {
		return errors.New("unsupported execution scope receipt schema")
	}
	if receipt.ObservedAt.IsZero() {
		return errors.New("execution scope receipt observation time is required")
	}
	if strings.TrimSpace(receipt.ReceiptDigest) == "" {
		return errors.New("execution scope receipt digest is required")
	}
	switch receipt.Status {
	case ExecutionScopeExecutionCompleted:
		if strings.TrimSpace(receipt.ScopeGrantDigest) == "" ||
			strings.TrimSpace(receipt.ExecutionReceiptDigest) == "" ||
			strings.TrimSpace(receipt.MissingStage) != "" {
			return errors.New("completed execution scope receipt is incomplete")
		}
	case ExecutionScopeReceiptUnknown:
		if strings.TrimSpace(receipt.MissingStage) == "" {
			return errors.New("unknown execution scope receipt requires a missing stage")
		}
	default:
		return fmt.Errorf("unsupported execution scope receipt status %q", receipt.Status)
	}
	return nil
}

func (receipt ExecutionScopeReceipt) computeDigest() (string, error) {
	return Digest(struct {
		Schema                 string
		ScopeGrantDigest       string
		ExecutionReceiptDigest string
		Status                 ExecutionScopeReceiptStatus
		MissingStage           string
		ObservedAt             time.Time
	}{
		Schema:                 receipt.Schema,
		ScopeGrantDigest:       receipt.ScopeGrantDigest,
		ExecutionReceiptDigest: receipt.ExecutionReceiptDigest,
		Status:                 receipt.Status,
		MissingStage:           receipt.MissingStage,
		ObservedAt:             receipt.ObservedAt.UTC(),
	})
}
