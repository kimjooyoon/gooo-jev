package decision

import "testing"

func clearRevisionActionMetric(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricBinding {
	t.Helper()
	projection := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPInput{
		Verification:        verifiedRevisionActionForLSP(t),
		EvidencePrefixDigest: "prefix-feedback-clear",
		NonAuthorizing:      true,
	})
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricInput{
		Projection:      projection,
		MetricName:      "jev-revision-action-verification",
		NonAuthorizing: true,
	})
}

func counterexampleRevisionActionMetric(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricBinding {
	t.Helper()
	projection := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPInput{
		Verification:        reviewRevisionActionForLSP(t),
		MissingStageIndex:   7,
		EvidencePrefixDigest: "prefix-feedback-counterexample",
		NonAuthorizing:      true,
	})
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricInput{
		Projection:      projection,
		MetricName:      "jev-revision-action-verification",
		NonAuthorizing: true,
	})
}

func unknownRevisionActionMetric(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricBinding {
	t.Helper()
	projection := ProjectExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSP(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionLSPInput{
		Verification:        verifiedRevisionActionForLSP(t),
		EvidencePrefixDigest: "prefix-feedback-unknown",
		NonAuthorizing:      true,
	})
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionMetricInput{
		Projection:             projection,
		MetricName:             "jev-revision-action-verification",
		AdditionalUnknownCount: 1,
		NonAuthorizing:         true,
	})
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackConfirmed(t *testing.T) {
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackInput{
		Metric:                  clearRevisionActionMetric(t),
		CandidateDigest:         "candidate-confirmed",
		ReplayObservationDigest: "replay-confirmed",
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

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackRefuted(t *testing.T) {
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackInput{
		Metric:                  counterexampleRevisionActionMetric(t),
		CandidateDigest:         "candidate-refuted",
		ReplayObservationDigest: "replay-refuted",
		NonAuthorizing:          true,
	})
	if output.Status != "bound" || output.FeedbackStatus != jevReplayFeedbackRefuted ||
		output.AggregationStatus != jevImprovementFeedbackNeedsRevision {
		t.Fatalf("output = %#v, want refuted needs-revision feedback", output)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackHoldsUnknown(t *testing.T) {
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackInput{
		Metric:                  unknownRevisionActionMetric(t),
		CandidateDigest:         "candidate-unknown",
		ReplayObservationDigest: "replay-unknown",
		NonAuthorizing:          true,
	})
	if output.Status != "bound" || output.FeedbackStatus != jevReplayFeedbackUnknown ||
		output.AggregationStatus != jevImprovementFeedbackHold {
		t.Fatalf("output = %#v, want unknown hold feedback", output)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackRequiresEvidence(t *testing.T) {
	metric := clearRevisionActionMetric(t)
	metric.MetricEvidenceDigest = "tampered"
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionFeedbackInput{
		Metric:                  metric,
		CandidateDigest:         "candidate-tampered",
		ReplayObservationDigest: "replay-tampered",
		NonAuthorizing:          true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-metric-validation" {
		t.Fatalf("output = %#v, want revision-action-metric-validation UNKNOWN", output)
	}
}