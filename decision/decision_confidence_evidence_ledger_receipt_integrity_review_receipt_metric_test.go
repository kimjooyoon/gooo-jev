package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetricIsOneHot(t *testing.T) {
	verified := ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetric(DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptVerification{
		ReceiptDigest:  "receipt",
		Status:         "verified",
		NonAuthorizing: true,
	})
	if verified.VerifiedCount != 1 || verified.MismatchCount != 0 || verified.UnknownCount != 0 || verified.Status != "verified" || !verified.NonAuthorizing {
		t.Fatalf("unexpected verified metric: %#v", verified)
	}
	unknown := ObserveDecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptMetric(DecisionConfidenceEvidenceLedgerReceiptIntegrityReviewReceiptVerification{
		ReceiptDigest:  "receipt",
		Status:         "unknown",
		NonAuthorizing: true,
	})
	if unknown.VerifiedCount != 0 || unknown.MismatchCount != 0 || unknown.UnknownCount != 1 || unknown.Status != "unknown" || !unknown.NonAuthorizing {
		t.Fatalf("unexpected unknown metric: %#v", unknown)
	}
}
