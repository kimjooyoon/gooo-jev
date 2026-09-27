package gooo

import "testing"

func TestObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation(
	t *testing.T,
) {
	signal, err := ObserveRevisionSelfImprovementCycleJEVTypedDecisionSignal(
		RevisionSelfImprovementCycleJEVTypedDecisionSignalObservationInput{
			QuestionID:       "question:route",
			QuestionKind:     "choice",
			DecisionDigest:   digestString("decision:route"),
			EvidenceDigest:   digestString("evidence:route"),
			ConfidenceMilli:  920,
		},
	)
	if err != nil {
		t.Fatalf("source signal returned an error: %v", err)
	}
	metric, err := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(signal)
	if err != nil {
		t.Fatalf("source metric returned an error: %v", err)
	}

	observation, err :=
		ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation(metric)
	if err != nil {
		t.Fatalf("reverse observation returned an error: %v", err)
	}
	if observation.Status != "BOUND" ||
		observation.ReverseSignal != "jev-decision-confidence-metric-reverse-complete" {
		t.Fatalf("unexpected reverse observation: %#v", observation)
	}
	if observation.ReconstructedMetricDigest != observation.MetricDigest {
		t.Fatalf("metric reconstruction was not equal: %#v", observation)
	}
	if !observation.NonExecuting || !observation.NonAuthorizing {
		t.Fatalf("reverse observation crossed an authority boundary: %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("reverse observation did not validate: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationPreservesUnknown(
	t *testing.T,
) {
	metric := RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation{
		Status:                    "UNKNOWN",
		MissingStage:              "revision-self-improvement-cycle-jev-decision-confidence-metric-input",
		MetricName:                "jev-typed-decision-confidence-milli",
		QuestionID:                "question:route",
		QuestionKind:              "choice",
		ConfidenceMilli:           920,
		ConfidenceBand:            "unknown",
		DecisionDigest:            digestString("decision:route"),
		EvidenceDigest:            digestString("evidence:route"),
		DecisionObservationDigest: digestString("observation:route"),
		MetricDigest:              digestString("metric:unknown"),
		MetricSignal:              "jev-decision-confidence-unknown",
		ReadOnly:                  true,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	observation, err :=
		ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation(metric)
	if err == nil {
		t.Fatal("expected unknown confidence metric to return an error")
	}
	if observation.Status != "UNKNOWN" || observation.MissingStage == "" {
		t.Fatalf("expected unknown reverse observation with a missing stage, got %#v", observation)
	}
}

func TestRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationRejectsTampering(
	t *testing.T,
) {
	signal, err := ObserveRevisionSelfImprovementCycleJEVTypedDecisionSignal(
		RevisionSelfImprovementCycleJEVTypedDecisionSignalObservationInput{
			QuestionID:       "question:route",
			QuestionKind:     "choice",
			DecisionDigest:   digestString("decision:route"),
			EvidenceDigest:   digestString("evidence:route"),
			ConfidenceMilli:  920,
		},
	)
	if err != nil {
		t.Fatalf("source signal returned an error: %v", err)
	}
	metric, err := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(signal)
	if err != nil {
		t.Fatalf("source metric returned an error: %v", err)
	}
	metric.MetricDigest = digestString("tampered")
	observation, err :=
		ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservation(metric)
	if err == nil {
		t.Fatal("expected tampered confidence metric to return an error")
	}
	if observation.Status != "UNKNOWN" ||
		observation.MissingStage == "" {
		t.Fatalf("expected unknown reverse observation after tampering, got %#v", observation)
	}
}
