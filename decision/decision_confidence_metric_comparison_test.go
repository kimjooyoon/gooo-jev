package decision

import "testing"

func TestCompareDecisionConfidenceAdmissionMetricsFailsClosed(t *testing.T) {
	_, err := CompareDecisionConfidenceAdmissionMetrics(
		DecisionConfidenceAdmissionMetricObservation{},
		DecisionConfidenceAdmissionMetricObservation{},
	)
	if err == nil {
		t.Fatal("expected invalid metrics to be rejected")
	}
}

func TestDecisionConfidenceMetricComparisonRejectsTampering(t *testing.T) {
	comparison := DecisionConfidenceMetricComparison{
		BaselineMetricDigest: "baseline",
		CurrentMetricDigest:  "current",
		Status:               DecisionConfidenceMetricImproved,
		NonAuthorizing:       true,
		ComparisonDigest:     "tampered",
	}
	if err := comparison.Validate(); err == nil {
		t.Fatal("expected tampered comparison to be rejected")
	}
}
