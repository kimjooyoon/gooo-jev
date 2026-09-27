package gooo

import "testing"

func testJEVSecurityWorkloadIdentityBoundaryInput() RevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundaryInput {
	return RevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundaryInput{
		WorkloadIdentity:                  "spiffe://example.test/workload/jev",
		WorkloadIdentityDigest:            digestString("workload-identity"),
		ExternalEvidenceDigest:            digestString("external-evidence"),
		Audience:                          "gooo-jev-observer",
		ExpiresAtUnix:                     1790000000,
		EvidenceStatus:                    "observed",
		PermissionState:                   "DEFER",
		GenerationTraceDigest:             digestString("generation-trace"),
		ReverseObservationCoverageDigest: digestString("reverse-observation-coverage"),
	}
}

func TestObserveRevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundary(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundary(
		testJEVSecurityWorkloadIdentityBoundaryInput(),
	)
	if observation.Status != "BOUND" || observation.MissingStage != "" {
		t.Fatalf("expected bounded security observation, got %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("expected valid security observation: %v", err)
	}
	if !observation.ReadOnly || !observation.NonExecuting || !observation.NonAuthorizing {
		t.Fatalf("security observation must not execute or authorize: %#v", observation)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundaryPreservesUnknown(
	t *testing.T,
) {
	input := testJEVSecurityWorkloadIdentityBoundaryInput()
	input.ExternalEvidenceDigest = ""
	observation := ObserveRevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundary(input)
	if observation.Status != "UNKNOWN" ||
		observation.MissingStage != "revision-self-improvement-cycle-jev-security-workload-identity-lineage" {
		t.Fatalf("expected missing lineage to remain unknown, got %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("expected valid unknown observation: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundaryRejectsTampering(
	t *testing.T,
) {
	observation := ObserveRevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundary(
		testJEVSecurityWorkloadIdentityBoundaryInput(),
	)
	observation.PermissionState = "ALLOW"
	if err := observation.Validate(); err == nil {
		t.Fatal("expected permission tampering to be rejected")
	}
	observation = ObserveRevisionSelfImprovementCycleJEVSecurityWorkloadIdentityBoundary(
		testJEVSecurityWorkloadIdentityBoundaryInput(),
	)
	observation.NonAuthorizing = false
	if err := observation.Validate(); err == nil {
		t.Fatal("expected authorization flag tampering to be rejected")
	}
}