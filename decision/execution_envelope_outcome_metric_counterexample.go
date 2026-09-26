package decision

import (
	"crypto/sha256"
	"encoding/hex"
)

// ExecutionEnvelopeOutcomeMetricCounterexampleInput binds an outcome metric
// comparison to a bounded disposition.
type ExecutionEnvelopeOutcomeMetricCounterexampleInput struct {
	PreviousOutcomeDigest string `json:"previous_outcome_digest"`
	CurrentOutcomeDigest  string `json:"current_outcome_digest"`
	Delta                 string `json:"delta"`
	NonAuthorizing        bool   `json:"non_authorizing"`
}

// ExecutionEnvelopeOutcomeMetricCounterexample records whether an outcome
// comparison requires review without asserting improvement.
type ExecutionEnvelopeOutcomeMetricCounterexample struct {
	Status          string `json:"status"`
	ReviewRequired  bool   `json:"review_required"`
	ComparisonDigest string `json:"comparison_digest"`
	NonAuthorizing  bool   `json:"non_authorizing"`
}

// ClassifyExecutionEnvelopeOutcomeMetricCounterexample is fail-closed for
// incomplete comparisons and authorization claims.
func ClassifyExecutionEnvelopeOutcomeMetricCounterexample(input ExecutionEnvelopeOutcomeMetricCounterexampleInput) ExecutionEnvelopeOutcomeMetricCounterexample {
	output := ExecutionEnvelopeOutcomeMetricCounterexample{Status: "UNKNOWN", NonAuthorizing: true}
	if !input.NonAuthorizing {
		output.NonAuthorizing = false
		return output
	}
	if input.PreviousOutcomeDigest == "" || input.CurrentOutcomeDigest == "" {
		return output
	}
	digest := sha256.Sum256([]byte(input.PreviousOutcomeDigest + "\x00" + input.CurrentOutcomeDigest + "\x00" + input.Delta))
	output.ComparisonDigest = hex.EncodeToString(digest[:])
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
