package decision

import "testing"

func capabilityBoundaryLSPPolicy() ExecutionEnvelopeCapabilityBoundaryInput {
	return ExecutionEnvelopeCapabilityBoundaryInput{
		SandboxProfile:         "sandbox-v1",
		WorkspaceDigest:        "workspace-digest",
		GatewayPolicyDigest:    "gateway-digest",
		ModelIdentity:          "model-v1",
		SPIFFEIdentity:         "spiffe://example/agent",
		NetworkAllowlistDigest: "network-digest",
		NonAuthorizing:         true,
	}
}

func TestProjectExecutionEnvelopeCapabilityBoundaryLSP(t *testing.T) {
	policy := capabilityBoundaryLSPPolicy()
	boundary := SealExecutionEnvelopeCapabilityBoundary(policy)
	got := ProjectExecutionEnvelopeCapabilityBoundaryLSP(ExecutionEnvelopeCapabilityBoundaryLSPInput{
		Boundary:       boundary,
		Policy:         policy,
		NonAuthorizing: true,
	})
	if got.Status != "bounded" || got.Code != "JEV_CAPABILITY_BOUNDARY_BOUNDED" || got.BoundaryDigest != boundary.BoundaryDigest || got.SPIFFEIdentity != policy.SPIFFEIdentity || got.EvidenceDigest == "" || !got.NonAuthorizing {
		t.Fatalf("got %+v", got)
	}
}

func TestProjectExecutionEnvelopeCapabilityBoundaryLSPRejectsTampering(t *testing.T) {
	policy := capabilityBoundaryLSPPolicy()
	boundary := SealExecutionEnvelopeCapabilityBoundary(policy)
	boundary.BoundaryDigest = "tampered"
	got := ProjectExecutionEnvelopeCapabilityBoundaryLSP(ExecutionEnvelopeCapabilityBoundaryLSPInput{
		Boundary:       boundary,
		Policy:         policy,
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingCapability != "boundary-integrity" || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestProjectExecutionEnvelopeCapabilityBoundaryLSPRejectsAuthorization(t *testing.T) {
	got := ProjectExecutionEnvelopeCapabilityBoundaryLSP(ExecutionEnvelopeCapabilityBoundaryLSPInput{})
	if got.Status != "UNKNOWN" || got.MissingCapability != "authorization-boundary" || got.NonAuthorizing || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}
