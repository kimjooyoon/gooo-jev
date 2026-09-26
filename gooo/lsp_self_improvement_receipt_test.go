package gooo

import "testing"

func lspSelfImprovementReceiptInputs(t *testing.T) (string, string, RevisionSelfImprovementReceipt) {
	t.Helper()
	applicationObservation, metricsBinding, applicationAssessment, generationAssessment := selfImprovementReceiptInputs(t)
	receipt, err := ObserveRevisionSelfImprovementReceipt(applicationObservation, metricsBinding, applicationAssessment, generationAssessment)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementReceipt() error = %v", err)
	}
	snapshot := Analyze(validContract)
	if err := snapshot.Validate(); err != nil {
		t.Fatalf("Analyze(validContract) Validate() error = %v", err)
	}
	if len(snapshot.Symbols) == 0 {
		t.Fatal("Analyze(validContract) returned no symbols")
	}
	return validContract, snapshot.Symbols[0].Name, receipt
}

func TestExplainSelfImprovementBindsDeclarationAndReceipt(t *testing.T) {
	source, symbolName, receipt := lspSelfImprovementReceiptInputs(t)
	result, err := ExplainSelfImprovement(source, symbolName, receipt)
	if err != nil {
		t.Fatalf("ExplainSelfImprovement() error = %v", err)
	}
	if result.Status != "BOUND" || result.SymbolName != symbolName || result.StageCount != 4 {
		t.Fatalf("unexpected LSP self-improvement receipt: %#v", result)
	}
	if result.ReceiptDigest != receipt.ReceiptDigest ||
		result.ApplicationDigest != receipt.ApplicationDigest ||
		result.GeneratedIRDigest != receipt.GeneratedIRDigest {
		t.Fatalf("LSP result lost receipt links: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestExplainSelfImprovementRetainsReceiptFailure(t *testing.T) {
	source, symbolName, receipt := lspSelfImprovementReceiptInputs(t)
	receipt.ReceiptDigest = digestString("tampered")
	result, err := ExplainSelfImprovement(source, symbolName, receipt)
	if err == nil {
		t.Fatal("ExplainSelfImprovement() error = nil, want receipt failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-receipt-receipt" {
		t.Fatalf("unexpected unknown LSP receipt result: %#v", result)
	}
}

func TestExplainSelfImprovementRetainsSourceLinkFailure(t *testing.T) {
	source, symbolName, receipt := lspSelfImprovementReceiptInputs(t)
	result, err := ExplainSelfImprovement(source+"\n", symbolName, receipt)
	if err == nil {
		t.Fatal("ExplainSelfImprovement() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-receipt-link" {
		t.Fatalf("unexpected unknown LSP link result: %#v", result)
	}
}

func TestExplainSelfImprovementRetainsMissingSymbol(t *testing.T) {
	source, _, receipt := lspSelfImprovementReceiptInputs(t)
	result, err := ExplainSelfImprovement(source, "missing_symbol", receipt)
	if err == nil {
		t.Fatal("ExplainSelfImprovement() error = nil, want missing symbol")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-receipt-symbol" {
		t.Fatalf("unexpected unknown LSP symbol result: %#v", result)
	}
}

func TestExplainSelfImprovementIsDeterministic(t *testing.T) {
	source, symbolName, receipt := lspSelfImprovementReceiptInputs(t)
	first, err := ExplainSelfImprovement(source, symbolName, receipt)
	if err != nil {
		t.Fatalf("first ExplainSelfImprovement() error = %v", err)
	}
	second, err := ExplainSelfImprovement(source, symbolName, receipt)
	if err != nil {
		t.Fatalf("second ExplainSelfImprovement() error = %v", err)
	}
	if first.ResultDigest != second.ResultDigest {
		t.Fatal("same receipt and symbol produced different result digest")
	}
}