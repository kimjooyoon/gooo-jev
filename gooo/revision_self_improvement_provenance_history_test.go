package gooo

import "testing"

func provenanceHistoryTransition(t *testing.T, metricDrift bool) RevisionSelfImprovementProvenanceTransitionObservation {
	t.Helper()
	previous, current := selfImprovementProvenanceTransitionInputs(t)
	if metricDrift {
		current.ChangedByteCount++
		current.ObservationDigest = digestRevisionSelfImprovementLifecycle(current)
	}
	transition, err := ObserveRevisionSelfImprovementProvenanceTransition(previous, current)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceTransition() error = %v", err)
	}
	return transition
}

func TestObserveRevisionSelfImprovementProvenanceHistoryBindsStableTransition(t *testing.T) {
	history, err := ObserveRevisionSelfImprovementProvenanceHistory([]RevisionSelfImprovementProvenanceTransitionObservation{
		provenanceHistoryTransition(t, false),
	})
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceHistory() error = %v", err)
	}
	if history.Status != "BOUND" || history.HistorySignal != "stable" ||
		history.StableCount != 1 || history.TransitionedCount != 0 {
		t.Fatalf("unexpected provenance history: %#v", history)
	}
	if err := history.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementProvenanceHistoryRecordsMetricTransition(t *testing.T) {
	history, err := ObserveRevisionSelfImprovementProvenanceHistory([]RevisionSelfImprovementProvenanceTransitionObservation{
		provenanceHistoryTransition(t, false),
		provenanceHistoryTransition(t, true),
	})
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceHistory() error = %v", err)
	}
	if history.Status != "BOUND" || history.HistorySignal != "mixed" ||
		history.StableCount != 1 || history.TransitionedCount != 1 ||
		history.MetricsChangeCount != 1 {
		t.Fatalf("unexpected metric provenance history: %#v", history)
	}
}

func TestObserveRevisionSelfImprovementProvenanceHistoryRetainsTransitionFailure(t *testing.T) {
	transition := provenanceHistoryTransition(t, false)
	transition.ObservationDigest = digestString("tampered")
	history, err := ObserveRevisionSelfImprovementProvenanceHistory([]RevisionSelfImprovementProvenanceTransitionObservation{
		transition,
	})
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementProvenanceHistory() error = nil, want transition failure")
	}
	if history.Status != "UNKNOWN" ||
		history.MissingStage != "revision-self-improvement-provenance-history-transition-0" {
		t.Fatalf("unexpected unknown provenance history: %#v", history)
	}
}

func TestObserveRevisionSelfImprovementProvenanceHistoryIsDeterministic(t *testing.T) {
	first, err := ObserveRevisionSelfImprovementProvenanceHistory([]RevisionSelfImprovementProvenanceTransitionObservation{
		provenanceHistoryTransition(t, false),
		provenanceHistoryTransition(t, true),
	})
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementProvenanceHistory() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementProvenanceHistory([]RevisionSelfImprovementProvenanceTransitionObservation{
		provenanceHistoryTransition(t, false),
		provenanceHistoryTransition(t, true),
	})
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementProvenanceHistory() error = %v", err)
	}
	if first.HistoryDigest != second.HistoryDigest {
		t.Fatal("same provenance transitions produced different history digest")
	}
}
