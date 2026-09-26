package decision

import "testing"

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceToLSP(t *testing.T) {
	input := evidenceFullProvenanceInput(t)
	complete := BindExecutionEnvelopeGoooEvidenceFullProvenance(input)
	ready := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceToLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceLSPInput{
		Binding: complete, NonAuthorizing: true,
	})
	if ready.Status != "ready" || ready.DiagnosticCode != "" || ready.MissingStageIndex != -1 {
		t.Fatalf("ready projection = %#v, want ready without diagnostic", ready)
	}
	if ready.EvidenceDigest != complete.EvidenceDigest || ready.EvidencePrefixDigest == "" ||
		ready.EvidenceBindingDigest == "" || !ready.NonExecuting || !ready.NonAuthorizing {
		t.Fatalf("ready projection = %#v, want complete evidence projection", ready)
	}

	changedInput := input
	changedInput.MetricSource += "\nchanged"
	changed := BindExecutionEnvelopeGoooEvidenceFullProvenance(changedInput)
	changedProjection := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceToLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceLSPInput{
		Binding: changed, NonAuthorizing: true,
	})
	if changedProjection.Status != "ready" || changedProjection.EvidencePrefixDigest == ready.EvidencePrefixDigest {
		t.Fatal("changed metric source must change the ready LSP prefix digest")
	}

	tampered := input
	tampered.EvidenceGeneration.EvidenceDigest = "tampered"
	unknown := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceToLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceLSPInput{
		Binding: BindExecutionEnvelopeGoooEvidenceFullProvenance(tampered),
		NonAuthorizing: true,
	})
	if unknown.Status != "UNKNOWN" || unknown.MissingStage != "evidence-ir-generation-replay" ||
		unknown.MissingStageIndex != 4 || unknown.EvidenceDigest != "" ||
		unknown.EvidencePrefixDigest == "" {
		t.Fatalf("unknown projection = %#v, want evidence replay prefix", unknown)
	}

	unauthorized := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceToLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceLSPInput{
		Binding: complete, NonAuthorizing: false,
	})
	if unauthorized.Status != "UNKNOWN" || unauthorized.NonAuthorizing ||
		unauthorized.DiagnosticCode != "lsp-diagnostic-authorization" {
		t.Fatalf("unauthorized projection = %#v, want non-authorizing UNKNOWN", unauthorized)
	}

	malformed := complete
	malformed.EvidenceBindingDigest = ""
	malformedProjection := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceToLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceLSPInput{
		Binding: malformed, NonAuthorizing: true,
	})
	if malformedProjection.Status != "UNKNOWN" || malformedProjection.MissingStage != "evidence-binding" ||
		malformedProjection.MissingStageIndex != 5 {
		t.Fatalf("malformed projection = %#v, want evidence-binding UNKNOWN", malformedProjection)
	}
}
