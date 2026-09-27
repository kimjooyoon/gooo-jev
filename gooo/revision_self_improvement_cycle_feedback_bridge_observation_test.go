package gooo

import "testing"

func cycleFeedbackBridgeInputs(t *testing.T) (
	RevisionSelfImprovementCycleObservation,
	RevisionSelfImprovementCycleMetricObservation,
	RevisionSelfImprovementCycleMetricFeedback,
) {
	t.Helper()
	cycle := cycleMetricCycle(t)
	metric, err := ObserveRevisionSelfImprovementCycleMetric(cycle, cycleMetricInput(t))
	if err != nil {
		t.Fatalf("observe cycle metric: %v", err)
	}
	feedback, err := ObserveRevisionSelfImprovementCycleMetricFeedback(metric)
	if err != nil {
		t.Fatalf("observe cycle metric feedback: %v", err)
	}
	return cycle, metric, feedback
}

func TestObserveRevisionSelfImprovementCycleFeedbackBridgeBindsLinks(t *testing.T) {
	cycle, metric, feedback := cycleFeedbackBridgeInputs(t)
	got, err := ObserveRevisionSelfImprovementCycleFeedbackBridge(cycle, metric, feedback)
	if err != nil {
		t.Fatalf("observe cycle feedback bridge: %v", err)
	}
	if got.Status != "BOUND" || got.BridgeSignal != "cycle-feedback-linked" || !got.SignalsAligned {
		t.Fatalf("unexpected cycle feedback bridge: %#v", got)
	}
	if got.CycleDigest != cycle.ObservationDigest ||
		got.MetricDigest != metric.ObservationDigest ||
		got.FeedbackDigest != feedback.FeedbackDigest {
		t.Fatalf("bridge lost stage links: %#v", got)
	}
	if got.FeedbackSignal != feedback.FeedbackSignal ||
		got.MetricSignal != metric.MetricSignal ||
		!got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("bridge lost feedback semantics: %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate cycle feedback bridge: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleFeedbackBridgePreservesUnknownLink(t *testing.T) {
	cycle, metric, feedback := cycleFeedbackBridgeInputs(t)
	feedback.CycleMetricDigest = digestString("tampered")
	feedback.FeedbackDigest = digestRevisionSelfImprovementCycleMetricFeedback(feedback)
	got, err := ObserveRevisionSelfImprovementCycleFeedbackBridge(cycle, metric, feedback)
	if err == nil {
		t.Fatal("expected feedback metric link error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "revision-self-improvement-cycle-feedback-bridge-feedback-metric-link" {
		t.Fatalf("unexpected unknown cycle feedback bridge: %#v", got)
	}
}

func TestObserveRevisionSelfImprovementCycleFeedbackBridgeIsDeterministic(t *testing.T) {
	cycle, metric, feedback := cycleFeedbackBridgeInputs(t)
	first, err := ObserveRevisionSelfImprovementCycleFeedbackBridge(cycle, metric, feedback)
	if err != nil {
		t.Fatalf("first cycle feedback bridge: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCycleFeedbackBridge(cycle, metric, feedback)
	if err != nil {
		t.Fatalf("second cycle feedback bridge: %v", err)
	}
	if first != second {
		t.Fatalf("cycle feedback bridges differ: %#v != %#v", first, second)
	}
}