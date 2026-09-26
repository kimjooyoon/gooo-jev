package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerCandidateReviewReceiptMetricIsOneHot(t *testing.T) {
	verified := ObserveDecisionConfidenceEvidenceLedgerCandidateReviewReceiptMetric(DecisionConfidenceEvidenceLedgerCandidateReviewReceiptVerification{
		ReceiptDigest:  "receipt",
		Status:         "verified",
		NonAuthorizing: true,
	})
	if verified.VerifiedCount != 1 || verified.MismatchCount != 0 || verified.UnknownCount != 0 || verified.Status != "verified" || !verified.NonAuthorizing {
		t.Fatalf("unexpected verified metric: %#v", verified)
	}
	mismatch := ObserveDecisionConfidenceEvidenceLedgerCandidateReviewReceiptMetric(DecisionConfidenceEvidenceLedgerCandidateReviewReceiptVerification{
		ReceiptDigest:  "receipt",
		Status:         "mismatch",
		NonAuthorizing: true,
	})
	if mismatch.VerifiedCount != 0 || mismatch.MismatchCount != 1 || mismatch.UnknownCount != 0 || mismatch.Status != "mismatch" || !mismatch.NonAuthorizing {
		t.Fatalf("unexpected mismatch metric: %#v", mismatch)
	}
}
