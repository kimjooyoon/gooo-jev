package gooo

import "testing"

const validDecisionSource = "package jevdecision\n" +
	"namespace jevdecision\n" +
	"entity DecisionSpec id \"gooo://jev/decision/spec\"\n" +
	"entity DecisionReceipt id \"gooo://jev/decision/receipt\"\n" +
	"decision ChooseReceipt kind choice id \"gooo://jev/decision/choose-receipt\"\n" +
	"activity ObserveDecision(DecisionSpec) -> DecisionReceipt\n"

func TestParseDecisionProducesTypedBoundIR(t *testing.T) {
	document, err := ParseDecision(validDecisionSource)
	if err != nil {
		t.Fatalf("ParseDecision() error = %v", err)
	}
	if document.Status != "BOUND" || document.MissingStage != "" {
		t.Fatalf("unexpected status: %#v", document)
	}
	if len(document.Decisions) != 1 || document.Decisions[0].Kind != ChoiceDecision {
		t.Fatalf("unexpected decisions: %#v", document.Decisions)
	}
	if document.Decisions[0].Position.Line != 5 {
		t.Fatalf("unexpected decision position: %#v", document.Decisions[0].Position)
	}
	if document.Base.Activities[0].Position.Line != 6 {
		t.Fatalf("base position was not preserved: %#v", document.Base.Activities[0].Position)
	}
	if document.SourceDigest == "" || document.IRDigest == "" || document.Decisions[0].Digest == "" {
		t.Fatalf("missing decision provenance: %#v", document)
	}
	if !document.NonExecuting || !document.NonAuthorizing {
		t.Fatalf("decision IR must not authorize execution")
	}
	if err := document.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestParseDecisionSupportsTypedKinds(t *testing.T) {
	source := "package jevdecision\n" +
		"namespace jevdecision\n" +
		"entity DecisionSpec id \"gooo://jev/decision/spec\"\n" +
		"entity DecisionReceipt id \"gooo://jev/decision/receipt\"\n" +
		"decision ChooseReceipt kind choice id \"gooo://jev/decision/choice\"\n" +
		"decision Confidence kind score id \"gooo://jev/decision/score\"\n" +
		"decision Approved kind boolean id \"gooo://jev/decision/boolean\"\n" +
		"activity ObserveDecision(DecisionSpec) -> DecisionReceipt\n"
	document, err := ParseDecision(source)
	if err != nil {
		t.Fatalf("ParseDecision() error = %v", err)
	}
	if len(document.Decisions) != 3 ||
		document.Decisions[0].Kind != ChoiceDecision ||
		document.Decisions[1].Kind != ScoreDecision ||
		document.Decisions[2].Kind != BooleanDecision {
		t.Fatalf("unexpected typed decisions: %#v", document.Decisions)
	}
}

func TestParseDecisionRetainsDecisionSyntaxUnknownStage(t *testing.T) {
	source := "package jevdecision\n" +
		"namespace jevdecision\n" +
		"entity DecisionSpec id \"gooo://jev/decision/spec\"\n" +
		"decision Choose kind invalid id \"gooo://jev/decision/invalid\"\n"
	document, err := ParseDecision(source)
	if err == nil {
		t.Fatal("ParseDecision() error = nil, want syntax error")
	}
	parseErr, ok := err.(ParseError)
	if !ok || parseErr.Stage != "decision-syntax" {
		t.Fatalf("unexpected parse error: %#v", err)
	}
	if document.Status != "UNKNOWN" || document.MissingStage != "decision-syntax" {
		t.Fatalf("unexpected unknown state: %#v", document)
	}
}

func TestParseDecisionRetainsBaseUnknownStage(t *testing.T) {
	source := "package jevdecision\n" +
		"namespace jevdecision\n" +
		"entity DecisionSpec id \"gooo://jev/decision/spec\"\n" +
		"decision Choose kind choice id \"gooo://jev/decision/choice\"\n" +
		"activity broken\n"
	document, err := ParseDecision(source)
	if err == nil {
		t.Fatal("ParseDecision() error = nil, want base parse error")
	}
	parseErr, ok := err.(ParseError)
	if !ok || parseErr.Stage != "decision-base-parse" {
		t.Fatalf("unexpected parse error: %#v", err)
	}
	if document.Status != "UNKNOWN" || document.MissingStage != "decision-base-parse" {
		t.Fatalf("unexpected unknown state: %#v", document)
	}
}

func TestValidateRejectsTamperedDecisionDigest(t *testing.T) {
	document, err := ParseDecision(validDecisionSource)
	if err != nil {
		t.Fatalf("ParseDecision() error = %v", err)
	}
	document.Decisions[0].Digest = "tampered"
	if err := document.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want tamper rejection")
	}
}

func TestParseDecisionSourceDigestChangesWithSource(t *testing.T) {
	first, err := ParseDecision(validDecisionSource)
	if err != nil {
		t.Fatalf("first ParseDecision() error = %v", err)
	}
	second, err := ParseDecision(validDecisionSource + "\n")
	if err != nil {
		t.Fatalf("second ParseDecision() error = %v", err)
	}
	if first.SourceDigest == second.SourceDigest || first.IRDigest == second.IRDigest {
		t.Fatal("source change did not change decision provenance")
	}
}
