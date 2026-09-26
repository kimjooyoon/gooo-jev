package decision

import "testing"

func TestVerifyDecisionConfidenceChangePlanFailsClosed(t *testing.T) {
	_, err := VerifyDecisionConfidenceChangePlan(DecisionConfidenceChangePlan{}, "verifier", "evidence", DecisionConfidenceChangePlanVerified)
	if err == nil {
		t.Fatal("expected invalid change plan to be rejected")
	}
}

func TestVerifyDecisionConfidenceChangePlanRejectsMissingEvidence(t *testing.T) {
	plan := DecisionConfidenceChangePlan{
		ProposalDigest:   "proposal",
		SourceReference:  "source",
		ChangeDigest:     "change",
		WriteSetDigest:   "write-set",
		NonExecuting:     true,
	}
	digest, err := Digest(plan)
	if err != nil {
		t.Fatal(err)
	}
	plan.ChangePlanDigest = digest
	if _, err := VerifyDecisionConfidenceChangePlan(plan, "verifier", "", DecisionConfidenceChangePlanVerified); err == nil {
		t.Fatal("expected verified result without evidence to be rejected")
	}
}

func TestDecisionConfidenceChangePlanVerificationRejectsTampering(t *testing.T) {
	verification := DecisionConfidenceChangePlanVerification{
		ChangePlanDigest:           "plan",
		VerifierReference:          "verifier",
		VerificationEvidenceDigest: "evidence",
		Status:                     DecisionConfidenceChangePlanVerified,
		NonExecuting:               true,
		VerificationDigest:         "tampered",
	}
	if err := verification.Validate(); err == nil {
		t.Fatal("expected tampered verification to be rejected")
	}
}
