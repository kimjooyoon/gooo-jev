package decision

import "strings"

// ExecutionEnvelopeCapabilityBoundarySourceInput connects exact declaration
// and capability policy sources without granting any capability.
type ExecutionEnvelopeCapabilityBoundarySourceInput struct {
	DeclarationID          string
	ContractID             string
	DeclarationSource      string
	SandboxProfile         string
	WorkspaceSource        string
	GatewayPolicySource    string
	ModelIdentity          string
	SPIFFEIdentity         string
	NetworkAllowlistSource string
	NonAuthorizing         bool
}

// ExecutionEnvelopeCapabilityBoundarySourceBinding records declaration and
// policy source digests alongside the existing bounded capability digest.
type ExecutionEnvelopeCapabilityBoundarySourceBinding struct {
	Status                 string
	DeclarationDigest      string
	BoundaryDigest         string
	WorkspaceDigest        string
	GatewayPolicyDigest    string
	NetworkAllowlistDigest string
	MissingCapability      string
	NonExecuting           bool
	NonAuthorizing         bool
}

func digestCapabilitySource(kind, source string) (string, error) {
	return Digest(struct {
		Kind   string
		Source string
	}{Kind: kind, Source: source})
}

// SealExecutionEnvelopeCapabilityBoundaryFromSource derives all policy
// evidence from exact source text before reusing the capability boundary seal.
func SealExecutionEnvelopeCapabilityBoundaryFromSource(input ExecutionEnvelopeCapabilityBoundarySourceInput) ExecutionEnvelopeCapabilityBoundarySourceBinding {
	output := ExecutionEnvelopeCapabilityBoundarySourceBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	declaration := ComputeExecutionEnvelopeDeclarationSourceDigest(ExecutionEnvelopeDeclarationSourceDigestInput{
		DeclarationID:  input.DeclarationID,
		ContractID:     input.ContractID,
		SourceText:     input.DeclarationSource,
		NonAuthorizing: input.NonAuthorizing,
	})
	if declaration.Status != "derived" {
		output.MissingCapability = declaration.MissingStage
		output.NonAuthorizing = declaration.NonAuthorizing
		return output
	}
	if strings.TrimSpace(input.SandboxProfile) == "" {
		output.MissingCapability = "sandbox-profile"
		return output
	}
	if strings.TrimSpace(input.ModelIdentity) == "" {
		output.MissingCapability = "model-identity"
		return output
	}
	if strings.TrimSpace(input.SPIFFEIdentity) == "" {
		output.MissingCapability = "spiffe-identity"
		return output
	}
	workspaceDigest, err := digestCapabilitySource("workspace", input.WorkspaceSource)
	if err != nil || strings.TrimSpace(input.WorkspaceSource) == "" {
		output.MissingCapability = "workspace"
		return output
	}
	gatewayDigest, err := digestCapabilitySource("gateway-policy", input.GatewayPolicySource)
	if err != nil || strings.TrimSpace(input.GatewayPolicySource) == "" {
		output.MissingCapability = "gateway-policy"
		return output
	}
	networkDigest, err := digestCapabilitySource("network-allowlist", input.NetworkAllowlistSource)
	if err != nil || strings.TrimSpace(input.NetworkAllowlistSource) == "" {
		output.MissingCapability = "network-allowlist"
		return output
	}
	boundary := SealExecutionEnvelopeCapabilityBoundary(ExecutionEnvelopeCapabilityBoundaryInput{
		SandboxProfile:         input.SandboxProfile,
		WorkspaceDigest:        workspaceDigest,
		GatewayPolicyDigest:    gatewayDigest,
		ModelIdentity:          input.ModelIdentity,
		SPIFFEIdentity:         input.SPIFFEIdentity,
		NetworkAllowlistDigest: networkDigest,
		NonAuthorizing:         true,
	})
	output.Status = boundary.Status
	output.DeclarationDigest = declaration.DeclarationDigest
	output.BoundaryDigest = boundary.BoundaryDigest
	output.WorkspaceDigest = workspaceDigest
	output.GatewayPolicyDigest = gatewayDigest
	output.NetworkAllowlistDigest = networkDigest
	output.MissingCapability = boundary.MissingCapability
	return output
}
