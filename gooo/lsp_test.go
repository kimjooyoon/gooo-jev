package gooo

import "testing"

func TestAnalyzeProjectsBoundDeclarations(t *testing.T) {
	snapshot := Analyze(validContract)
	if snapshot.Status != "BOUND" || snapshot.MissingStage != "" {
		t.Fatalf("unexpected snapshot status: %#v", snapshot)
	}
	if len(snapshot.Symbols) != 3 || len(snapshot.Diagnostics) != 0 {
		t.Fatalf("unexpected snapshot contents: %#v", snapshot)
	}
	if snapshot.Symbols[0].Kind != EntitySymbol || snapshot.Symbols[2].Kind != ActivitySymbol {
		t.Fatalf("unexpected symbol kinds: %#v", snapshot.Symbols)
	}
	if snapshot.SourceDigest == "" || snapshot.IRDigest == "" {
		t.Fatalf("missing snapshot provenance: %#v", snapshot)
	}
	if !snapshot.NonExecuting || !snapshot.NonAuthorizing {
		t.Fatalf("snapshot must not authorize execution")
	}
	if err := snapshot.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestAnalyzeRetainsUnknownSyntaxStage(t *testing.T) {
	source := "package jevdecision\nnamespace jevdecision\nactivity broken\n"
	snapshot := Analyze(source)
	if snapshot.Status != "UNKNOWN" || snapshot.MissingStage != "syntax" {
		t.Fatalf("unexpected unknown snapshot: %#v", snapshot)
	}
	if len(snapshot.Diagnostics) != 1 {
		t.Fatalf("diagnostic count = %d, want 1", len(snapshot.Diagnostics))
	}
	diagnostic := snapshot.Diagnostics[0]
	if diagnostic.Stage != "syntax" || diagnostic.Severity != SeverityError {
		t.Fatalf("unexpected diagnostic: %#v", diagnostic)
	}
	if diagnostic.Position.Line != 3 || diagnostic.SourceDigest != snapshot.SourceDigest {
		t.Fatalf("diagnostic provenance mismatch: %#v", diagnostic)
	}
	if err := snapshot.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestAnalyzeRetainsUnknownIRValidationStage(t *testing.T) {
	source := "package jevdecision\n" +
		"namespace jevdecision\n" +
		"entity DecisionReceipt id \"gooo://jev/decision/receipt\"\n" +
		"activity ObserveDecision(DecisionSpec) -> DecisionReceipt\n"
	snapshot := Analyze(source)
	if snapshot.Status != "UNKNOWN" || snapshot.MissingStage != "ir-validation" {
		t.Fatalf("unexpected unknown snapshot: %#v", snapshot)
	}
	if len(snapshot.Diagnostics) != 1 || snapshot.Diagnostics[0].Stage != "ir-validation" {
		t.Fatalf("unexpected diagnostics: %#v", snapshot.Diagnostics)
	}
}

func TestAnalyzeIsDeterministic(t *testing.T) {
	first := Analyze(validContract)
	second := Analyze(validContract)
	if first.SourceDigest != second.SourceDigest || first.IRDigest != second.IRDigest {
		t.Fatalf("same source produced different snapshot provenance")
	}
	if len(first.Symbols) != len(second.Symbols) {
		t.Fatalf("same source produced different symbol count")
	}
}
