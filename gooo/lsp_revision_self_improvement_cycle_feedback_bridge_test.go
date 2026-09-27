package gooo

import "testing"

func lspCycleFeedbackBridgeInput(t *testing.T) RevisionSelfImprovementCycleFeedbackBridgeObservation {
	t.Helper()
	cycle, metric, feedback := cycleFeedbackBridgeInputs(t)
	bridge, err := ObserveRevisionSelfImprovementCycleFeedbackBridge(cycle, metric, feedback)
	if err != nil {
		t.Fatalf("observe cycle feedback bridge: %v", err)
	}
	return bridge
}

func TestObserveLSPRevisionSelfImprovementCycleFeedbackBridgeBindsProjection(t *testing.T) {
	bridge := lspCycleFeedbackBridgeInput(t)
	got, err := ObserveLSPRevisionSelfImprovementCycleFeedbackBridge(bridge)
	if err != nil {
		t.Fatalf("observe lsp cycle feedback bridge: %v", err)
	}
	if got.Status != "BOUND" || got.MissingStage != "" {
		t.Fatalf("unexpected lsp cycle feedback bridge: %#v", got)
	}
	if got.CycleDigest != bridge.CycleDigest ||
		got.MetricDigest != bridge.MetricDigest ||
		got.FeedbackDigest != bridge.FeedbackDigest ||
		got.GeneratedIRDigest != bridge.GeneratedIRDigest ||
		got.ReverseObservationDigest != bridge.ReverseObservationDigest {
		t.Fatalf("projection lost cycle feedback evidence: %#v", got)
	}
	if got.FeedbackSignal != bridge.FeedbackSignal ||
		got.FeedbackReason != bridge.FeedbackReason ||
		got.BridgeSignal != "cycle-feedback-linked" ||
		!got.RequiresObservation || !got.SignalsAligned {
		t.Fatalf("projection lost feedback semantics: %#v", got)
	}
	if !got.ReadOnly || !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("projection safety flags = %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate lsp cycle feedback bridge: %v", err)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleFeedbackBridgePreservesUnknown(t *testing.T) {
	bridge := lspCycleFeedbackBridgeInput(t)
	bridge.FeedbackDigest = digestString("tampered")
	got, err := ObserveLSPRevisionSelfImprovementCycleFeedbackBridge(bridge)
	if err == nil {
		t.Fatal("expected tampered bridge error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "lsp-revision-self-improvement-cycle-feedback-bridge-input" {
		t.Fatalf("unexpected unknown lsp cycle feedback bridge: %#v", got)
	}
}

func TestObserveLSPRevisionSelfImprovementCycleFeedbackBridgeIsDeterministic(t *testing.T) {
	bridge := lspCycleFeedbackBridgeInput(t)
	first, err := ObserveLSPRevisionSelfImprovementCycleFeedbackBridge(bridge)
	if err != nil {
		t.Fatalf("first lsp cycle feedback bridge: %v", err)
	}
	second, err := ObserveLSPRevisionSelfImprovementCycleFeedbackBridge(bridge)
	if err != nil {
		t.Fatalf("second lsp cycle feedback bridge: %v", err)
	}
	if first != second {
		t.Fatalf("lsp cycle feedback bridge projections differ: %#v != %#v", first, second)
	}
}
