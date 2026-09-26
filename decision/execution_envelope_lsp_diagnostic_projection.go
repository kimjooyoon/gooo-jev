package decision

// ExecutionEnvelopeLSPDiagnosticInput is the editor-facing shape of a
// provenance diagnostic before publication.
type ExecutionEnvelopeLSPDiagnosticInput struct {
	Status              string
	Severity            string
	Code                string
	MissingStage        string
	MissingStageIndex   int
	EvidencePrefixDigest string
	NonAuthorizing      bool
}

// ExecutionEnvelopeLSPDiagnostic records whether a diagnostic is safe to
// publish without converting provenance into an authorization claim.
type ExecutionEnvelopeLSPDiagnostic struct {
	Status              string
	Publishable         bool
	Severity            string
	Code                string
	MissingStageIndex   int
	EvidencePrefixDigest string
	NonAuthorizing      bool
}

// ProjectExecutionEnvelopeLSPDiagnostic rejects incomplete editor evidence.
func ProjectExecutionEnvelopeLSPDiagnostic(input ExecutionEnvelopeLSPDiagnosticInput) ExecutionEnvelopeLSPDiagnostic {
	output := ExecutionEnvelopeLSPDiagnostic{
		Status: "UNKNOWN", MissingStageIndex: -1, NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.Code = "authorization-boundary"
		return output
	}
	if input.Status == "clear" {
		if input.Code != "provenance-complete" || input.EvidencePrefixDigest == "" {
			output.Code = "lsp-diagnostic-evidence"
			return output
		}
		output.Status = "clear"
		output.Severity = "info"
		output.Code = input.Code
		output.EvidencePrefixDigest = input.EvidencePrefixDigest
		return output
	}
	if input.Status != "diagnostic" ||
		(input.Severity != "error" && input.Severity != "warning") ||
		input.Code == "" || input.MissingStage == "" ||
		input.MissingStageIndex < 0 || input.EvidencePrefixDigest == "" {
		output.Code = "lsp-diagnostic-evidence"
		return output
	}
	output.Status = "publishable"
	output.Publishable = true
	output.Severity = input.Severity
	output.Code = input.Code
	output.MissingStageIndex = input.MissingStageIndex
	output.EvidencePrefixDigest = input.EvidencePrefixDigest
	return output
}