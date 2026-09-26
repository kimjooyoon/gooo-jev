package gooo

import "testing"

func lspSelfImprovementFeedbackInputs(t *testing.T) (string, string, RevisionSelfImprovementFeedback) {
	t.Helper()
	window := selfImprovementFeedbackWindow(t, 8, 2, false)
	feedback, err := ObserveRevisionSelfImprovementFeedback(window)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementFeedback() error = %v", err)
	}
	snapshot := Analyze(validContract)
	if err := snapshot.Validate(); err != nil {
		t.Fatalf("Analyze(validContract) Validate() error = %v", err)
	}
	if len(snapshot.Symbols) == 0 {
		t.Fatal("Analyze(validContract) returned no symbols")
	}
	return validContract, snapshot.Symbols[0].Name, feedback
}

func TestExplainSelfImprovementFeedbackBindsCandidateDeclaration(t *testing.T) {
	source, symbolName, feedback := lspSelfImprovementFeedbackInputs(t)
	result, err := ExplainSelfImprovementFeedback(source, symbolName, feedback)
	if err != nil {
		t.Fatalf("ExplainSelfImprovementFeedback() error = %v", err)
	}
	if result.Status != "BOUND" || result.SymbolName != symbolName ||
		result.FeedbackSignal != "remeasure" || !result.RequiresMeasurement {
		t.Fatalf("unexpected LSP self-improvement feedback: %#v", result)
	}
	if result.FeedbackDigest != feedback.FeedbackDigest ||
		result.WindowDigest != feedback.WindowDigest {
		t.Fatalf("LSP result lost feedback links: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestExplainSelfImprovementFeedbackRetainsFeedbackFailure(t *testing.T) {
	source, symbolName, feedback := lspSelfImprovementFeedbackInputs(t)
	feedback.FeedbackDigest = digestString("tampered")
	result, err := ExplainSelfImprovementFeedback(source, symbolName, feedback)
	if err == nil {
		t.Fatal("ExplainSelfImprovementFeedback() error = nil, want feedback failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-feedback-feedback" {
		t.Fatalf("unexpected unknown LSP feedback result: %#v", result)
	}
}

func TestExplainSelfImprovementFeedbackRetainsSourceLinkFailure(t *testing.T) {
	source, symbolName, feedback := lspSelfImprovementFeedbackInputs(t)
	result, err := ExplainSelfImprovementFeedback(source+"\n", symbolName, feedback)
	if err == nil {
		t.Fatal("ExplainSelfImprovementFeedback() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-feedback-link" {
		t.Fatalf("unexpected unknown LSP link result: %#v", result)
	}
}

func TestExplainSelfImprovementFeedbackRetainsMissingSymbol(t *testing.T) {
	source, _, feedback := lspSelfImprovementFeedbackInputs(t)
	result, err := ExplainSelfImprovementFeedback(source, "missing_symbol", feedback)
	if err == nil {
		t.Fatal("ExplainSelfImprovementFeedback() error = nil, want missing symbol")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-feedback-symbol" {
		t.Fatalf("unexpected unknown LSP symbol result: %#v", result)
	}
}

func TestExplainSelfImprovementFeedbackIsDeterministic(t *testing.T) {
	source, symbolName, feedback := lspSelfImprovementFeedbackInputs(t)
	first, err := ExplainSelfImprovementFeedback(source, symbolName, feedback)
	if err != nil {
		t.Fatalf("first ExplainSelfImprovementFeedback() error = %v", err)
	}
	second, err := ExplainSelfImprovementFeedback(source, symbolName, feedback)
	if err != nil {
		t.Fatalf("second ExplainSelfImprovementFeedback() error = %v", err)
	}
	if first.ResultDigest != second.ResultDigest {
		t.Fatal("same feedback and symbol produced different result digest")
	}
}