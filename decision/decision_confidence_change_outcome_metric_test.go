package decision

import "testing"

func TestObserveDecisionConfidenceChangeOutcomeMetricFailsClosed(t *testing.T) {
	_, err := ObserveDecisionConfidenceChangeOutcomeMetric(DecisionConfidenceChangePlanDispositionObservation{})
	if err == nil {
		t.Fatal("expected invalid disposition observation to be rejected")
	}
}

func TestDecisionConfidenceChangeOutcomeMetricRejectsNonOneHotCounts(t *testing.T) {
	metric := DecisionConfidenceChangeOutcomeMetric{
		ObservationDigest: "observation",
		AppliedCount:      1,
		UnknownCount:      1,
		NonAuthorizing:    true,
		MetricDigest:      "invalid",
	}
	if err := metric.Validate(); err == nil {
		t.Fatal("expected non-one-hot outcome metric to be rejected")
	}
}
