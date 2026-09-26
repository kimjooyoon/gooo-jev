package decision

import "strings"

// ExecutionEnvelopeGoooDeclarationProvenanceInput connects exact .gooo source
// and identity to the existing declaration-to-generation provenance binding.
type ExecutionEnvelopeGoooDeclarationProvenanceInput struct {
	DeclarationID  string
	ContractID     string
	SourceText     string
	NonAuthorizing bool
}

// BindExecutionEnvelopeDeclarationIRGenerationFromGooo parses the declaration
// first, then binds its source, IR, and generated source digests without
// executing or authorizing the declaration.
func BindExecutionEnvelopeDeclarationIRGenerationFromGooo(input ExecutionEnvelopeGoooDeclarationProvenanceInput) ExecutionEnvelopeDeclarationIRGenerationBinding {
	output := ExecutionEnvelopeDeclarationIRGenerationBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	derived := DeriveGoooDeclarationIRGeneration(GoooDeclarationIRGenerationInput{
		SourceText:     input.SourceText,
		NonAuthorizing: true,
	})
	if derived.Status != "ready" {
		output.MissingStage = derived.MissingStage
		return output
	}
	declaration := ComputeExecutionEnvelopeDeclarationSourceDigest(ExecutionEnvelopeDeclarationSourceDigestInput{
		DeclarationID:  input.DeclarationID,
		ContractID:     input.ContractID,
		SourceText:     input.SourceText,
		NonAuthorizing: true,
	})
	if declaration.Status != "derived" {
		output.MissingStage = declaration.MissingStage
		return output
	}
	if strings.TrimSpace(derived.IRDigest) == "" {
		output.MissingStage = "ir-digest"
		return output
	}
	if strings.TrimSpace(derived.GenerationDigest) == "" {
		output.MissingStage = "generation-digest"
		return output
	}
	return BindExecutionEnvelopeDeclarationIRGeneration(
		input.DeclarationID,
		input.ContractID,
		declaration.DeclarationDigest,
		derived.IRDigest,
		derived.GenerationDigest,
	)
}
