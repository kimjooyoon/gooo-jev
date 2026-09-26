package decision

import "testing"

func evidenceFullProvenanceInput(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceInput {
	t.Helper()
	derived := DeriveGoooEvidenceDeclarationIRGeneration(GoooEvidenceDeclarationIRGenerationInput{
		SourceText:     sampleGoooEvidenceDeclaration,
		NonAuthorizing: true,
	})
	if derived.Status != "ready" {
		t.Fatalf("evidence declaration did not derive: %#v", derived)
	}
	return ExecutionEnvelopeGoooEvidenceFullProvenanceInput{
		DeclarationID:            "decl-evidence",
		ContractID:               "contract-evidence",
		SourceText:               sampleGoooEvidenceDeclaration,
		EvidenceGeneration:       derived,
		ObservedStatus:            "ready",
		ExpectedStatus:            "ready",
		ObservedEvidenceDigest:   "evidence-evidence",
		ExpectedEvidenceDigest:   "evidence-evidence",
		ReverseObservationSource: "reverse_observation: ready evidence-evidence",
		MetricSource:              "metric: extended evidence source",
		NonAuthorizing:            true,
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenance(t *testing.T) {
	binding := BindExecutionEnvelopeGoooEvidenceFullProvenance(evidenceFullProvenanceInput(t))
	if binding.Status != "complete" || binding.EvidenceBindingDigest == "" ||
		binding.EvidenceDeclarationDigest == "" || binding.BaseBindingDigest == "" ||
		binding.ProvenanceEvidenceDigest == "" || binding.EvidenceDigest == "" {
		t.Fatalf("binding = %#v, want complete extended provenance", binding)
	}
	if err := binding.Validate(); err != nil {
		t.Fatalf("binding should validate: %v", err)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRejectsTampering(t *testing.T) {
	input := evidenceFullProvenanceInput(t)
	input.EvidenceGeneration.EvidenceDigest = "tampered"
	binding := BindExecutionEnvelopeGoooEvidenceFullProvenance(input)
	if binding.Status != "UNKNOWN" || binding.MissingStage != "evidence-ir-generation-replay" {
		t.Fatalf("binding = %#v, want evidence replay UNKNOWN", binding)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRequiresEvidence(t *testing.T) {
	input := evidenceFullProvenanceInput(t)
	input.SourceText = "package empty\nnamespace empty\nentity Example id \"gooo://example\"\n"
	input.EvidenceGeneration = DeriveGoooEvidenceDeclarationIRGeneration(GoooEvidenceDeclarationIRGenerationInput{
		SourceText:     input.SourceText,
		NonAuthorizing: true,
	})
	binding := BindExecutionEnvelopeGoooEvidenceFullProvenance(input)
	if binding.Status != "UNKNOWN" || binding.MissingStage != "evidence-declaration" {
		t.Fatalf("binding = %#v, want evidence-declaration UNKNOWN", binding)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRejectsAuthorization(t *testing.T) {
	input := evidenceFullProvenanceInput(t)
	input.NonAuthorizing = false
	binding := BindExecutionEnvelopeGoooEvidenceFullProvenance(input)
	if binding.Status != "UNKNOWN" || binding.MissingStage != "authorization-boundary" ||
		binding.NonAuthorizing {
		t.Fatalf("binding = %#v, want authorization UNKNOWN", binding)
	}
}
