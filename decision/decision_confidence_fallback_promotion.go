package decision

import "fmt"

type DecisionConfidenceFallbackPromotionStatus string

const (
	DecisionConfidenceFallbackPromotionEligible   DecisionConfidenceFallbackPromotionStatus = "eligible"
	DecisionConfidenceFallbackPromotionIneligible DecisionConfidenceFallbackPromotionStatus = "ineligible"
	DecisionConfidenceFallbackPromotionUnknown    DecisionConfidenceFallbackPromotionStatus = "unknown"
)

// DecisionConfidenceFallbackPromotionEvidence evaluates whether a reviewed
// fallback observation is eligible for a later proposal, never for execution.
type DecisionConfidenceFallbackPromotionEvidence struct {
	ObservationDigest string
	LinkDigest        string
	ReviewDigest      string
	Status            DecisionConfidenceFallbackPromotionStatus
	NonAuthorizing    bool
	PromotionDigest   string
}

func EvaluateDecisionConfidenceFallbackPromotion(
	observation DecisionConfidenceFallbackObservation,
	link DecisionConfidenceFallbackExecutionLink,
	review ImprovementReviewReceipt,
) (DecisionConfidenceFallbackPromotionEvidence, error) {
	if err := observation.Validate(); err != nil {
		return DecisionConfidenceFallbackPromotionEvidence{}, fmt.Errorf("validate fallback observation: %w", err)
	}
	if err := link.Validate(); err != nil {
		return DecisionConfidenceFallbackPromotionEvidence{}, fmt.Errorf("validate fallback execution link: %w", err)
	}
	if link.FallbackObservationDigest != observation.ObservationDigest {
		return DecisionConfidenceFallbackPromotionEvidence{}, fmt.Errorf("execution link does not reference observation")
	}
	if err := review.Validate(); err != nil {
		return DecisionConfidenceFallbackPromotionEvidence{}, fmt.Errorf("validate improvement review: %w", err)
	}
	if review.ExecutionGranted {
		return DecisionConfidenceFallbackPromotionEvidence{}, fmt.Errorf("promotion evidence cannot carry execution grant")
	}
	reviewDigest, err := Digest(review)
	if err != nil {
		return DecisionConfidenceFallbackPromotionEvidence{}, fmt.Errorf("digest improvement review: %w", err)
	}
	evidence := DecisionConfidenceFallbackPromotionEvidence{
		ObservationDigest: observation.ObservationDigest,
		LinkDigest:        link.LinkDigest,
		ReviewDigest:      reviewDigest,
		NonAuthorizing:    true,
	}
	switch {
	case observation.Status == DecisionConfidenceFallbackUnknown:
		evidence.Status = DecisionConfidenceFallbackPromotionUnknown
	case observation.Action == DecisionConfidenceFallbackNoAction:
		evidence.Status = DecisionConfidenceFallbackPromotionIneligible
	case review.Decision == ReviewPassed:
		evidence.Status = DecisionConfidenceFallbackPromotionEligible
	case review.Decision == ReviewFailed:
		evidence.Status = DecisionConfidenceFallbackPromotionIneligible
	default:
		evidence.Status = DecisionConfidenceFallbackPromotionUnknown
	}
	evidenceDigest, err := Digest(evidence)
	if err != nil {
		return DecisionConfidenceFallbackPromotionEvidence{}, fmt.Errorf("digest promotion evidence: %w", err)
	}
	evidence.PromotionDigest = evidenceDigest
	if err := evidence.Validate(); err != nil {
		return DecisionConfidenceFallbackPromotionEvidence{}, fmt.Errorf("validate promotion evidence: %w", err)
	}
	return evidence, nil
}

func (e DecisionConfidenceFallbackPromotionEvidence) Validate() error {
	if e.ObservationDigest == "" || e.LinkDigest == "" || e.ReviewDigest == "" {
		return fmt.Errorf("promotion evidence digests are required")
	}
	if e.Status != DecisionConfidenceFallbackPromotionEligible && e.Status != DecisionConfidenceFallbackPromotionIneligible && e.Status != DecisionConfidenceFallbackPromotionUnknown {
		return fmt.Errorf("unsupported promotion status %q", e.Status)
	}
	if !e.NonAuthorizing {
		return fmt.Errorf("promotion evidence must be non-authorizing")
	}
	if e.PromotionDigest == "" {
		return fmt.Errorf("promotion digest is required")
	}
	withoutDigest := e
	withoutDigest.PromotionDigest = ""
	digest, err := Digest(withoutDigest)
	if err != nil {
		return fmt.Errorf("digest promotion evidence: %w", err)
	}
	if digest != e.PromotionDigest {
		return fmt.Errorf("promotion digest mismatch")
	}
	return nil
}
