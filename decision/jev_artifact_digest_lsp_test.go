package decision

import "testing"

func TestProjectExecutionEnvelopeArtifactDigestLSP(t *testing.T) {
	for _, kind := range []string{"ir", "generation"} {
		t.Run(kind, func(t *testing.T) {
			artifact := ComputeExecutionEnvelopeArtifactDigest(ExecutionEnvelopeArtifactDigestInput{
				ArtifactKind: kind,
				ArtifactText: "exact " + kind + " source",
				NonAuthorizing: true,
			})
			got := ProjectExecutionEnvelopeArtifactDigestLSP(ExecutionEnvelopeArtifactDigestLSPInput{
				Artifact: artifact,
				NonAuthorizing: true,
			})
			if got.Status != "derived" || got.Code != "JEV_ARTIFACT_DIGEST_DERIVED" || got.ArtifactKind != kind || got.ArtifactDigest != artifact.ArtifactDigest || got.EvidenceDigest == "" || !got.NonExecuting || !got.NonAuthorizing {
				t.Fatalf("got %+v", got)
			}
		})
	}
}

func TestProjectExecutionEnvelopeArtifactDigestLSPPreservesUnknownStage(t *testing.T) {
	artifact := ExecutionEnvelopeArtifactDigest{
		Status: "UNKNOWN",
		ArtifactKind: "ir",
		MissingStage: "artifact-source",
		NonExecuting: true,
		NonAuthorizing: true,
	}
	got := ProjectExecutionEnvelopeArtifactDigestLSP(ExecutionEnvelopeArtifactDigestLSPInput{
		Artifact: artifact,
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "artifact-source" || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestProjectExecutionEnvelopeArtifactDigestLSPRejectsAuthorization(t *testing.T) {
	got := ProjectExecutionEnvelopeArtifactDigestLSP(ExecutionEnvelopeArtifactDigestLSPInput{})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}
