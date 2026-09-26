package decision

import "testing"

func TestPlanDecisionConfidenceFallbackFailsClosed(t *testing.T) {
	_, err := PlanDecisionConfidenceFallback(DecisionConfidenceImprovementSignal{})
	if err == nil {
		t.Fatal("expected invalid confidence signal to be rejected")
	}
}

func TestDecisionConfidenceFallbackPlanRejectsTampering(t *testing.T) {
	plan := DecisionConfidenceFallbackPlan{
		SignalDigest:   "signal",
		Action:         DecisionConfidenceFallbackReview,
		NonAuthorizing: true,
		PlanDigest:     "tampered",
	}
	if err := plan.Validate(); err == nil {
		t.Fatal("expected tampered fallback plan to be rejected")
	}
}
