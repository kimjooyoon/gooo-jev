package decision

import "strings"

// ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPInput projects a
// non-executing provenance binding into editor-facing evidence.
type ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPInput struct {
	Binding       ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackBinding
	NonAuthorizing bool
}

// ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPProjection is a
// deterministic diagnostic projection. It never authorizes or executes a
// revision and retains the first missing stage for UNKNOWN bindings.
type ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPProjection struct {
	Status         string
	Code           string
	Severity       string
	Message        string
	BindingDigest  string
	EvidenceDigest string
	MissingStage   string
	NonExecuting   bool
	NonAuthorizing bool
}

func digestExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPProjection(projection ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPProjection) (string, error) {
	return Digest(struct {
		Status         string
		Code           string
		Severity       string
		Message        string
		BindingDigest  string
		MissingStage   string
	}{
		Status:        projection.Status,
		Code:          projection.Code,
		Severity:      projection.Severity,
		Message:       projection.Message,
		BindingDigest: projection.BindingDigest,
		MissingStage:  projection.MissingStage,
	})
}

// ProjectExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSP preserves a
// bound digest or an UNKNOWN missing-stage diagnostic for LSP consumers.
func ProjectExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSP(input ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPInput) ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPProjection {
	output := ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPProjection{
		Status: "UNKNOWN", Code: "JEV_REVISION_FEEDBACK_UNKNOWN", Severity: "warning",
		NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Binding.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		output.Message = "JEV revision feedback provenance is UNKNOWN: missing authorization boundary"
		return finalizeExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSP(output)
	}
	if !input.Binding.NonExecuting {
		output.MissingStage = "execution-boundary"
		output.Message = "JEV revision feedback provenance is UNKNOWN: missing execution boundary"
		return finalizeExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSP(output)
	}
	if input.Binding.Status != "bound" || strings.TrimSpace(input.Binding.BindingDigest) == "" {
		output.MissingStage = input.Binding.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "cycle-revision-feedback-binding"
		}
		output.Message = "JEV revision feedback provenance is UNKNOWN: missing " + output.MissingStage
		return finalizeExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSP(output)
	}
	output.Status = "bound"
	output.Code = "JEV_REVISION_FEEDBACK_BOUND"
	output.Severity = "info"
	output.BindingDigest = input.Binding.BindingDigest
	output.Message = "JEV improvement cycle, revision candidate, and feedback evidence are bound for review"
	return finalizeExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSP(output)
}

func finalizeExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSP(output ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPProjection) ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPProjection {
	digest, err := digestExecutionEnvelopeJEVImprovementCycleRevisionFeedbackLSPProjection(output)
	if err != nil {
		output.Status = "UNKNOWN"
		output.Code = "JEV_REVISION_FEEDBACK_UNKNOWN"
		output.Severity = "warning"
		output.MissingStage = "lsp-evidence-digest"
		output.Message = "JEV revision feedback provenance is UNKNOWN: missing lsp-evidence-digest"
		output.EvidenceDigest = ""
		return output
	}
	output.EvidenceDigest = digest
	return output
}
