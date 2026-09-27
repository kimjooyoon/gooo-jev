package gooo

import "testing"

func lspSemanticTokenObservedInput() RevisionSelfImprovementCycleJEVSubsetConvergenceObservation {
	return ObserveRevisionSelfImprovementCycleJEVSubsetConvergence(
		RevisionSelfImprovementCycleJEVSubsetConvergenceInput{
			SubsetName:               "bounded-core",
			DeclarationDigest:        digestString("semantic-token-declaration"),
			IRDigest:                 digestString("semantic-token-ir"),
			GenerationDigest:         digestString("semantic-token-generation"),
			ReverseObservationDigest: digestString("semantic-token-reverse"),
			ConvergenceIteration:     1,
			Stable:                   true,
			ReverseObserved:           true,
			SourceLineCount:           12,
			GeneratedLineCount:        8,
		},
	)
}

func TestObserveLSPSemanticTokenProvenanceObserved(t *testing.T) {
	token := ObserveLSPSemanticTokenProvenance(lspSemanticTokenObservedInput())
	if token.Status != "BOUND" ||
		token.ProvenanceStatus != lspSemanticTokenProvenanceObserved ||
		token.TokenModifiers != lspSemanticTokenModifiersObserved ||
		token.MissingStage != "" {
		t.Fatalf("expected observed semantic token provenance, got %#v", token)
	}
	if err := token.Validate(); err != nil {
		t.Fatalf("expected valid observed semantic token provenance: %v", err)
	}
}

func TestObserveLSPSemanticTokenProvenancePreservesUnknown(t *testing.T) {
	input := ObserveRevisionSelfImprovementCycleJEVSubsetConvergence(
		RevisionSelfImprovementCycleJEVSubsetConvergenceInput{
			SubsetName:               "bounded-core",
			DeclarationDigest:        digestString("semantic-token-declaration"),
			IRDigest:                 "",
			GenerationDigest:         digestString("semantic-token-generation"),
			ReverseObservationDigest: digestString("semantic-token-reverse"),
			ConvergenceIteration:     1,
			Stable:                   true,
			ReverseObserved:           true,
			SourceLineCount:           12,
			GeneratedLineCount:        8,
		},
	)
	token := ObserveLSPSemanticTokenProvenance(input)
	if token.Status != "BOUND" ||
		token.ProvenanceStatus != lspSemanticTokenProvenanceUnknown ||
		token.MissingStage != "revision-self-improvement-cycle-jev-subset-convergence-lineage" ||
		token.TokenModifiers != lspSemanticTokenModifiersUnknown {
		t.Fatalf("expected unknown provenance to remain visible, got %#v", token)
	}
	if err := token.Validate(); err != nil {
		t.Fatalf("expected valid unknown semantic token provenance: %v", err)
	}
}

func TestObserveLSPSemanticTokenProvenanceRejectsInvalidInput(t *testing.T) {
	input := lspSemanticTokenObservedInput()
	input.ObservationDigest = "tampered"
	token := ObserveLSPSemanticTokenProvenance(input)
	if token.Status != "UNKNOWN" ||
		token.MissingStage != "lsp-semantic-token-provenance-input" ||
		token.ProvenanceStatus != lspSemanticTokenProvenanceUnknown {
		t.Fatalf("expected invalid input to remain unknown, got %#v", token)
	}
	if err := token.Validate(); err != nil {
		t.Fatalf("expected valid unknown projection: %v", err)
	}
}

func TestLSPSemanticTokenProvenanceRejectsTampering(t *testing.T) {
	token := ObserveLSPSemanticTokenProvenance(lspSemanticTokenObservedInput())
	token.TokenModifiers = lspSemanticTokenModifiersUnknown
	if err := token.Validate(); err == nil {
		t.Fatal("expected semantic token modifier tampering to be rejected")
	}
}
