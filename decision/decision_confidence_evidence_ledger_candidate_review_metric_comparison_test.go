package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerCandidateReviewMetricComparisonPreservesDecline(t *testing.T) {
	previous := DecisionConfidenceEvidenceLedgerCandidateReviewMetric{
		SignalDigest:   "previous",
		ReadyCount:     1,
		Status:         "ready-for-candidate-review",
		NonAuthorizing: true,
	}
	current := DecisionConfidenceEvidenceLedgerCandidateReviewMetric{
		SignalDigest:   "current",
		HoldCount:      1,
		Status:         "hold",
		NonAuthorizing: true,
	}
	comparison := CompareDecisionConfidenceEvidenceLedgerCandidateReviewMetric(previous, current)
	if comparison.Delta != "declined" || !comparison.NonAuthorizing {
		t.Fatalf("decline was not preserved: %#v", comparison)
	}

	invalid := current
	invalid.HoldCount = 0
	comparison = CompareDecisionConfidenceEvidenceLedgerCandidateReviewMetric(previous, invalid)
	if comparison.Delta != "inconclusive" || !comparison.NonAuthorizing {
		t.Fatalf("invalid metric was not inconclusive: %#v", comparison)
	}
}
