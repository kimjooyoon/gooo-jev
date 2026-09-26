package decision

import "testing"

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetric(t *testing.T) {
	projection := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPInput{
		Verification:        verifiedRevisionActionForLSP(t),
		EvidencePrefixDigest: "prefix-metric-clear",
		NonAuthorizing:      true,
	})
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricInput{
		Projection:      projection,
		MetricName:      "jev-revision-action-verification",
		NonAuthorizing: true,
	})
	if output.Status != "measured" || output.Total != 1 || output.VerifiedCount != 1 ||
		output.MetricEvidenceDigest == "" {
		t.Fatalf("output = %#v, want measured verified metric", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricRetainsCounterexample(t *testing.T) {
	projection := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPInput{
		Verification:        reviewRevisionActionForLSP(t),
		MissingStageIndex:   7,
		EvidencePrefixDigest: "prefix-metric-counterexample",
		NonAuthorizing:      true,
	})
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricInput{
		Projection:      projection,
		MetricName:      "jev-revision-action-verification",
		NonAuthorizing: true,
	})
	if output.Status != "measured-with-counterexample" || output.Total != 1 ||
		output.CounterexampleCount != 1 || output.UnknownCount != 0 {
		t.Fatalf("output = %#v, want counterexample metric", output)
	}
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricRetainsUnknown(t *testing.T) {
	projection := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPInput{
		Verification:        verifiedRevisionActionForLSP(t),
		EvidencePrefixDigest: "prefix-metric-unknown",
		NonAuthorizing:      true,
	})
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricInput{
		Projection:             projection,
		MetricName:             "jev-revision-action-verification",
		AdditionalUnknownCount: 2,
		NonAuthorizing:         true,
	})
	if output.Status != "measured-with-unknown" || output.Total != 3 ||
		output.VerifiedCount != 1 || output.UnknownCount != 2 {
		t.Fatalf("output = %#v, want unknown-retaining metric", output)
	}
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricFailsClosed(t *testing.T) {
	projection := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPInput{
		Verification:        verifiedRevisionActionForLSP(t),
		EvidencePrefixDigest: "prefix-metric-failure",
		NonAuthorizing:      true,
	})
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricInput{
		Projection:      projection,
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "metric-name" {
		t.Fatalf("output = %#v, want metric-name UNKNOWN", output)
	}
}