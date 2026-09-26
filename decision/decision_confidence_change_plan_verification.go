package decision

import "fmt"

type DecisionConfidenceChangePlanVerificationStatus string

const (
	DecisionConfidenceChangePlanVerified DecisionConfidenceChangePlanVerificationStatus = "verified"
	DecisionConfidenceChangePlanRejected DecisionConfidenceChangePlanVerificationStatus = "rejected"
	DecisionConfidenceChangePlanUnknown  DecisionConfidenceChangePlanVerificationStatus = "unknown"
)

// DecisionConfidenceChangePlanVerification records an independent check of
// a proposed write set without applying the plan.
type DecisionConfidenceChangePlanVerification struct {
	ChangePlanDigest          string
	VerifierReference         string
	VerificationEvidenceDigest string
	Status                    DecisionConfidenceChangePlanVerificationStatus
	NonExecuting              bool
	VerificationDigest        string
}

func VerifyDecisionConfidenceChangePlan(
	plan DecisionConfidenceChangePlan,
	verifierReference string,
	evidenceDigest string,
	status DecisionConfidenceChangePlanVerificationStatus,
) (DecisionConfidenceChangePlanVerification, error) {
	if err := plan.Validate(); err != nil {
		return DecisionConfidenceChangePlanVerification{}, fmt.Errorf("validate change plan: %w", err)
	}
	if verifierReference == "" {
		return DecisionConfidenceChangePlanVerification{}, fmt.Errorf("verifier reference is required")
	}
	if status != DecisionConfidenceChangePlanVerified && status != DecisionConfidenceChangePlanRejected && status != DecisionConfidenceChangePlanUnknown {
		return DecisionConfidenceChangePlanVerification{}, fmt.Errorf("unsupported verification status %q", status)
	}
	if status != DecisionConfidenceChangePlanUnknown && evidenceDigest == "" {
		return DecisionConfidenceChangePlanVerification{}, fmt.Errorf("non-unknown verification requires evidence digest")
	}
	verification := DecisionConfidenceChangePlanVerification{
		ChangePlanDigest:           plan.ChangePlanDigest,
		VerifierReference:          verifierReference,
		VerificationEvidenceDigest: evidenceDigest,
		Status:                     status,
		NonExecuting:               true,
	}
	digest, err := Digest(verification)
	if err != nil {
		return DecisionConfidenceChangePlanVerification{}, fmt.Errorf("digest change plan verification: %w", err)
	}
	verification.VerificationDigest = digest
	if err := verification.Validate(); err != nil {
		return DecisionConfidenceChangePlanVerification{}, fmt.Errorf("validate change plan verification: %w", err)
	}
	return verification, nil
}

func (v DecisionConfidenceChangePlanVerification) Validate() error {
	if v.ChangePlanDigest == "" || v.VerifierReference == "" {
		return fmt.Errorf("change plan and verifier references are required")
	}
	if v.Status != DecisionConfidenceChangePlanVerified && v.Status != DecisionConfidenceChangePlanRejected && v.Status != DecisionConfidenceChangePlanUnknown {
		return fmt.Errorf("unsupported verification status %q", v.Status)
	}
	if v.Status != DecisionConfidenceChangePlanUnknown && v.VerificationEvidenceDigest == "" {
		return fmt.Errorf("non-unknown verification requires evidence digest")
	}
	if !v.NonExecuting {
		return fmt.Errorf("change plan verification must be non-executing")
	}
	if v.VerificationDigest == "" {
		return fmt.Errorf("verification digest is required")
	}
	withoutDigest := v
	withoutDigest.VerificationDigest = ""
	digest, err := Digest(withoutDigest)
	if err != nil {
		return fmt.Errorf("digest change plan verification: %w", err)
	}
	if digest != v.VerificationDigest {
		return fmt.Errorf("verification digest mismatch")
	}
	return nil
}
