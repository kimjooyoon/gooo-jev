package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// SelfImprovementProvenanceTransitionSymbolResult is a read-only LSP
// projection of provenance stability or change between two lifecycles.
type SelfImprovementProvenanceTransitionSymbolResult struct {
	Status                             string
	MissingStage                       string
	SourceDigest                       string
	IRDigest                           string
	SymbolName                         string
	SymbolKind                         SymbolKind
	SymbolDigest                       string
	PreviousLifecycleObservationDigest string
	CurrentLifecycleObservationDigest  string
	PreviousSourceDigest               string
	CurrentSourceDigest                string
	PreviousProposedSourceDigest       string
	CurrentProposedSourceDigest        string
	PreviousInputIRDigest              string
	CurrentInputIRDigest               string
	PreviousProposedIRDigest            string
	CurrentProposedIRDigest             string
	PreviousGeneratedSourceDigest       string
	CurrentGeneratedSourceDigest        string
	PreviousGeneratedIRDigest           string
	CurrentGeneratedIRDigest            string
	PreviousStructureDigest             string
	CurrentStructureDigest              string
	PreviousMetricsBindingDigest        string
	CurrentMetricsBindingDigest         string
	PreviousGenerationAssessmentDigest  string
	CurrentGenerationAssessmentDigest   string
	PreviousChangedByteCount            int
	CurrentChangedByteCount             int
	PreviousChangedLineCount            int
	CurrentChangedLineCount             int
	PreviousIRChanged                   bool
	CurrentIRChanged                    bool
	PreviousExactSourceMatch            bool
	CurrentExactSourceMatch             bool
	PreviousStructureMatch              bool
	CurrentStructureMatch               bool
	SourceChanged                       bool
	ProposedSourceChanged               bool
	InputIRChanged                      bool
	ProposedIRChanged                   bool
	GeneratedSourceChanged              bool
	GeneratedIRChanged                 bool
	StructureChanged                   bool
	MetricsChanged                     bool
	GenerationChanged                  bool
	TransitionSignal                   string
	ResultDigest                       string
	NonExecuting                       bool
	NonAuthorizing                     bool
}

// ExplainSelfImprovementProvenanceTransition resolves one declaration against
// two validated lifecycles and their transition evidence.
func ExplainSelfImprovementProvenanceTransition(
	source, symbolName string,
	previous RevisionSelfImprovementLifecycleObservation,
	current RevisionSelfImprovementLifecycleObservation,
	transition RevisionSelfImprovementProvenanceTransitionObservation,
) (SelfImprovementProvenanceTransitionSymbolResult, error) {
	result := SelfImprovementProvenanceTransitionSymbolResult{
		Status:                             "UNKNOWN",
		MissingStage:                       "lsp-self-improvement-provenance-transition",
		SourceDigest:                       digestString(source),
		PreviousLifecycleObservationDigest: transition.PreviousLifecycleObservationDigest,
		CurrentLifecycleObservationDigest:  transition.CurrentLifecycleObservationDigest,
		PreviousSourceDigest:               transition.PreviousSourceDigest,
		CurrentSourceDigest:                transition.CurrentSourceDigest,
		PreviousProposedSourceDigest:       transition.PreviousProposedSourceDigest,
		CurrentProposedSourceDigest:         transition.CurrentProposedSourceDigest,
		PreviousInputIRDigest:              transition.PreviousInputIRDigest,
		CurrentInputIRDigest:               transition.CurrentInputIRDigest,
		PreviousProposedIRDigest:            transition.PreviousProposedIRDigest,
		CurrentProposedIRDigest:             transition.CurrentProposedIRDigest,
		PreviousGeneratedSourceDigest:       transition.PreviousGeneratedSourceDigest,
		CurrentGeneratedSourceDigest:        transition.CurrentGeneratedSourceDigest,
		PreviousGeneratedIRDigest:           transition.PreviousGeneratedIRDigest,
		CurrentGeneratedIRDigest:            transition.CurrentGeneratedIRDigest,
		PreviousStructureDigest:             transition.PreviousStructureDigest,
		CurrentStructureDigest:              transition.CurrentStructureDigest,
		PreviousMetricsBindingDigest:        transition.PreviousMetricsBindingDigest,
		CurrentMetricsBindingDigest:         transition.CurrentMetricsBindingDigest,
		PreviousGenerationAssessmentDigest:  transition.PreviousGenerationAssessmentDigest,
		CurrentGenerationAssessmentDigest:   transition.CurrentGenerationAssessmentDigest,
		PreviousChangedByteCount:            transition.PreviousChangedByteCount,
		CurrentChangedByteCount:             transition.CurrentChangedByteCount,
		PreviousChangedLineCount:            transition.PreviousChangedLineCount,
		CurrentChangedLineCount:             transition.CurrentChangedLineCount,
		PreviousIRChanged:                   transition.PreviousIRChanged,
		CurrentIRChanged:                    transition.CurrentIRChanged,
		PreviousExactSourceMatch:            transition.PreviousExactSourceMatch,
		CurrentExactSourceMatch:             transition.CurrentExactSourceMatch,
		PreviousStructureMatch:              transition.PreviousStructureMatch,
		CurrentStructureMatch:               transition.CurrentStructureMatch,
		SourceChanged:                       transition.SourceChanged,
		ProposedSourceChanged:               transition.ProposedSourceChanged,
		InputIRChanged:                      transition.InputIRChanged,
		ProposedIRChanged:                   transition.ProposedIRChanged,
		GeneratedSourceChanged:              transition.GeneratedSourceChanged,
		GeneratedIRChanged:                  transition.GeneratedIRChanged,
		StructureChanged:                    transition.StructureChanged,
		MetricsChanged:                      transition.MetricsChanged,
		GenerationChanged:                  transition.GenerationChanged,
		TransitionSignal:                    transition.TransitionSignal,
		NonExecuting:                        true,
		NonAuthorizing:                      true,
	}
	setResultDigest := func() {
		result.ResultDigest = digestLSPSelfImprovementProvenanceTransition(result)
	}
	setResultDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-transition-snapshot"
		setResultDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := previous.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-transition-previous"
		setResultDigest()
		return result, fmt.Errorf("previous self-improvement lifecycle is not valid: %w", err)
	}
	if err := current.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-transition-current"
		setResultDigest()
		return result, fmt.Errorf("current self-improvement lifecycle is not valid: %w", err)
	}
	if err := transition.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-transition-transition"
		setResultDigest()
		return result, fmt.Errorf("self-improvement provenance transition is not valid: %w", err)
	}
	if transition.PreviousLifecycleObservationDigest != previous.ObservationDigest ||
		transition.CurrentLifecycleObservationDigest != current.ObservationDigest {
		result.MissingStage = "lsp-self-improvement-provenance-transition-link"
		setResultDigest()
		return result, fmt.Errorf("transition is not linked to both lifecycle observations")
	}
	if snapshot.SourceDigest != current.SourceDigest {
		result.MissingStage = "lsp-self-improvement-provenance-transition-source-link"
		setResultDigest()
		return result, fmt.Errorf("LSP source digest does not match current lifecycle")
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
				result.MissingStage = "lsp-self-improvement-provenance-transition"
				setResultDigest()
				return result, fmt.Errorf("LSP provenance transition result is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-self-improvement-provenance-transition-symbol"
	setResultDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r SelfImprovementProvenanceTransitionSymbolResult) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("LSP provenance transition status is empty")
	}
	if r.Status == "BOUND" && r.MissingStage != "" {
		return fmt.Errorf("bound LSP provenance transition has a missing stage")
	}
	if r.Status == "UNKNOWN" && r.MissingStage == "" {
		return fmt.Errorf("unknown LSP provenance transition has no missing stage")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("LSP provenance transition must remain non-executing and non-authorizing")
	}
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" ||
		r.SymbolName == "" || r.SymbolDigest == "") {
		return fmt.Errorf("bound LSP provenance transition is incomplete")
	}
	if r.TransitionSignal == "" {
		return fmt.Errorf("LSP provenance transition signal is empty")
	}
	if digestLSPSelfImprovementProvenanceTransition(r) != r.ResultDigest {
		return fmt.Errorf("LSP provenance transition digest does not match its fields")
	}
	return nil
}

func digestLSPSelfImprovementProvenanceTransition(result SelfImprovementProvenanceTransitionSymbolResult) string {
	parts := []string{
		result.Status,
		result.MissingStage,
		result.SourceDigest,
		result.IRDigest,
		result.SymbolName,
		string(result.SymbolKind),
		result.SymbolDigest,
		result.PreviousLifecycleObservationDigest,
		result.CurrentLifecycleObservationDigest,
		result.PreviousSourceDigest,
		result.CurrentSourceDigest,
		result.PreviousProposedSourceDigest,
		result.CurrentProposedSourceDigest,
		result.PreviousInputIRDigest,
		result.CurrentInputIRDigest,
		result.PreviousProposedIRDigest,
		result.CurrentProposedIRDigest,
		result.PreviousGeneratedSourceDigest,
		result.CurrentGeneratedSourceDigest,
		result.PreviousGeneratedIRDigest,
		result.CurrentGeneratedIRDigest,
		result.PreviousStructureDigest,
		result.CurrentStructureDigest,
		result.PreviousMetricsBindingDigest,
		result.CurrentMetricsBindingDigest,
		result.PreviousGenerationAssessmentDigest,
		result.CurrentGenerationAssessmentDigest,
		strconv.Itoa(result.PreviousChangedByteCount),
		strconv.Itoa(result.CurrentChangedByteCount),
		strconv.Itoa(result.PreviousChangedLineCount),
		strconv.Itoa(result.CurrentChangedLineCount),
		strconv.FormatBool(result.PreviousIRChanged),
		strconv.FormatBool(result.CurrentIRChanged),
		strconv.FormatBool(result.PreviousExactSourceMatch),
		strconv.FormatBool(result.CurrentExactSourceMatch),
		strconv.FormatBool(result.PreviousStructureMatch),
		strconv.FormatBool(result.CurrentStructureMatch),
		strconv.FormatBool(result.SourceChanged),
		strconv.FormatBool(result.ProposedSourceChanged),
		strconv.FormatBool(result.InputIRChanged),
		strconv.FormatBool(result.ProposedIRChanged),
		strconv.FormatBool(result.GeneratedSourceChanged),
		strconv.FormatBool(result.GeneratedIRChanged),
		strconv.FormatBool(result.StructureChanged),
		strconv.FormatBool(result.MetricsChanged),
		strconv.FormatBool(result.GenerationChanged),
		result.TransitionSignal,
		strconv.FormatBool(result.NonExecuting),
		strconv.FormatBool(result.NonAuthorizing),
	}
	return digestString(strings.Join(parts, "|"))
}
