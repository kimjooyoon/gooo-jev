package decision

import "testing"

func TestProposeImprovementCandidateFromDecisionConfidenceFailsClosed(t *testing.T) {
	_, err := ProposeImprovementCandidateFromDecisionConfidence(
		DecisionConfidenceImprovementSignal{},
		ImprovementCycle{},
		ImprovementCycleReplay{},
		"candidate",
		"source",
	)
	if err == nil {
		t.Fatal("expected invalid confidence signal to be rejected")
	}
}
