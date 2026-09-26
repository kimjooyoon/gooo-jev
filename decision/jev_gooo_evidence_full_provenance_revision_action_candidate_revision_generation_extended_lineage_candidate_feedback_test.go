package decision

import "testing"

func confirmedExtendedLineageCandidateMetricForFeedback(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricBinding {
	t.Helper()
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricInput{
		Projection:      clearExtendedLineageCandidateForMetric(t),
		MetricName:      "jev-generated-extended-lineage-candidate-verification",
		NonAuthorizing: true,
	})
}

func refutedExtendedLineageCandidateMetricForFeedback(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricBinding {
	t.Helper()
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricInput{
		Projection:      counterexampleExtendedLineageCandidateForMetric(t),
		MetricName:      "jev-generated-extended-lineage-candidate-verification",
		NonAuthorizing: true,
	})
}

func unknownExtendedLineageCandidateMetricForFeedback(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricBinding {
	t.Helper()
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateMetricInput{
		Projection:             clearExtendedLineageCandidateForMetric(t),
		MetricName:             "jev-generated-extended-lineage-candidate-verification",
		AdditionalUnknownCount: 1,
		NonAuthorizing:         true,
	})
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedbackConfirmed(t *testing.T) {
	metric := confirmedExtendedLineageCandidateMetricForFeedback(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedbackInput{
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

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedbackRefuted(t *testing.T) {
	metric := refutedExtendedLineageCandidateMetricForFeedback(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedbackInput{
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

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedbackHoldsUnknown(t *testing.T) {
	metric := unknownExtendedLineageCandidateMetricForFeedback(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedbackInput{
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

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedbackRejectsMismatchedCandidate(t *testing.T) {
	metric := confirmedExtendedLineageCandidateMetricForFeedback(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         "different-revision-candidate",
		ReplayObservationDigest: "extended-lineage-candidate-feedback-mismatch",
		NonAuthorizing:          true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "candidate-digest-mismatch" {
		t.Fatalf("output = %#v, want candidate-digest-mismatch UNKNOWN", output)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedbackRequiresEvidence(t *testing.T) {
	metric := confirmedExtendedLineageCandidateMetricForFeedback(t)
	metric.MetricEvidenceDigest = "tampered"
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionGenerationExtendedLineageCandidateFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "extended-lineage-candidate-feedback-tampered",
		NonAuthorizing:          true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-lineage-candidate-metric-validation" {
		t.Fatalf("output = %#v, want metric-validation UNKNOWN", output)
	}
}