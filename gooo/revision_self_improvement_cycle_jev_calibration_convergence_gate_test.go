package gooo

import "testing"

func testJEVCalibrationConvergenceGateOutcomeDelta(
	signal string,
	coverageDigest string,
) RevisionSelfImprovementCycleJEVCalibrationOutcomeDeltaObservation {
	value := RevisionSelfImprovementCycleJEVCalibrationOutcomeDeltaObservation{
		Status:                           "BOUND",
		BaselineAccuracyMilli:            600,
		CurrentAccuracyMilli:             700,
		AccuracyDeltaMilli:              100,
		BaselineMeanConfidenceMilli:       650,
		CurrentMeanConfidenceMilli:        700,
		ConfidenceDeltaMilli:             50,
		BaselineOutcomeSamplesDigest:     digestString("baseline-outcomes"),
		CurrentOutcomeSamplesDigest:      digestString("current-outcomes"),
		SourceObservationDigest:          digestString("source-observation"),
		BaselineGenerationTraceDigest:    digestString("baseline-trace"),
		CurrentGenerationTraceDigest:     digestString("current-trace"),
		ReverseObservationCoverageDigest: coverageDigest,
		DeltaSignal:                      signal,
		ReadOnly:                         true,
		NonExecuting:                     true,
		NonAuthorizing:                   true,
	}
	value.ObservationDigest = digestRevisionSelfImprovementCycleJEVCalibrationOutcomeDelta(value)
	return value
}

func testJEVCalibrationConvergenceGateCoverage() RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetric {
	value := RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetric{
		Status:                   "BOUND",
		MetricName:               jevGenerationTraceCoverageMetricName,
		RequiredDigestCount:      6,
		LinkedDigestCount:        6,
		CoverageMilli:             1000,
		CoverageBand:              "high",
		MetricSignal:              jevGenerationTraceCoverageComplete,
		MetricDigest:              digestString("coverage-metric"),
		ReverseObservationDigest: digestString("reverse-observation"),
		NonExecuting:             true,
		NonAuthorizing:           true,
	}
	value.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricDigest(value)
	return value
}

func testJEVCalibrationConvergenceGateInput(
	deltaSignal string,
	choiceStatus string,
	choiceSignal string,
	choiceDecision string,
) RevisionSelfImprovementCycleJEVCalibrationConvergenceGateInput {
	coverage := testJEVCalibrationConvergenceGateCoverage()
	return RevisionSelfImprovementCycleJEVCalibrationConvergenceGateInput{
		OutcomeDelta:                           testJEVCalibrationConvergenceGateOutcomeDelta(deltaSignal, coverage.ObservationDigest),
		ChoiceSetSensitivityStatus:             choiceStatus,
		ChoiceSetSensitivitySignal:             choiceSignal,
		ChoiceSetDecisionSignal:                choiceDecision,
		ChoiceSetSensitivityObservationDigest:  digestString("choice-sensitivity"),
		Coverage:                                coverage,
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateConverged(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGate(
		testJEVCalibrationConvergenceGateInput(
			jevCalibrationOutcomeDeltaImproved,
			"BOUND",
			jevCalibrationConvergenceGateChoiceStable,
			jevCalibrationConvergenceGateChoiceCompare,
		),
	)
	if observation.Status != "BOUND" ||
		observation.GateDecision != jevCalibrationConvergenceGateConverged {
		t.Fatalf("expected converged gate, got %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("expected valid converged gate: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateRegressed(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGate(
		testJEVCalibrationConvergenceGateInput(
			jevCalibrationOutcomeDeltaRegressed,
			"BOUND",
			jevCalibrationConvergenceGateChoiceStable,
			jevCalibrationConvergenceGateChoiceCompare,
		),
	)
	if observation.GateDecision != jevCalibrationConvergenceGateRegressed {
		t.Fatalf("expected regressed gate, got %#v", observation)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGateDefersChangedChoiceSet(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGate(
		testJEVCalibrationConvergenceGateInput(
			jevCalibrationOutcomeDeltaImproved,
			"BOUND",
			jevCalibrationConvergenceGateChoiceChanged,
			jevCalibrationConvergenceGateChoiceDefer,
		),
	)
	if observation.Status != "BOUND" ||
		observation.GateDecision != jevCalibrationConvergenceGateDefer {
		t.Fatalf("changed choice evidence must defer, got %#v", observation)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGatePreservesUnknown(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGate(
		testJEVCalibrationConvergenceGateInput(
			jevCalibrationOutcomeDeltaImproved,
			"UNKNOWN",
			jevCalibrationConvergenceGateChoiceUnknown,
			jevCalibrationConvergenceGateChoiceDefer,
		),
	)
	if observation.Status != "UNKNOWN" || observation.GateDecision != jevCalibrationConvergenceGateUnknown {
		t.Fatalf("incomplete choice evidence must remain unknown, got %#v", observation)
	}
	if err := observation.Validate(); err != nil {
		t.Fatalf("expected valid unknown gate: %v", err)
	}
}

func TestRevisionSelfImprovementCycleJEVCalibrationConvergenceGateRejectsTampering(t *testing.T) {
	observation := ObserveRevisionSelfImprovementCycleJEVCalibrationConvergenceGate(
		testJEVCalibrationConvergenceGateInput(
			jevCalibrationOutcomeDeltaImproved,
			"BOUND",
			jevCalibrationConvergenceGateChoiceStable,
			jevCalibrationConvergenceGateChoiceCompare,
		),
	)
	observation.GateDecision = jevCalibrationConvergenceGateRegressed
	if err := observation.Validate(); err == nil {
		t.Fatal("expected gate decision tampering to be rejected")
	}
}