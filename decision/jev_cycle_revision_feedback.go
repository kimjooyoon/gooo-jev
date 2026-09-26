package decision

import "strings"

// ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackInput joins a validated
// plan improvement cycle with a revision candidate and observed feedback.
type ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackInput struct {
	CycleBinding      ExecutionEnvelopeJEVPlanImprovementCycleBinding
	RevisionCandidate JEVImprovementRevisionCandidate
	Feedback          JEVImprovementFeedbackAggregation
	NonAuthorizing    bool
}

// ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackBinding records the
// immutable evidence boundary without selecting, applying, or authorizing a
// revision.
type ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackBinding struct {
	Status                         string
	PlanID                         string
	CycleBindingDigest             string
	RevisionCandidateDigest        string
	RevisionCandidateEvidenceDigest string
	FeedbackStatus                string
	FeedbackEvidenceDigest        string
	BindingDigest                 string
	MissingStage                  string
	NonExecuting                  bool
	NonAuthorizing                bool
}

// BindExecutionEnvelopeJEVImprovementCycleToRevisionFeedback preserves the
// cycle outcome and feedback status as evidence for a later review boundary.
func BindExecutionEnvelopeJEVImprovementCycleToRevisionFeedback(input ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackInput) ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackBinding {
	output := ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackBinding{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.CycleBinding.NonAuthorizing || !input.RevisionCandidate.NonAuthorizing || !input.Feedback.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return output
	}
	if input.CycleBinding.Status != "bound" || strings.TrimSpace(input.CycleBinding.BindingDigest) == "" {
		output.MissingStage = input.CycleBinding.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "improvement-cycle"
		}
		return output
	}
	if err := input.RevisionCandidate.Validate(); err != nil {
		output.MissingStage = "revision-candidate"
		return output
	}
	if err := input.Feedback.Validate(); err != nil {
		output.MissingStage = "feedback-aggregation"
		return output
	}
	bindingDigest, err := Digest(struct {
		PlanID                          string
		CycleBindingDigest              string
		RevisionCandidateDigest         string
		RevisionCandidateEvidenceDigest string
		FeedbackStatus                  string
		FeedbackEvidenceDigest          string
	}{
		PlanID:                          input.CycleBinding.PlanID,
		CycleBindingDigest:              input.CycleBinding.BindingDigest,
		RevisionCandidateDigest:         input.RevisionCandidate.CandidateDigest,
		RevisionCandidateEvidenceDigest: input.RevisionCandidate.EvidenceDigest,
		FeedbackStatus:                  input.Feedback.Status,
		FeedbackEvidenceDigest:          input.Feedback.EvidenceDigest,
	})
	if err != nil {
		output.MissingStage = "cycle-revision-feedback-binding-digest"
		return output
	}
	output.Status = "bound"
	output.PlanID = input.CycleBinding.PlanID
	output.CycleBindingDigest = input.CycleBinding.BindingDigest
	output.RevisionCandidateDigest = input.RevisionCandidate.CandidateDigest
	output.RevisionCandidateEvidenceDigest = input.RevisionCandidate.EvidenceDigest
	output.FeedbackStatus = input.Feedback.Status
	output.FeedbackEvidenceDigest = input.Feedback.EvidenceDigest
	output.BindingDigest = bindingDigest
	return output
}
