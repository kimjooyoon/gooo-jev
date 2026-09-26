package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// SelfImprovementProvenanceHistorySymbolResult is a read-only LSP projection
// of ordered lifecycle provenance history for one declaration.
type SelfImprovementProvenanceHistorySymbolResult struct {
	Status                              string
	MissingStage                        string
	SourceDigest                        string
	IRDigest                            string
	SymbolName                          string
	SymbolKind                          SymbolKind
	SymbolDigest                        string
	CurrentLifecycleObservationDigest   string
	HistoryDigest                       string
	ObservationCount                   int
	TransitionDigests                  []string
	FirstTransitionDigest              string
	LastTransitionDigest               string
	FirstPreviousLifecycleDigest       string
	LastCurrentLifecycleDigest         string
	StableCount                        int
	TransitionedCount                  int
	MetricsChangeCount                 int
	GenerationChangeCount              int
	HistorySignal                      string
	ResultDigest                       string
	NonExecuting                      bool
	NonAuthorizing                    bool
}

// ExplainSelfImprovementProvenanceHistory resolves one declaration against
// ordered provenance history without executing or authorizing a change.
func ExplainSelfImprovementProvenanceHistory(
	source, symbolName string,
	current RevisionSelfImprovementLifecycleObservation,
	history RevisionSelfImprovementProvenanceHistoryObservation,
) (SelfImprovementProvenanceHistorySymbolResult, error) {
	result := SelfImprovementProvenanceHistorySymbolResult{
		Status:                            "UNKNOWN",
		MissingStage:                      "lsp-self-improvement-provenance-history",
		SourceDigest:                      digestString(source),
		CurrentLifecycleObservationDigest: current.ObservationDigest,
		HistoryDigest:                     history.HistoryDigest,
		ObservationCount:                  history.ObservationCount,
		TransitionDigests:                 append([]string(nil), history.TransitionDigests...),
		FirstTransitionDigest:             history.FirstTransitionDigest,
		LastTransitionDigest:              history.LastTransitionDigest,
		FirstPreviousLifecycleDigest:      history.FirstPreviousLifecycleDigest,
		LastCurrentLifecycleDigest:        history.LastCurrentLifecycleDigest,
		StableCount:                       history.StableCount,
		TransitionedCount:                 history.TransitionedCount,
		MetricsChangeCount:                history.MetricsChangeCount,
		GenerationChangeCount:             history.GenerationChangeCount,
		HistorySignal:                     history.HistorySignal,
		NonExecuting:                      true,
		NonAuthorizing:                    true,
	}
	setResultDigest := func() {
		result.ResultDigest = digestLSPSelfImprovementProvenanceHistory(result)
	}
	setResultDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-history-snapshot"
		setResultDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := current.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-history-current"
		setResultDigest()
		return result, fmt.Errorf("current self-improvement lifecycle is not valid: %w", err)
	}
	if err := history.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-history-history"
		setResultDigest()
		return result, fmt.Errorf("self-improvement provenance history is not valid: %w", err)
	}
	if snapshot.SourceDigest != current.SourceDigest {
		result.MissingStage = "lsp-self-improvement-provenance-history-source-link"
		setResultDigest()
		return result, fmt.Errorf("LSP source digest does not match current lifecycle")
	}
	if history.LastCurrentLifecycleDigest != current.ObservationDigest {
		result.MissingStage = "lsp-self-improvement-provenance-history-link"
		setResultDigest()
		return result, fmt.Errorf("provenance history is not linked to current lifecycle")
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
				result.MissingStage = "lsp-self-improvement-provenance-history"
				setResultDigest()
				return result, fmt.Errorf("LSP provenance history result is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-self-improvement-provenance-history-symbol"
	setResultDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r SelfImprovementProvenanceHistorySymbolResult) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("LSP provenance history status is empty")
	}
	if r.Status == "BOUND" && r.MissingStage != "" {
		return fmt.Errorf("bound LSP provenance history has a missing stage")
	}
	if r.Status == "UNKNOWN" && r.MissingStage == "" {
		return fmt.Errorf("unknown LSP provenance history has no missing stage")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("LSP provenance history must remain non-executing and non-authorizing")
	}
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" ||
		r.SymbolName == "" || r.SymbolDigest == "") {
		return fmt.Errorf("bound LSP provenance history is incomplete")
	}
	if r.ObservationCount < 1 || len(r.TransitionDigests) != r.ObservationCount {
		return fmt.Errorf("LSP provenance history transition count is invalid")
	}
	if r.HistorySignal != "stable" && r.HistorySignal != "transitioned" &&
		r.HistorySignal != "mixed" {
		return fmt.Errorf("LSP provenance history signal is invalid")
	}
	if digestLSPSelfImprovementProvenanceHistory(r) != r.ResultDigest {
		return fmt.Errorf("LSP provenance history digest does not match its fields")
	}
	return nil
}

func digestLSPSelfImprovementProvenanceHistory(result SelfImprovementProvenanceHistorySymbolResult) string {
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
		strings.Join(result.TransitionDigests, ","),
		result.FirstTransitionDigest,
		result.LastTransitionDigest,
		result.FirstPreviousLifecycleDigest,
		result.LastCurrentLifecycleDigest,
		strconv.Itoa(result.StableCount),
		strconv.Itoa(result.TransitionedCount),
		strconv.Itoa(result.MetricsChangeCount),
		strconv.Itoa(result.GenerationChangeCount),
		result.HistorySignal,
		strconv.FormatBool(result.NonExecuting),
		strconv.FormatBool(result.NonAuthorizing),
	}
	return digestString(strings.Join(parts, "|"))
}
