package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// SelfImprovementProvenanceReplanSymbolResult is a read-only LSP projection
// of the next execution plan in the provenance loop.
type SelfImprovementProvenanceReplanSymbolResult struct {
	Status                         string
	MissingStage                   string
	SourceDigest                   string
	IRDigest                      string
	SymbolName                    string
	SymbolKind                    SymbolKind
	SymbolDigest                  string
	ProvenanceDecisionDigest      string
	ExecutionPlanObservationDigest string
	DecisionDigest                string
	PreviousPlanDigest            string
	NextPlanDigest                string
	PlanObservationDigest         string
	MaterializationDigest         string
	BindingDigest                 string
	InputIRDigest                 string
	CandidateDigest               string
	EditDigest                   string
	DecisionSignal                string
	DecisionReason                string
	PlanningSignal                string
	PlanSignal                   string
	DecisionFeedbackSignal        string
	ReplanSignal                 string
	SignalsAligned               bool
	PlanObserved                 bool
	RequiresObservation          bool
	RequiresReview               bool
	RequiresInspection           bool
	RequiresMeasurement          bool
	ResultDigest                 string
	NonExecuting                 bool
	NonAuthorizing               bool
}

// ExplainSelfImprovementProvenanceReplan resolves a declaration against the
// next observed execution plan without executing or authorizing a revision.
func ExplainSelfImprovementProvenanceReplan(
	source, symbolName string,
	replan RevisionSelfImprovementProvenanceReplanObservation,
) (SelfImprovementProvenanceReplanSymbolResult, error) {
	result := SelfImprovementProvenanceReplanSymbolResult{
		Status:                          "UNKNOWN",
		MissingStage:                    "lsp-self-improvement-provenance-replan",
		SourceDigest:                    digestString(source),
		ProvenanceDecisionDigest:        replan.ProvenanceDecisionDigest,
		ExecutionPlanObservationDigest: replan.ExecutionPlanObservationDigest,
		DecisionDigest:                  replan.DecisionDigest,
		PreviousPlanDigest:              replan.PreviousPlanDigest,
		NextPlanDigest:                  replan.NextPlanDigest,
		PlanObservationDigest:           replan.PlanObservationDigest,
		MaterializationDigest:            replan.MaterializationDigest,
		BindingDigest:                   replan.BindingDigest,
		InputIRDigest:                   replan.InputIRDigest,
		CandidateDigest:                 replan.CandidateDigest,
		EditDigest:                      replan.EditDigest,
		DecisionSignal:                  replan.DecisionSignal,
		DecisionReason:                  replan.DecisionReason,
		PlanningSignal:                  replan.PlanningSignal,
		PlanSignal:                      replan.PlanSignal,
		DecisionFeedbackSignal:          replan.DecisionFeedbackSignal,
		ReplanSignal:                   replan.ReplanSignal,
		SignalsAligned:                 replan.SignalsAligned,
		PlanObserved:                   replan.PlanObserved,
		RequiresObservation:            replan.RequiresObservation,
		RequiresReview:                 replan.RequiresReview,
		RequiresInspection:             replan.RequiresInspection,
		RequiresMeasurement:            replan.RequiresMeasurement,
		NonExecuting:                   true,
		NonAuthorizing:                 true,
	}
	setDigest := func() { result.ResultDigest = digestLSPSelfImprovementProvenanceReplan(result) }
	setDigest()

	snapshot := Analyze(source)
	if err := snapshot.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-replan-snapshot"
		setDigest()
		return result, fmt.Errorf("language snapshot is not valid: %w", err)
	}
	if err := replan.Validate(); err != nil {
		result.MissingStage = "lsp-self-improvement-provenance-replan-replan"
		setDigest()
		return result, fmt.Errorf("provenance replan is not valid: %w", err)
	}
	if snapshot.SourceDigest != replan.SourceDigest {
		result.MissingStage = "lsp-self-improvement-provenance-replan-link"
		setDigest()
		return result, fmt.Errorf("LSP source digest does not match provenance replan source")
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
				result.MissingStage = "lsp-self-improvement-provenance-replan"
				setDigest()
				return result, fmt.Errorf("LSP provenance replan result is not valid: %w", err)
			}
			return result, nil
		}
	}
	result.MissingStage = "lsp-self-improvement-provenance-replan-symbol"
	setDigest()
	return result, fmt.Errorf("LSP symbol %q was not found", symbolName)
}

func (r SelfImprovementProvenanceReplanSymbolResult) Validate() error {
	if r.Status == "" { return fmt.Errorf("LSP provenance replan status is empty") }
	if r.Status == "BOUND" && r.MissingStage != "" { return fmt.Errorf("bound LSP provenance replan has a missing stage") }
	if r.Status == "UNKNOWN" && r.MissingStage == "" { return fmt.Errorf("unknown LSP provenance replan has no missing stage") }
	if r.DecisionSignal != "observe" && r.DecisionSignal != "remeasure" && r.DecisionSignal != "review" && r.DecisionSignal != "inspect" { return fmt.Errorf("LSP provenance replan decision signal is invalid") }
	if r.PlanningSignal != fmt.Sprintf("%s-plan", r.DecisionSignal) { return fmt.Errorf("LSP provenance replan planning signal is not linked") }
	if r.PlanSignal == "" || r.DecisionReason == "" { return fmt.Errorf("LSP provenance replan evidence is incomplete") }
	if r.DecisionFeedbackSignal != "provenance-decision-aligned" && r.DecisionFeedbackSignal != "provenance-decision-mismatch" { return fmt.Errorf("LSP provenance replan decision relationship is invalid") }
	if r.ReplanSignal != "provenance-replan-aligned" && r.ReplanSignal != "provenance-replan-mismatch" { return fmt.Errorf("LSP provenance replan signal is invalid") }
	if !r.PlanObserved || !r.RequiresObservation { return fmt.Errorf("LSP provenance replan must preserve observed required planning") }
	if !r.NonExecuting || !r.NonAuthorizing { return fmt.Errorf("LSP provenance replan must remain non-executing and non-authorizing") }
	if r.Status == "BOUND" && (r.SourceDigest == "" || r.IRDigest == "" || r.SymbolName == "" || r.SymbolDigest == "") { return fmt.Errorf("bound LSP provenance replan is incomplete") }
	if digestLSPSelfImprovementProvenanceReplan(r) != r.ResultDigest { return fmt.Errorf("LSP provenance replan digest does not match its fields") }
	return nil
}

func digestLSPSelfImprovementProvenanceReplan(r SelfImprovementProvenanceReplanSymbolResult) string {
	parts := []string{r.Status, r.MissingStage, r.SourceDigest, r.IRDigest, r.SymbolName, string(r.SymbolKind), r.SymbolDigest, r.ProvenanceDecisionDigest, r.ExecutionPlanObservationDigest, r.DecisionDigest, r.PreviousPlanDigest, r.NextPlanDigest, r.PlanObservationDigest, r.MaterializationDigest, r.BindingDigest, r.InputIRDigest, r.CandidateDigest, r.EditDigest, r.DecisionSignal, r.DecisionReason, r.PlanningSignal, r.PlanSignal, r.DecisionFeedbackSignal, r.ReplanSignal, strconv.FormatBool(r.SignalsAligned), strconv.FormatBool(r.PlanObserved), strconv.FormatBool(r.RequiresObservation), strconv.FormatBool(r.RequiresReview), strconv.FormatBool(r.RequiresInspection), strconv.FormatBool(r.RequiresMeasurement), strconv.FormatBool(r.NonExecuting), strconv.FormatBool(r.NonAuthorizing)}
	return digestString(strings.Join(parts, "|"))
}
