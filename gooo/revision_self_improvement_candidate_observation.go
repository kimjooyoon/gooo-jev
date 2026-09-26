package gooo

import "fmt"

// RevisionSelfImprovementCandidateObservation links the next candidate plan to
// a validated iteration without applying or authorizing that plan.
type RevisionSelfImprovementCandidateObservation struct {
	Status                       string
	MissingStage                 string
	IterationDigest              string
	CandidatePlanObservationDigest string
	CandidateDigest              string
	SourceDigest                 string
	MaterializationDigest        string
	PlanDigest                   string
	BindingDigest                string
	EditDigest                  string
	FeedbackSignal               string
	ObservationSignal            string
	PlanSignal                   string
	CandidatePlanObserved        bool
	ObservationDigest            string
	NonExecuting                 bool
	NonAuthorizing               bool
}

// ObserveRevisionSelfImprovementCandidate binds a candidate plan observation
// to the current self-improvement iteration without executing or authorizing it.
func ObserveRevisionSelfImprovementCandidate(iteration RevisionSelfImprovementIteration, planObservation RevisionCandidatePlanObservation) (RevisionSelfImprovementCandidateObservation, error) {
	result := RevisionSelfImprovementCandidateObservation{
		Status:                         "UNKNOWN",
		MissingStage:                   "revision-self-improvement-candidate",
		IterationDigest:                iteration.IterationDigest,
		CandidatePlanObservationDigest: planObservation.ObservationDigest,
		CandidateDigest:                planObservation.CandidateDigest,
		SourceDigest:                   planObservation.SourceDigest,
		MaterializationDigest:          planObservation.MaterializationDigest,
		PlanDigest:                     planObservation.PlanDigest,
		BindingDigest:                  planObservation.BindingDigest,
		EditDigest:                    planObservation.EditDigest,
		FeedbackSignal:                iteration.FeedbackSignal,
		ObservationSignal:             planObservation.ObservationSignal,
		PlanSignal:                    planObservation.PlanSignal,
		NonExecuting:                  true,
		NonAuthorizing:                true,
	}
	setObservationDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementCandidateObservation(result)
	}
	setObservationDigest()

	if err := iteration.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-candidate-iteration"
		setObservationDigest()
		return result, fmt.Errorf("self-improvement iteration is not valid: %w", err)
	}
	if err := planObservation.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-candidate-plan"
		setObservationDigest()
		return result, fmt.Errorf("candidate plan observation is not valid: %w", err)
	}
	if iteration.CandidateSourceDigest != planObservation.SourceDigest {
		result.MissingStage = "revision-self-improvement-candidate-source-link"
		setObservationDigest()
		return result, fmt.Errorf("candidate plan source does not match the iteration candidate source")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.CandidatePlanObserved = true
	setObservationDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-candidate"
		setObservationDigest()
		return result, fmt.Errorf("revision self-improvement candidate observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementCandidateObservation) Validate() error {
	if o.Status == "" {
		return fmt.Errorf("revision self-improvement candidate status is empty")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound revision self-improvement candidate has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown revision self-improvement candidate has no missing stage")
	}
	if o.Status == "BOUND" && !o.CandidatePlanObserved {
		return fmt.Errorf("bound revision self-improvement candidate has no observed plan")
	}
	if o.FeedbackSignal != "observe" && o.FeedbackSignal != "remeasure" &&
		o.FeedbackSignal != "review" && o.FeedbackSignal != "inspect" {
		return fmt.Errorf("revision self-improvement candidate feedback signal is invalid")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("revision self-improvement candidate must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementCandidateObservation(o) != o.ObservationDigest {
		return fmt.Errorf("revision self-improvement candidate observation digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementCandidateObservation(observation RevisionSelfImprovementCandidateObservation) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t|%t",
		observation.Status,
		observation.MissingStage,
		observation.IterationDigest,
		observation.CandidatePlanObservationDigest,
		observation.CandidateDigest,
		observation.SourceDigest,
		observation.MaterializationDigest,
		observation.PlanDigest,
		observation.BindingDigest,
		observation.EditDigest,
		observation.FeedbackSignal,
		observation.ObservationSignal,
		observation.PlanSignal,
		observation.CandidatePlanObserved,
		observation.NonExecuting,
		observation.NonAuthorizing,
	))
}