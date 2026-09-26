package decision

import "testing"

func TestSealExecutionEnvelopeCapabilityBoundaryFromSource(t *testing.T) {
	input := ExecutionEnvelopeCapabilityBoundarySourceInput{
		DeclarationID:          "capability-declaration",
		ContractID:             "execution-capability-contract",
		DeclarationSource:      "entity SafeExecution id gooo://safe-execution",
		SandboxProfile:         "restricted",
		WorkspaceSource:        "workspace: read-only /repo",
		GatewayPolicySource:    "gateway: approved-provider-only",
		ModelIdentity:          "jev-latest",
		SPIFFEIdentity:         "spiffe://example.org/ns/jev/sa/runner",
		NetworkAllowlistSource: "network: api.example.org:443",
		NonAuthorizing:         true,
	}
	bounded := SealExecutionEnvelopeCapabilityBoundaryFromSource(input)
	if bounded.Status != "bounded" || bounded.MissingCapability != "" {
		t.Fatalf("bounded = %#v, want bounded", bounded)
	}
	if bounded.DeclarationDigest == "" || bounded.BoundaryDigest == "" || bounded.WorkspaceDigest == "" ||
		bounded.GatewayPolicyDigest == "" || bounded.NetworkAllowlistDigest == "" {
		t.Fatal("bounded capability boundary must retain declaration and policy digests")
	}
	changedNetwork := input
	changedNetwork.NetworkAllowlistSource += ", backup.example.org:443"
	changed := SealExecutionEnvelopeCapabilityBoundaryFromSource(changedNetwork)
	if changed.Status != "bounded" || changed.BoundaryDigest == bounded.BoundaryDigest || changed.NetworkAllowlistDigest == bounded.NetworkAllowlistDigest {
		t.Fatal("changed allowlist source must change the bounded capability evidence")
	}
	missingSPIFFE := input
	missingSPIFFE.SPIFFEIdentity = ""
	spiffeUnknown := SealExecutionEnvelopeCapabilityBoundaryFromSource(missingSPIFFE)
	if spiffeUnknown.Status != "UNKNOWN" || spiffeUnknown.MissingCapability != "spiffe-identity" {
		t.Fatalf("missing SPIFFE identity = %#v, want UNKNOWN at spiffe-identity", spiffeUnknown)
	}
	missingNetwork := input
	missingNetwork.NetworkAllowlistSource = ""
	networkUnknown := SealExecutionEnvelopeCapabilityBoundaryFromSource(missingNetwork)
	if networkUnknown.Status != "UNKNOWN" || networkUnknown.MissingCapability != "network-allowlist" {
		t.Fatalf("missing network source = %#v, want UNKNOWN at network-allowlist", networkUnknown)
	}
	unauthorized := input
	unauthorized.NonAuthorizing = false
	unauthorizedBoundary := SealExecutionEnvelopeCapabilityBoundaryFromSource(unauthorized)
	if unauthorizedBoundary.Status != "UNKNOWN" || unauthorizedBoundary.NonAuthorizing || unauthorizedBoundary.MissingCapability != "authorization-boundary" {
		t.Fatalf("unauthorized capability source = %#v, want non-authorizing UNKNOWN", unauthorizedBoundary)
	}
}
