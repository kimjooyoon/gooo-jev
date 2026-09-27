package gooo

import "testing"

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageMetricBindsFourEvidenceDigests(t *testing.T) {
	cycle, metric, feedback := cycleFeedbackBridgeInputs(t)
	bridge, err := ObserveRevisionSelfImprovementCycleFeedbackBridge(cycle, metric, feedback)
	if err != nil {
		t.Fatalf("observe cycle feedback bridge: %v", err)
	}
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageMetric(bridge)
	if err != nil {
		t.Fatalf("observe cycle evidence coverage metric: %v", err)
	}
	if got.Status != "BOUND" || got.MetricName != "source-ir-generation-reverse-observation-coverage" || got.CoverageSignal != "evidence-complete" {
		t.Fatalf("unexpected evidence coverage metric: %#v", got)
	}
	if got.LinkedEvidenceCount != 4 || got.RequiredEvidenceCount != 4 {
		t.Fatalf("unexpected evidence coverage counts: %#v", got)
	}
	if got.SourceDigest != bridge.SourceDigest || got.CandidateSourceDigest != bridge.CandidateSourceDigest || got.GeneratedIRDigest != bridge.GeneratedIRDigest || got.ReverseObservationDigest != bridge.ReverseObservationDigest {
		t.Fatalf("metric lost source/ir/reverse evidence: %#v", got)
	}
	if !got.NonExecuting || !got.NonAuthorizing {
		t.Fatalf("metric safety flags = %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate evidence coverage metric: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageMetricPreservesUnknown(t *testing.T) {
	cycle, metric, feedback := cycleFeedbackBridgeInputs(t)
	bridge, err := ObserveRevisionSelfImprovementCycleFeedbackBridge(cycle, metric, feedback)
	if err != nil {
		t.Fatalf("observe cycle feedback bridge: %v", err)
	}
	bridge.GeneratedIRDigest = digestString("tampered")
	got, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageMetric(bridge)
	if err == nil {
		t.Fatal("expected tampered bridge error")
	}
	if got.Status != "UNKNOWN" || got.MissingStage != "revision-self-improvement-cycle-evidence-coverage-metric-bridge" {
		t.Fatalf("unexpected unknown evidence coverage metric: %#v", got)
	}
}

func TestObserveRevisionSelfImprovementCycleEvidenceCoverageMetricIsDeterministic(t *testing.T) {
	cycle, metric, feedback := cycleFeedbackBridgeInputs(t)
	bridge, err := ObserveRevisionSelfImprovementCycleFeedbackBridge(cycle, metric, feedback)
	if err != nil {
		t.Fatalf("observe cycle feedback bridge: %v", err)
	}
	first, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageMetric(bridge)
	if err != nil {
		t.Fatalf("first evidence coverage metric: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCycleEvidenceCoverageMetric(bridge)
	if err != nil {
		t.Fatalf("second evidence coverage metric: %v", err)
	}
	if first != second {
		t.Fatalf("evidence coverage metrics differ: %#v != %#v", first, second)
	}
}
