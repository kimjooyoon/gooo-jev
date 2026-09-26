package gooo

import "fmt"

// RevisionSelfImprovementExecutionPlanObservation binds a bounded decision to
// an observed application plan without permitting execution or authorization.
type RevisionSelfImprovementExecutionPlanObservation struct {
	Status                    string
	MissingStage              string
	DecisionDigest            string
	PlanObservationDigest     string
	MaterializationDigest     string
	PlanDigest                string
	BindingDigest             string
	SourceDigest              string
	InputIRDigest             string
	CandidateDigest           string
	EditDigest                string
	DecisionSignal             string
	DecisionReason             string
	PlanSignal                string
	PlanningSignal            string
	PlanObserved               bool
	RequiresObservation        bool
	RequiresReview             bool
	RequiresInspection         bool
	RequiresMeasurement        bool
	ObservationDigest          string
	NonExecuting               bool
	NonAuthorizing             bool
}

// ObserveRevisionSelfImprovementExecutionPlan binds a bounded decision to a
// validated candidate plan while keeping the next step observational only.
func ObserveRevisionSelfImprovementExecutionPlan(decision RevisionSelfImprovementDecisionObservation, planObservation RevisionCandidatePlanObservation, plan RevisionApplicationPlan) (RevisionSelfImprovementExecutionPlanObservation, error) {
	result := RevisionSelfImprovementExecutionPlanObservation{
		Status:                "UNKNOWN",
		MissingStage:          "revision-self-improvement-execution-plan",
		DecisionDigest:        decision.DecisionDigest,
		PlanObservationDigest: planObservation.ObservationDigest,
		MaterializationDigest: planObservation.MaterializationDigest,
		PlanDigest:            plan.PlanDigest,
		BindingDigest:         plan.BindingDigest,
		SourceDigest:          plan.SourceDigest,
		InputIRDigest:         plan.InputIRDigest,
		CandidateDigest:       plan.CandidateDigest,
		EditDigest:            plan.EditDigest,
		DecisionSignal:        decision.DecisionSignal,
		DecisionReason:        decision.DecisionReason,
		PlanSignal:            planObservation.PlanSignal,
		PlanningSignal:        fmt.Sprintf("%s-plan", decision.DecisionSignal),
		PlanObserved:          planObservation.PlanObserved,
		RequiresObservation:   decision.RequiresObservation,
		RequiresReview:        decision.RequiresReview,
		RequiresInspection:    decision.RequiresInspection,
		RequiresMeasurement:   decision.RequiresMeasurement,
		NonExecuting:          true,
		NonAuthorizing:        true,
	}
	setObservationDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementExecutionPlan(result)
	}
	setObservationDigest()

	if err := decision.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-execution-plan-decision"
		setObservationDigest()
		return result, fmt.Errorf("self-improvement decision is not valid: %w", err)
	}
	if err := planObservation.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-execution-plan-observation"
		setObservationDigest()
		return result, fmt.Errorf("candidate plan observation is not valid: %w", err)
	}
	if err := plan.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-execution-plan-plan"
		setObservationDigest()
		return result, fmt.Errorf("revision application plan is not valid: %w", err)
	}
	if decision.SourceDigest != planObservation.SourceDigest ||
		planObservation.SourceDigest != plan.SourceDigest {
		result.MissingStage = "revision-self-improvement-execution-plan-source-link"
		setObservationDigest()
		return result, fmt.Errorf("decision, candidate plan observation, and application plan sources are not linked")
	}
	if planObservation.CandidateDigest != plan.CandidateDigest ||
		planObservation.PlanDigest != plan.PlanDigest ||
		planObservation.BindingDigest != plan.BindingDigest ||
		planObservation.EditDigest != plan.EditDigest {
		result.MissingStage = "revision-self-improvement-execution-plan-plan-link"
		setObservationDigest()
		return result, fmt.Errorf("candidate plan observation and application plan are not linked")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setObservationDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-execution-plan"
		setObservationDigest()
		return result, fmt.Errorf("revision self-improvement execution plan observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementExecutionPlanObservation) Validate() error {
	if o.Status == "" {
		return fmt.Errorf("revision self-improvement execution plan status is empty")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound revision self-improvement execution plan has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown revision self-improvement execution plan has no missing stage")
	}
	for name, digest := range map[string]string{
		"decision":        o.DecisionDigest,
		"plan observation": o.PlanObservationDigest,
		"materialization": o.MaterializationDigest,
		"plan":            o.PlanDigest,
		"binding":         o.BindingDigest,
		"source":          o.SourceDigest,
		"input IR":        o.InputIRDigest,
		"candidate":       o.CandidateDigest,
		"edit":            o.EditDigest,
		"observation":     o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("revision self-improvement execution plan %s digest is invalid", name)
		}
	}
	if !o.PlanObserved {
		return fmt.Errorf("revision self-improvement execution plan must preserve plan observation")
	}
	if o.DecisionSignal != "observe" && o.DecisionSignal != "remeasure" &&
		o.DecisionSignal != "review" && o.DecisionSignal != "inspect" {
		return fmt.Errorf("revision self-improvement execution plan decision signal is invalid")
	}
	if o.PlanningSignal != fmt.Sprintf("%s-plan", o.DecisionSignal) {
		return fmt.Errorf("revision self-improvement execution plan planning signal is not linked")
	}
	if o.PlanSignal == "" {
		return fmt.Errorf("revision self-improvement execution plan plan signal is empty")
	}
	if o.Status == "BOUND" && o.DecisionReason == "" {
		return fmt.Errorf("bound revision self-improvement execution plan has no decision reason")
	}
	if !o.RequiresObservation {
		return fmt.Errorf("revision self-improvement execution plan must require observation")
	}
	switch o.DecisionSignal {
	case "observe":
		if o.RequiresReview || o.RequiresInspection || o.RequiresMeasurement {
			return fmt.Errorf("observe execution plan has incompatible requirements")
		}
	case "remeasure":
		if o.RequiresReview || o.RequiresInspection || !o.RequiresMeasurement {
			return fmt.Errorf("remeasure execution plan has incompatible requirements")
		}
	case "review":
		if o.RequiresMeasurement || o.RequiresInspection || !o.RequiresReview {
			return fmt.Errorf("review execution plan has incompatible requirements")
		}
	case "inspect":
		if o.RequiresReview || !o.RequiresMeasurement || !o.RequiresInspection {
			return fmt.Errorf("inspect execution plan has incompatible requirements")
		}
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("revision self-improvement execution plan must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementExecutionPlan(o) != o.ObservationDigest {
		return fmt.Errorf("revision self-improvement execution plan digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementExecutionPlan(observation RevisionSelfImprovementExecutionPlanObservation) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t|%t|%t|%t|%t|%t",
		observation.Status,
		observation.MissingStage,
		observation.DecisionDigest,
		observation.PlanObservationDigest,
		observation.MaterializationDigest,
		observation.PlanDigest,
		observation.BindingDigest,
		observation.SourceDigest,
		observation.InputIRDigest,
		observation.CandidateDigest,
		observation.EditDigest,
		observation.DecisionSignal,
		observation.DecisionReason,
		observation.PlanSignal,
		observation.PlanningSignal,
		observation.PlanObserved,
		observation.RequiresObservation,
		observation.RequiresReview,
		observation.RequiresInspection,
		observation.RequiresMeasurement,
		observation.NonExecuting,
		observation.NonAuthorizing,
	))
}
