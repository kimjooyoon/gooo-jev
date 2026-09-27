package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// RevisionSelfImprovementProvenanceReplanObservation records the handoff from
// a provenance decision bridge to the next observed execution plan.
type RevisionSelfImprovementProvenanceReplanObservation struct {
	Status                         string
	MissingStage                   string
	ProvenanceDecisionDigest       string
	ExecutionPlanObservationDigest string
	DecisionDigest                 string
	PreviousPlanDigest             string
	NextPlanDigest                 string
	PlanObservationDigest          string
	MaterializationDigest          string
	BindingDigest                  string
	SourceDigest                   string
	InputIRDigest                 string
	CandidateDigest                string
	EditDigest                    string
	DecisionSignal                 string
	DecisionReason                 string
	PlanningSignal                 string
	PlanSignal                    string
	DecisionFeedbackSignal         string
	ReplanSignal                   string
	SignalsAligned                 bool
	PlanObserved                   bool
	RequiresObservation            bool
	RequiresReview                 bool
	RequiresInspection             bool
	RequiresMeasurement            bool
	ObservationDigest              string
	NonExecuting                   bool
	NonAuthorizing                 bool
}

// ObserveRevisionSelfImprovementProvenanceReplan binds the next execution
// plan to the provenance decision without executing or authorizing it.
func ObserveRevisionSelfImprovementProvenanceReplan(
	decision RevisionSelfImprovementProvenanceReverseDecisionObservation,
	plan RevisionSelfImprovementExecutionPlanObservation,
) (RevisionSelfImprovementProvenanceReplanObservation, error) {
	result := RevisionSelfImprovementProvenanceReplanObservation{
		Status:                          "UNKNOWN",
		MissingStage:                    "revision-self-improvement-provenance-replan",
		ProvenanceDecisionDigest:        decision.ObservationDigest,
		ExecutionPlanObservationDigest:  plan.ObservationDigest,
		DecisionDigest:                  plan.DecisionDigest,
		PreviousPlanDigest:              decision.PlanDigest,
		NextPlanDigest:                  plan.PlanDigest,
		PlanObservationDigest:           plan.PlanObservationDigest,
		MaterializationDigest:           plan.MaterializationDigest,
		BindingDigest:                   plan.BindingDigest,
		SourceDigest:                    plan.SourceDigest,
		InputIRDigest:                   plan.InputIRDigest,
		CandidateDigest:                 plan.CandidateDigest,
		EditDigest:                      plan.EditDigest,
		DecisionSignal:                  plan.DecisionSignal,
		DecisionReason:                  plan.DecisionReason,
		PlanningSignal:                  plan.PlanningSignal,
		PlanSignal:                      plan.PlanSignal,
		DecisionFeedbackSignal:          decision.DecisionFeedbackSignal,
		ReplanSignal:                    "provenance-replan-mismatch",
		SignalsAligned:                  false,
		PlanObserved:                    plan.PlanObserved,
		RequiresObservation:             plan.RequiresObservation,
		RequiresReview:                  plan.RequiresReview,
		RequiresInspection:              plan.RequiresInspection,
		RequiresMeasurement:             plan.RequiresMeasurement,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	setDigest := func() { result.ObservationDigest = digestRevisionSelfImprovementProvenanceReplan(result) }
	setDigest()

	if err := decision.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-replan-decision"
		setDigest()
		return result, fmt.Errorf("provenance reverse decision is not valid: %w", err)
	}
	if err := plan.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-replan-plan"
		setDigest()
		return result, fmt.Errorf("execution plan observation is not valid: %w", err)
	}
	if decision.DecisionDigest != plan.DecisionDigest {
		result.MissingStage = "revision-self-improvement-provenance-replan-decision-link"
		setDigest()
		return result, fmt.Errorf("execution plan is not linked to provenance decision")
	}
	if decision.SourceDigest != plan.SourceDigest {
		result.MissingStage = "revision-self-improvement-provenance-replan-source-link"
		setDigest()
		return result, fmt.Errorf("execution plan source is not linked to provenance decision")
	}
	if decision.DecisionSignal != plan.DecisionSignal || decision.DecisionReason != plan.DecisionReason {
		result.MissingStage = "revision-self-improvement-provenance-replan-signal-link"
		setDigest()
		return result, fmt.Errorf("execution plan decision signal is not linked to provenance decision")
	}
	if decision.RequiresObservation != plan.RequiresObservation ||
		decision.RequiresReview != plan.RequiresReview ||
		decision.RequiresInspection != plan.RequiresInspection ||
		decision.RequiresMeasurement != plan.RequiresMeasurement {
		result.MissingStage = "revision-self-improvement-provenance-replan-requirement-link"
		setDigest()
		return result, fmt.Errorf("execution plan requirements are not linked to provenance decision")
	}
	if !plan.PlanObserved {
		result.MissingStage = "revision-self-improvement-provenance-replan-observed"
		setDigest()
		return result, fmt.Errorf("next execution plan was not observed")
	}
	result.SignalsAligned = decision.SignalsAligned && decision.DecisionSignal == plan.DecisionSignal
	if result.SignalsAligned {
		result.ReplanSignal = "provenance-replan-aligned"
	} else {
		result.ReplanSignal = "provenance-replan-mismatch"
	}
	result.Status = "BOUND"
	result.MissingStage = ""
	setDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-provenance-replan"
		result.ReplanSignal = "provenance-replan-mismatch"
		setDigest()
		return result, fmt.Errorf("provenance replan observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementProvenanceReplanObservation) Validate() error {
	if o.Status == "" { return fmt.Errorf("provenance replan status is empty") }
	if o.Status == "BOUND" && o.MissingStage != "" { return fmt.Errorf("bound provenance replan has a missing stage") }
	if o.Status == "UNKNOWN" && o.MissingStage == "" { return fmt.Errorf("unknown provenance replan has no missing stage") }
	for name, digest := range map[string]string{
		"provenance decision": o.ProvenanceDecisionDigest, "execution plan": o.ExecutionPlanObservationDigest,
		"decision": o.DecisionDigest, "previous plan": o.PreviousPlanDigest, "next plan": o.NextPlanDigest,
		"plan observation": o.PlanObservationDigest, "materialization": o.MaterializationDigest,
		"binding": o.BindingDigest, "source": o.SourceDigest, "input IR": o.InputIRDigest,
		"candidate": o.CandidateDigest, "edit": o.EditDigest, "observation": o.ObservationDigest,
	} {
		if !validDigest(digest) { return fmt.Errorf("provenance replan %s digest is invalid", name) }
	}
	if o.DecisionSignal != "observe" && o.DecisionSignal != "remeasure" && o.DecisionSignal != "review" && o.DecisionSignal != "inspect" { return fmt.Errorf("provenance replan decision signal is invalid") }
	if o.PlanningSignal != fmt.Sprintf("%s-plan", o.DecisionSignal) { return fmt.Errorf("provenance replan planning signal is not linked") }
	if o.PlanSignal == "" || o.DecisionReason == "" { return fmt.Errorf("provenance replan plan evidence is incomplete") }
	if o.DecisionFeedbackSignal != "provenance-decision-aligned" && o.DecisionFeedbackSignal != "provenance-decision-mismatch" { return fmt.Errorf("provenance replan decision relationship is invalid") }
	if o.ReplanSignal != "provenance-replan-aligned" && o.ReplanSignal != "provenance-replan-mismatch" { return fmt.Errorf("provenance replan signal is invalid") }
	if o.SignalsAligned != (o.ReplanSignal == "provenance-replan-aligned") { return fmt.Errorf("provenance replan alignment is inconsistent") }
	if !o.PlanObserved || !o.RequiresObservation { return fmt.Errorf("provenance replan must preserve observed required planning") }
	if !o.NonExecuting || !o.NonAuthorizing { return fmt.Errorf("provenance replan must remain non-executing and non-authorizing") }
	if digestRevisionSelfImprovementProvenanceReplan(o) != o.ObservationDigest { return fmt.Errorf("provenance replan digest does not match its fields") }
	return nil
}

func digestRevisionSelfImprovementProvenanceReplan(o RevisionSelfImprovementProvenanceReplanObservation) string {
	parts := []string{o.Status, o.MissingStage, o.ProvenanceDecisionDigest, o.ExecutionPlanObservationDigest, o.DecisionDigest, o.PreviousPlanDigest, o.NextPlanDigest, o.PlanObservationDigest, o.MaterializationDigest, o.BindingDigest, o.SourceDigest, o.InputIRDigest, o.CandidateDigest, o.EditDigest, o.DecisionSignal, o.DecisionReason, o.PlanningSignal, o.PlanSignal, o.DecisionFeedbackSignal, o.ReplanSignal, strconv.FormatBool(o.SignalsAligned), strconv.FormatBool(o.PlanObserved), strconv.FormatBool(o.RequiresObservation), strconv.FormatBool(o.RequiresReview), strconv.FormatBool(o.RequiresInspection), strconv.FormatBool(o.RequiresMeasurement), strconv.FormatBool(o.NonExecuting), strconv.FormatBool(o.NonAuthorizing)}
	return digestString(strings.Join(parts, "|"))
}
