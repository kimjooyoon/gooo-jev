package gooo

import "testing"

func iterationProvenanceInputs(t *testing.T, changed bool) (
	RevisionSelfImprovementIteration,
	RevisionSelfImprovementProvenanceHistoryObservation,
	RevisionSelfImprovementProvenanceFeedbackBridgeObservation,
) {
	t.Helper()
	metricCount, generationCount := 4, 4
	if changed {
		metricCount, generationCount = 8, 2
	}
	window := selfImprovementFeedbackWindow(t, metricCount, generationCount, changed)
	feedback, err := ObserveRevisionSelfImprovementFeedback(window)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementFeedback() error = %v", err)
	}
	legacyHistory, err := ObserveRevisionSelfImprovementHistory(
		[]RevisionSelfImprovementWindow{window},
		[]RevisionSelfImprovementFeedback{feedback},
	)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementHistory() error = %v", err)
	}
	iteration, err := ObserveRevisionSelfImprovementIteration(validContract, legacyHistory, feedback)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementIteration() error = %v", err)
	}

	transitions := []RevisionSelfImprovementProvenanceTransitionObservation{
		provenanceHistoryTransition(t, false),
	}
	if changed {
		transitions = append(transitions, provenanceHistoryTransition(t, true))
	}
	provenanceHistory, err := ObserveRevisionSelfImprovementProvenanceHistory(transitions)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceHistory() error = %v", err)
	}
	provenanceDecision, err := ObserveRevisionSelfImprovementProvenanceDecision(provenanceHistory)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceDecision() error = %v", err)
	}
	bridge, err := ObserveRevisionSelfImprovementProvenanceFeedbackBridge(
		provenanceHistory,
		provenanceDecision,
		feedback,
	)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementProvenanceFeedbackBridge() error = %v", err)
	}
	return iteration, provenanceHistory, bridge
}

func TestObserveRevisionSelfImprovementIterationProvenanceBindsStable(t *testing.T) {
	iteration, history, bridge := iterationProvenanceInputs(t, false)
	observation, err := ObserveRevisionSelfImprovementIterationProvenance(iteration, history, bridge)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementIterationProvenance() error = %v", err)
	}
	if observation.Status != "BOUND" ||
		observation.ObservationSignal != "iteration-provenance-bound" ||
		observation.ProvenanceHistorySignal != "stable" ||
		observation.DecisionSignal != "observe" ||
		observation.FeedbackSignal != "observe" ||
		observation.IRDigest == "" {
		t.Fatalf("unexpected stable iteration provenance observation: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestObserveRevisionSelfImprovementIterationProvenanceBindsChanged(t *testing.T) {
	iteration, history, bridge := iterationProvenanceInputs(t, true)
	observation, err := ObserveRevisionSelfImprovementIterationProvenance(iteration, history, bridge)
	if err != nil {
		t.Fatalf("ObserveRevisionSelfImprovementIterationProvenance() error = %v", err)
	}
	if observation.Status != "BOUND" ||
		observation.ProvenanceHistorySignal != "mixed" ||
		observation.LegacyHistorySignal != "mixed" ||
		observation.DecisionSignal != "inspect" ||
		observation.FeedbackSignal != "inspect" ||
		!observation.RequiresMeasurement ||
		observation.GenerationChangeCount == 0 {
		t.Fatalf("unexpected changed iteration provenance observation: %#v", observation)
	}
}

func TestObserveRevisionSelfImprovementIterationProvenanceRetainsFeedbackLinkFailure(t *testing.T) {
	iteration, history, bridge := iterationProvenanceInputs(t, false)
	iteration.FeedbackDigest = digestString("tampered")
	iteration.IterationDigest = digestRevisionSelfImprovementIteration(iteration)
	observation, err := ObserveRevisionSelfImprovementIterationProvenance(iteration, history, bridge)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementIterationProvenance() error = nil, want feedback link failure")
	}
	if observation.Status != "UNKNOWN" ||
		observation.MissingStage != "revision-self-improvement-iteration-provenance-feedback-link" {
		t.Fatalf("unexpected unknown feedback-link observation: %#v", observation)
	}
}

func TestObserveRevisionSelfImprovementIterationProvenanceRetainsHistoryFailure(t *testing.T) {
	iteration, history, bridge := iterationProvenanceInputs(t, false)
	history.HistoryDigest = digestString("tampered")
	observation, err := ObserveRevisionSelfImprovementIterationProvenance(iteration, history, bridge)
	if err == nil {
		t.Fatal("ObserveRevisionSelfImprovementIterationProvenance() error = nil, want history failure")
	}
	if observation.Status != "UNKNOWN" ||
		observation.MissingStage != "revision-self-improvement-iteration-provenance-history" {
		t.Fatalf("unexpected unknown history observation: %#v", observation)
	}
}

func TestObserveRevisionSelfImprovementIterationProvenanceIsDeterministic(t *testing.T) {
	iteration, history, bridge := iterationProvenanceInputs(t, true)
	first, err := ObserveRevisionSelfImprovementIterationProvenance(iteration, history, bridge)
	if err != nil {
		t.Fatalf("first ObserveRevisionSelfImprovementIterationProvenance() error = %v", err)
	}
	second, err := ObserveRevisionSelfImprovementIterationProvenance(iteration, history, bridge)
	if err != nil {
		t.Fatalf("second ObserveRevisionSelfImprovementIterationProvenance() error = %v", err)
	}
	if first.ObservationDigest != second.ObservationDigest {
		t.Fatal("same iteration, provenance history, and bridge produced different observation digest")
	}
}
