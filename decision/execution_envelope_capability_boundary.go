package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// ExecutionEnvelopeCapabilityBoundaryInput records the declared execution
// boundaries without granting any capability.
type ExecutionEnvelopeCapabilityBoundaryInput struct {
	SandboxProfile         string
	WorkspaceDigest        string
	GatewayPolicyDigest    string
	ModelIdentity          string
	SPIFFEIdentity         string
	NetworkAllowlistDigest string
	NonAuthorizing         bool
}

// ExecutionEnvelopeCapabilityBoundary records whether all declared execution
// boundary evidence is present.
type ExecutionEnvelopeCapabilityBoundary struct {
	Status            string
	BoundaryDigest    string
	MissingCapability string
	NonAuthorizing    bool
}

// SealExecutionEnvelopeCapabilityBoundary requires every boundary evidence
// value and produces a deterministic non-authorizing digest.
func SealExecutionEnvelopeCapabilityBoundary(input ExecutionEnvelopeCapabilityBoundaryInput) ExecutionEnvelopeCapabilityBoundary {
	output := ExecutionEnvelopeCapabilityBoundary{Status: "UNKNOWN", NonAuthorizing: true}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingCapability = "authorization-boundary"
		return output
	}
	values := []struct {
		name  string
		value string
	}{
		{name: "sandbox-profile", value: input.SandboxProfile},
		{name: "workspace", value: input.WorkspaceDigest},
		{name: "gateway-policy", value: input.GatewayPolicyDigest},
		{name: "model-identity", value: input.ModelIdentity},
		{name: "spiffe-identity", value: input.SPIFFEIdentity},
		{name: "network-allowlist", value: input.NetworkAllowlistDigest},
	}
	parts := make([]string, 0, len(values)*2)
	for _, boundary := range values {
		if boundary.value == "" {
			output.MissingCapability = boundary.name
			return output
		}
		parts = append(parts, boundary.name, boundary.value)
	}
	digest := sha256.Sum256([]byte(strings.Join(parts, string(rune(0)))))
	output.Status = "bounded"
	output.BoundaryDigest = hex.EncodeToString(digest[:])
	return output
}