package gooo

import "fmt"

// RevisionApplicationAssessment classifies an applied revision from bound evidence.
// It never executes, authorizes, or persists a proposed change.
type RevisionApplicationAssessment struct {
	Status               string
	MissingStage         string
	SourceDigest         string
	ReceiptDigest        string
	ApplicationDigest    string
	MetricsDigest        string
	BindingDigest        string
	ProposedSourceDigest string
	InputIRDigest        string
	ProposedIRDigest     string
	CandidateDigest      string
	EditDigest           string
	ChangedByteCount     int
	ChangedLineCount     int
	IRChanged            bool
	ChangeClass          string
	ReviewSignal         string
	AssessmentDigest     string
	NonExecuting         bool
	NonAuthorizing       bool
}

// AssessRevisionApplication binds the receipt and metrics binding before
// deriving a bounded, evidence-backed change class and review signal.
func AssessRevisionApplication(receipt RevisionApplicationReceipt, binding RevisionMetricsBinding) (RevisionApplicationAssessment, error) {
	assessment := RevisionApplicationAssessment{
		Status:               "UNKNOWN",
		MissingStage:         "revision-application-assessment",
		SourceDigest:         receipt.SourceDigest,
		ReceiptDigest:        receipt.ReceiptDigest,
		ApplicationDigest:    receipt.ApplicationDigest,
		MetricsDigest:        binding.MetricsDigest,
		BindingDigest:        binding.BindingDigest,
		ProposedSourceDigest: receipt.ProposedSourceDigest,
		InputIRDigest:        receipt.InputIRDigest,
		ProposedIRDigest:     receipt.ProposedIRDigest,
		CandidateDigest:      receipt.CandidateDigest,
		EditDigest:           receipt.EditDigest,
		ChangedByteCount:     binding.ChangedByteCount,
		ChangedLineCount:     binding.ChangedLineCount,
		IRChanged:            binding.IRChanged,
		NonExecuting:         true,
		NonAuthorizing:       true,
	}
	setAssessmentDigest := func() {
		assessment.AssessmentDigest = digestRevisionApplicationAssessment(assessment)
	}
	setAssessmentDigest()

	if err := receipt.Validate(); err != nil {
		assessment.MissingStage = "revision-application-assessment-receipt"
		setAssessmentDigest()
		return assessment, fmt.Errorf("revision application receipt is not valid: %w", err)
	}
	if err := binding.Validate(); err != nil {
		assessment.MissingStage = "revision-application-assessment-binding"
		setAssessmentDigest()
		return assessment, fmt.Errorf("revision metrics binding is not valid: %w", err)
	}
	if receipt.ReceiptDigest != binding.ReceiptDigest ||
		receipt.ApplicationDigest != binding.ApplicationDigest ||
		receipt.SourceDigest != binding.SourceDigest ||
		receipt.ProposedSourceDigest != binding.ProposedSourceDigest ||
		receipt.InputIRDigest != binding.InputIRDigest ||
		receipt.ProposedIRDigest != binding.ProposedIRDigest ||
		receipt.CandidateDigest != binding.CandidateDigest ||
		receipt.EditDigest != binding.EditDigest {
		assessment.MissingStage = "revision-application-assessment-link"
		setAssessmentDigest()
		return assessment, fmt.Errorf("revision application receipt and metrics binding are not linked")
	}

	assessment.Status = "BOUND"
	assessment.MissingStage = ""
	assessment.ChangeClass = revisionChangeClass(binding)
	assessment.ReviewSignal = revisionReviewSignal(binding)
	setAssessmentDigest()
	if err := assessment.Validate(); err != nil {
		assessment.Status = "UNKNOWN"
		assessment.MissingStage = "revision-application-assessment"
		setAssessmentDigest()
		return assessment, fmt.Errorf("revision application assessment is not valid: %w", err)
	}
	return assessment, nil
}

func (a RevisionApplicationAssessment) Validate() error {
	if a.Status == "" {
		return fmt.Errorf("revision application assessment status is empty")
	}
	if a.Status == "BOUND" && a.MissingStage != "" {
		return fmt.Errorf("bound revision application assessment has a missing stage")
	}
	if a.Status == "UNKNOWN" && a.MissingStage == "" {
		return fmt.Errorf("unknown revision application assessment has no missing stage")
	}
	if a.ChangedByteCount < 0 || a.ChangedLineCount < 0 {
		return fmt.Errorf("revision application assessment change counts must be non-negative")
	}
	if !a.NonExecuting || !a.NonAuthorizing {
		return fmt.Errorf("revision application assessment must remain non-executing and non-authorizing")
	}
	if digestRevisionApplicationAssessment(a) != a.AssessmentDigest {
		return fmt.Errorf("revision application assessment digest does not match its fields")
	}
	return nil
}

func revisionChangeClass(binding RevisionMetricsBinding) string {
	if binding.ChangedByteCount == 0 && binding.ChangedLineCount == 0 && !binding.IRChanged {
		return "none"
	}
	if binding.IRChanged {
		return "structural"
	}
	if binding.ChangedLineCount > 0 {
		return "localized"
	}
	return "textual"
}

func revisionReviewSignal(binding RevisionMetricsBinding) string {
	if binding.IRChanged {
		return "inspect"
	}
	if binding.ChangedByteCount > 0 || binding.ChangedLineCount > 0 {
		return "review"
	}
	return "observe"
}

func digestRevisionApplicationAssessment(assessment RevisionApplicationAssessment) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%d|%d|%t|%t|%t",
		assessment.Status,
		assessment.MissingStage,
		assessment.SourceDigest,
		assessment.ReceiptDigest,
		assessment.ApplicationDigest,
		assessment.MetricsDigest,
		assessment.BindingDigest,
		assessment.ProposedSourceDigest,
		assessment.InputIRDigest,
		assessment.ProposedIRDigest,
		assessment.CandidateDigest,
		assessment.EditDigest,
		assessment.ChangeClass,
		assessment.ReviewSignal,
		assessment.ChangedByteCount,
		assessment.ChangedLineCount,
		assessment.IRChanged,
		assessment.NonExecuting,
		assessment.NonAuthorizing,
	))
}