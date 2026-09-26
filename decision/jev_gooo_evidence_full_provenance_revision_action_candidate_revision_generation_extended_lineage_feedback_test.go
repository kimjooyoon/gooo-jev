package decision

import "testing"

func confirmedExtendedLineageCandidateRevisionMetricForFeedback(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageMetricBinding {
	t.Helper()
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageMetricInput{
		Projection:      clearExtendedLineageCandidateRevisionLSPForMetric(t),
		MetricName:      "jev-generated-extended-lineage-candidate-verification",
		NonAuthorizing: true,
	})
}

func refutedExtendedLineageCandidateRevisionMetricForFeedback(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageMetricBinding {
	t.Helper()
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageMetricInput{
		Projection:      counterexampleExtendedLineageCandidateRevisionLSPForMetric(t),
		MetricName:      "jev-generated-extended-lineage-candidate-verification",
		NonAuthorizing: true,
	})
}

func unknownExtendedLineageCandidateRevisionMetricForFeedback(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageMetricBinding {
	t.Helper()
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageMetricInput{
		Projection:             clearExtendedLineageCandidateRevisionLSPForMetric(t),
		MetricName:              "jev-generated-extended-lineage-candidate-verification",
		AdditionalUnknownCount: 1,
		NonAuthorizing:         true,
	})
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackConfirmed(t *testing.T) {
	metric := confirmedExtendedLineageCandidateRevisionMetricForFeedback(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "extended-lineage-candidate-feedback-confirmed",
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

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackRefuted(t *testing.T) {
	metric := refutedExtendedLineageCandidateRevisionMetricForFeedback(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "extended-lineage-candidate-feedback-refuted",
		NonAuthorizing:          true,
	})
	if output.Status != "bound" || output.FeedbackStatus != jevReplayFeedbackRefuted ||
		output.AggregationStatus != jevImprovementFeedbackNeedsRevision {
		t.Fatalf("output = %#v, want refuted needs-revision feedback", output)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackHoldsUnknown(t *testing.T) {
	metric := unknownExtendedLineageCandidateRevisionMetricForFeedback(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "extended-lineage-candidate-feedback-unknown",
		NonAuthorizing:          true,
	})
	if output.Status != "bound" || output.FeedbackStatus != jevReplayFeedbackUnknown ||
		output.AggregationStatus != jevImprovementFeedbackHold {
		t.Fatalf("output = %#v, want unknown hold feedback", output)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackRejectsMismatchedCandidate(t *testing.T) {
	metric := confirmedExtendedLineageCandidateRevisionMetricForFeedback(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         "different-revision-candidate",
		ReplayObservationDigest: "extended-lineage-candidate-feedback-mismatch",
		NonAuthorizing:          true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "candidate-digest-mismatch" {
		t.Fatalf("output = %#v, want candidate-digest-mismatch UNKNOWN", output)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackRequiresEvidence(t *testing.T) {
	metric := confirmedExtendedLineageCandidateRevisionMetricForFeedback(t)
	metric.MetricEvidenceDigest = "tampered"
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "extended-lineage-candidate-feedback-tampered",
		NonAuthorizing:          true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-lineage-metric-validation" {
		t.Fatalf("output = %#v, want metric-validation UNKNOWN", output)
	}
}
