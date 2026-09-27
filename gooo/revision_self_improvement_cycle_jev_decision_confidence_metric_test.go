package gooo

import "testing"

func TestObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(t *testing.T) {
	cases := []struct {
		name            string
		kind            string
		confidenceMilli int64
		band            string
	}{
		{name: "high-choice", kind: "choice", confidenceMilli: 920, band: "high"},
		{name: "medium-score", kind: "score", confidenceMilli: 730, band: "medium"},
		{name: "low-score", kind: "score", confidenceMilli: 120, band: "low"},
		{name: "noul", kind: "noul", confidenceMilli: -1, band: "none"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			signal, err := ObserveRevisionSelfImprovementCycleJEVTypedDecisionSignal(
				RevisionSelfImprovementCycleJEVTypedDecisionSignalObservationInput{
					QuestionID:       "question:" + tc.name,
					QuestionKind:     tc.kind,
					DecisionDigest:   digestString("decision:" + tc.name),
					EvidenceDigest:   digestString("evidence:" + tc.name),
					ConfidenceMilli:  tc.confidenceMilli,
				},
			)
			if err != nil {
				t.Fatalf("signal observation returned an error: %v", err)
			}
			metric, err := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(signal)
			if err != nil {
				t.Fatalf("confidence metric returned an error: %v", err)
			}
			if metric.Status != "BOUND" || metric.ConfidenceBand != tc.band {
				t.Fatalf("unexpected confidence metric: %#v", metric)
			}
			if metric.LinkedDigestCount != 3 || metric.RequiredDigestCount != 3 {
				t.Fatalf("expected three linked digests, got %#v", metric)
			}
			if err := metric.Validate(); err != nil {
				t.Fatalf("confidence metric did not validate: %v", err)
			}
		})
	}
}

func TestObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetricUnknown(t *testing.T) {
	signal := RevisionSelfImprovementCycleJEVTypedDecisionSignalObservation{
		Status:                    "UNKNOWN",
		MissingStage:              "revision-self-improvement-cycle-jev-typed-decision-signal-input",
		QuestionID:                "question:route",
		QuestionKind:              "choice",
		DecisionDigest:            digestString("decision:route"),
		EvidenceDigest:            digestString("evidence:route"),
		ConfidenceMilli:           920,
		DecisionSignal:            "jev-typed-decision-unknown",
		DecisionObservationDigest: digestString("unknown-observation"),
		NonExecuting:              true,
		NonAuthorizing:            true,
	}
	metric, err := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(signal)
	if err == nil {
		t.Fatal("expected unknown signal to be rejected")
	}
	if metric.Status != "UNKNOWN" || metric.MissingStage == "" {
		t.Fatalf("expected unknown metric with a missing stage, got %#v", metric)
	}
}

func TestRevisionSelfImprovementCycleJEVDecisionConfidenceMetricRejectsTampering(t *testing.T) {
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
		t.Fatalf("signal observation returned an error: %v", err)
	}
	metric, err := ObserveRevisionSelfImprovementCycleJEVDecisionConfidenceMetric(signal)
	if err != nil {
		t.Fatalf("confidence metric returned an error: %v", err)
	}
	metric.ConfidenceBand = "low"
	if err := metric.Validate(); err == nil {
		t.Fatal("expected tampered confidence band to be rejected")
	}
}
