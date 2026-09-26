package gooo

import "fmt"

// RevisionSelfImprovementReceipt is the composite, read-only provenance
// record for one observed candidate through application and generation.
type RevisionSelfImprovementReceipt struct {
	Status                         string
	MissingStage                   string
	ApplicationObservationDigest   string
	MetricsBindingDigest           string
	ApplicationAssessmentDigest    string
	GenerationAssessmentDigest     string
	SourceDigest                   string
	ProposedSourceDigest           string
	InputIRDigest                  string
	ProposedIRDigest               string
	ApplicationDigest              string
	CandidateDigest                string
	EditDigest                     string
	MetricsDigest                  string
	GeneratedSourceDigest          string
	GeneratedIRDigest              string
	StructureDigest                string
	ChangedByteCount               int
	ChangedLineCount               int
	IRChanged                      bool
	ExactSourceMatch               bool
	StructureMatch                 bool
	StageCount                     int
	ReceiptDigest                  string
	NonExecuting                   bool
	NonAuthorizing                 bool
}

// ObserveRevisionSelfImprovementReceipt links the lifecycle observations
// without executing, persisting, or authorizing the proposed change.
func ObserveRevisionSelfImprovementReceipt(applicationObservation RevisionApplicationObservation, metricsBinding RevisionMetricsBinding, applicationAssessment RevisionApplicationAssessment, generationAssessment RevisionGenerationAssessment) (RevisionSelfImprovementReceipt, error) {
	receipt := RevisionSelfImprovementReceipt{
		Status:                       "UNKNOWN",
		MissingStage:                 "revision-self-improvement-receipt",
		ApplicationObservationDigest: applicationObservation.ObservationDigest,
		MetricsBindingDigest:         metricsBinding.BindingDigest,
		ApplicationAssessmentDigest:  applicationAssessment.AssessmentDigest,
		GenerationAssessmentDigest:   generationAssessment.GenerationDigest,
		SourceDigest:                 applicationObservation.SourceDigest,
		ProposedSourceDigest:         applicationObservation.ProposedSourceDigest,
		InputIRDigest:                applicationObservation.InputIRDigest,
		ProposedIRDigest:             applicationObservation.ProposedIRDigest,
		ApplicationDigest:            applicationObservation.ApplicationDigest,
		CandidateDigest:              applicationObservation.CandidateDigest,
		EditDigest:                   applicationObservation.EditDigest,
		MetricsDigest:                metricsBinding.MetricsDigest,
		GeneratedSourceDigest:        generationAssessment.GeneratedSourceDigest,
		GeneratedIRDigest:            generationAssessment.GeneratedIRDigest,
		StructureDigest:              generationAssessment.StructureDigest,
		ChangedByteCount:             metricsBinding.ChangedByteCount,
		ChangedLineCount:             metricsBinding.ChangedLineCount,
		IRChanged:                    metricsBinding.IRChanged,
		ExactSourceMatch:             generationAssessment.ExactSourceMatch,
		StructureMatch:               generationAssessment.StructureMatch,
		StageCount:                   4,
		NonExecuting:                 true,
		NonAuthorizing:               true,
	}
	setReceiptDigest := func() {
		receipt.ReceiptDigest = digestRevisionSelfImprovementReceipt(receipt)
	}
	setReceiptDigest()

	if err := applicationObservation.Validate(); err != nil {
		receipt.MissingStage = "revision-self-improvement-receipt-application"
		setReceiptDigest()
		return receipt, fmt.Errorf("revision application observation is not valid: %w", err)
	}
	if err := metricsBinding.Validate(); err != nil {
		receipt.MissingStage = "revision-self-improvement-receipt-metrics"
		setReceiptDigest()
		return receipt, fmt.Errorf("revision metrics binding is not valid: %w", err)
	}
	if err := applicationAssessment.Validate(); err != nil {
		receipt.MissingStage = "revision-self-improvement-receipt-assessment"
		setReceiptDigest()
		return receipt, fmt.Errorf("revision application assessment is not valid: %w", err)
	}
	if err := generationAssessment.Validate(); err != nil {
		receipt.MissingStage = "revision-self-improvement-receipt-generation"
		setReceiptDigest()
		return receipt, fmt.Errorf("revision generation assessment is not valid: %w", err)
	}
	if applicationObservation.ApplicationDigest != metricsBinding.ApplicationDigest ||
		applicationObservation.SourceDigest != metricsBinding.SourceDigest ||
		applicationObservation.ProposedSourceDigest != metricsBinding.ProposedSourceDigest ||
		applicationObservation.InputIRDigest != metricsBinding.InputIRDigest ||
		applicationObservation.ProposedIRDigest != metricsBinding.ProposedIRDigest ||
		applicationObservation.CandidateDigest != metricsBinding.CandidateDigest ||
		applicationObservation.EditDigest != metricsBinding.EditDigest ||
		applicationAssessment.ApplicationDigest != applicationObservation.ApplicationDigest ||
		applicationAssessment.MetricsDigest != metricsBinding.MetricsDigest ||
		applicationAssessment.SourceDigest != applicationObservation.SourceDigest ||
		applicationAssessment.ProposedSourceDigest != applicationObservation.ProposedSourceDigest ||
		applicationAssessment.InputIRDigest != applicationObservation.InputIRDigest ||
		applicationAssessment.ProposedIRDigest != applicationObservation.ProposedIRDigest ||
		applicationAssessment.CandidateDigest != applicationObservation.CandidateDigest ||
		applicationAssessment.EditDigest != applicationObservation.EditDigest ||
		generationAssessment.AssessmentDigest != applicationAssessment.AssessmentDigest ||
		generationAssessment.ApplicationDigest != applicationObservation.ApplicationDigest ||
		generationAssessment.SourceDigest != applicationObservation.SourceDigest ||
		generationAssessment.ProposedSourceDigest != applicationObservation.ProposedSourceDigest {
		receipt.MissingStage = "revision-self-improvement-receipt-link"
		setReceiptDigest()
		return receipt, fmt.Errorf("revision lifecycle observations are not linked")
	}

	receipt.Status = "BOUND"
	receipt.MissingStage = ""
	setReceiptDigest()
	if err := receipt.Validate(); err != nil {
		receipt.Status = "UNKNOWN"
		receipt.MissingStage = "revision-self-improvement-receipt"
		setReceiptDigest()
		return receipt, fmt.Errorf("revision self-improvement receipt is not valid: %w", err)
	}
	return receipt, nil
}

func (r RevisionSelfImprovementReceipt) Validate() error {
	if r.Status == "" {
		return fmt.Errorf("revision self-improvement receipt status is empty")
	}
	if r.Status == "BOUND" && r.MissingStage != "" {
		return fmt.Errorf("bound revision self-improvement receipt has a missing stage")
	}
	if r.Status == "UNKNOWN" && r.MissingStage == "" {
		return fmt.Errorf("unknown revision self-improvement receipt has no missing stage")
	}
	if r.StageCount != 4 {
		return fmt.Errorf("revision self-improvement receipt stage count must be 4")
	}
	if r.ChangedByteCount < 0 || r.ChangedLineCount < 0 {
		return fmt.Errorf("revision self-improvement receipt change counts must be non-negative")
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("revision self-improvement receipt must remain non-executing and non-authorizing")
	}
	if digestRevisionSelfImprovementReceipt(r) != r.ReceiptDigest {
		return fmt.Errorf("revision self-improvement receipt digest does not match its fields")
	}
	return nil
}

func digestRevisionSelfImprovementReceipt(receipt RevisionSelfImprovementReceipt) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%d|%d|%t|%t|%t|%d|%t|%t",
		receipt.Status,
		receipt.MissingStage,
		receipt.ApplicationObservationDigest,
		receipt.MetricsBindingDigest,
		receipt.ApplicationAssessmentDigest,
		receipt.GenerationAssessmentDigest,
		receipt.SourceDigest,
		receipt.ProposedSourceDigest,
		receipt.InputIRDigest,
		receipt.ProposedIRDigest,
		receipt.ApplicationDigest,
		receipt.CandidateDigest,
		receipt.EditDigest,
		receipt.MetricsDigest,
		receipt.GeneratedSourceDigest,
		receipt.GeneratedIRDigest,
		receipt.StructureDigest,
		receipt.ChangedByteCount,
		receipt.ChangedLineCount,
		receipt.IRChanged,
		receipt.ExactSourceMatch,
		receipt.StructureMatch,
		receipt.StageCount,
		receipt.NonExecuting,
		receipt.NonAuthorizing,
	))
}