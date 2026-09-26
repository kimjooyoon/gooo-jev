package decision

// ExecutionEnvelopeFeedbackMetricInput is the explicit review feedback state
// observed as a one-hot metric without claiming improvement.
type ExecutionEnvelopeFeedbackMetricInput struct {
	Status         string `json:"status"`
	FeedbackDigest string `json:"feedback_digest"`
	ReviewRequired bool   `json:"review_required"`
	MissingStage   string `json:"missing_stage"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

// ExecutionEnvelopeFeedbackMetric preserves the first missing stage and the
// non-authorizing boundary while exposing one-hot feedback counts.
type ExecutionEnvelopeFeedbackMetric struct {
	FeedbackDigest       string `json:"feedback_digest"`
	AnalysisReadyCount   uint64 `json:"analysis_ready_count"`
	RejectedCount        uint64 `json:"rejected_count"`
	HoldCount            uint64 `json:"hold_count"`
	ObservationOnlyCount uint64 `json:"observation_only_count"`
	Status               string `json:"status"`
	MissingStage         string `json:"missing_stage"`
	NonAuthorizing       bool   `json:"non_authorizing"`
}

// ObserveExecutionEnvelopeFeedbackMetric never turns missing feedback into a
// positive metric and never turns provenance into authorization.
func ObserveExecutionEnvelopeFeedbackMetric(input ExecutionEnvelopeFeedbackMetricInput) ExecutionEnvelopeFeedbackMetric {
	metric := ExecutionEnvelopeFeedbackMetric{Status: "UNKNOWN", HoldCount: 1, NonAuthorizing: true}
	if !input.NonAuthorizing {
		metric.NonAuthorizing = false
		metric.MissingStage = "authorization-boundary"
		return metric
	}
	metric.FeedbackDigest = input.FeedbackDigest
	if input.FeedbackDigest == "" {
		metric.MissingStage = input.MissingStage
		if metric.MissingStage == "" {
			metric.MissingStage = "feedback-evidence"
		}
		return metric
	}
	switch input.Status {
	case "reviewed":
		metric.AnalysisReadyCount = 1
		metric.HoldCount = 0
		metric.Status = "analysis-ready"
	case "review-rejected":
		metric.RejectedCount = 1
		metric.HoldCount = 0
		metric.Status = "rejected"
	case "observation-only":
		metric.ObservationOnlyCount = 1
		metric.HoldCount = 0
		metric.Status = "observation-only"
	case "hold", "UNKNOWN":
		metric.Status = "hold"
	default:
		metric.MissingStage = "feedback-status"
	}
	return metric
}
