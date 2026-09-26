package decision

import "testing"

func clearExtendedLineageCandidateRevisionForMetric(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionLSPBinding {
	t.Helper()
	return ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionLSPInput{
		Verification: VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionVerificationInput{
			Guard: admittedExtendedLineageCandidateRevisionForRevisionVerification(t),
			ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
				Status:         "reproduced",
				EvidenceDigest: "extended-lineage-candidate-revision-metric-reproduced",
				NonAuthorizing: true,
			},
			NonAuthorizing: true,
		}),
		EvidencePrefixDigest: "extended-lineage-candidate-revision-metric-clear",
		MissingStageIndex:   -1,
		NonAuthorizing:      true,
	})
}

func counterexampleExtendedLineageCandidateRevisionForMetric(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionLSPBinding {
	t.Helper()
	return ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionLSPInput{
		Verification: VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevision(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionVerificationInput{
			Guard: admittedExtendedLineageCandidateRevisionForRevisionVerification(t),
			ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
				Status:         "counterexample",
				EvidenceDigest: "extended-lineage-candidate-revision-metric-counterexample",
				FirstMismatch:  "extended-lineage-output",
				NonAuthorizing: true,
			},
			NonAuthorizing: true,
		}),
		MissingStageIndex:   28,
		EvidencePrefixDigest: "extended-lineage-candidate-revision-metric-counterexample",
		NonAuthorizing:      true,
	})
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionMetric(t *testing.T) {
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionMetricInput{
		Projection:      clearExtendedLineageCandidateRevisionForMetric(t),
		MetricName:      "jev-generated-extended-lineage-candidate-revision-verification",
		NonAuthorizing: true,
	})
	if output.Status != "measured" || output.Total != 1 || output.VerifiedCount != 1 ||
		output.SourceCandidateID == "" || output.ParentCandidateDigest == "" ||
		output.CandidateDigest == "" || output.GuardEvidenceDigest == "" ||
		output.ReverseEvidenceDigest == "" {
		t.Fatalf("output = %#v, want measured extended lineage candidate revision metric", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionMetricRetainsCounterexample(t *testing.T) {
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionMetricInput{
		Projection:      counterexampleExtendedLineageCandidateRevisionForMetric(t),
		MetricName:      "jev-generated-extended-lineage-candidate-revision-verification",
		NonAuthorizing: true,
	})
	if output.Status != "measured-with-counterexample" || output.Total != 1 ||
		output.CounterexampleCount != 1 || output.UnknownCount != 0 {
		t.Fatalf("output = %#v, want counterexample extended lineage candidate revision metric", output)
	}
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionMetricRetainsUnknown(t *testing.T) {
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionMetricInput{
		Projection:             clearExtendedLineageCandidateRevisionForMetric(t),
		MetricName:             "jev-generated-extended-lineage-candidate-revision-verification",
		AdditionalUnknownCount: 2,
		NonAuthorizing:         true,
	})
	if output.Status != "measured-with-unknown" || output.Total != 3 ||
		output.VerifiedCount != 1 || output.UnknownCount != 2 {
		t.Fatalf("output = %#v, want unknown-retaining extended lineage candidate revision metric", output)
	}
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionMetricFailsClosed(t *testing.T) {
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateRevisionMetricInput{
		Projection:      clearExtendedLineageCandidateRevisionForMetric(t),
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "metric-name" {
		t.Fatalf("output = %#v, want metric-name UNKNOWN", output)
	}
}