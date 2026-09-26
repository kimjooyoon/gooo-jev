package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// SelfImprovementProvenanceDecisionSymbolResult is a read-only LSP
// projection of the next bounded provenance observation decision.
type SelfImprovementProvenanceDecisionSymbolResult struct {
	Status                            string
	MissingStage                      string
	SourceDigest                      string
	IRDigest                          string
	SymbolName                        string
	SymbolKind                        SymbolKind
	SymbolDigest                      string
	CurrentLifecycleObservationDigest string
	HistoryDigest                     string
	ObservationCount                  int
	StableCount                       int
	TransitionedCount                 int
	MetricsChangeCount                int
	GenerationChangeCount             int
	HistorySignal                     string
	DecisionSignal                    string
	DecisionReason                    string
	RequiresObservation               bool
	RequiresInspection                bool
	RequiresMeasurement               bool
	ResultDigest                      string
	NonExecuting                      bool
	NonAuthorizing                    bool
}

// ExplainSelfImprovementProvenanceDecision resolves one declaration against
// current lifecycle, provenance history, and bounded decision evidence.
func ExplainSelfImprovementProvenanceDecision(
	source, symbolName string,
	current RevisionSelfImprovementLifecycleObservation,
	history RevisionSelfImprovementProvenanceHistoryObservation,
	decision RevisionSelfImprovementProvenanceDecisionObservation,
) (SelfImprovementProvenanceDecisionSymbolResult, error) {
	result := SelfImprovementProvenanceDecisionSymbolResult{
		Status:                            "UNKNOWN",
		MissingStage:                      "lsp-self-improvement-provenance-decision",
		SourceDigest:                      digestString(source),
		CurrentLifecycleObservationDigest: current.ObservationDigest,
		HistoryDigest:                     decision.HistoryDigest,
		ObservationCount:                  decision.ObservationCount,
		StableCount:                       decision.StableCount,
		TransitionedCount:                 decision.TransitionedCount,
		MetricsChangeCount:                decision.MetricsChangeCount,
		GenerationChangeCount:             decision.GenerationChangeCount,
		HistorySignal:                     decision.HistorySignal,
		DecisionSignal:                    decision.DecisionSignal,
		DecisionReason:                    decision.DecisionReason,
		RequiresObservation:               decision.RequiresObservation,
		RequiresInspection:                decision.RequiresInspection,
		RequiresMeasurement:               decision.RequiresMeasurement,
		NonExecuting:                      true,
		NonAuthorizing:                    true,
	}
	setResultDigest := func() {
		result.ResultDigest = digestLSPSelfImprovementProvenanceDecision(result)
	}
	setResultDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-decision-snapshot"
		setResultDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := current.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-decision-current"
		setResultDigest()
		return result, fmt.Errorf("current self-improvement lifecycle is not valid: %w", err)
	}
	if err := history.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-decision-history"
		setResultDigest()
		return result, fmt.Errorf("self-improvement provenance history is not valid: %w", err)
	}
	if err := decision.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-decision-decision"
		setResultDigest()
		return result, fmt.Errorf("self-improvement provenance decision is not valid: %w", err)
	}
	if snapshot.SourceDigest != current.SourceDigest {
		result.MissingStage = "lsp-self-improvement-provenance-decision-source-link"
		setResultDigest()
		return result, fmt.Errorf("LSP source digest does not match current lifecycle")
	}
	if history.LastCurrentLifecycleDigest != current.ObservationDigest ||
		decision.HistoryDigest != history.HistoryDigest {
		result.MissingStage = "lsp-self-improvement-provenance-decision-link"
		setResultDigest()
		return result, fmt.Errorf("provenance decision is not linked to current history")
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
				result.MissingStage = "lsp-self-improvement-provenance-decision"
				setResultDigest()
				return result, fmt.Errorf("LSP provenance decision result is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-self-improvement-provenance-decision-symbol"
	setResultDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r SelfImprovementProvenanceDecisionSymbolResult) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("LSP provenance decision status is empty")
	}
	if r.Status == "BOUND" && r.MissingStage != "" {
		return fmt.Errorf("bound LSP provenance decision has a missing stage")
	}
	if r.Status == "UNKNOWN" && r.MissingStage == "" {
		return fmt.Errorf("unknown LSP provenance decision has no missing stage")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("LSP provenance decision must remain non-executing and non-authorizing")
	}
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" ||
		r.SymbolName == "" || r.SymbolDigest == "") {
		return fmt.Errorf("bound LSP provenance decision is incomplete")
	}
	if r.ObservationCount < 1 || r.StableCount < 0 || r.TransitionedCount < 0 ||
		r.StableCount+r.TransitionedCount != r.ObservationCount {
		return fmt.Errorf("LSP provenance decision counts are invalid")
	}
	if r.HistorySignal != "stable" && r.HistorySignal != "transitioned" &&
		r.HistorySignal != "mixed" {
		return fmt.Errorf("LSP provenance decision history signal is invalid")
	}
	if r.DecisionSignal != "observe" && r.DecisionSignal != "inspect" {
		return fmt.Errorf("LSP provenance decision signal is invalid")
	}
	if !r.RequiresObservation {
		return fmt.Errorf("LSP provenance decision must require observation")
	}
	if r.HistorySignal == "stable" && (r.DecisionSignal != "observe" ||
		r.RequiresInspection || r.RequiresMeasurement) {
		return fmt.Errorf("stable LSP provenance decision does not match its signal")
	}
	if (r.HistorySignal == "transitioned" || r.HistorySignal == "mixed") &&
		(r.DecisionSignal != "inspect" || !r.RequiresInspection) {
		return fmt.Errorf("transitioned LSP provenance decision does not match its signal")
	}
	if digestLSPSelfImprovementProvenanceDecision(r) != r.ResultDigest {
		return fmt.Errorf("LSP provenance decision digest does not match its fields")
	}
	return nil
}

func digestLSPSelfImprovementProvenanceDecision(result SelfImprovementProvenanceDecisionSymbolResult) string {
	parts := []string{
		result.Status,
		result.MissingStage,
		result.SourceDigest,
		result.IRDigest,
		result.SymbolName,
		string(result.SymbolKind),
		result.SymbolDigest,
		result.CurrentLifecycleObservationDigest,
		result.HistoryDigest,
		strconv.Itoa(result.ObservationCount),
		strconv.Itoa(result.StableCount),
		strconv.Itoa(result.TransitionedCount),
		strconv.Itoa(result.MetricsChangeCount),
		strconv.Itoa(result.GenerationChangeCount),
		result.HistorySignal,
		result.DecisionSignal,
		result.DecisionReason,
		strconv.FormatBool(result.RequiresObservation),
		strconv.FormatBool(result.RequiresInspection),
		strconv.FormatBool(result.RequiresMeasurement),
		strconv.FormatBool(result.NonExecuting),
		strconv.FormatBool(result.NonAuthorizing),
	}
	return digestString(strings.Join(parts, "|"))
}
