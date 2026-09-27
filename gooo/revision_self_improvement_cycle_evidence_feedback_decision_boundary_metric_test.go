package gooo

import "testing"

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricBindsFourLinks(t *testing.T) {
	cycle, evidenceFeedback, boundary := cycleEvidenceCoverageFeedbackDecisionBoundaryInputs(t)
	gatedBoundary, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(cycle, evidenceFeedback, boundary)
	if err != nil {
		t.Fatalf("observe cycle evidence feedback decision boundary: %v", err)
	}
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric(gatedBoundary)
	if err != nil {
		t.Fatalf("observe cycle evidence feedback decision boundary metric: %v", err)
	}
	if got.Status != "BOUND" ||
		got.MetricSignal != "boundary-complete" ||
		got.MetricName != "evidence-feedback-decision-boundary-provenance-link-count" {
		t.Fatalf("unexpected evidence feedback decision boundary metric: %#v", got)
	}
	if got.LinkedStageCount != 4 || got.RequiredStageCount != 4 || got.LinkedDigestCount != 4 {
		t.Fatalf("unexpected evidence feedback decision boundary metric counts: %#v", got)
	}
	if got.BoundaryDigest != gatedBoundary.ObservationDigest ||
		!got.NonExecuting ||
		!got.NonAuthorizing {
		t.Fatalf("metric lost provenance or safety: %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate evidence feedback decision boundary metric: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricPreservesUnknown(t *testing.T) {
	cycle, evidenceFeedback, boundary := cycleEvidenceCoverageFeedbackDecisionBoundaryInputs(t)
	gatedBoundary, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(cycle, evidenceFeedback, boundary)
	if err != nil {
		t.Fatalf("observe cycle evidence feedback decision boundary: %v", err)
	}
	gatedBoundary.ObservationDigest = digestString("tampered")
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric(gatedBoundary)
	if err == nil {
		t.Fatal("expected tampered boundary error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-boundary" {
		t.Fatalf("unexpected unknown evidence feedback decision boundary metric: %#v", got)
	}
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricIsDeterministic(t *testing.T) {
	cycle, evidenceFeedback, boundary := cycleEvidenceCoverageFeedbackDecisionBoundaryInputs(t)
	gatedBoundary, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(cycle, evidenceFeedback, boundary)
	if err != nil {
		t.Fatalf("observe cycle evidence feedback decision boundary: %v", err)
	}
	first, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric(gatedBoundary)
	if err != nil {
		t.Fatalf("first evidence feedback decision boundary metric: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric(gatedBoundary)
	if err != nil {
		t.Fatalf("second evidence feedback decision boundary metric: %v", err)
	}
	if first != second {
		t.Fatalf("evidence feedback decision boundary metrics differ: %#v != %#v", first, second)
	}
}
