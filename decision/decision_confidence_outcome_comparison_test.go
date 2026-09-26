package decision

import "testing"

func TestCompareDecisionConfidenceChangeOutcomesFailsClosed(t *testing.T) {
	_, err := CompareDecisionConfidenceChangeOutcomes(
		DecisionConfidenceChangeOutcomeMetric{},
		DecisionConfidenceChangeOutcomeMetric{},
	)
	if err == nil {
		t.Fatal("expected invalid outcome metrics to be rejected")
	}
}

func TestDecisionConfidenceChangeOutcomeComparisonRejectsTampering(t *testing.T) {
	comparison := DecisionConfidenceChangeOutcomeComparison{
		BaselineMetricDigest: "baseline",
		CurrentMetricDigest:  "current",
		Status:               DecisionConfidenceChangeOutcomeApplied,
		NonAuthorizing:       true,
		ComparisonDigest:     "tampered",
	}
	if err := comparison.Validate(); err == nil {
		t.Fatal("expected tampered outcome comparison to be rejected")
	}
}
