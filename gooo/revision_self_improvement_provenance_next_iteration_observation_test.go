package gooo

import "testing"

func provenanceNextIterationInputs(t *testing.T) (RevisionSelfImprovementProvenanceReplanObservation, RevisionSelfImprovementIteration) {
	t.Helper()
	decision, plan := provenanceReplanInputs(t)
	replan, err := ObserveRevisionSelfImprovementProvenanceReplan(decision, plan)
	if err != nil { t.Fatalf("ObserveRevisionSelfImprovementProvenanceReplan() error = %v", err) }
	source, history, feedback := selfImprovementIterationInputs(t)
	nextIteration, err := ObserveRevisionSelfImprovementIteration(source, history, feedback)
	if err != nil { t.Fatalf("ObserveRevisionSelfImprovementIteration() error = %v", err) }
	return replan, nextIteration
}

func TestObserveRevisionSelfImprovementProvenanceNextIterationBindsBoundary(t *testing.T) {
	replan, nextIteration := provenanceNextIterationInputs(t)
	result, err := ObserveRevisionSelfImprovementProvenanceNextIteration(replan, nextIteration)
	if err != nil { t.Fatalf("ObserveRevisionSelfImprovementProvenanceNextIteration() error = %v", err) }
	if result.Status != "BOUND" || result.TransitionSignal != "provenance-next-iteration-mismatch" || result.NextIterationDigest != nextIteration.IterationDigest || result.FeedbackSignal != nextIteration.FeedbackSignal { t.Fatalf("unexpected provenance next iteration: %#v", result) }
	if result.ReplanDigest != replan.ObservationDigest || result.SourceDigest != nextIteration.SourceDigest { t.Fatalf("next iteration lost replan links: %#v", result) }
	if err := result.Validate(); err != nil { t.Fatalf("Validate() error = %v", err) }
}

func TestObserveRevisionSelfImprovementProvenanceNextIterationRetainsReplanFailure(t *testing.T) {
	replan, nextIteration := provenanceNextIterationInputs(t)
	replan.ObservationDigest = digestString("tampered")
	result, err := ObserveRevisionSelfImprovementProvenanceNextIteration(replan, nextIteration)
	if err == nil { t.Fatal("ObserveRevisionSelfImprovementProvenanceNextIteration() error = nil, want replan failure") }
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-provenance-next-iteration-replan" { t.Fatalf("unexpected unknown replan boundary: %#v", result) }
}

func TestObserveRevisionSelfImprovementProvenanceNextIterationRetainsSourceFailure(t *testing.T) {
	replan, nextIteration := provenanceNextIterationInputs(t)
	nextIteration.SourceDigest = digestString("other-source")
	nextIteration.IterationDigest = digestRevisionSelfImprovementIteration(nextIteration)
	result, err := ObserveRevisionSelfImprovementProvenanceNextIteration(replan, nextIteration)
	if err == nil { t.Fatal("ObserveRevisionSelfImprovementProvenanceNextIteration() error = nil, want source link failure") }
	if result.Status != "UNKNOWN" || result.MissingStage != "revision-self-improvement-provenance-next-iteration-source-link" { t.Fatalf("unexpected unknown source boundary: %#v", result) }
}

func TestObserveRevisionSelfImprovementProvenanceNextIterationIsDeterministic(t *testing.T) {
	replan, nextIteration := provenanceNextIterationInputs(t)
	first, err := ObserveRevisionSelfImprovementProvenanceNextIteration(replan, nextIteration)
	if err != nil { t.Fatalf("first ObserveRevisionSelfImprovementProvenanceNextIteration() error = %v", err) }
	second, err := ObserveRevisionSelfImprovementProvenanceNextIteration(replan, nextIteration)
	if err != nil { t.Fatalf("second ObserveRevisionSelfImprovementProvenanceNextIteration() error = %v", err) }
	if first.ObservationDigest != second.ObservationDigest { t.Fatal("same replan and next iteration produced different boundary digest") }
}
