package decision

import "testing"

func TestProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateToLSP(t *testing.T) {
	input := evidenceFullProvenanceRevisionCandidateInput(t)
	complete := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidate(input)
	ready := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateToLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateLSPInput{
		Binding: complete, NonAuthorizing: true,
	})
	if ready.Status != "ready" || ready.DiagnosticCode != "" || ready.MissingStageIndex != -1 ||
		ready.CandidateDigest == "" || ready.EvidencePrefixDigest == "" ||
		!ready.NonExecuting || !ready.NonAuthorizing {
		t.Fatalf("ready projection = %#v, want ready candidate projection", ready)
	}

	changedInput := input
	changedInput.RevisionChangeDigest += "\nchanged"
	changed := GenerateExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidate(changedInput)
	changedProjection := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateToLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateLSPInput{
		Binding: changed, NonAuthorizing: true,
	})
	if changedProjection.Status != "ready" || changedProjection.EvidencePrefixDigest == ready.EvidencePrefixDigest {
		t.Fatal("changed revision input must change the LSP prefix digest")
	}

	malformed := complete
	malformed.CandidateDigest = ""
	malformedProjection := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateToLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateLSPInput{
		Binding: malformed, NonAuthorizing: true,
	})
	if malformedProjection.Status != "UNKNOWN" || malformedProjection.MissingStage != "lsp-source-binding" ||
		malformedProjection.DiagnosticCode != "lsp-diagnostic-location" {
		t.Fatalf("malformed projection = %#v, want location diagnostic", malformedProjection)
	}

	unauthorized := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateToLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionCandidateLSPInput{
		Binding: complete, NonAuthorizing: false,
	})
	if unauthorized.Status != "UNKNOWN" || unauthorized.NonAuthorizing ||
		unauthorized.DiagnosticCode != "lsp-diagnostic-authorization" {
		t.Fatalf("unauthorized projection = %#v, want authorization diagnostic", unauthorized)
	}
}
