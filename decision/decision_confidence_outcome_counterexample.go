package decision

import "fmt"

// RecordDecisionConfidenceOutcomeCounterexample maps only unsafe or
// unresolved external outcomes into the existing review counterexample path.
func RecordDecisionConfidenceOutcomeCounterexample(
	comparison DecisionConfidenceChangeOutcomeComparison,
) (DecisionConfidenceImprovementCounterexample, error) {
	if err := comparison.Validate(); err != nil {
		return DecisionConfidenceImprovementCounterexample{}, fmt.Errorf("validate outcome comparison: %w", err)
	}
	var reason DecisionConfidenceMetricComparisonStatus
	switch comparison.Status {
	case DecisionConfidenceChangeOutcomeRolledBack:
		reason = DecisionConfidenceMetricRegressed
	case DecisionConfidenceChangeOutcomeChanged:
		reason = DecisionConfidenceMetricInconclusive
	case DecisionConfidenceChangeOutcomeUnknown:
		reason = DecisionConfidenceMetricUnknown
	case DecisionConfidenceChangeOutcomeApplied, DecisionConfidenceChangeOutcomeUnchanged:
		return DecisionConfidenceImprovementCounterexample{}, fmt.Errorf("outcome status %q is not a counterexample", comparison.Status)
	default:
		return DecisionConfidenceImprovementCounterexample{}, fmt.Errorf("unsupported outcome status %q", comparison.Status)
	}
	counterexample := DecisionConfidenceImprovementCounterexample{
		ComparisonDigest: comparison.ComparisonDigest,
		ReasonStatus:     reason,
		RequiresReview:   true,
		NonAuthorizing:   true,
	}
	digest, err := Digest(counterexample)
	if err != nil {
		return DecisionConfidenceImprovementCounterexample{}, fmt.Errorf("digest outcome counterexample: %w", err)
	}
	counterexample.CounterexampleDigest = digest
	if err := counterexample.Validate(); err != nil {
		return DecisionConfidenceImprovementCounterexample{}, fmt.Errorf("validate outcome counterexample: %w", err)
	}
	return counterexample, nil
}
