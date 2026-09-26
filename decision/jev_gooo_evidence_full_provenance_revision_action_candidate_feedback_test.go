package decision

import "testing"

func confirmedActionDerivedCandidateMetric(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricBinding {
	t.Helper()
	projection := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPInput{
		Verification:        verifiedActionDerivedCandidateForLSP(t),
		EvidencePrefixDigest: "candidate-feedback-confirmed",
		NonAuthorizing:      true,
	})
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricInput{
		Projection:      projection,
		MetricName:      "jev-action-derived-candidate-verification",
		NonAuthorizing: true,
	})
}

func refutedActionDerivedCandidateMetric(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricBinding {
	t.Helper()
	projection := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPInput{
		Verification:        reviewActionDerivedCandidateForLSP(t),
		MissingStageIndex:   13,
		EvidencePrefixDigest: "candidate-feedback-refuted",
		NonAuthorizing:      true,
	})
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricInput{
		Projection:      projection,
		MetricName:      "jev-action-derived-candidate-verification",
		NonAuthorizing: true,
	})
}

func unknownActionDerivedCandidateMetric(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricBinding {
	t.Helper()
	projection := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateLSPInput{
		Verification:        verifiedActionDerivedCandidateForLSP(t),
		EvidencePrefixDigest: "candidate-feedback-unknown",
		NonAuthorizing:      true,
	})
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateMetricInput{
		Projection:             projection,
		MetricName:             "jev-action-derived-candidate-verification",
		AdditionalUnknownCount: 1,
		NonAuthorizing:         true,
	})
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackConfirmed(t *testing.T) {
	metric := confirmedActionDerivedCandidateMetric(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-replay-confirmed",
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

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackRefuted(t *testing.T) {
	metric := refutedActionDerivedCandidateMetric(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-replay-refuted",
		NonAuthorizing:          true,
	})
	if output.Status != "bound" || output.FeedbackStatus != jevReplayFeedbackRefuted ||
		output.AggregationStatus != jevImprovementFeedbackNeedsRevision {
		t.Fatalf("output = %#v, want refuted needs-revision feedback", output)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackHoldsUnknown(t *testing.T) {
	metric := unknownActionDerivedCandidateMetric(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-replay-unknown",
		NonAuthorizing:          true,
	})
	if output.Status != "bound" || output.FeedbackStatus != jevReplayFeedbackUnknown ||
		output.AggregationStatus != jevImprovementFeedbackHold {
		t.Fatalf("output = %#v, want unknown hold feedback", output)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackRejectsMismatchedCandidate(t *testing.T) {
	metric := confirmedActionDerivedCandidateMetric(t)
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         "different-candidate",
		ReplayObservationDigest: "candidate-replay-mismatch",
		NonAuthorizing:          true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "candidate-digest-mismatch" {
		t.Fatalf("output = %#v, want candidate-digest-mismatch UNKNOWN", output)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackRequiresEvidence(t *testing.T) {
	metric := confirmedActionDerivedCandidateMetric(t)
	metric.MetricEvidenceDigest = "tampered"
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         metric.RevisionCandidateDigest,
		ReplayObservationDigest: "candidate-replay-tampered",
		NonAuthorizing:          true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-metric-validation" {
		t.Fatalf("output = %#v, want metric-validation UNKNOWN", output)
	}
}