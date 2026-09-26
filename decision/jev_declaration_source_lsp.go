package decision

import "strings"

// ExecutionEnvelopeDeclarationSourceLSPInput adapts exact .gooo declaration
// source evidence to an editor-facing projection without authorizing execution.
type ExecutionEnvelopeDeclarationSourceLSPInput struct {
	Source         ExecutionEnvelopeDeclarationSourceDigest
	NonAuthorizing bool
}

// ExecutionEnvelopeDeclarationSourceLSPProjection preserves declaration and
// contract identity while making missing source stages explicit.
type ExecutionEnvelopeDeclarationSourceLSPProjection struct {
	Status            string
	Code              string
	Severity          string
	Message           string
	DeclarationID     string
	ContractID        string
	DeclarationDigest string
	EvidenceDigest    string
	MissingStage      string
	NonExecuting      bool
	NonAuthorizing    bool
}

func digestExecutionEnvelopeDeclarationSourceLSPProjection(projection ExecutionEnvelopeDeclarationSourceLSPProjection) (string, error) {
	return Digest(struct {
		Status            string
		Code              string
		Severity          string
		Message           string
		DeclarationID     string
		ContractID        string
		DeclarationDigest string
		MissingStage      string
	}{
		Status:            projection.Status,
		Code:              projection.Code,
		Severity:          projection.Severity,
		Message:           projection.Message,
		DeclarationID:     projection.DeclarationID,
		ContractID:        projection.ContractID,
		DeclarationDigest: projection.DeclarationDigest,
		MissingStage:      projection.MissingStage,
	})
}

// ProjectExecutionEnvelopeDeclarationSourceLSP publishes a verified source
// digest and preserves UNKNOWN when declaration provenance is incomplete.
func ProjectExecutionEnvelopeDeclarationSourceLSP(input ExecutionEnvelopeDeclarationSourceLSPInput) ExecutionEnvelopeDeclarationSourceLSPProjection {
	output := ExecutionEnvelopeDeclarationSourceLSPProjection{
		Status: "UNKNOWN", Code: "JEV_DECLARATION_SOURCE_UNKNOWN", Severity: "warning",
		NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Source.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		output.Message = "JEV declaration source is UNKNOWN: missing authorization boundary"
		return finalizeExecutionEnvelopeDeclarationSourceLSP(output)
	}
	if !input.Source.NonExecuting {
		output.MissingStage = "execution-boundary"
		output.Message = "JEV declaration source is UNKNOWN: missing execution boundary"
		return finalizeExecutionEnvelopeDeclarationSourceLSP(output)
	}
	if input.Source.Status != "derived" || strings.TrimSpace(input.Source.DeclarationDigest) == "" {
		output.MissingStage = input.Source.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "declaration-source"
		}
		output.Message = "JEV declaration source is UNKNOWN: missing " + output.MissingStage
		return finalizeExecutionEnvelopeDeclarationSourceLSP(output)
	}
	output.Status = "derived"
	output.Code = "JEV_DECLARATION_SOURCE_DERIVED"
	output.Severity = "info"
	output.Message = "JEV declaration source identity is derived from exact source text"
	output.DeclarationID = input.Source.DeclarationID
	output.ContractID = input.Source.ContractID
	output.DeclarationDigest = input.Source.DeclarationDigest
	return finalizeExecutionEnvelopeDeclarationSourceLSP(output)
}

func finalizeExecutionEnvelopeDeclarationSourceLSP(output ExecutionEnvelopeDeclarationSourceLSPProjection) ExecutionEnvelopeDeclarationSourceLSPProjection {
	digest, err := digestExecutionEnvelopeDeclarationSourceLSPProjection(output)
	if err != nil {
		output.Status = "UNKNOWN"
		output.Code = "JEV_DECLARATION_SOURCE_UNKNOWN"
		output.Severity = "warning"
		output.MissingStage = "lsp-evidence-digest"
		output.Message = "JEV declaration source is UNKNOWN: missing lsp-evidence-digest"
		output.EvidenceDigest = ""
		return output
	}
	output.EvidenceDigest = digest
	return output
}
