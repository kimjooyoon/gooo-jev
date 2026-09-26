package gooo

import "testing"

func lspSelfImprovementOutcomeInputs(t *testing.T) (string, string, SelfImprovementApplicationObservation, RevisionMetricsBinding, RevisionGenerationAssessment) {
	t.Helper()
	iteration, planObservation, applicationObservation := selfImprovementApplicationInputs(t)
	application, err := ObserveRevisionSelfImprovementApplication(iteration, planObservation, applicationObservation)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementApplication() error = %v", err)
	}
	_, metricsBinding, _, generationAssessment := selfImprovementReceiptInputs(t)
	snapshot := Analyze(validContract)
	if err := snapshot.Validate(); err != nil {
		t.Fatalf("Analyze(validContract) Validate() error = %v", err)
	}
	if len(snapshot.Symbols) == 0 {
		t.Fatal("Analyze(validContract) returned no symbols")
	}
	return validContract, snapshot.Symbols[0].Name, application, metricsBinding, generationAssessment
}

func TestExplainSelfImprovementOutcomeBindsDeclarationAndEvidence(t *testing.T) {
	source, symbolName, application, metrics, generation := lspSelfImprovementOutcomeInputs(t)
	result, err := ExplainSelfImprovementOutcome(source, symbolName, application, metrics, generation)
	if err != nil {
		t.Fatalf("ExplainSelfImprovementOutcome() error = %v", err)
	}
	if result.Status != "BOUND" || result.SymbolName != symbolName ||
		result.ApplicationDigest != application.ApplicationDigest ||
		result.MetricsDigest != metrics.MetricsDigest ||
		result.GeneratedIRDigest != generation.GeneratedIRDigest {
		t.Fatalf("unexpected LSP self-improvement outcome: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestExplainSelfImprovementOutcomeRetainsApplicationFailure(t *testing.T) {
	source, symbolName, application, metrics, generation := lspSelfImprovementOutcomeInputs(t)
	application.ObservationDigest = digestString("tampered")
	result, err := ExplainSelfImprovementOutcome(source, symbolName, application, metrics, generation)
	if err == nil {
		t.Fatal("ExplainSelfImprovementOutcome() error = nil, want application failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-outcome-application" {
		t.Fatalf("unexpected unknown LSP application result: %#v", result)
	}
}

func TestExplainSelfImprovementOutcomeRetainsMetricsFailure(t *testing.T) {
	source, symbolName, application, metrics, generation := lspSelfImprovementOutcomeInputs(t)
	metrics.BindingDigest = digestString("tampered")
	result, err := ExplainSelfImprovementOutcome(source, symbolName, application, metrics, generation)
	if err == nil {
		t.Fatal("ExplainSelfImprovementOutcome() error = nil, want metrics failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-outcome-metrics" {
		t.Fatalf("unexpected unknown LSP metrics result: %#v", result)
	}
}

func TestExplainSelfImprovementOutcomeRetainsSourceLinkFailure(t *testing.T) {
	source, symbolName, application, metrics, generation := lspSelfImprovementOutcomeInputs(t)
	result, err := ExplainSelfImprovementOutcome(source+"
", symbolName, application, metrics, generation)
	if err == nil {
		t.Fatal("ExplainSelfImprovementOutcome() error = nil, want source-link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-outcome-source-link" {
		t.Fatalf("unexpected unknown LSP source-link result: %#v", result)
	}
}

func TestExplainSelfImprovementOutcomeRetainsGenerationFailure(t *testing.T) {
	source, symbolName, application, metrics, generation := lspSelfImprovementOutcomeInputs(t)
	generation.GenerationDigest = digestString("tampered")
	result, err := ExplainSelfImprovementOutcome(source, symbolName, application, metrics, generation)
	if err == nil {
		t.Fatal("ExplainSelfImprovementOutcome() error = nil, want generation failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-outcome-generation" {
		t.Fatalf("unexpected unknown LSP generation result: %#v", result)
	}
}

func TestExplainSelfImprovementOutcomeRetainsMissingSymbol(t *testing.T) {
	source, _, application, metrics, generation := lspSelfImprovementOutcomeInputs(t)
	result, err := ExplainSelfImprovementOutcome(source, "missing_symbol", application, metrics, generation)
	if err == nil {
		t.Fatal("ExplainSelfImprovementOutcome() error = nil, want missing symbol")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-outcome-symbol" {
		t.Fatalf("unexpected unknown LSP symbol result: %#v", result)
	}
}

func TestExplainSelfImprovementOutcomeIsDeterministic(t *testing.T) {
	source, symbolName, application, metrics, generation := lspSelfImprovementOutcomeInputs(t)
	first, err := ExplainSelfImprovementOutcome(source, symbolName, application, metrics, generation)
	if err != nil {
		t.Fatalf("first ExplainSelfImprovementOutcome() error = %v", err)
	}
	second, err := ExplainSelfImprovementOutcome(source, symbolName, application, metrics, generation)
	if err != nil {
		t.Fatalf("second ExplainSelfImprovementOutcome() error = %v", err)
	}
	if first.ResultDigest != second.ResultDigest {
		t.Fatal("same evidence and symbol produced different result digest")
	}
}
