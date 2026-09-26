package gooo

import (
	"strings"
	"testing"
)

func TestHoverProjectsBoundSymbolProvenance(t *testing.T) {
	snapshot := Analyze(validContract)
	if len(snapshot.Symbols) == 0 {
		t.Fatal("Analyze() returned no symbols")
	}
	symbol := snapshot.Symbols[0]
	hover := Hover(validContract, symbol.Position)
	if hover.Status != "BOUND" || hover.MissingStage != "" {
		t.Fatalf("unexpected hover: %#v", hover)
	}
	if hover.SourceDigest != snapshot.SourceDigest || hover.IRDigest != snapshot.IRDigest {
		t.Fatalf("hover provenance mismatch: %#v", hover)
	}
	if hover.SymbolName != symbol.Name || hover.Kind != symbol.Kind || hover.SymbolDigest != symbol.Digest {
		t.Fatalf("hover symbol mismatch: %#v", hover)
	}
	if !strings.Contains(hover.Contents, symbol.Name) || !strings.Contains(hover.Contents, hover.Declaration) {
		t.Fatalf("hover contents missing declaration: %#v", hover)
	}
	if err := hover.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestHoverRetainsUnknownPositionStage(t *testing.T) {
	hover := Hover(validContract, Position{Line: 999, Column: 1})
	if hover.Status != "UNKNOWN" || hover.MissingStage != "lsp-hover-position" {
		t.Fatalf("unexpected unknown hover: %#v", hover)
	}
	if err := hover.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestHoverRetainsParserUnknownStage(t *testing.T) {
	source := "package jevdecision\nnamespace jevdecision\nactivity broken\n"
	hover := Hover(source, Position{Line: 3, Column: 1})
	if hover.Status != "UNKNOWN" || hover.MissingStage != "syntax" {
		t.Fatalf("unexpected parser hover: %#v", hover)
	}
	if err := hover.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestHoverRejectsTamperedBoundEvidence(t *testing.T) {
	snapshot := Analyze(validContract)
	hover := Hover(validContract, snapshot.Symbols[0].Position)
	if hover.Status != "BOUND" {
		t.Fatalf("unexpected hover: %#v", hover)
	}
	hover.Contents = "tampered"
	if err := hover.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want tamper rejection")
	}
}

func TestHoverIsDeterministic(t *testing.T) {
	snapshot := Analyze(validContract)
	first := Hover(validContract, snapshot.Symbols[0].Position)
	second := Hover(validContract, snapshot.Symbols[0].Position)
	if first.HoverDigest != second.HoverDigest {
		t.Fatal("same source position produced different hover digest")
	}
}
