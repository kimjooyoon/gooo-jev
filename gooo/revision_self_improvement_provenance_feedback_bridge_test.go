package gooo

import "testing"

func provenanceFeedbackBridgeInputs(t *testing.T, changed bool) (RevisionSelfImprovementProvenanceHistoryObservation, RevisionSelfImprovementProvenanceDecisionObservation, RevisionSelfImprovementFeedback) {
	t.Helper()
	transitions := []RevisionSelfImprovementProvenanceTransitionObservation{
		provenanceHistoryTransition(t, false),
	}
	var feedbackWindow RevisionSelfImprovementWindow
	if changed {
		transitions = append(transitions, provenanceHistoryTransition(t, true))
		feedbackWindow = selfImprovementFeedbackWindow(t, 8, 2, true)
	} else {
		feedbackWindow = selfImprovementFeedbackWindow(t, 4, 4, false)
	}
	history, err := ObserveRevisionSelfImprovementProvenanceHistory(transitions)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceHistory() error = %v", err)
	}
	decision, err := ObserveRevisionSelfImprovementProvenanceDecision(history)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceDecision() error = %v", err)
	}
	feedback, err := ObserveRevisionSelfImprovementFeedback(feedbackWindow)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementFeedback() error = %v", err)
	}
	return history, decision, feedback
}

func TestObserveRevisionSelfImprovementProvenanceFeedbackBridgeBindsStable(t *testing.T) {
	history, decision, feedback := provenanceFeedbackBridgeInputs(t, false)
	result, err := ObserveRevisionSelfImprovementProvenanceFeedbackBridge(history, decision, feedback)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceFeedbackBridge() error = %v", err)
	}
	if result.Status != "BOUND" || result.BridgeSignal != "provenance-feedback-bridge-bound" ||
		result.HistorySignal != "stable" || result.DecisionSignal != "observe" ||
		result.FeedbackSignal != "observe" {
		t.Fatalf("unexpected stable bridge: %#v", result)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementProvenanceFeedbackBridgeBindsTransition(t *testing.T) {
	history, decision, feedback := provenanceFeedbackBridgeInputs(t, true)
	result, err := ObserveRevisionSelfImprovementProvenanceFeedbackBridge(history, decision, feedback)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceFeedbackBridge() error = %v", err)
	}
	if result.Status != "BOUND" || result.HistorySignal != "mixed" ||
		result.ComparisonSignal != "mixed" || result.DecisionSignal != "inspect" ||
		result.FeedbackSignal != "inspect" || !result.RequiresMeasurement {
		t.Fatalf("unexpected transition bridge: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementProvenanceFeedbackBridgeRetainsMappingFailure(t *testing.T) {
	history, decision, _ := provenanceFeedbackBridgeInputs(t, false)
	_, _, feedback := provenanceFeedbackBridgeInputs(t, true)
	result, err := ObserveRevisionSelfImprovementProvenanceFeedbackBridge(history, decision, feedback)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementProvenanceFeedbackBridge() error = nil, want mapping failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-provenance-feedback-bridge-mapping" {
		t.Fatalf("unexpected unknown bridge: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementProvenanceFeedbackBridgeRetainsDecisionFailure(t *testing.T) {
	history, decision, feedback := provenanceFeedbackBridgeInputs(t, false)
	decision.DecisionDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementProvenanceFeedbackBridge(history, decision, feedback)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementProvenanceFeedbackBridge() error = nil, want decision failure")
	}
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-provenance-feedback-bridge-decision" {
		t.Fatalf("unexpected unknown decision bridge: %#v", result)
	}
}

func TestObserveRevisionSelfImprovementProvenanceFeedbackBridgeIsDeterministic(t *testing.T) {
	history, decision, feedback := provenanceFeedbackBridgeInputs(t, true)
	first, err := ObserveRevisionSelfImprovementProvenanceFeedbackBridge(history, decision, feedback)
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementProvenanceFeedbackBridge() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementProvenanceFeedbackBridge(history, decision, feedback)
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementProvenanceFeedbackBridge() error = %v", err)
	}
	if first.BridgeDigest != second.BridgeDigest {
		t.Fatal("same provenance history, decision, and feedback produced different bridge digest")
	}
}
