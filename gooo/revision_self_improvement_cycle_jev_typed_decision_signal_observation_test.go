package gooo

import "testing"

func TestObserveRevisionSelfImprovementCycleJEVTypedDecisionSignal(t *testing.T) {
	cases := []struct {
		name           string
		kind           string
		confidenceMilli int64
		signal         string
	}{
		{name: "choice", kind: "choice", confidenceMilli: 920, signal: "jev-choice-observed"},
		{name: "score", kind: "score", confidenceMilli: 730, signal: "jev-score-observed"},
		{name: "noul", kind: "noul", confidenceMilli: -1, signal: "jev-noul-observed"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			observation, err := ObserveRevisionSelfImprovementCycleJEVTypedDecisionSignal(
				RevisionSelfImprovementCycleJEVTypedDecisionSignalObservationInput{
					QuestionID:       "question:" + tc.name,
					QuestionKind:     tc.kind,
					DecisionDigest:   digestString("decision:" + tc.name),
					EvidenceDigest:   digestString("evidence:" + tc.name),
					ConfidenceMilli:  tc.confidenceMilli,
				},
			)
			if err != nil {
				t.Fatalf("observe returned an error: %v", err)
			}
			if observation.Status != "BOUND" || observation.DecisionSignal != tc.signal {
				t.Fatalf("unexpected bound observation: %#v", observation)
			}
			if err := observation.Validate(); err != nil {
				t.Fatalf("observation did not validate: %v", err)
			}
		})
	}
}

func TestObserveRevisionSelfImprovementCycleJEVTypedDecisionSignalUnknown(t *testing.T) {
	observation, err := ObserveRevisionSelfImprovementCycleJEVTypedDecisionSignal(
		RevisionSelfImprovementCycleJEVTypedDecisionSignalObservationInput{
			QuestionID:       "question:route",
			QuestionKind:     "choice",
			DecisionDigest:   digestString("decision:route"),
			EvidenceDigest:   digestString("evidence:route"),
			ConfidenceMilli:  1200,
		},
	)
	if err == nil {
		t.Fatal("expected out-of-range confidence to return an error")
	}
	if observation.Status != "UNKNOWN" || observation.MissingStage == "" {
		t.Fatalf("expected unknown observation with a missing stage, got %#v", observation)
	}
}

func TestRevisionSelfImprovementCycleJEVTypedDecisionSignalRejectsTampering(t *testing.T) {
	observation, err := ObserveRevisionSelfImprovementCycleJEVTypedDecisionSignal(
		RevisionSelfImprovementCycleJEVTypedDecisionSignalObservationInput{
			QuestionID:       "question:route",
			QuestionKind:     "choice",
			DecisionDigest:   digestString("decision:route"),
			EvidenceDigest:   digestString("evidence:route"),
			ConfidenceMilli:  920,
		},
	)
	if err != nil {
		t.Fatalf("observe returned an error: %v", err)
	}
	observation.ConfidenceMilli = 921
	if err := observation.Validate(); err == nil {
		t.Fatal("expected tampered confidence to be rejected")
	}
}
