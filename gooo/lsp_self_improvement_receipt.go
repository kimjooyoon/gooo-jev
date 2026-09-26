package gooo

import "fmt"

// SelfImprovementReceiptSymbolResult is a declaration-scoped, read-only LSP
// projection of the composite self-improvement receipt.
type SelfImprovementReceiptSymbolResult struct {
	Status              string
	MissingStage        string
	SourceDigest        string
	IRDigest            string
	SymbolName          string
	SymbolKind          SymbolKind
	SymbolDigest         string
	ReceiptDigest        string
	ApplicationDigest    string
	CandidateDigest      string
	MetricsDigest        string
	GeneratedIRDigest    string
	StructureDigest      string
	StageCount           int
	IRChanged            bool
	ExactSourceMatch     bool
	StructureMatch       bool
	ResultDigest         string
	NonExecuting         bool
	NonAuthorizing       bool
}

// ExplainSelfImprovement resolves one declaration against a validated
// RevisionSelfImprovementReceipt without executing or authorizing it.
func ExplainSelfImprovement(source, symbolName string, receipt RevisionSelfImprovementReceipt) (SelfImprovementReceiptSymbolResult, error) {
	result := SelfImprovementReceiptSymbolResult{
		Status:            "UNKNOWN",
		MissingStage:      "lsp-self-improvement-receipt",
		SourceDigest:      digestString(source),
		ReceiptDigest:     receipt.ReceiptDigest,
		ApplicationDigest: receipt.ApplicationDigest,
		CandidateDigest:   receipt.CandidateDigest,
		MetricsDigest:     receipt.MetricsDigest,
		GeneratedIRDigest: receipt.GeneratedIRDigest,
		StructureDigest:   receipt.StructureDigest,
		StageCount:        receipt.StageCount,
		IRChanged:         receipt.IRChanged,
		ExactSourceMatch:  receipt.ExactSourceMatch,
		StructureMatch:    receipt.StructureMatch,
		NonExecuting:      true,
		NonAuthorizing:    true,
	}
	setResultDigest := func() {
		result.ResultDigest = digestLSPSelfImprovementReceipt(result)
	}
	setResultDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-receipt-snapshot"
		setResultDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := receipt.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-receipt-receipt"
		setResultDigest()
		return result, fmt.Errorf("revision self-improvement receipt is not valid: %w", err)
	}
	if snapshot.SourceDigest != receipt.SourceDigest {
		result.MissingStage = "lsp-self-improvement-receipt-link"
		setResultDigest()
		return result, fmt.Errorf("LSP source digest does not match self-improvement receipt")
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
				result.MissingStage = "lsp-self-improvement-receipt"
				setResultDigest()
				return result, fmt.Errorf("LSP self-improvement receipt result is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-self-improvement-receipt-symbol"
	setResultDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r SelfImprovementReceiptSymbolResult) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("LSP self-improvement receipt status is empty")
	}
	if r.Status == "BOUND" && r.MissingStage != "" {
		return fmt.Errorf("bound LSP self-improvement receipt has a missing stage")
	}
	if r.Status == "UNKNOWN" && r.MissingStage == "" {
		return fmt.Errorf("unknown LSP self-improvement receipt has no missing stage")
	}
	if r.StageCount != 4 {
		return fmt.Errorf("LSP self-improvement receipt stage count must be 4")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("LSP self-improvement receipt must remain non-executing and non-authorizing")
	}
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" || r.SymbolName == "" || r.SymbolDigest == "") {
		return fmt.Errorf("bound LSP self-improvement receipt is incomplete")
	}
	if digestLSPSelfImprovementReceipt(r) != r.ResultDigest {
		return fmt.Errorf("LSP self-improvement receipt digest does not match its fields")
	}
	return nil
}

func digestLSPSelfImprovementReceipt(result SelfImprovementReceiptSymbolResult) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%d|%t|%t|%t|%t|%t",
		result.Status,
		result.MissingStage,
		result.SourceDigest,
		result.IRDigest,
		result.SymbolName,
		result.SymbolKind,
		result.SymbolDigest,
		result.ReceiptDigest,
		result.ApplicationDigest,
		result.CandidateDigest,
		result.MetricsDigest,
		result.GeneratedIRDigest,
		result.StructureDigest,
		result.StageCount,
		result.IRChanged,
		result.ExactSourceMatch,
		result.StructureMatch,
		result.NonExecuting,
		result.NonAuthorizing,
	))
}