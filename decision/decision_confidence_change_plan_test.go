package decision

import "testing"

func TestPlanDecisionConfidenceChangeFailsClosed(t *testing.T) {
	_, err := PlanDecisionConfidenceChange(DecisionConfidenceImprovementProposal{}, "source", "write-set")
	if err == nil {
		t.Fatal("expected invalid proposal to be rejected")
	}
}

func TestDecisionConfidenceChangePlanRejectsTampering(t *testing.T) {
	plan := DecisionConfidenceChangePlan{
		ProposalDigest:   "proposal",
		SourceReference:  "source",
		ChangeDigest:     "change",
		WriteSetDigest:   "write-set",
		NonExecuting:     true,
		ChangePlanDigest: "tampered",
	}
	if err := plan.Validate(); err == nil {
		t.Fatal("expected tampered change plan to be rejected")
	}
}
