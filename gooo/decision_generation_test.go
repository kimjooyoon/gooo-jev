package gooo

import (
	"strings"
	"testing"
)

func TestGenerateDecisionRecordsTypedRoundTrip(t *testing.T) {
	document, err := ParseDecision(validDecisionSource)
	if err != nil {
		t.Fatalf("ParseDecision() error = %v", err)
	}
	receipt, err := GenerateDecision(document)
	if err != nil {
		t.Fatalf("GenerateDecision() error = %v", err)
	}
	if receipt.Status != "BOUND" || receipt.MissingStage != "" {
		t.Fatalf("unexpected receipt: %#v", receipt)
	}
	if !receipt.ExactSourceMatch || !receipt.StructureMatch {
		t.Fatalf("canonical source was not observed exactly: %#v", receipt)
	}
	if receipt.GeneratedIRDigest == "" || receipt.StructureDigest == "" {
		t.Fatalf("missing generation evidence: %#v", receipt)
	}
	if !strings.Contains(receipt.GeneratedSource, "decision ChooseReceipt kind choice") {
		t.Fatalf("generated source is missing typed decision: %q", receipt.GeneratedSource)
	}
	observed, err := ParseDecision(receipt.GeneratedSource)
	if err != nil {
		t.Fatalf("generated ParseDecision() error = %v", err)
	}
	if observed.IRDigest != receipt.GeneratedIRDigest {
		t.Fatalf("generated IR digest = %q, parsed digest = %q", receipt.GeneratedIRDigest, observed.IRDigest)
	}
}

func TestGenerateDecisionSeparatesSourceFidelityFromStructure(t *testing.T) {
	source := strings.Replace(validDecisionSource, "namespace jevdecision\n", "namespace jevdecision\n# source comment\n", 1)
	document, err := ParseDecision(source)
	if err != nil {
		t.Fatalf("ParseDecision() error = %v", err)
	}
	receipt, err := GenerateDecision(document)
	if err != nil {
		t.Fatalf("GenerateDecision() error = %v", err)
	}
	if receipt.ExactSourceMatch {
		t.Fatal("comment-bearing source should not be exact canonical source")
	}
	if !receipt.StructureMatch {
		t.Fatal("comment-bearing source should preserve typed structure")
	}
}

func TestGenerateDecisionRetainsUnknownValidationStage(t *testing.T) {
	document, err := ParseDecision(validDecisionSource)
	if err != nil {
		t.Fatalf("ParseDecision() error = %v", err)
	}
	document.Decisions[0].Digest = "tampered"
	receipt, err := GenerateDecision(document)
	if err == nil {
		t.Fatal("GenerateDecision() error = nil, want validation error")
	}
	if receipt.Status != "UNKNOWN" || receipt.MissingStage != "decision-generation-validation" {
		t.Fatalf("unexpected failure receipt: %#v", receipt)
	}
}

func TestGenerateDecisionIsDeterministic(t *testing.T) {
	document, err := ParseDecision(validDecisionSource)
	if err != nil {
		t.Fatalf("ParseDecision() error = %v", err)
	}
	first, err := GenerateDecision(document)
	if err != nil {
		t.Fatalf("first GenerateDecision() error = %v", err)
	}
	second, err := GenerateDecision(document)
	if err != nil {
		t.Fatalf("second GenerateDecision() error = %v", err)
	}
	if first.GeneratedSourceDigest != second.GeneratedSourceDigest || first.StructureDigest != second.StructureDigest {
		t.Fatal("same typed decision IR produced different generation evidence")
	}
}
