package decision

// ExecutionEnvelopeCounterexampleDispositionInput turns an observed reverse
// boundary into a non-authorizing review disposition.
type ExecutionEnvelopeCounterexampleDispositionInput struct {
	Status         string `json:"status"`
	FirstMismatch  string `json:"first_mismatch"`
	EvidenceDigest string `json:"evidence_digest"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

// ExecutionEnvelopeCounterexampleDisposition records whether a human or
// policy review is required without claiming that a change is an improvement.
type ExecutionEnvelopeCounterexampleDisposition struct {
	Status         string `json:"status"`
	ReviewRequired bool   `json:"review_required"`
	MismatchStage  string `json:"mismatch_stage"`
	EvidenceDigest string `json:"evidence_digest"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

// ClassifyExecutionEnvelopeCounterexample is fail-closed for unknown or
// incomplete reverse observations.
func ClassifyExecutionEnvelopeCounterexample(input ExecutionEnvelopeCounterexampleDispositionInput) ExecutionEnvelopeCounterexampleDisposition {
	output := ExecutionEnvelopeCounterexampleDisposition{Status: "UNKNOWN", NonAuthorizing: true}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		output.Status = "UNKNOWN"
		output.MismatchStage = "authorization-boundary"
		return output
	}
	switch input.Status {
	case "counterexample":
		if input.FirstMismatch == "" || input.EvidenceDigest == "" {
			output.MismatchStage = "counterexample-evidence"
			return output
		}
		output.Status = "review-required"
		output.ReviewRequired = true
		output.MismatchStage = input.FirstMismatch
		output.EvidenceDigest = input.EvidenceDigest
	case "reproduced":
		if input.EvidenceDigest == "" {
			output.MismatchStage = "reproduced-evidence"
			return output
		}
		output.Status = "observation-only"
		output.EvidenceDigest = input.EvidenceDigest
	case "UNKNOWN":
		output.Status = "hold"
		output.MismatchStage = "reverse-observation"
	default:
		output.MismatchStage = "reverse-observation"
	}
	return output
}
