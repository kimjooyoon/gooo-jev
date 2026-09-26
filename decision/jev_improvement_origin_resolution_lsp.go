package decision

import "strings"

// JEVImprovementOriginResolutionLSPInput adapts verified origin resolution
// evidence to an editor-facing projection without authorizing execution.
type JEVImprovementOriginResolutionLSPInput struct {
	Resolution     JEVImprovementOriginResolution
	NonAuthorizing bool
}

// JEVImprovementOriginResolutionLSPProjection preserves the origin and
// evidence digests while making missing provenance stages explicit.
type JEVImprovementOriginResolutionLSPProjection struct {
	Status                    string
	Code                      string
	Severity                  string
	Message                   string
	OriginDigest              string
	EvidenceDigest            string
	DeclarationDigest         string
	IRDigest                  string
	GenerationDigest          string
	ReverseObservationDigest  string
	MetricDigest              string
	CandidateDigest           string
	CandidateSource           string
	ObservationEvidenceDigest string
	SelectionEvidenceDigest   string
	MissingStage               string
	NonExecuting               bool
	NonAuthorizing             bool
}

func digestJEVImprovementOriginResolutionLSPProjection(projection JEVImprovementOriginResolutionLSPProjection) (string, error) {
	return Digest(struct {
		Status                    string
		Code                      string
		Severity                  string
		Message                   string
		OriginDigest              string
		ResolutionEvidenceDigest  string
		MissingStage              string
	}{
		Status:                   projection.Status,
		Code:                     projection.Code,
		Severity:                 projection.Severity,
		Message:                  projection.Message,
		OriginDigest:              projection.OriginDigest,
		ResolutionEvidenceDigest: projection.EvidenceDigest,
		MissingStage:              projection.MissingStage,
	})
}

// ProjectJEVImprovementOriginResolutionLSP projects a resolved origin or an
// UNKNOWN missing-stage diagnostic and never authorizes or executes a change.
func ProjectJEVImprovementOriginResolutionLSP(input JEVImprovementOriginResolutionLSPInput) JEVImprovementOriginResolutionLSPProjection {
	output := JEVImprovementOriginResolutionLSPProjection{
		Status: "UNKNOWN", Code: "JEV_ORIGIN_RESOLUTION_UNKNOWN", Severity: "warning",
		NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Resolution.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		output.Message = "JEV origin resolution is UNKNOWN: missing authorization boundary"
		return finalizeJEVImprovementOriginResolutionLSP(output)
	}
	if !input.Resolution.NonExecuting {
		output.MissingStage = "execution-boundary"
		output.Message = "JEV origin resolution is UNKNOWN: missing execution boundary"
		return finalizeJEVImprovementOriginResolutionLSP(output)
	}
	if input.Resolution.Status != jevImprovementOriginResolved || strings.TrimSpace(input.Resolution.OriginDigest) == "" || strings.TrimSpace(input.Resolution.EvidenceDigest) == "" {
		output.MissingStage = input.Resolution.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "origin-resolution"
		}
		output.Message = "JEV origin resolution is UNKNOWN: missing " + output.MissingStage
		return finalizeJEVImprovementOriginResolutionLSP(output)
	}
	if err := input.Resolution.Validate(); err != nil {
		output.MissingStage = "origin-resolution-evidence"
		output.Message = "JEV origin resolution is UNKNOWN: invalid origin-resolution-evidence"
		return finalizeJEVImprovementOriginResolutionLSP(output)
	}
	output.Status = "resolved"
	output.Code = "JEV_ORIGIN_RESOLUTION_RESOLVED"
	output.Severity = "info"
	output.Message = "JEV improvement origin is resolved from declaration, IR, generation, reverse observation, metric, and candidate evidence"
	output.OriginDigest = input.Resolution.OriginDigest
	output.EvidenceDigest = input.Resolution.EvidenceDigest
	output.DeclarationDigest = input.Resolution.DeclarationDigest
	output.IRDigest = input.Resolution.IRDigest
	output.GenerationDigest = input.Resolution.GenerationDigest
	output.ReverseObservationDigest = input.Resolution.ReverseObservationDigest
	output.MetricDigest = input.Resolution.MetricDigest
	output.CandidateDigest = input.Resolution.CandidateDigest
	output.CandidateSource = input.Resolution.CandidateSource
	output.ObservationEvidenceDigest = input.Resolution.ObservationEvidenceDigest
	output.SelectionEvidenceDigest = input.Resolution.SelectionEvidenceDigest
	return finalizeJEVImprovementOriginResolutionLSP(output)
}

func finalizeJEVImprovementOriginResolutionLSP(output JEVImprovementOriginResolutionLSPProjection) JEVImprovementOriginResolutionLSPProjection {
	digest, err := digestJEVImprovementOriginResolutionLSPProjection(output)
	if err != nil {
		output.Status = "UNKNOWN"
		output.Code = "JEV_ORIGIN_RESOLUTION_UNKNOWN"
		output.Severity = "warning"
		output.MissingStage = "lsp-evidence-digest"
		output.Message = "JEV origin resolution is UNKNOWN: missing lsp-evidence-digest"
		output.EvidenceDigest = ""
		return output
	}
	output.EvidenceDigest = digest
	return output
}
