package gooo

import "testing"

func evidenceFeedbackDecisionObservationReverseInput(t *testing.T) RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricObservation {
	t.Helper()
	decision := evidenceFeedbackDecisionObservationMetricInput(t)
	metric, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(decision)
	if err != nil {
		t.Fatalf("observe decision observation metric: %v", err)
	}
	return metric
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationBinds(t *testing.T) {
	metric := evidenceFeedbackDecisionObservationReverseInput(t)
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation(metric)
	if err != nil {
		t.Fatalf("observe decision observation metric reverse observation: %v", err)
	}
	if got.Status != "BOUND" ||
		got.ReverseSignal != "decision-observation-metric-reverse-complete" ||
		got.MetricName != "evidence-feedback-decision-observation-provenance-link-count" {
		t.Fatalf("unexpected decision observation reverse observation: %#v", got)
	}
	if got.ReconstructedMetricDigest != metric.MetricDigest ||
		got.MetricObservationDigest != metric.ObservationDigest ||
		got.DecisionObservationDigest != metric.DecisionObservationDigest {
		t.Fatalf("reverse observation lost metric provenance: %#v", got)
	}
	if !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("reverse observation safety flags = %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate decision observation reverse observation: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationPreservesUnknown(t *testing.T) {
	metric := evidenceFeedbackDecisionObservationReverseInput(t)
	metric.MetricDigest = digestString("tampered")
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation(metric)
	if err == nil {
		t.Fatal("expected tampered metric error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation-input" {
		t.Fatalf("unexpected unknown decision observation reverse observation: %#v", got)
	}
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationIsDeterministic(t *testing.T) {
	metric := evidenceFeedbackDecisionObservationReverseInput(t)
	first, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation(metric)
	if err != nil {
		t.Fatalf("first decision observation reverse observation: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation(metric)
	if err != nil {
		t.Fatalf("second decision observation reverse observation: %v", err)
	}
	if first != second {
		t.Fatalf("decision observation reverse observations differ: %#v != %#v", first, second)
	}
}
