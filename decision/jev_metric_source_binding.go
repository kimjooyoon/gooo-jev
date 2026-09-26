package decision

import "strings"

// ExecutionEnvelopeProvenanceChainMetricSourceInput connects exact metric
// source text to the already validated declaration-to-reverse chain.
type ExecutionEnvelopeProvenanceChainMetricSourceInput struct {
	Binding         ExecutionEnvelopeProvenanceChainBinding
	SourceText      string
	NonAuthorizing  bool
}

// ExecutionEnvelopeProvenanceChainMetricSourceBinding records the metric
// coverage and the source digest that supplied its metric-stage evidence.
type ExecutionEnvelopeProvenanceChainMetricSourceBinding struct {
	Status              string
	ObservedStageCount  int
	ExpectedStageCount  int
	MissingStage        string
	CompletenessDigest  string
	EvidenceDigest      string
	SourceDigest        string
	NonExecuting        bool
	NonAuthorizing      bool
}

// MeasureExecutionEnvelopeProvenanceChainMetricFromSource derives the metric
// stage from exact source text without accepting an opaque metric digest.
func MeasureExecutionEnvelopeProvenanceChainMetricFromSource(input ExecutionEnvelopeProvenanceChainMetricSourceInput) ExecutionEnvelopeProvenanceChainMetricSourceBinding {
	output := ExecutionEnvelopeProvenanceChainMetricSourceBinding{
		Status: "UNKNOWN", ExpectedStageCount: 5,
		NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Binding.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if err := input.Binding.Validate(); err != nil {
		output.MissingStage = "provenance-validation"
		return output
	}
	if strings.TrimSpace(input.SourceText) == "" {
		output.MissingStage = "metric-source"
		return output
	}
	sourceDigest, err := Digest(struct {
		SourceText string
	}{SourceText: input.SourceText})
	if err != nil {
		output.MissingStage = "metric-source-digest"
		return output
	}
	declaration := ExecutionEnvelopeDeclarationIRGenerationBinding{
		Status:            "bound",
		DeclarationID:     input.Binding.DeclarationID,
		ContractID:        input.Binding.ContractID,
		DeclarationDigest: input.Binding.DeclarationDigest,
		IRDigest:          input.Binding.IRDigest,
		GenerationDigest:  input.Binding.GenerationDigest,
		BindingDigest:     input.Binding.BindingDigest,
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	chain := EvaluateExecutionEnvelopeProvenanceChainBinding(ExecutionEnvelopeProvenanceChainBindingInput{
		DeclarationIRGeneration:  declaration,
		ReverseObservationDigest: input.Binding.ReverseObservationDigest,
		MetricDigest:             sourceDigest,
		NonAuthorizing:           true,
	})
	output.SourceDigest = sourceDigest
	if chain.Status != "ready" {
		output.MissingStage = chain.MissingStage
		return output
	}
	metric := MeasureExecutionEnvelopeProvenanceChainMetric(chain)
	output.Status = metric.Status
	output.ObservedStageCount = metric.ObservedStageCount
	output.ExpectedStageCount = metric.ExpectedStageCount
	output.MissingStage = metric.MissingStage
	output.CompletenessDigest = metric.CompletenessDigest
	output.EvidenceDigest = metric.EvidenceDigest
	return output
}
