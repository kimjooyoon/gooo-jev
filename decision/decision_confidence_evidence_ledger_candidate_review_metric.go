package decision

// DecisionConfidenceEvidenceLedgerCandidateReviewMetric is a one-hot observation
// of review readiness, not a measure of improvement or authorization.
type DecisionConfidenceEvidenceLedgerCandidateReviewMetric struct {
	SignalDigest   string `json:"signal_digest"`
	ReadyCount     uint64 `json:"ready_count"`
	RejectedCount  uint64 `json:"rejected_count"`
	HoldCount      uint64 `json:"hold_count"`
	Status         string `json:"status"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

func ObserveDecisionConfidenceEvidenceLedgerCandidateReviewMetric(signal DecisionConfidenceEvidenceLedgerCandidateReviewSignal) DecisionConfidenceEvidenceLedgerCandidateReviewMetric {
	metric := DecisionConfidenceEvidenceLedgerCandidateReviewMetric{
		SignalDigest:   signal.ReviewOutcomeDigest,
		Status:         "hold",
		NonAuthorizing: true,
	}
	switch signal.Status {
	case "ready-for-candidate-review":
		metric.ReadyCount = 1
		metric.Status = "ready-for-candidate-review"
	case "rejected":
		metric.RejectedCount = 1
		metric.Status = "rejected"
	}
	if metric.SignalDigest == "" {
		metric.Status = "hold"
		metric.ReadyCount = 0
		metric.RejectedCount = 0
		metric.HoldCount = 1
	}
	if metric.Status == "hold" && metric.HoldCount == 0 {
		metric.HoldCount = 1
	}
	return metric
}
