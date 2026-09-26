package decision

import "testing"

func TestProjectExecutionEnvelopeDeclarationIRGenerationLSP(t *testing.T) {
	binding := BindExecutionEnvelopeDeclarationIRGeneration(
		"gooo://gooo-jev/declaration/example",
		"gooo://gooo-jev/contract/example",
		"declaration-digest",
		"ir-digest",
		"generation-digest",
	)
	got := ProjectExecutionEnvelopeDeclarationIRGenerationLSP(ExecutionEnvelopeDeclarationIRGenerationLSPInput{
		Binding:        binding,
		NonAuthorizing: true,
	})
	if got.Status != "bound" || got.Code != "JEV_DECLARATION_IR_GENERATION_BOUND" || got.BindingDigest != binding.BindingDigest || got.EvidenceDigest == "" || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("got %+v", got)
	}
}

func TestProjectExecutionEnvelopeDeclarationIRGenerationLSPPreservesUnknownStage(t *testing.T) {
	binding := ExecutionEnvelopeDeclarationIRGenerationBinding{
		Status: "UNKNOWN",
		MissingStage: "generation",
		NonExecuting: true,
		NonAuthorizing: true,
	}
	got := ProjectExecutionEnvelopeDeclarationIRGenerationLSP(ExecutionEnvelopeDeclarationIRGenerationLSPInput{
		Binding:        binding,
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "generation" || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}

func TestProjectExecutionEnvelopeDeclarationIRGenerationLSPRejectsAuthorization(t *testing.T) {
	got := ProjectExecutionEnvelopeDeclarationIRGenerationLSP(ExecutionEnvelopeDeclarationIRGenerationLSPInput{})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing || got.EvidenceDigest == "" {
		t.Fatalf("got %+v", got)
	}
}
