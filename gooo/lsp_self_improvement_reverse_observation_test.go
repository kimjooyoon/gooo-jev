package gooo

import "testing"

func lspSelfImprovementReverseObservationInputs(t *testing.T) (string, string, RevisionSelfImprovementLifecycleObservation, RevisionSelfImprovementReverseObservation) {
	t.Helper()
	source, symbolName, lifecycle := lspSelfImprovementLifecycleInputs(t)
	_, generation := generationAssessmentInputs(t)
	reverse, err := ObserveRevisionSelfImprovementReverseGeneration(lifecycle, generation)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementReverseGeneration() error = %v", err)
	}
	return source, symbolName, lifecycle, reverse
}

func TestExplainSelfImprovementReverseObservationBindsDeclarationAndProvenance(t *testing.T) {
	source, symbolName, lifecycle, reverse := lspSelfImprovementReverseObservationInputs(t)
	result, err := ExplainSelfImprovementReverseObservation(source, symbolName, lifecycle, reverse)
	if err != nil {
		t.Fatalf("ExplainSelfImprovementReverseObservation() error = %v", err)
	}
	if result.Status != "BOUND" || result.SymbolName != symbolName ||
		result.ReverseSignal != reverse.ReverseSignal ||
		!result.ExactSourceMatch || !result.ExactIRMatch || !result.ExactStructureMatch ||
		result.MetricsBindingDigest != lifecycle.MetricsBindingDigest {
		t.Fatalf("unexpected LSP self-improvement reverse result: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestExplainSelfImprovementReverseObservationRetainsReverseFailure(t *testing.T) {
	source, symbolName, lifecycle, reverse := lspSelfImprovementReverseObservationInputs(t)
	reverse.ObservationDigest = digestString("tampered")
	result, err := ExplainSelfImprovementReverseObservation(source, symbolName, lifecycle, reverse)
	if err == nil {
		t.Fatal("ExplainSelfImprovementReverseObservation() error = nil, want reverse failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-reverse-reverse" {
		t.Fatalf("unexpected unknown LSP reverse result: %#v", result)
	}
}

func TestExplainSelfImprovementReverseObservationRetainsSourceLinkFailure(t *testing.T) {
	source, symbolName, lifecycle, reverse := lspSelfImprovementReverseObservationInputs(t)
	result, err := ExplainSelfImprovementReverseObservation(source+"\n", symbolName, lifecycle, reverse)
	if err == nil {
		t.Fatal("ExplainSelfImprovementReverseObservation() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-reverse-link" {
		t.Fatalf("unexpected unknown LSP reverse link: %#v", result)
	}
}

func TestExplainSelfImprovementReverseObservationRetainsMissingSymbol(t *testing.T) {
	source, _, lifecycle, reverse := lspSelfImprovementReverseObservationInputs(t)
	result, err := ExplainSelfImprovementReverseObservation(source, "missing_symbol", lifecycle, reverse)
	if err == nil {
		t.Fatal("ExplainSelfImprovementReverseObservation() error = nil, want missing symbol")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-reverse-symbol" {
		t.Fatalf("unexpected unknown LSP reverse symbol: %#v", result)
	}
}

func TestExplainSelfImprovementReverseObservationIsDeterministic(t *testing.T) {
	source, symbolName, lifecycle, reverse := lspSelfImprovementReverseObservationInputs(t)
	first, err := ExplainSelfImprovementReverseObservation(source, symbolName, lifecycle, reverse)
	if err != nil {
		t.Fatalf("first ExplainSelfImprovementReverseObservation() error = %v", err)
	}
	second, err := ExplainSelfImprovementReverseObservation(source, symbolName, lifecycle, reverse)
	if err != nil {
		t.Fatalf("second ExplainSelfImprovementReverseObservation() error = %v", err)
	}
	if first.ResultDigest != second.ResultDigest {
		t.Fatal("same reverse observation and symbol produced different result digest")
	}
}
