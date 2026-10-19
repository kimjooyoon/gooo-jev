package gooo

import (
	"os"
	"strings"
	"testing"
)

const richDeclaration = "package support_triage\n" +
	"namespace operations\n" +
	"entity Ticket id \"support-ticket\"\n" +
	"  property state text\n" +
	"  property confidence number\n" +
	"entity Route id \"support-route\"\n" +
	"  property destination text\n" +
	"activity route_ticket(ticket Ticket, confidence number) -> Route\n"

func TestParseDiscoveryDeclarationDialectRoundTrip(t *testing.T) {
	document, err := Parse(richDeclaration)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(document.Entities) != 2 || len(document.Entities[0].Properties) != 2 || len(document.Activities) != 1 || len(document.Activities[0].Parameters) != 2 {
		t.Fatalf("rich declarations were lost: %#v", document)
	}
	if document.Activities[0].Parameters[0].Name != "ticket" || document.Activities[0].Parameters[1].Type != "number" {
		t.Fatalf("ordered parameters were lost: %#v", document.Activities[0].Parameters)
	}
	if document.Activities[0].Parameters[0].Position.Column >= document.Activities[0].Parameters[1].Position.Column {
		t.Fatalf("parameter positions were lost: %#v", document.Activities[0].Parameters)
	}
	receipt, err := Generate(document)
	if err != nil || !receipt.StructureMatch || !strings.Contains(receipt.GeneratedSource, "property confidence number") || !strings.Contains(receipt.GeneratedSource, "route_ticket(ticket Ticket, confidence number)") {
		t.Fatalf("round trip lost declaration structure: receipt=%#v, err=%v", receipt, err)
	}
	discovery := DiscoverUsage(richDeclaration, "pro")
	if discovery.Status != "BOUND" {
		t.Fatalf("discovery status = %q, stage = %q", discovery.Status, discovery.MissingStage)
	}
	if err := discovery.Validate(); err != nil {
		t.Fatalf("discovery validation: %v", err)
	}
}

func TestSupportTriageExampleIsDiscoverable(t *testing.T) {
	source, err := os.ReadFile("../examples/support-triage/usecase.gooo")
	if err != nil {
		t.Fatal(err)
	}
	discovery := DiscoverUsage(string(source), "")
	if discovery.Status != "BOUND" {
		t.Fatalf("published example is not bound: %s", discovery.MissingStage)
	}
}

func TestParseRichDeclarationPreservesUnknownBoundaries(t *testing.T) {
	for _, test := range []struct {
		name   string
		source string
		stage  string
	}{
		{"orphan property", strings.Replace(richDeclaration, "entity Ticket id \"support-ticket\"\n", "", 1), "property-owner"},
		{"duplicate property", strings.Replace(richDeclaration, "  property confidence number\n", "  property state number\n", 1), "duplicate-property"},
		{"undefined return", strings.Replace(richDeclaration, "-> Route", "-> Missing", 1), "ir-validation"},
		{"undefined parameter", strings.Replace(richDeclaration, "ticket Ticket, confidence number", "ticket Ticket, confidence Missing", 1), "ir-validation"},
		{"malformed parameter", strings.Replace(richDeclaration, "confidence number) -> Route", "confidence number,) -> Route", 1), "syntax"},
	} {
		t.Run(test.name, func(t *testing.T) {
			document, err := Parse(test.source)
			parseErr, ok := err.(ParseError)
			if !ok || parseErr.Stage != test.stage || document.Status != "UNKNOWN" {
				t.Fatalf("status=%s, error=%v; want stage %s", document.Status, err, test.stage)
			}
		})
	}
}

func TestValidateRejectsRichDeclarationMutation(t *testing.T) {
	for _, change := range []struct {
		name   string
		mutate func(*DocumentIR)
	}{
		{"property", func(document *DocumentIR) { document.Entities[0].Properties[0].Type = "number" }},
		{"parameter", func(document *DocumentIR) { document.Activities[0].Parameters[1].Type = "text" }},
		{"parameter order", func(document *DocumentIR) {
			document.Activities[0].Parameters[0], document.Activities[0].Parameters[1] = document.Activities[0].Parameters[1], document.Activities[0].Parameters[0]
		}},
	} {
		t.Run(change.name, func(t *testing.T) {
			document, err := Parse(richDeclaration)
			if err != nil {
				t.Fatal(err)
			}
			change.mutate(&document)
			if err := document.Validate(); err == nil {
				t.Fatal("mutation must invalidate declaration")
			}
		})
	}
}

func TestDecisionGenerationPreservesRichBase(t *testing.T) {
	source := richDeclaration + "decision RouteChoice kind choice id \"route-choice\"\n"
	document, err := ParseDecision(source)
	if err != nil {
		t.Fatal(err)
	}
	receipt, err := GenerateDecision(document)
	if err != nil || !receipt.StructureMatch {
		t.Fatalf("decision generation lost rich base: receipt=%#v, err=%v", receipt, err)
	}
}
