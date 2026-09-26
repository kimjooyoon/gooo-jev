package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// SelfImprovementReverseSymbolResult is a read-only LSP projection of
// reverse-generated source, IR, and structure evidence.
type SelfImprovementReverseSymbolResult struct {
	Status                         string
	MissingStage                   string
	SourceDigest                   string
	IRDigest                       string
	SymbolName                     string
	SymbolKind                     SymbolKind
	SymbolDigest                   string
	LifecycleObservationDigest     string
	MetricsBindingDigest           string
	GenerationAssessmentDigest     string
	ProposedSourceDigest           string
	InputIRDigest                  string
	ProposedIRDigest               string
	GeneratedSourceDigest          string
	GeneratedIRDigest              string
	StructureDigest                string
	ObservedGeneratedSourceDigest  string
	ObservedGeneratedIRDigest      string
	ObservedStructureDigest         string
	ExactSourceMatch                bool
	ExactIRMatch                    bool
	ExactStructureMatch             bool
	ReverseSignal                   string
	ResultDigest                    string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

// ExplainSelfImprovementReverseObservation resolves one declaration against
// reverse provenance without executing or authorizing a change.
func ExplainSelfImprovementReverseObservation(
	source, symbolName string,
	lifecycle RevisionSelfImprovementLifecycleObservation,
	reverse RevisionSelfImprovementReverseObservation,
) (SelfImprovementReverseSymbolResult, error) {
	result := SelfImprovementReverseSymbolResult{
		Status:                         "UNKNOWN",
		MissingStage:                   "lsp-self-improvement-reverse-observation",
		SourceDigest:                   digestString(source),
		LifecycleObservationDigest:     reverse.LifecycleObservationDigest,
		MetricsBindingDigest:           reverse.MetricsBindingDigest,
		GenerationAssessmentDigest:    reverse.GenerationAssessmentDigest,
		ProposedSourceDigest:           reverse.ProposedSourceDigest,
		InputIRDigest:                  reverse.InputIRDigest,
		ProposedIRDigest:               reverse.ProposedIRDigest,
		GeneratedSourceDigest:          reverse.GeneratedSourceDigest,
		GeneratedIRDigest:              reverse.GeneratedIRDigest,
		StructureDigest:                reverse.StructureDigest,
		ObservedGeneratedSourceDigest: reverse.ObservedGeneratedSourceDigest,
		ObservedGeneratedIRDigest:     reverse.ObservedGeneratedIRDigest,
		ObservedStructureDigest:        reverse.ObservedStructureDigest,
		ExactSourceMatch:               reverse.ExactSourceMatch,
		ExactIRMatch:                   reverse.ExactIRMatch,
		ExactStructureMatch:            reverse.ExactStructureMatch,
		ReverseSignal:                  reverse.ReverseSignal,
		NonExecuting:                   true,
		NonAuthorizing:                 true,
	}
	setResultDigest := func() {
		result.ResultDigest = digestLSPSelfImprovementReverseObservation(result)
	}
	setResultDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-reverse-snapshot"
		setResultDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := lifecycle.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-reverse-lifecycle"
		setResultDigest()
		return result, fmt.Errorf("self-improvement lifecycle is not valid: %w", err)
	}
	if err := reverse.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-reverse-reverse"
		setResultDigest()
		return result, fmt.Errorf("self-improvement reverse observation is not valid: %w", err)
	}
	if reverse.LifecycleObservationDigest != lifecycle.ObservationDigest ||
		reverse.SourceDigest != lifecycle.SourceDigest {
		result.MissingStage = "lsp-self-improvement-reverse-link"
		setResultDigest()
		return result, fmt.Errorf("reverse observation is not linked to lifecycle source")
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
				result.MissingStage = "lsp-self-improvement-reverse-observation"
				setResultDigest()
				return result, fmt.Errorf("LSP self-improvement reverse result is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-self-improvement-reverse-symbol"
	setResultDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r SelfImprovementReverseSymbolResult) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("LSP self-improvement reverse status is empty")
	}
	if r.Status == "BOUND" && r.MissingStage != "" {
		return fmt.Errorf("bound LSP self-improvement reverse result has a missing stage")
	}
	if r.Status == "UNKNOWN" && r.MissingStage == "" {
		return fmt.Errorf("unknown LSP self-improvement reverse result has no missing stage")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("LSP self-improvement reverse result must remain non-executing and non-authorizing")
	}
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" ||
		r.SymbolName == "" || r.SymbolDigest == "") {
		return fmt.Errorf("bound LSP self-improvement reverse result is incomplete")
	}
	if r.ReverseSignal == "" {
		return fmt.Errorf("LSP self-improvement reverse signal is empty")
	}
	if digestLSPSelfImprovementReverseObservation(r) != r.ResultDigest {
		return fmt.Errorf("LSP self-improvement reverse digest does not match its fields")
	}
	return nil
}

func digestLSPSelfImprovementReverseObservation(result SelfImprovementReverseSymbolResult) string {
	parts := []string{
		result.Status,
		result.MissingStage,
		result.SourceDigest,
		result.IRDigest,
		result.SymbolName,
		string(result.SymbolKind),
		result.SymbolDigest,
		result.LifecycleObservationDigest,
		result.MetricsBindingDigest,
		result.GenerationAssessmentDigest,
		result.ProposedSourceDigest,
		result.InputIRDigest,
		result.ProposedIRDigest,
		result.GeneratedSourceDigest,
		result.GeneratedIRDigest,
		result.StructureDigest,
		result.ObservedGeneratedSourceDigest,
		result.ObservedGeneratedIRDigest,
		result.ObservedStructureDigest,
		strconv.FormatBool(result.ExactSourceMatch),
		strconv.FormatBool(result.ExactIRMatch),
		strconv.FormatBool(result.ExactStructureMatch),
		result.ReverseSignal,
		strconv.FormatBool(result.NonExecuting),
		strconv.FormatBool(result.NonAuthorizing),
	}
	return digestString(strings.Join(parts, "|"))
}
