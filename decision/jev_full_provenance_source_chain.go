package decision

// ExecutionEnvelopeFullProvenanceSourceInput supplies every non-executing
// source boundary from declaration through reverse observation and metric.
type ExecutionEnvelopeFullProvenanceSourceInput struct {
	DeclarationID          string
	ContractID             string
	DeclarationSource      string
	IRSource               string
	GenerationSource       string
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

// ExecutionEnvelopeFullProvenanceSourceBinding is the flattened evidence
// envelope for the complete source-to-metric path.
type ExecutionEnvelopeFullProvenanceSourceBinding struct {
	Status                          string
	MissingStage                    string
	DeclarationID                   string
	ContractID                      string
	DeclarationDigest               string
	IRDigest                        string
	GenerationDigest                string
	BindingDigest                   string
	ReverseObservationSourceDigest string
	ReverseObservationDigest        string
	MetricDigest                    string
	EvidenceDigest                  string
	CompletenessDigest              string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

// BindExecutionEnvelopeFullProvenanceFromSource composes the individual
// source bindings while preserving the first unresolved boundary.
func BindExecutionEnvelopeFullProvenanceFromSource(input ExecutionEnvelopeFullProvenanceSourceInput) ExecutionEnvelopeFullProvenanceSourceBinding {
	output := ExecutionEnvelopeFullProvenanceSourceBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	declaration := BindExecutionEnvelopeDeclarationIRGenerationFromSources(ExecutionEnvelopeDeclarationIRGenerationArtifactSourceInput{
		DeclarationID:     input.DeclarationID,
		ContractID:        input.ContractID,
		DeclarationSource: input.DeclarationSource,
		IRSource:          input.IRSource,
		GenerationSource:  input.GenerationSource,
		NonAuthorizing:    input.NonAuthorizing,
	})
	if declaration.Status != "bound" {
		output.MissingStage = declaration.MissingStage
		output.NonAuthorizing = declaration.NonAuthorizing
		return output
	}
	reverse := BindExecutionEnvelopeReverseObservationFromSource(ExecutionEnvelopeReverseObservationSourceInput{
		ObservedStatus:         input.ObservedStatus,
		ExpectedStatus:         input.ExpectedStatus,
		ObservedMissingStage:   input.ObservedMissingStage,
		ExpectedMissingStage:   input.ExpectedMissingStage,
		ObservedEvidenceDigest: input.ObservedEvidenceDigest,
		ExpectedEvidenceDigest: input.ExpectedEvidenceDigest,
		SourceText:             input.ReverseObservationSource,
		NonAuthorizing:         input.NonAuthorizing,
	})
	if reverse.SourceDigest == "" {
		output.MissingStage = reverse.FirstMismatch
		output.NonAuthorizing = reverse.NonAuthorizing
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
		DeclarationIRGeneration: declaration,
		ReverseObservation:      reverseBinding,
		MetricDigest:            "",
		NonAuthorizing:          input.NonAuthorizing,
	})
	metric := MeasureExecutionEnvelopeProvenanceChainMetricFromSource(ExecutionEnvelopeProvenanceChainMetricSourceInput{
		Binding:        chain,
		SourceText:     input.MetricSource,
		NonAuthorizing: input.NonAuthorizing,
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
