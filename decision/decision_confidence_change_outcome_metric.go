package decision

import "fmt"

// DecisionConfidenceChangeOutcomeMetric records one exact external outcome
// without treating it as a quality improvement claim.
type DecisionConfidenceChangeOutcomeMetric struct {
	ObservationDigest string
	AppliedCount     uint64
	AbortedCount     uint64
	RolledBackCount  uint64
	NotAppliedCount  uint64
	UnknownCount     uint64
	NonAuthorizing   bool
	MetricDigest     string
}

func ObserveDecisionConfidenceChangeOutcomeMetric(
	observation DecisionConfidenceChangePlanDispositionObservation,
) (DecisionConfidenceChangeOutcomeMetric, error) {
	if err := observation.Validate(); err != nil {
		return DecisionConfidenceChangeOutcomeMetric{}, fmt.Errorf("validate disposition observation: %w", err)
	}
	metric := DecisionConfidenceChangeOutcomeMetric{
		ObservationDigest: observation.ObservationDigest,
		NonAuthorizing:    true,
	}
	switch observation.Status {
	case DecisionConfidenceChangePlanObservedApplied:
		metric.AppliedCount = 1
	case DecisionConfidenceChangePlanObservedAborted:
		metric.AbortedCount = 1
	case DecisionConfidenceChangePlanObservedRolledBack:
		metric.RolledBackCount = 1
	case DecisionConfidenceChangePlanObservedNotApplied:
		metric.NotAppliedCount = 1
	case DecisionConfidenceChangePlanObservedUnknown:
		metric.UnknownCount = 1
	default:
		return DecisionConfidenceChangeOutcomeMetric{}, fmt.Errorf("unsupported disposition observation status %q", observation.Status)
	}
	digest, err := Digest(metric)
	if err != nil {
		return DecisionConfidenceChangeOutcomeMetric{}, fmt.Errorf("digest outcome metric: %w", err)
	}
	metric.MetricDigest = digest
	if err := metric.Validate(); err != nil {
		return DecisionConfidenceChangeOutcomeMetric{}, fmt.Errorf("validate outcome metric: %w", err)
	}
	return metric, nil
}

func (m DecisionConfidenceChangeOutcomeMetric) Validate() error {
	if m.ObservationDigest == "" {
		return fmt.Errorf("observation digest is required")
	}
	if !m.NonAuthorizing {
		return fmt.Errorf("outcome metric must be non-authorizing")
	}
	if m.AppliedCount+m.AbortedCount+m.RolledBackCount+m.NotAppliedCount+m.UnknownCount != 1 {
		return fmt.Errorf("outcome metric must contain exactly one outcome")
	}
	if m.MetricDigest == "" {
		return fmt.Errorf("metric digest is required")
	}
	withoutDigest := m
	withoutDigest.MetricDigest = ""
	digest, err := Digest(withoutDigest)
	if err != nil {
		return fmt.Errorf("digest outcome metric: %w", err)
	}
	if digest != m.MetricDigest {
		return fmt.Errorf("metric digest mismatch")
	}
	return nil
}
