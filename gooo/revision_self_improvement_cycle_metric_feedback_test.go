package gooo

import "testing"

func cycleMetricObservation(t *testing.T) RevisionSelfImprovementCycleMetricObservation {
	t.Helper()
	cycle := cycleMetricCycle(t)
	metric, err := ObserveRevisionSelfImprovementCycleMetric(cycle, cycleMetricInput(t))
	if err != nil {
		t.Fatalf("observe cycle metric: %v", err)
	}
	return metric
}

func TestObserveRevisionSelfImprovementCycleMetricFeedbackImproved(t *testing.T) {
	metric := cycleMetricObservation(t)
	feedback, err := ObserveRevisionSelfImprovementCycleMetricFeedback(metric)
	if err != nil {
		t.Fatalf("observe improved metric feedback: %v", err)
	}
	if feedback.Status != "BOUND" || feedback.FeedbackSignal != "observe" ||
		feedback.FeedbackReason != "metric-improved-requires-follow-up-observation" ||
		feedback.RequiresReview || feedback.RequiresReplan || feedback.RequiresMeasurement {
		t.Fatalf("unexpected improved metric feedback: %#v", feedback)
	}
	if err := feedback.Validate(); err != nil {
		t.Fatalf("validate improved metric feedback: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleMetricFeedbackRegression(t *testing.T) {
	metric := cycleMetricObservation(t)
	metric.Direction = "lower-is-better"
	metric.BaselineValue = 12
	metric.CandidateValue = 15
	metric.MetricDelta = 3
	metric.MetricEvidenceDigest = digestRevisionSelfImprovementCycleMetricEvidence(metric)
	metric.ObservationDigest = digestRevisionSelfImprovementCycleMetric(metric)
	feedback, err := ObserveRevisionSelfImprovementCycleMetricFeedback(metric)
	if err != nil {
		t.Fatalf("observe regressed metric feedback: %v", err)
	}
	if feedback.FeedbackSignal != "replan" || !feedback.RequiresReplan ||
		feedback.RequiresReview || feedback.RequiresMeasurement {
		t.Fatalf("unexpected regressed metric feedback: %#v", feedback)
	}
}

func TestObserveRevisionSelfImprovementCycleMetricFeedbackStable(t *testing.T) {
	metric := cycleMetricObservation(t)
	metric.BaselineValue = 12
	metric.CandidateValue = 12
	metric.MetricDelta = 0
	metric.MetricSignal = "metric-stable"
	metric.MetricEvidenceDigest = digestRevisionSelfImprovementCycleMetricEvidence(metric)
	metric.ObservationDigest = digestRevisionSelfImprovementCycleMetric(metric)
	feedback, err := ObserveRevisionSelfImprovementCycleMetricFeedback(metric)
	if err != nil {
		t.Fatalf("observe stable metric feedback: %v", err)
	}
	if feedback.FeedbackSignal != "review" || !feedback.RequiresReview ||
		!feedback.RequiresMeasurement || feedback.RequiresReplan {
		t.Fatalf("unexpected stable metric feedback: %#v", feedback)
	}
}

func TestObserveRevisionSelfImprovementCycleMetricFeedbackPreservesUnknown(t *testing.T) {
	metric := cycleMetricObservation(t)
	metric.ObservationDigest = digestString("tampered")
	feedback, err := ObserveRevisionSelfImprovementCycleMetricFeedback(metric)
	if err == nil {
		t.Fatal("expected metric validation error")
	}
	if feedback.Status != "UNKNOWN" ||
		feedback.MissingStage != "revision-self-improvement-cycle-metric-feedback-metric" {
		t.Fatalf("unexpected unknown metric feedback: %#v", feedback)
	}
}

func TestObserveRevisionSelfImprovementCycleMetricFeedbackIsDeterministic(t *testing.T) {
	metric := cycleMetricObservation(t)
	first, err := ObserveRevisionSelfImprovementCycleMetricFeedback(metric)
	if err != nil {
		t.Fatalf("first metric feedback: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCycleMetricFeedback(metric)
	if err != nil {
		t.Fatalf("second metric feedback: %v", err)
	}
	if first != second {
		t.Fatalf("metric feedback differs: %#v != %#v", first, second)
	}
}