package decision

import "testing"

func suspendedPlanObservationForMetric(t *testing.T) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationBinding {
	t.Helper()
	return ObserveExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservation(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationInput{
		Plan:                      executionPlanForObservation(t),
		LifecycleEvent:            "suspended",
		ObservationEvidenceDigest: "observation-evidence-metric",
		ReverseObservationDigest:  "reverse-observation-metric",
		NonAuthorizing:            true,
	})
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricPreservesValue(t *testing.T) {
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricInput{
		Observation:          suspendedPlanObservationForMetric(t),
		MetricName:           "revision-latency",
		MetricValue:          42,
		MetricUnit:           "milliseconds",
		MetricEvidenceDigest: "metric-evidence-latency",
		NonAuthorizing:        true,
	})
	if output.Status != "bound" || output.MetricStatus != "observed" ||
		output.MetricValue != 42 || output.MetricUnit != "milliseconds" ||
		output.PermissionState != "not-authorized" {
		t.Fatalf("output = %#v, want externally measured metric", output)
	}
	if err := output.Validate(); err != nil {
		t.Fatalf("output should validate: %v", err)
	}
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricAllowsZero(t *testing.T) {
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricInput{
		Observation:          suspendedPlanObservationForMetric(t),
		MetricName:           "revision-cost",
		MetricValue:          0,
		MetricUnit:           "microcredits",
		MetricEvidenceDigest: "metric-evidence-zero",
		NonAuthorizing:        true,
	})
	if output.Status != "bound" || output.MetricValue != 0 {
		t.Fatalf("output = %#v, want zero metric preserved", output)
	}
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricRejectsNegativeValue(t *testing.T) {
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricInput{
		Observation:          suspendedPlanObservationForMetric(t),
		MetricName:           "revision-cost",
		MetricValue:          -1,
		MetricUnit:           "microcredits",
		MetricEvidenceDigest: "metric-evidence-negative",
		NonAuthorizing:        true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "metric-value" {
		t.Fatalf("output = %#v, want metric-value UNKNOWN", output)
	}
}

func TestMeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricRejectsObservationTampering(t *testing.T) {
	observation := suspendedPlanObservationForMetric(t)
	observation.ReverseObservationDigest = "tampered"
	output := MeasureExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetric(ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionCandidateRevisionExecutionPlanObservationMetricInput{
		Observation:          observation,
		MetricName:           "revision-latency",
		MetricValue:          42,
		MetricUnit:           "milliseconds",
		MetricEvidenceDigest: "metric-evidence-tampered",
		NonAuthorizing:        true,
	})
	if output.Status != "UNKNOWN" || output.MissingStage != "revision-action-candidate-generation-extended-lineage-candidate-execution-plan-observation-validation" {
		t.Fatalf("output = %#v, want observation validation UNKNOWN", output)
	}
}