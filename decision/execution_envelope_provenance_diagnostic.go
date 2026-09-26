package decision

// ExecutionEnvelopeProvenanceDiagnosticInput carries reverse-observation
// location data for an editor or LSP without turning it into success.
type ExecutionEnvelopeProvenanceDiagnosticInput struct {
	Status              string
	MissingStage        string
	MissingStageIndex   int
	EvidencePrefixDigest string
	NonAuthorizing      bool
}

// ExecutionEnvelopeProvenanceDiagnostic is a structured, fail-closed
// diagnostic that can be projected into an editor or LSP.
type ExecutionEnvelopeProvenanceDiagnostic struct {
	Status              string
	Severity            string
	Code                string
	MissingStage        string
	MissingStageIndex   int
	EvidencePrefixDigest string
	NonAuthorizing      bool
}

// DiagnoseExecutionEnvelopeProvenance preserves the first missing stage and
// its evidence prefix for deterministic editor-facing diagnostics.
func DiagnoseExecutionEnvelopeProvenance(input ExecutionEnvelopeProvenanceDiagnosticInput) ExecutionEnvelopeProvenanceDiagnostic {
	output := ExecutionEnvelopeProvenanceDiagnostic{
		Status: "UNKNOWN", MissingStageIndex: -1, NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.Code = "authorization-boundary"
		return output
	}
	if input.Status == "sealed" {
		if input.MissingStage != "" || input.MissingStageIndex != -1 ||
			input.EvidencePrefixDigest == "" {
			output.Code = "sealed-diagnostic-evidence"
			return output
		}
		output.Status = "clear"
		output.Severity = "info"
		output.Code = "provenance-complete"
		output.EvidencePrefixDigest = input.EvidencePrefixDigest
		return output
	}
	if input.MissingStage == "" || input.MissingStageIndex < 0 ||
		input.EvidencePrefixDigest == "" {
		output.Code = "ledger-diagnostic-evidence"
		return output
	}
	output.MissingStage = input.MissingStage
	output.MissingStageIndex = input.MissingStageIndex
	output.EvidencePrefixDigest = input.EvidencePrefixDigest
	switch input.Status {
	case "UNKNOWN":
		output.Status = "diagnostic"
		output.Severity = "error"
		output.Code = "provenance-incomplete"
	case "review-required":
		output.Status = "diagnostic"
		output.Severity = "warning"
		output.Code = "review-required"
	default:
		output.Status = "UNKNOWN"
		output.Code = "ledger-diagnostic-status"
	}
	return output
}