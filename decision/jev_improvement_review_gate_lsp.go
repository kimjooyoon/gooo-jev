package decision

import "strings"

// ExecutionEnvelopeJEVImprovementReviewGateLSPInput projects a review-gate
// decision into deterministic editor-facing evidence.
type ExecutionEnvelopeJEVImprovementReviewGateLSPInput struct {
	Decision       ExecutionEnvelopeJEVImprovementReviewGateDecision
	NonAuthorizing bool
}

// ExecutionEnvelopeJEVImprovementReviewGateLSPProjection is a non-executing
// diagnostic projection of the review boundary.
type ExecutionEnvelopeJEVImprovementReviewGateLSPProjection struct {
	Status         string
	Code           string
	Severity       string
	Message        string
	DecisionDigest string
	EvidenceDigest string
	MissingStage   string
	NonExecuting   bool
	NonAuthorizing bool
}

func digestExecutionEnvelopeJEVImprovementReviewGateLSPProjection(projection ExecutionEnvelopeJEVImprovementReviewGateLSPProjection) (string, error) {
	return Digest(struct {
		Status         string
		Code           string
		Severity       string
		Message        string
		DecisionDigest string
		MissingStage   string
	}{
		Status:         projection.Status,
		Code:           projection.Code,
		Severity:       projection.Severity,
		Message:        projection.Message,
		DecisionDigest: projection.DecisionDigest,
		MissingStage:   projection.MissingStage,
	})
}

// ProjectExecutionEnvelopeJEVImprovementReviewGateLSP exposes review status to
// language tooling without changing the decision or authorizing a revision.
func ProjectExecutionEnvelopeJEVImprovementReviewGateLSP(input ExecutionEnvelopeJEVImprovementReviewGateLSPInput) ExecutionEnvelopeJEVImprovementReviewGateLSPProjection {
	output := ExecutionEnvelopeJEVImprovementReviewGateLSPProjection{
		Status: "UNKNOWN", Code: "JEV_REVIEW_GATE_UNKNOWN", Severity: "warning",
		NonExecuting: true, NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Decision.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		output.Message = "JEV review gate is UNKNOWN: missing authorization boundary"
		return finalizeExecutionEnvelopeJEVImprovementReviewGateLSP(output)
	}
	if !input.Decision.NonExecuting {
		output.MissingStage = "execution-boundary"
		output.Message = "JEV review gate is UNKNOWN: missing execution boundary"
		return finalizeExecutionEnvelopeJEVImprovementReviewGateLSP(output)
	}
	if strings.TrimSpace(input.Decision.DecisionDigest) == "" {
		output.MissingStage = input.Decision.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "review-gate-decision"
		}
		output.Message = "JEV review gate is UNKNOWN: missing " + output.MissingStage
		return finalizeExecutionEnvelopeJEVImprovementReviewGateLSP(output)
	}
	switch input.Decision.Status {
	case "review-eligible":
		output.Status = "review-eligible"
		output.Code = "JEV_REVIEW_ELIGIBLE"
		output.Severity = "info"
		output.Message = "JEV improvement evidence is eligible for human review"
	case "revision-required":
		output.Status = "revision-required"
		output.Code = "JEV_REVISION_REQUIRED"
		output.Severity = "warning"
		output.Message = "JEV feedback requires revision before review can proceed"
	case "review-hold":
		output.Status = "review-hold"
		output.Code = "JEV_REVIEW_HOLD"
		output.Severity = "hint"
		output.Message = "JEV feedback is held until missing evidence is resolved"
	default:
		output.MissingStage = input.Decision.MissingStage
		if strings.TrimSpace(output.MissingStage) == "" {
			output.MissingStage = "review-disposition"
		}
		output.Message = "JEV review gate is UNKNOWN: missing " + output.MissingStage
		return finalizeExecutionEnvelopeJEVImprovementReviewGateLSP(output)
	}
	output.DecisionDigest = input.Decision.DecisionDigest
	return finalizeExecutionEnvelopeJEVImprovementReviewGateLSP(output)
}

func finalizeExecutionEnvelopeJEVImprovementReviewGateLSP(output ExecutionEnvelopeJEVImprovementReviewGateLSPProjection) ExecutionEnvelopeJEVImprovementReviewGateLSPProjection {
	digest, err := digestExecutionEnvelopeJEVImprovementReviewGateLSPProjection(output)
	if err != nil {
		output.Status = "UNKNOWN"
		output.Code = "JEV_REVIEW_GATE_UNKNOWN"
		output.Severity = "warning"
		output.MissingStage = "lsp-evidence-digest"
		output.Message = "JEV review gate is UNKNOWN: missing lsp-evidence-digest"
		output.EvidenceDigest = ""
		return output
	}
	output.EvidenceDigest = digest
	return output
}
