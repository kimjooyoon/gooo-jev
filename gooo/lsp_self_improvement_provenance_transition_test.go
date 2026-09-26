package gooo

import "testing"

func lspSelfImprovementProvenanceTransitionInputs(t *testing.T) (string, string, RevisionSelfImprovementLifecycleObservation, RevisionSelfImprovementLifecycleObservation, RevisionSelfImprovementProvenanceTransitionObservation) {
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
	return source, snapshot.Symbols[0].Name, previous, current, transition
}

func TestExplainSelfImprovementProvenanceTransitionBindsDeclarationAndEvidence(t *testing.T) {
	source, symbolName, previous, current, transition := lspSelfImprovementProvenanceTransitionInputs(t)
	result, err := ExplainSelfImprovementProvenanceTransition(source, symbolName, previous, current, transition)
	if err != nil {
		t.Fatalf("ExplainSelfImprovementProvenanceTransition() error = %v", err)
	}
	if result.Status != "BOUND" || result.SymbolName != symbolName ||
		result.TransitionSignal != "provenance-stable" ||
		result.PreviousLifecycleObservationDigest != previous.ObservationDigest ||
		result.CurrentLifecycleObservationDigest != current.ObservationDigest {
		t.Fatalf("unexpected LSP provenance transition: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestExplainSelfImprovementProvenanceTransitionRetainsTransitionFailure(t *testing.T) {
	source, symbolName, previous, current, transition := lspSelfImprovementProvenanceTransitionInputs(t)
	transition.ObservationDigest = digestString("tampered")
	result, err := ExplainSelfImprovementProvenanceTransition(source, symbolName, previous, current, transition)
	if err == nil {
		t.Fatal("ExplainSelfImprovementProvenanceTransition() error = nil, want transition failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-transition-transition" {
		t.Fatalf("unexpected unknown LSP transition: %#v", result)
	}
}

func TestExplainSelfImprovementProvenanceTransitionRetainsSourceLinkFailure(t *testing.T) {
	source, symbolName, previous, current, transition := lspSelfImprovementProvenanceTransitionInputs(t)
	result, err := ExplainSelfImprovementProvenanceTransition(source+"\n", symbolName, previous, current, transition)
	if err == nil {
		t.Fatal("ExplainSelfImprovementProvenanceTransition() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-transition-source-link" {
		t.Fatalf("unexpected unknown LSP transition source link: %#v", result)
	}
}

func TestExplainSelfImprovementProvenanceTransitionRetainsMissingSymbol(t *testing.T) {
	source, _, previous, current, transition := lspSelfImprovementProvenanceTransitionInputs(t)
	result, err := ExplainSelfImprovementProvenanceTransition(source, "missing_symbol", previous, current, transition)
	if err == nil {
		t.Fatal("ExplainSelfImprovementProvenanceTransition() error = nil, want missing symbol")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-transition-symbol" {
		t.Fatalf("unexpected unknown LSP transition symbol: %#v", result)
	}
}

func TestExplainSelfImprovementProvenanceTransitionIsDeterministic(t *testing.T) {
	source, symbolName, previous, current, transition := lspSelfImprovementProvenanceTransitionInputs(t)
	first, err := ExplainSelfImprovementProvenanceTransition(source, symbolName, previous, current, transition)
	if err != nil {
		t.Fatalf("first ExplainSelfImprovementProvenanceTransition() error = %v", err)
	}
	second, err := ExplainSelfImprovementProvenanceTransition(source, symbolName, previous, current, transition)
	if err != nil {
		t.Fatalf("second ExplainSelfImprovementProvenanceTransition() error = %v", err)
	}
	if first.ResultDigest != second.ResultDigest {
		t.Fatal("same provenance transition and symbol produced different result digest")
	}
}
