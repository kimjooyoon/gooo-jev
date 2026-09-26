package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerCandidateReviewMetricIsOneHot(t *testing.T) {
	ready := ObserveDecisionConfidenceEvidenceLedgerCandidateReviewMetric(DecisionConfidenceEvidenceLedgerCandidateReviewSignal{
		ReviewOutcomeDigest: "digest-ready",
		Status:              "ready-for-candidate-review",
		NonAuthorizing:      true,
	})
	if ready.ReadyCount != 1 || ready.RejectedCount != 0 || ready.HoldCount != 0 || ready.Status != "ready-for-candidate-review" || !ready.NonAuthorizing {
		t.Fatalf("unexpected ready metric: %#v", ready)
	}

	hold := ObserveDecisionConfidenceEvidenceLedgerCandidateReviewMetric(DecisionConfidenceEvidenceLedgerCandidateReviewSignal{Status: "hold", NonAuthorizing: true})
	if hold.ReadyCount != 0 || hold.RejectedCount != 0 || hold.HoldCount != 1 || hold.Status != "hold" || !hold.NonAuthorizing {
		t.Fatalf("unexpected hold metric: %#v", hold)
	}
}
