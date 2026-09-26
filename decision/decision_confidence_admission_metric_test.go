package decision

import "testing"

func TestObserveDecisionConfidenceAdmissionMetricFailsClosed(t *testing.T) {
	_, err := ObserveDecisionConfidenceAdmissionMetric(DecisionConfidenceExecutionAdmissionEvidence{})
	if err == nil {
		t.Fatal("expected invalid admission evidence to be rejected")
	}
}

func TestDecisionConfidenceAdmissionMetricRejectsNonOneHotCounts(t *testing.T) {
	metric := DecisionConfidenceAdmissionMetricObservation{
		AdmissionDigest: "admission",
		EligibleCount:  1,
		UnknownCount:   1,
		NonAuthorizing: true,
		MetricDigest:   "invalid",
	}
	if err := metric.Validate(); err == nil {
		t.Fatal("expected non-one-hot metric to be rejected")
	}
}
