package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// SelfImprovementProvenanceReverseDecisionSymbolResult is a read-only LSP
// projection of the next decision and its full provenance bridge.
type SelfImprovementProvenanceReverseDecisionSymbolResult struct {
	Status                    string
	MissingStage              string
	SourceDigest              string
	IRDigest                  string
	SymbolName                string
	SymbolKind                SymbolKind
	SymbolDigest              string
	ApplicationFeedbackDigest string
	ReverseObservationDigest  string
	DecisionDigest            string
	PlanDispositionDigest     string
	IterationProvenanceDigest string
	PlanDigest                string
	ApplicationDigest         string
	ReceiptDigest             string
	ProposedSourceDigest      string
	DecisionProposedSourceDigest string
	GeneratedIRDigest         string
	DecisionGeneratedIRDigest string
	FeedbackSignal            string
	ApplicationSignal         string
	FeedbackApplicationSignal string
	ReverseSignal             string
	ProvenanceReverseSignal   string
	DecisionSignal            string
	DecisionFeedbackSignal    string
	DecisionReason            string
	SignalsAligned            bool
	RequiresObservation       bool
	RequiresReview            bool
	RequiresInspection        bool
	RequiresMeasurement       bool
	ResultDigest              string
	NonExecuting              bool
	NonAuthorizing            bool
}

// ExplainSelfImprovementProvenanceReverseDecision resolves a declaration
// against the next decision evidence without executing or authorizing it.
func ExplainSelfImprovementProvenanceReverseDecision(
	source, symbolName string,
	decision RevisionSelfImprovementProvenanceReverseDecisionObservation,
) (SelfImprovementProvenanceReverseDecisionSymbolResult, error) {
	result := SelfImprovementProvenanceReverseDecisionSymbolResult{
		Status:                       "UNKNOWN",
		MissingStage:                 "lsp-self-improvement-provenance-reverse-decision",
		SourceDigest:                 digestString(source),
		ApplicationFeedbackDigest:    decision.ApplicationFeedbackDigest,
		ReverseObservationDigest:     decision.ReverseObservationDigest,
		DecisionDigest:               decision.DecisionDigest,
		PlanDispositionDigest:        decision.PlanDispositionDigest,
		IterationProvenanceDigest:    decision.IterationProvenanceDigest,
		PlanDigest:                   decision.PlanDigest,
		ApplicationDigest:            decision.ApplicationDigest,
		ReceiptDigest:                decision.ReceiptDigest,
		ProposedSourceDigest:         decision.ProposedSourceDigest,
		DecisionProposedSourceDigest: decision.DecisionProposedSourceDigest,
		GeneratedIRDigest:            decision.GeneratedIRDigest,
		DecisionGeneratedIRDigest:    decision.DecisionGeneratedIRDigest,
		FeedbackSignal:               decision.FeedbackSignal,
		ApplicationSignal:            decision.ApplicationSignal,
		FeedbackApplicationSignal:    decision.FeedbackApplicationSignal,
		ReverseSignal:                decision.ReverseSignal,
		ProvenanceReverseSignal:      decision.ProvenanceReverseSignal,
		DecisionSignal:               decision.DecisionSignal,
		DecisionFeedbackSignal:       decision.DecisionFeedbackSignal,
		DecisionReason:               decision.DecisionReason,
		SignalsAligned:               decision.SignalsAligned,
		RequiresObservation:          decision.RequiresObservation,
		RequiresReview:               decision.RequiresReview,
		RequiresInspection:           decision.RequiresInspection,
		RequiresMeasurement:          decision.RequiresMeasurement,
		NonExecuting:                 true,
		NonAuthorizing:               true,
	}
	setDigest := func() { result.ResultDigest = digestLSPSelfImprovementProvenanceReverseDecision(result) }
	setDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-reverse-decision-snapshot"
		setDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := decision.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-reverse-decision-decision"
		setDigest()
		return result, fmt.Errorf("provenance reverse decision is not valid: %w", err)
	}
	if snapshot.SourceDigest != decision.SourceDigest {
		result.MissingStage = "lsp-self-improvement-provenance-reverse-decision-link"
		setDigest()
		return result, fmt.Errorf("LSP source digest does not match provenance reverse decision source")
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
				result.MissingStage = "lsp-self-improvement-provenance-reverse-decision"
				setDigest()
				return result, fmt.Errorf("LSP provenance reverse decision result is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-self-improvement-provenance-reverse-decision-symbol"
	setDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r SelfImprovementProvenanceReverseDecisionSymbolResult) Validate() error {
	if r.Status == "" { return fmt.Errorf("LSP provenance reverse decision status is empty") }
	if r.Status == "BOUND" && r.MissingStage != "" { return fmt.Errorf("bound LSP provenance reverse decision has a missing stage") }
	if r.Status == "UNKNOWN" && r.MissingStage == "" { return fmt.Errorf("unknown LSP provenance reverse decision has no missing stage") }
	if r.FeedbackSignal != "observe" && r.FeedbackSignal != "remeasure" && r.FeedbackSignal != "review" && r.FeedbackSignal != "inspect" { return fmt.Errorf("LSP provenance reverse decision feedback signal is invalid") }
	if r.DecisionSignal != r.FeedbackSignal { return fmt.Errorf("LSP provenance reverse decision signal is not linked to feedback") }
	if r.DecisionReason == "" && r.Status == "BOUND" { return fmt.Errorf("bound LSP provenance reverse decision has no reason") }
	if !r.RequiresObservation { return fmt.Errorf("LSP provenance reverse decision must require observation") }
	if !r.NonExecuting || !r.NonAuthorizing { return fmt.Errorf("LSP provenance reverse decision must remain non-executing and non-authorizing") }
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" || r.SymbolName == "" || r.SymbolDigest == "") { return fmt.Errorf("bound LSP provenance reverse decision is incomplete") }
	if digestLSPSelfImprovementProvenanceReverseDecision(r) != r.ResultDigest { return fmt.Errorf("LSP provenance reverse decision digest does not match its fields") }
	return nil
}

func digestLSPSelfImprovementProvenanceReverseDecision(r SelfImprovementProvenanceReverseDecisionSymbolResult) string {
	parts := []string{r.Status, r.MissingStage, r.SourceDigest, r.IRDigest, r.SymbolName, string(r.SymbolKind), r.SymbolDigest, r.ApplicationFeedbackDigest, r.ReverseObservationDigest, r.DecisionDigest, r.PlanDispositionDigest, r.IterationProvenanceDigest, r.PlanDigest, r.ApplicationDigest, r.ReceiptDigest, r.ProposedSourceDigest, r.DecisionProposedSourceDigest, r.GeneratedIRDigest, r.DecisionGeneratedIRDigest, r.FeedbackSignal, r.ApplicationSignal, r.FeedbackApplicationSignal, r.ReverseSignal, r.ProvenanceReverseSignal, r.DecisionSignal, r.DecisionFeedbackSignal, r.DecisionReason, strconv.FormatBool(r.SignalsAligned), strconv.FormatBool(r.RequiresObservation), strconv.FormatBool(r.RequiresReview), strconv.FormatBool(r.RequiresInspection), strconv.FormatBool(r.RequiresMeasurement), strconv.FormatBool(r.NonExecuting), strconv.FormatBool(r.NonAuthorizing)}
	return digestString(strings.Join(parts, "|"))
}
