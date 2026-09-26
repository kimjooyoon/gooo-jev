package gooo

import "testing"

func lspRevisionEvidenceInputs(t *testing.T) (string, string, RevisionEvidenceSummary) {
	t.Helper()
	chain, binding, metrics := evidenceSummaryInputs(t)
	summary, err := SummarizeRevisionEvidence(chain, binding, metrics)
	if err != nil {
		t.Fatalf("SummarizeRevisionEvidence() error = %v", err)
	}
	snapshot := Analyze(validContract)
	if err := snapshot.Validate(); err != nil {
		t.Fatalf("Analyze(validContract) Validate() error = %v", err)
	}
	if len(snapshot.Symbols) == 0 {
		t.Fatal("Analyze(validContract) returned no symbols")
	}
	return validContract, snapshot.Symbols[0].Name, summary
}

func TestExplainRevisionEvidenceBindsSymbolAndProvenance(t *testing.T) {
	source, symbolName, summary := lspRevisionEvidenceInputs(t)
	result, err := ExplainRevisionEvidence(source, symbolName, summary)
	if err != nil {
		t.Fatalf("ExplainRevisionEvidence() error = %v", err)
	}
	if result.Status != "BOUND" || result.SymbolName != symbolName || result.EvidenceStageCount != 6 {
		t.Fatalf("unexpected LSP revision evidence: %#v", result)
	}
	if result.SourceDigest != summary.SourceDigest ||
		result.SummaryDigest != summary.SummaryDigest ||
		result.ChainDigest != summary.ChainDigest {
		t.Fatalf("LSP result lost provenance links: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestExplainRevisionEvidenceRetainsSourceLinkFailure(t *testing.T) {
	source, symbolName, summary := lspRevisionEvidenceInputs(t)
	result, err := ExplainRevisionEvidence(source+"\n", symbolName, summary)
	if err == nil {
		t.Fatal("ExplainRevisionEvidence() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-revision-evidence-link" {
		t.Fatalf("unexpected unknown LSP link result: %#v", result)
	}
}

func TestExplainRevisionEvidenceRetainsMissingSymbol(t *testing.T) {
	source, _, summary := lspRevisionEvidenceInputs(t)
	result, err := ExplainRevisionEvidence(source, "missing_symbol", summary)
	if err == nil {
		t.Fatal("ExplainRevisionEvidence() error = nil, want missing symbol")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "lsp-revision-evidence-symbol" {
		t.Fatalf("unexpected unknown LSP symbol result: %#v", result)
	}
}

func TestExplainRevisionEvidenceIsDeterministic(t *testing.T) {
	source, symbolName, summary := lspRevisionEvidenceInputs(t)
	first, err := ExplainRevisionEvidence(source, symbolName, summary)
	if err != nil {
		t.Fatalf("first ExplainRevisionEvidence() error = %v", err)
	}
	second, err := ExplainRevisionEvidence(source, symbolName, summary)
	if err != nil {
		t.Fatalf("second ExplainRevisionEvidence() error = %v", err)
	}
	if first.ResultDigest != second.ResultDigest {
		t.Fatal("same LSP evidence produced different result digest")
	}
}