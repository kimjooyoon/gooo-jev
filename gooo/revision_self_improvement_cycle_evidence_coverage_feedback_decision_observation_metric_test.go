package gooo

import "testing"

func evidenceFeedbackDecisionObservationMetricInput(t *testing.T) RevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecisionObservation {
	t.Helper()
	feedback := evidenceFeedbackDecisionObservationInput(t)
	decision, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionBoundaryMetricFeedbackDecision(feedback)
	if err != nil {
		t.Fatalf("observe evidence feedback decision: %v", err)
	}
	return decision
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricBindsMetric(t *testing.T) {
	decision := evidenceFeedbackDecisionObservationMetricInput(t)
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(decision)
	if err != nil {
		t.Fatalf("observe decision observation metric: %v", err)
	}
	if got.Status != "BOUND" ||
		got.MetricName != "evidence-feedback-decision-observation-provenance-link-count" ||
		got.MetricSignal != "decision-observation-complete" ||
		got.LinkedStageCount != 2 ||
		got.RequiredStageCount != 2 ||
		got.LinkedDigestCount != 3 {
		t.Fatalf("unexpected decision observation metric: %#v", got)
	}
	if got.DecisionObservationDigest != decision.ObservationDigest ||
		!got.NonExecuting ||
		!got.NonAuthorizing {
		t.Fatalf("metric lost decision observation provenance or safety: %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate decision observation metric: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricPreservesUnknown(t *testing.T) {
	decision := evidenceFeedbackDecisionObservationMetricInput(t)
	decision.DecisionDigest = digestString("tampered")
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(decision)
	if err == nil {
		t.Fatal("expected tampered decision error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "revision-self-improvement-cycle-evidence-coverage-feedback-decision-observation-metric-input" {
		t.Fatalf("unexpected unknown decision observation metric: %#v", got)
	}
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetricIsDeterministic(t *testing.T) {
	decision := evidenceFeedbackDecisionObservationMetricInput(t)
	first, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(decision)
	if err != nil {
		t.Fatalf("first decision observation metric: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageFeedbackDecisionObservationMetric(decision)
	if err != nil {
		t.Fatalf("second decision observation metric: %v", err)
	}
	if first != second {
		t.Fatalf("decision observation metrics differ: %#v != %#v", first, second)
	}
}
