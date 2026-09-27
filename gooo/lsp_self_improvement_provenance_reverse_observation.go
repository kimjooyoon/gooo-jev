package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// SelfImprovementProvenanceReverseSymbolResult is a read-only LSP projection
// of provenance-linked reverse observation evidence.
type SelfImprovementProvenanceReverseSymbolResult struct {
	Status                      string
	MissingStage                string
	SourceDigest                string
	IRDigest                    string
	SymbolName                  string
	SymbolKind                  SymbolKind
	SymbolDigest                string
	ApplicationFeedbackDigest   string
	ReverseObservationDigest    string
	PlanDispositionDigest       string
	IterationProvenanceDigest   string
	PlanDigest                  string
	ApplicationDigest           string
	ReceiptDigest               string
	ProposedSourceDigest        string
	ProposedIRDigest            string
	GeneratedIRDigest           string
	ObservedGeneratedIRDigest   string
	StructureDigest             string
	ObservedStructureDigest     string
	FeedbackSignal              string
	ApplicationSignal           string
	FeedbackApplicationSignal   string
	ReverseSignal               string
	ProvenanceReverseSignal     string
	SignalsAligned              bool
	ExactIRMatch                bool
	ExactStructureMatch         bool
	ResultDigest                string
	NonExecuting                bool
	NonAuthorizing              bool
}

// ExplainSelfImprovementProvenanceReverseObservation resolves a declaration
// against provenance-linked reverse evidence without executing or authorizing.
func ExplainSelfImprovementProvenanceReverseObservation(
	source, symbolName string,
	feedback RevisionSelfImprovementProvenanceApplicationFeedbackObservation,
	reverse RevisionSelfImprovementProvenanceReverseObservation,
) (SelfImprovementProvenanceReverseSymbolResult, error) {
	result := SelfImprovementProvenanceReverseSymbolResult{
		Status:                    "UNKNOWN",
		MissingStage:              "lsp-self-improvement-provenance-reverse-observation",
		SourceDigest:              digestString(source),
		ApplicationFeedbackDigest: feedback.ObservationDigest,
		ReverseObservationDigest:  reverse.ObservationDigest,
		PlanDispositionDigest:     reverse.PlanDispositionDigest,
		IterationProvenanceDigest: reverse.IterationProvenanceDigest,
		PlanDigest:                reverse.PlanDigest,
		ApplicationDigest:         reverse.ApplicationDigest,
		ReceiptDigest:             reverse.ReceiptDigest,
		ProposedSourceDigest:      reverse.ProposedSourceDigest,
		ProposedIRDigest:          reverse.ProposedIRDigest,
		GeneratedIRDigest:         reverse.GeneratedIRDigest,
		ObservedGeneratedIRDigest: reverse.ReverseGeneratedIRDigest,
		StructureDigest:           reverse.StructureDigest,
		ObservedStructureDigest:   reverse.ReverseStructureDigest,
		FeedbackSignal:            reverse.FeedbackSignal,
		ApplicationSignal:         reverse.ApplicationSignal,
		FeedbackApplicationSignal: reverse.FeedbackApplicationSignal,
		ReverseSignal:             reverse.ReverseSignal,
		ProvenanceReverseSignal:   reverse.ProvenanceReverseSignal,
		SignalsAligned:            reverse.SignalsAligned,
		ExactIRMatch:              reverse.ExactIRMatch,
		ExactStructureMatch:       reverse.ExactStructureMatch,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	setDigest := func() { result.ResultDigest = digestLSPSelfImprovementProvenanceReverseObservation(result) }
	setDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-reverse-snapshot"
		setDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := feedback.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-reverse-feedback"
		setDigest()
		return result, fmt.Errorf("provenance application feedback is not valid: %w", err)
	}
	if err := reverse.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-reverse-reverse"
		setDigest()
		return result, fmt.Errorf("provenance reverse observation is not valid: %w", err)
	}
	if reverse.ApplicationFeedbackDigest != feedback.ObservationDigest ||
		snapshot.SourceDigest != feedback.ApplicationSourceDigest ||
		reverse.ProposedSourceDigest != feedback.ApplicationProposedSourceDigest ||
		reverse.ProposedIRDigest != feedback.ApplicationProposedIRDigest {
		result.MissingStage = "lsp-self-improvement-provenance-reverse-link"
		setDigest()
		return result, fmt.Errorf("provenance reverse observation is not linked to application feedback")
	}
	for _, symbol := range snapshot.Symbols {
		if symbol.Name == symbolName {
			result.Status = "BOUND"
			result.MissingStage = ""
			result.IRDigest = snapshot.IRDigest
			result.SymbolName = symbol.Name
			result.SymbolKind = symbol.Kind
			result.SymbolDigest = symbol.Digest
			setDigest()
			if err := result.Validate(); err != nil {
				result.Status = "UNKNOWN"
				result.MissingStage = "lsp-self-improvement-provenance-reverse-observation"
				setDigest()
				return result, fmt.Errorf("LSP provenance reverse result is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-self-improvement-provenance-reverse-symbol"
	setDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r SelfImprovementProvenanceReverseSymbolResult) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("LSP provenance reverse status is empty")
	}
	if r.Status == "BOUND" && r.MissingStage != "" {
		return fmt.Errorf("bound LSP provenance reverse result has a missing stage")
	}
	if r.Status == "UNKNOWN" && r.MissingStage == "" {
		return fmt.Errorf("unknown LSP provenance reverse result has no missing stage")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("LSP provenance reverse result must remain non-executing and non-authorizing")
	}
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" || r.SymbolName == "" || r.SymbolDigest == "") {
		return fmt.Errorf("bound LSP provenance reverse result is incomplete")
	}
	if r.FeedbackSignal == "" || r.ApplicationSignal == "" || r.ReverseSignal == "" || r.ProvenanceReverseSignal == "" {
		return fmt.Errorf("LSP provenance reverse signals are incomplete")
	}
	if digestLSPSelfImprovementProvenanceReverseObservation(r) != r.ResultDigest {
		return fmt.Errorf("LSP provenance reverse digest does not match its fields")
	}
	return nil
}

func digestLSPSelfImprovementProvenanceReverseObservation(r SelfImprovementProvenanceReverseSymbolResult) string {
	parts := []string{r.Status, r.MissingStage, r.SourceDigest, r.IRDigest, r.SymbolName, string(r.SymbolKind), r.SymbolDigest, r.ApplicationFeedbackDigest, r.ReverseObservationDigest, r.PlanDispositionDigest, r.IterationProvenanceDigest, r.PlanDigest, r.ApplicationDigest, r.ReceiptDigest, r.ProposedSourceDigest, r.ProposedIRDigest, r.GeneratedIRDigest, r.ObservedGeneratedIRDigest, r.StructureDigest, r.ObservedStructureDigest, r.FeedbackSignal, r.ApplicationSignal, r.FeedbackApplicationSignal, r.ReverseSignal, r.ProvenanceReverseSignal, strconv.FormatBool(r.SignalsAligned), strconv.FormatBool(r.ExactIRMatch), strconv.FormatBool(r.ExactStructureMatch), strconv.FormatBool(r.NonExecuting), strconv.FormatBool(r.NonAuthorizing)}
	return digestString(strings.Join(parts, "|"))
}
