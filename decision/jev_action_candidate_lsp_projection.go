package decision

// JEVActionCandidateLSPProjectionInput adapts verified candidate evidence to
// an editor-facing diagnostic without authorizing execution.
type JEVActionCandidateLSPProjectionInput struct {
	Verification         JEVActionCandidateVerification
	MissingStageIndex    int
	EvidencePrefixDigest string
	NonAuthorizing       bool
}

// JEVActionCandidateLSPProjection preserves candidate evidence and publishable
// severity while keeping incomplete provenance UNKNOWN.
type JEVActionCandidateLSPProjection struct {
	Status              string
	Publishable         bool
	Severity            string
	Code                string
	CandidateID         string
	FirstMismatch       string
	MissingStageIndex   int
	EvidencePrefixDigest string
	NonAuthorizing      bool
}

// ProjectJEVActionCandidateLSP rejects incomplete editor evidence and never
// converts a missing provenance stage into a successful diagnostic.
func ProjectJEVActionCandidateLSP(input JEVActionCandidateLSPProjectionInput) JEVActionCandidateLSPProjection {
	output := JEVActionCandidateLSPProjection{
		Status: "UNKNOWN", MissingStageIndex: -1, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Verification.NonAuthorizing {
		output.NonAuthorizing = false
		output.Code = "authorization-boundary"
		return output
	}
	if !input.Verification.NonExecuting {
		output.Code = "execution-boundary"
		return output
	}
	if input.Verification.CandidateID == "" || input.Verification.EvidenceDigest == "" || input.EvidencePrefixDigest == "" {
		output.Code = "lsp-diagnostic-evidence"
		return output
	}
	output.CandidateID = input.Verification.CandidateID
	output.FirstMismatch = input.Verification.FirstMismatch
	output.EvidencePrefixDigest = input.EvidencePrefixDigest
	switch input.Verification.Status {
	case "verified":
		if input.Verification.FirstMismatch != "" {
			output.Status = "UNKNOWN"
			output.Code = "lsp-diagnostic-evidence"
			return output
		}
		output.Status = "clear"
		output.Severity = "info"
		output.Code = "jev-candidate-verified"
		return output
	case "review":
		if input.Verification.MissingStage == "" || input.MissingStageIndex < 0 {
			output.Status = "UNKNOWN"
			output.Code = "lsp-diagnostic-evidence"
			return output
		}
		output.Status = "publishable"
		output.Publishable = true
		output.MissingStageIndex = input.MissingStageIndex
		if input.Verification.FirstMismatch != "" {
			output.Severity = "error"
			output.Code = "jev-reverse-counterexample"
		} else {
			output.Severity = "warning"
			output.Code = "jev-candidate-review"
		}
		return output
	default:
		output.Status = "UNKNOWN"
		output.Code = "lsp-diagnostic-evidence"
		return output
	}
}
