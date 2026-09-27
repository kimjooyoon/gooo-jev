package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// SelfImprovementProvenanceApplicationFeedbackSymbolResult is a read-only
// LSP projection of plan and application feedback provenance.
type SelfImprovementProvenanceApplicationFeedbackSymbolResult struct {
	Status                    string
	MissingStage              string
	SourceDigest              string
	IRDigest                  string
	SymbolName                string
	SymbolKind                SymbolKind
	SymbolDigest              string
	ApplicationFeedbackDigest string
	PlanDispositionDigest     string
	IterationProvenanceDigest string
	PlanObservationDigest     string
	ApplicationObservationDigest string
	PlanDigest                string
	ApplicationDigest         string
	ReceiptDigest             string
	CandidateSourceDigest     string
	ProposedSourceDigest      string
	FeedbackSignal            string
	ApplicationSignal         string
	DispositionSignal          string
	FeedbackApplicationSignal string
	SignalsAligned            bool
	ApplicationObserved       bool
	ResultDigest              string
	NonExecuting              bool
	NonAuthorizing            bool
}

// ExplainSelfImprovementProvenanceApplicationFeedback resolves one declaration
// against application feedback provenance without executing or authorizing.
func ExplainSelfImprovementProvenanceApplicationFeedback(
	source, symbolName string,
	feedback RevisionSelfImprovementProvenanceApplicationFeedbackObservation,
) (SelfImprovementProvenanceApplicationFeedbackSymbolResult, error) {
	result := SelfImprovementProvenanceApplicationFeedbackSymbolResult{
		Status:                       "UNKNOWN",
		MissingStage:                 "lsp-self-improvement-provenance-application-feedback",
		SourceDigest:                 digestString(source),
		ApplicationFeedbackDigest:    feedback.ObservationDigest,
		PlanDispositionDigest:        feedback.PlanDispositionDigest,
		IterationProvenanceDigest:    feedback.IterationProvenanceDigest,
		PlanObservationDigest:        feedback.PlanObservationDigest,
		ApplicationObservationDigest: feedback.ApplicationObservationDigest,
		PlanDigest:                   feedback.PlanDigest,
		ApplicationDigest:            feedback.ApplicationDigest,
		ReceiptDigest:                feedback.ReceiptDigest,
		CandidateSourceDigest:        feedback.CandidateSourceDigest,
		ProposedSourceDigest:         feedback.ApplicationProposedSourceDigest,
		FeedbackSignal:               feedback.FeedbackSignal,
		ApplicationSignal:            feedback.ApplicationSignal,
		DispositionSignal:            feedback.DispositionSignal,
		FeedbackApplicationSignal:    feedback.FeedbackApplicationSignal,
		SignalsAligned:               feedback.SignalsAligned,
		ApplicationObserved:          feedback.ApplicationObserved,
		NonExecuting:                 true,
		NonAuthorizing:               true,
	}
	setDigest := func() { result.ResultDigest = digestLSPSelfImprovementProvenanceApplicationFeedback(result) }
	setDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-application-feedback-snapshot"
		setDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := feedback.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-application-feedback-feedback"
		setDigest()
		return result, fmt.Errorf("provenance application feedback is not valid: %w", err)
	}
	if snapshot.SourceDigest != feedback.ApplicationSourceDigest {
		result.MissingStage = "lsp-self-improvement-provenance-application-feedback-link"
		setDigest()
		return result, fmt.Errorf("LSP source digest does not match application feedback source")
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
				result.MissingStage = "lsp-self-improvement-provenance-application-feedback"
				setDigest()
				return result, fmt.Errorf("LSP provenance application feedback result is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-self-improvement-provenance-application-feedback-symbol"
	setDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r SelfImprovementProvenanceApplicationFeedbackSymbolResult) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("LSP provenance application feedback status is empty")
	}
	if r.Status == "BOUND" && r.MissingStage != "" {
		return fmt.Errorf("bound LSP provenance application feedback has a missing stage")
	}
	if r.Status == "UNKNOWN" && r.MissingStage == "" {
		return fmt.Errorf("unknown LSP provenance application feedback has no missing stage")
	}
	if r.FeedbackSignal != "observe" && r.FeedbackSignal != "remeasure" && r.FeedbackSignal != "review" && r.FeedbackSignal != "inspect" {
		return fmt.Errorf("LSP provenance application feedback signal is invalid")
	}
	if r.ApplicationSignal == "" || r.DispositionSignal == "" || r.FeedbackApplicationSignal == "" {
		return fmt.Errorf("LSP provenance application feedback signals are incomplete")
	}
	if !r.ApplicationObserved {
		return fmt.Errorf("LSP provenance application feedback was not observed")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("LSP provenance application feedback must remain non-executing and non-authorizing")
	}
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" || r.SymbolName == "" || r.SymbolDigest == "") {
		return fmt.Errorf("bound LSP provenance application feedback is incomplete")
	}
	if digestLSPSelfImprovementProvenanceApplicationFeedback(r) != r.ResultDigest {
		return fmt.Errorf("LSP provenance application feedback digest does not match its fields")
	}
	return nil
}

func digestLSPSelfImprovementProvenanceApplicationFeedback(r SelfImprovementProvenanceApplicationFeedbackSymbolResult) string {
	parts := []string{r.Status, r.MissingStage, r.SourceDigest, r.IRDigest, r.SymbolName, string(r.SymbolKind), r.SymbolDigest, r.ApplicationFeedbackDigest, r.PlanDispositionDigest, r.IterationProvenanceDigest, r.PlanObservationDigest, r.ApplicationObservationDigest, r.PlanDigest, r.ApplicationDigest, r.ReceiptDigest, r.CandidateSourceDigest, r.ProposedSourceDigest, r.FeedbackSignal, r.ApplicationSignal, r.DispositionSignal, r.FeedbackApplicationSignal, strconv.FormatBool(r.SignalsAligned), strconv.FormatBool(r.ApplicationObserved), strconv.FormatBool(r.NonExecuting), strconv.FormatBool(r.NonAuthorizing)}
	return digestString(strings.Join(parts, "|"))
}
