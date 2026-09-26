package gooo

import "testing"

func generationAssessmentInputs(t *testing.T) (RevisionApplicationAssessment, GenerationReceipt) {
	t.Helper()
	plan, application := appliedRevisionReceiptInputs(t)
	receipt, err := ObserveRevisionApplicationReceipt(plan, application)
	if err != nil {
		t.Fatalf("ObserveRevisionApplicationReceipt() error = %v", err)
	}
	metrics, err := MeasureRevision(validContract, application)
	if err != nil {
		t.Fatalf("MeasureRevision() error = %v", err)
	}
	binding, err := ObserveRevisionMetricsBinding(receipt, metrics)
	if err != nil {
		t.Fatalf("ObserveRevisionMetricsBinding() error = %v", err)
	}
	assessment, err := AssessRevisionApplication(receipt, binding)
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
	return assessment, generation
}

func TestAssessRevisionGenerationBindsReverseObservation(t *testing.T) {
	assessment, generation := generationAssessmentInputs(t)
	result, err := AssessRevisionGeneration(assessment, generation)
	if err != nil {
		t.Fatalf("AssessRevisionGeneration() error = %v", err)
	}
	if result.Status != "BOUND" || !result.StructureMatch || result.SourceMatchClass == "" {
		t.Fatalf("unexpected revision generation assessment: %#v", result)
	}
	if result.ProposedSourceDigest != generation.SourceDigest ||
		result.GeneratedIRDigest != generation.GeneratedIRDigest {
		t.Fatalf("assessment lost generation links: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestAssessRevisionGenerationRetainsReverseObservationFailure(t *testing.T) {
	assessment, generation := generationAssessmentInputs(t)
	generation.StructureMatch = false
	result, err := AssessRevisionGeneration(assessment, generation)
	if err == nil {
		t.Fatal("AssessRevisionGeneration() error = nil, want reverse observation failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-generation-assessment-generation" {
		t.Fatalf("unexpected unknown generation assessment: %#v", result)
	}
}

func TestAssessRevisionGenerationRetainsSourceLinkFailure(t *testing.T) {
	assessment, generation := generationAssessmentInputs(t)
	generation.SourceDigest = digestString("other-source")
	result, err := AssessRevisionGeneration(assessment, generation)
	if err == nil {
		t.Fatal("AssessRevisionGeneration() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-generation-assessment-link" {
		t.Fatalf("unexpected unknown link assessment: %#v", result)
	}
}

func TestAssessRevisionGenerationIsDeterministic(t *testing.T) {
	assessment, generation := generationAssessmentInputs(t)
	first, err := AssessRevisionGeneration(assessment, generation)
	if err != nil {
		t.Fatalf("first AssessRevisionGeneration() error = %v", err)
	}
	second, err := AssessRevisionGeneration(assessment, generation)
	if err != nil {
		t.Fatalf("second AssessRevisionGeneration() error = %v", err)
	}
	if first.GenerationDigest != second.GenerationDigest {
		t.Fatal("same generation evidence produced different assessment digest")
	}
}