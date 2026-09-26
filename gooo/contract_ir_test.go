package gooo

import (
	"strings"
	"testing"
)

const validContract = "package jevdecision\n" +
	"namespace jevdecision\n" +
	"\n" +
	"# comments are ignored\n" +
	"entity DecisionSpec id \"gooo://jev/decision/spec\"\n" +
	"entity DecisionReceipt id \"gooo://jev/decision/receipt\"\n" +
	"activity ObserveDecision(DecisionSpec) -> DecisionReceipt\n"

func TestParseProducesBoundNonExecutingIR(t *testing.T) {
	document, err := Parse(validContract)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if document.Status != "BOUND" || document.MissingStage != "" {
		t.Fatalf("unexpected status: %#v", document)
	}
	if document.Package != "jevdecision" || document.Namespace != "jevdecision" {
		t.Fatalf("unexpected identity: %#v", document)
	}
	if len(document.Entities) != 2 || len(document.Activities) != 1 {
		t.Fatalf("unexpected declarations: %#v", document)
	}
	if document.Entities[0].Position.Line != 5 || document.Activities[0].Position.Line != 7 {
		t.Fatalf("unexpected positions: %#v", document)
	}
	if document.SourceDigest == "" || document.IRDigest == "" || document.Entities[0].Digest == "" {
		t.Fatalf("missing provenance digest: %#v", document)
	}
	if !document.NonExecuting || !document.NonAuthorizing {
		t.Fatalf("parser must not authorize execution: %#v", document)
	}
	if err := document.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestParseDigestsAreDeterministicAndSourceBound(t *testing.T) {
	first, err := Parse(validContract)
	if err != nil {
		t.Fatalf("first Parse() error = %v", err)
	}
	second, err := Parse(validContract)
	if err != nil {
		t.Fatalf("second Parse() error = %v", err)
	}
	if first.SourceDigest != second.SourceDigest || first.IRDigest != second.IRDigest {
		t.Fatalf("same source produced different digests")
	}

	changed, err := Parse(validContract + "\n")
	if err != nil {
		t.Fatalf("changed Parse() error = %v", err)
	}
	if first.SourceDigest == changed.SourceDigest || first.IRDigest == changed.IRDigest {
		t.Fatalf("source change did not change provenance")
	}
}

func TestParseRejectsUndefinedActivityInput(t *testing.T) {
	source := strings.Replace(validContract, "DecisionSpec", "MissingSpec", 1)
	document, err := Parse(source)
	if err == nil {
		t.Fatal("Parse() error = nil, want validation error")
	}
	parseErr, ok := err.(ParseError)
	if !ok || parseErr.Stage != "ir-validation" {
		t.Fatalf("unexpected parse error: %#v", err)
	}
	if document.Status != "UNKNOWN" || document.MissingStage != "ir-validation" {
		t.Fatalf("unexpected UNKNOWN state: %#v", document)
	}
}

func TestParseRejectsDuplicateEntity(t *testing.T) {
	source := validContract + "entity DecisionSpec id \"gooo://jev/duplicate\"\n"
	_, err := Parse(source)
	parseErr, ok := err.(ParseError)
	if !ok || parseErr.Stage != "duplicate-entity" {
		t.Fatalf("unexpected parse error: %#v", err)
	}
}

func TestParseRejectsMalformedDeclaration(t *testing.T) {
	document, err := Parse("package jevdecision\nnamespace jevdecision\nactivity broken\n")
	if err == nil {
		t.Fatal("Parse() error = nil, want syntax error")
	}
	parseErr, ok := err.(ParseError)
	if !ok || parseErr.Stage != "syntax" {
		t.Fatalf("unexpected parse error: %#v", err)
	}
	if document.Status != "UNKNOWN" || document.MissingStage != "syntax" {
		t.Fatalf("unexpected UNKNOWN state: %#v", document)
	}
}

func TestValidateRejectsTamperedDeclarationDigest(t *testing.T) {
	document, err := Parse(validContract)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	document.Entities[0].Digest = "tampered"
	if err := document.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want tamper rejection")
	}
}

func TestParseRejectsEmptySource(t *testing.T) {
	document, err := Parse("")
	if err == nil {
		t.Fatal("Parse() error = nil, want empty-source")
	}
	parseErr, ok := err.(ParseError)
	if !ok || parseErr.Stage != "empty-source" {
		t.Fatalf("unexpected parse error: %#v", err)
	}
	if document.SourceDigest == "" || document.MissingStage != "empty-source" {
		t.Fatalf("unexpected empty-source state: %#v", document)
	}
}
