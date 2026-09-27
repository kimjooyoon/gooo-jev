package gooo

import "testing"

func evidenceFeedbackDecisionObservationReverseClosureInput(t *testing.T) RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation {
	t.Helper()
	metric := evidenceFeedbackDecisionObservationReverseInput(t)
	observation, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation(metric)
	if err != nil {
		t.Fatalf("observe decision observation metric reverse observation: %v", err)
	}
	return observation
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetricBinds(t *testing.T) {
	observation := evidenceFeedbackDecisionObservationReverseClosureInput(t)
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric(observation)
	if err != nil {
		t.Fatalf("observe reverse observation provenance closure metric: %v", err)
	}
	if got.Status != "BOUND" ||
		got.LinkedDigestCount != 5 ||
		got.RequiredDigestCount != 5 ||
		got.MetricSignal != "reverse-observation-provenance-complete" {
		t.Fatalf("unexpected reverse observation provenance closure metric: %#v", got)
	}
	if got.MetricObservationDigest != observation.MetricObservationDigest ||
		got.MetricDigest != observation.MetricDigest ||
		got.ReconstructedMetricDigest != observation.ReconstructedMetricDigest ||
		got.DecisionObservationDigest != observation.DecisionObservationDigest ||
		got.ReverseObservationDigest != observation.ObservationDigest {
		t.Fatalf("closure metric lost provenance: %#v", got)
	}
	if !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("closure metric safety flags = %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate reverse observation provenance closure metric: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetricPreservesUnknown(t *testing.T) {
	observation := evidenceFeedbackDecisionObservationReverseClosureInput(t)
	observation.MetricObservationDigest = digestString("tampered")
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric(observation)
	if err == nil {
		t.Fatal("expected tampered reverse observation error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation-provenance-closure-metric-input" {
		t.Fatalf("unexpected unknown reverse observation provenance closure metric: %#v", got)
	}
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetricIsDeterministic(t *testing.T) {
	observation := evidenceFeedbackDecisionObservationReverseClosureInput(t)
	first, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric(observation)
	if err != nil {
		t.Fatalf("first reverse observation provenance closure metric: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationProvenanceClosureMetric(observation)
	if err != nil {
		t.Fatalf("second reverse observation provenance closure metric: %v", err)
	}
	if first != second {
		t.Fatalf("reverse observation provenance closure metrics differ: %#v != %#v", first, second)
	}
}
