package gooo

import "testing"

func declarationIRGenerationLSPTestBinding() ExecutionEnvelopeDeclarationIRGenerationBinding {
	binding := ExecutionEnvelopeDeclarationIRGenerationBinding{
		Status:            ExecutionEnvelopeDeclarationIRGenerationLSPBound,
		DeclarationID:     "gooo://example/declaration",
		ContractID:        "gooo://example/contract",
		DeclarationDigest: digestString("declaration"),
		IRDigest:          digestString("ir"),
		GenerationDigest:  digestString("generation"),
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	binding.BindingDigest = executionEnvelopeDeclarationIRGenerationBindingDigest(binding)
	return binding
}

func TestProjectExecutionEnvelopeDeclarationIRGenerationLSPBound(t *testing.T) {
	projection := ProjectExecutionEnvelopeDeclarationIRGenerationLSP(
		ExecutionEnvelopeDeclarationIRGenerationLSPInput{
			Binding:        declarationIRGenerationLSPTestBinding(),
			NonAuthorizing: true,
		},
	)
	if projection.Status != ExecutionEnvelopeDeclarationIRGenerationLSPBound ||
		projection.Code != executionEnvelopeDeclarationIRGenerationLSPBoundCode ||
		projection.MissingStage != "" {
		t.Fatalf("bound declaration IR generation projection was not preserved: %+v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestProjectExecutionEnvelopeDeclarationIRGenerationLSPPreservesUnknown(t *testing.T) {
	binding := declarationIRGenerationLSPTestBinding()
	binding.Status = ExecutionEnvelopeDeclarationIRGenerationLSPUnknown
	binding.MissingStage = "generated_artifact"
	binding.BindingDigest = ""
	projection := ProjectExecutionEnvelopeDeclarationIRGenerationLSP(
		ExecutionEnvelopeDeclarationIRGenerationLSPInput{Binding: binding, NonAuthorizing: true},
	)
	if projection.Status != ExecutionEnvelopeDeclarationIRGenerationLSPUnknown ||
		projection.MissingStage != "generated_artifact" ||
		projection.Code != executionEnvelopeDeclarationIRGenerationLSPUnknownCode {
		t.Fatalf("unknown declaration IR generation stage was lost: %+v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestProjectExecutionEnvelopeDeclarationIRGenerationLSPRejectsTampering(t *testing.T) {
	binding := declarationIRGenerationLSPTestBinding()
	binding.IRDigest = "sha256:tampered"
	projection := ProjectExecutionEnvelopeDeclarationIRGenerationLSP(
		ExecutionEnvelopeDeclarationIRGenerationLSPInput{Binding: binding, NonAuthorizing: true},
	)
	if projection.Status != ExecutionEnvelopeDeclarationIRGenerationLSPError ||
		projection.Code != executionEnvelopeDeclarationIRGenerationLSPIntegrityCode ||
		projection.MissingStage != "binding_integrity" {
		t.Fatalf("tampered binding was not rejected: %+v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("error projection Validate() error = %v", err)
	}
}

func TestProjectExecutionEnvelopeDeclarationIRGenerationLSPFailsClosedOnCapabilityBoundary(t *testing.T) {
	binding := declarationIRGenerationLSPTestBinding()
	binding.NonAuthorizing = false
	projection := ProjectExecutionEnvelopeDeclarationIRGenerationLSP(
		ExecutionEnvelopeDeclarationIRGenerationLSPInput{Binding: binding, NonAuthorizing: true},
	)
	if projection.Status != ExecutionEnvelopeDeclarationIRGenerationLSPError ||
		projection.Code != executionEnvelopeDeclarationIRGenerationLSPBoundaryCode ||
		projection.MissingStage != "capability_boundary" {
		t.Fatalf("capability boundary was not preserved: %+v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("boundary projection Validate() error = %v", err)
	}
}

func TestExecutionEnvelopeDeclarationIRGenerationLSPRejectsProjectionTampering(t *testing.T) {
	projection := ProjectExecutionEnvelopeDeclarationIRGenerationLSP(
		ExecutionEnvelopeDeclarationIRGenerationLSPInput{
			Binding:        declarationIRGenerationLSPTestBinding(),
			NonAuthorizing: true,
		},
	)
	projection.ContractID = "gooo://tampered/contract"
	if err := projection.Validate(); err == nil {
		t.Fatal("tampered declaration IR generation projection was accepted")
	}
}