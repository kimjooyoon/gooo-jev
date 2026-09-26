package gooo

import "testing"

func evidenceGenerationBindingInputs(t *testing.T) (RevisionEvidenceChain, GenerationReceipt) {
	t.Helper()
	chain, _ := revisionTrendChainForReplacement(t, "lineage")
	document, err := Parse(validContract)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	generation, err := Generate(document)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	return chain, generation
}

func TestBindRevisionEvidenceGenerationBindsCanonicalObservation(t *testing.T) {
	chain, generation := evidenceGenerationBindingInputs(t)
	binding, err := BindRevisionEvidenceGeneration(chain, generation)
	if err != nil {
		t.Fatalf("BindRevisionEvidenceGeneration() error = %v", err)
	}
	if binding.Status != "BOUND" || !binding.StructureMatch || !binding.ReverseObserved {
		t.Fatalf("unexpected revision evidence generation binding: %#v", binding)
	}
	if binding.SourceDigest != chain.SourceDigest ||
		binding.GenerationSourceDigest != generation.SourceDigest ||
		binding.GeneratedIRDigest != generation.GeneratedIRDigest {
		t.Fatalf("binding lost generation links: %#v", binding)
	}
	if err := binding.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestBindRevisionEvidenceGenerationRetainsChainFailure(t *testing.T) {
	chain, generation := evidenceGenerationBindingInputs(t)
	chain.ChainDigest = digestString("tampered")
	binding, err := BindRevisionEvidenceGeneration(chain, generation)
	if err == nil {
		t.Fatal("BindRevisionEvidenceGeneration() error = nil, want chain failure")
	}
	if binding.Status != "UNKNOWN" || binding.MissingStage != "revision-evidence-generation-binding-chain" {
		t.Fatalf("unexpected unknown chain binding: %#v", binding)
	}
}

func TestBindRevisionEvidenceGenerationRetainsReverseObservationFailure(t *testing.T) {
	chain, generation := evidenceGenerationBindingInputs(t)
	generation.StructureMatch = false
	binding, err := BindRevisionEvidenceGeneration(chain, generation)
	if err == nil {
		t.Fatal("BindRevisionEvidenceGeneration() error = nil, want generation failure")
	}
	if binding.Status != "UNKNOWN" || binding.MissingStage != "revision-evidence-generation-binding-generation" {
		t.Fatalf("unexpected unknown generation binding: %#v", binding)
	}
}

func TestBindRevisionEvidenceGenerationRetainsSourceLinkFailure(t *testing.T) {
	chain, generation := evidenceGenerationBindingInputs(t)
	generation.SourceDigest = digestString("other-source")
	binding, err := BindRevisionEvidenceGeneration(chain, generation)
	if err == nil {
		t.Fatal("BindRevisionEvidenceGeneration() error = nil, want source link failure")
	}
	if binding.Status != "UNKNOWN" || binding.MissingStage != "revision-evidence-generation-binding-link" {
		t.Fatalf("unexpected unknown source binding: %#v", binding)
	}
}

func TestBindRevisionEvidenceGenerationIsDeterministic(t *testing.T) {
	chain, generation := evidenceGenerationBindingInputs(t)
	first, err := BindRevisionEvidenceGeneration(chain, generation)
	if err != nil {
		t.Fatalf("first BindRevisionEvidenceGeneration() error = %v", err)
	}
	second, err := BindRevisionEvidenceGeneration(chain, generation)
	if err != nil {
		t.Fatalf("second BindRevisionEvidenceGeneration() error = %v", err)
	}
	if first.BindingDigest != second.BindingDigest {
		t.Fatal("same chain and generation produced different binding digest")
	}
}