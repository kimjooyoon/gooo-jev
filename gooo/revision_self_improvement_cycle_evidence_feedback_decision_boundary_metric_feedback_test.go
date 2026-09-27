package gooo

import "testing"

func evidenceFeedbackDecisionBoundaryMetricInput(t *testing.T) RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricObservation {
	t.Helper()
	cycle, evidenceFeedback, boundary := cycleEvidenceCoverageFeedbackDecisionBoundaryInputs(t)
	gatedBoundary, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundary(cycle, evidenceFeedback, boundary)
	if err != nil {
		t.Fatalf("observe cycle evidence feedback decision boundary: %v", err)
	}
	metric, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetric(gatedBoundary)
	if err != nil {
		t.Fatalf("observe evidence feedback decision boundary metric: %v", err)
	}
	return metric
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackBindsObservation(t *testing.T) {
	metric := evidenceFeedbackDecisionBoundaryMetricInput(t)
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback(metric)
	if err != nil {
		t.Fatalf("observe evidence feedback decision boundary metric feedback: %v", err)
	}
	if got.Status != "BOUND" ||
		got.FeedbackSignal != "observe" ||
		got.FeedbackReason != "complete boundary requires follow-up observation" ||
		got.MetricSignal != "boundary-complete" ||
		got.RequiresReview ||
		got.RequiresReplan ||
		got.RequiresMeasurement {
		t.Fatalf("unexpected evidence feedback decision boundary metric feedback: %#v", got)
	}
	if got.BoundaryMetricDigest != metric.ObservationDigest ||
		!got.RequiresObservation ||
		!got.NonExecuting ||
		!got.NonAuthorizing {
		t.Fatalf("feedback lost metric provenance or safety: %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate evidence feedback decision boundary metric feedback: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackPreservesUnknown(t *testing.T) {
	metric := evidenceFeedbackDecisionBoundaryMetricInput(t)
	metric.ObservationDigest = digestString("tampered")
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback(metric)
	if err == nil {
		t.Fatal("expected tampered boundary metric error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "revision-self-improvement-cycle-evidence-coverage-feedback-decision-boundary-metric-feedback-metric" {
		t.Fatalf("unexpected unknown evidence feedback decision boundary metric feedback: %#v", got)
	}
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackIsDeterministic(t *testing.T) {
	metric := evidenceFeedbackDecisionBoundaryMetricInput(t)
	first, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback(metric)
	if err != nil {
		t.Fatalf("first evidence feedback decision boundary metric feedback: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedback(metric)
	if err != nil {
		t.Fatalf("second evidence feedback decision boundary metric feedback: %v", err)
	}
	if first != second {
		t.Fatalf("evidence feedback decision boundary metric feedback differs: %#v != %#v", first, second)
	}
}
