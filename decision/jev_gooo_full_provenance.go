package decision

import "strings"

// ExecutionEnvelopeGoooFullProvenanceInput connects parser output and source
// evidence to the complete declaration-to-metric provenance chain.
type ExecutionEnvelopeGoooFullProvenanceInput struct {
	DeclarationID          string
	ContractID             string
	SourceText             string
	IRGeneration           GoooDeclarationIRGeneration
	ObservedStatus         string
	ExpectedStatus         string
	ObservedMissingStage   string
	ExpectedMissingStage   string
	ObservedEvidenceDigest string
	ExpectedEvidenceDigest string
	ReverseObservationSource string
	MetricSource           string
	NonAuthorizing         bool
}

// BindExecutionEnvelopeFullProvenanceFromGooo replays parser output before
// composing declaration, IR, generation, reverse observation, and metric
// evidence. It never executes or authorizes the declaration.
func BindExecutionEnvelopeFullProvenanceFromGooo(input ExecutionEnvelopeGoooFullProvenanceInput) ExecutionEnvelopeFullProvenanceSourceBinding {
	output := ExecutionEnvelopeFullProvenanceSourceBinding{
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
	if derived.Status != "ready" || input.IRGeneration.Status != "ready" ||
		derived.IRDigest != input.IRGeneration.IRDigest ||
		derived.GenerationDigest != input.IRGeneration.GenerationDigest ||
		derived.GeneratedSource != input.IRGeneration.GeneratedSource {
		output.MissingStage = "ir-generation-replay"
		return output
	}
	declaration := BindExecutionEnvelopeDeclarationIRGenerationFromGooo(ExecutionEnvelopeGoooDeclarationProvenanceInput{
		DeclarationID:  input.DeclarationID,
		ContractID:     input.ContractID,
		SourceText:     input.SourceText,
		NonAuthorizing: true,
	})
	if declaration.Status != "bound" {
		output.MissingStage = declaration.MissingStage
		return output
	}
	reverse := BindExecutionEnvelopeReverseObservationFromSource(ExecutionEnvelopeReverseObservationSourceInput{
		ObservedStatus:         input.ObservedStatus,
		ExpectedStatus:         input.ExpectedStatus,
		ObservedMissingStage:   input.ObservedMissingStage,
		ExpectedMissingStage:   input.ExpectedMissingStage,
		ObservedEvidenceDigest: input.ObservedEvidenceDigest,
		ExpectedEvidenceDigest: input.ExpectedEvidenceDigest,
		SourceText:              input.ReverseObservationSource,
		NonAuthorizing:         true,
	})
	if strings.TrimSpace(reverse.SourceDigest) == "" {
		output.MissingStage = reverse.FirstMismatch
		return output
	}
	output.ReverseObservationSourceDigest = reverse.SourceDigest
	if reverse.Status != "reproduced" {
		output.MissingStage = "reverse-observation"
		return output
	}
	reverseBinding := ExecutionEnvelopeReverseObservationBinding{
		Status:                 reverse.Status,
		FirstMismatch:          reverse.FirstMismatch,
		ObservedEvidenceDigest: input.ObservedEvidenceDigest,
		ObservationDigest:      reverse.ObservationDigest,
		NonAuthorizing:         reverse.NonAuthorizing,
	}
	chain := EvaluateExecutionEnvelopeProvenanceChainBindingFromReverseObservation(ExecutionEnvelopeProvenanceChainReverseObservationInput{
		DeclarationIRGeneration:  declaration,
		ReverseObservation:      reverseBinding,
		MetricDigest:             "",
		NonAuthorizing:           true,
	})
	metric := MeasureExecutionEnvelopeProvenanceChainMetricFromSource(ExecutionEnvelopeProvenanceChainMetricSourceInput{
		Binding:        chain,
		SourceText:     input.MetricSource,
		NonAuthorizing: true,
	})
	output.DeclarationID = chain.DeclarationID
	output.ContractID = chain.ContractID
	output.DeclarationDigest = chain.DeclarationDigest
	output.IRDigest = chain.IRDigest
	output.GenerationDigest = chain.GenerationDigest
	output.BindingDigest = chain.BindingDigest
	output.ReverseObservationDigest = chain.ReverseObservationDigest
	output.MetricDigest = metric.SourceDigest
	output.EvidenceDigest = metric.EvidenceDigest
	output.CompletenessDigest = metric.CompletenessDigest
	output.MissingStage = metric.MissingStage
	if metric.Status == "complete" {
		output.Status = "complete"
		output.MissingStage = ""
	}
	return output
}
