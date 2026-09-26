package gooo

import "fmt"

// SelfImprovementFeedbackSymbolResult is a declaration-scoped, read-only LSP
// projection of the next-observation feedback signal.
type SelfImprovementFeedbackSymbolResult struct {
	Status                 string
	MissingStage           string
	SourceDigest           string
	IRDigest               string
	SymbolName             string
	SymbolKind             SymbolKind
	SymbolDigest           string
	WindowDigest           string
	FeedbackDigest         string
	ComparisonSignal       string
	FeedbackSignal         string
	FeedbackReason         string
	RequiresObservation    bool
	RequiresReview         bool
	RequiresInspection     bool
	RequiresMeasurement    bool
	ResultDigest            string
	NonExecuting           bool
	NonAuthorizing         bool
}

// ExplainSelfImprovementFeedback resolves the current candidate declaration
// against validated feedback without executing or authorizing a change.
func ExplainSelfImprovementFeedback(source, symbolName string, feedback RevisionSelfImprovementFeedback) (SelfImprovementFeedbackSymbolResult, error) {
	result := SelfImprovementFeedbackSymbolResult{
		Status:              "UNKNOWN",
		MissingStage:        "lsp-self-improvement-feedback",
		SourceDigest:        digestString(source),
		WindowDigest:        feedback.WindowDigest,
		FeedbackDigest:      feedback.FeedbackDigest,
		ComparisonSignal:    feedback.ComparisonSignal,
		FeedbackSignal:      feedback.FeedbackSignal,
		FeedbackReason:      feedback.FeedbackReason,
		RequiresObservation: feedback.RequiresObservation,
		RequiresReview:      feedback.RequiresReview,
		RequiresInspection:  feedback.RequiresInspection,
		RequiresMeasurement: feedback.RequiresMeasurement,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	setResultDigest := func() {
		result.ResultDigest = digestLSPSelfImprovementFeedback(result)
	}
	setResultDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-feedback-snapshot"
		setResultDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := feedback.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-feedback-feedback"
		setResultDigest()
		return result, fmt.Errorf("revision self-improvement feedback is not valid: %w", err)
	}
	if snapshot.SourceDigest != feedback.CandidateSourceDigest {
		result.MissingStage = "lsp-self-improvement-feedback-link"
		setResultDigest()
		return result, fmt.Errorf("LSP source digest does not match candidate feedback source")
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
				result.MissingStage = "lsp-self-improvement-feedback"
				setResultDigest()
				return result, fmt.Errorf("LSP self-improvement feedback result is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-self-improvement-feedback-symbol"
	setResultDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r SelfImprovementFeedbackSymbolResult) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("LSP self-improvement feedback status is empty")
	}
	if r.Status == "BOUND" && r.MissingStage != "" {
		return fmt.Errorf("bound LSP self-improvement feedback has a missing stage")
	}
	if r.Status == "UNKNOWN" && r.MissingStage == "" {
		return fmt.Errorf("unknown LSP self-improvement feedback has no missing stage")
	}
	if r.ComparisonSignal != "narrower" && r.ComparisonSignal != "wider" &&
		r.ComparisonSignal != "stable" && r.ComparisonSignal != "mixed" {
		return fmt.Errorf("LSP self-improvement feedback comparison signal is invalid")
	}
	if r.FeedbackSignal != "observe" && r.FeedbackSignal != "remeasure" &&
		r.FeedbackSignal != "review" && r.FeedbackSignal != "inspect" {
		return fmt.Errorf("LSP self-improvement feedback signal is invalid")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("LSP self-improvement feedback must remain non-executing and non-authorizing")
	}
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" ||
		r.SymbolName == "" || r.SymbolDigest == "") {
		return fmt.Errorf("bound LSP self-improvement feedback is incomplete")
	}
	if digestLSPSelfImprovementFeedback(r) != r.ResultDigest {
		return fmt.Errorf("LSP self-improvement feedback digest does not match its fields")
	}
	return nil
}

func digestLSPSelfImprovementFeedback(result SelfImprovementFeedbackSymbolResult) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t|%t|%t|%t|%t",
		result.Status,
		result.MissingStage,
		result.SourceDigest,
		result.IRDigest,
		result.SymbolName,
		result.SymbolKind,
		result.SymbolDigest,
		result.WindowDigest,
		result.FeedbackDigest,
		result.ComparisonSignal,
		result.FeedbackSignal,
		result.FeedbackReason,
		result.RequiresObservation,
		result.RequiresReview,
		result.RequiresInspection,
		result.RequiresMeasurement,
		result.NonExecuting,
		result.NonAuthorizing,
	))
}