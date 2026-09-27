package gooo

import "testing"

func lspSelfImprovementIterationProvenanceInputs(
	t *testing.T,
	changed bool,
) (string, string, RevisionSelfImprovementLifecycleObservation, RevisionSelfImprovementIterationProvenanceObservation) {
	t.Helper()
	source, symbolName, current := lspSelfImprovementLifecycleInputs(t)
	iteration, history, bridge := iterationProvenanceInputs(t, changed)
	observation, err := ObserveRevisionSelfImprovementIterationProvenance(iteration, history, bridge)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementIterationProvenance() error = %v", err)
	}
	return source, symbolName, current, observation
}

func TestExplainSelfImprovementIterationProvenanceBindsStableSymbol(t *testing.T) {
	source, symbolName, current, observation := lspSelfImprovementIterationProvenanceInputs(t, false)
	result, err := ExplainSelfImprovementIterationProvenance(source, symbolName, current, observation)
	if err != nil {
		t.Fatalf("ExplainSelfImprovementIterationProvenance() error = %v", err)
	}
	if result.Status != "BOUND" ||
		result.SymbolName != symbolName ||
		result.ProvenanceHistorySignal != "stable" ||
		result.DecisionSignal != "observe" ||
		result.FeedbackSignal != "observe" {
		t.Fatalf("unexpected stable LSP iteration provenance: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestExplainSelfImprovementIterationProvenanceBindsChangedSymbol(t *testing.T) {
	source, symbolName, current, observation := lspSelfImprovementIterationProvenanceInputs(t, true)
	result, err := ExplainSelfImprovementIterationProvenance(source, symbolName, current, observation)
	if err != nil {
		t.Fatalf("ExplainSelfImprovementIterationProvenance() error = %v", err)
	}
	if result.Status != "BOUND" ||
		result.ProvenanceHistorySignal != "mixed" ||
		result.LegacyHistorySignal != "mixed" ||
		result.DecisionSignal != "inspect" ||
		result.FeedbackSignal != "inspect" ||
		!result.RequiresMeasurement {
		t.Fatalf("unexpected changed LSP iteration provenance: %#v", result)
	}
}

func TestExplainSelfImprovementIterationProvenanceRetainsObservationFailure(t *testing.T) {
	source, symbolName, current, observation := lspSelfImprovementIterationProvenanceInputs(t, false)
	observation.ObservationDigest = digestString("tampered")
	result, err := ExplainSelfImprovementIterationProvenance(source, symbolName, current, observation)
	if err == nil {
		t.Fatal("ExplainSelfImprovementIterationProvenance() error = nil, want observation failure")
	}
	if result.Status != "UNKNOWN" ||
		result.MissingStage != "lsp-self-improvement-iteration-provenance-observation" {
		t.Fatalf("unexpected unknown LSP iteration provenance observation: %#v", result)
	}
}

func TestExplainSelfImprovementIterationProvenanceRetainsSourceLinkFailure(t *testing.T) {
	source, symbolName, current, observation := lspSelfImprovementIterationProvenanceInputs(t, false)
	result, err := ExplainSelfImprovementIterationProvenance(source+"\n", symbolName, current, observation)
	if err == nil {
		t.Fatal("ExplainSelfImprovementIterationProvenance() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" ||
		result.MissingStage != "lsp-self-improvement-iteration-provenance-source-link" {
		t.Fatalf("unexpected unknown LSP iteration provenance source link: %#v", result)
	}
}

func TestExplainSelfImprovementIterationProvenanceRetainsMissingSymbol(t *testing.T) {
	source, _, current, observation := lspSelfImprovementIterationProvenanceInputs(t, false)
	result, err := ExplainSelfImprovementIterationProvenance(source, "missing_symbol", current, observation)
	if err == nil {
		t.Fatal("ExplainSelfImprovementIterationProvenance() error = nil, want missing symbol")
	}
	if result.Status != "UNKNOWN" ||
		result.MissingStage != "lsp-self-improvement-iteration-provenance-symbol" {
		t.Fatalf("unexpected unknown LSP iteration provenance symbol: %#v", result)
	}
}

func TestExplainSelfImprovementIterationProvenanceIsDeterministic(t *testing.T) {
	source, symbolName, current, observation := lspSelfImprovementIterationProvenanceInputs(t, true)
	first, err := ExplainSelfImprovementIterationProvenance(source, symbolName, current, observation)
	if err != nil {
		t.Fatalf("first ExplainSelfImprovementIterationProvenance() error = %v", err)
	}
	second, err := ExplainSelfImprovementIterationProvenance(source, symbolName, current, observation)
	if err != nil {
		t.Fatalf("second ExplainSelfImprovementIterationProvenance() error = %v", err)
	}
	if first.ResultDigest != second.ResultDigest {
		t.Fatal("same iteration provenance and symbol produced different result digest")
	}
}
