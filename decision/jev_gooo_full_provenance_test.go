package decision

import "testing"

const fullProvenanceGoooSource = "package jevfull\n" +
	"namespace jevfull\n\n" +
	"entity Example id \"gooo://jevfull/example\"\n\n" +
	"activity Observe(Example) -> Example\n"

func fullProvenanceInput() ExecutionEnvelopeGoooFullProvenanceInput {
	irGeneration := DeriveGoooDeclarationIRGeneration(GoooDeclarationIRGenerationInput{
		SourceText:     fullProvenanceGoooSource,
		NonAuthorizing: true,
	})
	return ExecutionEnvelopeGoooFullProvenanceInput{
		DeclarationID:          "gooo://gooo-jev/declaration/full",
		ContractID:             "gooo://gooo-jev/contract/full",
		SourceText:             fullProvenanceGoooSource,
		IRGeneration:           irGeneration,
		ObservedStatus:         "ready",
		ExpectedStatus:         "ready",
		ObservedEvidenceDigest: "reverse-evidence",
		ExpectedEvidenceDigest: "reverse-evidence",
		ReverseObservationSource: "reverse observation source",
		MetricSource:             "metric source",
		NonAuthorizing:         true,
	}
}

func TestBindExecutionEnvelopeFullProvenanceFromGooo(t *testing.T) {
	got := BindExecutionEnvelopeFullProvenanceFromGooo(fullProvenanceInput())
	if got.Status != "complete" || got.MissingStage != "" || got.DeclarationDigest == "" || got.IRDigest == "" || got.GenerationDigest == "" || got.ReverseObservationDigest == "" || got.MetricDigest == "" || got.EvidenceDigest == "" || got.CompletenessDigest == "" || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("got %+v", got)
	}
}

func TestBindExecutionEnvelopeFullProvenanceFromGoooRejectsTamperedIRGeneration(t *testing.T) {
	input := fullProvenanceInput()
	input.IRGeneration.IRDigest = "tampered"
	got := BindExecutionEnvelopeFullProvenanceFromGooo(input)
	if got.Status != "UNKNOWN" || got.MissingStage != "ir-generation-replay" || got.EvidenceDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestBindExecutionEnvelopeFullProvenanceFromGoooPreservesReverseCounterexample(t *testing.T) {
	input := fullProvenanceInput()
	input.ExpectedStatus = "hold"
	got := BindExecutionEnvelopeFullProvenanceFromGooo(input)
	if got.Status != "UNKNOWN" || got.MissingStage != "reverse-observation" || got.EvidenceDigest != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestBindExecutionEnvelopeFullProvenanceFromGoooRejectsAuthorization(t *testing.T) {
	input := fullProvenanceInput()
	input.NonAuthorizing = false
	got := BindExecutionEnvelopeFullProvenanceFromGooo(input)
	if got.Status != "UNKNOWN" || got.MissingStage != "authorization-boundary" || got.NonAuthorizing {
		t.Fatalf("got %+v", got)
	}
}
