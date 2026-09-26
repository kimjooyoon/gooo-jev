package gooo

import "testing"

func lspSelfImprovementDecisionInputs(t *testing.T) (string, string, RevisionSelfImprovementDecisionObservation) {
	t.Helper()
	iteration, outcome, feedback := selfImprovementDecisionInputs(t)
	decision, err := ObserveRevisionSelfImprovementDecision(iteration, outcome, feedback)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementDecision() error = %v", err)
	}
	snapshot := Analyze(validContract)
	if err := snapshot.Validate(); err != nil {
		t.Fatalf("Analyze(validContract) Validate() error = %v", err)
	}
	if len(snapshot.Symbols) == 0 {
		t.Fatal("Analyze(validContract) returned no symbols")
	}
	return validContract, snapshot.Symbols[0].Name, decision
}

func TestExplainSelfImprovementDecisionBindsDeclarationAndDecision(t *testing.T) {
	source, symbolName, decision := lspSelfImprovementDecisionInputs(t)
	result, err := ExplainSelfImprovementDecision(source, symbolName, decision)
	if err != nil {
		t.Fatalf("ExplainSelfImprovementDecision() error = %v", err)
	}
	if result.Status != "BOUND" || result.SymbolName != symbolName ||
		result.DecisionSignal != decision.DecisionSignal ||
		result.DecisionReason != decision.DecisionReason ||
		result.OutcomeDigest != decision.OutcomeDigest {
		t.Fatalf("unexpected LSP self-improvement decision: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestExplainSelfImprovementDecisionRetainsDecisionFailure(t *testing.T) {
	source, symbolName, decision := lspSelfImprovementDecisionInputs(t)
	decision.DecisionDigest = digestString("tampered")
	result, err := ExplainSelfImprovementDecision(source, symbolName, decision)
	if err == nil {
		t.Fatal("ExplainSelfImprovementDecision() error = nil, want decision failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-decision-decision" {
		t.Fatalf("unexpected unknown LSP decision result: %#v", result)
	}
}

func TestExplainSelfImprovementDecisionRetainsSourceLinkFailure(t *testing.T) {
	source, symbolName, decision := lspSelfImprovementDecisionInputs(t)
	result, err := ExplainSelfImprovementDecision(source+"\n", symbolName, decision)
	if err == nil {
		t.Fatal("ExplainSelfImprovementDecision() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-decision-link" {
		t.Fatalf("unexpected unknown LSP decision link: %#v", result)
	}
}

func TestExplainSelfImprovementDecisionRetainsMissingSymbol(t *testing.T) {
	source, _, decision := lspSelfImprovementDecisionInputs(t)
	result, err := ExplainSelfImprovementDecision(source, "missing_symbol", decision)
	if err == nil {
		t.Fatal("ExplainSelfImprovementDecision() error = nil, want missing symbol")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-decision-symbol" {
		t.Fatalf("unexpected unknown LSP decision symbol: %#v", result)
	}
}

func TestExplainSelfImprovementDecisionIsDeterministic(t *testing.T) {
	source, symbolName, decision := lspSelfImprovementDecisionInputs(t)
	first, err := ExplainSelfImprovementDecision(source, symbolName, decision)
	if err != nil {
		t.Fatalf("first ExplainSelfImprovementDecision() error = %v", err)
	}
	second, err := ExplainSelfImprovementDecision(source, symbolName, decision)
	if err != nil {
		t.Fatalf("second ExplainSelfImprovementDecision() error = %v", err)
	}
	if first.ResultDigest != second.ResultDigest {
		t.Fatal("same decision and symbol produced different result digest")
	}
}
