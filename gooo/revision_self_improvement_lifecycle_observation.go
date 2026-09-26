package gooo

import "fmt"

// RevisionSelfImprovementLifecycleObservation binds planning, application,
// measured change, and reverse-generation evidence without executing it.
type RevisionSelfImprovementLifecycleObservation struct {
	Status                                string
	MissingStage                          string
	ExecutionApplicationObservationDigest string
	OutcomeObservationDigest              string
	ApplicationObservationDigest         string
	MetricsBindingDigest                 string
	GenerationAssessmentDigest           string
	PlanObservationDigest                string
	PlanDigest                            string
	ApplicationDigest                     string
	ReceiptDigest                         string
	SourceDigest                          string
	ProposedSourceDigest                  string
	InputIRDigest                         string
	ProposedIRDigest                      string
	CandidateDigest                       string
	EditDigest                            string
	ChangedByteCount                      int
	ChangedLineCount                      int
	IRChanged                             bool
	GeneratedSourceDigest                 string
	GeneratedIRDigest                     string
	StructureDigest                       string
	ExactSourceMatch                      bool
	StructureMatch                        bool
	DecisionSignal                        string
	FeedbackSignal                        string
	ApplicationSignal                     string
	PlanObserved                          bool
	ApplicationObserved                   bool
	LifecycleSignal                       string
	ObservationDigest                     string
	NonExecuting                          bool
	NonAuthorizing                        bool
}

// ObserveRevisionSelfImprovementLifecycle binds execution/application
// provenance to metrics and reverse-generation outcome evidence.
func ObserveRevisionSelfImprovementLifecycle(executionApplication RevisionSelfImprovementExecutionApplicationObservation, outcome RevisionSelfImprovementOutcomeObservation) (RevisionSelfImprovementLifecycleObservation, error) {
	result := RevisionSelfImprovementLifecycleObservation{
		Status:                                "UNKNOWN",
		MissingStage:                          "revision-self-improvement-lifecycle",
		ExecutionApplicationObservationDigest: executionApplication.ObservationDigest,
		OutcomeObservationDigest:              outcome.OutcomeDigest,
		ApplicationObservationDigest:          executionApplication.ApplicationObservationDigest,
		MetricsBindingDigest:                  outcome.MetricsBindingDigest,
		GenerationAssessmentDigest:            outcome.GenerationAssessmentDigest,
		PlanObservationDigest:                 executionApplication.PlanObservationDigest,
		PlanDigest:                            executionApplication.PlanDigest,
		ApplicationDigest:                     outcome.ApplicationDigest,
		ReceiptDigest:                         executionApplication.ReceiptDigest,
		SourceDigest:                          outcome.SourceDigest,
		ProposedSourceDigest:                  outcome.ProposedSourceDigest,
		InputIRDigest:                         outcome.InputIRDigest,
		ProposedIRDigest:                      outcome.ProposedIRDigest,
		CandidateDigest:                       outcome.CandidateDigest,
		EditDigest:                            outcome.EditDigest,
		ChangedByteCount:                      outcome.ChangedByteCount,
		ChangedLineCount:                      outcome.ChangedLineCount,
		IRChanged:                             outcome.IRChanged,
		GeneratedSourceDigest:                 outcome.GeneratedSourceDigest,
		GeneratedIRDigest:                     outcome.GeneratedIRDigest,
		StructureDigest:                       outcome.StructureDigest,
		ExactSourceMatch:                      outcome.ExactSourceMatch,
		StructureMatch:                        outcome.StructureMatch,
		DecisionSignal:                        executionApplication.DecisionSignal,
		FeedbackSignal:                        executionApplication.FeedbackSignal,
		ApplicationSignal:                     executionApplication.ApplicationSignal,
		PlanObserved:                          executionApplication.PlanObserved,
		ApplicationObserved:                   executionApplication.ApplicationObserved,
		LifecycleSignal:                       "plan-application-outcome-bound",
		NonExecuting:                          true,
		NonAuthorizing:                        true,
	}
	setObservationDigest := func() {
		result.ObservationDigest = digestRevisionSelfImprovementLifecycle(result)
	}
	setObservationDigest()

	if err := executionApplication.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-lifecycle-execution-application"
		setObservationDigest()
		return result, fmt.Errorf("execution application observation is not valid: %w", err)
	}
	if err := outcome.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-lifecycle-outcome"
		setObservationDigest()
		return result, fmt.Errorf("self-improvement outcome is not valid: %w", err)
	}
	if executionApplication.SourceDigest != outcome.SourceDigest ||
		executionApplication.ProposedSourceDigest != outcome.ProposedSourceDigest {
		result.MissingStage = "revision-self-improvement-lifecycle-source-link"
		setObservationDigest()
		return result, fmt.Errorf("execution application and outcome sources are not linked")
	}
	if outcome.ApplicationObservationDigest != executionApplication.ApplicationObservationDigest ||
		outcome.ApplicationDigest != executionApplication.ApplicationDigest ||
		outcome.InputIRDigest != executionApplication.InputIRDigest ||
		outcome.ProposedIRDigest != executionApplication.ProposedIRDigest ||
		outcome.CandidateDigest != executionApplication.CandidateDigest ||
		outcome.EditDigest != executionApplication.EditDigest ||
		outcome.SourceDigest != executionApplication.SourceDigest ||
		outcome.ProposedSourceDigest != executionApplication.ProposedSourceDigest {
		result.MissingStage = "revision-self-improvement-lifecycle-application-link"
		setObservationDigest()
		return result, fmt.Errorf("execution application and outcome application fields are not linked")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setObservationDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-lifecycle"
		setObservationDigest()
		return result, fmt.Errorf("revision self-improvement lifecycle observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementLifecycleObservation) Validate() error {
	if o.Status == "" {
		return fmt.Errorf("revision self-improvement lifecycle status is empty")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound revision self-improvement lifecycle has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown revision self-improvement lifecycle has no missing stage")
	}
	for name, digest := range map[string]string{
		"execution application":   o.ExecutionApplicationObservationDigest,
		"outcome":                 o.OutcomeObservationDigest,
		"application observation": o.ApplicationObservationDigest,
		"metrics binding":         o.MetricsBindingDigest,
		"generation assessment":   o.GenerationAssessmentDigest,
		"plan observation":        o.PlanObservationDigest,
		"plan":                    o.PlanDigest,
		"application":             o.ApplicationDigest,
		"receipt":                 o.ReceiptDigest,
		"source":                  o.SourceDigest,
		"proposed source":         o.ProposedSourceDigest,
		"input IR":                o.InputIRDigest,
		"proposed IR":             o.ProposedIRDigest,
		"candidate":               o.CandidateDigest,
		"edit":                    o.EditDigest,
		"generated source":        o.GeneratedSourceDigest,
		"generated IR":            o.GeneratedIRDigest,
		"structure":               o.StructureDigest,
		"observation":             o.ObservationDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("revision self-improvement lifecycle %s digest is invalid", name)
		}
	}
	if o.ChangedByteCount < 0 || o.ChangedLineCount < 0 {
		return fmt.Errorf("revision self-improvement lifecycle counts must be non-negative")
	}
	if o.DecisionSignal != "observe" && o.DecisionSignal != "remeasure" &&
		o.DecisionSignal != "review" && o.DecisionSignal != "inspect" {
		return fmt.Errorf("revision self-improvement lifecycle decision signal is invalid")
	}
	if o.FeedbackSignal != o.DecisionSignal {
		return fmt.Errorf("revision self-improvement lifecycle feedback signal is not linked")
	}
	if o.ApplicationSignal == "" {
		return fmt.Errorf("revision self-improvement lifecycle application signal is empty")
	}
	if !o.PlanObserved || !o.ApplicationObserved {
		return fmt.Errorf("revision self-improvement lifecycle must preserve plan and application observations")
	}
	if o.LifecycleSignal != "plan-application-outcome-bound" {
		return fmt.Errorf("revision self-improvement lifecycle signal is invalid")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("revision self-improvement lifecycle must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementLifecycle(o) != o.ObservationDigest {
		return fmt.Errorf("revision self-improvement lifecycle digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementLifecycle(observation RevisionSelfImprovementLifecycleObservation) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%d|%d|%t|%t|%t|%t|%t|%t|%t|%s|%t|%t",
		observation.Status,
		observation.MissingStage,
		observation.ExecutionApplicationObservationDigest,
		observation.OutcomeObservationDigest,
		observation.ApplicationObservationDigest,
		observation.MetricsBindingDigest,
		observation.GenerationAssessmentDigest,
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
		observation.GeneratedSourceDigest,
		observation.GeneratedIRDigest,
		observation.StructureDigest,
		observation.DecisionSignal,
		observation.FeedbackSignal,
		observation.ApplicationSignal,
		observation.ChangedByteCount,
		observation.ChangedLineCount,
		observation.IRChanged,
		observation.ExactSourceMatch,
		observation.StructureMatch,
		observation.PlanObserved,
		observation.ApplicationObserved,
		observation.NonExecuting,
		observation.NonAuthorizing,
		observation.LifecycleSignal,
		observation.NonExecuting,
		observation.NonAuthorizing,
	))
}
