package decision

import "testing"

func TestProposeDecisionConfidenceImprovementFailsClosed(t *testing.T) {
	_, err := ProposeDecisionConfidenceImprovement(
		DecisionConfidenceImprovementCounterexample{},
		ImprovementCandidate{},
		ImprovementReviewReceipt{},
		"change",
	)
	if err == nil {
		t.Fatal("expected invalid proposal inputs to be rejected")
	}
}

func TestDecisionConfidenceImprovementProposalRejectsTampering(t *testing.T) {
	proposal := DecisionConfidenceImprovementProposal{
		CounterexampleDigest:  "counterexample",
		CandidateDigest:      "candidate",
		ReviewDigest:         "review",
		ProposedChangeDigest: "change",
		NonExecuting:         true,
		ProposalDigest:       "tampered",
	}
	if err := proposal.Validate(); err == nil {
		t.Fatal("expected tampered proposal to be rejected")
	}
}
