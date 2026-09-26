package decision

import "testing"

func TestBindDecisionConfidenceProposalAdmissionFailsClosed(t *testing.T) {
	_, err := BindDecisionConfidenceProposalAdmission(
		DecisionConfidenceImprovementProposal{},
		DecisionConfidenceExecutionAdmissionEvidence{},
		DecisionConfidenceWorkloadIdentityReference{},
	)
	if err == nil {
		t.Fatal("expected invalid proposal admission inputs to be rejected")
	}
}

func TestDecisionConfidenceProposalAdmissionBindingRejectsTampering(t *testing.T) {
	binding := DecisionConfidenceProposalAdmissionBinding{
		ProposalDigest:          "proposal",
		AdmissionDigest:         "admission",
		IdentityReferenceDigest: "identity",
		Status:                  DecisionConfidenceProposalAdmissionReady,
		NonAuthorizing:          true,
		BindingDigest:           "tampered",
	}
	if err := binding.Validate(); err == nil {
		t.Fatal("expected tampered binding to be rejected")
	}
}
