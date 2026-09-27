package gooo

import "testing"

func TestObserveLSPRevisionSelfImprovementCycleJEVTypedDecisionSignal(t *testing.T) {
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
		t.Fatalf("source observation returned an error: %v", err)
	}

	projection, err := ObserveLSPRevisionSelfImprovementCycleJEVTypedDecisionSignal(signal)
	if err != nil {
		t.Fatalf("lsp projection returned an error: %v", err)
	}
	if projection.Status != "BOUND" ||
		projection.ProjectionSignal != "jev-typed-decision-lsp-projected" {
		t.Fatalf("unexpected lsp projection: %#v", projection)
	}
	if !projection.ReadOnly || !projection.NonExecuting || !projection.NonAuthorizing {
		t.Fatalf("lsp projection crossed an authority boundary: %#v", projection)
	}
	if err := projection.Validate(); err != nil {
		t.Fatalf("lsp projection did not validate: %v", err)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleJEVTypedDecisionSignalUnknown(
	t *testing.T,
) {
	signal, err := ObserveRevisionSelfImprovementCycleJEVTypedDecisionSignal(
		RevisionSelfImprovementCycleJEVTypedDecisionSignalObservationInput{
			QuestionID:       "question:route",
			QuestionKind:     "choice",
			DecisionDigest:   digestString("decision:route"),
			EvidenceDigest:   digestString("evidence:route"),
			ConfidenceMilli:  1200,
		},
	)
	if err == nil {
		t.Fatal("expected invalid source signal to return an error")
	}

	projection, projectionErr :=
		ObserveLSPRevisionSelfImprovementCycleJEVTypedDecisionSignal(signal)
	if projectionErr == nil {
		t.Fatal("expected unknown source signal to remain unprojectable")
	}
	if projection.Status != "UNKNOWN" || projection.MissingStage == "" {
		t.Fatalf("expected unknown lsp projection with a missing stage, got %#v", projection)
	}
}

func TestLSPRevisionSelfImprovementCycleJEVTypedDecisionSignalRejectsTampering(
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
		t.Fatalf("source observation returned an error: %v", err)
	}

	projection, err := ObserveLSPRevisionSelfImprovementCycleJEVTypedDecisionSignal(signal)
	if err != nil {
		t.Fatalf("lsp projection returned an error: %v", err)
	}
	projection.ConfidenceMilli = 921
	if err := projection.Validate(); err == nil {
		t.Fatal("expected tampered confidence to be rejected")
	}
}
