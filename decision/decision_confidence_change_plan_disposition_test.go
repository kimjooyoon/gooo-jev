package decision

import "testing"

func TestDispositionDecisionConfidenceChangePlanFailsClosed(t *testing.T) {
	_, err := DispositionDecisionConfidenceChangePlan(
		DecisionConfidenceChangePlan{},
		DecisionConfidenceChangePlanVerification{},
	)
	if err == nil {
		t.Fatal("expected invalid plan disposition inputs to be rejected")
	}
}

func TestDecisionConfidenceChangePlanDispositionRejectsTampering(t *testing.T) {
	disposition := DecisionConfidenceChangePlanDisposition{
		ChangePlanDigest:   "plan",
		VerificationDigest: "verification",
		Status:             DecisionConfidenceChangePlanAbort,
		NonExecuting:       true,
		DispositionDigest:  "tampered",
	}
	if err := disposition.Validate(); err == nil {
		t.Fatal("expected tampered disposition to be rejected")
	}
}
