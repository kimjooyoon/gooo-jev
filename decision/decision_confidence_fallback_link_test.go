package decision

import "testing"

func TestLinkDecisionConfidenceFallbackToReverseObservation(t *testing.T) {
	observation := DecisionConfidenceFallbackObservation{
		PlanDigest:     "plan",
		Action:         DecisionConfidenceFallbackReview,
		Status:         DecisionConfidenceFallbackUnknown,
		NonAuthorizing: true,
	}
	observationDigest, err := Digest(observation)
	if err != nil {
		t.Fatal(err)
	}
	observation.ObservationDigest = observationDigest

	link, err := LinkDecisionConfidenceFallbackToReverseObservation(observation, ReverseObservation{})
	if err != nil {
		t.Fatal(err)
	}
	if !link.NonAuthorizing || link.FallbackObservationDigest != observation.ObservationDigest {
		t.Fatalf("unexpected execution link: %+v", link)
	}
}

func TestDecisionConfidenceFallbackExecutionLinkRejectsTampering(t *testing.T) {
	link := DecisionConfidenceFallbackExecutionLink{
		FallbackObservationDigest: "fallback",
		ReverseObservationDigest:  "reverse",
		NonAuthorizing:            true,
		LinkDigest:                "tampered",
	}
	if err := link.Validate(); err == nil {
		t.Fatal("expected tampered execution link to be rejected")
	}
}
