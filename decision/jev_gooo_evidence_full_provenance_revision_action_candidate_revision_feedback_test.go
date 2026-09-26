package decision

import "testing"

func confirmedActionDerivedCandidateRevisionMetric(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricBinding {
	t.Helper()
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricInput{
		Projection:      clearActionDerivedCandidateRevisionLSPForMetric(t),
		MetricName:      "jev-action-derived-candidate-revision-verification",
		NonAuthorizing: true,
	})
}

func refutedActionDerivedCandidateRevisionMetric(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricBinding {
	t.Helper()
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricInput{
		Projection:      counterexampleActionDerivedCandidateRevisionLSPForMetric(t),
		MetricName:      "jev-action-derived-candidate-revision-verification",
		NonAuthorizing: true,
	})
}

func unknownActionDerivedCandidateRevisionMetric(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricBinding {
	t.Helper()
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionMetricInput{
		Projection:             clearActionDerivedCandidateRevisionLSPForMetric(t),
		MetricName:             "jev-action-derived-candidate-revision-verification",
		AdditionalUnknownCount: 1,
		NonAuthorizing:         true,
	})
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackConfirmed(t *testing.T) {
	metric := confirmedActionDerivedCandidateRevisionMetric(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-revision-feedback-confirmed",
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

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackRefuted(t *testing.T) {
	metric := refutedActionDerivedCandidateRevisionMetric(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-revision-feedback-refuted",
		NonAuthorizing:          true,
	})
	if output.Status != "bound" || output.FeedbackStatus != jevReplayFeedbackRefuted ||
		output.AggregationStatus != jevImprovementFeedbackNeedsRevision {
		t.Fatalf("output = %#v, want refuted needs-revision feedback", output)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackHoldsUnknown(t *testing.T) {
	metric := unknownActionDerivedCandidateRevisionMetric(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-revision-feedback-unknown",
		NonAuthorizing:          true,
	})
	if output.Status != "bound" || output.FeedbackStatus != jevReplayFeedbackUnknown ||
		output.AggregationStatus != jevImprovementFeedbackHold {
		t.Fatalf("output = %#v, want unknown hold feedback", output)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackRejectsMismatchedCandidate(t *testing.T) {
	metric := confirmedActionDerivedCandidateRevisionMetric(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         "different-revision-candidate",
		ReplayObservationDigest: "candidate-revision-feedback-mismatch",
		NonAuthorizing:          true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "candidate-digest-mismatch" {
		t.Fatalf("output = %#v, want candidate-digest-mismatch UNKNOWN", output)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackRequiresEvidence(t *testing.T) {
	metric := confirmedActionDerivedCandidateRevisionMetric(t)
	metric.MetricEvidenceDigest = "tampered"
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-revision-feedback-tampered",
		NonAuthorizing:          true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-revision-metric-validation" {
		t.Fatalf("output = %#v, want metric-validation UNKNOWN", output)
	}
}