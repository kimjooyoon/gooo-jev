package gooo

import "fmt"

// SelfImprovementDecisionSymbolResult is a declaration-scoped, read-only LSP
// projection of a bounded next-observation decision.
type SelfImprovementDecisionSymbolResult struct {
	Status                string
	MissingStage          string
	SourceDigest          string
	IRDigest              string
	SymbolName            string
	SymbolKind            SymbolKind
	SymbolDigest          string
	IterationDigest       string
	OutcomeDigest         string
	FeedbackDigest        string
	HistoryDigest         string
	ComparisonSignal      string
	FeedbackSignal        string
	DecisionSignal        string
	DecisionReason        string
	RequiresObservation   bool
	RequiresReview        bool
	RequiresInspection    bool
	RequiresMeasurement   bool
	ResultDigest          string
	NonExecuting          bool
	NonAuthorizing        bool
}

// ExplainSelfImprovementDecision resolves a declaration against a validated
// bounded decision without executing or authorizing a change.
func ExplainSelfImprovementDecision(source, symbolName string, decision RevisionSelfImprovementDecisionObservation) (SelfImprovementDecisionSymbolResult, error) {
	result := SelfImprovementDecisionSymbolResult{
		Status:              "UNKNOWN",
		MissingStage:        "lsp-self-improvement-decision",
		SourceDigest:        digestString(source),
		IterationDigest:     decision.IterationDigest,
		OutcomeDigest:       decision.OutcomeDigest,
		FeedbackDigest:      decision.FeedbackDigest,
		HistoryDigest:       decision.HistoryDigest,
		ComparisonSignal:    decision.ComparisonSignal,
		FeedbackSignal:      decision.FeedbackSignal,
		DecisionSignal:      decision.DecisionSignal,
		DecisionReason:      decision.DecisionReason,
		RequiresObservation: decision.RequiresObservation,
		RequiresReview:      decision.RequiresReview,
		RequiresInspection:  decision.RequiresInspection,
		RequiresMeasurement: decision.RequiresMeasurement,
		NonExecuting:        true,
		NonAuthorizing:      true,
	}
	setResultDigest := func() {
		result.ResultDigest = digestLSPSelfImprovementDecision(result)
	}
	setResultDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-decision-snapshot"
		setResultDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := decision.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-decision-decision"
		setResultDigest()
		return result, fmt.Errorf("self-improvement decision is not valid: %w", err)
	}
	if snapshot.SourceDigest != decision.SourceDigest {
		result.MissingStage = "lsp-self-improvement-decision-link"
		setResultDigest()
		return result, fmt.Errorf("LSP source digest does not match self-improvement decision source")
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
				result.MissingStage = "lsp-self-improvement-decision"
				setResultDigest()
				return result, fmt.Errorf("LSP self-improvement decision result is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-self-improvement-decision-symbol"
	setResultDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r SelfImprovementDecisionSymbolResult) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("LSP self-improvement decision status is empty")
	}
	if r.Status == "BOUND" && r.MissingStage != "" {
		return fmt.Errorf("bound LSP self-improvement decision has a missing stage")
	}
	if r.Status == "UNKNOWN" && r.MissingStage == "" {
		return fmt.Errorf("unknown LSP self-improvement decision has no missing stage")
	}
	if r.ComparisonSignal != "narrower" && r.ComparisonSignal != "wider" &&
		r.ComparisonSignal != "stable" && r.ComparisonSignal != "mixed" {
		return fmt.Errorf("LSP self-improvement decision comparison signal is invalid")
	}
	if r.FeedbackSignal != "observe" && r.FeedbackSignal != "remeasure" &&
		r.FeedbackSignal != "review" && r.FeedbackSignal != "inspect" {
		return fmt.Errorf("LSP self-improvement decision feedback signal is invalid")
	}
	if r.DecisionSignal != r.FeedbackSignal {
		return fmt.Errorf("LSP self-improvement decision signal is not linked to feedback")
	}
	if !r.RequiresObservation {
		return fmt.Errorf("LSP self-improvement decision must require observation")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("LSP self-improvement decision must remain non-executing and non-authorizing")
	}
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" ||
		r.SymbolName == "" || r.SymbolDigest == "" || r.DecisionReason == "") {
		return fmt.Errorf("bound LSP self-improvement decision is incomplete")
	}
	if digestLSPSelfImprovementDecision(r) != r.ResultDigest {
		return fmt.Errorf("LSP self-improvement decision digest does not match its fields")
	}
	return nil
}

func digestLSPSelfImprovementDecision(result SelfImprovementDecisionSymbolResult) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t|%t|%t|%t|%t",
		result.Status,
		result.MissingStage,
		result.SourceDigest,
		result.IRDigest,
		result.SymbolName,
		result.SymbolKind,
		result.SymbolDigest,
		result.IterationDigest,
		result.OutcomeDigest,
		result.FeedbackDigest,
		result.HistoryDigest,
		result.ComparisonSignal,
		result.FeedbackSignal,
		result.DecisionSignal,
		result.DecisionReason,
		result.RequiresObservation,
		result.RequiresReview,
		result.RequiresInspection,
		result.RequiresMeasurement,
		result.NonExecuting,
		result.NonAuthorizing,
	))
}
