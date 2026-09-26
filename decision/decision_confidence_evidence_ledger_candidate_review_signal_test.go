package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerCandidateReviewSignalDoesNotCreateCandidate(t *testing.T) {
	outcome := DecisionConfidenceEvidenceLedgerExternalReviewOutcome{
		HandoffDigest:  "handoff-digest",
		Decision:       "accepted-for-analysis",
		Status:         "accepted-for-analysis",
		NonAuthorizing: true,
	}
	signal, err := DeriveDecisionConfidenceEvidenceLedgerCandidateReviewSignal(outcome)
	if err != nil {
		t.Fatalf("derive candidate review signal: %v", err)
	}
	if signal.Status != "ready-for-candidate-review" || signal.ReviewOutcomeDigest == "" || !signal.NonAuthorizing {
		t.Fatalf("unexpected candidate review signal: %#v", signal)
	}

	outcome.Status = "unknown"
	signal, err = DeriveDecisionConfidenceEvidenceLedgerCandidateReviewSignal(outcome)
	if err != nil {
		t.Fatalf("derive held signal: %v", err)
	}
	if signal.Status != "hold" || !signal.NonAuthorizing {
		t.Fatalf("unknown review escaped hold: %#v", signal)
	}
}
