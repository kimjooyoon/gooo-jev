package decision

import "strings"

// ExecutionEnvelopeJEVImprovementReviewObservationInput records a review
// outcome against an immutable review-gate decision.
type ExecutionEnvelopeJEVImprovementReviewObservationInput struct {
	GateDecision         ExecutionEnvelopeJEVImprovementReviewGateDecision
	ReviewOutcome        string
	ReviewEvidenceDigest string
	NonAuthorizing       bool
}

// ExecutionEnvelopeJEVImprovementReviewObservation is evidence for the next
// feedback aggregation. It never applies a candidate or authorizes a change.
type ExecutionEnvelopeJEVImprovementReviewObservation struct {
	Status               string
	ReviewOutcome        string
	GateDecisionDigest   string
	ReviewEvidenceDigest string
	ObservationDigest    string
	MissingStage         string
	NonExecuting         bool
	NonAuthorizing       bool
}

func digestExecutionEnvelopeJEVImprovementReviewObservation(observation ExecutionEnvelopeJEVImprovementReviewObservation) (string, error) {
	return Digest(struct {
		Status               string
		ReviewOutcome        string
		GateDecisionDigest   string
		ReviewEvidenceDigest string
		MissingStage         string
	}{
		Status:               observation.Status,
		ReviewOutcome:        observation.ReviewOutcome,
		GateDecisionDigest:   observation.GateDecisionDigest,
		ReviewEvidenceDigest: observation.ReviewEvidenceDigest,
		MissingStage:         observation.MissingStage,
	})
}

// ObserveExecutionEnvelopeJEVImprovementReview validates a review observation
// as feedback-ready evidence without treating it as an authorization.
func ObserveExecutionEnvelopeJEVImprovementReview(input ExecutionEnvelopeJEVImprovementReviewObservationInput) ExecutionEnvelopeJEVImprovementReviewObservation {
	output := ExecutionEnvelopeJEVImprovementReviewObservation{
		Status: "UNKNOWN", NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.GateDecision.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		return finalizeExecutionEnvelopeJEVImprovementReviewObservation(output)
	}
	if !input.GateDecision.NonExecuting {
		output.MissingStage = "execution-boundary"
		return finalizeExecutionEnvelopeJEVImprovementReviewObservation(output)
	}
	if strings.TrimSpace(input.GateDecision.DecisionDigest) == "" {
		output.MissingStage = input.GateDecision.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "review-gate-decision"
		}
		return finalizeExecutionEnvelopeJEVImprovementReviewObservation(output)
	}
	switch input.GateDecision.Status {
	case "review-eligible", "revision-required", "review-hold":
	default:
		output.MissingStage = input.GateDecision.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "review-disposition"
		}
		return finalizeExecutionEnvelopeJEVImprovementReviewObservation(output)
	}
	switch input.ReviewOutcome {
	case "confirmed", "refuted", "unknown":
	default:
		output.MissingStage = "review-outcome"
		return finalizeExecutionEnvelopeJEVImprovementReviewObservation(output)
	}
	if strings.TrimSpace(input.ReviewEvidenceDigest) == "" {
		output.MissingStage = "review-evidence"
		return finalizeExecutionEnvelopeJEVImprovementReviewObservation(output)
	}
	output.Status = "observed"
	output.ReviewOutcome = input.ReviewOutcome
	output.GateDecisionDigest = input.GateDecision.DecisionDigest
	output.ReviewEvidenceDigest = input.ReviewEvidenceDigest
	return finalizeExecutionEnvelopeJEVImprovementReviewObservation(output)
}

func finalizeExecutionEnvelopeJEVImprovementReviewObservation(output ExecutionEnvelopeJEVImprovementReviewObservation) ExecutionEnvelopeJEVImprovementReviewObservation {
	digest, err := digestExecutionEnvelopeJEVImprovementReviewObservation(output)
	if err != nil {
		output.Status = "UNKNOWN"
		output.MissingStage = "review-observation-digest"
		output.ObservationDigest = ""
		return output
	}
	output.ObservationDigest = digest
	return output
}
