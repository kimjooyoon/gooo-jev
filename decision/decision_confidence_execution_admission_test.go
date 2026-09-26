package decision

import "testing"

func TestEvaluateDecisionConfidenceExecutionAdmissionFailsClosed(t *testing.T) {
	_, err := EvaluateDecisionConfidenceExecutionAdmission(
		DecisionConfidenceFallbackPromotionEvidence{},
		DecisionConfidenceWorkloadIdentityReference{},
	)
	if err == nil {
		t.Fatal("expected invalid admission inputs to be rejected")
	}
}

func TestDecisionConfidenceExecutionAdmissionRejectsTampering(t *testing.T) {
	evidence := DecisionConfidenceExecutionAdmissionEvidence{
		PromotionDigest:        "promotion",
		IdentityReferenceDigest: "identity",
		Status:                 DecisionConfidenceExecutionAdmissionEligible,
		NonAuthorizing:         true,
		AdmissionDigest:        "tampered",
	}
	if err := evidence.Validate(); err == nil {
		t.Fatal("expected tampered admission evidence to be rejected")
	}
}
