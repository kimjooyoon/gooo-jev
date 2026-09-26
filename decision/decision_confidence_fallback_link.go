package decision

import "fmt"

// DecisionConfidenceFallbackExecutionLink binds decision-layer reverse
// observation to the existing execution-envelope observation by digest only.
type DecisionConfidenceFallbackExecutionLink struct {
	FallbackObservationDigest string
	ReverseObservationDigest  string
	NonAuthorizing            bool
	LinkDigest                string
}

func LinkDecisionConfidenceFallbackToReverseObservation(
	observation DecisionConfidenceFallbackObservation,
	reverseObservation ReverseObservation,
) (DecisionConfidenceFallbackExecutionLink, error) {
	if err := observation.Validate(); err != nil {
		return DecisionConfidenceFallbackExecutionLink{}, fmt.Errorf("validate fallback observation: %w", err)
	}
	reverseDigest, err := Digest(reverseObservation)
	if err != nil {
		return DecisionConfidenceFallbackExecutionLink{}, fmt.Errorf("digest reverse observation: %w", err)
	}
	link := DecisionConfidenceFallbackExecutionLink{
		FallbackObservationDigest: observation.ObservationDigest,
		ReverseObservationDigest:  reverseDigest,
		NonAuthorizing:            true,
	}
	linkDigest, err := Digest(link)
	if err != nil {
		return DecisionConfidenceFallbackExecutionLink{}, fmt.Errorf("digest execution link: %w", err)
	}
	link.LinkDigest = linkDigest
	if err := link.Validate(); err != nil {
		return DecisionConfidenceFallbackExecutionLink{}, fmt.Errorf("validate execution link: %w", err)
	}
	return link, nil
}

func (l DecisionConfidenceFallbackExecutionLink) Validate() error {
	if l.FallbackObservationDigest == "" {
		return fmt.Errorf("fallback observation digest is required")
	}
	if l.ReverseObservationDigest == "" {
		return fmt.Errorf("reverse observation digest is required")
	}
	if !l.NonAuthorizing {
		return fmt.Errorf("execution link must be non-authorizing")
	}
	if l.LinkDigest == "" {
		return fmt.Errorf("link digest is required")
	}
	withoutDigest := l
	withoutDigest.LinkDigest = ""
	digest, err := Digest(withoutDigest)
	if err != nil {
		return fmt.Errorf("digest execution link: %w", err)
	}
	if digest != l.LinkDigest {
		return fmt.Errorf("link digest mismatch")
	}
	return nil
}
