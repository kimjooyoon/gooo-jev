package gooo

import "testing"

func lspSelfImprovementWindowInputs(t *testing.T) (string, string, RevisionSelfImprovementWindow) {
	t.Helper()
	baseline := selfImprovementReceiptWithProfile(t, 8, 3, false, true, true)
	candidate := selfImprovementReceiptWithProfile(t, 2, 1, false, true, true)
	window, err := ObserveRevisionSelfImprovementWindow(baseline, candidate)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementWindow() error = %v", err)
	}
	snapshot := Analyze(validContract)
	if err := snapshot.Validate(); err != nil {
		t.Fatalf("Analyze(validContract) Validate() error = %v", err)
	}
	if len(snapshot.Symbols) == 0 {
		t.Fatal("Analyze(validContract) returned no symbols")
	}
	return validContract, snapshot.Symbols[0].Name, window
}

func TestExplainSelfImprovementWindowBindsCandidateDeclaration(t *testing.T) {
	source, symbolName, window := lspSelfImprovementWindowInputs(t)
	result, err := ExplainSelfImprovementWindow(source, symbolName, window)
	if err != nil {
		t.Fatalf("ExplainSelfImprovementWindow() error = %v", err)
	}
	if result.Status != "BOUND" || result.SymbolName != symbolName ||
		result.ComparisonSignal != "narrower" || result.ByteDelta != -6 {
		t.Fatalf("unexpected LSP self-improvement window: %#v", result)
	}
	if result.WindowDigest != window.WindowDigest ||
		result.BaselineReceiptDigest != window.BaselineReceiptDigest ||
		result.CandidateReceiptDigest != window.CandidateReceiptDigest {
		t.Fatalf("LSP result lost window links: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestExplainSelfImprovementWindowRetainsWindowFailure(t *testing.T) {
	source, symbolName, window := lspSelfImprovementWindowInputs(t)
	window.WindowDigest = digestString("tampered")
	result, err := ExplainSelfImprovementWindow(source, symbolName, window)
	if err == nil {
		t.Fatal("ExplainSelfImprovementWindow() error = nil, want window failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-window-window" {
		t.Fatalf("unexpected unknown LSP window result: %#v", result)
	}
}

func TestExplainSelfImprovementWindowRetainsSourceLinkFailure(t *testing.T) {
	source, symbolName, window := lspSelfImprovementWindowInputs(t)
	result, err := ExplainSelfImprovementWindow(source+"\n", symbolName, window)
	if err == nil {
		t.Fatal("ExplainSelfImprovementWindow() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-window-link" {
		t.Fatalf("unexpected unknown LSP link result: %#v", result)
	}
}

func TestExplainSelfImprovementWindowRetainsMissingSymbol(t *testing.T) {
	source, _, window := lspSelfImprovementWindowInputs(t)
	result, err := ExplainSelfImprovementWindow(source, "missing_symbol", window)
	if err == nil {
		t.Fatal("ExplainSelfImprovementWindow() error = nil, want missing symbol")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-window-symbol" {
		t.Fatalf("unexpected unknown LSP symbol result: %#v", result)
	}
}

func TestExplainSelfImprovementWindowIsDeterministic(t *testing.T) {
	source, symbolName, window := lspSelfImprovementWindowInputs(t)
	first, err := ExplainSelfImprovementWindow(source, symbolName, window)
	if err != nil {
		t.Fatalf("first ExplainSelfImprovementWindow() error = %v", err)
	}
	second, err := ExplainSelfImprovementWindow(source, symbolName, window)
	if err != nil {
		t.Fatalf("second ExplainSelfImprovementWindow() error = %v", err)
	}
	if first.ResultDigest != second.ResultDigest {
		t.Fatal("same window and symbol produced different result digest")
	}
}