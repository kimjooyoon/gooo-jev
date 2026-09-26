package decision

// ExecutionEnvelopeOutcomeMetricInput observes an explicit review outcome as
// a one-hot metric without turning provenance into authorization.
type ExecutionEnvelopeOutcomeMetricInput struct {
	Status         string `json:"status"`
	OutcomeDigest  string `json:"outcome_digest"`
	MissingStage   string `json:"missing_stage"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

// ExecutionEnvelopeOutcomeMetric preserves accepted, rejected, and hold
// outcomes as separate evidence-linked counts.
type ExecutionEnvelopeOutcomeMetric struct {
	OutcomeDigest   string `json:"outcome_digest"`
	AcceptedCount   uint64 `json:"accepted_count"`
	RejectedCount   uint64 `json:"rejected_count"`
	HoldCount       uint64 `json:"hold_count"`
	Status          string `json:"status"`
	MissingStage    string `json:"missing_stage"`
	NonAuthorizing  bool   `json:"non_authorizing"`
}

// ObserveExecutionEnvelopeOutcomeMetric fails closed when outcome evidence is
// absent or the input claims authorization.
func ObserveExecutionEnvelopeOutcomeMetric(input ExecutionEnvelopeOutcomeMetricInput) ExecutionEnvelopeOutcomeMetric {
	metric := ExecutionEnvelopeOutcomeMetric{Status: "UNKNOWN", HoldCount: 1, NonAuthorizing: true}
	if !input.NonAuthorizing {
		metric.NonAuthorizing = false
		metric.MissingStage = "authorization-boundary"
		return metric
	}
	metric.OutcomeDigest = input.OutcomeDigest
	if input.OutcomeDigest == "" {
		metric.MissingStage = input.MissingStage
		if metric.MissingStage == "" {
			metric.MissingStage = "review-outcome"
		}
		return metric
	}
	switch input.Status {
	case "accepted":
		metric.AcceptedCount = 1
		metric.HoldCount = 0
		metric.Status = "accepted"
	case "rejected":
		metric.RejectedCount = 1
		metric.HoldCount = 0
		metric.Status = "rejected"
	case "hold", "UNKNOWN":
		metric.Status = "hold"
	default:
		metric.MissingStage = "review-outcome"
	}
	return metric
}
