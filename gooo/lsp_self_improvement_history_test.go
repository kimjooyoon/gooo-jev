package gooo

import "testing"

func lspSelfImprovementHistoryInputs(t *testing.T) (string, string, RevisionSelfImprovementHistory) {
	t.Helper()
	windows, feedback := selfImprovementHistoryInputs(t)
	history, err := ObserveRevisionSelfImprovementHistory(windows, feedback)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementHistory() error = %v", err)
	}
	snapshot := Analyze(validContract)
	if err := snapshot.Validate(); err != nil {
		t.Fatalf("Analyze(validContract) Validate() error = %v", err)
	}
	if len(snapshot.Symbols) == 0 {
		t.Fatal("Analyze(validContract) returned no symbols")
	}
	return validContract, snapshot.Symbols[0].Name, history
}

func TestExplainSelfImprovementHistoryBindsLastCandidateDeclaration(t *testing.T) {
	source, symbolName, history := lspSelfImprovementHistoryInputs(t)
	result, err := ExplainSelfImprovementHistory(source, symbolName, history)
	if err != nil {
		t.Fatalf("ExplainSelfImprovementHistory() error = %v", err)
	}
	if result.Status != "BOUND" || result.SymbolName != symbolName ||
		result.ObservationCount != history.ObservationCount ||
		result.LastCandidateSourceDigest != history.LastCandidateSourceDigest {
		t.Fatalf("unexpected LSP self-improvement history: %#v", result)
	}
	if result.HistoryDigest != history.HistoryDigest ||
		result.LastWindowDigest != history.LastWindowDigest ||
		result.LastFeedbackDigest != history.LastFeedbackDigest {
		t.Fatalf("LSP result lost history links: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestExplainSelfImprovementHistoryRetainsHistoryFailure(t *testing.T) {
	source, symbolName, history := lspSelfImprovementHistoryInputs(t)
	history.HistoryDigest = digestString("tampered")
	result, err := ExplainSelfImprovementHistory(source, symbolName, history)
	if err == nil {
		t.Fatal("ExplainSelfImprovementHistory() error = nil, want history failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-history-history" {
		t.Fatalf("unexpected unknown LSP history result: %#v", result)
	}
}

func TestExplainSelfImprovementHistoryRetainsSourceLinkFailure(t *testing.T) {
	source, symbolName, history := lspSelfImprovementHistoryInputs(t)
	result, err := ExplainSelfImprovementHistory(source+"\n", symbolName, history)
	if err == nil {
		t.Fatal("ExplainSelfImprovementHistory() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-history-link" {
		t.Fatalf("unexpected unknown LSP link result: %#v", result)
	}
}

func TestExplainSelfImprovementHistoryRetainsMissingSymbol(t *testing.T) {
	source, _, history := lspSelfImprovementHistoryInputs(t)
	result, err := ExplainSelfImprovementHistory(source, "missing_symbol", history)
	if err == nil {
		t.Fatal("ExplainSelfImprovementHistory() error = nil, want missing symbol")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-history-symbol" {
		t.Fatalf("unexpected unknown LSP symbol result: %#v", result)
	}
}

func TestExplainSelfImprovementHistoryIsDeterministic(t *testing.T) {
	source, symbolName, history := lspSelfImprovementHistoryInputs(t)
	first, err := ExplainSelfImprovementHistory(source, symbolName, history)
	if err != nil {
		t.Fatalf("first ExplainSelfImprovementHistory() error = %v", err)
	}
	second, err := ExplainSelfImprovementHistory(source, symbolName, history)
	if err != nil {
		t.Fatalf("second ExplainSelfImprovementHistory() error = %v", err)
	}
	if first.ResultDigest != second.ResultDigest {
		t.Fatal("same history and symbol produced different result digest")
	}
}