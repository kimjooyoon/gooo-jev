package gooo

import "fmt"

// RevisionImprovementSymbolResult exposes a bounded baseline/candidate
// observation at one declaration without executing or authorizing a revision.
type RevisionImprovementSymbolResult struct {
	Status                 string
	MissingStage           string
	SourceDigest           string
	IRDigest               string
	SymbolName             string
	SymbolKind             SymbolKind
	SymbolDigest            string
	BaselineSummaryDigest   string
	CandidateSummaryDigest  string
	ObservationDigest       string
	ObservationClass        string
	ReviewSignal            string
	SourceStable            bool
	ProposedIRChanged       bool
	GeneratedIRChanged      bool
	StructureStable         bool
	ResultDigest            string
	NonExecuting            bool
	NonAuthorizing          bool
}

// ExplainRevisionImprovement resolves one declaration and projects the exact
// baseline/candidate observation associated with it.
func ExplainRevisionImprovement(source, symbolName string, baseline, candidate RevisionEvidenceSummary, observation RevisionImprovementObservation) (RevisionImprovementSymbolResult, error) {
	result := RevisionImprovementSymbolResult{
		Status:                "UNKNOWN",
		MissingStage:          "lsp-revision-improvement",
		SourceDigest:          digestString(source),
		BaselineSummaryDigest: baseline.SummaryDigest,
		CandidateSummaryDigest: candidate.SummaryDigest,
		ObservationDigest:     observation.ObservationDigest,
		ObservationClass:      observation.ObservationClass,
		ReviewSignal:          observation.ReviewSignal,
		SourceStable:          observation.SourceStable,
		ProposedIRChanged:     observation.ProposedIRChanged,
		GeneratedIRChanged:    observation.GeneratedIRChanged,
		StructureStable:       observation.StructureStable,
		NonExecuting:          true,
		NonAuthorizing:        true,
	}
	setResultDigest := func() {
		result.ResultDigest = digestLSPRevisionImprovement(result)
	}
	setResultDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-revision-improvement-snapshot"
		setResultDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := baseline.Validate(); err != nil {
		result.MissingStage = "lsp-revision-improvement-baseline"
		setResultDigest()
		return result, fmt.Errorf("baseline revision evidence summary is not valid: %w", err)
	}
	if err := candidate.Validate(); err != nil {
		result.MissingStage = "lsp-revision-improvement-candidate"
		setResultDigest()
		return result, fmt.Errorf("candidate revision evidence summary is not valid: %w", err)
	}
	if err := observation.Validate(); err != nil {
		result.MissingStage = "lsp-revision-improvement-observation"
		setResultDigest()
		return result, fmt.Errorf("revision improvement observation is not valid: %w", err)
	}
	if observation.BaselineSummaryDigest != baseline.SummaryDigest ||
		observation.CandidateSummaryDigest != candidate.SummaryDigest ||
		snapshot.SourceDigest != candidate.SourceDigest {
		result.MissingStage = "lsp-revision-improvement-link"
		setResultDigest()
		return result, fmt.Errorf("LSP revision improvement evidence is not linked")
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
				result.MissingStage = "lsp-revision-improvement"
				setResultDigest()
				return result, fmt.Errorf("LSP revision improvement result is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-revision-improvement-symbol"
	setResultDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r RevisionImprovementSymbolResult) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("LSP revision improvement status is empty")
	}
	if r.Status == "BOUND" && r.MissingStage != "" {
		return fmt.Errorf("bound LSP revision improvement has a missing stage")
	}
	if r.Status == "UNKNOWN" && r.MissingStage == "" {
		return fmt.Errorf("unknown LSP revision improvement has no missing stage")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("LSP revision improvement must remain non-executing and non-authorizing")
	}
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" || r.SymbolName == "" || r.SymbolDigest == "" || r.ObservationClass == "" || r.ReviewSignal == "") {
		return fmt.Errorf("bound LSP revision improvement is incomplete")
	}
	if digestLSPRevisionImprovement(r) != r.ResultDigest {
		return fmt.Errorf("LSP revision improvement digest does not match its fields")
	}
	return nil
}

func digestLSPRevisionImprovement(result RevisionImprovementSymbolResult) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t|%t|%t|%t|%t",
		result.Status,
		result.MissingStage,
		result.SourceDigest,
		result.IRDigest,
		result.SymbolName,
		result.SymbolKind,
		result.SymbolDigest,
		result.BaselineSummaryDigest,
		result.CandidateSummaryDigest,
		result.ObservationDigest,
		result.ObservationClass,
		result.ReviewSignal,
		result.SourceStable,
		result.ProposedIRChanged,
		result.GeneratedIRChanged,
		result.StructureStable,
		result.NonExecuting,
		result.NonAuthorizing,
	))
}