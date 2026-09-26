package decision

import (
	"fmt"
	"strings"
)

// DecisionConfidenceChangePlanProvenanceBindingInput combines a non-executing
// improvement plan with its declaration-to-generation provenance boundary.
type DecisionConfidenceChangePlanProvenanceBindingInput struct {
	ChangePlan              DecisionConfidenceChangePlan
	DeclarationIRGeneration ExecutionEnvelopeDeclarationIRGenerationBinding
	NonAuthorizing          bool
}

// DecisionConfidenceChangePlanProvenanceBinding records the plan's origin
// without applying its write set or granting authorization.
type DecisionConfidenceChangePlanProvenanceBinding struct {
	Status                 string
	MissingStage           string
	ProposalDigest         string
	ChangePlanDigest       string
	SourceReference        string
	WriteSetDigest         string
	DeclarationID          string
	ContractID             string
	DeclarationDigest      string
	IRDigest               string
	GenerationDigest       string
	DeclarationBindingDigest string
	EvidenceDigest         string
	NonExecuting           bool
	NonAuthorizing         bool
}

// BindDecisionConfidenceChangePlanProvenance requires the plan and its
// declaration-to-generation binding to agree before recording a bound plan.
func BindDecisionConfidenceChangePlanProvenance(input DecisionConfidenceChangePlanProvenanceBindingInput) DecisionConfidenceChangePlanProvenanceBinding {
	output := DecisionConfidenceChangePlanProvenanceBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if err := input.ChangePlan.Validate(); err != nil {
		output.MissingStage = "change-plan"
		return output
	}
	if err := input.DeclarationIRGeneration.Validate(); err != nil {
		output.MissingStage = "declaration-ir-generation"
		return output
	}
	if input.ChangePlan.SourceReference != input.DeclarationIRGeneration.DeclarationID {
		output.Status = "review"
		output.MissingStage = "source-reference"
		return output
	}
	evidenceDigest, err := Digest(struct {
		ProposalDigest           string
		ChangePlanDigest         string
		SourceReference          string
		WriteSetDigest           string
		DeclarationBindingDigest string
		DeclarationDigest        string
		IRDigest                 string
		GenerationDigest         string
	}{
		ProposalDigest:           input.ChangePlan.ProposalDigest,
		ChangePlanDigest:         input.ChangePlan.ChangePlanDigest,
		SourceReference:          input.ChangePlan.SourceReference,
		WriteSetDigest:           input.ChangePlan.WriteSetDigest,
		DeclarationBindingDigest: input.DeclarationIRGeneration.BindingDigest,
		DeclarationDigest:        input.DeclarationIRGeneration.DeclarationDigest,
		IRDigest:                 input.DeclarationIRGeneration.IRDigest,
		GenerationDigest:         input.DeclarationIRGeneration.GenerationDigest,
	})
	if err != nil {
		output.MissingStage = "plan-provenance-evidence"
		return output
	}
	output.Status = "bound"
	output.ProposalDigest = input.ChangePlan.ProposalDigest
	output.ChangePlanDigest = input.ChangePlan.ChangePlanDigest
	output.SourceReference = input.ChangePlan.SourceReference
	output.WriteSetDigest = input.ChangePlan.WriteSetDigest
	output.DeclarationID = input.DeclarationIRGeneration.DeclarationID
	output.ContractID = input.DeclarationIRGeneration.ContractID
	output.DeclarationDigest = input.DeclarationIRGeneration.DeclarationDigest
	output.IRDigest = input.DeclarationIRGeneration.IRDigest
	output.GenerationDigest = input.DeclarationIRGeneration.GenerationDigest
	output.DeclarationBindingDigest = input.DeclarationIRGeneration.BindingDigest
	output.EvidenceDigest = evidenceDigest
	return output
}

func (binding DecisionConfidenceChangePlanProvenanceBinding) Validate() error {
	if binding.Status != "bound" || !binding.NonExecuting || !binding.NonAuthorizing ||
		strings.TrimSpace(binding.ProposalDigest) == "" ||
		strings.TrimSpace(binding.ChangePlanDigest) == "" ||
		strings.TrimSpace(binding.SourceReference) == "" ||
		strings.TrimSpace(binding.WriteSetDigest) == "" ||
		strings.TrimSpace(binding.DeclarationID) == "" ||
		strings.TrimSpace(binding.DeclarationBindingDigest) == "" ||
		strings.TrimSpace(binding.EvidenceDigest) == "" ||
		strings.TrimSpace(binding.MissingStage) != "" {
		return fmt.Errorf("change plan provenance binding is incomplete")
	}
	return nil
}
