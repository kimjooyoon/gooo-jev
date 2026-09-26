package decision

import "testing"

func TestRecordDecisionConfidenceOutcomeCounterexampleFailsClosed(t *testing.T) {
	_, err := RecordDecisionConfidenceOutcomeCounterexample(DecisionConfidenceChangeOutcomeComparison{})
	if err == nil {
		t.Fatal("expected invalid outcome comparison to be rejected")
	}
}

func TestRecordDecisionConfidenceOutcomeCounterexampleRejectsApplied(t *testing.T) {
	comparison := DecisionConfidenceChangeOutcomeComparison{
		BaselineMetricDigest: "baseline",
		CurrentMetricDigest:  "current",
		Status:               DecisionConfidenceChangeOutcomeApplied,
		NonAuthorizing:       true,
	}
	digest, err := Digest(comparison)
	if err != nil {
		t.Fatal(err)
	}
	comparison.ComparisonDigest = digest
	if _, err := RecordDecisionConfidenceOutcomeCounterexample(comparison); err == nil {
		t.Fatal("expected applied outcome not to become a counterexample")
	}
}
