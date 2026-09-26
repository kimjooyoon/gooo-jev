package decision

import "testing"

func TestProjectJEVImprovementOriginResolutionLSP(t *testing.T) {
	observation := originResolutionObservation()
	resolution := ResolveJEVImprovementOrigin(JEVImprovementOriginResolutionInput{
		Observation: observation,
		Selection: originResolutionSelection(observation),
		NonAuthorizing: true,
	})
	got := ProjectJEVImprovementOriginResolutionLSP(JEVImprovementOriginResolutionLSPInput{
		Resolution: resolution,
		NonAuthorizing: true,
	})
	if got.Status != "resolved" || got.Code != "JEV_ORIGIN_RESOLUTION_RESOLVED" || got.OriginDigest != resolution.OriginDigest || got.EvidenceDigest == "" || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("got %+v", got)
	}
}

func TestProjectJEVImprovementOriginResolutionLSPPreservesUnknownStage(t *testing.T) {
	resolution := JEVImprovementOriginResolution{
		Status: jevImprovementOriginUnknown,
		MissingStage: "candidate-selection",
		NonExecuting: true,
		NonAuthorizing: true,
	}
	got := ProjectJEVImprovementOriginResolutionLSP(JEVImprovementOriginResolutionLSPInput{
		Resolution: resolution,
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "candidate-selection" || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestProjectJEVImprovementOriginResolutionLSPRejectsAuthorization(t *testing.T) {
	got := ProjectJEVImprovementOriginResolutionLSP(JEVImprovementOriginResolutionLSPInput{})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}
