package decision

import "strings"

// ExecutionEnvelopeJEVImprovementReviewGateInput carries immutable cycle,
// candidate, and feedback evidence into the next review boundary.
type ExecutionEnvelopeJEVImprovementReviewGateInput struct {
	Binding        ExecutionEnvelopeJEVImprovementCycleRevisionFeedbackBinding
	NonAuthorizing bool
}

// ExecutionEnvelopeJEVImprovementReviewGateDecision records what the evidence
// permits a reviewer to consider. It never selects, applies, or authorizes a
// revision.
type ExecutionEnvelopeJEVImprovementReviewGateDecision struct {
	Status                     string
	FeedbackStatus             string
	RevisionCandidateDigest    string
	BindingDigest              string
	DecisionDigest             string
	MissingStage               string
	NonExecuting               bool
	NonAuthorizing             bool
}

func digestExecutionEnvelopeJEVImprovementReviewGateDecision(decision ExecutionEnvelopeJEVImprovementReviewGateDecision) (string, error) {
	return Digest(struct {
		Status                  string
		FeedbackStatus          string
		RevisionCandidateDigest string
		BindingDigest           string
		MissingStage            string
	}{
		Status:                  decision.Status,
		FeedbackStatus:          decision.FeedbackStatus,
		RevisionCandidateDigest: decision.RevisionCandidateDigest,
		BindingDigest:           decision.BindingDigest,
		MissingStage:            decision.MissingStage,
	})
}

// EvaluateExecutionEnvelopeJEVImprovementReviewGate maps observed feedback to
// a review disposition while preserving the evidence boundary.
func EvaluateExecutionEnvelopeJEVImprovementReviewGate(input ExecutionEnvelopeJEVImprovementReviewGateInput) ExecutionEnvelopeJEVImprovementReviewGateDecision {
	output := ExecutionEnvelopeJEVImprovementReviewGateDecision{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Binding.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return finalizeExecutionEnvelopeJEVImprovementReviewGateDecision(output)
	}
	if !input.Binding.NonExecuting {
		output.MissingStage = "execution-boundary"
		return finalizeExecutionEnvelopeJEVImprovementReviewGateDecision(output)
	}
	if input.Binding.Status != "bound" || strings.TrimSpace(input.Binding.BindingDigest) == "" {
		output.MissingStage = input.Binding.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "cycle-revision-feedback-binding"
		}
		return finalizeExecutionEnvelopeJEVImprovementReviewGateDecision(output)
	}
	if strings.TrimSpace(input.Binding.RevisionCandidateDigest) == "" || strings.TrimSpace(input.Binding.RevisionCandidateEvidenceDigest) == "" {
		output.MissingStage = "revision-candidate"
		return finalizeExecutionEnvelopeJEVImprovementReviewGateDecision(output)
	}
	if strings.TrimSpace(input.Binding.FeedbackEvidenceDigest) == "" {
		output.MissingStage = "feedback-aggregation"
		return finalizeExecutionEnvelopeJEVImprovementReviewGateDecision(output)
	}
	output.FeedbackStatus = input.Binding.FeedbackStatus
	output.RevisionCandidateDigest = input.Binding.RevisionCandidateDigest
	output.BindingDigest = input.Binding.BindingDigest
	switch input.Binding.FeedbackStatus {
	case "stable-for-review":
		output.Status = "review-eligible"
	case "needs-revision":
		output.Status = "revision-required"
	case "hold":
		output.Status = "review-hold"
	default:
		output.FeedbackStatus = ""
		output.MissingStage = "feedback-status"
		return finalizeExecutionEnvelopeJEVImprovementReviewGateDecision(output)
	}
	return finalizeExecutionEnvelopeJEVImprovementReviewGateDecision(output)
}

func finalizeExecutionEnvelopeJEVImprovementReviewGateDecision(output ExecutionEnvelopeJEVImprovementReviewGateDecision) ExecutionEnvelopeJEVImprovementReviewGateDecision {
	digest, err := digestExecutionEnvelopeJEVImprovementReviewGateDecision(output)
	if err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "review-gate-decision-digest"
		output.DecisionDigest = ""
		return output
	}
	output.DecisionDigest = digest
	return output
}
