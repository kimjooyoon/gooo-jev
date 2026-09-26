package decision

import "strings"

type goooEvidenceFullProvenanceLSPStage struct {
	Name  string
	Value string
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceLSPInput adapts the extended
// evidence-aware provenance binding to editor-facing diagnostics.
type ExecutionEnvelopeGoooEvidenceFullProvenanceLSPInput struct {
	Binding        ExecutionEnvelopeGoooEvidenceFullProvenanceBinding
	NonAuthorizing bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceLSPProjection exposes ordered
// evidence stages without granting execution or authorization capabilities.
type ExecutionEnvelopeGoooEvidenceFullProvenanceLSPProjection struct {
	Status                    string
	DiagnosticCode            string
	MissingStage              string
	MissingStageIndex         int
	DeclarationDigest         string
	IRDigest                  string
	GenerationDigest          string
	BaseBindingDigest         string
	EvidenceDeclarationDigest string
	EvidenceBindingDigest     string
	ReverseObservationDigest  string
	MetricDigest              string
	ProvenanceEvidenceDigest  string
	CompletenessDigest        string
	EvidenceDigest            string
	EvidencePrefixDigest      string
	NonExecuting              bool
	NonAuthorizing            bool
}

func goooEvidenceFullProvenanceLSPStageIndex(stage string) (int, bool) {
	switch stage {
	case "authorization-boundary", "execution-boundary", "declaration-id", "contract-id", "declaration-source", "declaration-source-digest":
		return 0, true
	case "ir-source", "artifact-kind", "artifact-source", "artifact-digest", "legacy-generation", "legacy-ir-generation":
		return 1, true
	case "generation-source":
		return 2, true
	case "evidence-declaration", "evidence-ir-generation-replay":
		return 4, true
	case "evidence-binding":
		return 5, true
	case "reverse-observation-source", "reverse-observation", "reverse-observation-binding":
		return 6, true
	case "metric-source", "metric-source-digest":
		return 7, true
	case "completeness":
		return 9, true
	default:
		return -1, false
	}
}

func digestGoooEvidenceFullProvenanceLSPPrefix(stages []goooEvidenceFullProvenanceLSPStage, missingStage string, end int) (string, error) {
	if end < 0 || end > len(stages) {
		end = 0
	}
	return Digest(struct {
		MissingStage string
		Stages       []goooEvidenceFullProvenanceLSPStage
	}{MissingStage: missingStage, Stages: stages[:end]})
}

func ProjectExecutionEnvelopeGoooEvidenceFullProvenanceToLSP(input ExecutionEnvelopeGoooEvidenceFullProvenanceLSPInput) ExecutionEnvelopeGoooEvidenceFullProvenanceLSPProjection {
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceLSPProjection{
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
	if !input.Binding.NonExecuting {
		output.MissingStage = "execution-boundary"
		output.MissingStageIndex = 0
		output.DiagnosticCode = "lsp-diagnostic-execution"
		return output
	}
	stages := []goooEvidenceFullProvenanceLSPStage{
		{Name: "declaration", Value: input.Binding.DeclarationDigest},
		{Name: "ir", Value: input.Binding.IRDigest},
		{Name: "generation", Value: input.Binding.GenerationDigest},
		{Name: "base_binding", Value: input.Binding.BaseBindingDigest},
		{Name: "evidence_declaration", Value: input.Binding.EvidenceDeclarationDigest},
		{Name: "evidence_binding", Value: input.Binding.EvidenceBindingDigest},
		{Name: "reverse_observation", Value: input.Binding.ReverseObservationDigest},
		{Name: "metric", Value: input.Binding.MetricDigest},
		{Name: "provenance_evidence", Value: input.Binding.ProvenanceEvidenceDigest},
		{Name: "completeness", Value: input.Binding.CompletenessDigest},
	}
	if input.Binding.Status == "complete" && strings.TrimSpace(input.Binding.MissingStage) == "" {
		if err := input.Binding.Validate(); err != nil {
			output.MissingStage = "evidence-binding"
			output.MissingStageIndex = 5
			output.DiagnosticCode = "lsp-diagnostic-evidence-binding"
			return output
		} else {
			for _, stage := range stages {
				if strings.TrimSpace(stage.Value) == "" {
					output.MissingStage = "lsp-source-binding"
					output.DiagnosticCode = "lsp-diagnostic-location"
					return output
				}
			}
			prefix, err := digestGoooEvidenceFullProvenanceLSPPrefix(stages, "", len(stages))
			if err != nil {
				output.MissingStage = "lsp-prefix-digest"
				output.DiagnosticCode = "lsp-diagnostic-prefix"
				return output
			}
			output.Status = "ready"
			output.DiagnosticCode = ""
			output.EvidenceDigest = input.Binding.EvidenceDigest
			output.EvidencePrefixDigest = prefix
			output.DeclarationDigest = input.Binding.DeclarationDigest
			output.IRDigest = input.Binding.IRDigest
			output.GenerationDigest = input.Binding.GenerationDigest
			output.BaseBindingDigest = input.Binding.BaseBindingDigest
			output.EvidenceDeclarationDigest = input.Binding.EvidenceDeclarationDigest
			output.EvidenceBindingDigest = input.Binding.EvidenceBindingDigest
			output.ReverseObservationDigest = input.Binding.ReverseObservationDigest
			output.MetricDigest = input.Binding.MetricDigest
			output.ProvenanceEvidenceDigest = input.Binding.ProvenanceEvidenceDigest
			output.CompletenessDigest = input.Binding.CompletenessDigest
			return output
		}
	}
	output.MissingStage = input.Binding.MissingStage
	if strings.TrimSpace(output.MissingStage) == "" {
		output.MissingStage = "lsp-source-binding"
	}
	index, ok := goooEvidenceFullProvenanceLSPStageIndex(output.MissingStage)
	if !ok {
		output.DiagnosticCode = "lsp-diagnostic-location"
		return output
	}
	output.MissingStageIndex = index
	output.DiagnosticCode = "lsp-diagnostic-" + strings.ReplaceAll(output.MissingStage, "_", "-")
	prefix, err := digestGoooEvidenceFullProvenanceLSPPrefix(stages, output.MissingStage, index)
	if err != nil {
		output.MissingStage = "lsp-prefix-digest"
		output.DiagnosticCode = "lsp-diagnostic-prefix"
		output.MissingStageIndex = -1
		return output
	}
	output.EvidencePrefixDigest = prefix
	return output
}
