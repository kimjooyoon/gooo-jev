package decision

import (
	"testing"
	"time"
)

func executionFloat64Ptr(value float64) *float64 {
	return &value
}

func TestExecutionEnvelopeBindsGrantIdentityAndTerminalResult(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	boundary, err := NewCapabilityBoundary(
		"spiffe://example.org/ns/prod/sa/jev",
		"decision-review",
		"decision:review",
		"evidence-digest-1",
		now.Add(-time.Minute),
		now.Add(time.Hour),
	)
	if err != nil {
		t.Fatalf("NewCapabilityBoundary() error = %v", err)
	}
	identity, err := NewWorkloadIdentityObservation(
		"spiffe://example.org/ns/prod/sa/jev",
		"identity-evidence-1",
		WorkloadIdentityObserved,
		now,
	)
	if err != nil {
		t.Fatalf("NewWorkloadIdentityObservation() error = %v", err)
	}
	request := CapabilityRequest{
		Subject:    "spiffe://example.org/ns/prod/sa/jev",
		Audience:   "decision-review",
		Capability: "decision:review",
	}
	grant, err := NewCapabilityGrant(request, boundary, identity, "grant-evidence-1")
	if err != nil {
		t.Fatalf("NewCapabilityGrant() error = %v", err)
	}
	spec := Spec{
		ID:             "review",
		Question:       "Should this candidate be reviewed?",
		Kind:           KindChoice,
		AllowedChoices: []string{"accept", "review"},
		PolicyDigest:   "policy-digest-1",
	}
	result := Result{
		SpecID:          "review",
		Kind:            KindChoice,
		Value:           Value{Choice: "review"},
		Confidence:      executionFloat64Ptr(0.9),
		EvidenceDigest:  "jev-evidence-1",
		Provider:        "provider-neutral",
		Status:          StatusObserved,
		ObservedAt:      now,
	}
	decisionReceipt, err := Observe(spec, State{Digest: "state-digest-1"}, result)
	if err != nil {
		t.Fatalf("Observe() error = %v", err)
	}
	execution, err := NewExecutionReceipt(
		grant,
		boundary,
		identity,
		decisionReceipt,
		"result-digest-1",
		ExecutionCompleted,
		now,
	)
	if err != nil {
		t.Fatalf("NewExecutionReceipt() error = %v", err)
	}
	if execution.Status != ExecutionCompleted {
		t.Fatalf("execution status = %q", execution.Status)
	}
	if err := execution.Validate(); err != nil {
		t.Fatalf("execution Validate() error = %v", err)
	}
	reverse, err := NewReverseObservation(
		execution,
		"output-digest-1",
		"verifier-digest-1",
		"",
		ProvenanceObserved,
		now,
	)
	if err != nil {
		t.Fatalf("NewReverseObservation() error = %v", err)
	}
	if err := reverse.Validate(); err != nil {
		t.Fatalf("reverse Validate() error = %v", err)
	}
}

func TestExecutionEnvelopeFailsClosedWithUnknownStages(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	boundary, err := NewCapabilityBoundary(
		"spiffe://example.org/ns/prod/sa/jev",
		"decision-review",
		"decision:review",
		"evidence-digest-1",
		now.Add(time.Hour),
		now.Add(2*time.Hour),
	)
	if err != nil {
		t.Fatalf("NewCapabilityBoundary() error = %v", err)
	}
	identity, err := NewWorkloadIdentityObservation(
		"spiffe://example.org/ns/prod/sa/jev",
		"identity-evidence-1",
		WorkloadIdentityObserved,
		now,
	)
	if err != nil {
		t.Fatalf("NewWorkloadIdentityObservation() error = %v", err)
	}
	unknown, err := NewExecutionReceipt(
		CapabilityGrant{},
		boundary,
		identity,
		Receipt{},
		"",
		ExecutionCompleted,
		now,
	)
	if err != nil {
		t.Fatalf("NewExecutionReceipt() unknown error = %v", err)
	}
	if unknown.Status != ExecutionUnknown || unknown.UnknownReason != "missing-or-invalid-grant" {
		t.Fatalf("unknown execution = %#v", unknown)
	}
	if err := unknown.Validate(); err != nil {
		t.Fatalf("unknown execution Validate() error = %v", err)
	}
	reverse, err := NewReverseObservation(
		unknown,
		"",
		"",
		"",
		ProvenanceObserved,
		now,
	)
	if err != nil {
		t.Fatalf("NewReverseObservation() unknown error = %v", err)
	}
	if reverse.Status != ProvenanceUnknown || reverse.MissingStage != "terminal-execution" {
		t.Fatalf("unknown reverse = %#v", reverse)
	}
	if err := reverse.Validate(); err != nil {
		t.Fatalf("unknown reverse Validate() error = %v", err)
	}
}

func TestExecutionEnvelopeRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0).UTC()
	identity, err := NewWorkloadIdentityObservation(
		"spiffe://example.org/ns/prod/sa/jev",
		"identity-evidence-1",
		WorkloadIdentityObserved,
		now,
	)
	if err != nil {
		t.Fatalf("NewWorkloadIdentityObservation() error = %v", err)
	}
	tampered := identity
	tampered.SPIFFEID = "spiffe://example.org/ns/prod/sa/other"
	if err := tampered.Validate(); err == nil {
		t.Fatal("tampered identity unexpectedly validated")
	}

	request := CapabilityRequest{
		Subject:    "spiffe://example.org/ns/prod/sa/jev",
		Audience:   "decision-review",
		Capability: "decision:review",
	}
	if err := request.Validate(); err != nil {
		t.Fatalf("CapabilityRequest.Validate() error = %v", err)
	}
}

