package gooo

import "fmt"

// RevisionSelfImprovementApplicationObservation links an observed candidate
// plan application to the current iteration without executing or authorizing it.
type RevisionSelfImprovementApplicationObservation struct {
	Status                    string
	MissingStage              string
	IterationDigest            string
	PlanObservationDigest      string
	ApplicationObservationDigest string
	PlanDigest                string
	ApplicationDigest          string
	ReceiptDigest              string
	SourceDigest               string
	ProposedSourceDigest       string
	InputIRDigest              string
	ProposedIRDigest           string
	CandidateDigest            string
	EditDigest                string
	FeedbackSignal             string
	ApplicationSignal          string
	ApplicationObserved        bool
	ObservationDigest          string
	NonExecuting               bool
	NonAuthorizing             bool
}

// ObserveRevisionSelfImprovementApplication links an observed application
// receipt to a candidate plan and current iteration without applying anything.
func ObserveRevisionSelfImprovementApplication(iteration RevisionSelfImprovementIteration, planObservation RevisionCandidatePlanObservation, applicationObservation RevisionApplicationObservation) (RevisionSelfImprovementApplicationObservation, error) {
	result := RevisionSelfImprovementApplicationObservation{
		Status:                      "UNKNOWN",
		MissingStage:                "revision-self-improvement-application",
		IterationDigest:             iteration.IterationDigest,
		PlanObservationDigest:       applicationObservation.PlanObservationDigest,
		ApplicationObservationDigest: applicationObservation.ObservationDigest,
		PlanDigest:                  applicationObservation.PlanDigest,
		ApplicationDigest:           applicationObservation.ApplicationDigest,
		ReceiptDigest:               applicationObservation.ReceiptDigest,
		SourceDigest:                applicationObservation.SourceDigest,
		ProposedSourceDigest:        applicationObservation.ProposedSourceDigest,
		InputIRDigest:               applicationObservation.InputIRDigest,
		ProposedIRDigest:            applicationObservation.ProposedIRDigest,
		CandidateDigest:             applicationObservation.CandidateDigest,
		EditDigest:                 applicationObservation.EditDigest,
		FeedbackSignal:              iteration.FeedbackSignal,
		ApplicationSignal:           applicationObservation.ApplicationSignal,
		ApplicationObserved:         applicationObservation.ApplicationObserved,
		NonExecuting:               true,
		NonAuthorizing:             true,
	}
	setObservationDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementApplicationObservation(result)
	}
	setObservationDigest()

	if err := iteration.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-application-iteration"
		setObservationDigest()
		return result, fmt.Errorf("self-improvement iteration is not valid: %w", err)
	}
	if err := planObservation.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-application-plan"
		setObservationDigest()
		return result, fmt.Errorf("candidate plan observation is not valid: %w", err)
	}
	if err := applicationObservation.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-application-observation"
		setObservationDigest()
		return result, fmt.Errorf("application observation is not valid: %w", err)
	}
	if planObservation.ObservationDigest != applicationObservation.PlanObservationDigest ||
		planObservation.PlanDigest != applicationObservation.PlanDigest ||
		planObservation.SourceDigest != applicationObservation.SourceDigest ||
		planObservation.CandidateDigest != applicationObservation.CandidateDigest ||
		planObservation.EditDigest != applicationObservation.EditDigest {
		result.MissingStage = "revision-self-improvement-application-plan-link"
		setObservationDigest()
		return result, fmt.Errorf("candidate plan and application observation are not linked")
	}
	if iteration.CandidateSourceDigest != applicationObservation.SourceDigest {
		result.MissingStage = "revision-self-improvement-application-source-link"
		setObservationDigest()
		return result, fmt.Errorf("application source does not match the iteration candidate source")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setObservationDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-application"
		setObservationDigest()
		return result, fmt.Errorf("revision self-improvement application observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementApplicationObservation) Validate() error {
	if o.Status == "" {
		return fmt.Errorf("revision self-improvement application status is empty")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound revision self-improvement application has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown revision self-improvement application has no missing stage")
	}
	if !o.ApplicationObserved {
		return fmt.Errorf("revision self-improvement application has not been observed")
	}
	if o.FeedbackSignal != "observe" && o.FeedbackSignal != "remeasure" &&
		o.FeedbackSignal != "review" && o.FeedbackSignal != "inspect" {
		return fmt.Errorf("revision self-improvement application feedback signal is invalid")
	}
	if o.ApplicationSignal == "" {
		return fmt.Errorf("revision self-improvement application signal is empty")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("revision self-improvement application must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementApplicationObservation(o) != o.ObservationDigest {
		return fmt.Errorf("revision self-improvement application digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementApplicationObservation(observation RevisionSelfImprovementApplicationObservation) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t|%t",
		observation.Status,
		observation.MissingStage,
		observation.IterationDigest,
		observation.PlanObservationDigest,
		observation.ApplicationObservationDigest,
		observation.PlanDigest,
		observation.ApplicationDigest,
		observation.ReceiptDigest,
		observation.SourceDigest,
		observation.ProposedSourceDigest,
		observation.InputIRDigest,
		observation.ProposedIRDigest,
		observation.CandidateDigest,
		observation.EditDigest,
		observation.FeedbackSignal,
		observation.ApplicationSignal,
		observation.ApplicationObserved,
		observation.NonExecuting,
		observation.NonAuthorizing,
	))
}