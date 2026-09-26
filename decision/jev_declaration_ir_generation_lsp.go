package decision

import "strings"

// ExecutionEnvelopeDeclarationIRGenerationLSPInput adapts a validated
// declaration-to-generation binding to editor-facing evidence.
type ExecutionEnvelopeDeclarationIRGenerationLSPInput struct {
	Binding        ExecutionEnvelopeDeclarationIRGenerationBinding
	NonAuthorizing bool
}

// ExecutionEnvelopeDeclarationIRGenerationLSPProjection preserves exact
// declaration, IR, generation, and binding digests without enabling execution.
type ExecutionEnvelopeDeclarationIRGenerationLSPProjection struct {
	Status            string
	Code              string
	Severity          string
	Message           string
	DeclarationID     string
	ContractID        string
	DeclarationDigest string
	IRDigest          string
	GenerationDigest  string
	BindingDigest     string
	EvidenceDigest    string
	MissingStage      string
	NonExecuting      bool
	NonAuthorizing    bool
}

func digestExecutionEnvelopeDeclarationIRGenerationLSPProjection(projection ExecutionEnvelopeDeclarationIRGenerationLSPProjection) (string, error) {
	return Digest(struct {
		Status            string
		Code              string
		Severity          string
		Message           string
		DeclarationID     string
		ContractID        string
		DeclarationDigest string
		IRDigest          string
		GenerationDigest  string
		BindingDigest     string
		MissingStage      string
	}{
		Status:            projection.Status,
		Code:              projection.Code,
		Severity:          projection.Severity,
		Message:           projection.Message,
		DeclarationID:     projection.DeclarationID,
		ContractID:        projection.ContractID,
		DeclarationDigest: projection.DeclarationDigest,
		IRDigest:          projection.IRDigest,
		GenerationDigest:  projection.GenerationDigest,
		BindingDigest:     projection.BindingDigest,
		MissingStage:      projection.MissingStage,
	})
}

// ProjectExecutionEnvelopeDeclarationIRGenerationLSP publishes only a
// validated binding and preserves the first missing stage as UNKNOWN.
func ProjectExecutionEnvelopeDeclarationIRGenerationLSP(input ExecutionEnvelopeDeclarationIRGenerationLSPInput) ExecutionEnvelopeDeclarationIRGenerationLSPProjection {
	output := ExecutionEnvelopeDeclarationIRGenerationLSPProjection{
		Status: "UNKNOWN", Code: "JEV_DECLARATION_IR_GENERATION_UNKNOWN", Severity: "warning",
		NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Binding.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		output.Message = "JEV declaration IR generation is UNKNOWN: missing authorization boundary"
		return finalizeExecutionEnvelopeDeclarationIRGenerationLSP(output)
	}
	if !input.Binding.NonExecuting {
		output.MissingStage = "execution-boundary"
		output.Message = "JEV declaration IR generation is UNKNOWN: missing execution boundary"
		return finalizeExecutionEnvelopeDeclarationIRGenerationLSP(output)
	}
	if input.Binding.Status != "bound" || strings.TrimSpace(input.Binding.BindingDigest) == "" {
		output.MissingStage = input.Binding.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "declaration-ir-generation"
		}
		output.Message = "JEV declaration IR generation is UNKNOWN: missing " + output.MissingStage
		return finalizeExecutionEnvelopeDeclarationIRGenerationLSP(output)
	}
	if err := input.Binding.Validate(); err != nil {
		output.MissingStage = "binding-evidence"
		output.Message = "JEV declaration IR generation is UNKNOWN: invalid binding-evidence"
		return finalizeExecutionEnvelopeDeclarationIRGenerationLSP(output)
	}
	output.Status = "bound"
	output.Code = "JEV_DECLARATION_IR_GENERATION_BOUND"
	output.Severity = "info"
	output.Message = "JEV declaration, IR, and generation evidence are bound for review"
	output.DeclarationID = input.Binding.DeclarationID
	output.ContractID = input.Binding.ContractID
	output.DeclarationDigest = input.Binding.DeclarationDigest
	output.IRDigest = input.Binding.IRDigest
	output.GenerationDigest = input.Binding.GenerationDigest
	output.BindingDigest = input.Binding.BindingDigest
	return finalizeExecutionEnvelopeDeclarationIRGenerationLSP(output)
}

func finalizeExecutionEnvelopeDeclarationIRGenerationLSP(output ExecutionEnvelopeDeclarationIRGenerationLSPProjection) ExecutionEnvelopeDeclarationIRGenerationLSPProjection {
	digest, err := digestExecutionEnvelopeDeclarationIRGenerationLSPProjection(output)
	if err != nil {
		output.Status = "UNKNOWN"
		output.Code = "JEV_DECLARATION_IR_GENERATION_UNKNOWN"
		output.Severity = "warning"
		output.MissingStage = "lsp-evidence-digest"
		output.Message = "JEV declaration IR generation is UNKNOWN: missing lsp-evidence-digest"
		output.EvidenceDigest = ""
		return output
	}
	output.EvidenceDigest = digest
	return output
}
