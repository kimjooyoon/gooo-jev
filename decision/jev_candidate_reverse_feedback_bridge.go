package decision

import "strings"

// ExecutionEnvelopeJEVCandidateReverseObservationFeedbackInput converts a
// sealed candidate reverse observation into the existing feedback contract.
type ExecutionEnvelopeJEVCandidateReverseObservationFeedbackInput struct {
	GateBinding    ExecutionEnvelopeJEVCandidateReverseObservationGateBinding
	MetricDigest   string
	NonAuthorizing bool
}

// DeriveJEVImprovementReplayFeedbackFromCandidateReverseObservation preserves
// reproduced and counterexample evidence as confirmed/refuted feedback. It
// never treats either result as execution or authorization.
func DeriveJEVImprovementReplayFeedbackFromCandidateReverseObservation(input ExecutionEnvelopeJEVCandidateReverseObservationFeedbackInput) JEVImprovementReplayFeedback {
	output := JEVImprovementReplayFeedback{
		Status: jevReplayFeedbackUnknown, NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.GateBinding.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.GateBinding.NonExecuting {
		output.MissingStage = "execution-boundary"
		return output
	}
	if input.GateBinding.Status != "bound" || strings.TrimSpace(input.GateBinding.CandidateDigest) == "" || strings.TrimSpace(input.GateBinding.ReverseObservationDigest) == "" {
		output.MissingStage = input.GateBinding.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "candidate-reverse-observation"
		}
		return output
	}
	if strings.TrimSpace(input.MetricDigest) == "" {
		output.MissingStage = "metric"
		return output
	}
	output.CandidateDigest = input.GateBinding.CandidateDigest
	output.ReplayObservationDigest = input.GateBinding.ReverseObservationDigest
	output.MetricDigest = input.MetricDigest
	switch input.GateBinding.ReverseObservationStatus {
	case "reproduced":
		output.Status = jevReplayFeedbackConfirmed
		output.FeedbackKind = jevReplayFeedbackConfirmed
	case "counterexample":
		output.Status = jevReplayFeedbackRefuted
		output.FeedbackKind = jevReplayFeedbackRefuted
	default:
		output.MissingStage = "reverse-observation-status"
		return output
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
