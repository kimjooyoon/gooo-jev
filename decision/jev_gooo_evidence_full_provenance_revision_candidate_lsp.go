package decision

import "strings"

type goooEvidenceRevisionCandidateLSPStage struct {
	Name  string
	Value string
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateLSPInput
// adapts a review-only revision candidate binding to LSP evidence.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateLSPInput struct {
	Binding        ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateBinding
	NonAuthorizing bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateLSPProjection
// exposes ordered candidate evidence without enabling application.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateLSPProjection struct {
	Status                    string
	DiagnosticCode            string
	MissingStage              string
	MissingStageIndex         int
	EvidenceBindingDigest     string
	ProvenanceEvidenceDigest  string
	LedgerEvidenceDigest      string
	LedgerEntryDigest         string
	CandidateDigest           string
	CandidateEvidenceDigest   string
	BoundRevisionChangeDigest string
	BindingDigest             string
	EvidencePrefixDigest      string
	NonExecuting              bool
	NonAuthorizing            bool
}

func goooEvidenceRevisionCandidateLSPStageIndex(stage string) (int, bool) {
	switch stage {
	case "authorization-boundary", "execution-boundary", "evidence-full-provenance", "evidence-full-provenance-validation":
		return 0, true
	case "replay-feedback-ledger", "feedback-unknown", "feedback-direction-binding":
		return 1, true
	case "revision-source", "revision-change", "revision-change-digest":
		return 2, true
	case "revision-candidate":
		return 3, true
	case "revision-candidate-binding-digest":
		return 4, true
	default:
		return -1, false
	}
}

func digestGoooEvidenceRevisionCandidateLSPPrefix(stages []goooEvidenceRevisionCandidateLSPStage, missingStage string, end int) (string, error) {
	if end < 0 || end > len(stages) {
		end = 0
	}
	return Digest(struct {
		MissingStage string
		Stages       []goooEvidenceRevisionCandidateLSPStage
	}{MissingStage: missingStage, Stages: stages[:end]})
}

func ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateToLSP(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateLSPInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateLSPProjection {
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateLSPProjection{
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
	stages := []goooEvidenceRevisionCandidateLSPStage{
		{Name: "evidence-binding", Value: input.Binding.EvidenceBindingDigest},
		{Name: "provenance-evidence", Value: input.Binding.ProvenanceEvidenceDigest},
		{Name: "ledger-evidence", Value: input.Binding.LedgerEvidenceDigest},
		{Name: "ledger-entry", Value: input.Binding.LedgerEntryDigest},
		{Name: "candidate", Value: input.Binding.CandidateDigest},
		{Name: "candidate-evidence", Value: input.Binding.CandidateEvidenceDigest},
		{Name: "bound-revision-change", Value: input.Binding.BoundRevisionChangeDigest},
		{Name: "binding", Value: input.Binding.BindingDigest},
	}
	if input.Binding.Status == "bound" && strings.TrimSpace(input.Binding.MissingStage) == "" {
		for _, stage := range stages {
			if strings.TrimSpace(stage.Value) == "" {
				output.MissingStage = "lsp-source-binding"
				output.DiagnosticCode = "lsp-diagnostic-location"
				return output
			}
		}
		prefix, err := digestGoooEvidenceRevisionCandidateLSPPrefix(stages, "", len(stages))
		if err != nil {
			output.MissingStage = "lsp-prefix-digest"
			output.DiagnosticCode = "lsp-diagnostic-prefix"
			return output
		}
		output.Status = "ready"
		output.EvidencePrefixDigest = prefix
		output.EvidenceBindingDigest = input.Binding.EvidenceBindingDigest
		output.ProvenanceEvidenceDigest = input.Binding.ProvenanceEvidenceDigest
		output.LedgerEvidenceDigest = input.Binding.LedgerEvidenceDigest
		output.LedgerEntryDigest = input.Binding.LedgerEntryDigest
		output.CandidateDigest = input.Binding.CandidateDigest
		output.CandidateEvidenceDigest = input.Binding.CandidateEvidenceDigest
		output.BoundRevisionChangeDigest = input.Binding.BoundRevisionChangeDigest
		output.BindingDigest = input.Binding.BindingDigest
		return output
	}
	output.MissingStage = input.Binding.MissingStage
	if strings.TrimSpace(output.MissingStage) == "" {
		output.MissingStage = "lsp-source-binding"
	}
	index, ok := goooEvidenceRevisionCandidateLSPStageIndex(output.MissingStage)
	if !ok {
		output.DiagnosticCode = "lsp-diagnostic-location"
		return output
	}
	output.MissingStageIndex = index
	output.DiagnosticCode = "lsp-diagnostic-" + strings.ReplaceAll(output.MissingStage, "_", "-")
	prefix, err := digestGoooEvidenceRevisionCandidateLSPPrefix(stages, output.MissingStage, index)
	if err != nil {
		output.MissingStage = "lsp-prefix-digest"
		output.DiagnosticCode = "lsp-diagnostic-prefix"
		output.MissingStageIndex = -1
		return output
	}
	output.EvidencePrefixDigest = prefix
	return output
}
