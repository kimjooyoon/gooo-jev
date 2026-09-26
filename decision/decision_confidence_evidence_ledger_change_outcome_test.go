package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerChangeOutcomeDoesNotClaimImprovement(t *testing.T) {
	plan := DecisionConfidenceEvidenceLedgerChangePlanObservation{
		ReviewOutcomeDigest: "review-digest",
		WriteSetDigest:      "write-set-digest",
		PlanDigest:          "plan-digest",
		Status:              "plan-observed",
		NonAuthorizing:      true,
	}
	outcome, err := ObserveDecisionConfidenceEvidenceLedgerChangeOutcome(plan, "applied")
	if err != nil {
		t.Fatalf("observe applied outcome: %v", err)
	}
	if outcome.Status != "applied" || outcome.PlanDigest == "" || !outcome.NonAuthorizing {
		t.Fatalf("unexpected applied outcome: %#v", outcome)
	}

	outcome, err = ObserveDecisionConfidenceEvidenceLedgerChangeOutcome(plan, "improved")
	if err != nil {
		t.Fatalf("observe invalid outcome: %v", err)
	}
	if outcome.Status != "unknown" || !outcome.NonAuthorizing {
		t.Fatalf("invalid outcome was promoted: %#v", outcome)
	}

	plan.Status = "unknown"
	outcome, err = ObserveDecisionConfidenceEvidenceLedgerChangeOutcome(plan, "rolled-back")
	if err != nil {
		t.Fatalf("observe held outcome: %v", err)
	}
	if outcome.Status != "unknown" || !outcome.NonAuthorizing {
		t.Fatalf("held plan escaped unknown: %#v", outcome)
	}
}
