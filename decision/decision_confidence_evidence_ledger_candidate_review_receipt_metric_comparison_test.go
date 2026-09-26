package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerCandidateReviewReceiptMetricComparisonPreservesChange(t *testing.T) {
	previous := DecisionConfidenceEvidenceLedgerCandidateReviewReceiptMetric{
		ReceiptDigest:  "receipt",
		VerifiedCount:  1,
		Status:         "verified",
		NonAuthorizing: true,
	}
	current := DecisionConfidenceEvidenceLedgerCandidateReviewReceiptMetric{
		ReceiptDigest:  "receipt",
		MismatchCount:  1,
		Status:         "mismatch",
		NonAuthorizing: true,
	}
	comparison := CompareDecisionConfidenceEvidenceLedgerCandidateReviewReceiptMetric(previous, current)
	if comparison.Delta != "changed" || !comparison.NonAuthorizing {
		t.Fatalf("receipt degradation was not preserved: %#v", comparison)
	}

	current.Status = "verified"
	current.MismatchCount = 0
	current.VerifiedCount = 1
	comparison = CompareDecisionConfidenceEvidenceLedgerCandidateReviewReceiptMetric(previous, current)
	if comparison.Delta != "unchanged" || !comparison.NonAuthorizing {
		t.Fatalf("unchanged receipt was not preserved: %#v", comparison)
	}
}
