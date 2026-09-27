package gooo

import "testing"

func TestComputeExecutionEnvelopeDeclarationSourceDigestAndLSP(t *testing.T) {
	source := ComputeExecutionEnvelopeDeclarationSourceDigest(ExecutionEnvelopeDeclarationSourceDigestInput{
		DeclarationID:  "gooo://gooo-jev/declaration/example",
		ContractID:     "gooo://gooo-jev/contract/example",
		SourceText:     "entity Example",
		NonAuthorizing: true,
	})
	projection := ProjectExecutionEnvelopeDeclarationSourceLSP(ExecutionEnvelopeDeclarationSourceLSPInput{
		Source:         source,
		NonAuthorizing: true,
	})
	if projection.Status != GoooDeclarationSourceDerived ||
		projection.Code != goooDeclarationSourceDerivedCode ||
		projection.DeclarationDigest != source.DeclarationDigest {
		t.Fatalf("declaration source evidence was not projected: %+v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestProjectExecutionEnvelopeDeclarationSourceLSPPreservesUnknown(t *testing.T) {
	projection := ProjectExecutionEnvelopeDeclarationSourceLSP(ExecutionEnvelopeDeclarationSourceLSPInput{
		Source: ExecutionEnvelopeDeclarationSourceDigest{
			Status:         GoooDeclarationSourceUnknown,
			MissingStage:   "declaration-source",
			NonExecuting:   true,
			NonAuthorizing: true,
		},
		NonAuthorizing: true,
	})
	if projection.Status != GoooDeclarationSourceUnknown || projection.MissingStage != "declaration-source" {
		t.Fatalf("unknown declaration source was lost: %+v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestProjectExecutionEnvelopeDeclarationSourceLSPFailsClosedOnAuthorization(t *testing.T) {
	projection := ProjectExecutionEnvelopeDeclarationSourceLSP(ExecutionEnvelopeDeclarationSourceLSPInput{})
	if projection.Status != GoooDeclarationSourceError ||
		projection.Code != goooDeclarationSourceBoundaryCode ||
		projection.MissingStage != "authorization-boundary" {
		t.Fatalf("authorization boundary was not preserved: %+v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestProjectExecutionEnvelopeDeclarationSourceLSPRejectsProjectionTampering(t *testing.T) {
	source := ComputeExecutionEnvelopeDeclarationSourceDigest(ExecutionEnvelopeDeclarationSourceDigestInput{
		DeclarationID:  "gooo://gooo-jev/declaration/example",
		ContractID:     "gooo://gooo-jev/contract/example",
		SourceText:     "entity Example",
		NonAuthorizing: true,
	})
	projection := ProjectExecutionEnvelopeDeclarationSourceLSP(ExecutionEnvelopeDeclarationSourceLSPInput{Source: source, NonAuthorizing: true})
	projection.ContractID = "gooo://tampered/contract"
	if err := projection.Validate(); err == nil {
		t.Fatal("tampered declaration source projection was accepted")
	}
}