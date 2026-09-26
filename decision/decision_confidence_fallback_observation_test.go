package decision

import "testing"

func TestObserveDecisionConfidenceFallbackRejectsMissingEvidence(t *testing.T) {
	plan := DecisionConfidenceFallbackPlan{
		SignalDigest:   "signal",
		Action:         DecisionConfidenceFallbackReview,
		NonAuthorizing: true,
	}
	planDigest, err := Digest(plan)
	if err != nil {
		t.Fatal(err)
	}
	plan.PlanDigest = planDigest

	_, err = ObserveDecisionConfidenceFallback(
		plan,
		DecisionConfidenceFallbackReview,
		DecisionConfidenceFallbackObserved,
		"",
	)
	if err == nil {
		t.Fatal("expected observed fallback without evidence to be rejected")
	}
}

func TestObserveDecisionConfidenceFallbackPreservesUnknown(t *testing.T) {
	plan := DecisionConfidenceFallbackPlan{
		SignalDigest:   "signal",
		Action:         DecisionConfidenceFallbackReview,
		NonAuthorizing: true,
	}
	planDigest, err := Digest(plan)
	if err != nil {
		t.Fatal(err)
	}
	plan.PlanDigest = planDigest

	observation, err := ObserveDecisionConfidenceFallback(
		plan,
		DecisionConfidenceFallbackReview,
		DecisionConfidenceFallbackUnknown,
		"",
	)
	if err != nil {
		t.Fatal(err)
	}
	if observation.Status != DecisionConfidenceFallbackUnknown || !observation.NonAuthorizing {
		t.Fatalf("unexpected observation: %+v", observation)
	}
}
