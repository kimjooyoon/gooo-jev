package decision

import "testing"

func evidenceFullProvenanceMetricInput(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceMetricInput {
	t.Helper()
	base := evidenceFullProvenanceInput(t)
	return ExecutionEnvelopeGoooEvidenceFullProvenanceMetricInput{
		Binding:                BindExecutionEnvelopeGoooEvidenceFullProvenance(base),
		ObservedStatus:         base.ObservedStatus,
		ExpectedStatus:         base.ExpectedStatus,
		ObservedMissingStage:   base.ObservedMissingStage,
		ExpectedMissingStage:   base.ExpectedMissingStage,
		ObservedEvidenceDigest: base.ObservedEvidenceDigest,
		ExpectedEvidenceDigest: base.ExpectedEvidenceDigest,
		MetricSource:           base.MetricSource,
		NonAuthorizing:         true,
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceMetricConfirmed(t *testing.T) {
	metric := BindExecutionEnvelopeGoooEvidenceFullProvenanceMetric(evidenceFullProvenanceMetricInput(t))
	if metric.Status != "ready" || metric.Classification != "CONFIRMED" ||
		metric.ConfirmedCount != 1 || metric.ConfirmedPermille != 1000 ||
		metric.UnknownCount != 0 || metric.EvidenceDigest == "" {
		t.Fatalf("metric = %#v, want confirmed ready metric", metric)
	}
	if err := metric.Validate(); err != nil {
		t.Fatalf("metric should validate: %v", err)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceMetricRefuted(t *testing.T) {
	input := evidenceFullProvenanceMetricInput(t)
	input.ObservedEvidenceDigest = "different-observation"
	metric := BindExecutionEnvelopeGoooEvidenceFullProvenanceMetric(input)
	if metric.Status != "ready" || metric.Classification != "REFUTED" ||
		metric.RefutedCount != 1 || metric.RefutedPermille != 1000 {
		t.Fatalf("metric = %#v, want refuted ready metric", metric)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceMetricUnknown(t *testing.T) {
	input := evidenceFullProvenanceMetricInput(t)
	input.MetricSource = ""
	metric := BindExecutionEnvelopeGoooEvidenceFullProvenanceMetric(input)
	if metric.Status != "UNKNOWN" || metric.MissingStage != "metric-source" ||
		metric.UnknownCount != 1 || metric.UnknownPermille != 1000 ||
		metric.EvidenceDigest != "" {
		t.Fatalf("metric = %#v, want metric-source UNKNOWN", metric)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceMetricRejectsTampering(t *testing.T) {
	input := evidenceFullProvenanceMetricInput(t)
	input.Binding.EvidenceDigest = "tampered"
	metric := BindExecutionEnvelopeGoooEvidenceFullProvenanceMetric(input)
	if metric.Status != "UNKNOWN" || metric.MissingStage != "binding-validation" {
		t.Fatalf("metric = %#v, want binding-validation UNKNOWN", metric)
	}
}
