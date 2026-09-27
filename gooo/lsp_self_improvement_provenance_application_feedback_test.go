package gooo

import "testing"

func lspProvenanceApplicationFeedbackInputs(t *testing.T) (string, string, RevisionSelfImprovementProvenanceApplicationFeedbackObservation) {
	t.Helper()
	disposition, application := provenanceApplicationFeedbackInputs(t)
	feedback, err := ObserveRevisionSelfImprovementProvenanceApplicationFeedback(disposition, application)
	if err != nil { t.Fatalf("ObserveRevisionSelfImprovementProvenanceApplicationFeedback() error = %v", err) }
	snapshot := Analyze(validContract)
	if err := snapshot.Validate(); err != nil { t.Fatalf("Analyze(validContract) Validate() error = %v", err) }
	if len(snapshot.Symbols) == 0 { t.Fatal("Analyze(validContract) returned no symbols") }
	return validContract, snapshot.Symbols[0].Name, feedback
}

func TestExplainSelfImprovementProvenanceApplicationFeedbackBindsDeclaration(t *testing.T) {
	source, symbolName, feedback := lspProvenanceApplicationFeedbackInputs(t)
	result, err := ExplainSelfImprovementProvenanceApplicationFeedback(source, symbolName, feedback)
	if err != nil { t.Fatalf("ExplainSelfImprovementProvenanceApplicationFeedback() error = %v", err) }
	if result.Status != "BOUND" || result.SymbolName != symbolName || result.DispositionSignal != feedback.DispositionSignal || result.FeedbackApplicationSignal != "provenance-application-mismatch" { t.Fatalf("unexpected LSP provenance application feedback: %#v", result) }
	if result.ApplicationFeedbackDigest != feedback.ObservationDigest || result.ReceiptDigest != feedback.ReceiptDigest { t.Fatalf("LSP application feedback lost links: %#v", result) }
	if err := result.Validate(); err != nil { t.Fatalf("Validate() error = %v", err) }
}

func TestExplainSelfImprovementProvenanceApplicationFeedbackRetainsFeedbackFailure(t *testing.T) {
	source, symbolName, feedback := lspProvenanceApplicationFeedbackInputs(t)
	feedback.ObservationDigest = digestString("tampered")
	result, err := ExplainSelfImprovementProvenanceApplicationFeedback(source, symbolName, feedback)
	if err == nil { t.Fatal("ExplainSelfImprovementProvenanceApplicationFeedback() error = nil, want feedback failure") }
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-application-feedback-feedback" { t.Fatalf("unexpected unknown LSP provenance application feedback: %#v", result) }
}

func TestExplainSelfImprovementProvenanceApplicationFeedbackRetainsLinkFailure(t *testing.T) {
	source, symbolName, feedback := lspProvenanceApplicationFeedbackInputs(t)
	result, err := ExplainSelfImprovementProvenanceApplicationFeedback(source+"\n", symbolName, feedback)
	if err == nil { t.Fatal("ExplainSelfImprovementProvenanceApplicationFeedback() error = nil, want link failure") }
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-application-feedback-link" { t.Fatalf("unexpected unknown LSP provenance application feedback link: %#v", result) }
}

func TestExplainSelfImprovementProvenanceApplicationFeedbackRetainsMissingSymbol(t *testing.T) {
	source, _, feedback := lspProvenanceApplicationFeedbackInputs(t)
	result, err := ExplainSelfImprovementProvenanceApplicationFeedback(source, "missing_symbol", feedback)
	if err == nil { t.Fatal("ExplainSelfImprovementProvenanceApplicationFeedback() error = nil, want missing symbol") }
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-application-feedback-symbol" { t.Fatalf("unexpected unknown LSP provenance application feedback symbol: %#v", result) }
}

func TestExplainSelfImprovementProvenanceApplicationFeedbackIsDeterministic(t *testing.T) {
	source, symbolName, feedback := lspProvenanceApplicationFeedbackInputs(t)
	first, err := ExplainSelfImprovementProvenanceApplicationFeedback(source, symbolName, feedback)
	if err != nil { t.Fatalf("first ExplainSelfImprovementProvenanceApplicationFeedback() error = %v", err) }
	second, err := ExplainSelfImprovementProvenanceApplicationFeedback(source, symbolName, feedback)
	if err != nil { t.Fatalf("second ExplainSelfImprovementProvenanceApplicationFeedback() error = %v", err) }
	if first.ResultDigest != second.ResultDigest { t.Fatal("same application feedback and symbol produced different result digest") }
}
