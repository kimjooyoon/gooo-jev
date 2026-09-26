package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerChangeOutcomeFeedbackPreservesRegression(t *testing.T) {
	feedback, err := DeriveDecisionConfidenceEvidenceLedgerChangeOutcomeFeedback(DecisionConfidenceEvidenceLedgerChangeOutcomeComparison{
		PreviousPlanDigest: "plan",
		CurrentPlanDigest:  "plan",
		Delta:              "rolled-back",
		NonAuthorizing:     true,
	})
	if err != nil {
		t.Fatalf("derive rollback feedback: %v", err)
	}
	if feedback.Status != "review-required" || feedback.PlanDigest == "" || !feedback.NonAuthorizing {
		t.Fatalf("unexpected rollback feedback: %#v", feedback)
	}

	feedback, err = DeriveDecisionConfidenceEvidenceLedgerChangeOutcomeFeedback(DecisionConfidenceEvidenceLedgerChangeOutcomeComparison{
		PreviousPlanDigest: "plan",
		CurrentPlanDigest:  "plan",
		Delta:              "unknown",
		NonAuthorizing:     true,
	})
	if err != nil {
		t.Fatalf("derive unknown feedback: %v", err)
	}
	if feedback.Status != "hold" || !feedback.NonAuthorizing {
		t.Fatalf("unknown outcome escaped hold: %#v", feedback)
	}
}
