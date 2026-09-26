package decision

import (
	"strings"
	"testing"
)

func TestSealExecutionEnvelopeCapabilityBoundary(t *testing.T) {
	input := ExecutionEnvelopeCapabilityBoundaryInput{
		SandboxProfile:         "sandbox-v1",
		WorkspaceDigest:        "workspace-digest",
		GatewayPolicyDigest:    "gateway-digest",
		ModelIdentity:          "model-v1",
		SPIFFEIdentity:         "spiffe://example/agent",
		NetworkAllowlistDigest: "network-digest",
		NonAuthorizing:         true,
	}
	output := SealExecutionEnvelopeCapabilityBoundary(input)
	if output.Status != "bounded" || output.MissingCapability != "" ||
		output.BoundaryDigest == "" || output.NonAuthorizing != true {
		t.Fatalf("unexpected boundary: %#v", output)
	}
	if len(output.BoundaryDigest) != 64 || strings.Trim(output.BoundaryDigest, "0123456789abcdef") != "" {
		t.Fatalf("unexpected digest: %#v", output)
	}

	input.NetworkAllowlistDigest = ""
	output = SealExecutionEnvelopeCapabilityBoundary(input)
	if output.Status != "UNKNOWN" || output.MissingCapability != "network-allowlist" ||
		output.BoundaryDigest != "" {
		t.Fatalf("unexpected missing boundary: %#v", output)
	}

	input.NetworkAllowlistDigest = "network-digest"
	input.NonAuthorizing = false
	output = SealExecutionEnvelopeCapabilityBoundary(input)
	if output.Status != "UNKNOWN" || output.MissingCapability != "authorization-boundary" ||
		output.NonAuthorizing != false {
		t.Fatalf("unexpected authorization boundary: %#v", output)
	}
}