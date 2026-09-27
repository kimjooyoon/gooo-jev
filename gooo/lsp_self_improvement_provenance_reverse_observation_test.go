package gooo

import "testing"

func lspProvenanceReverseObservationInputs(t *testing.T) (string, string, RevisionSelfImprovementProvenanceApplicationFeedbackObservation, RevisionSelfImprovementProvenanceReverseObservation) {
	t.Helper()
	feedback, reverse := provenanceReverseObservationInputs(t)
	snapshot := Analyze(validContract)
	if err := snapshot.Validate(); err != nil { t.Fatalf("Analyze(validContract) Validate() error = %v", err) }
	if len(snapshot.Symbols) == 0 { t.Fatal("Analyze(validContract) returned no symbols") }
	return validContract, snapshot.Symbols[0].Name, feedback, reverse
}

func TestExplainSelfImprovementProvenanceReverseObservationBindsDeclaration(t *testing.T) {
	source, symbolName, feedback, reverse := lspProvenanceReverseObservationInputs(t)
	result, err := ExplainSelfImprovementProvenanceReverseObservation(source, symbolName, feedback, reverse)
	if err != nil { t.Fatalf("ExplainSelfImprovementProvenanceReverseObservation() error = %v", err) }
	if result.Status != "BOUND" || result.SymbolName != symbolName || result.ReverseSignal != "reverse-observed" || result.ProvenanceReverseSignal != "provenance-reverse-mismatch" || !result.ExactIRMatch || !result.ExactStructureMatch { t.Fatalf("unexpected LSP provenance reverse result: %#v", result) }
	if result.ApplicationFeedbackDigest != feedback.ObservationDigest || result.ReverseObservationDigest != reverse.ObservationDigest { t.Fatalf("LSP provenance reverse result lost links: %#v", result) }
	if err := result.Validate(); err != nil { t.Fatalf("Validate() error = %v", err) }
}

func TestExplainSelfImprovementProvenanceReverseObservationRetainsReverseFailure(t *testing.T) {
	source, symbolName, feedback, reverse := lspProvenanceReverseObservationInputs(t)
	reverse.ObservationDigest = digestString("tampered")
	result, err := ExplainSelfImprovementProvenanceReverseObservation(source, symbolName, feedback, reverse)
	if err == nil { t.Fatal("ExplainSelfImprovementProvenanceReverseObservation() error = nil, want reverse failure") }
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-reverse-reverse" { t.Fatalf("unexpected unknown LSP provenance reverse result: %#v", result) }
}

func TestExplainSelfImprovementProvenanceReverseObservationRetainsLinkFailure(t *testing.T) {
	source, symbolName, feedback, reverse := lspProvenanceReverseObservationInputs(t)
	result, err := ExplainSelfImprovementProvenanceReverseObservation(source+"\n", symbolName, feedback, reverse)
	if err == nil { t.Fatal("ExplainSelfImprovementProvenanceReverseObservation() error = nil, want link failure") }
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-reverse-link" { t.Fatalf("unexpected unknown LSP provenance reverse link: %#v", result) }
}

func TestExplainSelfImprovementProvenanceReverseObservationRetainsMissingSymbol(t *testing.T) {
	source, _, feedback, reverse := lspProvenanceReverseObservationInputs(t)
	result, err := ExplainSelfImprovementProvenanceReverseObservation(source, "missing_symbol", feedback, reverse)
	if err == nil { t.Fatal("ExplainSelfImprovementProvenanceReverseObservation() error = nil, want missing symbol") }
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-self-improvement-provenance-reverse-symbol" { t.Fatalf("unexpected unknown LSP provenance reverse symbol: %#v", result) }
}

func TestExplainSelfImprovementProvenanceReverseObservationIsDeterministic(t *testing.T) {
	source, symbolName, feedback, reverse := lspProvenanceReverseObservationInputs(t)
	first, err := ExplainSelfImprovementProvenanceReverseObservation(source, symbolName, feedback, reverse)
	if err != nil { t.Fatalf("first ExplainSelfImprovementProvenanceReverseObservation() error = %v", err) }
	second, err := ExplainSelfImprovementProvenanceReverseObservation(source, symbolName, feedback, reverse)
	if err != nil { t.Fatalf("second ExplainSelfImprovementProvenanceReverseObservation() error = %v", err) }
	if first.ResultDigest != second.ResultDigest { t.Fatal("same provenance reverse evidence and symbol produced different result digest") }
}
