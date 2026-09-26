package decision

// ExecutionEnvelopeProvenanceChainLSPBindingInput adapts the complete
// provenance chain to the existing editor diagnostic projection.
type ExecutionEnvelopeProvenanceChainLSPBindingInput struct {
	Chain                ExecutionEnvelopeProvenanceChainBinding
	MissingStageIndex    int
	EvidencePrefixDigest string
	NonAuthorizing       bool
}

// ExecutionEnvelopeProvenanceChainLSPBinding preserves the exact missing
// stage and chain evidence alongside the editor-facing diagnostic.
type ExecutionEnvelopeProvenanceChainLSPBinding struct {
	Status              string
	Publishable         bool
	Severity            string
	Code                string
	MissingStage        string
	MissingStageIndex   int
	EvidencePrefixDigest string
	ChainEvidenceDigest string
	NonAuthorizing      bool
}

// ProjectExecutionEnvelopeProvenanceChainLSP binds a chain result to LSP only
// when the stage index and evidence prefix make the diagnostic auditable.
func ProjectExecutionEnvelopeProvenanceChainLSP(input ExecutionEnvelopeProvenanceChainLSPBindingInput) ExecutionEnvelopeProvenanceChainLSPBinding {
	output := ExecutionEnvelopeProvenanceChainLSPBinding{
		Status: "UNKNOWN", MissingStageIndex: -1, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Chain.NonAuthorizing {
		output.NonAuthorizing = false
		output.Code = "authorization-boundary"
		output.MissingStage = "authorization-boundary"
		return output
	}
	if input.Chain.Status == "ready" {
		projected := ProjectExecutionEnvelopeLSPDiagnostic(ExecutionEnvelopeLSPDiagnosticInput{
			Status:              "clear",
			Code:                "provenance-complete",
			EvidencePrefixDigest: input.EvidencePrefixDigest,
			NonAuthorizing:      true,
		})
		output.Status = projected.Status
		output.Publishable = projected.Publishable
		output.Severity = projected.Severity
		output.Code = projected.Code
		output.EvidencePrefixDigest = projected.EvidencePrefixDigest
		output.ChainEvidenceDigest = input.Chain.EvidenceDigest
		return output
	}
	if input.Chain.MissingStage == "" {
		output.Code = "lsp-diagnostic-evidence"
		return output
	}
	projected := ProjectExecutionEnvelopeLSPDiagnostic(ExecutionEnvelopeLSPDiagnosticInput{
		Status:              "diagnostic",
		Severity:            "error",
		Code:                "provenance-chain",
		MissingStage:        input.Chain.MissingStage,
		MissingStageIndex:   input.MissingStageIndex,
		EvidencePrefixDigest: input.EvidencePrefixDigest,
		NonAuthorizing:      true,
	})
	output.Status = projected.Status
	output.Publishable = projected.Publishable
	output.Severity = projected.Severity
	output.Code = projected.Code
	output.MissingStage = input.Chain.MissingStage
	output.MissingStageIndex = projected.MissingStageIndex
	output.EvidencePrefixDigest = projected.EvidencePrefixDigest
	output.ChainEvidenceDigest = input.Chain.EvidenceDigest
	return output
}
