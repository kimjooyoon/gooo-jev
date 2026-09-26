package gooo

import "testing"

func lspSelfImprovementProvenanceHistoryInputs(t *testing.T) (string, string, RevisionSelfImprovementLifecycleObservation, RevisionSelfImprovementProvenanceHistoryObservation) {
	t.Helper()
	source := validContract
	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		t.Fatalf("Analyze(validContract) Validate() error = %v", err)
	}
	if len(snapshot.Symbols) == 0 {
		t.Fatal("Analyze(validContract) returned no symbols")
	}
	previous, current := selfImprovementProvenanceTransitionInputs(t)
	transition, err := ObserveRevisionSelfImprovementProvenanceTransition(previous, current)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceTransition() error = %v", err)
	}
	history, err := ObserveRevisionSelfImprovementProvenanceHistory([]RevisionSelfImprovementProvenanceTransitionObservation{transition})
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceHistory() error = %v", err)
	}
	return source, snapshot.Symbols[0].Name, current, history
}

func TestExplainSelfImprovementProvenanceHistoryBindsDeclarationAndHistory(t *testing.T) {
	source, symbolName, current, history := lspSelfImprovementProvenanceHistoryInputs(t)
	result, err := ExplainSelfImprovementProvenanceHistory(source, symbolName, current, history)
	if err != nil {
		t.Fatalf("ExplainSelfImprovementProvenanceHistory() error = %v", err)
	}
	if result.Status != "BOUND" || result.SymbolName != symbolName ||
		result.HistorySignal != "stable" ||
		result.LastCurrentLifecycleDigest != current.ObservationDigest ||
		result.ObservationCount != 1 {
		t.Fatalf("unexpected LSP provenance history: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestExplainSelfImprovementProvenanceHistoryRetainsHistoryFailure(t *testing.T) {
	source, symbolName, current, history := lspSelfImprovementProvenanceHistoryInputs(t)
	history.HistoryDigest = digestString("tampered")
	result, err := ExplainSelfImprovementProvenanceHistory(source, symbolName, current, history)
	if err == nil {
		t.Fatal("ExplainSelfImprovementProvenanceHistory() error = nil, want history failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-history-history" {
		t.Fatalf("unexpected unknown LSP history: %#v", result)
	}
}

func TestExplainSelfImprovementProvenanceHistoryRetainsSourceLinkFailure(t *testing.T) {
	source, symbolName, current, history := lspSelfImprovementProvenanceHistoryInputs(t)
	result, err := ExplainSelfImprovementProvenanceHistory(source+"\n", symbolName, current, history)
	if err == nil {
		t.Fatal("ExplainSelfImprovementProvenanceHistory() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-history-source-link" {
		t.Fatalf("unexpected unknown LSP history source link: %#v", result)
	}
}

func TestExplainSelfImprovementProvenanceHistoryRetainsMissingSymbol(t *testing.T) {
	source, _, current, history := lspSelfImprovementProvenanceHistoryInputs(t)
	result, err := ExplainSelfImprovementProvenanceHistory(source, "missing_symbol", current, history)
	if err == nil {
		t.Fatal("ExplainSelfImprovementProvenanceHistory() error = nil, want missing symbol")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-history-symbol" {
		t.Fatalf("unexpected unknown LSP history symbol: %#v", result)
	}
}

func TestExplainSelfImprovementProvenanceHistoryIsDeterministic(t *testing.T) {
	source, symbolName, current, history := lspSelfImprovementProvenanceHistoryInputs(t)
	first, err := ExplainSelfImprovementProvenanceHistory(source, symbolName, current, history)
	if err != nil {
		t.Fatalf("first ExplainSelfImprovementProvenanceHistory() error = %v", err)
	}
	second, err := ExplainSelfImprovementProvenanceHistory(source, symbolName, current, history)
	if err != nil {
		t.Fatalf("second ExplainSelfImprovementProvenanceHistory() error = %v", err)
	}
	if first.ResultDigest != second.ResultDigest {
		t.Fatal("same provenance history and symbol produced different result digest")
	}
}
