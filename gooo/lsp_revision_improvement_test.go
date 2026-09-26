package gooo

import "testing"

func lspRevisionImprovementInputs(t *testing.T) (string, string, RevisionEvidenceSummary, RevisionEvidenceSummary, RevisionImprovementObservation) {
	t.Helper()
	baseline := summaryForReplacement(t, "lineage")
	candidate := summaryForReplacement(t, "lineage-expanded")
	observation, err := ObserveRevisionImprovement(baseline, candidate)
	if err != nil {
		t.Fatalf("ObserveRevisionImprovement() error = %v", err)
	}
	snapshot := Analyze(validContract)
	if err := snapshot.Validate(); err != nil {
		t.Fatalf("Analyze(validContract) Validate() error = %v", err)
	}
	if len(snapshot.Symbols) == 0 {
		t.Fatal("Analyze(validContract) returned no symbols")
	}
	return validContract, snapshot.Symbols[0].Name, baseline, candidate, observation
}

func TestExplainRevisionImprovementBindsSymbolAndObservation(t *testing.T) {
	source, symbolName, baseline, candidate, observation := lspRevisionImprovementInputs(t)
	result, err := ExplainRevisionImprovement(source, symbolName, baseline, candidate, observation)
	if err != nil {
		t.Fatalf("ExplainRevisionImprovement() error = %v", err)
	}
	if result.Status != "BOUND" || result.SymbolName != symbolName || result.ObservationClass == "" {
		t.Fatalf("unexpected LSP revision improvement: %#v", result)
	}
	if result.BaselineSummaryDigest != baseline.SummaryDigest ||
		result.CandidateSummaryDigest != candidate.SummaryDigest ||
		result.ObservationDigest != observation.ObservationDigest {
		t.Fatalf("LSP result lost observation links: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestExplainRevisionImprovementRetainsObservationFailure(t *testing.T) {
	source, symbolName, baseline, candidate, observation := lspRevisionImprovementInputs(t)
	observation.ObservationDigest = digestString("tampered")
	result, err := ExplainRevisionImprovement(source, symbolName, baseline, candidate, observation)
	if err == nil {
		t.Fatal("ExplainRevisionImprovement() error = nil, want observation failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-revision-improvement-observation" {
		t.Fatalf("unexpected unknown observation result: %#v", result)
	}
}

func TestExplainRevisionImprovementRetainsMissingSymbol(t *testing.T) {
	source, _, baseline, candidate, observation := lspRevisionImprovementInputs(t)
	result, err := ExplainRevisionImprovement(source, "missing_symbol", baseline, candidate, observation)
	if err == nil {
		t.Fatal("ExplainRevisionImprovement() error = nil, want missing symbol")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-revision-improvement-symbol" {
		t.Fatalf("unexpected unknown symbol result: %#v", result)
	}
}

func TestExplainRevisionImprovementIsDeterministic(t *testing.T) {
	source, symbolName, baseline, candidate, observation := lspRevisionImprovementInputs(t)
	first, err := ExplainRevisionImprovement(source, symbolName, baseline, candidate, observation)
	if err != nil {
		t.Fatalf("first ExplainRevisionImprovement() error = %v", err)
	}
	second, err := ExplainRevisionImprovement(source, symbolName, baseline, candidate, observation)
	if err != nil {
		t.Fatalf("second ExplainRevisionImprovement() error = %v", err)
	}
	if first.ResultDigest != second.ResultDigest {
		t.Fatal("same LSP observation produced different result digest")
	}
}