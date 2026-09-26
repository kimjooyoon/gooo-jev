package decision

import "testing"

func TestProposeImprovementCandidateFromCounterexampleFailsClosed(t *testing.T) {
	_, err := ProposeImprovementCandidateFromCounterexample(
		DecisionConfidenceImprovementCounterexample{},
		ImprovementCycle{},
		ImprovementCycleReplay{},
		"candidate",
		"source",
	)
	if err == nil {
		t.Fatal("expected invalid counterexample to be rejected")
	}
}
