package gooo

import "fmt"

// RevisionEvidenceSymbolResult is a source- and symbol-bound LSP projection
// of revision provenance. It cannot execute or authorize a revision.
type RevisionEvidenceSymbolResult struct {
	Status                  string
	MissingStage            string
	SourceDigest            string
	IRDigest                string
	SymbolName              string
	SymbolKind              SymbolKind
	SymbolDigest            string
	SummaryDigest            string
	ChainDigest              string
	GenerationBindingDigest  string
	EvidenceStageCount       int
	ExactSourceMatch         bool
	StructureMatch           bool
	ReverseObserved          bool
	ResultDigest             string
	NonExecuting             bool
	NonAuthorizing           bool
}

// ExplainRevisionEvidence resolves one declaration from a parsed source and
// returns the exact revision evidence summary bound to that source.
func ExplainRevisionEvidence(source, symbolName string, summary RevisionEvidenceSummary) (RevisionEvidenceSymbolResult, error) {
	result := RevisionEvidenceSymbolResult{
		Status:                 "UNKNOWN",
		MissingStage:           "lsp-revision-evidence",
		SourceDigest:           digestString(source),
		SummaryDigest:          summary.SummaryDigest,
		ChainDigest:            summary.ChainDigest,
		GenerationBindingDigest: summary.GenerationBindingDigest,
		EvidenceStageCount:     summary.EvidenceStageCount,
		ExactSourceMatch:       summary.ExactSourceMatch,
		StructureMatch:         summary.StructureMatch,
		ReverseObserved:        summary.ReverseObserved,
		NonExecuting:           true,
		NonAuthorizing:         true,
	}
	setResultDigest := func() {
		result.ResultDigest = digestRevisionEvidenceSymbolResult(result)
	}
	setResultDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-revision-evidence-snapshot"
		setResultDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := summary.Validate(); err != nil {
		result.MissingStage = "lsp-revision-evidence-summary"
		setResultDigest()
		return result, fmt.Errorf("revision evidence summary is not valid: %w", err)
	}
	if snapshot.SourceDigest != summary.SourceDigest {
		result.MissingStage = "lsp-revision-evidence-link"
		setResultDigest()
		return result, fmt.Errorf("LSP source digest does not match revision evidence summary")
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
				result.MissingStage = "lsp-revision-evidence"
				setResultDigest()
				return result, fmt.Errorf("LSP revision evidence result is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-revision-evidence-symbol"
	setResultDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r RevisionEvidenceSymbolResult) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("LSP revision evidence status is empty")
	}
	if r.Status == "BOUND" && r.MissingStage != "" {
		return fmt.Errorf("bound LSP revision evidence has a missing stage")
	}
	if r.Status == "UNKNOWN" && r.MissingStage == "" {
		return fmt.Errorf("unknown LSP revision evidence has no missing stage")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("LSP revision evidence must remain non-executing and non-authorizing")
	}
	if r.EvidenceStageCount != 6 {
		return fmt.Errorf("LSP revision evidence stage count must be 6")
	}
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" || r.SymbolName == "" || r.SymbolDigest == "") {
		return fmt.Errorf("bound LSP revision evidence is incomplete")
	}
	if digestRevisionEvidenceSymbolResult(r) != r.ResultDigest {
		return fmt.Errorf("LSP revision evidence digest does not match its fields")
	}
	return nil
}

func digestRevisionEvidenceSymbolResult(result RevisionEvidenceSymbolResult) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%d|%t|%t|%t|%t|%t|%s",
		result.Status,
		result.MissingStage,
		result.SourceDigest,
		result.IRDigest,
		result.SymbolName,
		result.SymbolKind,
		result.SymbolDigest,
		result.SummaryDigest,
		result.ChainDigest,
		result.EvidenceStageCount,
		result.ExactSourceMatch,
		result.StructureMatch,
		result.ReverseObserved,
		result.NonExecuting,
		result.NonAuthorizing,
		result.ResultDigest,
	))
}