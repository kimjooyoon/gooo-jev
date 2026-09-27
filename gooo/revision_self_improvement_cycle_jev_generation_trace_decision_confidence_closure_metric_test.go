package gooo

import "testing"

func generationTraceDecisionConfidenceClosureInput() RevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetricInput {
	generation := RevisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetric{
		Status:                  "BOUND",
		MetricName:              jevGenerationTraceCoverageMetricName,
		RequiredDigestCount:     6,
		LinkedDigestCount:       6,
		CoverageMilli:           1000,
		CoverageBand:            "high",
		MetricSignal:            jevGenerationTraceCoverageComplete,
		MetricDigest:             digestString("generation-metric"),
		ReverseObservationDigest: digestString("generation-reverse"),
		NonExecuting:             true,
		NonAuthorizing:           true,
	}
	generation.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricDigest(generation)

	confidence := RevisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetric{
		Status:                   "BOUND",
		MetricName:               jevDecisionConfidenceMetricName,
		MetricDigest:             digestString("confidence-metric"),
		ReverseObservationDigest: digestString("confidence-reverse"),
		ClosureSignal:             jevDecisionConfidenceReverseClosureComplete,
		NonExecuting:             true,
		NonAuthorizing:           true,
	}
	confidence.ObservationDigest = revisionSelfImprovementCycleJEVDecisionConfidenceMetricReverseObservationProvenanceClosureMetricDigest(confidence)

	return RevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetricInput{
		GenerationTraceMetric:   generation,
		ConfidenceClosureMetric: confidence,
	}
}

func TestObserveRevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetricBinds(t *testing.T) {
	closure := ObserveRevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetric(
		generationTraceDecisionConfidenceClosureInput(),
	)
	if closure.Status != "BOUND" ||
		closure.LinkedMetricCount != 2 ||
		closure.MissingStage != "" ||
		closure.GenerationTraceMetricDigest == "" ||
		closure.ConfidenceClosureMetricDigest == "" {
		t.Fatalf("unexpected bound closure: %#v", closure)
	}
	if err := closure.Validate(); err != nil {
		t.Fatalf("validate closure: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetricPreservesMissingGenerationStage(t *testing.T) {
	input := generationTraceDecisionConfidenceClosureInput()
	input.GenerationTraceMetric.Status = "UNKNOWN"
	input.GenerationTraceMetric.MissingStage = "generated_artifact_digest"
	input.GenerationTraceMetric.LinkedDigestCount = 0
	input.GenerationTraceMetric.CoverageMilli = 0
	input.GenerationTraceMetric.CoverageBand = "unknown"
	input.GenerationTraceMetric.MetricSignal = jevGenerationTraceCoverageUnknown
	input.GenerationTraceMetric.MetricDigest = ""
	input.GenerationTraceMetric.ReverseObservationDigest = ""
	input.GenerationTraceMetric.ObservationDigest = revisionSelfImprovementCycleJEVGenerationTraceReverseObservationCoverageMetricDigest(input.GenerationTraceMetric)

	closure := ObserveRevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetric(input)
	if closure.Status != "UNKNOWN" || closure.MissingStage != "generated_artifact_digest" || closure.LinkedMetricCount != 0 {
		t.Fatalf("unexpected missing generation stage closure: %#v", closure)
	}
	if err := closure.Validate(); err != nil {
		t.Fatalf("validate unknown closure: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetricRejectsTamperedConfidenceMetric(t *testing.T) {
	input := generationTraceDecisionConfidenceClosureInput()
	input.ConfidenceClosureMetric.MetricDigest = digestString("tampered")
	closure := ObserveRevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetric(input)

	if closure.Status != "UNKNOWN" || closure.MissingStage != "decision_confidence_closure_metric" {
		t.Fatalf("tampered confidence metric was not rejected: %#v", closure)
	}
	if err := closure.Validate(); err != nil {
		t.Fatalf("validate tampered closure: %v", err)
	}
}

func TestObserveRevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetricFailsClosedOnCapabilityBoundary(t *testing.T) {
	input := generationTraceDecisionConfidenceClosureInput()
	input.ConfidenceClosureMetric.NonAuthorizing = false
	closure := ObserveRevisionSelfImprovementCycleJEVGenerationTraceDecisionConfidenceClosureMetric(input)

	if closure.Status != "UNKNOWN" || closure.MissingStage != "capability_boundary" || closure.NonAuthorizing {
		t.Fatalf("capability boundary was not preserved: %#v", closure)
	}
}