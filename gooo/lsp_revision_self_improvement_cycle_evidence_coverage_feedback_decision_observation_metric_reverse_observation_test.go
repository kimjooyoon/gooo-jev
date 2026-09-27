package gooo

import "testing"

func lspEvidenceFeedbackDecisionObservationReverseInput(t *testing.T) RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation {
	t.Helper()
	metric := evidenceFeedbackDecisionObservationReverseInput(t)
	observation, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation(metric)
	if err != nil {
		t.Fatalf("observe decision observation metric reverse observation: %v", err)
	}
	return observation
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationBindsProjection(t *testing.T) {
	observation := lspEvidenceFeedbackDecisionObservationReverseInput(t)
	got, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation(observation)
	if err != nil {
		t.Fatalf("observe lsp decision observation metric reverse observation: %v", err)
	}
	if got.Status != "BOUND" ||
		got.ReverseSignal != "decision-observation-metric-reverse-complete" ||
		got.MetricName != "evidence-feedback-decision-observation-provenance-link-count" {
		t.Fatalf("unexpected lsp decision observation metric reverse observation: %#v", got)
	}
	if got.MetricObservationDigest != observation.MetricObservationDigest ||
		got.MetricDigest != observation.MetricDigest ||
		got.ReconstructedMetricDigest != observation.ReconstructedMetricDigest ||
		got.DecisionObservationDigest != observation.DecisionObservationDigest ||
		got.ObservationDigest != observation.ObservationDigest {
		t.Fatalf("projection lost reverse observation provenance: %#v", got)
	}
	if !got.ReadOnly || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("projection safety flags = %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate lsp decision observation metric reverse observation: %v", err)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationPreservesUnknown(t *testing.T) {
	observation := lspEvidenceFeedbackDecisionObservationReverseInput(t)
	observation.ObservationDigest = digestString("tampered")
	got, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation(observation)
	if err == nil {
		t.Fatal("expected tampered reverse observation error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "lsp-revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-reverse-observation-input" {
		t.Fatalf("unexpected unknown lsp decision observation metric reverse observation: %#v", got)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservationIsDeterministic(t *testing.T) {
	observation := lspEvidenceFeedbackDecisionObservationReverseInput(t)
	first, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation(observation)
	if err != nil {
		t.Fatalf("first lsp decision observation metric reverse observation: %v", err)
	}
	second, err := ObserveLSPRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricReverseObservation(observation)
	if err != nil {
		t.Fatalf("second lsp decision observation metric reverse observation: %v", err)
	}
	if first != second {
		t.Fatalf("lsp decision observation metric reverse observations differ: %#v != %#v", first, second)
	}
}
