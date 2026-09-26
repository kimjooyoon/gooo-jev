package decision

import "testing"

func TestBindDecisionConfidencePromotionRejectsInvalidSPIFFEID(t *testing.T) {
	promotion := DecisionConfidenceFallbackPromotionEvidence{
		ObservationDigest: "observation",
		LinkDigest:        "link",
		ReviewDigest:      "review",
		Status:            DecisionConfidenceFallbackPromotionEligible,
		NonAuthorizing:    true,
	}
	promotionDigest, err := Digest(promotion)
	if err != nil {
		t.Fatal(err)
	}
	promotion.PromotionDigest = promotionDigest

	_, err = BindDecisionConfidencePromotionToWorkloadIdentity(
		promotion,
		"https://example.invalid/workload",
		"svid",
		"bundle",
		"attestation",
	)
	if err == nil {
		t.Fatal("expected non-SPIFFE identity to be rejected")
	}
}

func TestDecisionConfidenceWorkloadIdentityReferenceRejectsTampering(t *testing.T) {
	reference := DecisionConfidenceWorkloadIdentityReference{
		PromotionDigest:   "promotion",
		PromotionStatus:   DecisionConfidenceFallbackPromotionUnknown,
		SPIFFEID:          "spiffe://example/workload",
		SVIDDigest:        "svid",
		TrustBundleDigest: "bundle",
		AttestationDigest: "attestation",
		NonAuthorizing:    true,
		ReferenceDigest:   "tampered",
	}
	if err := reference.Validate(); err == nil {
		t.Fatal("expected tampered workload identity reference to be rejected")
	}
}
