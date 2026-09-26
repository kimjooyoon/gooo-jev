package decision

import "fmt"

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
	ChainBindingDigest  string
	NonAuthorizing      bool
}

// ProvenanceChainMissingStageIndex derives the zero-based stage location
// from the chain itself instead of trusting an editor-supplied index.
func ProvenanceChainMissingStageIndex(stage string) (int, bool) {
	stages := []string{"declaration", "ir", "generation", "reverse_observation", "metric"}
	for index, candidate := range stages {
		if stage == candidate {
			return index, true
		}
	}
	return -1, false
}

type provenanceChainEvidencePrefixStage struct {
	Name  string
	Value string
}

// DeriveProvenanceChainEvidencePrefixDigest computes the canonical evidence
// prefix before the missing stage, or the complete chain when ready.
func DeriveProvenanceChainEvidencePrefixDigest(chain ExecutionEnvelopeProvenanceChainBinding, missingStageIndex int) (string, error) {
	if err := chain.Validate(); err != nil {
		return "", fmt.Errorf("chain integrity: %w", err)
	}
	stages := []provenanceChainEvidencePrefixStage{
		{Name: "declaration", Value: chain.DeclarationDigest},
		{Name: "ir", Value: chain.IRDigest},
		{Name: "generation", Value: chain.GenerationDigest},
		{Name: "reverse_observation", Value: chain.ReverseObservationDigest},
		{Name: "metric", Value: chain.MetricDigest},
	}
	limit := len(stages)
	if chain.Status == "ready" {
		if missingStageIndex != -1 {
			return "", fmt.Errorf("ready chain must use complete prefix index")
		}
	} else {
		derivedIndex, ok := ProvenanceChainMissingStageIndex(chain.MissingStage)
		if !ok || derivedIndex != missingStageIndex {
			return "", fmt.Errorf("missing stage index does not match chain")
		}
		limit = missingStageIndex
	}
	digest, err := Digest(struct {
		MissingStageIndex int
		Stages            []provenanceChainEvidencePrefixStage
	}{
		MissingStageIndex: missingStageIndex,
		Stages:            stages[:limit],
	})
	if err != nil {
		return "", err
	}
	return digest, nil
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
	if err := input.Chain.Validate(); err != nil {
		output.Code = "chain-integrity"
		return output
	}
	missingStageIndex := -1
	if input.Chain.Status != "ready" {
		if input.Chain.MissingStage == "" {
			output.Code = "lsp-diagnostic-evidence"
			return output
		}
		var ok bool
		missingStageIndex, ok = ProvenanceChainMissingStageIndex(input.Chain.MissingStage)
		if !ok {
			output.Code = "lsp-diagnostic-location"
			return output
		}
	}
	prefixDigest, err := DeriveProvenanceChainEvidencePrefixDigest(input.Chain, missingStageIndex)
	if err != nil {
		output.Code = "lsp-diagnostic-prefix"
		return output
	}
	if input.EvidencePrefixDigest != "" && input.EvidencePrefixDigest != prefixDigest {
		output.Code = "lsp-diagnostic-prefix"
		return output
	}
	if input.Chain.Status == "ready" {
		projected := ProjectExecutionEnvelopeLSPDiagnostic(ExecutionEnvelopeLSPDiagnosticInput{
			Status:              "clear",
			Code:                "provenance-complete",
			EvidencePrefixDigest: prefixDigest,
			NonAuthorizing:      true,
		})
		output.Status = projected.Status
		output.Publishable = projected.Publishable
		output.Severity = projected.Severity
		output.Code = projected.Code
		output.EvidencePrefixDigest = projected.EvidencePrefixDigest
		output.ChainEvidenceDigest = input.Chain.EvidenceDigest
		output.ChainBindingDigest = input.Chain.BindingDigest
		return output
	}
	projected := ProjectExecutionEnvelopeLSPDiagnostic(ExecutionEnvelopeLSPDiagnosticInput{
		Status:              "diagnostic",
		Severity:            "error",
		Code:                "provenance-chain",
		MissingStage:        input.Chain.MissingStage,
		MissingStageIndex:   missingStageIndex,
		EvidencePrefixDigest: prefixDigest,
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
	output.ChainBindingDigest = input.Chain.BindingDigest
	return output
}
