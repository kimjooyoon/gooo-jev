package decision

import "fmt"

type DecisionConfidenceProposalAdmissionStatus string

const (
	DecisionConfidenceProposalAdmissionReady    DecisionConfidenceProposalAdmissionStatus = "ready-for-external-authorization"
	DecisionConfidenceProposalAdmissionNotReady DecisionConfidenceProposalAdmissionStatus = "not-ready"
	DecisionConfidenceProposalAdmissionUnknown  DecisionConfidenceProposalAdmissionStatus = "unknown"
)

// DecisionConfidenceProposalAdmissionBinding joins proposal, admission, and
// workload identity evidence without granting execution.
type DecisionConfidenceProposalAdmissionBinding struct {
	ProposalDigest         string
	AdmissionDigest        string
	IdentityReferenceDigest string
	Status                 DecisionConfidenceProposalAdmissionStatus
	NonAuthorizing         bool
	BindingDigest          string
}

func BindDecisionConfidenceProposalAdmission(
	proposal DecisionConfidenceImprovementProposal,
	admission DecisionConfidenceExecutionAdmissionEvidence,
	identity DecisionConfidenceWorkloadIdentityReference,
) (DecisionConfidenceProposalAdmissionBinding, error) {
	if err := proposal.Validate(); err != nil {
		return DecisionConfidenceProposalAdmissionBinding{}, fmt.Errorf("validate improvement proposal: %w", err)
	}
	if err := admission.Validate(); err != nil {
		return DecisionConfidenceProposalAdmissionBinding{}, fmt.Errorf("validate admission evidence: %w", err)
	}
	if err := identity.Validate(); err != nil {
		return DecisionConfidenceProposalAdmissionBinding{}, fmt.Errorf("validate workload identity reference: %w", err)
	}
	if admission.IdentityReferenceDigest != identity.ReferenceDigest {
		return DecisionConfidenceProposalAdmissionBinding{}, fmt.Errorf("admission does not reference workload identity")
	}
	binding := DecisionConfidenceProposalAdmissionBinding{
		ProposalDigest:          proposal.ProposalDigest,
		AdmissionDigest:         admission.AdmissionDigest,
		IdentityReferenceDigest: identity.ReferenceDigest,
		NonAuthorizing:          true,
	}
	switch admission.Status {
	case DecisionConfidenceExecutionAdmissionEligible:
		binding.Status = DecisionConfidenceProposalAdmissionReady
	case DecisionConfidenceExecutionAdmissionIneligible:
		binding.Status = DecisionConfidenceProposalAdmissionNotReady
	case DecisionConfidenceExecutionAdmissionUnknown:
		binding.Status = DecisionConfidenceProposalAdmissionUnknown
	default:
		return DecisionConfidenceProposalAdmissionBinding{}, fmt.Errorf("unsupported admission status %q", admission.Status)
	}
	digest, err := Digest(binding)
	if err != nil {
		return DecisionConfidenceProposalAdmissionBinding{}, fmt.Errorf("digest proposal admission binding: %w", err)
	}
	binding.BindingDigest = digest
	if err := binding.Validate(); err != nil {
		return DecisionConfidenceProposalAdmissionBinding{}, fmt.Errorf("validate proposal admission binding: %w", err)
	}
	return binding, nil
}

func (b DecisionConfidenceProposalAdmissionBinding) Validate() error {
	if b.ProposalDigest == "" || b.AdmissionDigest == "" || b.IdentityReferenceDigest == "" {
		return fmt.Errorf("binding digests are required")
	}
	if b.Status != DecisionConfidenceProposalAdmissionReady && b.Status != DecisionConfidenceProposalAdmissionNotReady && b.Status != DecisionConfidenceProposalAdmissionUnknown {
		return fmt.Errorf("unsupported binding status %q", b.Status)
	}
	if !b.NonAuthorizing {
		return fmt.Errorf("binding must be non-authorizing")
	}
	if b.BindingDigest == "" {
		return fmt.Errorf("binding digest is required")
	}
	withoutDigest := b
	withoutDigest.BindingDigest = ""
	digest, err := Digest(withoutDigest)
	if err != nil {
		return fmt.Errorf("digest proposal admission binding: %w", err)
	}
	if digest != b.BindingDigest {
		return fmt.Errorf("binding digest mismatch")
	}
	return nil
}
