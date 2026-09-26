package gooo

import "testing"

func selfImprovementOutcomeInputs(t *testing.T) (RevisionSelfImprovementIteration, RevisionSelfImprovementApplicationObservation, RevisionMetricsBinding, RevisionGenerationAssessment) {
	t.Helper()
	iteration, planObservation, applicationObservation := selfImprovementApplicationInputs(t)
	application, err := ObserveRevisionSelfImprovementApplication(iteration, planObservation, applicationObservation)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementApplication() error = %v", err)
	}
	_, metricsBinding, _, generationAssessment := selfImprovementReceiptInputs(t)
	return iteration, application, metricsBinding, generationAssessment
}

func TestObserveRevisionSelfImprovementOutcomeBindsLifecycle(t *testing.T) {
	iteration, application, metrics, generation := selfImprovementOutcomeInputs(t)
	result, err := ObserveRevisionSelfImprovementOutcome(iteration, application, metrics, generation)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementOutcome() error = %v", err)
	}
	if result.Status != "BOUND" || result.OutcomeSignal != "metrics-generation-bound" ||
		result.ApplicationDigest != application.ApplicationDigest ||
		result.MetricsDigest != metrics.MetricsDigest ||
		result.GeneratedIRDigest != generation.GeneratedIRDigest {
		t.Fatalf("unexpected self-improvement outcome: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementOutcomeRetainsApplicationFailure(t *testing.T) {
	iteration, application, metrics, generation := selfImprovementOutcomeInputs(t)
	application.ApplicationObservationDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementOutcome(iteration, application, metrics, generation)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementOutcome() error = nil, want application failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-outcome-application" {
		t.Fatalf("unexpected unknown application outcome: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementOutcomeRetainsMetricsFailure(t *testing.T) {
	iteration, application, metrics, generation := selfImprovementOutcomeInputs(t)
	metrics.BindingDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementOutcome(iteration, application, metrics, generation)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementOutcome() error = nil, want metrics failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-outcome-metrics" {
		t.Fatalf("unexpected unknown metrics outcome: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementOutcomeRetainsGenerationFailure(t *testing.T) {
	iteration, application, metrics, generation := selfImprovementOutcomeInputs(t)
	generation.GenerationDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementOutcome(iteration, application, metrics, generation)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementOutcome() error = nil, want generation failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-outcome-generation" {
		t.Fatalf("unexpected unknown generation outcome: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementOutcomeRetainsApplicationMetricsLinkFailure(t *testing.T) {
	iteration, application, metrics, generation := selfImprovementOutcomeInputs(t)
	metrics.ApplicationDigest = digestString("other-application")
	metrics.BindingDigest = digestRevisionMetricsBinding(metrics)
	result, err := ObserveRevisionSelfImprovementOutcome(iteration, application, metrics, generation)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementOutcome() error = nil, want application-metrics link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-outcome-application-metrics-link" {
		t.Fatalf("unexpected unknown application-metrics outcome: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementOutcomeRetainsGenerationLinkFailure(t *testing.T) {
	iteration, application, metrics, generation := selfImprovementOutcomeInputs(t)
	generation.GenerationSourceDigest = digestString("other-source")
	generation.GenerationDigest = digestRevisionGenerationAssessment(generation)
	result, err := ObserveRevisionSelfImprovementOutcome(iteration, application, metrics, generation)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementOutcome() error = nil, want generation link failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-outcome-generation-link" {
		t.Fatalf("unexpected unknown generation-link outcome: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementOutcomeIsDeterministic(t *testing.T) {
	iteration, application, metrics, generation := selfImprovementOutcomeInputs(t)
	first, err := ObserveRevisionSelfImprovementOutcome(iteration, application, metrics, generation)
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementOutcome() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementOutcome(iteration, application, metrics, generation)
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementOutcome() error = %v", err)
	}
	if first.OutcomeDigest != second.OutcomeDigest {
		t.Fatal("same lifecycle evidence produced different outcome digest")
	}
}
