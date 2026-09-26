package decision

import (
	"crypto/sha256"
	"encoding/hex"
)

// ExecutionEnvelopeReviewHandoffInput binds a counterexample disposition to
// a review handoff without treating the handoff as authorization.
type ExecutionEnvelopeReviewHandoffInput struct {
	Status          string `json:"status"`
	ReviewRequired  bool   `json:"review_required"`
	ComparisonDigest string `json:"comparison_digest"`
	NonAuthorizing  bool   `json:"non_authorizing"`
}

// ExecutionEnvelopeReviewHandoff records a durable handoff boundary.
type ExecutionEnvelopeReviewHandoff struct {
	Status         string `json:"status"`
	HandoffDigest  string `json:"handoff_digest"`
	ReviewRequired bool   `json:"review_required"`
	NonAuthorizing bool   `json:"non_authorizing"`
}

// CreateExecutionEnvelopeReviewHandoff makes only an explicit review-required
// disposition handoffable and keeps other states non-authorizing.
func CreateExecutionEnvelopeReviewHandoff(input ExecutionEnvelopeReviewHandoffInput) ExecutionEnvelopeReviewHandoff {
	output := ExecutionEnvelopeReviewHandoff{Status: "UNKNOWN", NonAuthorizing: true}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		return output
	}
	switch input.Status {
	case "review-required":
		if !input.ReviewRequired || input.ComparisonDigest == "" {
			return output
		}
		digest := sha256.Sum256([]byte(input.ComparisonDigest + "\x00review-required"))
		output.Status = "handoff-required"
		output.HandoffDigest = hex.EncodeToString(digest[:])
		output.ReviewRequired = true
	case "observation-only":
		output.Status = "observation-only"
		output.ReviewRequired = false
		output.HandoffDigest = input.ComparisonDigest
	case "hold":
		output.Status = "hold"
	default:
		output.Status = "UNKNOWN"
	}
	return output
}
