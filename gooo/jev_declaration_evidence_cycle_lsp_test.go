package gooo

import "testing"

const (
	declarationEvidenceCycleDeclarationDigest = "sha256:1111111111111111111111111111111111111111111111111111111111111111"
	declarationEvidenceCycleIRDigest          = "sha256:2222222222222222222222222222222222222222222222222222222222222222"
	declarationEvidenceCycleGenerationDigest  = "sha256:3333333333333333333333333333333333333333333333333333333333333333"
	declarationEvidenceCycleBindingDigest     = "sha256:4444444444444444444444444444444444444444444444444444444444444444"
	declarationEvidenceCycleReverseDigest     = "sha256:5555555555555555555555555555555555555555555555555555555555555555"
)

func TestProjectExecutionEnvelopeDeclarationEvidenceCycleLSPBound(t *testing.T) {
	projection := ProjectExecutionEnvelopeDeclarationEvidenceCycleLSP(ExecutionEnvelopeDeclarationEvidenceCycleInput{
		SourceStatus:             GoooDeclarationSourceDerived,
		DeclarationID:            "decl://example",
		ContractID:               "contract://example",
		DeclarationDigest:        declarationEvidenceCycleDeclarationDigest,
		IRStatus:                 ExecutionEnvelopeDeclarationIRGenerationLSPBound,
		IRDigest:                 declarationEvidenceCycleIRDigest,
		GenerationDigest:         declarationEvidenceCycleGenerationDigest,
		BindingDigest:            declarationEvidenceCycleBindingDigest,
		ReverseStatus:            "BOUND",
		ReverseObservationDigest: declarationEvidenceCycleReverseDigest,
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
		SourceStatus:       GoooDeclarationSourceDerived,
		DeclarationID:      "decl://example",
		ContractID:         "contract://example",
		DeclarationDigest:  declarationEvidenceCycleDeclarationDigest,
		IRStatus:           ExecutionEnvelopeDeclarationIRGenerationLSPUnknown,
		IRDigest:           declarationEvidenceCycleIRDigest,
		GenerationDigest:   declarationEvidenceCycleGenerationDigest,
		BindingDigest:      declarationEvidenceCycleBindingDigest,
		ReverseStatus:      "UNKNOWN",
		NonExecuting:       true,
		NonAuthorizing:     true,
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
		DeclarationDigest:         declarationEvidenceCycleDeclarationDigest,
		IRStatus:                 ExecutionEnvelopeDeclarationIRGenerationLSPBound,
		IRDigest:                 declarationEvidenceCycleIRDigest,
		GenerationDigest:          declarationEvidenceCycleGenerationDigest,
		BindingDigest:             declarationEvidenceCycleBindingDigest,
		ReverseStatus:             "DEFERRED",
		ReverseObservationDigest: declarationEvidenceCycleReverseDigest,
		ReverseFirstMismatch:     "producer-deferred",
		NonExecuting:              true,
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
		DeclarationDigest:         declarationEvidenceCycleDeclarationDigest,
		IRStatus:                 ExecutionEnvelopeDeclarationIRGenerationLSPBound,
		IRDigest:                 declarationEvidenceCycleIRDigest,
		GenerationDigest:          declarationEvidenceCycleGenerationDigest,
		BindingDigest:             declarationEvidenceCycleBindingDigest,
		ReverseStatus:             "BOUND",
		ReverseObservationDigest: declarationEvidenceCycleReverseDigest,
		NonExecuting:              true,
		NonAuthorizing:            true,
	})
	projection.EvidenceDigest = "sha256:tampered"
	if err := projection.Validate(); err == nil {
		t.Fatal("expected tampered evidence digest to fail validation")
	}
}

