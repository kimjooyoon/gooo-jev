package gooo

import "testing"

func declarationEvidenceCycleMetricBoundProjection() ExecutionEnvelopeDeclarationEvidenceCycleProjection {
	return ProjectExecutionEnvelopeDeclarationEvidenceCycleLSP(ExecutionEnvelopeDeclarationEvidenceCycleInput{
		SourceStatus:              GoooDeclarationSourceDerived,
		DeclarationID:             "decl://example",
		ContractID:                "contract://example",
		DeclarationDigest:         "1111111111111111111111111111111111111111111111111111111111111111",
		IRStatus:                  ExecutionEnvelopeDeclarationIRGenerationLSPBound,
		IRDigest:                  "2222222222222222222222222222222222222222222222222222222222222222",
		GenerationDigest:          "3333333333333333333333333333333333333333333333333333333333333333",
		BindingDigest:             "4444444444444444444444444444444444444444444444444444444444444444",
		ReverseStatus:             "BOUND",
		ReverseObservationDigest: "5555555555555555555555555555555555555555555555555555555555555555",
		NonExecuting:              true,
		NonAuthorizing:            true,
	})
}

func TestMeasureDeclarationEvidenceCycleMetricBound(t *testing.T) {
	metric := MeasureExecutionEnvelopeDeclarationEvidenceCycleMetric(declarationEvidenceCycleMetricBoundProjection())
	if metric.Status != ExecutionEnvelopeDeclarationEvidenceCycleMetricBound ||
		metric.StageCount != 3 || metric.StageTotal != 3 || metric.Coverage != 1 || metric.FirstMissingStage != "" {
		t.Fatalf("metric = %#v", metric)
	}
	if err := metric.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestMeasureDeclarationEvidenceCycleMetricPreservesUnknownStage(t *testing.T) {
	projection := ProjectExecutionEnvelopeDeclarationEvidenceCycleLSP(ExecutionEnvelopeDeclarationEvidenceCycleInput{
		SourceStatus:      GoooDeclarationSourceDerived,
		DeclarationID:     "decl://example",
		ContractID:        "contract://example",
		DeclarationDigest: "1111111111111111111111111111111111111111111111111111111111111111",
		IRStatus:          ExecutionEnvelopeDeclarationIRGenerationLSPUnknown,
		IRDigest:          "2222222222222222222222222222222222222222222222222222222222222222",
		GenerationDigest:  "3333333333333333333333333333333333333333333333333333333333333333",
		BindingDigest:     "4444444444444444444444444444444444444444444444444444444444444444",
		ReverseStatus:     "UNKNOWN",
		NonExecuting:      true,
		NonAuthorizing:    true,
	})
	metric := MeasureExecutionEnvelopeDeclarationEvidenceCycleMetric(projection)
	if metric.Status != ExecutionEnvelopeDeclarationEvidenceCycleMetricUnknown ||
		metric.FirstMissingStage != "declaration_ir_generation" || metric.StageCount != 1 {
		t.Fatalf("metric = %#v", metric)
	}
	if err := metric.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestMeasureDeclarationEvidenceCycleMetricPreservesDeferredReverse(t *testing.T) {
	projection := declarationEvidenceCycleMetricBoundProjection()
	projection.Status = ExecutionEnvelopeDeclarationEvidenceCycleDeferred
	projection.Code = executionEnvelopeDeclarationEvidenceCycleDeferredCode
	projection.MissingStage = "reverse_observation"
	projection.FirstMismatch = "producer-deferred"
	projection.Message = "declaration evidence cycle is deferred at reverse observation"
	projection.EvidenceDigest = executionEnvelopeDeclarationEvidenceCycleDigest(projection)
	metric := MeasureExecutionEnvelopeDeclarationEvidenceCycleMetric(projection)
	if metric.Status != ExecutionEnvelopeDeclarationEvidenceCycleMetricDeferred || metric.StageCount != 2 {
		t.Fatalf("metric = %#v", metric)
	}
	if err := metric.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestMeasureDeclarationEvidenceCycleMetricRejectsProjectionTampering(t *testing.T) {
	projection := declarationEvidenceCycleMetricBoundProjection()
	projection.EvidenceDigest = "tampered"
	metric := MeasureExecutionEnvelopeDeclarationEvidenceCycleMetric(projection)
	if metric.Status != ExecutionEnvelopeDeclarationEvidenceCycleMetricError || metric.FirstMissingStage != "projection_integrity" {
		t.Fatalf("metric = %#v", metric)
	}
	if err := metric.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}