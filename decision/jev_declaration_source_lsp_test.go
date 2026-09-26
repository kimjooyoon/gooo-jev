package decision

import "testing"

func TestProjectExecutionEnvelopeDeclarationSourceLSP(t *testing.T) {
	source := ComputeExecutionEnvelopeDeclarationSourceDigest(ExecutionEnvelopeDeclarationSourceDigestInput{
		DeclarationID:  "gooo://gooo-jev/declaration/example",
		ContractID:     "gooo://gooo-jev/contract/example",
		SourceText:     "entity Example id "gooo://example"",
		NonAuthorizing: true,
	})
	got := ProjectExecutionEnvelopeDeclarationSourceLSP(ExecutionEnvelopeDeclarationSourceLSPInput{
		Source:         source,
		NonAuthorizing: true,
	})
	if got.Status != "derived" || got.Code != "JEV_DECLARATION_SOURCE_DERIVED" || got.DeclarationDigest != source.DeclarationDigest || got.EvidenceDigest == "" || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("got %+v", got)
	}
}

func TestProjectExecutionEnvelopeDeclarationSourceLSPPreservesUnknownStage(t *testing.T) {
	source := ExecutionEnvelopeDeclarationSourceDigest{
		Status: "UNKNOWN",
		MissingStage: "declaration-source",
		NonExecuting: true,
		NonAuthorizing: true,
	}
	got := ProjectExecutionEnvelopeDeclarationSourceLSP(ExecutionEnvelopeDeclarationSourceLSPInput{
		Source:         source,
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "declaration-source" || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestProjectExecutionEnvelopeDeclarationSourceLSPRejectsAuthorization(t *testing.T) {
	got := ProjectExecutionEnvelopeDeclarationSourceLSP(ExecutionEnvelopeDeclarationSourceLSPInput{})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}
