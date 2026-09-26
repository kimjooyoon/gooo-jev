package decision

import "fmt"

// DecisionConfidenceImprovementCounterexample preserves a comparison that
// cannot safely support improvement and routes it back to review.
type DecisionConfidenceImprovementCounterexample struct {
	ComparisonDigest string
	ReasonStatus    DecisionConfidenceMetricComparisonStatus
	RequiresReview  bool
	NonAuthorizing  bool
	CounterexampleDigest string
}

func RecordDecisionConfidenceComparisonCounterexample(
	comparison DecisionConfidenceMetricComparison,
) (DecisionConfidenceImprovementCounterexample, error) {
	if err := comparison.Validate(); err != nil {
		return DecisionConfidenceImprovementCounterexample{}, fmt.Errorf("validate metric comparison: %w", err)
	}
	if comparison.Status == DecisionConfidenceMetricImproved || comparison.Status == DecisionConfidenceMetricUnchanged {
		return DecisionConfidenceImprovementCounterexample{}, fmt.Errorf("comparison status %q is not a counterexample", comparison.Status)
	}
	counterexample := DecisionConfidenceImprovementCounterexample{
		ComparisonDigest: comparison.ComparisonDigest,
		ReasonStatus:    comparison.Status,
		RequiresReview:  true,
		NonAuthorizing:  true,
	}
	digest, err := Digest(counterexample)
	if err != nil {
		return DecisionConfidenceImprovementCounterexample{}, fmt.Errorf("digest counterexample: %w", err)
	}
	counterexample.CounterexampleDigest = digest
	if err := counterexample.Validate(); err != nil {
		return DecisionConfidenceImprovementCounterexample{}, fmt.Errorf("validate counterexample: %w", err)
	}
	return counterexample, nil
}

func (c DecisionConfidenceImprovementCounterexample) Validate() error {
	if c.ComparisonDigest == "" {
		return fmt.Errorf("comparison digest is required")
	}
	if c.ReasonStatus != DecisionConfidenceMetricRegressed && c.ReasonStatus != DecisionConfidenceMetricInconclusive && c.ReasonStatus != DecisionConfidenceMetricUnknown {
		return fmt.Errorf("unsupported counterexample reason %q", c.ReasonStatus)
	}
	if !c.RequiresReview || !c.NonAuthorizing {
		return fmt.Errorf("counterexample must require review and be non-authorizing")
	}
	if c.CounterexampleDigest == "" {
		return fmt.Errorf("counterexample digest is required")
	}
	withoutDigest := c
	withoutDigest.CounterexampleDigest = ""
	digest, err := Digest(withoutDigest)
	if err != nil {
		return fmt.Errorf("digest counterexample: %w", err)
	}
	if digest != c.CounterexampleDigest {
		return fmt.Errorf("counterexample digest mismatch")
	}
	return nil
}
