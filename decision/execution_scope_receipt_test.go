package decision

import (
	"testing"
	"time"
)

func TestRecordScopedExecutionClosesProvenanceBoundary(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	grant := ExecutionScopeGrant{
		Schema:                 ExecutionScopeGrantSchemaV1,
		ScopeDigest:            "scope-digest",
		CapabilityGrantDigest:  "capability-grant-digest",
		DecisionDigest:         "decision-digest",
		Status:                 ExecutionScopeGranted,
		GrantedAt:              now,
	}
	if err := grant.assignDigest(); err != nil {
		t.Fatalf("grant.assignDigest() error = %v", err)
	}
	execution := ExecutionReceipt{
		Schema:                    ExecutionReceiptSchemaV1,
		GrantDigest:               "capability-grant-digest",
		CapabilityBoundaryDigest:  "boundary-digest",
		DecisionReceiptDigest:     "decision-digest",
		WorkloadIdentityDigest:    "identity-digest",
		ResultDigest:              "result-digest",
		Status:                    ExecutionCompleted,
		ObservedAt:                now,
	}
	if err := func() error {
		digest, err := execution.computeDigest()
		execution.ReceiptDigest = digest
		return err
	}(); err != nil {
		t.Fatalf("execution digest error = %v", err)
	}
	receipt, err := RecordScopedExecution(grant, execution, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("RecordScopedExecution() error = %v", err)
	}
	if receipt.Status != ExecutionScopeExecutionCompleted || receipt.MissingStage != "" {
		t.Fatalf("receipt = %#v", receipt)
	}
	if err := receipt.Validate(); err != nil {
		t.Fatalf("receipt Validate() error = %v", err)
	}
}

func TestRecordScopedExecutionPreservesNonTerminalExecution(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	grant := ExecutionScopeGrant{
		Schema:                ExecutionScopeGrantSchemaV1,
		ScopeDigest:           "scope-digest",
		CapabilityGrantDigest: "capability-grant-digest",
		DecisionDigest:        "decision-digest",
		Status:                ExecutionScopeGranted,
		GrantedAt:             now,
	}
	if err := grant.assignDigest(); err != nil {
		t.Fatalf("grant.assignDigest() error = %v", err)
	}
	execution := ExecutionReceipt{
		Schema:                    ExecutionReceiptSchemaV1,
		GrantDigest:               "capability-grant-digest",
		CapabilityBoundaryDigest:  "boundary-digest",
		DecisionReceiptDigest:     "decision-digest",
		WorkloadIdentityDigest:    "identity-digest",
		Status:                    ExecutionNonTerminal,
		UnknownReason:             "waiting-for-gateway",
		ObservedAt:                now,
	}
	if err := func() error {
		digest, err := execution.computeDigest()
		execution.ReceiptDigest = digest
		return err
	}(); err != nil {
		t.Fatalf("execution digest error = %v", err)
	}
	receipt, err := RecordScopedExecution(grant, execution, now.Add(time.Minute))
	if err != nil {
		t.Fatalf("RecordScopedExecution() error = %v", err)
	}
	if receipt.Status != ExecutionScopeReceiptUnknown || receipt.MissingStage != "terminal-execution" {
		t.Fatalf("unknown receipt = %#v", receipt)
	}
	if err := receipt.Validate(); err != nil {
		t.Fatalf("unknown receipt Validate() error = %v", err)
	}
}

func TestExecutionScopeReceiptRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	receipt := ExecutionScopeReceipt{
		Schema:                 ExecutionScopeReceiptSchemaV1,
		ScopeGrantDigest:       "scope-grant-digest",
		ExecutionReceiptDigest: "execution-receipt-digest",
		Status:                 ExecutionScopeExecutionCompleted,
		ObservedAt:             now,
	}
	if err := receipt.assignDigest(); err != nil {
		t.Fatalf("receipt.assignDigest() error = %v", err)
	}
	receipt.ExecutionReceiptDigest = "tampered-execution"
	if err := receipt.Validate(); err == nil {
		t.Fatal("tampered execution scope receipt unexpectedly validated")
	}
}
