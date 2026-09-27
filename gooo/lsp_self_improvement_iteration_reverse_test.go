package gooo

import "testing"

func lspSelfImprovementIterationReverseInputs(
	t *testing.T,
) (string, string, RevisionSelfImprovementLifecycleObservation, RevisionSelfImprovementIterationProvenanceObservation, RevisionSelfImprovementReverseObservation) {
	t.Helper()
	source, symbolName, current := lspSelfImprovementLifecycleInputs(t)
	iteration, history, bridge := iterationProvenanceInputs(t, false)
	iterationProvenance, err := ObserveRevisionSelfImprovementIterationProvenance(iteration, history, bridge)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementIterationProvenance() error = %v", err)
	}
	_, generation := generationAssessmentInputs(t)
	reverse, err := ObserveRevisionSelfImprovementReverseGeneration(current, generation)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementReverseGeneration() error = %v", err)
	}
	return source, symbolName, current, iterationProvenance, reverse
}

func TestExplainSelfImprovementIterationReverseBindsEvidence(t *testing.T) {
	source, symbolName, current, iteration, reverse := lspSelfImprovementIterationReverseInputs(t)
	result, err := ExplainSelfImprovementIterationReverseObservation(source, symbolName, current, iteration, reverse)
	if err != nil {
		t.Fatalf("ExplainSelfImprovementIterationReverseObservation() error = %v", err)
	}
	if result.Status != "BOUND" ||
		result.SymbolName != symbolName ||
		result.RelationSignal == "iteration-reverse-unknown" ||
		result.ReverseSignal != reverse.ReverseSignal ||
		!result.ExactIRMatch {
		t.Fatalf("unexpected LSP iteration reverse result: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestExplainSelfImprovementIterationReverseRetainsIterationFailure(t *testing.T) {
	source, symbolName, current, iteration, reverse := lspSelfImprovementIterationReverseInputs(t)
	iteration.ObservationDigest = digestString("tampered")
	result, err := ExplainSelfImprovementIterationReverseObservation(source, symbolName, current, iteration, reverse)
	if err == nil {
		t.Fatal("ExplainSelfImprovementIterationReverseObservation() error = nil, want iteration failure")
	}
	if result.Status != "UNKNOWN" ||
		result.MissingStage != "lsp-self-improvement-iteration-reverse-iteration" {
		t.Fatalf("unexpected unknown LSP iteration reverse iteration: %#v", result)
	}
}

func TestExplainSelfImprovementIterationReverseRetainsSourceLinkFailure(t *testing.T) {
	source, symbolName, current, iteration, reverse := lspSelfImprovementIterationReverseInputs(t)
	result, err := ExplainSelfImprovementIterationReverseObservation(source+"\n", symbolName, current, iteration, reverse)
	if err == nil {
		t.Fatal("ExplainSelfImprovementIterationReverseObservation() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" ||
		result.MissingStage != "lsp-self-improvement-iteration-reverse-source-link" {
		t.Fatalf("unexpected unknown LSP iteration reverse source link: %#v", result)
	}
}

func TestExplainSelfImprovementIterationReverseRetainsMissingSymbol(t *testing.T) {
	source, _, current, iteration, reverse := lspSelfImprovementIterationReverseInputs(t)
	result, err := ExplainSelfImprovementIterationReverseObservation(source, "missing_symbol", current, iteration, reverse)
	if err == nil {
		t.Fatal("ExplainSelfImprovementIterationReverseObservation() error = nil, want missing symbol")
	}
	if result.Status != "UNKNOWN" ||
		result.MissingStage != "lsp-self-improvement-iteration-reverse-symbol" {
		t.Fatalf("unexpected unknown LSP iteration reverse symbol: %#v", result)
	}
}

func TestExplainSelfImprovementIterationReverseIsDeterministic(t *testing.T) {
	source, symbolName, current, iteration, reverse := lspSelfImprovementIterationReverseInputs(t)
	first, err := ExplainSelfImprovementIterationReverseObservation(source, symbolName, current, iteration, reverse)
	if err != nil {
		t.Fatalf("first ExplainSelfImprovementIterationReverseObservation() error = %v", err)
	}
	second, err := ExplainSelfImprovementIterationReverseObservation(source, symbolName, current, iteration, reverse)
	if err != nil {
		t.Fatalf("second ExplainSelfImprovementIterationReverseObservation() error = %v", err)
	}
	if first.ResultDigest != second.ResultDigest {
		t.Fatal("same iteration, reverse observation, and symbol produced different result digest")
	}
}
