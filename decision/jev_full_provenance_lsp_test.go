package decision

import "testing"

func TestProjectExecutionEnvelopeFullProvenanceToLSP(t *testing.T) {
	input := ExecutionEnvelopeFullProvenanceSourceInput{
		DeclarationID:            "decl-1",
		ContractID:               "contract-1",
		DeclarationSource:        "entity Example id gooo://example",
		IRSource:                 "ir: Example declaration",
		GenerationSource:         "generation: Example artifact",
		ObservedStatus:           "ready",
		ExpectedStatus:           "ready",
		ObservedEvidenceDigest:   "evidence-1",
		ExpectedEvidenceDigest:   "evidence-1",
		ReverseObservationSource: "reverse_observation: ready evidence-1",
		MetricSource:             "metric: all five provenance stages observed",
		NonAuthorizing:           true,
	}
	complete := BindExecutionEnvelopeFullProvenanceFromSource(input)
	ready := ProjectExecutionEnvelopeFullProvenanceToLSP(ExecutionEnvelopeFullProvenanceLSPInput{
		Binding: complete, NonAuthorizing: true,
	})
	if ready.Status != "ready" || ready.DiagnosticCode != "" || ready.MissingStageIndex != -1 {
		t.Fatalf("ready projection = %#v, want ready without diagnostic", ready)
	}
	if ready.EvidenceDigest == "" || ready.EvidencePrefixDigest == "" {
		t.Fatal("ready projection must retain evidence and full prefix digests")
	}
	changed := input
	changed.IRSource += "\n"
	changedProjection := ProjectExecutionEnvelopeFullProvenanceToLSP(ExecutionEnvelopeFullProvenanceLSPInput{
		Binding: BindExecutionEnvelopeFullProvenanceFromSource(changed), NonAuthorizing: true,
	})
	if changedProjection.Status != "ready" || changedProjection.EvidencePrefixDigest == ready.EvidencePrefixDigest {
		t.Fatal("changed IR source must change the ready LSP prefix digest")
	}
	missingMetric := input
	missingMetric.MetricSource = ""
	unknown := ProjectExecutionEnvelopeFullProvenanceToLSP(ExecutionEnvelopeFullProvenanceLSPInput{
		Binding: BindExecutionEnvelopeFullProvenanceFromSource(missingMetric), NonAuthorizing: true,
	})
	if unknown.Status != "UNKNOWN" || unknown.MissingStage != "metric-source" || unknown.MissingStageIndex != 4 {
		t.Fatalf("unknown projection = %#v, want metric-source index 4", unknown)
	}
	if unknown.EvidenceDigest != "" || unknown.EvidencePrefixDigest == "" || unknown.DiagnosticCode != "lsp-diagnostic-metric-source" {
		t.Fatal("unknown projection must keep only its prefix evidence and diagnostic")
	}
	unauthorized := ProjectExecutionEnvelopeFullProvenanceToLSP(ExecutionEnvelopeFullProvenanceLSPInput{
		Binding: complete, NonAuthorizing: false,
	})
	if unauthorized.Status != "UNKNOWN" || unauthorized.NonAuthorizing || unauthorized.DiagnosticCode != "lsp-diagnostic-authorization" {
		t.Fatalf("unauthorized projection = %#v, want non-authorizing UNKNOWN", unauthorized)
	}
	malformed := complete
	malformed.IRDigest = ""
	malformedProjection := ProjectExecutionEnvelopeFullProvenanceToLSP(ExecutionEnvelopeFullProvenanceLSPInput{
		Binding: malformed, NonAuthorizing: true,
	})
	if malformedProjection.Status != "UNKNOWN" || malformedProjection.DiagnosticCode != "lsp-diagnostic-location" {
		t.Fatalf("malformed complete projection = %#v, want location diagnostic", malformedProjection)
	}
}
