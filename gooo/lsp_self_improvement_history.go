package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// SelfImprovementHistorySymbolResult is a declaration-scoped, read-only LSP
// projection of ordered self-improvement history.
type SelfImprovementHistorySymbolResult struct {
	Status                            string
	MissingStage                      string
	SourceDigest                      string
	IRDigest                          string
	SymbolName                        string
	SymbolKind                        SymbolKind
	SymbolDigest                      string
	ObservationCount                  int
	FirstWindowDigest                 string
	LastWindowDigest                  string
	FirstFeedbackDigest               string
	LastFeedbackDigest                string
	LastCandidateSourceDigest         string
	LastCandidateProposedSourceDigest string
	LastCandidateGeneratedIRDigest    string
	StableCount                       int
	NarrowerCount                     int
	WiderCount                        int
	MixedCount                        int
	ObserveCount                      int
	RemeasureCount                    int
	ReviewCount                       int
	InspectCount                     int
	HistorySignal                     string
	HistoryDigest                     string
	ResultDigest                      string
	NonExecuting                     bool
	NonAuthorizing                   bool
}

// ExplainSelfImprovementHistory resolves the current candidate declaration
// against validated ordered history without executing or authorizing it.
func ExplainSelfImprovementHistory(source, symbolName string, history RevisionSelfImprovementHistory) (SelfImprovementHistorySymbolResult, error) {
	result := SelfImprovementHistorySymbolResult{
		Status:                            "UNKNOWN",
		MissingStage:                      "lsp-self-improvement-history",
		SourceDigest:                      digestString(source),
		ObservationCount:                  history.ObservationCount,
		FirstWindowDigest:                 history.FirstWindowDigest,
		LastWindowDigest:                  history.LastWindowDigest,
		FirstFeedbackDigest:               history.FirstFeedbackDigest,
		LastFeedbackDigest:                history.LastFeedbackDigest,
		LastCandidateSourceDigest:         history.LastCandidateSourceDigest,
		LastCandidateProposedSourceDigest: history.LastCandidateProposedSourceDigest,
		LastCandidateGeneratedIRDigest:    history.LastCandidateGeneratedIRDigest,
		StableCount:                       history.StableCount,
		NarrowerCount:                     history.NarrowerCount,
		WiderCount:                        history.WiderCount,
		MixedCount:                        history.MixedCount,
		ObserveCount:                      history.ObserveCount,
		RemeasureCount:                    history.RemeasureCount,
		ReviewCount:                       history.ReviewCount,
		InspectCount:                      history.InspectCount,
		HistorySignal:                     history.HistorySignal,
		HistoryDigest:                     history.HistoryDigest,
		NonExecuting:                      true,
		NonAuthorizing:                    true,
	}
	setResultDigest := func() {
		result.ResultDigest = digestLSPSelfImprovementHistory(result)
	}
	setResultDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-history-snapshot"
		setResultDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := history.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-history-history"
		setResultDigest()
		return result, fmt.Errorf("revision self-improvement history is not valid: %w", err)
	}
	if snapshot.SourceDigest != history.LastCandidateSourceDigest {
		result.MissingStage = "lsp-self-improvement-history-link"
		setResultDigest()
		return result, fmt.Errorf("LSP source digest does not match last candidate history source")
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
				result.MissingStage = "lsp-self-improvement-history"
				setResultDigest()
				return result, fmt.Errorf("LSP self-improvement history result is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-self-improvement-history-symbol"
	setResultDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r SelfImprovementHistorySymbolResult) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("LSP self-improvement history status is empty")
	}
	if r.Status == "BOUND" && r.MissingStage != "" {
		return fmt.Errorf("bound LSP self-improvement history has a missing stage")
	}
	if r.Status == "UNKNOWN" && r.MissingStage == "" {
		return fmt.Errorf("unknown LSP self-improvement history has no missing stage")
	}
	if r.HistorySignal != "stable" && r.HistorySignal != "narrower" &&
		r.HistorySignal != "wider" && r.HistorySignal != "mixed" {
		return fmt.Errorf("LSP self-improvement history signal is invalid")
	}
	if r.ObservationCount < 1 {
		return fmt.Errorf("LSP self-improvement history observation count must be positive")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("LSP self-improvement history must remain non-executing and non-authorizing")
	}
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" ||
		r.SymbolName == "" || r.SymbolDigest == "") {
		return fmt.Errorf("bound LSP self-improvement history is incomplete")
	}
	if digestLSPSelfImprovementHistory(r) != r.ResultDigest {
		return fmt.Errorf("LSP self-improvement history digest does not match its fields")
	}
	return nil
}

func digestLSPSelfImprovementHistory(result SelfImprovementHistorySymbolResult) string {
	parts := []string{
		result.Status,
		result.MissingStage,
		result.SourceDigest,
		result.IRDigest,
		result.SymbolName,
		fmt.Sprint(result.SymbolKind),
		result.SymbolDigest,
		strconv.Itoa(result.ObservationCount),
		result.FirstWindowDigest,
		result.LastWindowDigest,
		result.FirstFeedbackDigest,
		result.LastFeedbackDigest,
		result.LastCandidateSourceDigest,
		result.LastCandidateProposedSourceDigest,
		result.LastCandidateGeneratedIRDigest,
		strconv.Itoa(result.StableCount),
		strconv.Itoa(result.NarrowerCount),
		strconv.Itoa(result.WiderCount),
		strconv.Itoa(result.MixedCount),
		strconv.Itoa(result.ObserveCount),
		strconv.Itoa(result.RemeasureCount),
		strconv.Itoa(result.ReviewCount),
		strconv.Itoa(result.InspectCount),
		result.HistorySignal,
		result.HistoryDigest,
		strconv.FormatBool(result.NonExecuting),
		strconv.FormatBool(result.NonAuthorizing),
	}
	return digestString(strings.Join(parts, "|"))
}