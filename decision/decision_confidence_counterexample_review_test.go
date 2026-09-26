package decision

import "testing"

func TestReviewDecisionConfidenceCounterexampleFailsClosed(t *testing.T) {
	_, err := ReviewDecisionConfidenceCounterexample(
		DecisionConfidenceImprovementCounterexample{},
		ImprovementCycle{},
		ImprovementCycleReplay{},
		"candidate",
		"source",
		"reviewer",
		"evidence",
		ReviewPassed,
	)
	if err == nil {
		t.Fatal("expected invalid counterexample to be rejected")
	}
}
