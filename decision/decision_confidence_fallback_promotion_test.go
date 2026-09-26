package decision

import "testing"

func TestEvaluateDecisionConfidenceFallbackPromotionFailsClosed(t *testing.T) {
	_, err := EvaluateDecisionConfidenceFallbackPromotion(
		DecisionConfidenceFallbackObservation{},
		DecisionConfidenceFallbackExecutionLink{},
		ImprovementReviewReceipt{},
	)
	if err == nil {
		t.Fatal("expected invalid promotion inputs to be rejected")
	}
}

func TestDecisionConfidenceFallbackPromotionEvidenceRejectsTampering(t *testing.T) {
	evidence := DecisionConfidenceFallbackPromotionEvidence{
		ObservationDigest: "observation",
		LinkDigest:        "link",
		ReviewDigest:      "review",
		Status:            DecisionConfidenceFallbackPromotionEligible,
		NonAuthorizing:    true,
		PromotionDigest:   "tampered",
	}
	if err := evidence.Validate(); err == nil {
		t.Fatal("expected tampered promotion evidence to be rejected")
	}
}
