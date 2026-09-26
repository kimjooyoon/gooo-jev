package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricComparisonPreservesChange(t *testing.T) {
	previous := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetric{ReceiptDigest: "receipt", VerifiedCount: 1, Status: "verified", NonAuthorizing: true}
	current := DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetric{ReceiptDigest: "receipt", MismatchCount: 1, Status: "mismatch", NonAuthorizing: true}
	comparison := CompareDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetric(previous, current)
	if comparison.Delta != "changed" || !comparison.NonAuthorizing {
		t.Fatalf("integrity degradation was not preserved: %#v", comparison)
	}

	current.Status = "verified"
	current.MismatchCount = 0
	current.VerifiedCount = 1
	comparison = CompareDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetric(previous, current)
	if comparison.Delta != "unchanged" || !comparison.NonAuthorizing {
		t.Fatalf("unchanged integrity was not preserved: %#v", comparison)
	}
}
