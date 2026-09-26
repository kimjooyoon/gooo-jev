package decision

import "testing"

func metricForObservationFeedback(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricBinding {
	t.Helper()
	return MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricInput{
		Observation:          suspendedPlanObservationForMetric(t),
		MetricName:           "revision-latency",
		MetricValue:          42,
		MetricUnit:            "milliseconds",
		MetricEvidenceDigest: "metric-evidence-feedback",
		NonAuthorizing:        true,
	})
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackConfirmed(t *testing.T) {
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackInput{
		Metric:                metricForObservationFeedback(t),
		FeedbackChoice:        "confirmed",
		FeedbackEvidenceDigest: "feedback-evidence-confirmed",
		NonAuthorizing:        true,
	})
	if output.Status != "bound" || output.FeedbackStatus != "confirmed" ||
		output.Confirmed != 1 || output.Refuted != 0 || output.Unknown != 0 ||
		output.PermissionState != "not-authorized" {
		t.Fatalf("output = %#v, want confirmed feedback aggregation", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackUnknown(t *testing.T) {
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackInput{
		Metric:                metricForObservationFeedback(t),
		FeedbackChoice:        "unknown",
		FeedbackEvidenceDigest: "feedback-evidence-unknown",
		NonAuthorizing:        true,
	})
	if output.Status != "bound" || output.FeedbackStatus != "unknown" || output.Unknown != 1 {
		t.Fatalf("output = %#v, want unknown feedback aggregation", output)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackRejectsChoice(t *testing.T) {
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackInput{
		Metric:                metricForObservationFeedback(t),
		FeedbackChoice:        "maybe",
		FeedbackEvidenceDigest: "feedback-evidence-invalid",
		NonAuthorizing:        true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "feedback-choice" {
		t.Fatalf("output = %#v, want feedback-choice UNKNOWN", output)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackRequiresEvidence(t *testing.T) {
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackInput{
		Metric:         metricForObservationFeedback(t),
		FeedbackChoice: "refuted",
		NonAuthorizing: true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "feedback-evidence" {
		t.Fatalf("output = %#v, want feedback-evidence UNKNOWN", output)
	}
}

func TestBindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackRejectsMetricTampering(t *testing.T) {
	metric := metricForObservationFeedback(t)
	metric.MetricEvidenceDigest = "tampered"
	output := BindExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedback(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationFeedbackInput{
		Metric:                metric,
		FeedbackChoice:        "refuted",
		FeedbackEvidenceDigest: "feedback-evidence-tampered",
		NonAuthorizing:        true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-lineage-candidate-execution-plan-observation-metric-validation" {
		t.Fatalf("output = %#v, want metric validation UNKNOWN", output)
	}
}