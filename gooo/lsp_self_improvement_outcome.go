package gooo

import "fmt"

// SelfImprovementOutcomeSymbolResult is a read-only LSP projection of
// application, metrics, and generation evidence for one declaration.
type SelfImprovementOutcomeSymbolResult struct {
	Status                       string
	MissingStage                 string
	SourceDigest                 string
	IRDigest                     string
	SymbolName                   string
	SymbolKind                   SymbolKind
	SymbolDigest                 string
	ApplicationObservationDigest string
	ApplicationDigest            string
	MetricsBindingDigest          string
	MetricsDigest                string
	ChangedByteCount             int
	ChangedLineCount             int
	IRChanged                    bool
	GeneratedSourceDigest        string
	GeneratedIRDigest            string
	StructureDigest              string
	ExactSourceMatch              bool
	StructureMatch                bool
	ResultDigest                 string
	NonExecuting                 bool
	NonAuthorizing               bool
}

// ExplainSelfImprovementOutcome resolves one declaration against bound
// application, metrics, and generation evidence without executing it.
func ExplainSelfImprovementOutcome(source, symbolName string, application RevisionSelfImprovementApplicationObservation, metrics RevisionMetricsBinding, generation RevisionGenerationAssessment) (SelfImprovementOutcomeSymbolResult, error) {
	result := SelfImprovementOutcomeSymbolResult{
		Status:                       "UNKNOWN",
		MissingStage:                 "lsp-self-improvement-outcome",
		SourceDigest:                 digestString(source),
		ApplicationObservationDigest: application.ObservationDigest,
		ApplicationDigest:             application.ApplicationDigest,
		MetricsBindingDigest:          metrics.BindingDigest,
		MetricsDigest:                metrics.MetricsDigest,
		ChangedByteCount:             metrics.ChangedByteCount,
		ChangedLineCount:             metrics.ChangedLineCount,
		IRChanged:                    metrics.IRChanged,
		GeneratedSourceDigest:        generation.GeneratedSourceDigest,
		GeneratedIRDigest:             generation.GeneratedIRDigest,
		StructureDigest:              generation.StructureDigest,
		ExactSourceMatch:              generation.ExactSourceMatch,
		StructureMatch:               generation.StructureMatch,
		NonExecuting:                 true,
		NonAuthorizing:               true,
	}
	setResultDigest := func() {
		result.ResultDigest = digestLSPSelfImprovementOutcome(result)
	}
	setResultDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-outcome-snapshot"
		setResultDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := application.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-outcome-application"
		setResultDigest()
		return result, fmt.Errorf("self-improvement application observation is not valid: %w", err)
	}
	if err := metrics.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-outcome-metrics"
		setResultDigest()
		return result, fmt.Errorf("revision metrics binding is not valid: %w", err)
	}
	if err := generation.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-outcome-generation"
		setResultDigest()
		return result, fmt.Errorf("revision generation assessment is not valid: %w", err)
	}
	if snapshot.SourceDigest != application.SourceDigest {
		result.MissingStage = "lsp-self-improvement-outcome-source-link"
		setResultDigest()
		return result, fmt.Errorf("LSP source digest does not match self-improvement application")
	}
	if application.ReceiptDigest != metrics.ReceiptDigest ||
		application.ApplicationDigest != metrics.ApplicationDigest ||
		application.SourceDigest != metrics.SourceDigest ||
		application.ProposedSourceDigest != metrics.ProposedSourceDigest ||
		application.InputIRDigest != metrics.InputIRDigest ||
		application.ProposedIRDigest != metrics.ProposedIRDigest ||
		application.CandidateDigest != metrics.CandidateDigest ||
		application.EditDigest != metrics.EditDigest {
		result.MissingStage = "lsp-self-improvement-outcome-application-metrics-link"
		setResultDigest()
		return result, fmt.Errorf("application observation and metrics binding are not linked")
	}
	if generation.ApplicationDigest != application.ApplicationDigest ||
		generation.SourceDigest != application.SourceDigest ||
		generation.ProposedSourceDigest != application.ProposedSourceDigest ||
		generation.GenerationSourceDigest != application.ProposedSourceDigest {
		result.MissingStage = "lsp-self-improvement-outcome-generation-link"
		setResultDigest()
		return result, fmt.Errorf("generation assessment is not linked to the application source")
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
				result.MissingStage = "lsp-self-improvement-outcome"
				setResultDigest()
				return result, fmt.Errorf("LSP self-improvement outcome is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-self-improvement-outcome-symbol"
	setResultDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r SelfImprovementOutcomeSymbolResult) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("LSP self-improvement outcome status is empty")
	}
	if r.Status == "BOUND" && r.MissingStage != "" {
		return fmt.Errorf("bound LSP self-improvement outcome has a missing stage")
	}
	if r.Status == "UNKNOWN" && r.MissingStage == "" {
		return fmt.Errorf("unknown LSP self-improvement outcome has no missing stage")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("LSP self-improvement outcome must remain non-executing and non-authorizing")
	}
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" ||
		r.SymbolName == "" || r.SymbolDigest == "") {
		return fmt.Errorf("bound LSP self-improvement outcome is incomplete")
	}
	if r.ChangedByteCount < 0 || r.ChangedLineCount < 0 {
		return fmt.Errorf("LSP self-improvement outcome counts must be non-negative")
	}
	if digestLSPSelfImprovementOutcome(r) != r.ResultDigest {
		return fmt.Errorf("LSP self-improvement outcome digest does not match its fields")
	}
	return nil
}

func digestLSPSelfImprovementOutcome(result SelfImprovementOutcomeSymbolResult) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%d|%d|%t|%t|%t|%t|%t",
		result.Status,
		result.MissingStage,
		result.SourceDigest,
		result.IRDigest,
		result.SymbolName,
		result.SymbolKind,
		result.SymbolDigest,
		result.ApplicationObservationDigest,
		result.ApplicationDigest,
		result.MetricsBindingDigest,
		result.MetricsDigest,
		result.GeneratedSourceDigest,
		result.GeneratedIRDigest,
		result.StructureDigest,
		result.ChangedByteCount,
		result.ChangedLineCount,
		result.IRChanged,
		result.ExactSourceMatch,
		result.StructureMatch,
		result.NonExecuting,
		result.NonAuthorizing,
	))
}
