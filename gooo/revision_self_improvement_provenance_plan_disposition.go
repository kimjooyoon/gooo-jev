package gooo

import (
	"fmt"
	"strconv"
	"strings"
)

// RevisionSelfImprovementProvenancePlanDispositionObservation records the
// bounded handoff from provenance iteration evidence to an observed plan.
type RevisionSelfImprovementProvenancePlanDispositionObservation struct {
	Status                          string
	MissingStage                    string
	IterationProvenanceDigest       string
	PlanObservationDigest           string
	PlanDigest                      string
	DecisionDigest                  string
	CurrentSourceDigest             string
	CandidateSourceDigest           string
	PlanSourceDigest                string
	CandidateProposedSourceDigest   string
	CandidateGeneratedIRDigest     string
	PlanInputIRDigest               string
	ProvenanceHistorySignal         string
	ProvenanceDecisionSignal        string
	DecisionSignal                  string
	PlanningSignal                  string
	PlanSignal                      string
	DispositionSignal               string
	PlanObserved                    bool
	ProvenanceRequiresObservation   bool
	ProvenanceRequiresReview        bool
	ProvenanceRequiresInspection    bool
	ProvenanceRequiresMeasurement   bool
	SignalsAligned                  bool
	RequiresObservation             bool
	RequiresReview                  bool
	RequiresInspection              bool
	RequiresMeasurement             bool
	ObservationDigest               string
	NonExecuting                    bool
	NonAuthorizing                  bool
}

// ObserveRevisionSelfImprovementProvenancePlanDisposition links a validated
// provenance iteration to an observed execution plan without executing or
// authorizing the plan.
func ObserveRevisionSelfImprovementProvenancePlanDisposition(
	iteration RevisionSelfImprovementIterationProvenanceObservation,
	plan RevisionSelfImprovementExecutionPlanObservation,
) (RevisionSelfImprovementProvenancePlanDispositionObservation, error) {
	result := RevisionSelfImprovementProvenancePlanDispositionObservation{
		Status:                        "UNKNOWN",
		MissingStage:                  "revision-self-improvement-provenance-plan-disposition",
		IterationProvenanceDigest:     iteration.ObservationDigest,
		PlanObservationDigest:         plan.PlanObservationDigest,
		PlanDigest:                    plan.PlanDigest,
		DecisionDigest:                plan.DecisionDigest,
		CurrentSourceDigest:           iteration.SourceDigest,
		CandidateSourceDigest:         iteration.CandidateSourceDigest,
		PlanSourceDigest:              plan.SourceDigest,
		CandidateProposedSourceDigest: iteration.CandidateProposedSourceDigest,
		CandidateGeneratedIRDigest:   iteration.CandidateGeneratedIRDigest,
		PlanInputIRDigest:             plan.InputIRDigest,
		ProvenanceHistorySignal:       iteration.ProvenanceHistorySignal,
		ProvenanceDecisionSignal:      iteration.DecisionSignal,
		DecisionSignal:                plan.DecisionSignal,
		PlanningSignal:                plan.PlanningSignal,
		PlanSignal:                    plan.PlanSignal,
		DispositionSignal:             "provenance-plan-unknown",
		PlanObserved:                  plan.PlanObserved,
		ProvenanceRequiresObservation: iteration.RequiresObservation,
		ProvenanceRequiresReview:      false,
		ProvenanceRequiresInspection:  iteration.RequiresInspection,
		ProvenanceRequiresMeasurement:  iteration.RequiresMeasurement,
		SignalsAligned:                 plan.DecisionSignal == iteration.DecisionSignal &&
			plan.RequiresObservation == iteration.RequiresObservation &&
			plan.RequiresInspection == iteration.RequiresInspection &&
			plan.RequiresMeasurement == iteration.RequiresMeasurement,
		RequiresObservation:           plan.RequiresObservation,
		RequiresReview:                plan.RequiresReview,
		RequiresInspection:            plan.RequiresInspection,
		RequiresMeasurement:           plan.RequiresMeasurement,
		NonExecuting:                  true,
		NonAuthorizing:                true,
	}
	setObservationDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementProvenancePlanDisposition(result)
	}
	setObservationDigest()

	if err := iteration.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-plan-disposition-iteration"
		setObservationDigest()
		return result, fmt.Errorf("iteration provenance is not valid: %w", err)
	}
	if err := plan.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-provenance-plan-disposition-plan"
		setObservationDigest()
		return result, fmt.Errorf("execution plan observation is not valid: %w", err)
	}
	if plan.SourceDigest != iteration.CandidateSourceDigest {
		result.MissingStage = "revision-self-improvement-provenance-plan-disposition-source-link"
		setObservationDigest()
		return result, fmt.Errorf("execution plan source is not linked to the iteration candidate source")
	}
	if result.SignalsAligned {
		switch iteration.ProvenanceHistorySignal {
		case "stable":
			if iteration.DecisionSignal != "observe" ||
				plan.PlanningSignal != "observe-plan" ||
				plan.RequiresReview || plan.RequiresInspection || plan.RequiresMeasurement {
				result.MissingStage = "revision-self-improvement-provenance-plan-disposition-mapping"
				setObservationDigest()
				return result, fmt.Errorf("stable provenance does not map to observe plan")
			}
		case "transitioned", "mixed":
			if iteration.DecisionSignal != "inspect" ||
				plan.PlanningSignal != "inspect-plan" ||
				!plan.RequiresInspection {
				result.MissingStage = "revision-self-improvement-provenance-plan-disposition-mapping"
				setObservationDigest()
				return result, fmt.Errorf("changed provenance does not map to inspect plan")
			}
		default:
			result.MissingStage = "revision-self-improvement-provenance-plan-disposition-signal"
			setObservationDigest()
			return result, fmt.Errorf("provenance history signal is not recognized")
		}
		result.DispositionSignal = fmt.Sprintf("provenance-%s-plan", plan.DecisionSignal)
	} else {
		result.DispositionSignal = "provenance-plan-mismatch"
	}
	if !plan.PlanObserved {
		result.MissingStage = "revision-self-improvement-provenance-plan-disposition-observed"
		setObservationDigest()
		return result, fmt.Errorf("execution plan was not observed")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setObservationDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-provenance-plan-disposition"
		result.DispositionSignal = "provenance-plan-unknown"
		setObservationDigest()
		return result, fmt.Errorf("provenance plan disposition is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementProvenancePlanDispositionObservation) Validate() error {
	if o.Status == "" {
		return fmt.Errorf("provenance plan disposition status is empty")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound provenance plan disposition has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown provenance plan disposition has no missing stage")
	}
	for name, digest := range map[string]string{
		"iteration provenance": o.IterationProvenanceDigest,
		"plan observation":     o.PlanObservationDigest,
		"plan":                 o.PlanDigest,
		"decision":             o.DecisionDigest,
		"current source":       o.CurrentSourceDigest,
		"candidate source":     o.CandidateSourceDigest,
		"plan source":          o.PlanSourceDigest,
		"candidate proposed":   o.CandidateProposedSourceDigest,
		"candidate generated":  o.CandidateGeneratedIRDigest,
		"plan input IR":        o.PlanInputIRDigest,
		"observation":          o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("provenance plan disposition %s digest is invalid", name)
		}
	}
	if o.ProvenanceHistorySignal != "stable" &&
		o.ProvenanceHistorySignal != "transitioned" &&
		o.ProvenanceHistorySignal != "mixed" {
		return fmt.Errorf("provenance plan disposition history signal is invalid")
	}
	if o.ProvenanceDecisionSignal != "observe" && o.ProvenanceDecisionSignal != "inspect" {
		return fmt.Errorf("provenance plan disposition provenance decision signal is invalid")
	}
	if o.DecisionSignal != "observe" && o.DecisionSignal != "remeasure" &&
		o.DecisionSignal != "review" && o.DecisionSignal != "inspect" {
		return fmt.Errorf("provenance plan disposition decision signal is invalid")
	}
	if o.PlanningSignal != fmt.Sprintf("%s-plan", o.DecisionSignal) {
		return fmt.Errorf("provenance plan disposition planning signal is invalid")
	}
	if o.PlanSignal == "" {
		return fmt.Errorf("provenance plan disposition plan signal is empty")
	}
	if o.DispositionSignal != "provenance-observe-plan" &&
		o.DispositionSignal != "provenance-inspect-plan" &&
		o.DispositionSignal != "provenance-plan-mismatch" &&
		o.DispositionSignal != "provenance-plan-unknown" {
		return fmt.Errorf("provenance plan disposition signal is invalid")
	}
	expectedAligned := o.ProvenanceDecisionSignal == o.DecisionSignal &&
		o.ProvenanceRequiresObservation == o.RequiresObservation &&
		o.ProvenanceRequiresInspection == o.RequiresInspection &&
		o.ProvenanceRequiresMeasurement == o.RequiresMeasurement
	if o.SignalsAligned != expectedAligned {
		return fmt.Errorf("provenance plan disposition alignment is inconsistent")
	}
	if o.SignalsAligned && o.DispositionSignal != fmt.Sprintf("provenance-%s-plan", o.DecisionSignal) {
		return fmt.Errorf("aligned provenance plan disposition signal is invalid")
	}
	if !o.SignalsAligned && o.DispositionSignal != "provenance-plan-mismatch" {
		return fmt.Errorf("mismatched provenance plan disposition signal is invalid")
	}
	if !o.PlanObserved || !o.RequiresObservation {
		return fmt.Errorf("provenance plan disposition must preserve observed required planning")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("provenance plan disposition must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementProvenancePlanDisposition(o) != o.ObservationDigest {
		return fmt.Errorf("provenance plan disposition digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementProvenancePlanDisposition(
	observation RevisionSelfImprovementProvenancePlanDispositionObservation,
) string {
	parts := []string{
		observation.Status,
		observation.MissingStage,
		observation.IterationProvenanceDigest,
		observation.PlanObservationDigest,
		observation.PlanDigest,
		observation.DecisionDigest,
		observation.CurrentSourceDigest,
		observation.CandidateSourceDigest,
		observation.PlanSourceDigest,
		observation.CandidateProposedSourceDigest,
		observation.CandidateGeneratedIRDigest,
		observation.PlanInputIRDigest,
		observation.ProvenanceHistorySignal,
		observation.ProvenanceDecisionSignal,
		observation.DecisionSignal,
		observation.PlanningSignal,
		observation.PlanSignal,
		observation.DispositionSignal,
		strconv.FormatBool(observation.PlanObserved),
		strconv.FormatBool(observation.ProvenanceRequiresObservation),
		strconv.FormatBool(observation.ProvenanceRequiresReview),
		strconv.FormatBool(observation.ProvenanceRequiresInspection),
		strconv.FormatBool(observation.ProvenanceRequiresMeasurement),
		strconv.FormatBool(observation.SignalsAligned),
		strconv.FormatBool(observation.RequiresObservation),
		strconv.FormatBool(observation.RequiresReview),
		strconv.FormatBool(observation.RequiresInspection),
		strconv.FormatBool(observation.RequiresMeasurement),
		strconv.FormatBool(observation.NonExecuting),
		strconv.FormatBool(observation.NonAuthorizing),
	}
	return digestString(strings.Join(parts, "|"))
}
