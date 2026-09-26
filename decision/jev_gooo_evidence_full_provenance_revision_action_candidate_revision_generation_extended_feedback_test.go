package decision

import "testing"

func confirmedGeneratedExtendedCandidateRevisionMetricForFeedback(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedMetricBinding {
	t.Helper()
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedMetricInput{
		Projection:      clearGeneratedExtendedCandidateRevisionLSPForMetric(t),
		MetricName:      "jev-generated-extended-candidate-verification",
		NonAuthorizing: true,
	})
}

func refutedGeneratedExtendedCandidateRevisionMetricForFeedback(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedMetricBinding {
	t.Helper()
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedMetricInput{
		Projection:      counterexampleGeneratedExtendedCandidateRevisionLSPForMetric(t),
		MetricName:      "jev-generated-extended-candidate-verification",
		NonAuthorizing: true,
	})
}

func unknownGeneratedExtendedCandidateRevisionMetricForFeedback(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedMetricBinding {
	t.Helper()
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedMetricInput{
		Projection:             clearGeneratedExtendedCandidateRevisionLSPForMetric(t),
		MetricName:              "jev-generated-extended-candidate-verification",
		AdditionalUnknownCount: 1,
		NonAuthorizing:         true,
	})
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedbackConfirmed(t *testing.T) {
	metric := confirmedGeneratedExtendedCandidateRevisionMetricForFeedback(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "generated-extended-candidate-revision-feedback-confirmed",
		NonAuthorizing:          true,
	})
	if output.Status != "bound" || output.FeedbackStatus != jevReplayFeedbackConfirmed ||
		output.AggregationStatus != jevImprovementFeedbackStableForReview {
		t.Fatalf("output = %#v, want confirmed stable-for-review feedback", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedbackRefuted(t *testing.T) {
	metric := refutedGeneratedExtendedCandidateRevisionMetricForFeedback(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "generated-extended-candidate-revision-feedback-refuted",
		NonAuthorizing:          true,
	})
	if output.Status != "bound" || output.FeedbackStatus != jevReplayFeedbackRefuted ||
		output.AggregationStatus != jevImprovementFeedbackNeedsRevision {
		t.Fatalf("output = %#v, want refuted needs-revision feedback", output)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedbackHoldsUnknown(t *testing.T) {
	metric := unknownGeneratedExtendedCandidateRevisionMetricForFeedback(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "generated-extended-candidate-revision-feedback-unknown",
		NonAuthorizing:          true,
	})
	if output.Status != "bound" || output.FeedbackStatus != jevReplayFeedbackUnknown ||
		output.AggregationStatus != jevImprovementFeedbackHold {
		t.Fatalf("output = %#v, want unknown hold feedback", output)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedbackRejectsMismatchedCandidate(t *testing.T) {
	metric := confirmedGeneratedExtendedCandidateRevisionMetricForFeedback(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         "different-revision-candidate",
		ReplayObservationDigest: "generated-extended-candidate-revision-feedback-mismatch",
		NonAuthorizing:          true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "candidate-digest-mismatch" {
		t.Fatalf("output = %#v, want candidate-digest-mismatch UNKNOWN", output)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedbackRequiresEvidence(t *testing.T) {
	metric := confirmedGeneratedExtendedCandidateRevisionMetricForFeedback(t)
	metric.MetricEvidenceDigest = "tampered"
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "generated-extended-candidate-revision-feedback-tampered",
		NonAuthorizing:          true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-metric-validation" {
		t.Fatalf("output = %#v, want metric-validation UNKNOWN", output)
	}
}
