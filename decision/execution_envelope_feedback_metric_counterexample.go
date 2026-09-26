package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// ExecutionEnvelopeFeedbackMetricCounterexampleInput binds a metric comparison
// to a review disposition without treating a decline as an authorization.
type ExecutionEnvelopeFeedbackMetricCounterexampleInput struct {
	PreviousFeedbackDigest string `json:"previous_feedback_digest"`
	CurrentFeedbackDigest  string `json:"current_feedback_digest"`
	Delta                  string `json:"delta"`
	NonAuthorizing         bool   `json:"non_authorizing"`
}

// ExecutionEnvelopeFeedbackMetricCounterexample records the next bounded
// action for a metric comparison.
type ExecutionEnvelopeFeedbackMetricCounterexample struct {
	Status          string `json:"status"`
	ReviewRequired  bool   `json:"review_required"`
	ComparisonDigest string `json:"comparison_digest"`
	NonAuthorizing  bool   `json:"non_authorizing"`
}

// ClassifyExecutionEnvelopeFeedbackMetricCounterexample keeps declines
// reviewable and never infers improvement from a comparison alone.
func ClassifyExecutionEnvelopeFeedbackMetricCounterexample(input ExecutionEnvelopeFeedbackMetricCounterexampleInput) ExecutionEnvelopeFeedbackMetricCounterexample {
	output := ExecutionEnvelopeFeedbackMetricCounterexample{Status: "UNKNOWN", NonAuthorizing: true}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		return output
	}
	if input.PreviousFeedbackDigest == "" || input.CurrentFeedbackDigest == "" {
		return output
	}
	comparisonDigest := sha256.Sum256([]byte(strings.Join([]string{input.PreviousFeedbackDigest, input.CurrentFeedbackDigest, input.Delta}, "\x00")))
	output.ComparisonDigest = hex.EncodeToString(comparisonDigest[:])
	switch input.Delta {
	case "declined":
		output.Status = "review-required"
		output.ReviewRequired = true
	case "increased", "unchanged":
		output.Status = "observation-only"
	case "inconclusive":
		output.Status = "hold"
	case "UNKNOWN":
		output.Status = "UNKNOWN"
	default:
		output.Status = "UNKNOWN"
	}
	return output
}
