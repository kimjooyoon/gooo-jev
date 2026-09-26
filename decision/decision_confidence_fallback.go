package decision

import "fmt"

type DecisionConfidenceFallbackAction string

const (
	DecisionConfidenceFallbackReview   DecisionConfidenceFallbackAction = "review"
	DecisionConfidenceFallbackNoAction DecisionConfidenceFallbackAction = "no-action"
)

// DecisionConfidenceFallbackPlan records a deterministic, non-authorizing
// response to a confidence signal when execution is not yet permitted.
type DecisionConfidenceFallbackPlan struct {
	SignalDigest   string
	Action         DecisionConfidenceFallbackAction
	NonAuthorizing bool
	PlanDigest     string
}

func PlanDecisionConfidenceFallback(signal DecisionConfidenceImprovementSignal) (DecisionConfidenceFallbackPlan, error) {
	if err := signal.Validate(); err != nil {
		return DecisionConfidenceFallbackPlan{}, fmt.Errorf("validate confidence improvement signal: %w", err)
	}

	plan := DecisionConfidenceFallbackPlan{
		SignalDigest:   Digest(signal),
		NonAuthorizing: true,
	}
	switch string(signal.Status) {
	case "candidate", "unknown":
		plan.Action = DecisionConfidenceFallbackReview
	case "no-change":
		plan.Action = DecisionConfidenceFallbackNoAction
	default:
		return DecisionConfidenceFallbackPlan{}, fmt.Errorf("unsupported confidence improvement signal status %q", signal.Status)
	}
	if err := plan.Validate(); err != nil {
		return DecisionConfidenceFallbackPlan{}, fmt.Errorf("validate fallback plan: %w", err)
	}
	return plan, nil
}

func (p DecisionConfidenceFallbackPlan) Validate() error {
	if p.SignalDigest == "" {
		return fmt.Errorf("signal digest is required")
	}
	if p.Action != DecisionConfidenceFallbackReview && p.Action != DecisionConfidenceFallbackNoAction {
		return fmt.Errorf("unsupported fallback action %q", p.Action)
	}
	if !p.NonAuthorizing {
		return fmt.Errorf("fallback plan must be non-authorizing")
	}
	if p.PlanDigest == "" {
		return fmt.Errorf("plan digest is required")
	}
	withoutDigest := p
	withoutDigest.PlanDigest = ""
	if Digest(withoutDigest) != p.PlanDigest {
		return fmt.Errorf("plan digest mismatch")
	}
	return nil
}
