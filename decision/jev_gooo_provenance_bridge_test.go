package decision

import "testing"

const provenanceBridgeGoooSource = "package jevbridge\n" +
	"namespace jevbridge\n\n" +
	"entity Example id \"gooo://jevbridge/example\"\n\n" +
	"activity Observe(Example) -> Example\n"

func TestBindExecutionEnvelopeDeclarationIRGenerationFromGooo(t *testing.T) {
	got := BindExecutionEnvelopeDeclarationIRGenerationFromGooo(ExecutionEnvelopeGoooDeclarationProvenanceInput{
		DeclarationID:  "gooo://gooo-jev/declaration/bridge",
		ContractID:     "gooo://gooo-jev/contract/bridge",
		SourceText:     provenanceBridgeGoooSource,
		NonAuthorizing: true,
	})
	if got.Status != "bound" || got.DeclarationDigest == "" || got.IRDigest == "" || got.GenerationDigest == "" || got.BindingDigest == "" || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("got %+v", got)
	}
	source := ComputeExecutionEnvelopeDeclarationSourceDigest(ExecutionEnvelopeDeclarationSourceDigestInput{
		DeclarationID:  "gooo://gooo-jev/declaration/bridge",
		ContractID:     "gooo://gooo-jev/contract/bridge",
		SourceText:     provenanceBridgeGoooSource,
		NonAuthorizing: true,
	})
	if got.DeclarationDigest != source.DeclarationDigest {
		t.Fatalf("declaration digest mismatch: got %q want %q", got.DeclarationDigest, source.DeclarationDigest)
	}
	if err := got.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestBindExecutionEnvelopeDeclarationIRGenerationFromGoooPreservesSyntaxUnknown(t *testing.T) {
	got := BindExecutionEnvelopeDeclarationIRGenerationFromGooo(ExecutionEnvelopeGoooDeclarationProvenanceInput{
		DeclarationID:  "gooo://gooo-jev/declaration/bridge",
		ContractID:     "gooo://gooo-jev/contract/bridge",
		SourceText:     "package broken\nnamespace broken\nproperty Missing string\n",
		NonAuthorizing: true,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "syntax" || got.BindingDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestBindExecutionEnvelopeDeclarationIRGenerationFromGoooRejectsAuthorization(t *testing.T) {
	got := BindExecutionEnvelopeDeclarationIRGenerationFromGooo(ExecutionEnvelopeGoooDeclarationProvenanceInput{
		SourceText:     provenanceBridgeGoooSource,
		NonAuthorizing: false,
	})
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing {
		t.Fatalf("got %+v", got)
	}
}
