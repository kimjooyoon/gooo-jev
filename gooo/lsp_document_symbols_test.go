package gooo

import "testing"

func TestDocumentSymbolsBindsAnalyzeProjection(t *testing.T) {
	result := DocumentSymbols(validContract)
	if result.Status != "BOUND" || result.MissingStage != "" {
		t.Fatalf("unexpected document symbols result: %#v", result)
	}
	if result.SourceDigest != digestString(validContract) || result.IRDigest == "" {
		t.Fatalf("document symbols lost source or IR provenance: %#v", result)
	}
	if len(result.Symbols) == 0 || result.SymbolsDigest == "" {
		t.Fatalf("document symbols projection is empty: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestDocumentSymbolsRejectsTamperedProjection(t *testing.T) {
	result := DocumentSymbols(validContract)
	if result.Status != "BOUND" || len(result.Symbols) == 0 {
		t.Fatalf("unexpected document symbols result: %#v", result)
	}
	result.Symbols[0].Name = "tampered"
	if err := result.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want tampered symbol failure")
	}
}

func TestDocumentSymbolsIsDeterministic(t *testing.T) {
	first := DocumentSymbols(validContract)
	second := DocumentSymbols(validContract)
	if first.Status != "BOUND" || second.Status != "BOUND" {
		t.Fatalf("unexpected document symbols status: %#v %#v", first, second)
	}
	if first.SymbolsDigest != second.SymbolsDigest {
		t.Fatal("same source produced different document symbol digest")
	}
}
