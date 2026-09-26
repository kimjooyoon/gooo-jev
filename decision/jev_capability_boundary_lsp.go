package decision

import "strings"

// ExecutionEnvelopeCapabilityBoundaryLSPInput carries the original policy
// inputs and their sealed result for read-only LSP integrity checking.
type ExecutionEnvelopeCapabilityBoundaryLSPInput struct {
	Boundary       ExecutionEnvelopeCapabilityBoundary
	Policy         ExecutionEnvelopeCapabilityBoundaryInput
	NonAuthorizing bool
}

// ExecutionEnvelopeCapabilityBoundaryLSPProjection exposes security policy
// identities only after re-sealing and comparing the boundary digest.
type ExecutionEnvelopeCapabilityBoundaryLSPProjection struct {
	Status                 string
	Code                   string
	Severity               string
	Message                string
	SandboxProfile         string
	ModelIdentity          string
	SPIFFEIdentity         string
	WorkspaceDigest        string
	GatewayPolicyDigest    string
	NetworkAllowlistDigest string
	BoundaryDigest         string
	EvidenceDigest         string
	MissingCapability      string
	NonAuthorizing         bool
}

func digestExecutionEnvelopeCapabilityBoundaryLSPProjection(projection ExecutionEnvelopeCapabilityBoundaryLSPProjection) (string, error) {
	return Digest(struct {
		Status                 string
		Code                   string
		Severity               string
		Message                string
		SandboxProfile         string
		ModelIdentity           string
		SPIFFEIdentity          string
		WorkspaceDigest         string
		GatewayPolicyDigest     string
		NetworkAllowlistDigest  string
		BoundaryDigest          string
		MissingCapability       string
	}{
		Status:                 projection.Status,
		Code:                   projection.Code,
		Severity:                projection.Severity,
		Message:                 projection.Message,
		SandboxProfile:          projection.SandboxProfile,
		ModelIdentity:           projection.ModelIdentity,
		SPIFFEIdentity:          projection.SPIFFEIdentity,
		WorkspaceDigest:         projection.WorkspaceDigest,
		GatewayPolicyDigest:     projection.GatewayPolicyDigest,
		NetworkAllowlistDigest:  projection.NetworkAllowlistDigest,
		BoundaryDigest:          projection.BoundaryDigest,
		MissingCapability:       projection.MissingCapability,
	})
}

// ProjectExecutionEnvelopeCapabilityBoundaryLSP validates the sealed policy
// before publishing security evidence and never grants a capability.
func ProjectExecutionEnvelopeCapabilityBoundaryLSP(input ExecutionEnvelopeCapabilityBoundaryLSPInput) ExecutionEnvelopeCapabilityBoundaryLSPProjection {
	output := ExecutionEnvelopeCapabilityBoundaryLSPProjection{
		Status: "UNKNOWN", Code: "JEV_CAPABILITY_BOUNDARY_UNKNOWN", Severity: "warning",
		NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Policy.NonAuthorizing || !input.Boundary.NonAuthorizing {
		if !input.NonAuthorizing || !input.Policy.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingCapability = "authorization-boundary"
		output.Message = "JEV capability boundary is UNKNOWN: missing authorization boundary"
		return finalizeExecutionEnvelopeCapabilityBoundaryLSP(output)
	}
	if input.Boundary.Status != "bounded" || strings.TrimSpace(input.Boundary.BoundaryDigest) == "" {
		output.MissingCapability = input.Boundary.MissingCapability
		if strings.TrimSpace(output.MissingCapability) == "" {
			output.MissingCapability = "capability-boundary"
		}
		output.Message = "JEV capability boundary is UNKNOWN: missing " + output.MissingCapability
		return finalizeExecutionEnvelopeCapabilityBoundaryLSP(output)
	}
	expected := SealExecutionEnvelopeCapabilityBoundary(input.Policy)
	if expected.Status != "bounded" || expected.BoundaryDigest != input.Boundary.BoundaryDigest {
		output.MissingCapability = "boundary-integrity"
		output.Message = "JEV capability boundary is UNKNOWN: invalid boundary-integrity"
		return finalizeExecutionEnvelopeCapabilityBoundaryLSP(output)
	}
	output.Status = "bounded"
	output.Code = "JEV_CAPABILITY_BOUNDARY_BOUNDED"
	output.Severity = "info"
	output.Message = "JEV capability boundary is sealed for review; no capability is granted by this projection"
	output.SandboxProfile = input.Policy.SandboxProfile
	output.ModelIdentity = input.Policy.ModelIdentity
	output.SPIFFEIdentity = input.Policy.SPIFFEIdentity
	output.WorkspaceDigest = input.Policy.WorkspaceDigest
	output.GatewayPolicyDigest = input.Policy.GatewayPolicyDigest
	output.NetworkAllowlistDigest = input.Policy.NetworkAllowlistDigest
	output.BoundaryDigest = input.Boundary.BoundaryDigest
	return finalizeExecutionEnvelopeCapabilityBoundaryLSP(output)
}

func finalizeExecutionEnvelopeCapabilityBoundaryLSP(output ExecutionEnvelopeCapabilityBoundaryLSPProjection) ExecutionEnvelopeCapabilityBoundaryLSPProjection {
	digest, err := digestExecutionEnvelopeCapabilityBoundaryLSPProjection(output)
	if err != nil {
		output.Status = "UNKNOWN"
		output.Code = "JEV_CAPABILITY_BOUNDARY_UNKNOWN"
		output.Severity = "warning"
		output.MissingCapability = "lsp-evidence-digest"
		output.Message = "JEV capability boundary is UNKNOWN: missing lsp-evidence-digest"
		output.EvidenceDigest = ""
		return output
	}
	output.EvidenceDigest = digest
	return output
}
