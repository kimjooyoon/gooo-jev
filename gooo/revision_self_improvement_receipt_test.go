package gooo

import "testing"

func selfImprovementReceiptInputs(t *testing.T) (RevisionApplicationObservation, RevisionMetricsBinding, RevisionApplicationAssessment, RevisionGenerationAssessment) {
	t.Helper()
	planObservation, receipt := revisionApplicationObservationInputs(t)
	application, err := ApplyRevision(validContract, receipt.SourceDigest, mustCandidateFromReceipt(t, receipt), mustEditFromReceipt(t, receipt))
	if err != nil {
		t.Fatalf("ApplyRevision() error = %v", err)
	}
	metrics, err := MeasureRevision(validContract, application)
	if err != nil {
		t.Fatalf("MeasureRevision() error = %v", err)
	}
	metricsBinding, err := ObserveRevisionMetricsBinding(receipt, metrics)
	if err != nil {
		t.Fatalf("ObserveRevisionMetricsBinding() error = %v", err)
	}
	applicationAssessment, err := AssessRevisionApplication(receipt, metricsBinding)
	if err != nil {
		t.Fatalf("AssessRevisionApplication() error = %v", err)
	}
	proposed, err := Parse(application.ProposedSource)
	if err != nil {
		t.Fatalf("Parse(proposed source) error = %v", err)
	}
	generation, err := Generate(proposed)
	if err != nil {
		t.Fatalf("Generate(proposed IR) error = %v", err)
	}
	generationAssessment, err := AssessRevisionGeneration(applicationAssessment, generation)
	if err != nil {
		t.Fatalf("AssessRevisionGeneration() error = %v", err)
	}
	applicationObservation, err := ObserveRevisionApplication(planObservation, receipt)
	if err != nil {
		t.Fatalf("ObserveRevisionApplication() error = %v", err)
	}
	return applicationObservation, metricsBinding, applicationAssessment, generationAssessment
}

func mustCandidateFromReceipt(t *testing.T, receipt RevisionApplicationReceipt) RevisionCandidate {
	t.Helper()
	_, plan := candidatePlanObservationInputs(t)
	if plan.CandidateDigest != receipt.CandidateDigest {
		t.Fatalf("fixture receipt candidate mismatch: plan=%s receipt=%s", plan.CandidateDigest, receipt.CandidateDigest)
	}
	return plan.Candidate
}

func mustEditFromReceipt(t *testing.T, receipt RevisionApplicationReceipt) SourceEdit {
	t.Helper()
	_, plan := candidatePlanObservationInputs(t)
	if plan.EditDigest != receipt.EditDigest {
		t.Fatalf("fixture receipt edit mismatch: plan=%s receipt=%s", plan.EditDigest, receipt.EditDigest)
	}
	return plan.Edit
}

func TestObserveRevisionSelfImprovementReceiptBindsLifecycle(t *testing.T) {
	applicationObservation, metricsBinding, applicationAssessment, generationAssessment := selfImprovementReceiptInputs(t)
	receipt, err := ObserveRevisionSelfImprovementReceipt(applicationObservation, metricsBinding, applicationAssessment, generationAssessment)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementReceipt() error = %v", err)
	}
	if receipt.Status != "BOUND" || receipt.StageCount != 4 || !receipt.StructureMatch {
		t.Fatalf("unexpected self-improvement receipt: %#v", receipt)
	}
	if receipt.ApplicationDigest != applicationObservation.ApplicationDigest ||
		receipt.MetricsDigest != metricsBinding.MetricsDigest ||
		receipt.GenerationAssessmentDigest != generationAssessment.GenerationDigest {
		t.Fatalf("receipt lost lifecycle links: %#v", receipt)
	}
	if err := receipt.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementReceiptRetainsApplicationFailure(t *testing.T) {
	applicationObservation, metricsBinding, applicationAssessment, generationAssessment := selfImprovementReceiptInputs(t)
	applicationObservation.ObservationDigest = digestString("tampered")
	receipt, err := ObserveRevisionSelfImprovementReceipt(applicationObservation, metricsBinding, applicationAssessment, generationAssessment)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementReceipt() error = nil, want application failure")
	}
	if receipt.Status != "UNKNOWN" || receipt.MissingStage != "revision-self-improvement-receipt-application" {
		t.Fatalf("unexpected unknown application receipt: %#v", receipt)
	}
}

func TestObserveRevisionSelfImprovementReceiptRetainsGenerationFailure(t *testing.T) {
	applicationObservation, metricsBinding, applicationAssessment, generationAssessment := selfImprovementReceiptInputs(t)
	generationAssessment.GenerationDigest = digestString("tampered")
	receipt, err := ObserveRevisionSelfImprovementReceipt(applicationObservation, metricsBinding, applicationAssessment, generationAssessment)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementReceipt() error = nil, want generation failure")
	}
	if receipt.Status != "UNKNOWN" || receipt.MissingStage != "revision-self-improvement-receipt-generation" {
		t.Fatalf("unexpected unknown generation receipt: %#v", receipt)
	}
}

func TestObserveRevisionSelfImprovementReceiptRetainsLinkFailure(t *testing.T) {
	applicationObservation, metricsBinding, applicationAssessment, generationAssessment := selfImprovementReceiptInputs(t)
	generationAssessment.ApplicationDigest = digestString("other-application")
	generationAssessment.GenerationDigest = digestRevisionGenerationAssessment(generationAssessment)
	receipt, err := ObserveRevisionSelfImprovementReceipt(applicationObservation, metricsBinding, applicationAssessment, generationAssessment)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementReceipt() error = nil, want link failure")
	}
	if receipt.Status != "UNKNOWN" || receipt.MissingStage != "revision-self-improvement-receipt-link" {
		t.Fatalf("unexpected unknown link receipt: %#v", receipt)
	}
}

func TestObserveRevisionSelfImprovementReceiptIsDeterministic(t *testing.T) {
	applicationObservation, metricsBinding, applicationAssessment, generationAssessment := selfImprovementReceiptInputs(t)
	first, err := ObserveRevisionSelfImprovementReceipt(applicationObservation, metricsBinding, applicationAssessment, generationAssessment)
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementReceipt() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementReceipt(applicationObservation, metricsBinding, applicationAssessment, generationAssessment)
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementReceipt() error = %v", err)
	}
	if first.ReceiptDigest != second.ReceiptDigest {
		t.Fatal("same lifecycle evidence produced different receipt digest")
	}
}