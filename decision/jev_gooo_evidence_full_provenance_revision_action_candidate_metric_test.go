package decision

import "testing"

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetric(t *testing.T) {
	projection := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPInput{
		Verification:        verifiedActionDerivedCandidateForLSP(t),
		EvidencePrefixDigest: "candidate-metric-clear",
		NonAuthorizing:      true,
	})
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricInput{
		Projection:      projection,
		MetricName:      "jev-action-derived-candidate-verification",
		NonAuthorizing: true,
	})
	if output.Status != "measured" || output.Total != 1 || output.VerifiedCount != 1 ||
		output.CandidateID == "" || output.MetricEvidenceDigest == "" {
		t.Fatalf("output = %#v, want measured verified candidate metric", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricRetainsCounterexample(t *testing.T) {
	projection := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPInput{
		Verification:        reviewActionDerivedCandidateForLSP(t),
		MissingStageIndex:   13,
		EvidencePrefixDigest: "candidate-metric-counterexample",
		NonAuthorizing:      true,
	})
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricInput{
		Projection:      projection,
		MetricName:      "jev-action-derived-candidate-verification",
		NonAuthorizing: true,
	})
	if output.Status != "measured-with-counterexample" || output.Total != 1 ||
		output.CounterexampleCount != 1 || output.UnknownCount != 0 ||
		output.CandidateEvidenceDigest == "" {
		t.Fatalf("output = %#v, want counterexample candidate metric", output)
	}
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricRetainsUnknown(t *testing.T) {
	projection := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPInput{
		Verification:        verifiedActionDerivedCandidateForLSP(t),
		EvidencePrefixDigest: "candidate-metric-unknown",
		NonAuthorizing:      true,
	})
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricInput{
		Projection:             projection,
		MetricName:             "jev-action-derived-candidate-verification",
		AdditionalUnknownCount: 2,
		NonAuthorizing:         true,
	})
	if output.Status != "measured-with-unknown" || output.Total != 3 ||
		output.VerifiedCount != 1 || output.UnknownCount != 2 {
		t.Fatalf("output = %#v, want unknown-retaining candidate metric", output)
	}
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricFailsClosed(t *testing.T) {
	projection := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPInput{
		Verification:        verifiedActionDerivedCandidateForLSP(t),
		EvidencePrefixDigest: "candidate-metric-failure",
		NonAuthorizing:      true,
	})
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricInput{
		Projection:      projection,
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "metric-name" {
		t.Fatalf("output = %#v, want metric-name UNKNOWN", output)
	}
}