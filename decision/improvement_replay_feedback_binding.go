package decision

// ImprovementReplayFeedbackBindingInput connects a replay result to a
// validated, non-authorizing feedback observation.
type ImprovementReplayFeedbackBindingInput struct {
	ReplayBindingStatus string
	CycleDigest         string
	ReplayDigest        string
	CandidateDigest     string
	Feedback            FeedbackObservation
	NonAuthorizing      bool
}

// ImprovementReplayFeedbackBinding preserves replay provenance and feedback
// classification without granting execution or authorization.
type ImprovementReplayFeedbackBinding struct {
	Status          string
	FeedbackKind    FeedbackKind
	MetricName      string
	FeedbackDigest  string
	CycleDigest     string
	ReplayDigest    string
	CandidateDigest string
	EvidenceDigest  string
	MissingStage    string
	NonExecuting    bool
	NonAuthorizing  bool
}

// BindImprovementReplayFeedback requires recorded replay evidence and valid
// feedback before recording confirmed or refuted feedback.
func BindImprovementReplayFeedback(input ImprovementReplayFeedbackBindingInput) ImprovementReplayFeedbackBinding {
	output := ImprovementReplayFeedbackBinding{
		Status:         "UNKNOWN",
		FeedbackKind:   FeedbackUnknown,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if input.ReplayBindingStatus == "diverged" {
		output.Status = "review"
		output.MissingStage = "replay-diverged"
		return output
	}
	if input.ReplayBindingStatus == "review" {
		output.Status = "review"
		output.MissingStage = "replay-evidence"
		return output
	}
	if input.ReplayBindingStatus != "replayed" {
		output.MissingStage = "replay-binding"
		return output
	}
	if input.CycleDigest == "" || input.ReplayDigest == "" || input.CandidateDigest == "" {
		output.MissingStage = "replay-evidence"
		return output
	}
	if err := input.Feedback.Validate(); err != nil {
		output.MissingStage = "feedback-observation"
		return output
	}
	evidenceDigest, err := Digest(struct {
		CycleDigest     string
		ReplayDigest    string
		CandidateDigest string
		FeedbackDigest  string
	}{
		CycleDigest:     input.CycleDigest,
		ReplayDigest:    input.ReplayDigest,
		CandidateDigest: input.CandidateDigest,
		FeedbackDigest:  input.Feedback.FeedbackDigest,
	})
	if err != nil {
		output.MissingStage = "feedback-evidence"
		return output
	}
	output.FeedbackKind = input.Feedback.Kind
	output.MetricName = input.Feedback.MetricName
	output.FeedbackDigest = input.Feedback.FeedbackDigest
	output.CycleDigest = input.CycleDigest
	output.ReplayDigest = input.ReplayDigest
	output.CandidateDigest = input.CandidateDigest
	output.EvidenceDigest = evidenceDigest
	switch input.Feedback.Kind {
	case FeedbackConfirmed, FeedbackRefuted:
		output.Status = "feedback-recorded"
	case FeedbackUnknown:
		output.Status = "review"
		output.MissingStage = "feedback-unknown"
	default:
		output.Status = "UNKNOWN"
		output.MissingStage = "feedback-kind"
	}
	return output
}
