package decision

import "testing"

func TestMeasureExecutionEnvelopeProvenanceChainMetricFromSource(t *testing.T) {
	declaration := BindExecutionEnvelopeDeclarationIRGeneration(
		"decl-1", "contract-1", "decl-digest", "ir-digest", "generation-digest",
	)
	reverse := BindExecutionEnvelopeReverseObservation(ObserveExecutionEnvelopeProvenanceReverse(ExecutionEnvelopeReverseObservationInput{
		ObservedStatus:         "ready",
		ExpectedStatus:         "ready",
		ObservedEvidenceDigest: "evidence-1",
		ExpectedEvidenceDigest: "evidence-1",
		NonAuthorizing:         true,
	}))
	chain := EvaluateExecutionEnvelopeProvenanceChainBinding(ExecutionEnvelopeProvenanceChainBindingInput{
		DeclarationIRGeneration:  declaration,
		ReverseObservationDigest: reverse.ObservationDigest,
		MetricDigest:             "legacy-metric",
		NonAuthorizing:           true,
	})
	if chain.Status != "ready" {
		t.Fatalf("chain status = %q, want ready", chain.Status)
	}
	input := ExecutionEnvelopeProvenanceChainMetricSourceInput{
		Binding:        chain,
		SourceText:     "metric: all five provenance stages observed",
		NonAuthorizing: true,
	}
	metric := MeasureExecutionEnvelopeProvenanceChainMetricFromSource(input)
	if metric.Status != "complete" || metric.ObservedStageCount != 5 {
		t.Fatalf("metric = %#v, want complete five-stage metric", metric)
	}
	if metric.SourceDigest == "" || metric.EvidenceDigest == "" || metric.CompletenessDigest == "" {
		t.Fatal("complete metric must retain source, evidence, and completeness digests")
	}
	changed := input
	changed.SourceText += "\n"
	changedMetric := MeasureExecutionEnvelopeProvenanceChainMetricFromSource(changed)
	if changedMetric.SourceDigest == metric.SourceDigest || changedMetric.EvidenceDigest == metric.EvidenceDigest {
		t.Fatal("changed metric source must change source and evidence digests")
	}
	missingSource := input
	missingSource.SourceText = ""
	missingMetric := MeasureExecutionEnvelopeProvenanceChainMetricFromSource(missingSource)
	if missingMetric.Status != "UNKNOWN" || missingMetric.MissingStage != "metric-source" {
		t.Fatalf("missing metric source = %#v, want UNKNOWN at metric-source", missingMetric)
	}
	unauthorized := input
	unauthorized.NonAuthorizing = false
	unauthorizedMetric := MeasureExecutionEnvelopeProvenanceChainMetricFromSource(unauthorized)
	if unauthorizedMetric.Status != "UNKNOWN" || unauthorizedMetric.NonAuthorizing || unauthorizedMetric.MissingStage != "authorization-boundary" {
		t.Fatalf("unauthorized metric source = %#v, want non-authorizing UNKNOWN", unauthorizedMetric)
	}
}
