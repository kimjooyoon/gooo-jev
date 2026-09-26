package gooo

import "testing"

func TestObserveRevisionSelfImprovementProvenanceDecisionBindsStableHistory(t *testing.T) {
	history, err := ObserveRevisionSelfImprovementProvenanceHistory([]RevisionSelfImprovementProvenanceTransitionObservation{
		provenanceHistoryTransition(t, false),
	})
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceHistory() error = %v", err)
	}
	result, err := ObserveRevisionSelfImprovementProvenanceDecision(history)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceDecision() error = %v", err)
	}
	if result.Status != "BOUND" || result.DecisionSignal != "observe" ||
		result.DecisionReason != "provenance-stable-observation" ||
		result.RequiresInspection || result.RequiresMeasurement {
		t.Fatalf("unexpected stable provenance decision: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementProvenanceDecisionRequestsInspection(t *testing.T) {
	history, err := ObserveRevisionSelfImprovementProvenanceHistory([]RevisionSelfImprovementProvenanceTransitionObservation{
		provenanceHistoryTransition(t, false),
		provenanceHistoryTransition(t, true),
	})
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceHistory() error = %v", err)
	}
	result, err := ObserveRevisionSelfImprovementProvenanceDecision(history)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceDecision() error = %v", err)
	}
	if result.Status != "BOUND" || result.DecisionSignal != "inspect" ||
		!result.RequiresInspection || !result.RequiresMeasurement {
		t.Fatalf("unexpected transitioned provenance decision: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementProvenanceDecisionRetainsHistoryFailure(t *testing.T) {
	history, err := ObserveRevisionSelfImprovementProvenanceHistory([]RevisionSelfImprovementProvenanceTransitionObservation{
		provenanceHistoryTransition(t, false),
	})
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceHistory() error = %v", err)
	}
	history.HistoryDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementProvenanceDecision(history)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementProvenanceDecision() error = nil, want history failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-provenance-decision-history" {
		t.Fatalf("unexpected unknown provenance decision: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementProvenanceDecisionIsDeterministic(t *testing.T) {
	history, err := ObserveRevisionSelfImprovementProvenanceHistory([]RevisionSelfImprovementProvenanceTransitionObservation{
		provenanceHistoryTransition(t, false),
		provenanceHistoryTransition(t, true),
	})
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceHistory() error = %v", err)
	}
	first, err := ObserveRevisionSelfImprovementProvenanceDecision(history)
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementProvenanceDecision() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementProvenanceDecision(history)
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementProvenanceDecision() error = %v", err)
	}
	if first.DecisionDigest != second.DecisionDigest {
		t.Fatal("same provenance history produced different decision digest")
	}
}
