package gooo

import "testing"

func selfImprovementReverseObservationInputs(t *testing.T) (RevisionSelfImprovementLifecycleObservation, GenerationReceipt) {
	t.Helper()
	executionApplication, outcome := selfImprovementLifecycleInputs(t)
	lifecycle, err := ObserveRevisionSelfImprovementLifecycle(executionApplication, outcome)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementLifecycle() error = %v", err)
	}
	_, generation := generationAssessmentInputs(t)
	return lifecycle, generation
}

func TestObserveRevisionSelfImprovementReverseGenerationBindsEvidence(t *testing.T) {
	lifecycle, generation := selfImprovementReverseObservationInputs(t)
	result, err := ObserveRevisionSelfImprovementReverseGeneration(lifecycle, generation)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementReverseGeneration() error = %v", err)
	}
	if result.Status != "BOUND" || result.ReverseSignal != "reverse-observed" ||
		!result.ExactSourceMatch || !result.ExactIRMatch || !result.ExactStructureMatch ||
		result.MetricsBindingDigest != lifecycle.MetricsBindingDigest {
		t.Fatalf("unexpected self-improvement reverse observation: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementReverseGenerationRetainsSourceLinkFailure(t *testing.T) {
	lifecycle, generation := selfImprovementReverseObservationInputs(t)
	generation.SourceDigest = digestString("other-source")
	result, err := ObserveRevisionSelfImprovementReverseGeneration(lifecycle, generation)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementReverseGeneration() error = nil, want source link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-reverse-source-link" {
		t.Fatalf("unexpected unknown source-link observation: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementReverseGenerationRetainsSourceIntegrityFailure(t *testing.T) {
	lifecycle, generation := selfImprovementReverseObservationInputs(t)
	generation.GeneratedSource += "\n"
	result, err := ObserveRevisionSelfImprovementReverseGeneration(lifecycle, generation)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementReverseGeneration() error = nil, want generated source integrity failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-reverse-generated-source-integrity" {
		t.Fatalf("unexpected unknown source-integrity observation: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementReverseGenerationRetainsStructureFailure(t *testing.T) {
	lifecycle, generation := selfImprovementReverseObservationInputs(t)
	generation.StructureDigest = digestString("other-structure")
	result, err := ObserveRevisionSelfImprovementReverseGeneration(lifecycle, generation)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementReverseGeneration() error = nil, want structure failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-reverse-structure" {
		t.Fatalf("unexpected unknown structure observation: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementReverseGenerationIsDeterministic(t *testing.T) {
	lifecycle, generation := selfImprovementReverseObservationInputs(t)
	first, err := ObserveRevisionSelfImprovementReverseGeneration(lifecycle, generation)
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementReverseGeneration() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementReverseGeneration(lifecycle, generation)
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementReverseGeneration() error = %v", err)
	}
	if first.ObservationDigest != second.ObservationDigest {
		t.Fatal("same lifecycle and generation evidence produced different reverse observation digest")
	}
}
