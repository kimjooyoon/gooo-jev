package decision

import "fmt"

// DecisionConfidenceAdmissionMetricObservation records one exact admission
// outcome as a one-hot count linked to its source evidence.
type DecisionConfidenceAdmissionMetricObservation struct {
	AdmissionDigest string
	EligibleCount  uint64
	IneligibleCount uint64
	UnknownCount   uint64
	NonAuthorizing bool
	MetricDigest   string
}

func ObserveDecisionConfidenceAdmissionMetric(
	evidence DecisionConfidenceExecutionAdmissionEvidence,
) (DecisionConfidenceAdmissionMetricObservation, error) {
	if err := evidence.Validate(); err != nil {
		return DecisionConfidenceAdmissionMetricObservation{}, fmt.Errorf("validate admission evidence: %w", err)
	}
	metric := DecisionConfidenceAdmissionMetricObservation{
		AdmissionDigest: evidence.AdmissionDigest,
		NonAuthorizing:  true,
	}
	switch evidence.Status {
	case DecisionConfidenceExecutionAdmissionEligible:
		metric.EligibleCount = 1
	case DecisionConfidenceExecutionAdmissionIneligible:
		metric.IneligibleCount = 1
	case DecisionConfidenceExecutionAdmissionUnknown:
		metric.UnknownCount = 1
	default:
		return DecisionConfidenceAdmissionMetricObservation{}, fmt.Errorf("unsupported admission status %q", evidence.Status)
	}
	digest, err := Digest(metric)
	if err != nil {
		return DecisionConfidenceAdmissionMetricObservation{}, fmt.Errorf("digest admission metric: %w", err)
	}
	metric.MetricDigest = digest
	if err := metric.Validate(); err != nil {
		return DecisionConfidenceAdmissionMetricObservation{}, fmt.Errorf("validate admission metric: %w", err)
	}
	return metric, nil
}

func (m DecisionConfidenceAdmissionMetricObservation) Validate() error {
	if m.AdmissionDigest == "" {
		return fmt.Errorf("admission digest is required")
	}
	if !m.NonAuthorizing {
		return fmt.Errorf("admission metric must be non-authorizing")
	}
	if m.EligibleCount+m.IneligibleCount+m.UnknownCount != 1 {
		return fmt.Errorf("admission metric must contain exactly one outcome")
	}
	if m.MetricDigest == "" {
		return fmt.Errorf("metric digest is required")
	}
	withoutDigest := m
	withoutDigest.MetricDigest = ""
	digest, err := Digest(withoutDigest)
	if err != nil {
		return fmt.Errorf("digest admission metric: %w", err)
	}
	if digest != m.MetricDigest {
		return fmt.Errorf("metric digest mismatch")
	}
	return nil
}
