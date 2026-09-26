package decision

import (
	"fmt"
	"strings"
)

// ExecutionEnvelopeDeclarationIRGenerationBinding binds a .gooo declaration
// identity to the IR and generation artifacts derived from it.
type ExecutionEnvelopeDeclarationIRGenerationBinding struct {
	Status            string
	DeclarationID     string
	ContractID        string
	DeclarationDigest string
	IRDigest          string
	GenerationDigest  string
	BindingDigest     string
	NonExecuting      bool
	NonAuthorizing    bool
	MissingStage      string
}

// BindExecutionEnvelopeDeclarationIRGeneration creates a content-addressed
// declaration-to-IR-to-generation boundary without executing the plan.
func BindExecutionEnvelopeDeclarationIRGeneration(
	declarationID string,
	contractID string,
	declarationDigest string,
	irDigest string,
	generationDigest string,
) ExecutionEnvelopeDeclarationIRGenerationBinding {
	output := ExecutionEnvelopeDeclarationIRGenerationBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	stages := []struct {
		name  string
		value string
	}{
		{name: "declaration-id", value: declarationID},
		{name: "contract-id", value: contractID},
		{name: "declaration", value: declarationDigest},
		{name: "ir", value: irDigest},
		{name: "generation", value: generationDigest},
	}
	for _, stage := range stages {
		if strings.TrimSpace(stage.value) == "" {
			output.MissingStage = stage.name
			return output
		}
	}
	bindingDigest, err := Digest(struct {
		DeclarationID     string
		ContractID        string
		DeclarationDigest string
		IRDigest          string
		GenerationDigest  string
	}{
		DeclarationID:     declarationID,
		ContractID:        contractID,
		DeclarationDigest: declarationDigest,
		IRDigest:          irDigest,
		GenerationDigest:  generationDigest,
	})
	if err != nil {
		output.MissingStage = "binding-evidence"
		return output
	}
	output.Status = "bound"
	output.DeclarationID = declarationID
	output.ContractID = contractID
	output.DeclarationDigest = declarationDigest
	output.IRDigest = irDigest
	output.GenerationDigest = generationDigest
	output.BindingDigest = bindingDigest
	return output
}

func (binding ExecutionEnvelopeDeclarationIRGenerationBinding) Validate() error {
	if binding.Status != "bound" || !binding.NonExecuting || !binding.NonAuthorizing ||
		strings.TrimSpace(binding.DeclarationID) == "" ||
		strings.TrimSpace(binding.ContractID) == "" ||
		strings.TrimSpace(binding.DeclarationDigest) == "" ||
		strings.TrimSpace(binding.IRDigest) == "" ||
		strings.TrimSpace(binding.GenerationDigest) == "" ||
		strings.TrimSpace(binding.BindingDigest) == "" ||
		strings.TrimSpace(binding.MissingStage) != "" {
		return fmt.Errorf("declaration IR generation binding is incomplete")
	}
	expected, err := Digest(struct {
		DeclarationID     string
		ContractID        string
		DeclarationDigest string
		IRDigest          string
		GenerationDigest  string
	}{
		DeclarationID:     binding.DeclarationID,
		ContractID:        binding.ContractID,
		DeclarationDigest: binding.DeclarationDigest,
		IRDigest:          binding.IRDigest,
		GenerationDigest:  binding.GenerationDigest,
	})
	if err != nil {
		return err
	}
	if expected != binding.BindingDigest {
		return fmt.Errorf("declaration IR generation binding digest mismatch")
	}
	return nil
}

// ExecutionEnvelopeProvenanceChainBindingInput connects the declaration/IR/
// generation binding to reverse observation and metric evidence.
type ExecutionEnvelopeProvenanceChainBindingInput struct {
	DeclarationIRGeneration ExecutionEnvelopeDeclarationIRGenerationBinding
	ReverseObservationDigest string
	MetricDigest             string
	NonAuthorizing           bool
}

// ExecutionEnvelopeProvenanceChainBinding records the full evidence chain
// without making an authorization or execution claim.
type ExecutionEnvelopeProvenanceChainBinding struct {
	Status                  string
	MissingStage            string
	DeclarationID           string
	ContractID              string
	DeclarationDigest       string
	IRDigest                string
	GenerationDigest        string
	BindingDigest           string
	ReverseObservationDigest string
	MetricDigest            string
	EvidenceDigest          string
	NonExecuting             bool
	NonAuthorizing           bool
}

// EvaluateExecutionEnvelopeProvenanceChainBinding reuses the existing
// five-stage provenance gate after validating the declaration-to-generation
// binding.
func EvaluateExecutionEnvelopeProvenanceChainBinding(input ExecutionEnvelopeProvenanceChainBindingInput) ExecutionEnvelopeProvenanceChainBinding {
	output := ExecutionEnvelopeProvenanceChainBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if err := input.DeclarationIRGeneration.Validate(); err != nil {
		output.MissingStage = "declaration-ir-generation"
		return output
	}
	binding := input.DeclarationIRGeneration
	gate := EvaluateExecutionEnvelopeProvenanceGate(ExecutionEnvelopeProvenanceGateInput{
		DeclarationDigest:        binding.DeclarationDigest,
		IRDigest:                 binding.IRDigest,
		GenerationDigest:         binding.GenerationDigest,
		ReverseObservationDigest: input.ReverseObservationDigest,
		MetricDigest:             input.MetricDigest,
		NonAuthorizing:           true,
	})
	output.Status = gate.Status
	output.MissingStage = gate.MissingStage
	output.DeclarationID = binding.DeclarationID
	output.ContractID = binding.ContractID
	output.DeclarationDigest = binding.DeclarationDigest
	output.IRDigest = binding.IRDigest
	output.GenerationDigest = binding.GenerationDigest
	output.BindingDigest = binding.BindingDigest
	output.ReverseObservationDigest = input.ReverseObservationDigest
	output.MetricDigest = input.MetricDigest
	output.EvidenceDigest = gate.EvidenceDigest
	return output
}
