package gooo

import "fmt"

// SelfImprovementWindowSymbolResult is a declaration-scoped, read-only LSP
// projection of a baseline/candidate self-improvement comparison.
type SelfImprovementWindowSymbolResult struct {
	Status                 string
	MissingStage           string
	SourceDigest           string
	IRDigest               string
	SymbolName             string
	SymbolKind             SymbolKind
	SymbolDigest            string
	BaselineReceiptDigest   string
	CandidateReceiptDigest  string
	ByteDelta              int
	LineDelta              int
	ComparisonSignal        string
	SourceStable            bool
	CandidateStable         bool
	IRChangeStable          bool
	StructureStable         bool
	ExactSourceStable       bool
	GeneratedIRStable       bool
	WindowDigest            string
	ResultDigest            string
	NonExecuting            bool
	NonAuthorizing          bool
}

// ExplainSelfImprovementWindow resolves the current candidate declaration
// against a validated receipt window without executing or authorizing it.
func ExplainSelfImprovementWindow(source, symbolName string, window RevisionSelfImprovementWindow) (SelfImprovementWindowSymbolResult, error) {
	result := SelfImprovementWindowSymbolResult{
		Status:                "UNKNOWN",
		MissingStage:          "lsp-self-improvement-window",
		SourceDigest:          digestString(source),
		BaselineReceiptDigest:  window.BaselineReceiptDigest,
		CandidateReceiptDigest: window.CandidateReceiptDigest,
		ByteDelta:             window.ByteDelta,
		LineDelta:             window.LineDelta,
		ComparisonSignal:      window.ComparisonSignal,
		SourceStable:          window.SourceStable,
		CandidateStable:       window.CandidateStable,
		IRChangeStable:        window.IRChangeStable,
		StructureStable:       window.StructureStable,
		ExactSourceStable:     window.ExactSourceStable,
		GeneratedIRStable:     window.GeneratedIRStable,
		WindowDigest:          window.WindowDigest,
		NonExecuting:          true,
		NonAuthorizing:        true,
	}
	setResultDigest := func() {
		result.ResultDigest = digestLSPSelfImprovementWindow(result)
	}
	setResultDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-window-snapshot"
		setResultDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := window.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-window-window"
		setResultDigest()
		return result, fmt.Errorf("revision self-improvement window is not valid: %w", err)
	}
	if snapshot.SourceDigest != window.CandidateSourceDigest {
		result.MissingStage = "lsp-self-improvement-window-link"
		setResultDigest()
		return result, fmt.Errorf("LSP source digest does not match candidate window source")
	}
	for _, symbol := range snapshot.Symbols {
		if symbol.Name == symbolName {
			result.Status = "BOUND"
			result.MissingStage = ""
			result.IRDigest = snapshot.IRDigest
			result.SymbolName = symbol.Name
			result.SymbolKind = symbol.Kind
			result.SymbolDigest = symbol.Digest
			setResultDigest()
			if err := result.Validate(); err != nil {
				result.Status = "UNKNOWN"
				result.MissingStage = "lsp-self-improvement-window"
				setResultDigest()
				return result, fmt.Errorf("LSP self-improvement window result is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-self-improvement-window-symbol"
	setResultDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r SelfImprovementWindowSymbolResult) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("LSP self-improvement window status is empty")
	}
	if r.Status == "BOUND" && r.MissingStage != "" {
		return fmt.Errorf("bound LSP self-improvement window has a missing stage")
	}
	if r.Status == "UNKNOWN" && r.MissingStage == "" {
		return fmt.Errorf("unknown LSP self-improvement window has no missing stage")
	}
	if r.ComparisonSignal != "narrower" && r.ComparisonSignal != "wider" &&
		r.ComparisonSignal != "stable" && r.ComparisonSignal != "mixed" {
		return fmt.Errorf("LSP self-improvement window signal is invalid")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("LSP self-improvement window must remain non-executing and non-authorizing")
	}
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" ||
		r.SymbolName == "" || r.SymbolDigest == "") {
		return fmt.Errorf("bound LSP self-improvement window is incomplete")
	}
	if digestLSPSelfImprovementWindow(r) != r.ResultDigest {
		return fmt.Errorf("LSP self-improvement window digest does not match its fields")
	}
	return nil
}

func digestLSPSelfImprovementWindow(result SelfImprovementWindowSymbolResult) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%d|%d|%t|%t|%t|%t|%t|%t|%s|%t|%t",
		result.Status,
		result.MissingStage,
		result.SourceDigest,
		result.IRDigest,
		result.SymbolName,
		result.SymbolKind,
		result.SymbolDigest,
		result.BaselineReceiptDigest,
		result.CandidateReceiptDigest,
		result.ComparisonSignal,
		result.WindowDigest,
		result.ByteDelta,
		result.LineDelta,
		result.SourceStable,
		result.CandidateStable,
		result.IRChangeStable,
		result.StructureStable,
		result.ExactSourceStable,
		result.GeneratedIRStable,
		result.NonExecuting,
		result.NonAuthorizing,
	))
}