package decision

import "strings"

// ExecutionEnvelopeJEVCycleCandidateReplayFeedbackInput connects a cycle,
// its generated candidate, and a sealed reverse observation.
type ExecutionEnvelopeJEVCycleCandidateReplayFeedbackInput struct {
	Cycle          JEVImprovementCycleObservation
	Candidate      ExecutionEnvelopeJEVCycleRevisionCandidateGenerationBinding
	ReverseBinding ExecutionEnvelopeReverseObservationBinding
	MetricDigest   string
	NonAuthorizing bool
}

// DeriveJEVImprovementReplayFeedbackFromCycleCandidateReverseObservation
// admits feedback only when the candidate points to the same cycle evidence.
func DeriveJEVImprovementReplayFeedbackFromCycleCandidateReverseObservation(input ExecutionEnvelopeJEVCycleCandidateReplayFeedbackInput) JEVImprovementReplayFeedback {
	output := JEVImprovementReplayFeedback{
		Status: jevReplayFeedbackUnknown, NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Cycle.NonAuthorizing || !input.Candidate.NonAuthorizing || !input.ReverseBinding.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.Cycle.NonExecuting || !input.Candidate.NonExecuting {
		output.MissingStage = "execution-boundary"
		return output
	}
	if err := input.Cycle.Validate(); err != nil {
		output.MissingStage = input.Cycle.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "improvement-cycle-observation"
		}
		return output
	}
	if input.Candidate.Status != "bound" ||
		input.Candidate.CycleStatus != input.Cycle.Status ||
		strings.TrimSpace(input.Candidate.CycleEvidenceDigest) == "" ||
		input.Candidate.CycleEvidenceDigest != input.Cycle.EvidenceDigest ||
		strings.TrimSpace(input.Candidate.CandidateDigest) == "" ||
		strings.TrimSpace(input.Candidate.CandidateEvidenceDigest) == "" ||
		strings.TrimSpace(input.Candidate.BoundRevisionChangeDigest) == "" ||
		strings.TrimSpace(input.Candidate.BindingDigest) == "" {
		output.MissingStage = "candidate-cycle-binding"
		return output
	}
	return DeriveJEVImprovementReplayFeedbackFromCycleRevisionCandidate(ExecutionEnvelopeJEVCycleCandidateFeedbackAdapterInput{
		CandidateBinding: input.Candidate,
		ReverseBinding:   input.ReverseBinding,
		MetricDigest:     input.MetricDigest,
		NonAuthorizing:   true,
	})
}
