package gooo

import "testing"

func TestDefinitionBindsIndexedSymbolPosition(t *testing.T) {
	snapshot := Analyze(validContract)
	if snapshot.Status != "BOUND" || len(snapshot.Symbols) == 0 {
		t.Fatalf("unexpected language snapshot: %#v", snapshot)
	}
	symbol := snapshot.Symbols[0]
	result := Definition(validContract, symbol.Position)
	if result.Status != "BOUND" || result.MissingStage != "" {
		t.Fatalf("unexpected definition result: %#v", result)
	}
	if result.SourceDigest != snapshot.SourceDigest || result.IRDigest != snapshot.IRDigest ||
		result.SymbolName != symbol.Name || result.SymbolDigest != symbol.Digest ||
		result.DefinitionPosition != symbol.Position {
		t.Fatalf("definition lost symbol provenance: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestDefinitionRetainsUnknownSymbolStage(t *testing.T) {
	result := Definition(validContract, Position{Line: 1, Column: 1})
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-definition-symbol" {
		t.Fatalf("unexpected unknown definition: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestDefinitionRejectsTamperedDigest(t *testing.T) {
	snapshot := Analyze(validContract)
	if snapshot.Status != "BOUND" || len(snapshot.Symbols) == 0 {
		t.Fatalf("unexpected language snapshot: %#v", snapshot)
	}
	result := Definition(validContract, snapshot.Symbols[0].Position)
	if result.Status != "BOUND" {
		t.Fatalf("unexpected definition result: %#v", result)
	}
	result.DefinitionPosition.Column++
	if err := result.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want tampered definition failure")
	}
}

func TestDefinitionIsDeterministic(t *testing.T) {
	snapshot := Analyze(validContract)
	if snapshot.Status != "BOUND" || len(snapshot.Symbols) == 0 {
		t.Fatalf("unexpected language snapshot: %#v", snapshot)
	}
	first := Definition(validContract, snapshot.Symbols[0].Position)
	second := Definition(validContract, snapshot.Symbols[0].Position)
	if first.DefinitionDigest != second.DefinitionDigest {
		t.Fatal("same source position produced different definition digest")
	}
}
