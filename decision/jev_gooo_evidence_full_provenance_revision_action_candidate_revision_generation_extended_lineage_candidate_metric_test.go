package decision

import "testing"

func clearExtendedLineageCandidateForMetric(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateLSPBinding {
	t.Helper()
	return ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateLSPInput{
		Verification: VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateVerificationInput{
			Guard: admittedExtendedLineageCandidateForCandidateVerification(t),
			ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
				Status:         "reproduced",
				EvidenceDigest: "extended-lineage-candidate-metric-reproduced",
				NonAuthorizing: true,
			},
			NonAuthorizing: true,
		}),
		EvidencePrefixDigest: "extended-lineage-candidate-metric-clear",
		MissingStageIndex:   -1,
		NonAuthorizing:      true,
	})
}

func counterexampleExtendedLineageCandidateForMetric(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateLSPBinding {
	t.Helper()
	return ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateLSPInput{
		Verification: VerifyExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidate(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateVerificationInput{
			Guard: admittedExtendedLineageCandidateForCandidateVerification(t),
			ReverseObservation: ExecutionEnvelopeReverseObservationOutput{
				Status:         "counterexample",
				EvidenceDigest: "extended-lineage-candidate-metric-counterexample",
				FirstMismatch:  "extended-lineage-output",
				NonAuthorizing: true,
			},
			NonAuthorizing: true,
		}),
		MissingStageIndex:   27,
		EvidencePrefixDigest: "extended-lineage-candidate-metric-counterexample",
		NonAuthorizing:      true,
	})
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetric(t *testing.T) {
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricInput{
		Projection:      clearExtendedLineageCandidateForMetric(t),
		MetricName:      "jev-generated-extended-lineage-candidate-verification",
		NonAuthorizing: true,
	})
	if output.Status != "measured" || output.Total != 1 || output.VerifiedCount != 1 ||
		output.SourceCandidateID == "" || output.RevisionCandidateDigest == "" ||
		output.CandidateDigest == "" || output.GuardEvidenceDigest == "" ||
		output.ReverseEvidenceDigest == "" {
		t.Fatalf("output = %#v, want measured extended lineage candidate metric", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricRetainsCounterexample(t *testing.T) {
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricInput{
		Projection:      counterexampleExtendedLineageCandidateForMetric(t),
		MetricName:      "jev-generated-extended-lineage-candidate-verification",
		NonAuthorizing: true,
	})
	if output.Status != "measured-with-counterexample" || output.Total != 1 ||
		output.CounterexampleCount != 1 || output.UnknownCount != 0 {
		t.Fatalf("output = %#v, want counterexample extended lineage candidate metric", output)
	}
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricRetainsUnknown(t *testing.T) {
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricInput{
		Projection:             clearExtendedLineageCandidateForMetric(t),
		MetricName:             "jev-generated-extended-lineage-candidate-verification",
		AdditionalUnknownCount: 2,
		NonAuthorizing:         true,
	})
	if output.Status != "measured-with-unknown" || output.Total != 3 ||
		output.VerifiedCount != 1 || output.UnknownCount != 2 {
		t.Fatalf("output = %#v, want unknown-retaining extended lineage candidate metric", output)
	}
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricFailsClosed(t *testing.T) {
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricInput{
		Projection:      clearExtendedLineageCandidateForMetric(t),
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "metric-name" {
		t.Fatalf("output = %#v, want metric-name UNKNOWN", output)
	}
}