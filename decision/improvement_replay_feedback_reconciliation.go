package decision

// ImprovementReplayFeedbackReconciliationInput combines a replay feedback
// binding with its validated non-authorizing feedback history.
type ImprovementReplayFeedbackReconciliationInput struct {
	ReplayFeedback  ImprovementReplayFeedbackBinding
	FeedbackHistory FeedbackHistory
	NonAuthorizing  bool
}

// ImprovementReplayFeedbackReconciliation records a bounded feedback
// disposition without claiming that a metric improved or authorizing execution.
type ImprovementReplayFeedbackReconciliation struct {
	Status                   string
	MetricName               string
	HistoryDigest            string
	ReplayEvidenceDigest     string
	LatestSummaryDigest      string
	EvidenceDigest           string
	MissingStage             string
	NonExecuting             bool
	NonAuthorizing           bool
}

// ReconcileImprovementReplayFeedback requires a recorded replay feedback
// binding and a valid history before classifying the latest observation window.
func ReconcileImprovementReplayFeedback(input ImprovementReplayFeedbackReconciliationInput) ImprovementReplayFeedbackReconciliation {
	output := ImprovementReplayFeedbackReconciliation{
		Status:         "UNKNOWN",
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.ReplayFeedback.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.ReplayFeedback.NonExecuting {
		output.NonExecuting = false
		output.MissingStage = "execution-boundary"
		return output
	}
	if input.ReplayFeedback.Status != "feedback-recorded" {
		output.Status = "review"
		output.MissingStage = "replay-feedback"
		return output
	}
	if input.ReplayFeedback.FeedbackKind != FeedbackConfirmed && input.ReplayFeedback.FeedbackKind != FeedbackRefuted {
		output.MissingStage = "replay-feedback-kind"
		return output
	}
	if input.ReplayFeedback.MetricName == "" ||
		input.ReplayFeedback.FeedbackDigest == "" ||
		input.ReplayFeedback.EvidenceDigest == "" {
		output.MissingStage = "replay-feedback-evidence"
		return output
	}
	if err := input.FeedbackHistory.Validate(); err != nil {
		output.MissingStage = "feedback-history"
		return output
	}
	if input.ReplayFeedback.MetricName != input.FeedbackHistory.MetricName {
		output.MissingStage = "feedback-metric"
		return output
	}
	latest := input.FeedbackHistory.Summaries[len(input.FeedbackHistory.Summaries)-1]
	evidenceDigest, err := Digest(struct {
		HistoryDigest        string
		ReplayEvidenceDigest string
		LatestSummaryDigest  string
		MetricName           string
	}{
		HistoryDigest:        input.FeedbackHistory.HistoryDigest,
		ReplayEvidenceDigest: input.ReplayFeedback.EvidenceDigest,
		LatestSummaryDigest:  latest.SummaryDigest,
		MetricName:           input.FeedbackHistory.MetricName,
	})
	if err != nil {
		output.MissingStage = "reconciliation-evidence"
		return output
	}
	output.MetricName = input.FeedbackHistory.MetricName
	output.HistoryDigest = input.FeedbackHistory.HistoryDigest
	output.ReplayEvidenceDigest = input.ReplayFeedback.EvidenceDigest
	output.LatestSummaryDigest = latest.SummaryDigest
	output.EvidenceDigest = evidenceDigest
	switch {
	case latest.UnknownCount > 0:
		output.Status = "review"
		output.MissingStage = "feedback-unknown"
	case latest.ConfirmedCount > latest.RefutedCount:
		output.Status = "confirmed"
	case latest.RefutedCount > latest.ConfirmedCount:
		output.Status = "refuted"
	default:
		output.Status = "review"
		output.MissingStage = "feedback-tie"
	}
	return output
}
