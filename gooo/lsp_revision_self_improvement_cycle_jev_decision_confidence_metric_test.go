package gooo

import "testing"

func TestObserveLSPRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(
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

	projection, err :=
		ObserveLSPRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(metric)
	if err != nil {
		t.Fatalf("lsp projection returned an error: %v", err)
	}
	if projection.Status != "BOUND" ||
		projection.ProjectionSignal != "jev-decision-confidence-lsp-projected" ||
		projection.ConfidenceBand != "high" {
		t.Fatalf("unexpected lsp confidence projection: %#v", projection)
	}
	if !projection.ReadOnly || !projection.NonExecuting || !projection.NonAuthorizing {
		t.Fatalf("lsp projection crossed an authority boundary: %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("lsp projection did not validate: %v", err)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleJEVDecisionConfidenceMetricUnknown(
	t *testing.T,
) {
	metric := RevisionSelfImprovementCycleJEVDecisionConfidenceMetricObservation{
		Status:                    "UNKNOWN",
		MissingStage:              "revision-self-improvement-cycle-jev-decision-confidence-metric-input",
		MetricName:                "jev-typed-decision-confidence-milli",
		QuestionID:                "question:route",
		QuestionKind:              "choice",
		ConfidenceMilli:           1200,
		ConfidenceBand:            "unknown",
		DecisionDigest:            digestString("decision:route"),
		EvidenceDigest:            digestString("evidence:route"),
		DecisionObservationDigest: digestString("observation:route"),
		MetricSignal:               "jev-decision-confidence-unknown",
		MetricDigest:              digestString("metric:unknown"),
		ReadOnly:                  true,
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	projection, err :=
		ObserveLSPRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(metric)
	if err == nil {
		t.Fatal("expected unknown metric to remain unprojectable")
	}
	if projection.Status != "UNKNOWN" || projection.MissingStage == "" {
		t.Fatalf("expected unknown lsp projection with a missing stage, got %#v", projection)
	}
}

func TestLSPRevisionSelfImprovementCycleJEVDecisionConfidenceMetricRejectsTampering(
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
	projection, err :=
		ObserveLSPRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(metric)
	if err != nil {
		t.Fatalf("lsp projection returned an error: %v", err)
	}

	projection.ConfidenceBand = "low"
	if err := projection.Validate(); err == nil {
		t.Fatal("expected tampered confidence band to be rejected")
	}
}
