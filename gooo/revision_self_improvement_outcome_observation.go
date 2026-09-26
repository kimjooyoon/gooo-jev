package gooo

import "fmt"

// RevisionSelfImprovementOutcomeObservation binds application, metrics, and
// generation evidence for one iteration without executing or authorizing it.
type RevisionSelfImprovementOutcomeObservation struct {
	Status                       string
	MissingStage                 string
	IterationDigest              string
	ApplicationObservationDigest string
	MetricsBindingDigest         string
	GenerationAssessmentDigest   string
	SourceDigest                 string
	ProposedSourceDigest         string
	InputIRDigest                string
	ProposedIRDigest             string
	CandidateDigest              string
	EditDigest                  string
	ApplicationDigest            string
	MetricsDigest                string
	GeneratedSourceDigest        string
	GeneratedIRDigest            string
	StructureDigest              string
	ChangedByteCount             int
	ChangedLineCount             int
	IRChanged                    bool
	ExactSourceMatch              bool
	StructureMatch                bool
	OutcomeSignal                 string
	OutcomeDigest                 string
	NonExecuting                  bool
	NonAuthorizing                bool
}

// ObserveRevisionSelfImprovementOutcome binds the observed application,
// measured change, and reverse generation evidence to one iteration.
func ObserveRevisionSelfImprovementOutcome(iteration RevisionSelfImprovementIteration, application RevisionSelfImprovementApplicationObservation, metrics RevisionMetricsBinding, generation RevisionGenerationAssessment) (RevisionSelfImprovementOutcomeObservation, error) {
	result := RevisionSelfImprovementOutcomeObservation{
		Status:                       "UNKNOWN",
		MissingStage:                 "revision-self-improvement-outcome",
		IterationDigest:              iteration.IterationDigest,
		ApplicationObservationDigest: application.ObservationDigest,
		MetricsBindingDigest:         metrics.BindingDigest,
		GenerationAssessmentDigest:   generation.GenerationDigest,
		SourceDigest:                 application.SourceDigest,
		ProposedSourceDigest:         application.ProposedSourceDigest,
		InputIRDigest:                application.InputIRDigest,
		ProposedIRDigest:             application.ProposedIRDigest,
		CandidateDigest:              application.CandidateDigest,
		EditDigest:                   application.EditDigest,
		ApplicationDigest:            application.ApplicationDigest,
		MetricsDigest:                metrics.MetricsDigest,
		GeneratedSourceDigest:        generation.GeneratedSourceDigest,
		GeneratedIRDigest:             generation.GeneratedIRDigest,
		StructureDigest:              generation.StructureDigest,
		ChangedByteCount:             metrics.ChangedByteCount,
		ChangedLineCount:             metrics.ChangedLineCount,
		IRChanged:                    metrics.IRChanged,
		ExactSourceMatch:              generation.ExactSourceMatch,
		StructureMatch:               generation.StructureMatch,
		OutcomeSignal:                "metrics-generation-bound",
		NonExecuting:                 true,
		NonAuthorizing:               true,
	}
	setOutcomeDigest := func() {
		result.OutcomeDigest = digestRevisionSelfImprovementOutcomeObservation(result)
	}
	setOutcomeDigest()

	if err := iteration.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-outcome-iteration"
		setOutcomeDigest()
		return result, fmt.Errorf("self-improvement iteration is not valid: %w", err)
	}
	if err := application.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-outcome-application"
		setOutcomeDigest()
		return result, fmt.Errorf("self-improvement application observation is not valid: %w", err)
	}
	if err := metrics.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-outcome-metrics"
		setOutcomeDigest()
		return result, fmt.Errorf("revision metrics binding is not valid: %w", err)
	}
	if err := generation.Validate(); err != nil {
		result.MissingStage = "revision-self-improvement-outcome-generation"
		setOutcomeDigest()
		return result, fmt.Errorf("revision generation assessment is not valid: %w", err)
	}
	if iteration.CandidateSourceDigest != application.SourceDigest ||
		iteration.CandidateProposedSourceDigest != application.ProposedSourceDigest {
		result.MissingStage = "revision-self-improvement-outcome-iteration-link"
		setOutcomeDigest()
		return result, fmt.Errorf("iteration candidate sources are not linked to the application observation")
	}
	if application.ReceiptDigest != metrics.ReceiptDigest ||
		application.ApplicationDigest != metrics.ApplicationDigest ||
		application.SourceDigest != metrics.SourceDigest ||
		application.ProposedSourceDigest != metrics.ProposedSourceDigest ||
		application.InputIRDigest != metrics.InputIRDigest ||
		application.ProposedIRDigest != metrics.ProposedIRDigest ||
		application.CandidateDigest != metrics.CandidateDigest ||
		application.EditDigest != metrics.EditDigest {
		result.MissingStage = "revision-self-improvement-outcome-application-metrics-link"
		setOutcomeDigest()
		return result, fmt.Errorf("application observation and metrics binding are not linked")
	}
	if generation.ApplicationDigest != application.ApplicationDigest ||
		generation.SourceDigest != application.SourceDigest ||
		generation.ProposedSourceDigest != application.ProposedSourceDigest ||
		generation.GenerationSourceDigest != application.ProposedSourceDigest {
		result.MissingStage = "revision-self-improvement-outcome-generation-link"
		setOutcomeDigest()
		return result, fmt.Errorf("generation assessment is not linked to the proposed application source")
	}

	result.Status = "BOUND"
	result.MissingStage = ""
	setOutcomeDigest()
	if err := result.Validate(); err != nil {
		result.Status = "UNKNOWN"
		result.MissingStage = "revision-self-improvement-outcome"
		setOutcomeDigest()
		return result, fmt.Errorf("revision self-improvement outcome observation is not valid: %w", err)
	}
	return result, nil
}

func (o RevisionSelfImprovementOutcomeObservation) Validate() error {
	if o.Status == "" {
		return fmt.Errorf("revision self-improvement outcome status is empty")
	}
	if o.Status == "BOUND" && o.MissingStage != "" {
		return fmt.Errorf("bound revision self-improvement outcome has a missing stage")
	}
	if o.Status == "UNKNOWN" && o.MissingStage == "" {
		return fmt.Errorf("unknown revision self-improvement outcome has no missing stage")
	}
	for name, digest := range map[string]string{
		"iteration": o.IterationDigest,
		"application observation": o.ApplicationObservationDigest,
		"metrics binding": o.MetricsBindingDigest,
		"generation assessment": o.GenerationAssessmentDigest,
		"source": o.SourceDigest,
		"proposed source": o.ProposedSourceDigest,
		"input IR": o.InputIRDigest,
		"proposed IR": o.ProposedIRDigest,
		"candidate": o.CandidateDigest,
		"edit": o.EditDigest,
		"application": o.ApplicationDigest,
		"metrics": o.MetricsDigest,
		"generated source": o.GeneratedSourceDigest,
		"generated IR": o.GeneratedIRDigest,
		"structure": o.StructureDigest,
	} {
		if !validDigest(digest) {
			return fmt.Errorf("revision self-improvement outcome %s digest is invalid", name)
		}
	}
	if o.ChangedByteCount < 0 || o.ChangedLineCount < 0 {
		return fmt.Errorf("revision self-improvement outcome counts must be non-negative")
	}
	if o.OutcomeSignal != "metrics-generation-bound" {
		return fmt.Errorf("revision self-improvement outcome signal is invalid")
	}
	if !o.NonExecuting || !o.NonAuthorizing {
		return fmt.Errorf("revision self-improvement outcome must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementOutcomeObservation(o) != o.OutcomeDigest {
		return fmt.Errorf("revision self-improvement outcome digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementOutcomeObservation(observation RevisionSelfImprovementOutcomeObservation) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%d|%d|%t|%t|%t|%s|%t|%t",
		observation.Status,
		observation.MissingStage,
		observation.IterationDigest,
		observation.ApplicationObservationDigest,
		observation.MetricsBindingDigest,
		observation.GenerationAssessmentDigest,
		observation.SourceDigest,
		observation.ProposedSourceDigest,
		observation.InputIRDigest,
		observation.ProposedIRDigest,
		observation.CandidateDigest,
		observation.EditDigest,
		observation.ApplicationDigest,
		observation.MetricsDigest,
		observation.GeneratedSourceDigest,
		observation.GeneratedIRDigest,
		observation.StructureDigest,
		observation.ChangedByteCount,
		observation.ChangedLineCount,
		observation.IRChanged,
		observation.ExactSourceMatch,
		observation.StructureMatch,
		observation.OutcomeSignal,
		observation.NonExecuting,
		observation.NonAuthorizing,
	))
}
