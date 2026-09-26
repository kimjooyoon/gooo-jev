package decision

import "strings"

type fullProvenanceLSPStage struct {
	Name  string
	Value string
}

// ExecutionEnvelopeFullProvenanceLSPInput projects a source-bound envelope
// without allowing an editor diagnostic to authorize execution.
type ExecutionEnvelopeFullProvenanceLSPInput struct {
	Binding         ExecutionEnvelopeFullProvenanceSourceBinding
	NonAuthorizing  bool
}

// ExecutionEnvelopeFullProvenanceLSPBinding is the stable LSP-facing state
// with an explicit first unresolved stage and evidence prefix digest.
type ExecutionEnvelopeFullProvenanceLSPBinding struct {
	Status              string
	DiagnosticCode      string
	MissingStage        string
	MissingStageIndex   int
	EvidenceDigest      string
	EvidencePrefixDigest string
	NonExecuting        bool
	NonAuthorizing      bool
}

func fullProvenanceLSPStageIndex(stage string) (int, bool) {
	switch stage {
	case "authorization-boundary", "declaration-id", "contract-id", "declaration-source", "declaration-source-digest":
		return 0, true
	case "ir-source", "artifact-kind", "artifact-source", "artifact-digest":
		return 1, true
	case "generation-source":
		return 2, true
	case "reverse-observation-source", "reverse-observation", "reverse-observation-binding":
		return 3, true
	case "metric-source", "metric-source-digest":
		return 4, true
	default:
		return -1, false
	}
}

func digestFullProvenanceLSPPrefix(stages []fullProvenanceLSPStage, missingStage string, end int) (string, error) {
	if end < 0 || end > len(stages) {
		end = 0
	}
	return Digest(struct {
		MissingStage string
		Stages       []fullProvenanceLSPStage
		}{MissingStage: missingStage, Stages: stages[:end]})
}

// ProjectExecutionEnvelopeFullProvenanceToLSP preserves UNKNOWN and emits
// only the evidence prefix that precedes the first unresolved stage.
func ProjectExecutionEnvelopeFullProvenanceToLSP(input ExecutionEnvelopeFullProvenanceLSPInput) ExecutionEnvelopeFullProvenanceLSPBinding {
	output := ExecutionEnvelopeFullProvenanceLSPBinding{
		Status: "UNKNOWN", MissingStageIndex: -1,
		NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Binding.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		output.MissingStageIndex = 0
		output.DiagnosticCode = "lsp-diagnostic-authorization"
		return output
	}
	stages := []fullProvenanceLSPStage{
		{Name: "declaration", Value: input.Binding.DeclarationDigest},
		{Name: "ir", Value: input.Binding.IRDigest},
		{Name: "generation", Value: input.Binding.GenerationDigest},
		{Name: "reverse_observation", Value: input.Binding.ReverseObservationDigest},
		{Name: "metric", Value: input.Binding.MetricDigest},
	}
	if input.Binding.Status == "complete" && strings.TrimSpace(input.Binding.MissingStage) == "" {
		for _, stage := range stages {
			if strings.TrimSpace(stage.Value) == "" {
				output.MissingStage = "lsp-source-binding"
				output.DiagnosticCode = "lsp-diagnostic-location"
				return output
			}
		}
		prefix, err := digestFullProvenanceLSPPrefix(stages, "", len(stages))
		if err != nil {
			output.MissingStage = "lsp-prefix-digest"
			output.DiagnosticCode = "lsp-diagnostic-prefix"
			return output
		}
		output.Status = "ready"
		output.EvidenceDigest = input.Binding.EvidenceDigest
		output.EvidencePrefixDigest = prefix
		return output
	}
	output.MissingStage = input.Binding.MissingStage
	if strings.TrimSpace(output.MissingStage) == "" {
		output.MissingStage = "lsp-source-binding"
	}
	index, ok := fullProvenanceLSPStageIndex(output.MissingStage)
	if !ok {
		output.DiagnosticCode = "lsp-diagnostic-location"
		return output
	}
	output.MissingStageIndex = index
	output.DiagnosticCode = "lsp-diagnostic-" + strings.ReplaceAll(output.MissingStage, "_", "-")
	prefix, err := digestFullProvenanceLSPPrefix(stages, output.MissingStage, index)
	if err != nil {
		output.MissingStage = "lsp-prefix-digest"
		output.DiagnosticCode = "lsp-diagnostic-prefix"
		output.MissingStageIndex = -1
		return output
	}
	output.EvidencePrefixDigest = prefix
	return output
}
