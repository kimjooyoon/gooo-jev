package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// SelfImprovementIterationProvenanceSymbolResult is a read-only LSP
// projection of iteration provenance and its declaration symbol.
type SelfImprovementIterationProvenanceSymbolResult struct {
	Status                            string
	MissingStage                      string
	SourceDigest                      string
	IRDigest                          string
	SymbolName                        string
	SymbolKind                        SymbolKind
	SymbolDigest                      string
	CurrentLifecycleObservationDigest string
	IterationDigest                   string
	LegacyHistoryDigest               string
	ProvenanceHistoryDigest           string
	WindowDigest                      string
	FeedbackDigest                    string
	BridgeDigest                      string
	CandidateGeneratedIRDigest        string
	LegacyHistorySignal               string
	ProvenanceHistorySignal           string
	DecisionSignal                    string
	FeedbackSignal                    string
	BridgeSignal                      string
	StageCount                        int
	MetricsChangeCount                int
	GenerationChangeCount             int
	RequiresObservation               bool
	RequiresInspection                bool
	RequiresMeasurement               bool
	ObservationDigest                 string
	ResultDigest                      string
	NonExecuting                      bool
	NonAuthorizing                    bool
}

// ExplainSelfImprovementIterationProvenance resolves one declaration against
// the current lifecycle and read-only iteration provenance evidence.
func ExplainSelfImprovementIterationProvenance(
	source, symbolName string,
	current RevisionSelfImprovementLifecycleObservation,
	observation RevisionSelfImprovementIterationProvenanceObservation,
) (SelfImprovementIterationProvenanceSymbolResult, error) {
	result := SelfImprovementIterationProvenanceSymbolResult{
		Status:                            "UNKNOWN",
		MissingStage:                      "lsp-self-improvement-iteration-provenance",
		SourceDigest:                      digestString(source),
		CurrentLifecycleObservationDigest: current.ObservationDigest,
		IterationDigest:                   observation.IterationDigest,
		LegacyHistoryDigest:               observation.LegacyHistoryDigest,
		ProvenanceHistoryDigest:           observation.ProvenanceHistoryDigest,
		WindowDigest:                      observation.WindowDigest,
		FeedbackDigest:                    observation.FeedbackDigest,
		BridgeDigest:                      observation.BridgeDigest,
		CandidateGeneratedIRDigest:        observation.CandidateGeneratedIRDigest,
		LegacyHistorySignal:               observation.LegacyHistorySignal,
		ProvenanceHistorySignal:           observation.ProvenanceHistorySignal,
		DecisionSignal:                    observation.DecisionSignal,
		FeedbackSignal:                    observation.FeedbackSignal,
		BridgeSignal:                      observation.BridgeSignal,
		StageCount:                        observation.StageCount,
		MetricsChangeCount:                observation.MetricsChangeCount,
		GenerationChangeCount:             observation.GenerationChangeCount,
		RequiresObservation:               observation.RequiresObservation,
		RequiresInspection:               observation.RequiresInspection,
		RequiresMeasurement:               observation.RequiresMeasurement,
		ObservationDigest:                 observation.ObservationDigest,
		NonExecuting:                      true,
		NonAuthorizing:                    true,
	}
	setResultDigest := func() {
		result.ResultDigest = digestLSPSelfImprovementIterationProvenance(result)
	}
	setResultDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-iteration-provenance-snapshot"
		setResultDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := current.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-iteration-provenance-current"
		setResultDigest()
		return result, fmt.Errorf("current self-improvement lifecycle is not valid: %w", err)
	}
	if err := observation.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-iteration-provenance-observation"
		setResultDigest()
		return result, fmt.Errorf("iteration provenance observation is not valid: %w", err)
	}
	if snapshot.SourceDigest != current.SourceDigest ||
		observation.SourceDigest != snapshot.SourceDigest {
		result.MissingStage = "lsp-self-improvement-iteration-provenance-source-link"
		setResultDigest()
		return result, fmt.Errorf("LSP source digest is not linked to current iteration provenance")
	}
	if observation.IRDigest != snapshot.IRDigest {
		result.MissingStage = "lsp-self-improvement-iteration-provenance-ir-link"
		setResultDigest()
		return result, fmt.Errorf("LSP IR digest is not linked to current iteration provenance")
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
				result.MissingStage = "lsp-self-improvement-iteration-provenance"
				setResultDigest()
				return result, fmt.Errorf("LSP iteration provenance result is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-self-improvement-iteration-provenance-symbol"
	setResultDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r SelfImprovementIterationProvenanceSymbolResult) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("LSP iteration provenance status is empty")
	}
	if r.Status == "BOUND" && r.MissingStage != "" {
		return fmt.Errorf("bound LSP iteration provenance has a missing stage")
	}
	if r.Status == "UNKNOWN" && r.MissingStage == "" {
		return fmt.Errorf("unknown LSP iteration provenance has no missing stage")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("LSP iteration provenance must remain non-executing and non-authorizing")
	}
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" ||
		r.SymbolName == "" || r.SymbolDigest == "") {
		return fmt.Errorf("bound LSP iteration provenance is incomplete")
	}
	for name, digest := range map[string]string{
		"source":                 r.SourceDigest,
		"ir":                     r.IRDigest,
		"symbol":                 r.SymbolDigest,
		"lifecycle":              r.CurrentLifecycleObservationDigest,
		"iteration":             r.IterationDigest,
		"legacy-history":         r.LegacyHistoryDigest,
		"provenance-history":     r.ProvenanceHistoryDigest,
		"window":                 r.WindowDigest,
		"feedback":               r.FeedbackDigest,
		"bridge":                 r.BridgeDigest,
		"candidate-generated-ir": r.CandidateGeneratedIRDigest,
		"observation":            r.ObservationDigest,
		"result":                 r.ResultDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("LSP iteration provenance %s digest is invalid", name)
		}
	}
	if r.StageCount != 4 {
		return fmt.Errorf("LSP iteration provenance stage count must be 4")
	}
	if r.LegacyHistorySignal != "stable" && r.LegacyHistorySignal != "narrower" &&
		r.LegacyHistorySignal != "wider" && r.LegacyHistorySignal != "mixed" {
		return fmt.Errorf("LSP iteration provenance legacy history signal is invalid")
	}
	if r.ProvenanceHistorySignal != "stable" &&
		r.ProvenanceHistorySignal != "transitioned" &&
		r.ProvenanceHistorySignal != "mixed" {
		return fmt.Errorf("LSP iteration provenance provenance history signal is invalid")
	}
	if r.DecisionSignal != "observe" && r.DecisionSignal != "inspect" {
		return fmt.Errorf("LSP iteration provenance decision signal is invalid")
	}
	if r.FeedbackSignal != "observe" && r.FeedbackSignal != "inspect" {
		return fmt.Errorf("LSP iteration provenance feedback signal is invalid")
	}
	if r.BridgeSignal != "provenance-feedback-bridge-bound" &&
		r.BridgeSignal != "provenance-feedback-bridge-unknown" {
		return fmt.Errorf("LSP iteration provenance bridge signal is invalid")
	}
	if r.MetricsChangeCount < 0 || r.GenerationChangeCount < 0 {
		return fmt.Errorf("LSP iteration provenance change counts are invalid")
	}
	if !r.RequiresObservation {
		return fmt.Errorf("LSP iteration provenance must require observation")
	}
	if r.ProvenanceHistorySignal == "stable" &&
		(r.LegacyHistorySignal != "stable" || r.DecisionSignal != "observe" ||
			r.FeedbackSignal != "observe") {
		return fmt.Errorf("stable LSP iteration provenance mapping is invalid")
	}
	if (r.ProvenanceHistorySignal == "transitioned" ||
		r.ProvenanceHistorySignal == "mixed") &&
		(r.LegacyHistorySignal != "mixed" || r.DecisionSignal != "inspect" ||
			r.FeedbackSignal != "inspect") {
		return fmt.Errorf("changed LSP iteration provenance mapping is invalid")
	}
	if digestLSPSelfImprovementIterationProvenance(r) != r.ResultDigest {
		return fmt.Errorf("LSP iteration provenance digest does not match its fields")
	}
	return nil
}

func digestLSPSelfImprovementIterationProvenance(result SelfImprovementIterationProvenanceSymbolResult) string {
	parts := []string{
		result.Status,
		result.MissingStage,
		result.SourceDigest,
		result.IRDigest,
		result.SymbolName,
		string(result.SymbolKind),
		result.SymbolDigest,
		result.CurrentLifecycleObservationDigest,
		result.IterationDigest,
		result.LegacyHistoryDigest,
		result.ProvenanceHistoryDigest,
		result.WindowDigest,
		result.FeedbackDigest,
		result.BridgeDigest,
		result.CandidateGeneratedIRDigest,
		result.LegacyHistorySignal,
		result.ProvenanceHistorySignal,
		result.DecisionSignal,
		result.FeedbackSignal,
		result.BridgeSignal,
		strconv.Itoa(result.StageCount),
		strconv.Itoa(result.MetricsChangeCount),
		strconv.Itoa(result.GenerationChangeCount),
		strconv.FormatBool(result.RequiresObservation),
		strconv.FormatBool(result.RequiresInspection),
		strconv.FormatBool(result.RequiresMeasurement),
		result.ObservationDigest,
		strconv.FormatBool(result.NonExecuting),
		strconv.FormatBool(result.NonAuthorizing),
	}
	return digestString(strings.Join(parts, "|"))
}
