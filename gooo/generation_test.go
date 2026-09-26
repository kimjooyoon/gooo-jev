package gooo

import (
	"strings"
	"testing"
)

func TestGenerateRecordsCanonicalAndReverseEvidence(t *testing.T) {
	document, err := Parse(validContract)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	receipt, err := Generate(document)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if receipt.Status != "BOUND" || receipt.MissingStage != "" {
		t.Fatalf("unexpected receipt status: %#v", receipt)
	}
	if !receipt.StructureMatch {
		t.Fatalf("generated structure did not match: %#v", receipt)
	}
	if receipt.ExactSourceMatch {
		t.Fatalf("comment-bearing source should not be reported as exact canonical source")
	}
	if receipt.SourceDigest == "" || receipt.GeneratedSourceDigest == "" || receipt.GeneratedIRDigest == "" || receipt.StructureDigest == "" {
		t.Fatalf("missing generation evidence: %#v", receipt)
	}
	if !receipt.NonExecuting || !receipt.NonAuthorizing {
		t.Fatalf("generation must remain non-executing and non-authorizing")
	}
	if !strings.Contains(receipt.GeneratedSource, "activity ObserveDecision(DecisionSpec) -> DecisionReceipt") {
		t.Fatalf("canonical source is missing activity: %q", receipt.GeneratedSource)
	}
	observed, err := Parse(receipt.GeneratedSource)
	if err != nil {
		t.Fatalf("generated source Parse() error = %v", err)
	}
	if observed.IRDigest != receipt.GeneratedIRDigest {
		t.Fatalf("receipt IR digest = %q, parsed digest = %q", receipt.GeneratedIRDigest, observed.IRDigest)
	}
}

func TestGenerateCanProveExactCanonicalSourceMatch(t *testing.T) {
	source := "package jevdecision\n" +
		"namespace jevdecision\n" +
		"entity DecisionSpec id \"gooo://jev/decision/spec\"\n" +
		"entity DecisionReceipt id \"gooo://jev/decision/receipt\"\n" +
		"activity ObserveDecision(DecisionSpec) -> DecisionReceipt\n"
	document, err := Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	receipt, err := Generate(document)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !receipt.ExactSourceMatch {
		t.Fatalf("canonical source was not recognized as exact")
	}
}

func TestGeneratePreservesUnknownStageOnTampering(t *testing.T) {
	document, err := Parse(validContract)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	document.Entities[0].Digest = "tampered"

	receipt, err := Generate(document)
	if err == nil {
		t.Fatal("Generate() error = nil, want validation error")
	}
	if receipt.Status != "UNKNOWN" || receipt.MissingStage != "generation-validation" {
		t.Fatalf("unexpected failure receipt: %#v", receipt)
	}
}
