package gooo

import "fmt"

const (
	ExecutionEnvelopeDeclarationIRGenerationLSPBound   = "BOUND"
	ExecutionEnvelopeDeclarationIRGenerationLSPUnknown = "UNKNOWN"
	ExecutionEnvelopeDeclarationIRGenerationLSPError   = "ERROR"

	executionEnvelopeDeclarationIRGenerationLSPBoundCode      = "gooo.declaration_ir_generation.bound"
	executionEnvelopeDeclarationIRGenerationLSPUnknownCode    = "gooo.declaration_ir_generation.unknown"
	executionEnvelopeDeclarationIRGenerationLSPIntegrityCode  = "gooo.declaration_ir_generation.integrity"
	executionEnvelopeDeclarationIRGenerationLSPBoundaryCode   = "gooo.declaration_ir_generation.capability_boundary"
)

type ExecutionEnvelopeDeclarationIRGenerationBinding struct {
	Status           string
	DeclarationID    string
	ContractID       string
	DeclarationDigest string
	IRDigest         string
	GenerationDigest string
	BindingDigest    string
	MissingStage     string
	NonExecuting     bool
	NonAuthorizing   bool
}

type ExecutionEnvelopeDeclarationIRGenerationLSPInput struct {
	Binding        ExecutionEnvelopeDeclarationIRGenerationBinding
	NonAuthorizing bool
}

type ExecutionEnvelopeDeclarationIRGenerationLSPProjection struct {
	Status           string
	Code             string
	Severity         string
	Message          string
	DeclarationID    string
	ContractID       string
	DeclarationDigest string
	IRDigest         string
	GenerationDigest string
	BindingDigest    string
	EvidenceDigest   string
	MissingStage     string
	NonExecuting     bool
	NonAuthorizing   bool
}

// ProjectExecutionEnvelopeDeclarationIRGenerationLSP preserves unresolved
// stages and exposes only non-executing, non-authorizing editor evidence.
func ProjectExecutionEnvelopeDeclarationIRGenerationLSP(
	input ExecutionEnvelopeDeclarationIRGenerationLSPInput,
) ExecutionEnvelopeDeclarationIRGenerationLSPProjection {
	binding := input.Binding
	output := ExecutionEnvelopeDeclarationIRGenerationLSPProjection{
		Status:           ExecutionEnvelopeDeclarationIRGenerationLSPUnknown,
		Code:             executionEnvelopeDeclarationIRGenerationLSPUnknownCode,
		Severity:         "error",
		Message:          "declaration IR generation evidence is missing or unresolved",
		DeclarationID:    binding.DeclarationID,
		ContractID:       binding.ContractID,
		DeclarationDigest: binding.DeclarationDigest,
		IRDigest:         binding.IRDigest,
		GenerationDigest: binding.GenerationDigest,
		BindingDigest:    binding.BindingDigest,
		MissingStage:     binding.MissingStage,
		NonExecuting:     true,
		NonAuthorizing:   true,
	}
	if output.MissingStage == "" {
		output.MissingStage = "declaration_ir_generation_binding"
	}
	if !input.NonAuthorizing || !binding.NonExecuting || !binding.NonAuthorizing {
		output.Status = ExecutionEnvelopeDeclarationIRGenerationLSPError
		output.Code = executionEnvelopeDeclarationIRGenerationLSPBoundaryCode
		output.Message = "declaration IR generation crossed a capability boundary"
		output.MissingStage = "capability_boundary"
		output.EvidenceDigest = executionEnvelopeDeclarationIRGenerationLSPDigest(output)
		return output
	}
	if err := binding.Validate(); err != nil {
		output.Status = ExecutionEnvelopeDeclarationIRGenerationLSPError
		output.Code = executionEnvelopeDeclarationIRGenerationLSPIntegrityCode
		output.Message = "declaration IR generation binding failed integrity validation"
		output.MissingStage = "binding_integrity"
		output.EvidenceDigest = executionEnvelopeDeclarationIRGenerationLSPDigest(output)
		return output
	}
	if binding.Status != ExecutionEnvelopeDeclarationIRGenerationLSPBound {
		output.EvidenceDigest = executionEnvelopeDeclarationIRGenerationLSPDigest(output)
		return output
	}

	output.Status = ExecutionEnvelopeDeclarationIRGenerationLSPBound
	output.Code = executionEnvelopeDeclarationIRGenerationLSPBoundCode
	output.Severity = "info"
	output.Message = "declaration IR generation evidence is available for inspection"
	output.MissingStage = ""
	output.EvidenceDigest = executionEnvelopeDeclarationIRGenerationLSPDigest(output)
	return output
}

func (binding ExecutionEnvelopeDeclarationIRGenerationBinding) Validate() error {
	if binding.Status != ExecutionEnvelopeDeclarationIRGenerationLSPUnknown &&
		binding.Status != ExecutionEnvelopeDeclarationIRGenerationLSPBound {
		return fmt.Errorf("invalid declaration IR generation binding status %q", binding.Status)
	}
	if !binding.NonExecuting || !binding.NonAuthorizing {
		return fmt.Errorf("declaration IR generation binding crossed a forbidden boundary")
	}
	if binding.Status == ExecutionEnvelopeDeclarationIRGenerationLSPUnknown {
		if binding.MissingStage == "" {
			return fmt.Errorf("unknown declaration IR generation binding must preserve a missing stage")
		}
		return nil
	}
	if binding.MissingStage != "" || binding.DeclarationID == "" || binding.ContractID == "" {
		return fmt.Errorf("bound declaration IR generation binding is incomplete")
	}
	for name, digest := range map[string]string{
		"declaration": binding.DeclarationDigest,
		"ir":          binding.IRDigest,
		"generation":  binding.GenerationDigest,
		"binding":     binding.BindingDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("invalid %s digest", name)
		}
	}
	if binding.BindingDigest != executionEnvelopeDeclarationIRGenerationBindingDigest(binding) {
		return fmt.Errorf("declaration IR generation binding digest mismatch")
	}
	return nil
}

func executionEnvelopeDeclarationIRGenerationBindingDigest(
	binding ExecutionEnvelopeDeclarationIRGenerationBinding,
) string {
	return digestString(fmt.Sprintf(
		"gooo-declaration-ir-generation-binding|%s|%s|%s|%s|%s|%s|%s|%t|%t",
		binding.Status,
		binding.DeclarationID,
		binding.ContractID,
		binding.DeclarationDigest,
		binding.IRDigest,
		binding.GenerationDigest,
		binding.MissingStage,
		binding.NonExecuting,
		binding.NonAuthorizing,
	))
}

func (projection ExecutionEnvelopeDeclarationIRGenerationLSPProjection) Validate() error {
	switch projection.Status {
	case ExecutionEnvelopeDeclarationIRGenerationLSPBound,
		ExecutionEnvelopeDeclarationIRGenerationLSPUnknown,
		ExecutionEnvelopeDeclarationIRGenerationLSPError:
	default:
		return fmt.Errorf("invalid declaration IR generation LSP status %q", projection.Status)
	}
	if !projection.NonExecuting || !projection.NonAuthorizing {
		return fmt.Errorf("declaration IR generation LSP must remain non-executing and non-authorizing")
	}
	switch projection.Status {
	case ExecutionEnvelopeDeclarationIRGenerationLSPBound:
		if projection.MissingStage != "" ||
			projection.Code != executionEnvelopeDeclarationIRGenerationLSPBoundCode ||
			projection.Severity != "info" ||
			projection.DeclarationID == "" || projection.ContractID == "" ||
			!validDigest(projection.DeclarationDigest) ||
			!validDigest(projection.IRDigest) ||
			!validDigest(projection.GenerationDigest) ||
			!validDigest(projection.BindingDigest) {
			return fmt.Errorf("bound declaration IR generation LSP is incomplete")
		}
	case ExecutionEnvelopeDeclarationIRGenerationLSPUnknown:
		if projection.MissingStage == "" || projection.Code != executionEnvelopeDeclarationIRGenerationLSPUnknownCode {
			return fmt.Errorf("unknown declaration IR generation LSP must preserve its missing stage")
		}
	case ExecutionEnvelopeDeclarationIRGenerationLSPError:
		if projection.MissingStage == "" ||
			(projection.Code != executionEnvelopeDeclarationIRGenerationLSPIntegrityCode &&
				projection.Code != executionEnvelopeDeclarationIRGenerationLSPBoundaryCode) {
			return fmt.Errorf("error declaration IR generation LSP must preserve its failure stage")
		}
	}
	if projection.EvidenceDigest != executionEnvelopeDeclarationIRGenerationLSPDigest(projection) {
		return fmt.Errorf("declaration IR generation LSP evidence digest mismatch")
	}
	return nil
}

func executionEnvelopeDeclarationIRGenerationLSPDigest(
	projection ExecutionEnvelopeDeclarationIRGenerationLSPProjection,
) string {
	return digestString(fmt.Sprintf(
		"gooo-declaration-ir-generation-lsp|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t",
		projection.Status,
		projection.Code,
		projection.Severity,
		projection.Message,
		projection.DeclarationID,
		projection.ContractID,
		projection.DeclarationDigest,
		projection.IRDigest,
		projection.GenerationDigest,
		projection.BindingDigest,
		projection.MissingStage,
		projection.NonExecuting,
		projection.NonAuthorizing,
	))
}