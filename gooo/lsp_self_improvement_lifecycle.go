package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// SelfImprovementLifecycleSymbolResult is a declaration-scoped, read-only LSP
// projection of the complete self-improvement lifecycle provenance.
type SelfImprovementLifecycleSymbolResult struct {
	Status                              string
	MissingStage                        string
	SourceDigest                         string
	IRDigest                             string
	SymbolName                           string
	SymbolKind                           SymbolKind
	SymbolDigest                         string
	ExecutionApplicationObservationDigest string
	OutcomeObservationDigest              string
	ApplicationObservationDigest          string
	MetricsBindingDigest                  string
	GenerationAssessmentDigest            string
	PlanObservationDigest                 string
	PlanDigest                            string
	ApplicationDigest                     string
	ReceiptDigest                         string
	ProposedSourceDigest                  string
	InputIRDigest                         string
	ProposedIRDigest                      string
	CandidateDigest                       string
	EditDigest                            string
	ChangedByteCount                      int
	ChangedLineCount                      int
	IRChanged                             bool
	GeneratedSourceDigest                 string
	GeneratedIRDigest                     string
	StructureDigest                       string
	ExactSourceMatch                      bool
	StructureMatch                        bool
	DecisionSignal                        string
	FeedbackSignal                       string
	ApplicationSignal                     string
	PlanObserved                          bool
	ApplicationObserved                   bool
	LifecycleSignal                       string
	ResultDigest                          string
	NonExecuting                          bool
	NonAuthorizing                        bool
}

// ExplainSelfImprovementLifecycle resolves one declaration against a
// validated lifecycle without executing or authorizing a change.
func ExplainSelfImprovementLifecycle(source, symbolName string, lifecycle RevisionSelfImprovementLifecycleObservation) (SelfImprovementLifecycleSymbolResult, error) {
	result := SelfImprovementLifecycleSymbolResult{
		Status:                                "UNKNOWN",
		MissingStage:                          "lsp-self-improvement-lifecycle",
		SourceDigest:                          digestString(source),
		ExecutionApplicationObservationDigest: lifecycle.ExecutionApplicationObservationDigest,
		OutcomeObservationDigest:              lifecycle.OutcomeObservationDigest,
		ApplicationObservationDigest:          lifecycle.ApplicationObservationDigest,
		MetricsBindingDigest:                  lifecycle.MetricsBindingDigest,
		GenerationAssessmentDigest:            lifecycle.GenerationAssessmentDigest,
		PlanObservationDigest:                 lifecycle.PlanObservationDigest,
		PlanDigest:                            lifecycle.PlanDigest,
		ApplicationDigest:                     lifecycle.ApplicationDigest,
		ReceiptDigest:                         lifecycle.ReceiptDigest,
		ProposedSourceDigest:                 lifecycle.ProposedSourceDigest,
		InputIRDigest:                         lifecycle.InputIRDigest,
		ProposedIRDigest:                      lifecycle.ProposedIRDigest,
		CandidateDigest:                       lifecycle.CandidateDigest,
		EditDigest:                            lifecycle.EditDigest,
		ChangedByteCount:                      lifecycle.ChangedByteCount,
		ChangedLineCount:                      lifecycle.ChangedLineCount,
		IRChanged:                             lifecycle.IRChanged,
		GeneratedSourceDigest:                 lifecycle.GeneratedSourceDigest,
		GeneratedIRDigest:                     lifecycle.GeneratedIRDigest,
		StructureDigest:                       lifecycle.StructureDigest,
		ExactSourceMatch:                      lifecycle.ExactSourceMatch,
		StructureMatch:                        lifecycle.StructureMatch,
		DecisionSignal:                        lifecycle.DecisionSignal,
		FeedbackSignal:                        lifecycle.FeedbackSignal,
		ApplicationSignal:                     lifecycle.ApplicationSignal,
		PlanObserved:                          lifecycle.PlanObserved,
		ApplicationObserved:                   lifecycle.ApplicationObserved,
		LifecycleSignal:                       lifecycle.LifecycleSignal,
		NonExecuting:                          true,
		NonAuthorizing:                        true,
	}
	setResultDigest := func() {
		result.ResultDigest = digestLSPSelfImprovementLifecycle(result)
	}
	setResultDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-lifecycle-snapshot"
		setResultDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := lifecycle.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-lifecycle-lifecycle"
		setResultDigest()
		return result, fmt.Errorf("self-improvement lifecycle is not valid: %w", err)
	}
	if snapshot.SourceDigest != lifecycle.SourceDigest {
		result.MissingStage = "lsp-self-improvement-lifecycle-link"
		setResultDigest()
		return result, fmt.Errorf("LSP source digest does not match self-improvement lifecycle")
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
				result.MissingStage = "lsp-self-improvement-lifecycle"
				setResultDigest()
				return result, fmt.Errorf("LSP self-improvement lifecycle result is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-self-improvement-lifecycle-symbol"
	setResultDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r SelfImprovementLifecycleSymbolResult) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("LSP self-improvement lifecycle status is empty")
	}
	if r.Status == "BOUND" && r.MissingStage != "" {
		return fmt.Errorf("bound LSP self-improvement lifecycle has a missing stage")
	}
	if r.Status == "UNKNOWN" && r.MissingStage == "" {
		return fmt.Errorf("unknown LSP self-improvement lifecycle has no missing stage")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("LSP self-improvement lifecycle must remain non-executing and non-authorizing")
	}
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" ||
		r.SymbolName == "" || r.SymbolDigest == "") {
		return fmt.Errorf("bound LSP self-improvement lifecycle is incomplete")
	}
	if r.ChangedByteCount < 0 || r.ChangedLineCount < 0 {
		return fmt.Errorf("LSP self-improvement lifecycle counts must be non-negative")
	}
	if r.DecisionSignal != "observe" && r.DecisionSignal != "remeasure" &&
		r.DecisionSignal != "review" && r.DecisionSignal != "inspect" {
		return fmt.Errorf("LSP self-improvement lifecycle decision signal is invalid")
	}
	if r.FeedbackSignal != r.DecisionSignal {
		return fmt.Errorf("LSP self-improvement lifecycle feedback signal is not linked")
	}
	if r.ApplicationSignal == "" || r.LifecycleSignal == "" {
		return fmt.Errorf("LSP self-improvement lifecycle signals are incomplete")
	}
	if !r.PlanObserved || !r.ApplicationObserved {
		return fmt.Errorf("LSP self-improvement lifecycle must preserve observations")
	}
	if digestLSPSelfImprovementLifecycle(r) != r.ResultDigest {
		return fmt.Errorf("LSP self-improvement lifecycle digest does not match its fields")
	}
	return nil
}

func digestLSPSelfImprovementLifecycle(result SelfImprovementLifecycleSymbolResult) string {
	parts := []string{
		result.Status,
		result.MissingStage,
		result.SourceDigest,
		result.IRDigest,
		result.SymbolName,
		string(result.SymbolKind),
		result.SymbolDigest,
		result.ExecutionApplicationObservationDigest,
		result.OutcomeObservationDigest,
		result.ApplicationObservationDigest,
		result.MetricsBindingDigest,
		result.GenerationAssessmentDigest,
		result.PlanObservationDigest,
		result.PlanDigest,
		result.ApplicationDigest,
		result.ReceiptDigest,
		result.ProposedSourceDigest,
		result.InputIRDigest,
		result.ProposedIRDigest,
		result.CandidateDigest,
		result.EditDigest,
		result.GeneratedSourceDigest,
		result.GeneratedIRDigest,
		result.StructureDigest,
		result.DecisionSignal,
		result.FeedbackSignal,
		result.ApplicationSignal,
		result.LifecycleSignal,
		strconv.Itoa(result.ChangedByteCount),
		strconv.Itoa(result.ChangedLineCount),
		strconv.FormatBool(result.IRChanged),
		strconv.FormatBool(result.ExactSourceMatch),
		strconv.FormatBool(result.StructureMatch),
		strconv.FormatBool(result.PlanObserved),
		strconv.FormatBool(result.ApplicationObserved),
		strconv.FormatBool(result.NonExecuting),
		strconv.FormatBool(result.NonAuthorizing),
	}
	return digestString(strings.Join(parts, "|"))
}
