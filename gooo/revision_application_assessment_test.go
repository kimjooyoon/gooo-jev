package gooo

import "testing"

func assessmentInputs(t *testing.T) (RevisionApplicationReceipt, RevisionMetricsBinding) {
	t.Helper()
	receipt, metrics := metricsBindingInputs(t)
	binding, err := ObserveRevisionMetricsBinding(receipt, metrics)
	if err != nil {
		t.Fatalf("ObserveRevisionMetricsBinding() error = %v", err)
	}
	return receipt, binding
}

func TestAssessRevisionApplicationBindsEvidenceAndClassifiesChange(t *testing.T) {
	receipt, binding := assessmentInputs(t)
	assessment, err := AssessRevisionApplication(receipt, binding)
	if err != nil {
		t.Fatalf("AssessRevisionApplication() error = %v", err)
	}
	if assessment.Status != "BOUND" || assessment.ChangeClass == "" || assessment.ReviewSignal == "" {
		t.Fatalf("unexpected revision application assessment: %#v", assessment)
	}
	if assessment.SourceDigest != receipt.SourceDigest ||
		assessment.BindingDigest != binding.BindingDigest ||
		assessment.ChangedByteCount != binding.ChangedByteCount {
		t.Fatalf("assessment lost evidence links: %#v", assessment)
	}
	if err := assessment.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestAssessRevisionApplicationRetainsReceiptFailure(t *testing.T) {
	receipt, binding := assessmentInputs(t)
	receipt.ReceiptDigest = digestString("tampered")
	assessment, err := AssessRevisionApplication(receipt, binding)
	if err == nil {
		t.Fatal("AssessRevisionApplication() error = nil, want receipt failure")
	}
	if assessment.Status != "UNKNOWN" || assessment.MissingStage != "revision-application-assessment-receipt" {
		t.Fatalf("unexpected unknown assessment: %#v", assessment)
	}
}

func TestAssessRevisionApplicationRetainsLinkFailure(t *testing.T) {
	receipt, binding := assessmentInputs(t)
	binding.ReceiptDigest = digestString("other-receipt")
	binding.BindingDigest = digestRevisionMetricsBinding(binding)
	assessment, err := AssessRevisionApplication(receipt, binding)
	if err == nil {
		t.Fatal("AssessRevisionApplication() error = nil, want link failure")
	}
	if assessment.Status != "UNKNOWN" || assessment.MissingStage != "revision-application-assessment-link" {
		t.Fatalf("unexpected unknown link assessment: %#v", assessment)
	}
}

func TestAssessRevisionApplicationIsDeterministic(t *testing.T) {
	receipt, binding := assessmentInputs(t)
	first, err := AssessRevisionApplication(receipt, binding)
	if err != nil {
		t.Fatalf("first AssessRevisionApplication() error = %v", err)
	}
	second, err := AssessRevisionApplication(receipt, binding)
	if err != nil {
		t.Fatalf("second AssessRevisionApplication() error = %v", err)
	}
	if first.AssessmentDigest != second.AssessmentDigest {
		t.Fatal("same evidence produced different assessment digest")
	}
}