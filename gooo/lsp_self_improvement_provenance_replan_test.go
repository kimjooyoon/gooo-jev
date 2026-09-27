package gooo

import "testing"

func lspProvenanceReplanInputs(t *testing.T) (string, string, RevisionSelfImprovementProvenanceReplanObservation) {
	t.Helper()
	decision, plan := provenanceReplanInputs(t)
	replan, err := ObserveRevisionSelfImprovementProvenanceReplan(decision, plan)
	if err != nil { t.Fatalf("ObserveRevisionSelfImprovementProvenanceReplan() error = %v", err) }
	snapshot := Analyze(validContract)
	if err := snapshot.Validate(); err != nil { t.Fatalf("Analyze(validContract) Validate() error = %v", err) }
	if len(snapshot.Symbols) == 0 { t.Fatal("Analyze(validContract) returned no symbols") }
	return validContract, snapshot.Symbols[0].Name, replan
}

func TestExplainSelfImprovementProvenanceReplanBindsDeclaration(t *testing.T) {
	source, symbolName, replan := lspProvenanceReplanInputs(t)
	result, err := ExplainSelfImprovementProvenanceReplan(source, symbolName, replan)
	if err != nil { t.Fatalf("ExplainSelfImprovementProvenanceReplan() error = %v", err) }
	if result.Status != "BOUND" || result.SymbolName != symbolName || result.ReplanSignal != "provenance-replan-mismatch" || result.NextPlanDigest != replan.NextPlanDigest { t.Fatalf("unexpected LSP provenance replan: %#v", result) }
	if result.ProvenanceDecisionDigest != replan.ProvenanceDecisionDigest || result.ExecutionPlanObservationDigest != replan.ExecutionPlanObservationDigest { t.Fatalf("LSP replan lost links: %#v", result) }
	if err := result.Validate(); err != nil { t.Fatalf("Validate() error = %v", err) }
}

func TestExplainSelfImprovementProvenanceReplanRetainsFailure(t *testing.T) {
	source, symbolName, replan := lspProvenanceReplanInputs(t)
	replan.ObservationDigest = digestString("tampered")
	result, err := ExplainSelfImprovementProvenanceReplan(source, symbolName, replan)
	if err == nil { t.Fatal("ExplainSelfImprovementProvenanceReplan() error = nil, want replan failure") }
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-replan-replan" { t.Fatalf("unexpected unknown LSP replan: %#v", result) }
}

func TestExplainSelfImprovementProvenanceReplanRetainsLinkFailure(t *testing.T) {
	source, symbolName, replan := lspProvenanceReplanInputs(t)
	result, err := ExplainSelfImprovementProvenanceReplan(source+"\n", symbolName, replan)
	if err == nil { t.Fatal("ExplainSelfImprovementProvenanceReplan() error = nil, want link failure") }
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-replan-link" { t.Fatalf("unexpected unknown LSP replan link: %#v", result) }
}

func TestExplainSelfImprovementProvenanceReplanRetainsMissingSymbol(t *testing.T) {
	source, _, replan := lspProvenanceReplanInputs(t)
	result, err := ExplainSelfImprovementProvenanceReplan(source, "missing_symbol", replan)
	if err == nil { t.Fatal("ExplainSelfImprovementProvenanceReplan() error = nil, want missing symbol") }
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-replan-symbol" { t.Fatalf("unexpected unknown LSP replan symbol: %#v", result) }
}

func TestExplainSelfImprovementProvenanceReplanIsDeterministic(t *testing.T) {
	source, symbolName, replan := lspProvenanceReplanInputs(t)
	first, err := ExplainSelfImprovementProvenanceReplan(source, symbolName, replan)
	if err != nil { t.Fatalf("first ExplainSelfImprovementProvenanceReplan() error = %v", err) }
	second, err := ExplainSelfImprovementProvenanceReplan(source, symbolName, replan)
	if err != nil { t.Fatalf("second ExplainSelfImprovementProvenanceReplan() error = %v", err) }
	if first.ResultDigest != second.ResultDigest { t.Fatal("same provenance replan and symbol produced different result digest") }
}
