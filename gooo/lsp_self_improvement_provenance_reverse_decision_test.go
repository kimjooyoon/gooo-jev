package gooo

import "testing"

func lspProvenanceReverseDecisionInputs(t *testing.T) (string, string, RevisionSelfImprovementProvenanceReverseDecisionObservation) {
	t.Helper()
	_, _, decision := provenanceReverseDecisionInputs(t)
	feedback, reverse, decisionInput := provenanceReverseDecisionInputs(t)
	result, err := ObserveRevisionSelfImprovementProvenanceReverseDecision(feedback, reverse, decisionInput)
	if err != nil { t.Fatalf("ObserveRevisionSelfImprovementProvenanceReverseDecision() error = %v", err) }
	_ = decision
	snapshot := Analyze(validContract)
	if err := snapshot.Validate(); err != nil { t.Fatalf("Analyze(validContract) Validate() error = %v", err) }
	if len(snapshot.Symbols) == 0 { t.Fatal("Analyze(validContract) returned no symbols") }
	return validContract, snapshot.Symbols[0].Name, result
}

func TestExplainSelfImprovementProvenanceReverseDecisionBindsDeclaration(t *testing.T) {
	source, symbolName, decision := lspProvenanceReverseDecisionInputs(t)
	result, err := ExplainSelfImprovementProvenanceReverseDecision(source, symbolName, decision)
	if err != nil { t.Fatalf("ExplainSelfImprovementProvenanceReverseDecision() error = %v", err) }
	if result.Status != "BOUND" || result.SymbolName != symbolName || result.DecisionSignal != decision.DecisionSignal || result.DecisionFeedbackSignal != "provenance-decision-mismatch" { t.Fatalf("unexpected LSP provenance reverse decision: %#v", result) }
	if result.ApplicationFeedbackDigest != decision.ApplicationFeedbackDigest || result.ReverseObservationDigest != decision.ReverseObservationDigest { t.Fatalf("LSP decision lost provenance links: %#v", result) }
	if err := result.Validate(); err != nil { t.Fatalf("Validate() error = %v", err) }
}

func TestExplainSelfImprovementProvenanceReverseDecisionRetainsDecisionFailure(t *testing.T) {
	source, symbolName, decision := lspProvenanceReverseDecisionInputs(t)
	decision.DecisionDigest = digestString("tampered")
	result, err := ExplainSelfImprovementProvenanceReverseDecision(source, symbolName, decision)
	if err == nil { t.Fatal("ExplainSelfImprovementProvenanceReverseDecision() error = nil, want decision failure") }
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-reverse-decision-decision" { t.Fatalf("unexpected unknown LSP provenance decision: %#v", result) }
}

func TestExplainSelfImprovementProvenanceReverseDecisionRetainsLinkFailure(t *testing.T) {
	source, symbolName, decision := lspProvenanceReverseDecisionInputs(t)
	result, err := ExplainSelfImprovementProvenanceReverseDecision(source+"\n", symbolName, decision)
	if err == nil { t.Fatal("ExplainSelfImprovementProvenanceReverseDecision() error = nil, want link failure") }
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-reverse-decision-link" { t.Fatalf("unexpected unknown LSP provenance decision link: %#v", result) }
}

func TestExplainSelfImprovementProvenanceReverseDecisionRetainsMissingSymbol(t *testing.T) {
	source, _, decision := lspProvenanceReverseDecisionInputs(t)
	result, err := ExplainSelfImprovementProvenanceReverseDecision(source, "missing_symbol", decision)
	if err == nil { t.Fatal("ExplainSelfImprovementProvenanceReverseDecision() error = nil, want missing symbol") }
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-reverse-decision-symbol" { t.Fatalf("unexpected unknown LSP provenance decision symbol: %#v", result) }
}

func TestExplainSelfImprovementProvenanceReverseDecisionIsDeterministic(t *testing.T) {
	source, symbolName, decision := lspProvenanceReverseDecisionInputs(t)
	first, err := ExplainSelfImprovementProvenanceReverseDecision(source, symbolName, decision)
	if err != nil { t.Fatalf("first ExplainSelfImprovementProvenanceReverseDecision() error = %v", err) }
	second, err := ExplainSelfImprovementProvenanceReverseDecision(source, symbolName, decision)
	if err != nil { t.Fatalf("second ExplainSelfImprovementProvenanceReverseDecision() error = %v", err) }
	if first.ResultDigest != second.ResultDigest { t.Fatal("same provenance decision and symbol produced different result digest") }
}
