package decision

import (
	"fmt"
	"strings"
)

// DecisionConfidenceWorkloadIdentityReference records SPIFFE/SPIRE-related
// identity evidence without turning it into an execution authorization.
type DecisionConfidenceWorkloadIdentityReference struct {
	PromotionDigest    string
	PromotionStatus    DecisionConfidenceFallbackPromotionStatus
	SPIFFEID           string
	SVIDDigest         string
	TrustBundleDigest  string
	AttestationDigest  string
	NonAuthorizing     bool
	ReferenceDigest    string
}

func BindDecisionConfidencePromotionToWorkloadIdentity(
	promotion DecisionConfidenceFallbackPromotionEvidence,
	spiffeID string,
	svidDigest string,
	trustBundleDigest string,
	attestationDigest string,
) (DecisionConfidenceWorkloadIdentityReference, error) {
	if err := promotion.Validate(); err != nil {
		return DecisionConfidenceWorkloadIdentityReference{}, fmt.Errorf("validate promotion evidence: %w", err)
	}
	if !strings.HasPrefix(spiffeID, "spiffe://") || len(spiffeID) <= len("spiffe://") {
		return DecisionConfidenceWorkloadIdentityReference{}, fmt.Errorf("invalid SPIFFE ID")
	}
	if svidDigest == "" || trustBundleDigest == "" || attestationDigest == "" {
		return DecisionConfidenceWorkloadIdentityReference{}, fmt.Errorf("SVID, trust bundle, and attestation digests are required")
	}
	reference := DecisionConfidenceWorkloadIdentityReference{
		PromotionDigest:   promotion.PromotionDigest,
		PromotionStatus:   promotion.Status,
		SPIFFEID:          spiffeID,
		SVIDDigest:        svidDigest,
		TrustBundleDigest: trustBundleDigest,
		AttestationDigest: attestationDigest,
		NonAuthorizing:    true,
	}
	digest, err := Digest(reference)
	if err != nil {
		return DecisionConfidenceWorkloadIdentityReference{}, fmt.Errorf("digest workload identity reference: %w", err)
	}
	reference.ReferenceDigest = digest
	if err := reference.Validate(); err != nil {
		return DecisionConfidenceWorkloadIdentityReference{}, fmt.Errorf("validate workload identity reference: %w", err)
	}
	return reference, nil
}

func (r DecisionConfidenceWorkloadIdentityReference) Validate() error {
	if r.PromotionDigest == "" || r.SPIFFEID == "" || r.SVIDDigest == "" || r.TrustBundleDigest == "" || r.AttestationDigest == "" {
		return fmt.Errorf("workload identity reference fields are required")
	}
	if !strings.HasPrefix(r.SPIFFEID, "spiffe://") || len(r.SPIFFEID) <= len("spiffe://") {
		return fmt.Errorf("invalid SPIFFE ID")
	}
	if r.PromotionStatus != DecisionConfidenceFallbackPromotionEligible && r.PromotionStatus != DecisionConfidenceFallbackPromotionIneligible && r.PromotionStatus != DecisionConfidenceFallbackPromotionUnknown {
		return fmt.Errorf("unsupported promotion status %q", r.PromotionStatus)
	}
	if !r.NonAuthorizing {
		return fmt.Errorf("workload identity reference must be non-authorizing")
	}
	if r.ReferenceDigest == "" {
		return fmt.Errorf("reference digest is required")
	}
	withoutDigest := r
	withoutDigest.ReferenceDigest = ""
	digest, err := Digest(withoutDigest)
	if err != nil {
		return fmt.Errorf("digest workload identity reference: %w", err)
	}
	if digest != r.ReferenceDigest {
		return fmt.Errorf("reference digest mismatch")
	}
	return nil
}
