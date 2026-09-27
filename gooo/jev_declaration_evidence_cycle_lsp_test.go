package gooo

import "testing"

func TestProjectExecutionEnvelopeDeclarationEvidenceCycleLSPBound(t *testing.T) {
	projection := ProjectExecutionEnvelopeDeclarationEvidenceCycleLSP(ExecutionEnvelopeDeclarationEvidenceCycleInput{
		SourceStatus:             GoooDeclarationSourceDerived,
		DeclarationID:            "decl://example",
		ContractID:               "contract://example",
		DeclarationDigest:        "sha256:declaration",
		IRStatus:                 ExecutionEnvelopeDeclarationIRGenerationLSPBound,
		IRDigest:                 "sha256:ir",
		GenerationDigest:         "sha256:generation",
		BindingDigest:            "sha256:binding",
		ReverseStatus:            "BOUND",
		ReverseObservationDigest: "sha256:reverse",
		NonExecuting:             true,
		NonAuthorizing:           true,
	})
	if projection.Status != ExecutionEnvelopeDeclarationEvidenceCycleBound {
		t.Fatalf("status = %s, want %s", projection.Status, ExecutionEnvelopeDeclarationEvidenceCycleBound)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("bound projection invalid: %v", err)
	}
}

func TestProjectExecutionEnvelopeDeclarationEvidenceCycleLSPPreservesFirstMissingStage(t *testing.T) {
	projection := ProjectExecutionEnvelopeDeclarationEvidenceCycleLSP(ExecutionEnvelopeDeclarationEvidenceCycleInput{
		SourceStatus:      GoooDeclarationSourceDerived,
		DeclarationID:     "decl://example",
		ContractID:        "contract://example",
		DeclarationDigest: "sha256:declaration",
		IRStatus:          ExecutionEnvelopeDeclarationIRGenerationLSPUnknown,
		IRDigest:          "sha256:ir",
		GenerationDigest:  "sha256:generation",
		BindingDigest:     "sha256:binding",
		ReverseStatus:     "UNKNOWN",
		NonExecuting:      true,
		NonAuthorizing:    true,
	})
	if projection.Status != ExecutionEnvelopeDeclarationEvidenceCycleUnknown ||
		projection.MissingStage != "declaration_ir_generation" {
		t.Fatalf("projection = %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("unknown projection invalid: %v", err)
	}
}

func TestProjectExecutionEnvelopeDeclarationEvidenceCycleLSPPreservesDeferredReverseObservation(t *testing.T) {
	projection := ProjectExecutionEnvelopeDeclarationEvidenceCycleLSP(ExecutionEnvelopeDeclarationEvidenceCycleInput{
		SourceStatus:             GoooDeclarationSourceDerived,
		DeclarationID:            "decl://example",
		ContractID:               "contract://example",
		DeclarationDigest:        "sha256:declaration",
		IRStatus:                 ExecutionEnvelopeDeclarationIRGenerationLSPBound,
		IRDigest:                 "sha256:ir",
		GenerationDigest:         "sha256:generation",
		BindingDigest:            "sha256:binding",
		ReverseStatus:            "DEFERRED",
		ReverseObservationDigest: "sha256:reverse",
		ReverseFirstMismatch:     "producer-deferred",
		NonExecuting:             true,
		NonAuthorizing:           true,
	})
	if projection.Status != ExecutionEnvelopeDeclarationEvidenceCycleDeferred ||
		projection.FirstMismatch != "producer-deferred" {
		t.Fatalf("projection = %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("deferred projection invalid: %v", err)
	}
}

func TestProjectExecutionEnvelopeDeclarationEvidenceCycleLSPRejectsTampering(t *testing.T) {
	projection := ProjectExecutionEnvelopeDeclarationEvidenceCycleLSP(ExecutionEnvelopeDeclarationEvidenceCycleInput{
		SourceStatus:             GoooDeclarationSourceDerived,
		DeclarationID:            "decl://example",
		ContractID:               "contract://example",
		DeclarationDigest:         "sha256:declaration",
		IRStatus:                 ExecutionEnvelopeDeclarationIRGenerationLSPBound,
		IRDigest:                 "sha256:ir",
		GenerationDigest:         "sha256:generation",
		BindingDigest:            "sha256:binding",
		ReverseStatus:             "BOUND",
		ReverseObservationDigest: "sha256:reverse",
		NonExecuting:              true,
		NonAuthorizing:           true,
	})
	projection.EvidenceDigest = "sha256:tampered"
	if err := projection.Validate(); err == nil {
		t.Fatal("expected tampered evidence digest to fail validation")
	}
}

