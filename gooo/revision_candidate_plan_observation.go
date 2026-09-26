package gooo

import "fmt"

// RevisionCandidatePlanObservation links a candidate observation to an
// application plan while remaining non-executing and non-authorizing.
type RevisionCandidatePlanObservation struct {
	Status                    string
	MissingStage              string
	CandidateObservationDigest string
	CandidateDigest           string
	SourceDigest              string
	MaterializationDigest     string
	PlanDigest                string
	BindingDigest             string
	EditDigest                string
	ObservationSignal         string
	PlanSignal                string
	PlanObserved              bool
	ObservationDigest         string
	NonExecuting              bool
	NonAuthorizing            bool
}

// ObserveRevisionCandidatePlan binds a candidate observation to a validated
// application plan without applying or authorizing that plan.
func ObserveRevisionCandidatePlan(observation RevisionCandidateObservation, plan RevisionApplicationPlan) (RevisionCandidatePlanObservation, error) {
	result := RevisionCandidatePlanObservation{
		Status:                     "UNKNOWN",
		MissingStage:               "revision-candidate-plan-observation",
		CandidateObservationDigest: observation.ResultDigest,
		CandidateDigest:            observation.CandidateDigest,
		SourceDigest:               observation.SourceDigest,
		MaterializationDigest:      observation.MaterializationDigest,
		PlanDigest:                 plan.PlanDigest,
		BindingDigest:              plan.BindingDigest,
		EditDigest:                 plan.EditDigest,
		ObservationSignal:          observation.MaterializationSignal,
		NonExecuting:               true,
		NonAuthorizing:             true,
	}
	setObservationDigest := func() {
		result.ObservationDigest = digestRevisionCandidatePlanObservation(result)
	}
	setObservationDigest()

	if err := observation.Validate(); err != nil {
		result.MissingStage = "revision-candidate-plan-observation-candidate"
		setObservationDigest()
		return result, fmt.Errorf("revision candidate observation is not valid: %w", err)
	}
	if err := plan.Validate(); err != nil {
		result.MissingStage = "revision-candidate-plan-observation-plan"
		setObservationDigest()
		return result, fmt.Errorf("revision application plan is not valid: %w", err)
	}
	if observation.CandidateDigest != plan.CandidateDigest || observation.SourceDigest != plan.SourceDigest {
		result.MissingStage = "revision-candidate-plan-observation-link"
		setObservationDigest()
		return result, fmt.Errorf("candidate observation and application plan are not linked")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	result.PlanObserved = true
	if observation.MaterializationSignal == "inspect-candidate" {
		result.PlanSignal = "inspect-plan"
	} else {
		result.PlanSignal = "plan-observed"
	}
	setObservationDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-candidate-plan-observation"
		setObservationDigest()
		return result, fmt.Errorf("revision candidate plan observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionCandidatePlanObservation) Validate() error {
	if o.Status == "" {
		return fmt.Errorf("revision candidate plan observation status is empty")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound revision candidate plan observation has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown revision candidate plan observation has no missing stage")
	}
	if o.Status == "BOUND" && (!o.PlanObserved || o.PlanSignal == "") {
		return fmt.Errorf("bound revision candidate plan observation is incomplete")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("revision candidate plan observation must remain non-executing and non-authorizing")
	}
	if digestRevisionCandidatePlanObservation(o) != o.ObservationDigest {
		return fmt.Errorf("revision candidate plan observation digest does not match its fields")
	}
	return nil
}

func digestRevisionCandidatePlanObservation(observation RevisionCandidatePlanObservation) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t|%t|%t",
		observation.Status,
		observation.MissingStage,
		observation.CandidateObservationDigest,
		observation.CandidateDigest,
		observation.SourceDigest,
		observation.MaterializationDigest,
		observation.PlanDigest,
		observation.BindingDigest,
		observation.EditDigest,
		observation.ObservationSignal,
		observation.PlanSignal,
		observation.PlanObserved,
		observation.NonExecuting,
		observation.NonAuthorizing,
		observation.PlanObserved,
	))
}