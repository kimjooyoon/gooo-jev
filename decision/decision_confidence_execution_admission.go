package decision

import "fmt"

type DecisionConfidenceExecutionAdmissionStatus string

const (
	DecisionConfidenceExecutionAdmissionEligible   DecisionConfidenceExecutionAdmissionStatus = "eligible"
	DecisionConfidenceExecutionAdmissionIneligible DecisionConfidenceExecutionAdmissionStatus = "ineligible"
	DecisionConfidenceExecutionAdmissionUnknown    DecisionConfidenceExecutionAdmissionStatus = "unknown"
)

// DecisionConfidenceExecutionAdmissionEvidence records whether the existing
// evidence is sufficient for a later, separately authorized execution step.
type DecisionConfidenceExecutionAdmissionEvidence struct {
	PromotionDigest       string
	IdentityReferenceDigest string
	Status                DecisionConfidenceExecutionAdmissionStatus
	NonAuthorizing        bool
	AdmissionDigest       string
}

func EvaluateDecisionConfidenceExecutionAdmission(
	promotion DecisionConfidenceFallbackPromotionEvidence,
	identity DecisionConfidenceWorkloadIdentityReference,
) (DecisionConfidenceExecutionAdmissionEvidence, error) {
	if err := promotion.Validate(); err != nil {
		return DecisionConfidenceExecutionAdmissionEvidence{}, fmt.Errorf("validate promotion evidence: %w", err)
	}
	if err := identity.Validate(); err != nil {
		return DecisionConfidenceExecutionAdmissionEvidence{}, fmt.Errorf("validate workload identity reference: %w", err)
	}
	if identity.PromotionDigest != promotion.PromotionDigest || identity.PromotionStatus != promotion.Status {
		return DecisionConfidenceExecutionAdmissionEvidence{}, fmt.Errorf("identity reference does not match promotion evidence")
	}
	evidence := DecisionConfidenceExecutionAdmissionEvidence{
		PromotionDigest:        promotion.PromotionDigest,
		IdentityReferenceDigest: identity.ReferenceDigest,
		NonAuthorizing:         true,
	}
	switch {
	case promotion.Status == DecisionConfidenceFallbackPromotionUnknown:
		evidence.Status = DecisionConfidenceExecutionAdmissionUnknown
	case promotion.Status == DecisionConfidenceFallbackPromotionEligible:
		evidence.Status = DecisionConfidenceExecutionAdmissionEligible
	default:
		evidence.Status = DecisionConfidenceExecutionAdmissionIneligible
	}
	digest, err := Digest(evidence)
	if err != nil {
		return DecisionConfidenceExecutionAdmissionEvidence{}, fmt.Errorf("digest admission evidence: %w", err)
	}
	evidence.AdmissionDigest = digest
	if err := evidence.Validate(); err != nil {
		return DecisionConfidenceExecutionAdmissionEvidence{}, fmt.Errorf("validate admission evidence: %w", err)
	}
	return evidence, nil
}

func (e DecisionConfidenceExecutionAdmissionEvidence) Validate() error {
	if e.PromotionDigest == "" || e.IdentityReferenceDigest == "" {
		return fmt.Errorf("admission evidence digests are required")
	}
	if e.Status != DecisionConfidenceExecutionAdmissionEligible && e.Status != DecisionConfidenceExecutionAdmissionIneligible && e.Status != DecisionConfidenceExecutionAdmissionUnknown {
		return fmt.Errorf("unsupported admission status %q", e.Status)
	}
	if !e.NonAuthorizing {
		return fmt.Errorf("admission evidence must be non-authorizing")
	}
	if e.AdmissionDigest == "" {
		return fmt.Errorf("admission digest is required")
	}
	withoutDigest := e
	withoutDigest.AdmissionDigest = ""
	digest, err := Digest(withoutDigest)
	if err != nil {
		return fmt.Errorf("digest admission evidence: %w", err)
	}
	if digest != e.AdmissionDigest {
		return fmt.Errorf("admission digest mismatch")
	}
	return nil
}
