package gooo

import "fmt"

// RevisionSelfImprovementExecutionApplicationObservation binds an observed
// execution plan to its application provenance without executing it.
type RevisionSelfImprovementExecutionApplicationObservation struct {
	Status                           string
	MissingStage                     string
	ExecutionPlanObservationDigest   string
	ApplicationObservationDigest    string
	PlanObservationDigest            string
	PlanDigest                       string
	ApplicationDigest                string
	ReceiptDigest                    string
	SourceDigest                     string
	ProposedSourceDigest             string
	InputIRDigest                    string
	ProposedIRDigest                 string
	CandidateDigest                  string
	EditDigest                       string
	DecisionSignal                   string
	FeedbackSignal                   string
	ApplicationSignal                string
	PlanObserved                     bool
	ApplicationObserved              bool
	ObservationDigest                string
	NonExecuting                     bool
	NonAuthorizing                   bool
}

// ObserveRevisionSelfImprovementExecutionApplication binds a plan observation
// to a validated application observation without applying or authorizing it.
func ObserveRevisionSelfImprovementExecutionApplication(plan RevisionSelfImprovementExecutionPlanObservation, application RevisionSelfImprovementApplicationObservation) (RevisionSelfImprovementExecutionApplicationObservation, error) {
	result := RevisionSelfImprovementExecutionApplicationObservation{
		Status:                         "UNKNOWN",
		MissingStage:                   "revision-self-improvement-execution-application",
		ExecutionPlanObservationDigest: plan.ObservationDigest,
		ApplicationObservationDigest:   application.ObservationDigest,
		PlanObservationDigest:          application.PlanObservationDigest,
		PlanDigest:                     application.PlanDigest,
		ApplicationDigest:              application.ApplicationDigest,
		ReceiptDigest:                  application.ReceiptDigest,
		SourceDigest:                   application.SourceDigest,
		ProposedSourceDigest:           application.ProposedSourceDigest,
		InputIRDigest:                  application.InputIRDigest,
		ProposedIRDigest:               application.ProposedIRDigest,
		CandidateDigest:                application.CandidateDigest,
		EditDigest:                    application.EditDigest,
		DecisionSignal:                 plan.DecisionSignal,
		FeedbackSignal:                 application.FeedbackSignal,
		ApplicationSignal:              application.ApplicationSignal,
		PlanObserved:                   plan.PlanObserved,
		ApplicationObserved:             application.ApplicationObserved,
		NonExecuting:                  true,
		NonAuthorizing:                true,
	}
	setObservationDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementExecutionApplication(result)
	}
	setObservationDigest()

	if err := plan.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-execution-application-plan"
		setObservationDigest()
		return result, fmt.Errorf("self-improvement execution plan is not valid: %w", err)
	}
	if err := application.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-execution-application-observation"
		setObservationDigest()
		return result, fmt.Errorf("self-improvement application observation is not valid: %w", err)
	}
	if plan.SourceDigest != application.SourceDigest {
		result.MissingStage = "revision-self-improvement-execution-application-source-link"
		setObservationDigest()
		return result, fmt.Errorf("execution plan and application observation sources are not linked")
	}
	if plan.PlanObservationDigest != application.PlanObservationDigest ||
		plan.PlanDigest != application.PlanDigest ||
		plan.InputIRDigest != application.InputIRDigest ||
		plan.CandidateDigest != application.CandidateDigest ||
		plan.EditDigest != application.EditDigest {
		result.MissingStage = "revision-self-improvement-execution-application-plan-link"
		setObservationDigest()
		return result, fmt.Errorf("execution plan and application observation plan fields are not linked")
	}
	if plan.DecisionSignal != application.FeedbackSignal {
		result.MissingStage = "revision-self-improvement-execution-application-signal-link"
		setObservationDigest()
		return result, fmt.Errorf("execution plan decision signal and application feedback signal are not linked")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setObservationDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-execution-application"
		setObservationDigest()
		return result, fmt.Errorf("revision self-improvement execution application observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementExecutionApplicationObservation) Validate() error {
	if o.Status == "" {
		return fmt.Errorf("revision self-improvement execution application status is empty")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound revision self-improvement execution application has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown revision self-improvement execution application has no missing stage")
	}
	for name, digest := range map[string]string{
		"execution plan observation": o.ExecutionPlanObservationDigest,
		"application observation":    o.ApplicationObservationDigest,
		"plan observation":           o.PlanObservationDigest,
		"plan":                       o.PlanDigest,
		"application":                o.ApplicationDigest,
		"receipt":                    o.ReceiptDigest,
		"source":                     o.SourceDigest,
		"proposed source":            o.ProposedSourceDigest,
		"input IR":                   o.InputIRDigest,
		"proposed IR":                o.ProposedIRDigest,
		"candidate":                  o.CandidateDigest,
		"edit":                       o.EditDigest,
		"observation":                o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("revision self-improvement execution application %s digest is invalid", name)
		}
	}
	if !o.PlanObserved || !o.ApplicationObserved {
		return fmt.Errorf("revision self-improvement execution application must preserve observations")
	}
	if o.DecisionSignal != "observe" && o.DecisionSignal != "remeasure" &&
		o.DecisionSignal != "review" && o.DecisionSignal != "inspect" {
		return fmt.Errorf("revision self-improvement execution application decision signal is invalid")
	}
	if o.FeedbackSignal != o.DecisionSignal {
		return fmt.Errorf("revision self-improvement execution application feedback signal is not linked")
	}
	if o.ApplicationSignal == "" {
		return fmt.Errorf("revision self-improvement execution application signal is empty")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("revision self-improvement execution application must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementExecutionApplication(o) != o.ObservationDigest {
		return fmt.Errorf("revision self-improvement execution application digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementExecutionApplication(observation RevisionSelfImprovementExecutionApplicationObservation) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t|%t|%t",
		observation.Status,
		observation.MissingStage,
		observation.ExecutionPlanObservationDigest,
		observation.ApplicationObservationDigest,
		observation.PlanObservationDigest,
		observation.PlanDigest,
		observation.ApplicationDigest,
		observation.ReceiptDigest,
		observation.SourceDigest,
		observation.ProposedSourceDigest,
		observation.InputIRDigest,
		observation.ProposedIRDigest,
		observation.CandidateDigest,
		observation.EditDigest,
		observation.DecisionSignal,
		observation.FeedbackSignal,
		observation.ApplicationSignal,
		observation.PlanObserved,
		observation.ApplicationObserved,
		observation.NonExecuting,
		observation.NonAuthorizing,
	))
}
