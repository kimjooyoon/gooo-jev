package decision

import "fmt"

type DecisionConfidenceFallbackObservationStatus string

const (
	DecisionConfidenceFallbackObserved DecisionConfidenceFallbackObservationStatus = "observed"
	DecisionConfidenceFallbackUnknown  DecisionConfidenceFallbackObservationStatus = "unknown"
)

// DecisionConfidenceFallbackObservation binds a fallback plan to what was
// actually observed without authorizing execution.
type DecisionConfidenceFallbackObservation struct {
	PlanDigest         string
	Action             DecisionConfidenceFallbackAction
	Status             DecisionConfidenceFallbackObservationStatus
	EvidenceDigest     string
	NonAuthorizing     bool
	ObservationDigest  string
}

func ObserveDecisionConfidenceFallback(
	plan DecisionConfidenceFallbackPlan,
	action DecisionConfidenceFallbackAction,
	status DecisionConfidenceFallbackObservationStatus,
	evidenceDigest string,
) (DecisionConfidenceFallbackObservation, error) {
	if err := plan.Validate(); err != nil {
		return DecisionConfidenceFallbackObservation{}, fmt.Errorf("validate fallback plan: %w", err)
	}
	if action != plan.Action {
		return DecisionConfidenceFallbackObservation{}, fmt.Errorf("observed action %q does not match planned action %q", action, plan.Action)
	}
	if status != DecisionConfidenceFallbackObserved && status != DecisionConfidenceFallbackUnknown {
		return DecisionConfidenceFallbackObservation{}, fmt.Errorf("unsupported observation status %q", status)
	}
	if status == DecisionConfidenceFallbackObserved && evidenceDigest == "" {
		return DecisionConfidenceFallbackObservation{}, fmt.Errorf("observed fallback requires evidence digest")
	}
	observation := DecisionConfidenceFallbackObservation{
		PlanDigest:     plan.PlanDigest,
		Action:         action,
		Status:         status,
		EvidenceDigest: evidenceDigest,
		NonAuthorizing: true,
	}
	digest, err := Digest(observation)
	if err != nil {
		return DecisionConfidenceFallbackObservation{}, fmt.Errorf("digest fallback observation: %w", err)
	}
	observation.ObservationDigest = digest
	if err := observation.Validate(); err != nil {
		return DecisionConfidenceFallbackObservation{}, fmt.Errorf("validate fallback observation: %w", err)
	}
	return observation, nil
}

func (o DecisionConfidenceFallbackObservation) Validate() error {
	if o.PlanDigest == "" {
		return fmt.Errorf("plan digest is required")
	}
	if o.Action != DecisionConfidenceFallbackReview && o.Action != DecisionConfidenceFallbackNoAction {
		return fmt.Errorf("unsupported fallback action %q", o.Action)
	}
	if o.Status != DecisionConfidenceFallbackObserved && o.Status != DecisionConfidenceFallbackUnknown {
		return fmt.Errorf("unsupported observation status %q", o.Status)
	}
	if o.Status == DecisionConfidenceFallbackObserved && o.EvidenceDigest == "" {
		return fmt.Errorf("observed fallback requires evidence digest")
	}
	if !o.NonAuthorizing {
		return fmt.Errorf("fallback observation must be non-authorizing")
	}
	if o.ObservationDigest == "" {
		return fmt.Errorf("observation digest is required")
	}
	withoutDigest := o
	withoutDigest.ObservationDigest = ""
	digest, err := Digest(withoutDigest)
	if err != nil {
		return fmt.Errorf("digest fallback observation: %w", err)
	}
	if digest != o.ObservationDigest {
		return fmt.Errorf("observation digest mismatch")
	}
	return nil
}
