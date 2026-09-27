package gooo

import "testing"

func TestObserveRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetricBindsLinkCount(t *testing.T) {
	cycle, bridge, link := cycleFeedbackDecisionBoundaryInputs(t)
	boundary, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionBoundary(cycle, bridge, link)
	if err != nil {
		t.Fatalf("observe cycle feedback decision boundary: %v", err)
	}
	got, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetric(boundary)
	if err != nil {
		t.Fatalf("observe cycle feedback decision boundary metric: %v", err)
	}
	if got.Status != "BOUND" || got.MetricSignal != "boundary-complete" || got.MetricName != "provenance-link-count" {
		t.Fatalf("unexpected boundary metric: %#v", got)
	}
	if got.LinkedStageCount != 3 || got.RequiredStageCount != 3 || got.LinkedDigestCount != 3 {
		t.Fatalf("unexpected boundary metric counts: %#v", got)
	}
	if got.BoundaryDigest != boundary.ObservationDigest || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("boundary metric lost provenance or safety: %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate boundary metric: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetricPreservesUnknown(t *testing.T) {
	cycle, bridge, link := cycleFeedbackDecisionBoundaryInputs(t)
	boundary, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionBoundary(cycle, bridge, link)
	if err != nil {
		t.Fatalf("observe cycle feedback decision boundary: %v", err)
	}
	boundary.ObservationDigest = digestString("tampered")
	got, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetric(boundary)
	if err == nil {
		t.Fatal("expected tampered boundary error")
	}
	if got.Status != "UNKNOWN" || got.MissingStage != "revision-self-improvement-cycle-feedback-decision-boundary-metric-boundary" {
		t.Fatalf("unexpected unknown boundary metric: %#v", got)
	}
}

func TestObserveRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetricIsDeterministic(t *testing.T) {
	cycle, bridge, link := cycleFeedbackDecisionBoundaryInputs(t)
	boundary, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionBoundary(cycle, bridge, link)
	if err != nil {
		t.Fatalf("observe cycle feedback decision boundary: %v", err)
	}
	first, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetric(boundary)
	if err != nil {
		t.Fatalf("first boundary metric: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionBoundaryMetric(boundary)
	if err != nil {
		t.Fatalf("second boundary metric: %v", err)
	}
	if first != second {
		t.Fatalf("boundary metrics differ: %#v != %#v", first, second)
	}
}
