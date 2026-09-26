package decision

// ExecutionEnvelopeReverseObservationTransitionMetricInput contains counts
// of observed reverse-evidence transitions.
type ExecutionEnvelopeReverseObservationTransitionMetricInput struct {
	StableCount       uint64
	ChangedCount      uint64
	UnknownCount      uint64
	NonAuthorizing    bool
}

// ExecutionEnvelopeReverseObservationTransitionMetric records transition
// counts without inferring product improvement or regression.
type ExecutionEnvelopeReverseObservationTransitionMetric struct {
	Status            string
	Total             uint64
	StableCount       uint64
	ChangedCount      uint64
	UnknownCount      uint64
	NonAuthorizing    bool
}

// MeasureExecutionEnvelopeReverseObservationTransitionMetric preserves
// unknown observations as data rather than dropping them.
func MeasureExecutionEnvelopeReverseObservationTransitionMetric(input ExecutionEnvelopeReverseObservationTransitionMetricInput) ExecutionEnvelopeReverseObservationTransitionMetric {
	output := ExecutionEnvelopeReverseObservationTransitionMetric{
		Status: "UNKNOWN", NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		return output
	}
	output.Total = input.StableCount + input.ChangedCount + input.UnknownCount
	if output.Total == 0 {
		return output
	}
	output.StableCount = input.StableCount
	output.ChangedCount = input.ChangedCount
	output.UnknownCount = input.UnknownCount
	output.Status = "measured"
	if output.UnknownCount > 0 {
		output.Status = "measured-with-unknown"
	}
	return output
}