package gooo

import "testing"

func TestParseWorkloadIdentityBindsSpiffeEvidence(t *testing.T) {
	evidence := digestString("spiffe-attestation")
	identity, err := ParseWorkloadIdentity("spiffe://example.org/ns/prod/sa/gooo", evidence)
	if err != nil {
		t.Fatalf("ParseWorkloadIdentity() error = %v", err)
	}
	if identity.Status != "BOUND" || identity.MissingStage != "" {
		t.Fatalf("unexpected identity: %#v", identity)
	}
	if identity.Scheme != "spiffe" || identity.TrustDomain != "example.org" || identity.Path != "/ns/prod/sa/gooo" {
		t.Fatalf("unexpected identity fields: %#v", identity)
	}
	if identity.IdentityDigest == "" || !identity.NonExecuting || !identity.NonAuthorizing {
		t.Fatalf("missing identity safety evidence: %#v", identity)
	}
	if err := identity.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestParseWorkloadIdentityRejectsNonSpiffeScheme(t *testing.T) {
	identity, err := ParseWorkloadIdentity("https://example.org/workload", digestString("evidence"))
	if err == nil {
		t.Fatal("ParseWorkloadIdentity() error = nil, want scheme rejection")
	}
	if identity.Status != "UNKNOWN" || identity.MissingStage != "identity-scheme" {
		t.Fatalf("unexpected scheme failure: %#v", identity)
	}
}

func TestParseWorkloadIdentityRetainsMissingEvidenceStage(t *testing.T) {
	identity, err := ParseWorkloadIdentity("spiffe://example.org/workload", "")
	if err == nil {
		t.Fatal("ParseWorkloadIdentity() error = nil, want evidence rejection")
	}
	if identity.Status != "UNKNOWN" || identity.MissingStage != "identity-evidence" {
		t.Fatalf("unexpected evidence failure: %#v", identity)
	}
}

func TestParseWorkloadIdentityRejectsPathTraversal(t *testing.T) {
	identity, err := ParseWorkloadIdentity("spiffe://example.org/ns/prod/../sa/gooo", digestString("evidence"))
	if err == nil {
		t.Fatal("ParseWorkloadIdentity() error = nil, want path rejection")
	}
	if identity.MissingStage != "identity-path" {
		t.Fatalf("unexpected path failure: %#v", identity)
	}
}

func TestValidateRejectsTamperedWorkloadIdentity(t *testing.T) {
	identity, err := ParseWorkloadIdentity("spiffe://example.org/ns/prod/sa/gooo", digestString("evidence"))
	if err != nil {
		t.Fatalf("ParseWorkloadIdentity() error = %v", err)
	}
	identity.TrustDomain = "tampered.example"
	if err := identity.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want tamper rejection")
	}
}

func TestParseWorkloadIdentityIsDeterministic(t *testing.T) {
	evidence := digestString("evidence")
	first, err := ParseWorkloadIdentity("spiffe://example.org/ns/prod/sa/gooo", evidence)
	if err != nil {
		t.Fatalf("first ParseWorkloadIdentity() error = %v", err)
	}
	second, err := ParseWorkloadIdentity("spiffe://example.org/ns/prod/sa/gooo", evidence)
	if err != nil {
		t.Fatalf("second ParseWorkloadIdentity() error = %v", err)
	}
	if first.IdentityDigest != second.IdentityDigest {
		t.Fatal("same identity produced different digest")
	}
}
