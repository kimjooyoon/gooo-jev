package gooo

import "testing"

func cycleMetricInput(t *testing.T) RevisionSelfImprovementCycleMetricInput {
	t.Helper()
	nextIteration, receipt, reverse := selfImprovementCycleInputs(t)
	cycle, err := ObserveRevisionSelfImprovementCycle(nextIteration, receipt, reverse)
	if err != nil {
		t.Fatalf("observe self-improvement cycle: %v", err)
	}
	return RevisionSelfImprovementCycleMetricInput{
		CycleDigest:              cycle.ObservationDigest,
		MetricName:               "reverse-observation-fidelity",
		SourceDigest:             cycle.SourceDigest,
		CandidateSourceDigest:    cycle.NextSourceDigest,
		GeneratedIRDigest:        cycle.NextIRDigest,
		ReverseObservationDigest: cycle.ReverseObservationDigest,
		Direction:                "higher-is-better",
		BaselineValue:            10,
		CandidateValue:           12,
	}
}

func cycleMetricCycle(t *testing.T) RevisionSelfImprovementCycleObservation {
	t.Helper()
	nextIteration, receipt, reverse := selfImprovementCycleInputs(t)
	cycle, err := ObserveRevisionSelfImprovementCycle(nextIteration, receipt, reverse)
	if err != nil {
		t.Fatalf("observe self-improvement cycle: %v", err)
	}
	return cycle
}

func TestObserveRevisionSelfImprovementCycleMetricBindsImprovement(t *testing.T) {
	cycle := cycleMetricCycle(t)
	got, err := ObserveRevisionSelfImprovementCycleMetric(cycle, cycleMetricInput(t))
	if err != nil {
		t.Fatalf("observe cycle metric: %v", err)
	}
	if got.Status != "BOUND" || got.MetricSignal != "metric-improved" {
		t.Fatalf("unexpected improved metric: %#v", got)
	}
	if got.MetricDelta != 2 || got.MetricEvidenceDigest == "" {
		t.Fatalf("unexpected metric evidence: %#v", got)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("validate cycle metric: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleMetricBindsRegressionWithDirection(t *testing.T) {
	cycle := cycleMetricCycle(t)
	input := cycleMetricInput(t)
	input.Direction = "lower-is-better"
	input.BaselineValue = 12
	input.CandidateValue = 15
	got, err := ObserveRevisionSelfImprovementCycleMetric(cycle, input)
	if err != nil {
		t.Fatalf("observe regressed cycle metric: %v", err)
	}
	if got.MetricSignal != "metric-regressed" || got.MetricDelta != 3 {
		t.Fatalf("unexpected regressed metric: %#v", got)
	}
}

func TestObserveRevisionSelfImprovementCycleMetricPreservesUnknownLink(t *testing.T) {
	cycle := cycleMetricCycle(t)
	input := cycleMetricInput(t)
	input.SourceDigest = digestString("tampered")
	got, err := ObserveRevisionSelfImprovementCycleMetric(cycle, input)
	if err == nil {
		t.Fatal("expected metric evidence link error")
	}
	if got.Status != "UNKNOWN" ||
		got.MissingStage != "revision-self-improvement-cycle-metric-evidence-link" {
		t.Fatalf("unexpected unknown metric: %#v", got)
	}
}

func TestObserveRevisionSelfImprovementCycleMetricIsDeterministic(t *testing.T) {
	cycle := cycleMetricCycle(t)
	input := cycleMetricInput(t)
	first, err := ObserveRevisionSelfImprovementCycleMetric(cycle, input)
	if err != nil {
		t.Fatalf("first cycle metric: %v", err)
	}
	second, err := ObserveRevisionSelfImprovementCycleMetric(cycle, input)
	if err != nil {
		t.Fatalf("second cycle metric: %v", err)
	}
	if first != second {
		t.Fatalf("cycle metrics differ: %#v != %#v", first, second)
	}
}