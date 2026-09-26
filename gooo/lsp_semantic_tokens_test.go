package gooo

import "testing"

func TestSemanticTokensProjectBoundSymbols(t *testing.T) {
	snapshot := Analyze(validContract)
	result := SemanticTokens(validContract)
	if result.Status != "BOUND" || result.MissingStage != "" {
		t.Fatalf("unexpected semantic token result: %#v", result)
	}
	if len(result.Tokens) != len(snapshot.Symbols) || len(result.Tokens) == 0 {
		t.Fatalf("token count = %d, symbol count = %d", len(result.Tokens), len(snapshot.Symbols))
	}
	if result.SourceDigest != snapshot.SourceDigest || result.IRDigest != snapshot.IRDigest {
		t.Fatalf("token provenance mismatch: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestSemanticTokensRetainUnknownSyntaxStage(t *testing.T) {
	source := "package jevdecision\nnamespace jevdecision\nactivity broken\n"
	result := SemanticTokens(source)
	if result.Status != "UNKNOWN" || result.MissingStage != "syntax" {
		t.Fatalf("unexpected unknown tokens: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestSemanticTokensRejectTamperedToken(t *testing.T) {
	result := SemanticTokens(validContract)
	if result.Status != "BOUND" {
		t.Fatalf("unexpected semantic token result: %#v", result)
	}
	result.Tokens[0].Length++
	if err := result.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want token tamper rejection")
	}
}

func TestSemanticTokensAreDeterministic(t *testing.T) {
	first := SemanticTokens(validContract)
	second := SemanticTokens(validContract)
	if first.TokensDigest != second.TokensDigest {
		t.Fatal("same source produced different semantic token digest")
	}
}
