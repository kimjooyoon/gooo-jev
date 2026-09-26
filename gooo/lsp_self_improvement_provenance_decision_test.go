package gooo

import "testing"

func lspSelfImprovementProvenanceDecisionInputs(t *testing.T) (string, string, RevisionSelfImprovementLifecycleObservation, RevisionSelfImprovementProvenanceHistoryObservation, RevisionSelfImprovementProvenanceDecisionObservation) {
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
	decision, err := ObserveRevisionSelfImprovementProvenanceDecision(history)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceDecision() error = %v", err)
	}
	return source, snapshot.Symbols[0].Name, current, history, decision
}

func TestExplainSelfImprovementProvenanceDecisionBindsEvidence(t *testing.T) {
	source, symbolName, current, history, decision := lspSelfImprovementProvenanceDecisionInputs(t)
	result, err := ExplainSelfImprovementProvenanceDecision(source, symbolName, current, history, decision)
	if err != nil {
		t.Fatalf("ExplainSelfImprovementProvenanceDecision() error = %v", err)
	}
	if result.Status != "BOUND" || result.SymbolName != symbolName ||
		result.DecisionSignal != "observe" ||
		result.HistoryDigest != history.HistoryDigest ||
		result.CurrentLifecycleObservationDigest != current.ObservationDigest {
		t.Fatalf("unexpected LSP provenance decision: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestExplainSelfImprovementProvenanceDecisionRetainsDecisionFailure(t *testing.T) {
	source, symbolName, current, history, decision := lspSelfImprovementProvenanceDecisionInputs(t)
	decision.DecisionDigest = digestString("tampered")
	result, err := ExplainSelfImprovementProvenanceDecision(source, symbolName, current, history, decision)
	if err == nil {
		t.Fatal("ExplainSelfImprovementProvenanceDecision() error = nil, want decision failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-decision-decision" {
		t.Fatalf("unexpected unknown LSP provenance decision: %#v", result)
	}
}

func TestExplainSelfImprovementProvenanceDecisionRetainsSourceLinkFailure(t *testing.T) {
	source, symbolName, current, history, decision := lspSelfImprovementProvenanceDecisionInputs(t)
	result, err := ExplainSelfImprovementProvenanceDecision(source+"\n", symbolName, current, history, decision)
	if err == nil {
		t.Fatal("ExplainSelfImprovementProvenanceDecision() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-decision-source-link" {
		t.Fatalf("unexpected unknown LSP provenance decision source link: %#v", result)
	}
}

func TestExplainSelfImprovementProvenanceDecisionRetainsMissingSymbol(t *testing.T) {
	source, _, current, history, decision := lspSelfImprovementProvenanceDecisionInputs(t)
	result, err := ExplainSelfImprovementProvenanceDecision(source, "missing_symbol", current, history, decision)
	if err == nil {
		t.Fatal("ExplainSelfImprovementProvenanceDecision() error = nil, want missing symbol")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-decision-symbol" {
		t.Fatalf("unexpected unknown LSP provenance decision symbol: %#v", result)
	}
}

func TestExplainSelfImprovementProvenanceDecisionIsDeterministic(t *testing.T) {
	source, symbolName, current, history, decision := lspSelfImprovementProvenanceDecisionInputs(t)
	first, err := ExplainSelfImprovementProvenanceDecision(source, symbolName, current, history, decision)
	if err != nil {
		t.Fatalf("first ExplainSelfImprovementProvenanceDecision() error = %v", err)
	}
	second, err := ExplainSelfImprovementProvenanceDecision(source, symbolName, current, history, decision)
	if err != nil {
		t.Fatalf("second ExplainSelfImprovementProvenanceDecision() error = %v", err)
	}
	if first.ResultDigest != second.ResultDigest {
		t.Fatal("same provenance decision and symbol produced different result digest")
	}
}
