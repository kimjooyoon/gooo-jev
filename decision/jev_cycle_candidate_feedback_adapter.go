package decision

import "strings"

// ExecutionEnvelopeJEVCycleCandidateFeedbackAdapterInput connects cycle-bound
// candidate generation to a sealed reverse observation and metric identity.
type ExecutionEnvelopeJEVCycleCandidateFeedbackAdapterInput struct {
	CandidateBinding ExecutionEnvelopeJEVCycleRevisionCandidateGenerationBinding
	ReverseBinding   ExecutionEnvelopeReverseObservationBinding
	MetricDigest     string
	NonAuthorizing   bool
}

// DeriveJEVImprovementReplayFeedbackFromCycleRevisionCandidate returns the
// existing feedback type while preserving cycle evidence in the metric digest.
func DeriveJEVImprovementReplayFeedbackFromCycleRevisionCandidate(input ExecutionEnvelopeJEVCycleCandidateFeedbackAdapterInput) JEVImprovementReplayFeedback {
	output := JEVImprovementReplayFeedback{
		Status: jevReplayFeedbackUnknown, NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.CandidateBinding.NonAuthorizing || !input.ReverseBinding.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.CandidateBinding.NonExecuting {
		output.MissingStage = "execution-boundary"
		return output
	}
	if input.CandidateBinding.Status != "bound" || strings.TrimSpace(input.CandidateBinding.CandidateDigest) == "" || strings.TrimSpace(input.CandidateBinding.CandidateEvidenceDigest) == "" || strings.TrimSpace(input.CandidateBinding.CycleEvidenceDigest) == "" {
		output.MissingStage = input.CandidateBinding.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "cycle-revision-candidate"
		}
		return output
	}
	if err := input.ReverseBinding.Validate(); err != nil {
		output.MissingStage = "reverse-observation"
		return output
	}
	if input.ReverseBinding.Status != "reproduced" && input.ReverseBinding.Status != "counterexample" {
		output.MissingStage = input.ReverseBinding.FirstMismatch
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "reverse-observation-status"
		}
		return output
	}
	if strings.TrimSpace(input.MetricDigest) == "" {
		output.MissingStage = "metric"
		return output
	}
	boundMetricDigest, err := Digest(struct {
		MetricDigest        string
		CycleEvidenceDigest string
	}{
		MetricDigest:        input.MetricDigest,
		CycleEvidenceDigest: input.CandidateBinding.CycleEvidenceDigest,
	})
	if err != nil {
		output.MissingStage = "cycle-metric-binding"
		return output
	}
	output.CandidateDigest = input.CandidateBinding.CandidateDigest
	output.ReplayObservationDigest = input.ReverseBinding.ObservationDigest
	output.MetricDigest = boundMetricDigest
	switch input.ReverseBinding.Status {
	case "reproduced":
		output.Status = jevReplayFeedbackConfirmed
		output.FeedbackKind = jevReplayFeedbackConfirmed
	case "counterexample":
		output.Status = jevReplayFeedbackRefuted
		output.FeedbackKind = jevReplayFeedbackRefuted
	}
	output.FeedbackEvidenceDigest = digestJEVImprovementReplayFeedbackEvidence(output.ReplayObservationDigest, output.MetricDigest, output.FeedbackKind)
	output.EvidenceDigest = digestJEVImprovementReplayFeedback(output.Status, output.FeedbackKind, output.CandidateDigest, output.ReplayObservationDigest, output.MetricDigest, output.FeedbackEvidenceDigest)
	if err := output.Validate(); err != nil {
		output.Status = jevReplayFeedbackUnknown
		output.MissingStage = "feedback-evidence"
		output.EvidenceDigest = ""
	}
	return output
}
