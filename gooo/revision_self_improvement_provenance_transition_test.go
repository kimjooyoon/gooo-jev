package gooo

import "testing"

func selfImprovementProvenanceTransitionInputs(t *testing.T) (RevisionSelfImprovementLifecycleObservation, RevisionSelfImprovementLifecycleObservation) {
	t.Helper()
	executionApplication, outcome := selfImprovementLifecycleInputs(t)
	lifecycle, err := ObserveRevisionSelfImprovementLifecycle(executionApplication, outcome)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementLifecycle() error = %v", err)
	}
	return lifecycle, lifecycle
}

func TestObserveRevisionSelfImprovementProvenanceTransitionRetainsStableEvidence(t *testing.T) {
	previous, current := selfImprovementProvenanceTransitionInputs(t)
	result, err := ObserveRevisionSelfImprovementProvenanceTransition(previous, current)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceTransition() error = %v", err)
	}
	if result.Status != "BOUND" || result.TransitionSignal != "provenance-stable" ||
		result.MetricsChanged || result.GenerationChanged || result.SourceChanged {
		t.Fatalf("unexpected stable provenance transition: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementProvenanceTransitionRecordsMetricDrift(t *testing.T) {
	previous, current := selfImprovementProvenanceTransitionInputs(t)
	current.ChangedByteCount++
	current.ObservationDigest = digestRevisionSelfImprovementLifecycle(current)
	result, err := ObserveRevisionSelfImprovementProvenanceTransition(previous, current)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceTransition() error = %v", err)
	}
	if result.Status != "BOUND" || result.TransitionSignal != "provenance-transition-bound" ||
		!result.MetricsChanged || result.GenerationChanged {
		t.Fatalf("unexpected metric provenance transition: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementProvenanceTransitionRetainsPreviousFailure(t *testing.T) {
	previous, current := selfImprovementProvenanceTransitionInputs(t)
	previous.ObservationDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementProvenanceTransition(previous, current)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementProvenanceTransition() error = nil, want previous failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-provenance-transition-previous" {
		t.Fatalf("unexpected unknown previous transition: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementProvenanceTransitionRetainsCurrentFailure(t *testing.T) {
	previous, current := selfImprovementProvenanceTransitionInputs(t)
	current.ObservationDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementProvenanceTransition(previous, current)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementProvenanceTransition() error = nil, want current failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-provenance-transition-current" {
		t.Fatalf("unexpected unknown current transition: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementProvenanceTransitionIsDeterministic(t *testing.T) {
	previous, current := selfImprovementProvenanceTransitionInputs(t)
	first, err := ObserveRevisionSelfImprovementProvenanceTransition(previous, current)
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementProvenanceTransition() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementProvenanceTransition(previous, current)
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementProvenanceTransition() error = %v", err)
	}
	if first.ObservationDigest != second.ObservationDigest {
		t.Fatal("same lifecycle snapshots produced different transition digest")
	}
}
