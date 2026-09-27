package gooo

import "testing"

func securityWorkloadIdentityBoundaryInput(t *testing.T) RevisionSelfImprovementCycleSecurityWorkloadIdentityBoundaryInput {
	t.Helper()
	return RevisionSelfImprovementCycleSecurityWorkloadIdentityBoundaryInput{
		WorkloadID:   "spiffe://gooo.example/workload/compiler",
		IssuerID:     "spiffe://gooo.example/trust/issuer",
		AudienceID:   "spiffe://gooo.example/service/lsp",
		PolicyDigest: digestString("gooo-security-policy-v1"),
	}
}

func TestObserveRevisionSelfImprovementCycleSecurityWorkloadIdentityBoundaryBinds(t *testing.T) {
	input := securityWorkloadIdentityBoundaryInput(t)
	got, err := ObserveRevisionSelfImprovementCycleSecurityWorkloadIdentityBoundary(input)
	if err != nil {
		t.Fatalf("observe security workload identity boundary: %v", err)
	}
	if got.Status != "BOUND" ||
		got.BoundarySignal != "workload-identity-boundary-observed" ||
		got.WorkloadID != input.WorkloadID ||
		got.IssuerID != input.IssuerID ||
		got.AudienceID != input.AudienceID ||
		got.PolicyDigest != input.PolicyDigest {
		t.Fatalf("unexpected security workload identity boundary: %#v", got)
	}
	if !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("security boundary safety flags = %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate security workload identity boundary: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleSecurityWorkloadIdentityBoundaryPreservesUnknown(t *testing.T) {
	input := securityWorkloadIdentityBoundaryInput(t)
	input.WorkloadID = "https://untrusted.example/compiler"
	got, err := ObserveRevisionSelfImprovementCycleSecurityWorkloadIdentityBoundary(input)
	if err == nil {
		t.Fatal("expected invalid SPIFFE ID error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "revision-self-improvement-cycle-security-workload-identity-boundary-input" ||
		got.BoundarySignal != "workload-identity-boundary-unknown" {
		t.Fatalf("unexpected unknown security workload identity boundary: %#v", got)
	}
}

func TestObserveRevisionSelfImprovementCycleSecurityWorkloadIdentityBoundaryRejectsTamperedDigest(t *testing.T) {
	input := securityWorkloadIdentityBoundaryInput(t)
	got, err := ObserveRevisionSelfImprovementCycleSecurityWorkloadIdentityBoundary(input)
	if err != nil {
		t.Fatalf("observe security workload identity boundary: %v", err)
	}
	got.BoundaryDigest = digestString("tampered")
	if err := got.Validate(); err == nil {
		t.Fatal("expected tampered boundary digest error")
	}
}

func TestObserveRevisionSelfImprovementCycleSecurityWorkloadIdentityBoundaryIsDeterministic(t *testing.T) {
	input := securityWorkloadIdentityBoundaryInput(t)
	first, err := ObserveRevisionSelfImprovementCycleSecurityWorkloadIdentityBoundary(input)
	if err != nil {
		t.Fatalf("first security workload identity boundary: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCycleSecurityWorkloadIdentityBoundary(input)
	if err != nil {
		t.Fatalf("second security workload identity boundary: %v", err)
	}
	if first != second {
		t.Fatalf("security workload identity boundaries differ: %#v != %#v", first, second)
	}
}
