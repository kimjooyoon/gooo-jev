package gooo

import "testing"

func cycleFeedbackDecisionBoundaryInputs(t *testing.T) (
	RevisionSelfImprovementCycleObservation,
	RevisionSelfImprovementCycleFeedbackBridgeObservation,
	RevisionSelfImprovementCycleFeedbackDecisionLinkObservation,
) {
	t.Helper()
	cycle := cycleMetricCycle(t)
	bridge, decision, questionDigest := cycleFeedbackDecisionLinkInputs(t)
	link, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionLink(bridge, decision, questionDigest)
	if err != nil {
		t.Fatalf("observe cycle feedback decision link: %v", err)
	}
	return cycle, bridge, link
}

func TestObserveRevisionSelfImprovementCycleFeedbackDecisionBoundaryBindsCycle(t *testing.T) {
	cycle, bridge, link := cycleFeedbackDecisionBoundaryInputs(t)
	got, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionBoundary(cycle, bridge, link)
	if err != nil {
		t.Fatalf("observe cycle feedback decision boundary: %v", err)
	}
	if got.Status != "BOUND" || got.BoundarySignal != "cycle-feedback-decision-boundary-linked" || !got.BoundaryAligned {
		t.Fatalf("unexpected cycle feedback decision boundary: %#v", got)
	}
	if got.CycleDigest != cycle.ObservationDigest || got.BridgeDigest != bridge.ObservationDigest || got.DecisionLinkDigest != link.ObservationDigest {
		t.Fatalf("boundary lost provenance links: %#v", got)
	}
	if got.DecisionSignal != revisionSelfImprovementDecisionReceiptRequireReplan || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("boundary lost decision semantics: %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate cycle feedback decision boundary: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleFeedbackDecisionBoundaryPreservesUnknown(t *testing.T) {
	cycle, bridge, link := cycleFeedbackDecisionBoundaryInputs(t)
	bridge.CycleDigest = digestString("tampered")
	got, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionBoundary(cycle, bridge, link)
	if err == nil {
		t.Fatal("expected tampered bridge error")
	}
	if got.Status != "UNKNOWN" || got.MissingStage != "revision-self-improvement-cycle-feedback-decision-boundary-bridge" {
		t.Fatalf("unexpected unknown cycle feedback decision boundary: %#v", got)
	}
}

func TestObserveRevisionSelfImprovementCycleFeedbackDecisionBoundaryIsDeterministic(t *testing.T) {
	cycle, bridge, link := cycleFeedbackDecisionBoundaryInputs(t)
	first, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionBoundary(cycle, bridge, link)
	if err != nil {
		t.Fatalf("first cycle feedback decision boundary: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCycleFeedbackDecisionBoundary(cycle, bridge, link)
	if err != nil {
		t.Fatalf("second cycle feedback decision boundary: %v", err)
	}
	if first != second {
		t.Fatalf("cycle feedback decision boundaries differ: %#v != %#v", first, second)
	}
}
