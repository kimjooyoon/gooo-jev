package decision

import "strings"

// ExecutionEnvelopeArtifactDigestLSPInput adapts an exact IR or generation
// artifact digest to an editor-facing projection without authorizing use.
type ExecutionEnvelopeArtifactDigestLSPInput struct {
	Artifact       ExecutionEnvelopeArtifactDigest
	NonAuthorizing bool
}

// ExecutionEnvelopeArtifactDigestLSPProjection preserves artifact kind and
// content digest while keeping incomplete source provenance UNKNOWN.
type ExecutionEnvelopeArtifactDigestLSPProjection struct {
	Status         string
	Code           string
	Severity       string
	Message        string
	ArtifactKind   string
	ArtifactDigest string
	EvidenceDigest string
	MissingStage   string
	NonExecuting   bool
	NonAuthorizing bool
}

func digestExecutionEnvelopeArtifactDigestLSPProjection(projection ExecutionEnvelopeArtifactDigestLSPProjection) (string, error) {
	return Digest(struct {
		Status         string
		Code           string
		Severity       string
		Message        string
		ArtifactKind   string
		ArtifactDigest string
		MissingStage   string
	}{
		Status:         projection.Status,
		Code:           projection.Code,
		Severity:        projection.Severity,
		Message:        projection.Message,
		ArtifactKind:   projection.ArtifactKind,
		ArtifactDigest: projection.ArtifactDigest,
		MissingStage:   projection.MissingStage,
	})
}

// ProjectExecutionEnvelopeArtifactDigestLSP publishes only a derived IR or
// generation artifact identity and never executes or authorizes that artifact.
func ProjectExecutionEnvelopeArtifactDigestLSP(input ExecutionEnvelopeArtifactDigestLSPInput) ExecutionEnvelopeArtifactDigestLSPProjection {
	output := ExecutionEnvelopeArtifactDigestLSPProjection{
		Status: "UNKNOWN", Code: "JEV_ARTIFACT_DIGEST_UNKNOWN", Severity: "warning",
		NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Artifact.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		output.Message = "JEV artifact digest is UNKNOWN: missing authorization boundary"
		return finalizeExecutionEnvelopeArtifactDigestLSP(output)
	}
	if !input.Artifact.NonExecuting {
		output.MissingStage = "execution-boundary"
		output.Message = "JEV artifact digest is UNKNOWN: missing execution boundary"
		return finalizeExecutionEnvelopeArtifactDigestLSP(output)
	}
	if input.Artifact.Status != "derived" || strings.TrimSpace(input.Artifact.ArtifactDigest) == "" {
		output.MissingStage = input.Artifact.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "artifact-source"
		}
		output.Message = "JEV artifact digest is UNKNOWN: missing " + output.MissingStage
		return finalizeExecutionEnvelopeArtifactDigestLSP(output)
	}
	if input.Artifact.ArtifactKind != "ir" && input.Artifact.ArtifactKind != "generation" {
		output.MissingStage = "artifact-kind"
		output.Message = "JEV artifact digest is UNKNOWN: missing artifact-kind"
		return finalizeExecutionEnvelopeArtifactDigestLSP(output)
	}
	output.Status = "derived"
	output.Code = "JEV_ARTIFACT_DIGEST_DERIVED"
	output.Severity = "info"
	output.Message = "JEV " + input.Artifact.ArtifactKind + " artifact identity is derived from exact source text"
	output.ArtifactKind = input.Artifact.ArtifactKind
	output.ArtifactDigest = input.Artifact.ArtifactDigest
	return finalizeExecutionEnvelopeArtifactDigestLSP(output)
}

func finalizeExecutionEnvelopeArtifactDigestLSP(output ExecutionEnvelopeArtifactDigestLSPProjection) ExecutionEnvelopeArtifactDigestLSPProjection {
	digest, err := digestExecutionEnvelopeArtifactDigestLSPProjection(output)
	if err != nil {
		output.Status = "UNKNOWN"
		output.Code = "JEV_ARTIFACT_DIGEST_UNKNOWN"
		output.Severity = "warning"
		output.MissingStage = "lsp-evidence-digest"
		output.Message = "JEV artifact digest is UNKNOWN: missing lsp-evidence-digest"
		output.EvidenceDigest = ""
		return output
	}
	output.EvidenceDigest = digest
	return output
}
