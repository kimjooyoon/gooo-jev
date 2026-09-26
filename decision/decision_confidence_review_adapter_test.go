package decision

import "testing"

func TestReviewDecisionConfidenceCandidateFailsClosed(t *testing.T) {
	_, err := ReviewDecisionConfidenceCandidate(
		DecisionConfidenceImprovementSignal{},
		ImprovementCycle{},
		ImprovementCycleReplay{},
		"candidate",
		"source",
		"reviewer",
		"evidence",
		ReviewPassed,
	)
	if err == nil {
		t.Fatal("expected invalid confidence signal to be rejected")
	}
}
