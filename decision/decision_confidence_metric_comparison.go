package decision

import "fmt"

type DecisionConfidenceMetricComparisonStatus string

const (
	DecisionConfidenceMetricImproved      DecisionConfidenceMetricComparisonStatus = "improved"
	DecisionConfidenceMetricRegressed     DecisionConfidenceMetricComparisonStatus = "regressed"
	DecisionConfidenceMetricUnchanged     DecisionConfidenceMetricComparisonStatus = "unchanged"
	DecisionConfidenceMetricInconclusive  DecisionConfidenceMetricComparisonStatus = "inconclusive"
	DecisionConfidenceMetricUnknown       DecisionConfidenceMetricComparisonStatus = "unknown"
)

// DecisionConfidenceMetricComparison compares two exact observations without
// treating an inconclusive or unknown transition as improvement.
type DecisionConfidenceMetricComparison struct {
	BaselineMetricDigest string
	CurrentMetricDigest  string
	Status               DecisionConfidenceMetricComparisonStatus
	NonAuthorizing       bool
	ComparisonDigest     string
}

func CompareDecisionConfidenceAdmissionMetrics(
	baseline DecisionConfidenceAdmissionMetricObservation,
	current DecisionConfidenceAdmissionMetricObservation,
) (DecisionConfidenceMetricComparison, error) {
	if err := baseline.Validate(); err != nil {
		return DecisionConfidenceMetricComparison{}, fmt.Errorf("validate baseline metric: %w", err)
	}
	if err := current.Validate(); err != nil {
		return DecisionConfidenceMetricComparison{}, fmt.Errorf("validate current metric: %w", err)
	}
	comparison := DecisionConfidenceMetricComparison{
		BaselineMetricDigest: baseline.MetricDigest,
		CurrentMetricDigest:  current.MetricDigest,
		NonAuthorizing:      true,
	}
	baselineStatus := metricOutcome(baseline)
	currentStatus := metricOutcome(current)
	switch {
	case baselineStatus == "unknown" || currentStatus == "unknown":
		comparison.Status = DecisionConfidenceMetricUnknown
	case baselineStatus == currentStatus:
		comparison.Status = DecisionConfidenceMetricUnchanged
	case currentStatus == "eligible" && baselineStatus != "eligible":
		comparison.Status = DecisionConfidenceMetricImproved
	case baselineStatus == "eligible" && currentStatus != "eligible":
		comparison.Status = DecisionConfidenceMetricRegressed
	default:
		comparison.Status = DecisionConfidenceMetricInconclusive
	}
	digest, err := Digest(comparison)
	if err != nil {
		return DecisionConfidenceMetricComparison{}, fmt.Errorf("digest metric comparison: %w", err)
	}
	comparison.ComparisonDigest = digest
	if err := comparison.Validate(); err != nil {
		return DecisionConfidenceMetricComparison{}, fmt.Errorf("validate metric comparison: %w", err)
	}
	return comparison, nil
}

func metricOutcome(metric DecisionConfidenceAdmissionMetricObservation) string {
	switch {
	case metric.EligibleCount == 1:
		return "eligible"
	case metric.IneligibleCount == 1:
		return "ineligible"
	default:
		return "unknown"
	}
}

func (c DecisionConfidenceMetricComparison) Validate() error {
	if c.BaselineMetricDigest == "" || c.CurrentMetricDigest == "" {
		return fmt.Errorf("metric digests are required")
	}
	if c.Status != DecisionConfidenceMetricImproved && c.Status != DecisionConfidenceMetricRegressed && c.Status != DecisionConfidenceMetricUnchanged && c.Status != DecisionConfidenceMetricInconclusive && c.Status != DecisionConfidenceMetricUnknown {
		return fmt.Errorf("unsupported comparison status %q", c.Status)
	}
	if !c.NonAuthorizing {
		return fmt.Errorf("metric comparison must be non-authorizing")
	}
	if c.ComparisonDigest == "" {
		return fmt.Errorf("comparison digest is required")
	}
	withoutDigest := c
	withoutDigest.ComparisonDigest = ""
	digest, err := Digest(withoutDigest)
	if err != nil {
		return fmt.Errorf("digest metric comparison: %w", err)
	}
	if digest != c.ComparisonDigest {
		return fmt.Errorf("comparison digest mismatch")
	}
	return nil
}
