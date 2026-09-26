package decision

import "testing"

func TestDecisionConfidenceEvidenceLedgerChangeOutcomeComparisonPreservesRollback(t *testing.T) {
	previous := DecisionConfidenceEvidenceLedgerChangeOutcome{PlanDigest: "plan", Outcome: "applied", Status: "applied", NonAuthorizing: true}
	current := DecisionConfidenceEvidenceLedgerChangeOutcome{PlanDigest: "plan", Outcome: "rolled-back", Status: "rolled-back", NonAuthorizing: true}
	comparison := CompareDecisionConfidenceEvidenceLedgerChangeOutcome(previous, current)
	if comparison.Delta != "rolled-back" || !comparison.NonAuthorizing {
		t.Fatalf("rollback was not preserved: %#v", comparison)
	}

	current.PlanDigest = "other-plan"
	comparison = CompareDecisionConfidenceEvidenceLedgerChangeOutcome(previous, current)
	if comparison.Delta != "unknown" || !comparison.NonAuthorizing {
		t.Fatalf("different plan escaped unknown: %#v", comparison)
	}
}
