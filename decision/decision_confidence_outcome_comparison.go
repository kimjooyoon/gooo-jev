package decision

import "fmt"

type DecisionConfidenceChangeOutcomeComparisonStatus string

const (
	DecisionConfidenceChangeOutcomeApplied    DecisionConfidenceChangeOutcomeComparisonStatus = "applied"
	DecisionConfidenceChangeOutcomeRolledBack DecisionConfidenceChangeOutcomeComparisonStatus = "rolled-back"
	DecisionConfidenceChangeOutcomeChanged    DecisionConfidenceChangeOutcomeComparisonStatus = "changed"
	DecisionConfidenceChangeOutcomeUnchanged  DecisionConfidenceChangeOutcomeComparisonStatus = "unchanged"
	DecisionConfidenceChangeOutcomeUnknown    DecisionConfidenceChangeOutcomeComparisonStatus = "unknown"
)

// DecisionConfidenceChangeOutcomeComparison compares execution outcomes
// without calling an applied result an improvement.
type DecisionConfidenceChangeOutcomeComparison struct {
	BaselineMetricDigest string
	CurrentMetricDigest  string
	Status               DecisionConfidenceChangeOutcomeComparisonStatus
	NonAuthorizing       bool
	ComparisonDigest     string
}

func CompareDecisionConfidenceChangeOutcomes(
	baseline DecisionConfidenceChangeOutcomeMetric,
	current DecisionConfidenceChangeOutcomeMetric,
) (DecisionConfidenceChangeOutcomeComparison, error) {
	if err := baseline.Validate(); err != nil {
		return DecisionConfidenceChangeOutcomeComparison{}, fmt.Errorf("validate baseline outcome metric: %w", err)
	}
	if err := current.Validate(); err != nil {
		return DecisionConfidenceChangeOutcomeComparison{}, fmt.Errorf("validate current outcome metric: %w", err)
	}
	comparison := DecisionConfidenceChangeOutcomeComparison{
		BaselineMetricDigest: baseline.MetricDigest,
		CurrentMetricDigest:  current.MetricDigest,
		NonAuthorizing:      true,
	}
	baselineOutcome := changeOutcome(baseline)
	currentOutcome := changeOutcome(current)
	switch {
	case baselineOutcome == "unknown" || currentOutcome == "unknown":
		comparison.Status = DecisionConfidenceChangeOutcomeUnknown
	case baselineOutcome == currentOutcome:
		comparison.Status = DecisionConfidenceChangeOutcomeUnchanged
	case currentOutcome == "applied":
		comparison.Status = DecisionConfidenceChangeOutcomeApplied
	case currentOutcome == "rolled-back":
		comparison.Status = DecisionConfidenceChangeOutcomeRolledBack
	default:
		comparison.Status = DecisionConfidenceChangeOutcomeChanged
	}
	digest, err := Digest(comparison)
	if err != nil {
		return DecisionConfidenceChangeOutcomeComparison{}, fmt.Errorf("digest outcome comparison: %w", err)
	}
	comparison.ComparisonDigest = digest
	if err := comparison.Validate(); err != nil {
		return DecisionConfidenceChangeOutcomeComparison{}, fmt.Errorf("validate outcome comparison: %w", err)
	}
	return comparison, nil
}

func changeOutcome(metric DecisionConfidenceChangeOutcomeMetric) string {
	switch {
	case metric.AppliedCount == 1:
		return "applied"
	case metric.AbortedCount == 1:
		return "aborted"
	case metric.RolledBackCount == 1:
		return "rolled-back"
	case metric.NotAppliedCount == 1:
		return "not-applied"
	default:
		return "unknown"
	}
}

func (c DecisionConfidenceChangeOutcomeComparison) Validate() error {
	if c.BaselineMetricDigest == "" || c.CurrentMetricDigest == "" {
		return fmt.Errorf("outcome metric digests are required")
	}
	if c.Status != DecisionConfidenceChangeOutcomeApplied && c.Status != DecisionConfidenceChangeOutcomeRolledBack && c.Status != DecisionConfidenceChangeOutcomeChanged && c.Status != DecisionConfidenceChangeOutcomeUnchanged && c.Status != DecisionConfidenceChangeOutcomeUnknown {
		return fmt.Errorf("unsupported outcome comparison status %q", c.Status)
	}
	if !c.NonAuthorizing {
		return fmt.Errorf("outcome comparison must be non-authorizing")
	}
	if c.ComparisonDigest == "" {
		return fmt.Errorf("comparison digest is required")
	}
	withoutDigest := c
	withoutDigest.ComparisonDigest = ""
	digest, err := Digest(withoutDigest)
	if err != nil {
		return fmt.Errorf("digest outcome comparison: %w", err)
	}
	if digest != c.ComparisonDigest {
		return fmt.Errorf("comparison digest mismatch")
	}
	return nil
}
