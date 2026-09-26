package decision

import "testing"

func TestRecordDecisionConfidenceCounterexampleFailsClosed(t *testing.T) {
	_, err := RecordDecisionConfidenceComparisonCounterexample(DecisionConfidenceMetricComparison{})
	if err == nil {
		t.Fatal("expected invalid comparison to be rejected")
	}
}

func TestDecisionConfidenceCounterexampleRejectsImprovement(t *testing.T) {
	comparison := DecisionConfidenceMetricComparison{
		BaselineMetricDigest: "baseline",
		CurrentMetricDigest:  "current",
		Status:               DecisionConfidenceMetricImproved,
		NonAuthorizing:       true,
	}
	digest, err := Digest(comparison)
	if err != nil {
		t.Fatal(err)
	}
	comparison.ComparisonDigest = digest
	if _, err := RecordDecisionConfidenceComparisonCounterexample(comparison); err == nil {
		t.Fatal("expected improvement to avoid counterexample creation")
	}
}

func TestDecisionConfidenceCounterexampleRejectsTampering(t *testing.T) {
	counterexample := DecisionConfidenceImprovementCounterexample{
		ComparisonDigest:     "comparison",
		ReasonStatus:         DecisionConfidenceMetricRegressed,
		RequiresReview:       true,
		NonAuthorizing:       true,
		CounterexampleDigest: "tampered",
	}
	if err := counterexample.Validate(); err == nil {
		t.Fatal("expected tampered counterexample to be rejected")
	}
}
