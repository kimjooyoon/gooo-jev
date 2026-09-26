package decision

import "testing"

func clearActionDerivedCandidateRevisionLSPForMetric(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPBinding {
	t.Helper()
	return ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPInput{
		Verification: VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionVerificationInput{
			Guard: verifiedActionDerivedCandidateRevisionForVerification(t),
			ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
				Status:         "reproduced",
				EvidenceDigest: "candidate-revision-metric-reproduced",
				NonAuthorizing: true,
			},
			NonAuthorizing: true,
		}),
		EvidencePrefixDigest: "candidate-revision-metric-clear",
		NonAuthorizing:      true,
	})
}

func counterexampleActionDerivedCandidateRevisionLSPForMetric(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPBinding {
	t.Helper()
	return ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionLSPInput{
		Verification:        reviewActionDerivedCandidateRevisionForLSP(t),
		MissingStageIndex:   17,
		EvidencePrefixDigest: "candidate-revision-metric-counterexample",
		NonAuthorizing:      true,
	})
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetric(t *testing.T) {
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricInput{
		Projection:      clearActionDerivedCandidateRevisionLSPForMetric(t),
		MetricName:      "jev-action-derived-candidate-revision-verification",
		NonAuthorizing: true,
	})
	if output.Status != "measured" || output.Total != 1 || output.VerifiedCount != 1 ||
		output.SourceCandidateID == "" || output.RevisionCandidateDigest == "" {
		t.Fatalf("output = %#v, want measured verified candidate revision metric", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricRetainsCounterexample(t *testing.T) {
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricInput{
		Projection:      counterexampleActionDerivedCandidateRevisionLSPForMetric(t),
		MetricName:      "jev-action-derived-candidate-revision-verification",
		NonAuthorizing: true,
	})
	if output.Status != "measured-with-counterexample" || output.Total != 1 ||
		output.CounterexampleCount != 1 || output.UnknownCount != 0 {
		t.Fatalf("output = %#v, want counterexample candidate revision metric", output)
	}
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricRetainsUnknown(t *testing.T) {
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricInput{
		Projection:             clearActionDerivedCandidateRevisionLSPForMetric(t),
		MetricName:             "jev-action-derived-candidate-revision-verification",
		AdditionalUnknownCount: 2,
		NonAuthorizing:         true,
	})
	if output.Status != "measured-with-unknown" || output.Total != 3 ||
		output.VerifiedCount != 1 || output.UnknownCount != 2 {
		t.Fatalf("output = %#v, want unknown-retaining candidate revision metric", output)
	}
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricFailsClosed(t *testing.T) {
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricInput{
		Projection:      clearActionDerivedCandidateRevisionLSPForMetric(t),
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "metric-name" {
		t.Fatalf("output = %#v, want metric-name UNKNOWN", output)
	}
}