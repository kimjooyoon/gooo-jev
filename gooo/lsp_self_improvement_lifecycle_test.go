package gooo

import "testing"

func lspSelfImprovementLifecycleInputs(t *testing.T) (string, string, RevisionSelfImprovementLifecycleObservation) {
	t.Helper()
	executionApplication, outcome := selfImprovementLifecycleInputs(t)
	lifecycle, err := ObserveRevisionSelfImprovementLifecycle(executionApplication, outcome)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementLifecycle() error = %v", err)
	}
	snapshot := Analyze(validContract)
	if err := snapshot.Validate(); err != nil {
		t.Fatalf("Analyze(validContract) Validate() error = %v", err)
	}
	if len(snapshot.Symbols) == 0 {
		t.Fatal("Analyze(validContract) returned no symbols")
	}
	return validContract, snapshot.Symbols[0].Name, lifecycle
}

func TestExplainSelfImprovementLifecycleBindsDeclarationAndProvenance(t *testing.T) {
	source, symbolName, lifecycle := lspSelfImprovementLifecycleInputs(t)
	result, err := ExplainSelfImprovementLifecycle(source, symbolName, lifecycle)
	if err != nil {
		t.Fatalf("ExplainSelfImprovementLifecycle() error = %v", err)
	}
	if result.Status != "BOUND" || result.SymbolName != symbolName ||
		result.LifecycleSignal != lifecycle.LifecycleSignal ||
		result.MetricsBindingDigest != lifecycle.MetricsBindingDigest ||
		result.GeneratedIRDigest != lifecycle.GeneratedIRDigest ||
		result.PlanDigest != lifecycle.PlanDigest {
		t.Fatalf("unexpected LSP self-improvement lifecycle: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestExplainSelfImprovementLifecycleRetainsLifecycleFailure(t *testing.T) {
	source, symbolName, lifecycle := lspSelfImprovementLifecycleInputs(t)
	lifecycle.ObservationDigest = digestString("tampered")
	result, err := ExplainSelfImprovementLifecycle(source, symbolName, lifecycle)
	if err == nil {
		t.Fatal("ExplainSelfImprovementLifecycle() error = nil, want lifecycle failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-lifecycle-lifecycle" {
		t.Fatalf("unexpected unknown LSP lifecycle result: %#v", result)
	}
}

func TestExplainSelfImprovementLifecycleRetainsSourceLinkFailure(t *testing.T) {
	source, symbolName, lifecycle := lspSelfImprovementLifecycleInputs(t)
	result, err := ExplainSelfImprovementLifecycle(source+"\n", symbolName, lifecycle)
	if err == nil {
		t.Fatal("ExplainSelfImprovementLifecycle() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-lifecycle-link" {
		t.Fatalf("unexpected unknown LSP lifecycle link: %#v", result)
	}
}

func TestExplainSelfImprovementLifecycleRetainsMissingSymbol(t *testing.T) {
	source, _, lifecycle := lspSelfImprovementLifecycleInputs(t)
	result, err := ExplainSelfImprovementLifecycle(source, "missing_symbol", lifecycle)
	if err == nil {
		t.Fatal("ExplainSelfImprovementLifecycle() error = nil, want missing symbol")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-lifecycle-symbol" {
		t.Fatalf("unexpected unknown LSP lifecycle symbol: %#v", result)
	}
}

func TestExplainSelfImprovementLifecycleIsDeterministic(t *testing.T) {
	source, symbolName, lifecycle := lspSelfImprovementLifecycleInputs(t)
	first, err := ExplainSelfImprovementLifecycle(source, symbolName, lifecycle)
	if err != nil {
		t.Fatalf("first ExplainSelfImprovementLifecycle() error = %v", err)
	}
	second, err := ExplainSelfImprovementLifecycle(source, symbolName, lifecycle)
	if err != nil {
		t.Fatalf("second ExplainSelfImprovementLifecycle() error = %v", err)
	}
	if first.ResultDigest != second.ResultDigest {
		t.Fatal("same lifecycle and symbol produced different result digest")
	}
}
